# Deploy to a VPS

These commands assume Ubuntu 24.04 and a non-root user with `sudo` access.

1. On the server, install Docker:

   ```sh
   sudo apt update
   sudo apt install -y ca-certificates curl
   curl -fsSL https://get.docker.com | sudo sh
   sudo usermod -aG docker "$USER"
   ```

   Sign out and back in for the group change to take effect.

2. Copy this project to the server, for example to `/opt/service-devops`. If the
   repository is hosted, clone it there instead.

   ```sh
   sudo mkdir -p /opt/service-devops
   sudo chown "$USER":"$USER" /opt/service-devops
   rsync -az --delete ./ user@SERVER_IP:/opt/service-devops/
   ```

3. Start it:

   ```sh
   cd /opt/service-devops
   docker compose up -d --build
   docker compose ps
   curl http://127.0.0.1:8080/healthz
   ```

4. Allow TCP port 8080 in the VPS provider firewall and, if UFW is enabled:

   ```sh
   sudo ufw allow 8080/tcp
   ```

The service will be available at `http://SERVER_IP:8080/` and will restart
after a server reboot. For a public production service, put Caddy or Nginx in
front of it and expose only ports 80/443; that also enables TLS certificates.
