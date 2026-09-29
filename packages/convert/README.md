# @bdfkit/convert

Worker client for running a TinyGo WebAssembly converter runtime in the browser. This package is the common API used by the converter presets and is useful when building a custom Wasm runtime with a selected set of formats.

```sh
npm install @bdfkit/convert
```

Pass the URL of a compatible converter Wasm module to `createConverter`:

```ts
import { createConverter } from "@bdfkit/convert";

const converter = createConverter(new URL("./converter.wasm", import.meta.url));
const result = await converter.convert(inputBytes, { name: "report.pdf" });
converter.terminate();
```

The package starts a dedicated Worker and exposes `convert`, streaming `open`/`page`/`finish`, and `close`. It does not include a converter Wasm file by itself. Use a preset package or follow the [Wasm runtime build guide](https://github.com/shibukawa/bdf/blob/main/docs/build-wasm-runtime.md).

## License

MIT. See [LICENSE](LICENSE).
