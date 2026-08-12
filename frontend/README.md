# Vue 3 + Vite

This template should help get you started developing with Vue 3 in Vite. The template uses Vue 3 `<script setup>` SFCs,
check out the [script setup docs](https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup) to learn more.

## Recommended IDE Setup

- [VS Code](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar)

## Testing

Install the Chromium runtime once after installing dependencies:

```bash
npm run test:e2e:install
```

Run the frontend state tests and browser E2E tests:

```bash
npm test
npm run test:e2e
```

`test:e2e` starts the Wails v2.10.2 development bridge, uses a temporary home
directory for game data, and runs serially against the real Go backend. On a
machine whose system Node.js is older than Node 20, the runner uses an isolated
Node 22 runtime automatically. Use `npm run test:e2e:ui` for Playwright UI mode.

The native GNOME/Wayland smoke test expects a current desktop binary:

```bash
npm run build
cd ..
go run github.com/wailsapp/wails/v2/cmd/wails@v2.10.2 build -tags webkit2_41 -s -skipbindings
cd frontend
npm run test:native
```

The native test uses AT-SPI to inspect and interact with the real Wails window.
It temporarily enables GNOME toolkit accessibility when necessary and restores
the previous setting before exiting.
