# Display boards and devices

A display is a node that shows jobs instead of printing them: on the devices enrolled at it — the tablets of a warehouse, the screen at a dock. A job sent to a display becomes a slip on its board; what nobody opens in time is printed on paper, on the display's fallback printer, so that shipping does not stop.

## The display {#display}

A display's settings are in its sidebar on the Topology page ([Display](nodes.md#display)): its fallback printer and after how long an unopened slip is printed, how long ended slips stay on the board, and how long a device may be silent and still count as connected. While no device of the display is connected, every unopened slip is printed. A display takes `.bdf` documents only; put a convert node before it for files as received.

The jobs sent to a display are its group's ([Scopes of nodes](roles.md#scopes)); a display with no group of its own gives its jobs its fallback printer's group.

## Enrolling a device {#enroll}

A device is enrolled with a one-time code an administrator makes for the display:

1. On the Topology page, open the display and press **Add a device (QR code)** under Devices. The dialog shows a code of ten characters, its QR code, the URL `/tablet?code=…`, and until when the code is valid (15 minutes).
2. On the device, open the URL (scan the QR code, or open `/tablet` and type the code), give the device a name — "tablet at the shipping dock" — and enroll it.
3. Add the page to the home screen: it is the device's board from now on.

The device's page makes a key that cannot leave the browser (WebCrypto, kept in the browser's own storage) and sends only its public half: from then on the device signs in by signing a challenge, and a URL or a cookie copied to another device does not let that device in. A device signs in again by itself when its session ends (every ten minutes); a device the manager no longer knows is asked to enroll again. The code is used once.

The display's sidebar lists its devices, their names and when each was last seen, and Remove takes one off: it is out at once, and must be enrolled again to come back. Codes, enrollments and removals are in the audit log.

## Devices on the network {#network}

A device makes its key with WebCrypto, which a browser gives only to a page on HTTPS, or on localhost (the computer itself). An enrollment's URL is therefore at an origin the devices reach:

- `docflow.device_url` (`DOCFLOW_DEVICE_URL`) when it is set: the manager's HTTPS address, such as `https://docflow.example.com`.
- Otherwise the manager as the administrator reached it. When that is localhost — the administrator's browser on the manager's computer — the URL is at the computer's address on the network instead (`192.168.1.7`); and when the manager listens on `docflow.https_port`, the URL goes there.

`docflow.https_port` (`DOCFLOW_HTTPS_PORT`) is an HTTPS listener that serves what the HTTP one does, with the certificate of `docflow.tls_cert_file` and `docflow.tls_key_file`. In development it goes without them: the manager makes a certificate of the computer signed by itself and keeps it in the user's cache (`bdf/docflow-dev-tls.pem`), so a device's browser warns about it the first time; go on to the page (Advanced, then Proceed). The development settings (`config.dev.toml`) listen on 8443. The dialog of a code warns when its URL is plain http on the network, and the board says so when it is opened there.

## The board {#board}

The board (`/tablet`) shows the slips of the display under two tabs: **Open**, the slips new and opened, and **Done and canceled**. A new slip arrives with a sound and a notice; the board polls the manager every few seconds and shows the display's title.

- Tap a slip to open it: its pages, drawn large. Opening it tells the manager, so that it is not printed on paper.
- **Done** marks the slip done; it moves to Done and canceled, and stays there for the display's keep.
- **Print here** prints the slip on the display's fallback printer, when it has one; the slip stays open.
- The name in the bar is the worker's, written on what they do.

On the topology, the display's sidebar shows the slips on its board — those of the jobs you see; every one, for the role topology —, open and the latest done or canceled, each with its state and its first page; a click shows its pages ([Looking at a document](topology.md#lightbox)).

A job canceled on the Jobs pages leaves the board. A slip the fallback printed on paper — because nobody opened it in time, or no device was connected — is marked "Printed on paper" and stays open until someone marks it done.
