# spinclass

Demo app to spin up a specified number of OpenStack Compute instances, monitor and time their build status

## Install

1. Download a binary from the most recent [release](https://github.com/sivel/spinclass/releases)

## Example Execution

```
./spinclass -username my-api-user -apikey my-api-key -region IAD -image 753a7703-4960-488b-aab4-a3cdd4b276dc -flavor performance1-1
```

## Configuration File

The optional configuration file should be placed in the same location as the `spinclass` binary and named `spinclass.yaml`

Full configuration example:

```yaml
server:
    port: ":8000"
    cert: /path/to/ssl.crt
    key: /path/to/ssl.key
openstack:
    identity: "https://identity.api.rackspacecloud.com/v2.0"
    username: my-api-user
    apikey: my-api-key
    password: password-if-not-using-api-key
    region: IAD
    image: 753a7703-4960-488b-aab4-a3cdd4b276dc
    flavor: performance1-1
```
