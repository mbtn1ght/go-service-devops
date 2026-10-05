# Deploy to the VPS

The normal deployment path is driven from the development machine:

```sh
git add .
git commit -m "Describe the change"
make deploy
```

`make deploy` pushes the current clean branch and runs `make update`, which
connects to the VPS, resets `/opt/service-devops` to the deployed branch,
rebuilds
the Docker image, recreates the service, and checks `/healthz`.

## One-time server setup

On the VPS, install Docker and Git, then clone the repository to the expected
location:

```sh
apt update
apt install -y ca-certificates curl git
curl -fsSL https://get.docker.com | sh
git clone https://github.com/mbtn1ght/go-service-devops.git /opt/service-devops
```

From the development machine, install your SSH public key for non-interactive
deployments:

```sh
ssh-copy-id root@193.233.246.215
```

For the current VPS, this setup is already complete. The deployment checkout
uses the public HTTPS repository URL, so it needs no GitHub credential to pull.

## Operations

```sh
make update  # rebuild and restart the version already in GitHub
make logs    # local logs
```

On the VPS:

```sh
cd /opt/service-devops
docker compose logs -f app
docker compose ps
```
