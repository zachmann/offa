# OFFA - Openid Federation Forward Auth

<img alt="logo" height="200" src="logos/offa-text.svg"/>

OFFA offers easy to use OpenID Federation Authentication and Authorisation for existing services.
OFFA can be deployed along existing services and handle all OpenID Federation communication for your services.

OFFA implements Forward Authentication usable with
[Traefik](https://doc.traefik.io/traefik/middlewares/http/forwardauth/),
[NGINX](https://docs.nginx.com/nginx/admin-guide/security-controls/configuring-subrequest-authentication/),
[Caddy](https://caddyserver.com/docs/caddyfile/directives/forward_auth),
and maybe other reverse proxies.

OFFA also implements [Auth MemCookie](https://zenprojects.github.io/Apache-Authmemcookie-Module/) usable with Apache.


## OFFA Configuration
TODO, for now please refer to [internal/config/config.go](internal/config/config.go).


## How to Deploy With
In this section we detail how to deploy OFFA with your reverse proxy.

- [Traefik](#traefik)
- [NGINX](#nginx)
- [Caddy](#caddy)
- [Apache](#apache)

### Traefik

Still needs testing

### NGINX

Still needs testing

### Caddy

The following example configuration can be used (tweak as needed):

We assume the following project layout:
```ascii
.
├── caddy
│   ├── Caddyfile
│   ├── config
│   └── data
├── docker-compose.yaml
├── offa
│   └── config.yaml
```


#### docker-compose.yaml
```yaml
services:
  caddy:
    image: caddy:latest
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./caddy/Caddyfile:/etc/caddy/Caddyfile
      - ./caddy/data:/data
      - ./caddy/config:/config
    networks:
      caddy:

  offa:
    image: myoidc/offa
    restart: unless-stopped
    volumes:
      - ./offa/config.yaml:/config.yaml:ro
      - ./offa:/data
    networks:
      caddy:

  # This would be your service
  hello:
    image: plippe/hello-world-web-service
    restart: unless-stopped

networks:
  caddy:
```

#### caddy/Caddyfile

```
offa.example.com {
  reverse_proxy offa:15661
}

hello.example.com {
    forward_auth offa:15661 {
		uri /auth
		copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
	}

    reverse_proxy hello:5000
}
```

#### offa/config.yaml

```yaml
server:

logging:
  access:
    stderr: true
  internal:
    level: info
    stderr: true

sessions:
  ttl: 3600
  cookie_domain: example.com

auth:
  - domain: hello.example.com
    require:
      groups: users

federation:
  entity_id: https://offa.example.com
  trust_anchors:
    - entity_id: https://ta.example.com
  authority_hints:
    - https://ta.example.com
  logo_uri: https://offa.example.com/static/img/offa-text.svg
  key_storage: /data
  use_resolve_endpoint: true
  use_entity_collection_endpoint: true
```

For more information about the offa config file, please refer to [OFFA Configuration](#offa-configuration).


### Apache

Still needs testing

## Docker Images

Docker images are available at [docker hub oidfed/offa](https://hub.docker.com/r/oidfed/offa/tags).

## Implementation State

This is currently a Proof of Concept, that still needs some improvements and tweaking.

The following is a list of TODOs:

- Query userinfo endpoint for user information in additional to id token.
- Show Home page with user information
    - Also use this as the default redirect target if no `next` is given
- Verify usage with:
    - Traefik
    - NGINX
    - Apache
- Other things will probably be added with further testing
- Extended documentation

---
You do not have to use OFFA, it's just an offer.

<img alt="logo" height="100" src="logos/offa.svg"/>
