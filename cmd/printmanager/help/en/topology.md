# Topology

The topology is the network the documents flow through: the nodes, of every kind, and the edges between them. Everyone signed in sees it on the Topology page; the role topology lays it out on Admin → Topology.

## The view {#view}

The Topology page draws the nodes where the administrator placed them, with their kinds and the edges, and the activity of the jobs you may see — yours and your groups'; an administrator sees everyone's:

- a count at each node: items held in an inbox, slips open on a display, deliveries waiting at a printer or an archive, deliveries that failed, and the jobs that reached the node in the last hour;
- dots running along the edges: the legs of the last ten minutes.

Drag the background to pan, turn the wheel or pinch to zoom, and Fit brings the whole topology into view. Click a node to list your jobs there, and to print there when you may send to it (the same form as on [Jobs](jobs.md#print)). The page polls the manager, and the badge beside the title says when the activity last changed. Nothing is changed from this page.

## The editor {#editor}

Admin → Topology is the editor: the canvas, and the sidebar for the node selected.

- **+ Add a node** opens the palette of kinds. Choose one, and the new node appears on the canvas with the sidebar open on it: give it a name and its settings, and press Add; Escape drops it.
- **Click a node** to open it in the sidebar: its settings, and what is at it ([The sidebar of a node](#sidebar)). Save writes the node; Delete, or the Delete key, removes it. A node another node sends to cannot be deleted: remove the edge first.
- **Drag a node** to move it. Where the nodes are drawn is saved for everyone (the layout), apart from the nodes themselves.
- **Drag from a node's right port onto another node** to connect them. The edge is added to the node's `to` (a queue, an endpoint, a CUPS server, an inbox, a template or a convert node), to its `members` (a class), to its `otherwise` (a filter: a rule's nodes are set in the rule), or it sets the fallback printer (a display). Click an edge to remove it, after a confirmation.
- Zoom, Fit and Refresh, as on the view.

A node is checked when it is saved: the nodes it names exist; the name is not another kind's; no path leads back to the node; and every node takes the kinds of payload that reach it ([Kinds of payload](flow.md#kinds)). The ports and the edges are colored by what flows along them — blue a `.bdf` document, green data, grey any file as received — and an edge that would bring a kind a node does not take is refused, with a message that names the path: "laser (printer) takes bdf, but any reaches it: raw → fan → laser". A line under the settings says what the node takes and what it sends on, as inferred from the entries.

### Names and titles {#names}

A node's name is letters, digits, `.`, `_` and `-` (up to 64), and it does not change later: it is the node's IPP destination, `/ipp/print/<name>`, and what applications and other nodes refer to. The title is shown instead of the name where there is room. A display, a printer and a sink keep their kind; a node of any other kind may be changed into another of those kinds, keeping its name.

### The sidebar of a node {#sidebar}

Besides the settings of its kind ([Node kinds](nodes.md)), the sidebar of a saved node shows:

- its IPP URL, and links to the jobs sent to it (Admin → Jobs), to the inbox it is (Admin → Inboxes) and to the deliveries waiting at it (Admin → Deliveries);
- of an endpoint, its entry (the URL applications POST to) and its key: Rotate key makes a new one and shows it once ([Endpoints](entries.md#endpoints));
- of a template node, its templates: upload one, the fields each looks up, its file, delete ([Templates and data](templates.md));
- of a display, its devices, when each was last seen, and Add a device, which shows a one-time code and its QR code ([Display boards](displays.md));
- of a CUPS server node, the URL computers add the printer with, and the command that lists it ([CUPS server nodes](entries.md#cups));
- of a sink, the latest documents it kept, with their pages;
- of a printer with a preview, what the fake printer received ([Printer](nodes.md#printer)).
