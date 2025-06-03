# Config
OFFA is configured through a single configuration file named `config.yaml`.

## Config File Location

OFFA will search for this file at startup at different locations, the first 
file that is found will be used. Supported locations are:

- `config.yaml`
- `config/config.yaml`
- `/config/config.yaml`
- `/offa/config/config.yaml`
- `/offa/config.yaml`
- `/data/config/config.yaml`
- `/data/config.yaml`
- `/etc/offa/config.yaml`

## Small Example Config File
The following is a small example config file:

??? file "config.yaml"

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