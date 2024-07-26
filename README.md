# Environmental Monitoring and Response System

To be brief I've listed the quick start here so you can jump right in. 

Following the quick start is a more in-depth look at emrs.

## Quick start 


1. clone the repo
2. run make (may be packages need to be gotten - not tested from blank install)
3. set `EMRS_HOME` environment variable to point to an empty or non-existing directory
4. run `./bin/emrs server --new` and follow the prompts
5. run the following commands to install the example actions:
```
    ./bin/emrs action --new "_actions/alert.go" --name "alert"
    ./bin/emrs action --new "_actions/logger.go" --name "logger"
```

Now, if you `tree $EMRS_HOME`, you should see something similar to:

```
/Users/bosley/.emrs-server
├── actions
│   ├── alert.go
│   └── logger.go
├── server.cfg
└── storage
    └── datastore.db

3 directories, 4 files
```

Any go file that you want to use to handle requests from the server can be installed with this
method (more below.)


6. create a few assets that can be used to identify something sending data to the server

```
./bin/emrs asset --new "asset-0"
./bin/emrs asset --new "asset-1"
./bin/emrs asset --new "asset-2"
```

Show the generated assets:

```
 ./bin/emrs asset --list
```

Which will yield something similar (ids will vary) to:

```
     0 | cf070dbe-a24c-8b4a-ac57-023a98e62c73 | asset-0
     1 | eecec5a4-858d-e1b1-67ac-93a8fa205611 | asset-1
     2 | 56821c8e-3a5d-29f0-3ada-eb325443e387 | asset-2
```

7. create tokens that can be used to validate submissions

```
./bin/emrs tokens --count 3 -duration "24h"
```

This will create a json list of tokens (basically JWTs) that can be submitted along
with the data to validate the submission and permit the request to be executed.
These tokens, or sometimes mentioned as "vouchers" will be valid for 24 hours as-per
the command above.

8. Send requests to the server:

Using emrs/api, the `HttpSubmissions` function can be used to get the `SubmissionApi`,
or you can submit an http request as follows:

```
  HTTPS POST to <URL>/submit/event 

  Header:
    Content-Type: octet-stream
    EMRS-API-Version: <VERSION>           [API Version (not yet utilized)]
    origin: <known asset UUID>            [UUID of reporting asset - must be known to EMRS]
    route: <emrs url proc path>           [example:   log.Log ] (more below)
    token: <authentication token>         [valid emrs token (from 7.)]
  Body:
    optional: binary data stream

```

9. Read more below to fill in the knowledge gaps left by the get-up-and-go list

## EMRS - About 

If `EMRS_HOME` is not set in the environment `--home` should be specified whenever executing
emrs

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

## Installation

Set `$EMRS_HOME` to an empty or non-existing directory in the environment.

Run `make` and then the following:

```
    ./bin/emrs server --new
```

Follow the prompts.

EMRS is now installed.

## Configuration

In `EMRS_HOME/server.cfg` there is a `key` and `cert` field. If these contain
paths to a valid key and cert then HTTPS will be enabled.

Other than that the config file be mostly untouched by hand.

## Startup

```
   ./bin/emrs server
```

Use `--release` to enable `release` mode.

## Asset Management

### List assets

```
    ./bin/emrs asset --list
```

### Add asset

```
    ./bin/emrs asset --new "asset-0"
```

### Remove asset

```
    ./bin/emrs asset --remove eecec5a4-858d-e1b1-67ac-93a8fa205611
                              ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ Assset UUID
```

### Update asset

As of now modifying an asset only entails renaming it, as its UUID is fixed

`--name` must be specified or the name will be erased (set to the nothing given to it)

```
    ./bin/emrs asset --update "56821c8e-3a5d-29f0-3ada-eb325443e387" --name "orangie"
```

## Command and Control

### Shutdown

EMRS Can be shutdown remotely provided that the server identity and user key is in the `--home` path
and specified. The following show shutting down a running EMRS instance from the CLI

```
    ./bin/emrs cnc --down
```

When running from the cli, you will be prompted for the password you set during installation.

This is just a demo method used to build-out the cnc api auth, and is likely to be removed.

## Event Submissions

Submissions to the server at the moment only take the form of "events." These "events"
are instances where something happened somewhere and now something needs to happen with
data from the event and digested into the operational context that EMRS finds itself in.

### EMRS URL

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

The first section of the path indicates which action to run. This action must be known
by the server. The next section is the function within the action that will handle the
data from the submission.

The remaining sections are free-use and contain no unreasonable upper-limit.

### Submitting Data

The data submission is a POST request to the server to `/submit/event`
with the following data:

```
  HTTPS POST format for Submission and CNC endpoints
  
  Header:
    Content-Type: octet-stream
    EMRS-API-Version: <VERSION>           [API Version (not yet utilized)]
    origin: <known asset UUID>            [UUID of reporting asset - must be known to EMRS]
    route <emrs url proc path>            [log.Log]
    token <authentication token>          [emrs token (see 'emrs tokens --help']
  Body:
    optional: binary data stream          [raw data to submit]
```

Example submission (assuming logger action is installed):

```
./bin/emrs submit -to cf070dbe-a24c-8b4a-ac57-023a98e62c73:logger.Log.example.sections@http://localhost:8080 --data test
```

The server output should be similar to:

```
logger: 2024-07-25 19:35:01.64112 -0400 EDT m=+172.858977959

    origin: cf070dbe-a24c-8b4a-ac57-023a98e62c73
     route: example.sections
      data: test
```

## Handling Data

Actions, once installed, can be used after a server restart. 

The `name` given to the action on install _must_ match the `package` that it declares in its source.

Functions that are expected to be triggered via submissions should be public and take the signature:

```
func MySpecialHandler(origin string, route []string, data []byte) error
```

The origin will be the UUID reported on submission, the route will be any remaining sections
on the route (that is, it excludes the action name and function sections.)
The data will be the raw data submitted to the endpoint to do with as you please.

### EMRS Runtime

***WARNING*** This is the next active part of development, and is not yet usable at all.

In an installed action, the package "emrs" is made available for import. (TODO: make a testing package do logic can be tested outside of running server)

This package contains functionality that permits interaction with the emrs server to
manage specific states or to execute some desired function.

Example:

```
package mycoolaction

import(
    "emrs"
)

func something() {
    emrs.Log("whoop-there-it-is")
}

```

The specific functions available to the system are currently:

Log, Emit, and Signal" though they are not fully implemented (aside log)

## Next Steps

Once the emrs runtime is to a point where the software is functional and at-least potentially-usefull, a GUI is going
to be developed that interacts with the EMRS server via the established apis found in /api:

- CNC
- Submissions
- Stats
