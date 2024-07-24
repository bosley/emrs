# Environmental Monitoring and Response System

## Creating EMRS instance:

If `EMRS_HOME` is not set in the environment `--home` should be specified whenever executing
emrs

```
    go run cmd/cli/*.go --new
```

Helpful flags

```

--force 
        (used in conjunction with `--new` to overwrite existing EMRS instances

--stat [https://HOST:PORT]
        Get uptime of server

--vouchers [n]
        Generate EMRS equivalent of JWTs (vouchers.) Give --duration to set the 
        time that the voucher(s) should be food for

```

Within the configuration file generated for the EMRS instance in the specified `home`
directory there exists a "key" and "cert" entry. If these contain a path to a valid
key and cert, EMRS will attempt to run with HTTPS.

Remember that the URL/ port combination used in the key generation must reflect
that which the binding indicates within the config (localhost vs 127.0.0.1 vs example.server.com)
lest ye invoke net/http errors.

Note: 

As a general rule if something is listed as a "binding" its just a PORT:ADDRESS pair, and if
something is refenced as a URL it will require a prefixed "http" or "https." Usually where
this comes into play EMRS attemps to auto generate URLs from the given bindings, but it is
possible that an error occurs in some instances as this thing is not even off the ground yet.
Just make sure that when you type in a EMRS URL (below) that you include the correct schema
prefix (http/https) and don't omit them.

## Asset Management

### List assets

```
    go run cmd/cli/*.go --list-assets
```

### Add asset

```
    go run cmd/cli/*.go --new-asset "some-name"
```

### Remove asset

```
    go run cmd/cli/*.go --remove-asset 55adedda-2ae0-61dc-5bc6-44e29f7468d8
                                       ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ Assset UUID
```

### Update asset

As of now modifying an asset only entails renaming it, as its UUID is fixed

`--name` must be specified or the name will be erased (set to the nothing given to it)

```
 go run cmd/cli/*.go --update-asset 70b29cd0-5e19-fa07-e18e-b309856d2ab5 --name "orangie"
```

## Command and Control

### Shutdown

EMRS Can be shutdown remotely provided that the server identity and user key is in the `--home` path
and specified. The following show shutting down a running EMRS instance from the CLI

```
 go run cmd/cli/*.go --down

```

## Event Submissions

Submissions to the server at the moment only take the form of "events." These "events"
are instances where something happened somewhere and now something needs to happen with
data from the event and digested into the operational context that EMRS finds itself in.

**EMRS URL**

The URL associated with events are formatted as follows:

```
   01808f90-097b-9c5a-7427-077ebc2254af:example.action.path@https://localhost:8080

   |__________________________________| |_________________| |_____| |_______| |___|

        EMRS Asset UUID                   Processing Path    Schema  Address   Port
```

The EMRS Asset UUID is auto-generated, but the Processing Path is the pipeline
that the data will be sent through. Each segment within the path, seperated by `.`
are referred to internally as "chunk" but may referred to as "section." The UUID and
path chunks all follow the same general fomatting rule for validation; they must
be alphanumeric, with the allowance of `_` and `-`. Any other symbols contained in
a processing path or UUID will be flagged as invalid.

The chunks within the processing path as-of right now are arbitrary and the entire path
is mapped to a configurable action.

In the future, the processing path will determine how the data is pipelined through the system

Example submission:

```
 go run cmd/cli/*.go --submit [EMRS URL] 
```

The above submission will send an empty data field 


## Disk Layout

```
EMRS_HOME
├── actions
│   └── init.go                 <----   Where EMRS Url processing path is computerd and request delegated
├── server.cfg                  <----   Server identity, https key setup, etc
└── storage
    └── datastore.db            <----   EMRS Database

```
