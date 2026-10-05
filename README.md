# Service DevOps

Go service with a Vite/React admin dashboard. The dashboard is built into the
same Docker image and is served by the Go process.

## Local workflow

```sh
make run       # frontend build + Go server at http://localhost:8080
make up        # full Docker Compose stack in background
make logs      # follow application logs
make down      # stop the stack
make test
```

`make run` uses `corepack pnpm` and will install the locked frontend
dependencies before building it. Docker builds use the same lockfile.

## Deployment workflow

The VPS is a disposable checkout of the branch passed to `make deploy` at
`/opt/service-devops`. Do not edit application files directly on it: `make
update` resets that directory to the remote branch and removes untracked
files.

Configure key-based SSH access once (the account must already be authorized
to access the repository from the VPS):

```sh
ssh-copy-id root@193.233.246.215
```

Then commit your changes and use:

```sh
make push      # push the current clean branch
make update    # pull the current branch on the VPS and rebuild/restart Docker
make deploy    # push, update the VPS, and run its health check
```

Targets accept overrides, for example:

```sh
make deploy VPS_HOST=203.0.113.10 BRANCH=main
```
