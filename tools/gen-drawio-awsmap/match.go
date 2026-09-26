package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Matching: an old name is matched with the entries of the current
// palette by normalized names. The names of an old shape are its name in
// the stencil file (or its JavaScript name), its local name, and the
// titles and tags of the palette entries that use it; those of a current
// entry are its icon's name, its title and its tags. Names are normalized
// by dropping punctuation, case, a variant number ("topic_2") and the
// words "Amazon" and "AWS", and by replacing the old names of renamed
// services with their current ones (aliases). The first name of an old
// shape that equals names of current entries picks among them: an entry
// of the corresponding category, then the strongest match (icon name,
// then title, then tags), then a service (resource icon) over a
// resource, then palette order. When no name matches, the names are tried
// again without a leading service name ("EC2 Instance" → "Instance").
// Groups are never matched automatically.

// target is an entry of the current palette that old names map to.
type target struct {
	entry *paletteEntry
	// kind is "res" (a resourceIcon entry), "plain" (a stencil entry),
	// "group", "rect" (a group drawn as a rectangle), and for overrides
	// "groupicon" (a group's icon alone) or "builtin" (a draw.io style)
	kind     string
	id       string // the icon's local name (resIcon, grIcon or shape), the title for rect, the style for builtin
	category string // the palette's category, e.g. "Compute"
	order    int
	// idKeys, titleKey and tagsKey are the normalized names
	idKeys            []string
	titleKey, tagsKey string
	// iconW, iconH are the size of the stencil of a plain or groupicon target
	iconW, iconH float64
}

// label is how a target is written in the overrides file and reports.
func (t *target) label() string {
	if t.kind == "plain" {
		return t.id
	}
	return t.kind + ":" + t.id
}

// strength returns how well key names the target: 3 its icon, 2 its
// title, 1 its tags, 0 not at all.
func (t *target) strength(key string) int {
	switch {
	case contains(t.idKeys, key):
		return 3
	case t.titleKey == key:
		return 2
	case t.tagsKey == key:
		return 1
	}
	return 0
}

// targets returns the entries of the current palette as targets.
func (s *sources) targets() []*target {
	var out []*target
	for i := range s.aws4 {
		e := &s.aws4[i]
		st := parseStyle(e.style)
		t := &target{entry: e, category: strings.TrimPrefix(e.palette, "aws4"), order: i}
		shape := strings.TrimPrefix(st["shape"], "mxgraph.aws4.")
		switch {
		case shape == "resourceIcon":
			t.kind, t.id = "res", strings.TrimPrefix(st["resIcon"], "mxgraph.aws4.")
		case shape == "group" || shape == "groupCenter":
			t.kind, t.id = "group", strings.TrimPrefix(st["grIcon"], "mxgraph.aws4.")
		case st["shape"] == "":
			t.kind, t.id = "rect", strings.ReplaceAll(strings.ToLower(e.title), " ", "_")
		default:
			t.kind, t.id = "plain", shape
			info := s.aws4Stencils[st["shape"]]
			t.iconW, t.iconH = info.w, info.h
		}
		t.idKeys = uniq(normalize(t.id))
		if strings.Contains(t.id, "_") {
			// a variant number: cloudwatch_2, instance_with_cloudwatch2
			t.idKeys = uniq(normalize(t.id), stripVariant(normalize(t.id)))
		}
		t.titleKey, t.tagsKey = normalize(e.title), normalize(e.tags)
		out = append(out, t)
	}
	return out
}

var variantRe = regexp.MustCompile(`^(.*[a-z]{3}) ?\d$`)

// stripVariant drops a variant number from a normalized name: "topic 2"
// and "ebs2" but not "ec2", "s3" or "route 53".
func stripVariant(k string) string {
	if m := variantRe.FindStringSubmatch(k); m != nil {
		return m[1]
	}
	return k
}

// categoryOf maps the categories of the older palettes and stencil files
// to those of the current palette.
var categoryOf = map[string]string{
	"analytics": "Analytics", "application services": "Application Integration",
	"app services": "Application Integration", "artificial intelligence": "Artificial Intelligence",
	"business productivity": "Business Applications", "compute": "Compute",
	"compute and networking": "Compute", "contact center": "Contact Center", "database": "Database",
	"desktop and app streaming": "End User Computing", "developer tools": "Developer Tools",
	"game development": "Games", "general": "General Resources", "non service specific": "General Resources",
	"groups": "Groups", "misc": "Groups", "internet of things": "Internet of Things",
	"management tools": "Management Governance", "deployment management": "Management Governance",
	"deployment and management": "Management Governance", "administration and security": "Management Governance",
	"messaging": "Application Integration", "migration": "Migration Modernization",
	"mobile services": "Front End Web Mobile", "networking and content delivery": "Network Content Delivery",
	"networking": "Network Content Delivery", "content delivery": "Network Content Delivery",
	"sdks": "Developer Tools", "security identity and compliance": "Security Identity Compliance",
	"security and identity": "Security Identity Compliance", "storage": "Storage",
	"storage and content delivery": "Storage", "enterprise applications": "Business Applications",
}

