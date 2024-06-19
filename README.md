# emrs

Environmental Monitoring and Response System

This thing is under heavy development and isn't really doing anything at the moment.

__Check back later!__


# Getting Started

Start by ensuring you have docker installed on your system.

Once that is done, get the server up by running:

```
    ./dev.sh rel build
    ./dev.sh rel run
```

Thats it! The server is now local on port 8080.

## Development

For a more development-friendly version of the server change `rel` to `dev`
in the above command as so:

```
    ./dev.sh dev build
    ./dev.sh dev run
```

Thats it! The server is now local on port 8080.

### Note:

On MacOS you may need to do the following:

```
        EDIT:   ~/.docker/config.json

      REMOVE:   "credsStore": "desktop"
```
see: [here](https://serverfault.com/questions/1130018/how-to-fix-error-internal-load-metadata-for-docker-io-error-while-using-dock) for more details.

