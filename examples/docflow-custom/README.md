# docflow-custom

BDF Docflow composed by a program of its own: the web application of
`docflow/web` with every input and output of `docflow/io`, and a kind of
node of this program's own, the **folder** exit (`kinds/folder`), which
writes what reaches it to a directory of the server. It is what a program
that imports bdf and adds kinds of node, settings and pages of its own
looks like (`docflow gen-project` writes one).

- A kind of node of one's own is registered with `docflow.RegisterKind`
  in an `init`, with its name and help in each language, its icon, its
  chapter of the manual and the **fields** of its settings
  (`docflow.Field`): the editor lays them out as a form, marks the
  required ones, and `PutNode` checks a node's settings against them
  before the kind's own `Check`; `Send` reads them with
  `docflow.SettingsOf`. The program imports the package to have the kind.
- The program is a Popcorn Web project (`popcornweb.toml` beside
  `go.mod`): `pw dev` runs it with pw's development identity provider
  (its people are `devidp.toml`, and `config.dev.toml` makes
  `docflow-admins` the administrators), and `pw build` makes
  `bin/docflow-custom`. Its `main` is `web.RegisterConfig()` (and the
  program's own settings, with `pw.RegisterConfig`) and
  `web.Main(web.Options{…})`.
- Its own public tree (`public`, served at `/public`) is empty: Docflow's
  scripts and style are served by `docflow/web` at `/_docflow/`.

```sh
go mod tidy
go test ./...
pw dev
```

On `/topology`, the palette lists **Folder** under Sink with its fields;
a job sent to an uploader in front of it is written to the directory.
