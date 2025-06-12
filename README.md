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


## Documentation

For more information please refer to the Documentation at https://go-oidfed.github.io/offa/



## OFFA Configuration

Configuration of OFFA is explained in details at
https://go-oidfed.github.io/offa/config/.



## Docker Images

Docker images are available at [docker hub oidfed/offa](https://hub.docker.com/r/oidfed/offa/tags).

## Implementation State

This is currently a Proof of Concept, that still needs some improvements and tweaking.

The following is a list of TODOs:

- Query userinfo endpoint for user information in additional to id token.
- Show Home page with user information
    - Also use this as the default redirect target if no `next` is given
- Other things will probably be added with further testing

## Fun Facts about OFFA

- The default port `15661` represents the name `OFFA`
  - `O` is the 15th letter of the alphabet, `F` the sixth, `A` the first
- The elements in the logo have a meaning:
  - You might have noticed that OFFA sounds a lot like offer. The open hand 
    offers the feather.
  - The word `federation` contains the German word `Feder` which means 
    `feather`. Therefore, the feather.
  - Putting it together: A hand offering a feather.


---
You do not have to use OFFA, it's just an offer.

<img alt="logo" height="100" src="logos/offa.svg"/>
