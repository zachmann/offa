
The following `docker-compose.yaml` is an example on how to set up offa next to your service with traefik (it
assumes an already running `traefik`, but it can also be added to the docker compose file).

This assumes that the `docker-compose.yaml` lies next to your offa's `config.yaml` file.

```
services:

  offa:
    image: oidfed/offa:main
    restart: unless-stopped
    volumes:
      - ./config.yaml:/config.yaml:ro
      - ./:/data
    expose:
      - 15661
    labels:
      - traefik.enable=true
      - traefik.port=15661
      - traefik.http.routers.https-offa.entryPoints=https
      - traefik.http.routers.https-offa.rule=Host(`offa.example.com`)
      - traefik.http.routers.https-offa.tls=true
      - traefik.http.routers.https-offa.tls.certresolver=le
      - traefik.http.middlewares.offa.forwardauth.address=https://offa.example.com/auth
      - traefik.http.middlewares.offa.forwardauth.trustForwardHeader=true
      - traefik.http.middlewares.offa.forwardauth.authResponseHeaders=X-Forwarded-User,X-Forwarded-Groups,X-Forwarded-Name,X-ForwardedEmail,X-Forwarded-Provider,X-Forwarded-Subject

  whoami:
    image: containous/whoami
    labels:
      - traefik.enable=true
      - traefik.http.routers.https-whoami.rule=Host(`whoami.example.com`)
      - traefik.http.routers.https-whoami.entryPoints=https
      - traefik.http.routers.https-whoami.tls=true
      - traefik.http.routers.https-whoami.tls.certresolver=le
      - traefik.http.routers.https-whoami.middlewares=offa@docker
    restart: unless-stopped

```
