# @bdfkit/convert-all

Browser converter preset containing every browser converter in BDF. It bundles the largest TinyGo WebAssembly runtime and is intended for demos or applications that need broad format coverage.

```sh
npm install @bdfkit/convert-all
```

```ts
import { createAllConverter } from "@bdfkit/convert-all";

const converter = createAllConverter();
const result = await converter.convert(bytes, { name: "input.pptx" });
```

Use it with `Viewer.openFile` from `@bdfkit/viewer`, or use the converter directly. If bundle size matters, choose `@bdfkit/convert-pdf-epub`, `@bdfkit/convert-office`, or build a custom runtime with the [Wasm runtime build guide](https://github.com/shibukawa/bdf/blob/main/docs/build-wasm-runtime.md). Music converters are included in this preset; custom builds can omit them.

## License

MIT. See [LICENSE](LICENSE).
