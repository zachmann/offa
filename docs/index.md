# OFFA - Openid Federation Forward Auth

<img alt="logo" height="200" src="assets/offa-text.svg"/>

OFFA offers easy to use OpenID Federation Authentication and Authorisation for existing services.
OFFA can be deployed along existing services and handle all OpenID Federation communication for your services.

OFFA implements Forward Authentication usable with
[Traefik](https://doc.traefik.io/traefik/middlewares/http/forwardauth/),
[NGINX](https://docs.nginx.com/nginx/admin-guide/security-controls/configuring-subrequest-authentication/),
[Caddy](https://caddyserver.com/docs/caddyfile/directives/forward_auth),
and maybe other reverse proxies.

OFFA also implements [Auth MemCookie](https://zenprojects.github.io/Apache-Authmemcookie-Module/) usable with Apache.

See [Proxies](proxies/index.md) on how to deploy offa with your favorite reverse proxy.