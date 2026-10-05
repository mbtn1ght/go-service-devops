# Frontend

The UI is based on [shadcn-admin](https://github.com/satnaing/shadcn-admin),
released under the MIT license (see `LICENSE`). It is served by the Go service
in production.

## Local development

Start the API in one terminal:

```sh
go run ./cmd
```

Then run the Vite development server in another:

```sh
cd web
corepack enable
pnpm install
pnpm dev
```

Build the deployable application from the repository root:

```sh
docker compose up --build
```
