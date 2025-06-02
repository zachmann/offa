
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
    └── config.yaml
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
    image: oidfed/offa:main
    restart: unless-stopped
    volumes:
      - ./offa/config.yaml:/config.yaml:ro
      - ./offa:/data
    networks:
      caddy:

  # This would be your service
  whoami:
    image: containous/whoami
    restart: unless-stopped

networks:
  caddy:
```

#### caddy/Caddyfile

```
offa.example.com {
  reverse_proxy offa:15661
}

whoami.example.com {
    forward_auth offa:15661 {
		uri /auth
		copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
	}

    reverse_proxy whoami:80
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
  - domain: whoami.example.com
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

For more information about the offa config file, please refer to [OFFA Configuration](../config.md).
