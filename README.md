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
