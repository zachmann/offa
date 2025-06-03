---
title: debug_auth
icon: fontawesome/solid/person-digging
---

<span class="badge badge-purple" title="Value Type">boolean</span>
<span class="badge badge-blue" title="Default Value">`false`</span>
<span class="badge badge-green" title="If this option is required or optional">optional</span>

The `debug_auth` config option is used to enable additional output at the
`auth` endpoint. If set to `true` OFFA prints additional output for each
request to the `auth` endpoint that includes all received HTTP headers and
their values.
This can be useful to debug whether the reverse proxy sends the necessary
headers.

??? file "config.yaml"

    ```yaml
    debug_auth: true
    ```
