# jqplay

[![OpenCollective](https://opencollective.com/jqplay/backers/badge.svg)](#backers) [![OpenCollective](https://opfcfollective.com/jqplay/sponsors/badge.svg)](#sponsors)

[jqplay](https://jqplay.org) is a playground for [jq](https://github.com/jqlang/jq). Please put it into good use.

This version runs entirely in the browser using [jq-wasm](https://github.com/owenthereal/jq-wasm) - no backend required!

## Development

To develop `jqplay`, you need to have [Bun](https://bun.sh/) installed.

### Install dependencies

```bash
bun install
```

### Start development server

```bash
bun run dev
```

Point your browser to [`http://localhost:8080/`](http://localhost:8080/).

### Build for production

```bash
bun run build
```

The built files will be in the `dist/` directory.

### Preview production build

```bash
bun run preview
```

## How it works

This is a pure frontend application that uses:

- [jq-wasm](https://github.com/owenthereal/jq-wasm) - WebAssembly-powered jq running in the browser
- [Ace editor](https://ace.c9.io/) - Code editor for input/output
- [Bootstrap](https://getbootstrap.com/) - UI framework
- [Vite](https://vitejs.dev/) - Build tool and dev server

All jq processing happens client-side in WebAssembly, so there's no backend server to maintain.
