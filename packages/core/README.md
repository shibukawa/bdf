# @bdfkit/core

Low-level TypeScript APIs for reading BDF (Browser-specific Document Format) data.

```sh
npm install @bdfkit/core
```

Use this package when an application needs direct access to the BDF container, manifest, objects, text, search index, or music data. It does not create a DOM viewer and has no browser UI. The public API is exported from the package root:

```ts
import { BdfDocument, BufferSource, extractText } from "@bdfkit/core";
```

For a ready-to-use renderer, see [`@bdfkit/render`](https://github.com/shibukawa/bdf/tree/main/packages/render) and [`@bdfkit/viewer`](https://github.com/shibukawa/bdf/tree/main/packages/viewer). See the [npm integration guide](https://github.com/shibukawa/bdf/blob/main/docs/npm.md) for package combinations and server conversion.

## License

MIT. See [LICENSE](LICENSE).