// oldEntries returns the palette entries that draw an old name by itself
// (not the cells of composite entries, whose titles name the whole).
func (s *sources) oldEntries(n oldName) []*paletteEntry {
	var out []*paletteEntry
	for i := range s.old {
		e := &s.old[i]
		if !e.group && parseStyle(e.style)["shape"] == n.name {
			out = append(out, e)
		}
	}
	return out
}

// oldCategory returns the current category of an old name: that of its
// first palette entry, or of its stencil file.
func (s *sources) oldCategory(n oldName) string {
	if es := s.oldEntries(n); len(es) > 0 {
		p := es[0].palette
		for _, g := range []string{"aws2", "aws3d", "aws3"} {
			p = strings.TrimPrefix(p, g)
		}
		return categoryOf[strings.ToLower(strings.ReplaceAll(p, "-", " "))]
	}
	parts := strings.Split(n.name, ".")
	if len(parts) == 4 {
		return categoryOf[strings.NewReplacer("_", " ", "-", " ").Replace(parts[2])]
	}
	return ""
}

// oldKeys returns the normalized names of an old shape in the order they
// are tried.
func (s *sources) oldKeys(n oldName) []string {
	local := n.name[strings.LastIndex(n.name, ".")+1:]
	display := n.display
	if n.gen == "aws3d" {
		local, display = splitCamel(local), splitCamel(display)
	}
	keys := []string{stripVariant(normalize(display)), stripVariant(normalize(local))}
	for _, e := range s.oldEntries(n) {
		keys = append(keys, normalize(e.title), normalize(e.tags))
	}
	return uniq(keys...)
}

// servicePrefixes are service names that old resource names start with
// ("EC2 Instance", "S3 Bucket"), tried away when nothing else matches.
var servicePrefixes = []string{"ec2", "rds", "s3", "iam", "vpc", "sns", "sqs", "ses", "emr", "ebs",
	"cloudfront", "route 53", "dynamodb", "opsworks", "cloudformation", "cloudwatch",
	"elastic beanstalk", "kinesis", "storage gateway", "redshift", "elasticache", "s3 glacier"}

// match returns the target an old name matches automatically and the
// name it matched by, or nil.
func (s *sources) match(n oldName, targets []*target) (*target, string) {
	cat := s.oldCategory(n)
	keys := s.oldKeys(n)
	var stripped []string
	for _, k := range keys {
		for _, p := range servicePrefixes {
			if rest, ok := strings.CutPrefix(k, p+" "); ok {
				stripped = append(stripped, rest)
			}
		}
	}
	for _, k := range append(keys, uniq(stripped...)...) {
		var best *target
		bestStrength := 0
		for _, t := range targets {
			if t.kind == "group" || t.kind == "rect" {
				continue
			}
			st := t.strength(k)
			if st == 0 {
				continue
			}
			if best == nil || better(t, st, best, bestStrength, cat) {
				best, bestStrength = t, st
			}
		}
		if best != nil {
			return best, k
		}
	}
	return nil, ""
}

// better reports whether t (matching with strength ts) is preferred over
// u (strength us) for an old name of category cat.
func better(t *target, ts int, u *target, us int, cat string) bool {
	if (t.category == cat) != (u.category == cat) {
		return t.category == cat
	}
	if ts != us {
		return ts > us
	}
	if (t.kind == "res") != (u.kind == "res") {
		return t.kind == "res"
	}
	return t.order < u.order
}

