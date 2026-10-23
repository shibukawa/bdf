# Topology

The topology is the network the documents flow through: the nodes, of every kind, and the edges between them. It is one page, Topology, for everyone signed in: a person with the role topology edits it there, and anyone else looks at it and prints where they may. Admin in the bar, and the Topology tab of the administrator's pages, lead to the same page.

## The canvas {#canvas}

The canvas draws the nodes where the administrator placed them. A node is a white box with a title bar colored by its role, an icon of its kind, its name, its `[kind]` and title under it, and the first lines of its description:

| Role | Kinds |
|---|---|
| Source | endpoint, CUPS server |
| Control flow | queue, class, inbox |
| Filter | filter |
| Transform | template, convert |
| Sink | display, printer, archive, sink, HTTP (pull) |

The ports and the edges are colored by what flows along them ([Kinds of payload](flow.md#kinds)): blue a `.bdf` document, green data, amber any file as received, cyan a stocked file (from an inbox or an archive to HTTP (pull)), purple several kinds. A port not typed yet is hollow and dashed: a queue, a class, a filter or an inbox takes the kind of its first edge. An edge drawn dashed and red no longer connects ([Broken edges](flow.md#broken)). An export has a dashed arrow going out: applications and people fetch what it lists. The legend under the canvas names the colors.

The activity of the jobs you may see runs on it — yours and your groups'; an administrator sees everyone's:

- a count at each node: items held in an inbox, slips open on a display, deliveries waiting at a printer or an archive, deliveries that failed, entries an export lists, and the jobs that reached the node in the last hour;
- dots running along the edges: the legs of the last ten minutes.

Drag the background to pan, turn the wheel or pinch to zoom, and Fit brings the whole topology into view. The page polls the manager every few seconds; the badge beside the title says when the activity last changed, or that the manager is not reached.

## Looking {#view}

Without the role topology, the page only looks. Click a node, and the sidebar shows its description, the counts of your jobs there, where it sends, what it takes and sends on, and:

- **Print here**, when you may send to the node ([Scopes of nodes](roles.md#scopes)): a file and a name for the job, as on [Jobs](jobs.md#print); otherwise "You may not print here";
- **Jobs here**: your jobs that were sent to the node or reached it, with their state, who sent them, when and their route; Document downloads one, and Reprint here sends it again to this node, as a new job of yours.

A sink also shows the latest documents it kept, and an export the entries listed there that you may read ([HTTP (pull)](exports.md)).

## The editor {#editor}

With the role topology, the same page is the editor: the canvas, and the sidebar for the node selected, with the administrator's tabs above.

- **+ Add a node** opens the palette of kinds, under four headings: Source, Control (the roles control flow and filter), Convert (transform) and Sink. Choose one, and the new node appears on the canvas with the sidebar open on it: give it a name and its settings, and press Add; Escape drops it.
- **Click a node** to open it in the sidebar: its settings, and what is at it ([The sidebar of a node](#sidebar)). Save writes the node; Delete, or the Delete key, removes it. A node another node sends to cannot be deleted: remove the edge first.
- **Drag a node** to move it. Where the nodes are drawn is saved for everyone (the layout), apart from the nodes themselves.
- **Drag from a node's right port onto another node** to connect them. While you drag, the nodes that would not take what the edge carries fade, the wire turns red over one with the reason, and it snaps to the port of one that takes it; a drop on one that does not is refused before the server is asked. The edge is added to the node's `to` (a queue, an endpoint, a CUPS server, an inbox, a template or a convert node, an archive), to its `members` (a class), to its `otherwise` (a filter: a rule's nodes are set in the rule), or it sets the fallback printer (a display). Click an edge to remove it, after a confirmation.
- Zoom, Fit and Refresh.

A node is checked when it is saved: the nodes it names exist; the name is not another kind's; no path leads back to the node; and a new edge brings no kind of payload the node after it does not take, with a message that names the path: "laser (printer) takes bdf, but any reaches it: raw → fan → laser". A change of a node's settings alone is saved even when it leaves an edge carrying what the node after it does not take: the edge stays, broken, and the save says so ([Broken edges](flow.md#broken)). A line under the settings says what the node takes, what may be sent to it now, and what it sends on; a broken edge from it is explained there, with the ways to mend it.

### Names, titles and descriptions {#names}

A node's name is letters, digits, `.`, `_` and `-` (up to 64), and it does not change later: it is the node's IPP destination, `/ipp/print/<name>`, and what applications and other nodes refer to. The title is shown under the name. The description is a note of what the node is for — "the invoices of the ERP, PDF only; printed on the 3rd floor" — up to 1000 characters, as the description of an element of a C4 diagram: its first lines in the node on the canvas, all of it in the sidebar. A display, a printer, a sink and an export keep their kind; a node of any other kind may be changed into another of those kinds, keeping its name.

### The sidebar of a node {#sidebar}

Besides the settings of its kind ([Node kinds](nodes.md)), the sidebar of a saved node shows:

- its IPP URL, and links to the jobs sent to it (Admin → Jobs), to the inbox it is (Admin → Inboxes) and to the deliveries waiting at it (Admin → Deliveries);
- of an endpoint, its entry (the URL applications POST to) and its key: Rotate key makes a new one and shows it once ([Endpoints](entries.md#endpoints));
- of a template node, its templates: upload one, the fields each looks up, its file, delete ([Templates and data](templates.md));
- of a display, its devices, when each was last seen, and Add a device, which shows a one-time code and its QR code ([Display boards](displays.md));
- of a CUPS server node, the URL computers add the printer with, and the command that lists it ([CUPS server nodes](entries.md#cups));
- of a sink, the latest documents it kept, with their pages;
- of an export, the URLs readers fetch from, and the entries listed there ([HTTP (pull)](exports.md));
- of a printer with a preview, what the fake printer received ([Printer](nodes.md#printer)).

Help at the top of the sidebar opens the section of the node's kind in this manual.
