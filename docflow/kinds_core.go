package docflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/shibukawa/yuruidb"
)

// The core's kinds of node. Their settings are checked here; what a job
// that reaches one does is in route.go and the file of the kind.
func init() {
	RegisterKind(Kind{Name: KindUpload,
		Titles: Text{"en": "Uploader", "ja": "アップローダー"}, Help: Text{"en": "a REST entry for applications", "ja": "アプリケーション向けの REST の入口"},
		Icon: []string{"M14 4h4a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2h-4", "M3 12h10", "M9 8l4 4-4 4"}, Role: RoleEntry, PassesOn: true})
	RegisterKind(Kind{Name: KindQueue,
		Titles: Text{"en": "Queue", "ja": "キュー"}, Help: Text{"en": "forwards to every node it points to", "ja": "指すすべてのノードへ送る"},
		Icon: []string{"M12 3l9 5-9 5-9-5z", "M3 13l9 5 9-5", "M3 17l9 5 9-5"}, Role: RoleRoute, PassesOn: true, Check: checkForward})
	RegisterKind(Kind{Name: KindClass,
		Titles: Text{"en": "Class", "ja": "クラス"}, Help: Text{"en": "one member in turn, or the first", "ja": "メンバーの 1 つへ順番に、または先頭へ"},
		Icon: []string{"M17 3l4 4-4 4", "M3 11V9a2 2 0 0 1 2-2h16", "M7 21l-4-4 4-4", "M21 13v2a2 2 0 0 1-2 2H3"}, Role: RoleRoute, PassesOn: true, Check: checkClass, route: (*Manager).routeClass})
	RegisterKind(Kind{Name: KindFilter,
		Titles: Text{"en": "Filter", "ja": "フィルタ"}, Help: Text{"en": "sends on by rules on what a job is, and its metadata", "ja": "ジョブの性質とメタデータのルールで振り分ける"},
		Icon: []string{"M4 4h16l-6 8v6l-4 2v-8z"}, Picks: true, Role: RoleRoute, PassesOn: true, Check: checkFilter, route: (*Manager).routeFilter})
	RegisterKind(Kind{Name: KindAssert,
		Titles: Text{"en": "Assert", "ja": "アサート"}, Help: Text{"en": "checks the metadata; right after an entry, refuses the sender what fails", "ja": "メタデータを検査する。入口の直後なら、通らないものを送り手に拒否する"},
		Icon: []string{"M12 3l8 3v6c0 5-3.5 8-8 9-4.5-1-8-4-8-9V6z", "M9 12l2 2 4-4"}, Picks: true, Role: RoleRoute, PassesOn: true, Check: checkAssert, route: (*Manager).routeAssert})
	RegisterKind(Kind{Name: KindBranch,
		Titles: Text{"en": "Branch", "ja": "ブランチ"}, Help: Text{"en": "a switch: the first case whose checks hold, else default", "ja": "switch 文。チェックを満たす最初のケース、なければ default へ"},
		Icon: []string{"M3 12h6", "M9 12c4 0 4-6 8-6h4", "M9 12c4 0 4 6 8 6h4", "M18 3l3 3-3 3", "M18 15l3 3-3 3"}, Picks: true, Role: RoleRoute, PassesOn: true, Check: checkBranch, route: (*Manager).routeBranch})
	RegisterKind(Kind{Name: KindInbox,
		Titles: Text{"en": "Inbox", "ja": "インボックス"}, Help: Text{"en": "held until a person or an app releases it", "ja": "人かアプリが送り出すまで留める"},
		Icon: []string{"M5.5 5h13l3.5 7v7H2v-7z", "M2 12h5l2 3h6l2-3h5"}, Role: RoleRoute, PassesOn: true, Stocks: true, Check: checkInbox,
		route: func(m *Manager, ctx context.Context, r *routing, name, from string, depth int) {
			m.routeInbox(ctx, r, name, from, depth)
		}})
	RegisterKind(Kind{Name: KindTemplate,
		Titles: Text{"en": "Template", "ja": "テンプレート"}, Help: Text{"en": "fills a template with data, into a .bdf document", "ja": "データでテンプレートを埋めて .bdf 文書にする"},
		Icon: []string{"M14 3H7a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8z", "M14 3v5h5", "M9 13h6", "M9 17h4"}, Role: RoleStep, Takes: []string{PayloadData}, Makes: PayloadBDF, Check: checkTemplate, route: routeStep(KindTemplate)})
	RegisterKind(Kind{Name: KindConvert,
		Titles: Text{"en": "Convert", "ja": "変換"}, Help: Text{"en": "makes a .bdf document of a file received as it is", "ja": "受け取ったままのファイルから .bdf 文書を作る"},
		Icon: []string{"M21 12a9 9 0 0 1-15 6.7L3 16", "M3 12a9 9 0 0 1 15-6.7L21 8", "M3 21v-5h5", "M21 3v5h-5"}, Role: RoleStep, Makes: PayloadBDF, Check: checkConvert, route: routeStep(KindConvert)})
	RegisterKind(Kind{Name: KindSign,
		Titles: Text{"en": "Sign", "ja": "署名"}, Help: Text{"en": "signs a .bdf document as a witness of who sent it", "ja": "誰が送ったかの証人として .bdf 文書に署名する"},
		Icon: []string{"M3 21h18", "M14.5 4.5l5 5L9 20H4v-5z", "M12.5 6.5l5 5"}, Role: RoleStep, Takes: []string{PayloadBDF}, Makes: PayloadBDF, Check: checkSign, route: routeStep(KindSign)})
	RegisterKind(Kind{Name: KindArchive,
		Titles: Text{"en": "Archive", "ja": "アーカイブ"}, Help: Text{"en": "written once, never changed", "ja": "1 度だけ書き、変更しない"},
		Icon: []string{"M3 4h18v4H3z", "M5 8v12h14V8", "M10 12h4"}, Role: RoleExit, PassesOn: true, Stocks: true, Check: checkArchive, route: routeStep(KindArchive)})
	RegisterKind(Kind{Name: KindDisplay,
		Titles: Text{"en": "Display", "ja": "ディスプレイ"}, Help: Text{"en": "shown on the devices enrolled there", "ja": "そこに登録したデバイスに表示する"},
		Icon: []string{"M3 4h18v12H3z", "M8 21h8", "M12 16v5"}, Role: RoleExit, Takes: []string{PayloadBDF}, Makes: PayloadBDF, Fixed: true, Check: checkDisplay,
		route: func(m *Manager, ctx context.Context, r *routing, name, from string, _ int) {
			m.routeDisplay(ctx, r, name, from)
		}})
	RegisterKind(Kind{Name: KindSink,
		Titles: Text{"en": "Sink", "ja": "シンク"}, Help: Text{"en": "keeps the latest documents, to look at", "ja": "最新の文書を見られるように保持する"},
		Icon: []string{"M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7z", "M12 9a3 3 0 1 0 0 6 3 3 0 0 0 0-6z"}, Role: RoleExit, Fixed: true, Check: checkSink,
		route: func(m *Manager, ctx context.Context, r *routing, name, from string, _ int) {
			m.routeSink(ctx, r, name, from)
		}})
	RegisterKind(Kind{Name: KindDownload,
		Titles: Text{"en": "Downloader", "ja": "ダウンローダー"}, Help: Text{"en": "fetched by applications and people from an inbox or an archive", "ja": "インボックスかアーカイブのものを、アプリケーションや人が取りに来る"},
		Icon: []string{"M4 14v5a1 1 0 0 0 1 1h14a1 1 0 0 0 1-1v-5", "M12 3v12", "M7 8l5-5 5 5"}, Role: RoleExit, Takes: []string{PayloadStocked}, Fixed: true, Check: func(ctx context.Context, m *Manager, in Node, n *Node) error { return m.checkDownload(ctx, in, n) },
		route: func(m *Manager, ctx context.Context, r *routing, name, from string, _ int) {
			m.routeDownload(ctx, r, name, from)
		}})
}