// aliases replace old or alternative names of services with the names
// the current palette uses (whole normalized words).
var aliases = [][2]string{
	{"simple storage service", "s3"},
	{"simple notification service", "sns"},
	{"simple queue service", "sqs"},
	{"simple email service", "ses"},
	{"simple workflow service", "swf"},
	{"elastic compute cloud", "ec2"},
	{"elastic map reduce", "emr"},
	{"elastic mapreduce", "emr"},
	{"elastic load balancer", "elastic load balancing"},
	{"elb", "elastic load balancing"},
	{"elasticsearch service", "opensearch service"},
	{"elastic search service", "opensearch service"},
	{"elasticsearch", "opensearch service"},
	{"elastic search", "opensearch service"},
	{"relational database service", "rds"},
	{"rds db instance", "rds instance"},
	{"identity and access management", "iam"},
	{"identity access management", "iam"},
	{"key management service", "kms"},
	{"elastic block store", "ebs"},
	{"elastic file system", "efs"},
	{"kinesis firehose", "data firehose"},
	{"kinesis streams", "kinesis data streams"},
	{"kinesis analytics", "kinesis data analytics"},
	{"ec2 container service", "elastic container service"},
	{"ec2 container registry", "elastic container registry"},
	{"ecs", "elastic container service"},
	{"ecr", "elastic container registry"},
	{"elasticcache", "elasticache"},
	{"elastic cache", "elasticache"},
	{"dynamo db", "dynamodb"},
	{"simple db", "simpledb"},
	{"code commit", "codecommit"},
	{"code deploy", "codedeploy"},
	{"code pipeline", "codepipeline"},
	{"code build", "codebuild"},
	{"code star", "codestar"},
	{"cloud hsm", "cloudhsm"},
	{"cloud search", "cloudsearch"},
	{"cloud watch", "cloudwatch"},
	{"cloud trail", "cloudtrail"},
	{"cloud formation", "cloudformation"},
	{"cloud front", "cloudfront"},
	{"clouddirectory", "cloud directory"},
	{"route53", "route 53"},
	{"acm", "certificate manager"},
	{"webapp firewall", "waf"},
	{"web application firewall", "waf"},
	{"mobile analytics", "pinpoint"},
	{"mobile hub", "amplify"},
	{"x ray", "xray"},
	{"game lift", "gamelift"},
	{"ops works", "opsworks"},
	{"work docs", "workdocs"},
	{"work mail", "workmail"},
	{"work spaces", "workspaces"},
	{"quick sight", "quicksight"},
	{"ec2 systems manager", "systems manager"},
	{"import export", "snowball import export"},
	{"glacier", "s3 glacier"},
	{"cli", "command line interface"},
	{"sts", "security token service"},
	{"flowlogs", "flow logs"},
	{"vpcnat", "vpc nat"},
	{"hostedzone", "hosted zone"},
	{"routetable", "route table"},
	{"beanstalk", "elastic beanstalk"},
}

var aliasRes = func() []*regexp.Regexp {
	var out []*regexp.Regexp
	for _, a := range aliases {
		out = append(out, regexp.MustCompile(`(^| )`+regexp.QuoteMeta(a[0])+`( |$)`))
	}
	return out
}()

// normalize returns the comparable form of a name.
func normalize(s string) string {
	var words []string
	for _, w := range strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
	}) {
		if w != "amazon" && w != "aws" {
			words = append(words, w)
		}
	}
	out := strings.Join(words, " ")
	for i, re := range aliasRes {
		if re.MatchString(out) && !strings.Contains(out, aliases[i][1]) {
			out = re.ReplaceAllString(out, "${1}"+aliases[i][1]+"${2}")
		}
	}
	return out
}

// splitCamel splits the words of a camel case identifier (dataCenter →
// data Center).
func splitCamel(s string) string {
	var b strings.Builder
	prev := rune(0)
	for _, r := range s {
		if unicode.IsUpper(r) && unicode.IsLower(prev) {
			b.WriteByte(' ')
		}
		b.WriteRune(r)
		prev = r
	}
	return b.String()
}

func uniq(ss ...string) []string {
	var out []string
	for _, s := range ss {
		if s != "" && !contains(out, s) {
			out = append(out, s)
		}
	}
	return out
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

// findTarget returns the target an override names: res:<resIcon>,
// <stencil>, group:<grIcon> or rect:<title>, optionally followed by
// @<Category> (spaces as _) to pick the entry of that palette.
func findTarget(targets []*target, spec string) (*target, error) {
	spec, cat, _ := strings.Cut(spec, "@")
	cat = strings.ReplaceAll(cat, "_", " ")
	kind, id := "plain", spec
	if k, v, ok := strings.Cut(spec, ":"); ok {
		kind, id = k, v
	}
	for _, t := range targets {
		if t.kind == kind && t.id == id && (cat == "" || t.category == cat) {
			return t, nil
		}
	}
	return nil, fmt.Errorf("no palette entry %q", spec)
}
