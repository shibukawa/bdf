package main

import "testing"

func TestNormalize(t *testing.T) {
	for in, want := range map[string]string{
		"AWS Elastic Load Balancer":     "elastic load balancing",
		"Amazon Kinesis Firehose":       "data firehose",
		"ElastiCache":                   "elasticache",
		"EC2 Container Registry":        "elastic container registry",
		"Amazon Simple Storage Service": "s3",
		"Simple Storage Service (S3)":   "simple storage service s3", // the name is there already
		"Route53":                       "route 53",
		"Import/Export":                 "snowball import export",
	} {
		if got := normalize(in); got != want {
			t.Errorf("normalize(%q) = %q, want %q", in, got, want)
		}
	}
	for in, want := range map[string]string{"topic 2": "topic", "ebs2": "ebs", "ec2": "ec2", "s3": "s3", "route 53": "route 53",
		"certificate manager 2": "certificate manager", "emr engine mapr m3": "emr engine mapr m3"} {
		if got := stripVariant(in); got != want {
			t.Errorf("stripVariant(%q) = %q, want %q", in, got, want)
		}
	}
	if got := splitCamel("elasticLoadBalancing"); got != "elastic Load Balancing" {
		t.Errorf("splitCamel: %q", got)
	}
}

// TestMatch checks the ranking of candidates: the category of the old
// shape first, then the strength of the match, then services over
// resources.
func TestMatch(t *testing.T) {
	s := &sources{
		aws4: []paletteEntry{
			{palette: "aws4Application Integration", title: "Topic", style: "shape=mxgraph.aws4.topic;", w: 78, h: 67},
			{palette: "aws4Internet of Things", title: "Topic", style: "shape=mxgraph.aws4.topic_2;", w: 53, h: 78},
			{palette: "aws4General Resources", title: "Internet", tags: "internet gateway", style: "shape=mxgraph.aws4.internet_alt2;", w: 78, h: 78},
			{palette: "aws4Network Content Delivery", title: "Internet Gateway", style: "shape=mxgraph.aws4.internet_gateway;", w: 78, h: 78},
			{palette: "aws4Compute", title: "Auto Scaling", style: "shape=mxgraph.aws4.auto_scaling2;", w: 48, h: 48},
			{palette: "aws4Compute", title: "Auto Scaling", style: "shape=mxgraph.aws4.resourceIcon;resIcon=mxgraph.aws4.auto_scaling3;", w: 78, h: 78},
			{palette: "aws4Groups", title: "Region", style: "shape=mxgraph.aws4.group;grIcon=mxgraph.aws4.group_region;", w: 130, h: 130},
		},
		old: []paletteEntry{
			{palette: "aws3Internet of Things", title: "Topic", style: "shape=mxgraph.aws3.topic;"},
			{palette: "aws3Messaging", title: "Topic", style: "shape=mxgraph.aws3.topic_2;"},
			{palette: "aws3Compute", title: "Internet Gateway", style: "shape=mxgraph.aws3.internet_gateway;"},
			{palette: "aws3Compute", title: "Auto Scaling", style: "shape=mxgraph.aws3.auto_scaling;"},
		},
		aws4Stencils: map[string]stencilInfo{},
	}
	targets := s.targets()
	for name, want := range map[string]string{
		"mxgraph.aws3.topic":            "topic_2", // the IoT topic, by category
		"mxgraph.aws3.topic_2":          "topic",   // the SNS topic: the variant number goes
		"mxgraph.aws3.internet_gateway": "internet_gateway",
		"mxgraph.aws3.auto_scaling":     "res:auto_scaling3",
		"mxgraph.aws3.region":           "", // groups are not matched
	} {
		n := oldName{gen: "aws3", name: name, display: name[len("mxgraph.aws3."):]}
		got, _ := s.match(n, targets)
		label := ""
		if got != nil {
			label = got.label()
		}
		if label != want {
			t.Errorf("%s: %q, want %q", name, label, want)
		}
	}
}

func TestParseOverrides(t *testing.T) {
	ovs, err := parseOverrides("# comment\n\nmxgraph.aws3.a  res:x  # why\nmxgraph.aws3.b -\n")
	if err != nil || len(ovs) != 2 || ovs[0].spec != "res:x" || ovs[0].comment != "why" || ovs[1].spec != "-" || ovs[1].line != 4 {
		t.Fatalf("%+v %v", ovs, err)
	}
	for _, bad := range []string{"mxgraph.aws3.a\n", "mxgraph.aws3.a x y\n", "mxgraph.aws3.a x\nmxgraph.aws3.a y\n"} {
		if _, err := parseOverrides(bad); err == nil {
			t.Errorf("%q accepted", bad)
		}
	}
}
