# GoChat web

GoChat のフロントエンド（React + TypeScript、[Vite](https://vite.dev/) でビルド）。

## Available Scripts

### `npm start` / `npm run dev`

開発サーバーを [http://localhost:3000](http://localhost:3000) で起動する。\
安全のため 127.0.0.1 のみで待ち受ける（`vite.config.ts` の `server.host`）。

### `npm test`

[Vitest](https://vitest.dev/) でテストを実行する（watch モード）。一度だけ実行する場合は `npx vitest run`。

### `npm run build`

型チェック（`tsc --noEmit`）の後、本番用ビルドを `dist/` に出力する。

### `npm run preview`

`dist/` のビルド結果をローカルで確認する。