// routeStep makes a job that reaches a node of the kind a delivery: the
// step of a template or a convert node, an archive's write.
func routeStep(kind string) func(m *Manager, ctx context.Context, r *routing, name, from string, _ int) {
	return func(m *Manager, ctx context.Context, r *routing, name, from string, _ int) {
		m.routeDelivery(ctx, r, Leg{Node: name, Kind: kind, From: from}, "", "", kind)
	}
}

// checkForward takes where a node forwards to. A node is put before the
// nodes after it, so none is fine: a job that reaches it meanwhile gets
// an error there (route).
func checkForward(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.To = cleanNames(in.To)
	return nil
}

// checkEntry checks the settings every entry has: where it forwards, what
// it accepts and sends on, the formats it takes, the metadata it adds,
// its own port (of a kind that serves one: a CUPS server's IPP; the
// others are reached at the API), and the apps that may send to it with
// their own credentials (an uploader's REST, a CUPS server's IPP Basic).
func checkEntry(ctx context.Context, m *Manager, k *Kind, in Node, n *Node) error {
	n.To = cleanNames(in.To)
	n.Accept = or(in.Accept, AcceptFile)
	if n.Accept != AcceptFile && n.Accept != AcceptData {
		return fmt.Errorf("%w: accept %q: file or data", ErrInvalid, in.Accept)
	}
	n.Output = or(in.Output, OutputBDF)
	if n.Output != OutputBDF && n.Output != OutputFile {
		return fmt.Errorf("%w: output %q: bdf (the document made of the file) or file (as received)", ErrInvalid, in.Output)
	}
	if n.Accept == AcceptData {
		n.Output = "" // data goes on as it is
	}
	n.Formats = cleanFormats(in.Formats)
	if err := checkFormats(n.Formats); err != nil {
		return err
	}
	var err error
	if n.Meta, err = claimedMeta(in.Meta); err != nil {
		return err
	}
	if ServesPort(k.Name) {
		n.Port = in.Port
	}
	return m.checkApps(ctx, in, n)
}

