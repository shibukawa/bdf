# Topology

The topology is the network the documents flow through: the nodes, of every kind, and the edges between them. It is one page, Topology, for everyone signed in: a person with the role topology edits it there, and anyone else looks at it and prints where they may. Admin in the bar, and the Topology tab of the administrator's pages, lead to the same page.

## The canvas {#canvas}

The canvas draws the nodes where the administrator placed them. A node is a white box with a title bar colored by its role, an icon of its kind, its name, its `[kind]` and title under it, and the first lines of its description:

| Role | Kinds |
|---|---|
| Source | uploader, CUPS server |
| Control flow | queue, class, inbox |
| Filter | filter |
| Transform | template, convert |
| Sink | display, printer, archive, sink, downloader |

The ports and the edges are colored by what flows along them ([Kinds of payload](flow.md#kinds)): blue a `.bdf` document, green data, amber any file as received, cyan a stocked file (from an inbox or an archive to a downloader), purple several kinds. A port not typed yet is hollow and dashed: a queue, a class, a filter or an inbox takes the kind of its first edge. An edge drawn dashed and red no longer connects ([Broken edges](flow.md#broken)). A downloader has a dashed arrow going out: applications and people fetch what it lists. The legend under the canvas names the colors.

The activity of the jobs you may see runs on it — yours and your groups'; an administrator sees everyone's:

- a count at each node: items held in an inbox, slips open on a display, deliveries waiting at a printer or an archive, deliveries that failed, entries a downloader lists, and the jobs that reached the node in the last hour;
- the first page of the newest document at a sink (kept), a display (a slip open on its board), an inbox (an item held) and an archive (written), of those you see;
- dots running along the edges: the legs of the last ten minutes.

Drag the background to pan, turn the wheel or pinch to zoom, and Fit brings the whole topology into view. The page polls the manager every few seconds; the badge beside the title says when the activity last changed, or that the manager is not reached.

## Looking {#view}

Without the role topology, the page only looks. Click a node, and the sidebar shows its description, the counts of your jobs there, where it sends, what it takes and sends on, and:

- **Print here**, at a source you may send to ([Scopes of nodes](roles.md#scopes)): a file and a name for the job, as on [Jobs](jobs.md#print); at a source you may not send to, "You may not print here"; at another node, that documents come in at a source;
- **Jobs here**: your jobs that were sent to the node or reached it, with their state, who sent them, when and their route; Document downloads one, and Reprint here sends it again to this node, as a new job of yours.

Under the counts, a node shows what is at it now, read again as the activity changes: a sink the latest documents it kept; a display the slips on its board that you see, open and lately done or canceled ([Display boards](displays.md#board)); an inbox what it holds, with Release and Discard when you work it ([Inboxes](inboxes.md#inboxes)); an archive the latest documents it wrote that you see; a downloader the entries listed there that you may read ([Downloaders](downloads.md)). Each document has its first page: a click shows every page ([Looking at a document](#lightbox)). A CUPS server (source) shows how a computer adds it as a printer ([CUPS server nodes](entries.md#cups)).

### Looking at a document {#lightbox}

A click on a document's first page, in the sidebar or on the Inboxes page, shows every page large, one under another, with Download beside its title, and of an inbox's item Release and Discard. The bar under the title zooms: − and + step out and in, the percentage is the page's width as a part of the box's, **Fit** shows each page whole in the box (as it opens) and **Width** as wide as the box. The keys +, − and 0 (Fit), Ctrl or ⌘ with the wheel, and a pinch on a trackpad or a screen zoom too, around where they point; drag a page larger than the box to move it. Zoomed in, a page is drawn again, sharper. Escape, ✕ or a click beside the box closes it.

## The editor {#editor}

With the role topology, the same page is the editor: the canvas, and the sidebar for the node selected, with the administrator's tabs above.

- **+ Add a node** opens the palette of kinds, under four headings: Source, Control (the roles control flow and filter), Convert (transform) and Sink. Choose one, and the new node appears on the canvas with the sidebar open on it: give it a title, a name and its settings, and press Add; Escape drops it. A node needs no edge to be added: add the nodes in any order, and connect them once they are there. A job that reaches an entry, a queue or a step that sends to no node yet gets an error there, "sends to no node yet", on its route.
- **Click a node** to open it in the sidebar: its settings under Settings, how it is reached and what is at it under Usage ([The sidebar of a node](#sidebar)). Save writes the node; Delete, or the Delete key, removes it. A node another node sends to cannot be deleted: remove the edge first.
- **Drag a node** to move it. Where the nodes are drawn is saved for everyone (the layout), apart from the nodes themselves.
- **Drag from a node's right port onto another node** to connect them. While you drag, the nodes that would not take what the edge carries fade, the wire turns red over one with the reason, and it snaps to the port of one that takes it; a drop on one that does not is refused before the server is asked. The edge is added to the node's `to` (a queue, an uploader, a CUPS server, an inbox, a template or a convert node, an archive), to its `members` (a class), to its `otherwise` (a filter: a rule's nodes are set in the rule), or it sets the fallback printer (a display). A branch has a port for each case and one for its default, an assert one for what passes (`to`) and one for what fails (`on_error`), each labelled: the edge goes to the arm of the port it is dragged from. Click an edge to remove it, after a confirmation.
- Zoom, Fit and Refresh.

A node is checked when it is saved: the nodes it names exist; the name is not another kind's; no path leads back to the node; and a new edge brings no kind of payload the node after it does not take, with a message that names the path: "laser (printer) takes bdf, but any reaches it: raw → fan → laser". A change of a node's settings alone is saved even when it leaves an edge carrying what the node after it does not take: the edge stays, broken, and the save says so ([Broken edges](flow.md#broken)). A line under the settings says what the node takes, what may be sent to it now, and what it sends on; a broken edge from it is explained there, with the ways to mend it.

### What is not saved yet {#unsaved}

A node added from the palette is a draft until Add puts it: drawn dashed, "not saved yet", and kept only in this page. Clicking elsewhere before Add leaves it on the canvas with what was typed for it: click it to go on, then Add it, or Cancel it (Escape). Several drafts may wait at once; a draft has no edges until it is added. A saved node whose settings were changed and not saved is kept the same way: it is ringed in amber, and opened again it shows the changes, with Save and Revert (back to the node as it is). An edge dragged to or from it meanwhile stays: only the fields you changed are kept. Leaving the page with drafts or changes not saved asks first; they are not kept past the page.

A required field left empty is marked in amber once you leave it, with what it needs under it: the name, a printer's host, an archive's store, a CUPS server (sink)'s host and queue, an S3 Trigger's bucket, and the settings an adapter's kind needs (a mail's `url`, `from` and `recipients`; a notify node's `credential` by slack). A value that is wrong stays red. Add and Save mark every required field left empty and go to the first, without sending the node. A draft or a node left with required fields empty has a badge on the canvas counting them, and opened again, all of them are marked.

### Names, titles and descriptions {#names}

A node has a title and a name. The title is what people read, in any characters: it is shown in the node's title bar on the canvas (the name, while it has none). The name is the system's: letters, digits, `.`, `_` and `-`, up to 64, starting with a letter or a digit, and no other node's; the edges of other nodes, the URLs (`/api/upload/<name>`, `/api/inbox/<name>`, `/api/download/<name>`, a CUPS server's `/printers/<name>`) and the records call the node by it. The Name field is checked as you type: a space, another character, or a name another node has turns it red, with what is wrong under it, and the node is not saved until it is right.

Changing the name of a saved node renames it when it is saved, after a confirmation: the edges of the nodes that send to it, where it is drawn, and what is kept under its name move with it — a template node's templates, a display's devices and the slips it still shows, an inbox's items, a sink's documents, a downloader's entries — and an archive keeps writing to its collection. Where it is reached from outside changes with it: computers and applications set up with the old name must be set up again. The jobs that reached it keep the old name. A node with work on its way (deliveries waiting or failed, items held in an inbox, slips open on a display) is renamed once that is done; so is one the manager's `docflow.on_error` names, once the configuration names the new one.

The description is a note of what the node is for — "the invoices of the ERP, PDF only; printed on the 3rd floor" — up to 1000 characters, as the description of an element of a C4 diagram: its first lines in the node on the canvas, all of it in the sidebar. A display, a printer, a sink and a downloader keep their kind; a node of any other kind may be changed into another of those kinds, keeping its name.

### The sidebar of a node {#sidebar}

The sidebar of a saved node has two tabs. The one chosen stays as you click from node to node.

**Settings**: the title, the name and the description, the settings of its kind ([Node kinds](nodes.md)), its scope and owner, Save and Delete; of a template node, its templates too: upload one, the fields each looks up, its file, delete ([Templates and data](templates.md)).

**Usage**: how the node is reached from outside, and what is at it:

- of a CUPS server (source), the URL computers add the printer with, how macOS, Windows and CUPS add it, and the commands ([CUPS server nodes](entries.md#cups));
- of a printer or a CUPS server (sink), where this system prints the jobs, over IPP (a printer also over 9100 or LPD) ([CUPS server (sink)](nodes.md#cups-out));
- of an uploader, the request applications send and the apps allowed, and its key: Rotate key makes a new one and shows it once ([Uploaders](entries.md#uploaders)); of an S3 trigger, the request an S3 event's reference is posted with, and its key;
- of a display, the slips on its board, then its devices, when each was last seen, and Add a device, which shows a one-time code and its QR code ([Display boards](displays.md));
- of an inbox, what it holds, with Release and Discard, as on the Inboxes page, then the API that applications work it with ([Inboxes](inboxes.md));
- of an archive, the latest documents it wrote that you see, with their pages, drawn from its store;
- of a downloader, the URLs readers fetch from, and the entries listed there ([Downloaders](downloads.md));
- links to the jobs sent to it (Admin → Jobs), to the inbox it is (Inboxes), to the deliveries waiting at it (Admin → Deliveries), and the OpenAPI document of an uploader, an inbox or a downloader;
- of a sink, the latest documents it kept, with their pages; of a printer with a preview, what the fake printer received ([Printer](nodes.md#printer)).

What is at a node is read again as the activity changes, and what you chose in it (the nodes an item goes to) stays. The Usage tab of an inbox or a display counts what waits there: items held, slips open. The other nodes are reached along the edges alone, and their Usage has the links. Help at the top of the sidebar opens the section of the node's kind in this manual.
