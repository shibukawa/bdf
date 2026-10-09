# @bdfkit/convert-office

Browser converter preset for Word, Excel, PowerPoint, Visio, InDesign (IDML), CSV, TSV, Parquet, PDF, and EPUB. It bundles a TinyGo WebAssembly runtime for this format group.

```sh
npm install @bdfkit/convert-office
```

```ts
import { createOfficeConverter } from "@bdfkit/convert-office";

const converter = createOfficeConverter();
const result = await converter.convert(bytes, { name: "report.xlsx" });
```

Use it with `Viewer.openFile` from `@bdfkit/viewer`, or use the converter directly. For a smaller PDF/EPUB runtime or a custom format combination, see [`@bdfkit/convert-pdf-epub`](https://github.com/shibukawa/bdf/tree/main/packages/convert-pdf-epub) and the [Wasm runtime build guide](https://github.com/shibukawa/bdf/blob/main/docs/build-wasm-runtime.md).

## License

MIT. See [LICENSE](LICENSE).