// checkApps takes the apps a node lists: those that send to an entry,
// read a downloader or work an inbox with their own credentials.
func (m *Manager) checkApps(ctx context.Context, in Node, n *Node) error {
	n.Apps = cleanNames(in.Apps)
	for _, id := range n.Apps {
		if _, err := m.App(ctx, id); err != nil {
			return fmt.Errorf("%w: no app %q", ErrInvalid, id)
		}
	}
	return nil
}

func checkTemplate(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.To = cleanNames(in.To)
	// its templates are uploaded once it is there: none named yet is
	// fine, and a job that reaches it meanwhile fails there (fill)
	n.Template, n.TemplateFrom = strings.TrimSpace(in.Template), strings.TrimSpace(in.TemplateFrom)
	if n.Template != "" && n.TemplateFrom != "" {
		return fmt.Errorf("%w: a template node names the template it fills (template), or the metadata key that names one (template_from), not both", ErrInvalid)
	}
	if n.Template != "" && !templateIDRE.MatchString(n.Template) {
		return fmt.Errorf("%w: template id %q: letters, digits, _ and -", ErrInvalid, n.Template)
	}
	n.ListName, n.Strict = strings.TrimSpace(in.ListName), in.Strict
	n.Original = cleanNames(in.Original)
	return nil
}

func checkConvert(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.To = cleanNames(in.To)
	n.Formats = cleanFormats(in.Formats)
	n.Original = cleanNames(in.Original)
	return checkFormats(n.Formats)
}

func checkSign(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.To = cleanNames(in.To)
	n.SignRole = strings.TrimSpace(in.SignRole)
	if len(n.SignRole) > 64 {
		return fmt.Errorf("%w: sign_role is 64 bytes at most", ErrInvalid)
	}
	return nil
}

func checkSink(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.Latest = in.Latest
	if n.Latest < 0 || n.Latest > 1000 {
		return fmt.Errorf("%w: latest %d: a sink keeps 1 to 1000 documents", ErrInvalid, in.Latest)
	}
	return nil
}

func checkClass(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.Members = cleanNames(in.Members)
	n.Policy = or(in.Policy, PolicyRoundRobin)
	if n.Policy != PolicyRoundRobin && n.Policy != PolicyFailover {
		return fmt.Errorf("%w: policy %q: round-robin or failover", ErrInvalid, in.Policy)
	}
	return nil
}

func checkFilter(ctx context.Context, m *Manager, in Node, n *Node) error {
	rules, err := checkRules(in.Rules)
	if err != nil {
		return err
	}
	n.Rules, n.Otherwise = rules, cleanNames(in.Otherwise)
	if len(n.Rules) == 0 && len(n.Otherwise) == 0 {
		return fmt.Errorf("%w: a filter has rules, or nodes it sends to otherwise", ErrInvalid)
	}
	return nil
}

func checkInbox(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.To, n.Keep = cleanNames(in.To), in.Keep
	return m.checkApps(ctx, in, n)
}

func checkArchive(ctx context.Context, m *Manager, in Node, n *Node) error {
	// what it wrote goes on to To, once written (deliver)
	n.To = cleanNames(in.To)
	n.Store, n.Collection = strings.TrimSpace(in.Store), strings.TrimSpace(in.Collection)
	if n.Store == "" {
		return fmt.Errorf("%w: an archive has a store (bbolt, s3, r2 or gcs)", ErrInvalid)
	}
	if n.Collection != "" {
		if err := yuruidb.ValidateCollection(n.Collection); err != nil {
			return fmt.Errorf("%w: collection %q: letters, digits, _ and -", ErrInvalid, n.Collection)
		}
	}
	// the store is opened now, so that a wrong URL is refused here
	if _, err := m.st.archive(ctx, n.Store, n.collection()); err != nil {
		return fmt.Errorf("%w: store: %v", ErrInvalid, err)
	}
	return nil
}

// checkDisplay takes a fallback that is an exit sending .bdf out (a
// printer).
func checkDisplay(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.Keep, n.OfflineAfter = in.Keep, in.OfflineAfter
	if in.Fallback != nil && strings.TrimSpace(in.Fallback.Printer) != "" {
		n.Fallback = &Fallback{Printer: strings.TrimSpace(in.Fallback.Printer), After: in.Fallback.After}
		k, err := m.kindOf(ctx, n.Fallback.Printer)
		if x := kindNamed(k); err != nil || x == nil || x.Send == nil || !x.takes().has(PayloadBDF) {
			return fmt.Errorf("%w: fallback printer %q is not a printer", ErrInvalid, n.Fallback.Printer)
		}
	}
	return nil
}
