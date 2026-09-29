# @bdfkit/convert-pdf-epub

Browser converter preset for PDF and EPUB files. It bundles the TinyGo WebAssembly runtime, so consumers only need Node/npm and their frontend bundler.

```sh
npm install @bdfkit/convert-pdf-epub
```

```ts
import { createPdfEpubConverter } from "@bdfkit/convert-pdf-epub";

const converter = createPdfEpubConverter();
const result = await converter.convert(bytes, { name: "book.epub" });
```

Use it with `Viewer.openFile` from `@bdfkit/viewer`, or use the converter directly when the application owns the document surface. See the [frontend integration guide](https://github.com/shibukawa/bdf/blob/main/docs/integrating-to-frontend.md).

## License

MIT. See [LICENSE](LICENSE).
