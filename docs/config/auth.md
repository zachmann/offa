---
icon: material/security
---
<span class="badge badge-purple" title="Value Type">list of Auth Rules</span>
<span class="badge badge-orange" title="If this option is required or optional">required, unless apache is used</span>

Under the `auth` option Auth Rules are configured. Each Auth Rule controls 
access to a resource / service / URI.

When OFFA receives a request to the
[forward auth endpoint](server.md#forward_auth), the reverse proxy includes 
information in the request on which site the user wants to access. OFFA uses 
the configured Auth Rules to find a matching rule and evaluate if the user 
is authorised to do so or not.
(If no rule matches, the user is not authorised).

The following options are available for each Auth Rule:

## `domain`
<span class="badge badge-purple" title="Value Type">string</span>
<span class="badge badge-red" title="If this option is required or optional">required, unless `domain_regex` is given</span>

The `domain` option is used to set the domain that is used to match a 
request with the Auth Rule.

??? file "config.yaml"

    ```yaml
    auth:
        - domain: foobar.example.com
    ```

## `domain_regex`
<span class="badge badge-purple" title="Value Type">string</span>
<span class="badge badge-red" title="If this option is required or optional">required, unless `domain` is given</span>

The `domain_regex` option is used to set a regular expression to match the 
request's domain with the Auth Rule.

The regex must be in the `Golang` flavor. We recommend https://regex101.com/ 
to try out your regexes.

??? file "config.yaml"

    ```yaml
    auth:
        - domain_regex: '^(pub|img)-data\.example\.com$'
    ```

## `path`
<span class="badge badge-purple" title="Value Type">string</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `path` option is used to set an url path that is used to match a
request with the Auth Rule.
If not set, any path will match.

!!! warning

    Using the `path` option requires an exact match. Sub-paths are not 
    matched. To do so, [`path_regex`](#path-regex) must be used.

    !!! example

        `/private` does only match `/private`, but not `/private/` or 
        `/private/foo`.

??? file "config.yaml"

    ```yaml
    auth:
        - domain: foobar.example.com
          path: "/private"
    ```

## `path_regex`
<span class="badge badge-purple" title="Value Type">string</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `path_regex` option is used to set a regular expression to match the
request's url path with the Auth Rule.

The regex must be in the `Golang` flavor. We recommend https://regex101.com/
to try out your regexes.

??? file "config.yaml"

    ```yaml
    auth:
        - domain: foobar.example.com
          path_regex: '^/private(/?].*)?$'
    ```

## `require`
<span class="badge badge-purple" title="Value Type">list</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `require` option is used to define authorisation requirements, i.e. 
requirements a user has to fulfill in order to get access.

The syntax for this option is a bit more complex to allow for flexibility.
Generally, the value for `require` is a list. Each entry is a 'require option',
i.e. as soon as one of the options matches the user, the user will be 
authorised -- to say it differently: those are logically ORed.

If only one option is given the list can be skipped, and it can be given as a 
single option.

If no options are given, every authenticated user will be authorised.

## `require` options

Each option entry is a mapping where the keys are OIDC Claim names and the 
value is a list of values.
The option matches for a user when the user fulfills all specified claims. 
Only claims where the value is a string or array of strings can be used.
If the claim value is a string, the specified claim value must be equal to 
the user claim value in order to fulfill the claim. 
If the claim value is an array of strings, the  specified claim values must 
be a subset of the user claim values in order to fulfill the claim.

If only a single claim value is required (nevertheless if the claim value 
type is string or array), it can be specified as a single string (skipping 
the list).

The following `require` options are all valid and all equivalent and all 
require that users are in the `admin` group:

=== ":material-file-code: `config.yaml`"

    ```yaml
    auth:
      - domain: foobar.example.com
        require:
           - groups:
                - admin
    ```

=== ":material-file-code: `config.yaml`"

    ```yaml
    auth:
      - domain: foobar.example.com
        require:
            groups:
                - admin
    ```

=== ":material-file-code: `config.yaml`"

    ```yaml
    auth:
      - domain: foobar.example.com
        require:
           - groups: admin
    ```

=== ":material-file-code: `config.yaml`"

    ```yaml
    auth:
      - domain: foobar.example.com
        require:
            groups: admin
    ```

The following is a more complex example with four different require options:

!!! file "config.yaml"

    ```yaml
    auth:
      - domain: foobar.example.com
        require:
            - groups: admin
            - groups:
                - dev
                - foobar
            - sub: john
            - groups: dev
              foo: bar  
    ```

## `forward_headers`
<span class="badge badge-purple" title="Value Type">mapping / object</span>
<span class="badge badge-blue" title="Default Value">see file example</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `forward_headers` option is used to specify HTTP Headers that should be 
forwarded to the client / site and also from which OIDC claims the information should
be obtained. The `forward_headers` is a mapping where keys are http header 
names, and the value are oidc claims.

!!! info

    OIDC Claims can be given as a single string or a list of strings. If a 
    list is given OFFA will use the value from the first non-empty claim.

    !!! example

        In the config below `X-Forwarded-User` will be populated with the 
        value in `preferred_username` if that is set, or `sub` otherwise.

The default mapping is as listed in the following `config.yaml` example.

!!! file "config.yaml"

    ```yaml
    auth:
      - domain: foobar.example.com
        forward_headers:
            X-Forwarded-User:
                - preferred_username
                - sub
            X-Forwarded-Groups:
                - entitlements
                - groups
            X-Forwarded-Email: email
            X-Forwarded-Name: name
            X-Forwarded-Provider: iss
            X-Forwarded-Subject: sub
    ```

## `redirect_status`
<span class="badge badge-purple" title="Value Type">integer</span>
<span class="badge badge-blue" title="Default Value">303</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

If the user needs to be authenticated, OFFA will redirect the browser to the 
[login page](server.md#login). By default, this uses the HTTP status code 
`303` See Other. However, this might not be supported by the reverse proxy, 
(e.g. nginx only accepts `401` and `403` responses to authentication 
subrequests). 
The `redirect_status` option is used to change the status code for such 
redirects.

??? file "config.yaml"

    ```yaml
    auth:
      - domain: foobar.example.com
        redirect_status: 401
    ```
