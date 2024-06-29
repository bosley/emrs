

FINISH USER PART OF DATABASE


FIGURE OUT HOW TO MODEL ASSETS AND HOW THEIR EVENTS
ACAN BE "GROUPED" FOR CONSUMPTION BY WATCHDOGS. 

DETERMINE HOW TO INTEGRATE LUA AND ALL THAT. NERV WOULD
BE GOOD TO USE AROUND THIS TIME FOR EVENT KICKING OFF AND SUBMISSION
FROM HTTPS


MAYBE JUST WRITE DOWN ALL THE THINGS THAT NEED TO BE DONE TO GET TO MVP THEN JUST ORDER THEM
BY WHAT MAKES SENSE 




- Read the discord for info on the overall design we will go for

essentially:
    define submission endpoints with optional auth
    define "pipelines" that utilize a "decoder" function written in interpreted go that yeilds a json object,
        these pipelines will take-in the json object. Pipelines are essentially nerv topics
    watchdog scripts written in golang (interpreted) will be run in nerv consumers and execute every n-events.
    When "n" is hit, the event(s) are given to the watchdog for execution. Watchdog settings should be made available
    to permit/deny async execution. They may want to block while processing and buffer the next events until its complete.
    
These watchdog's operational environments should permit them to store a state that persists between executions.
We may also want to have these environments be made available to others. 

This hints at each asset having its own internal persistent environment, and have a means to
communicate with whatever sector its in. Calls to a sector obejct (sentinal?) could have events
posted to other assets in the sector, or to a parent sector, to kick off alerts or store metrics.


```

    home (S)
        |                     /-- Light (rx)
        | -------- Garage (S) ------- Sensor Cluster (tx)
        |
        | -------- Back yard (S)
        |            |
        |            | -------------- Wood Shed (S) ----- Light (rx)
        |            |                      | ----------- Sensor Cluster (tx)
        |            |                      | ----------- Alarm (rx)
        |            | ------ Sauna (S)
        |            |          | ------ Light (rx)
        |            |
        |            | ------ Camera (tx/rx)
        |
        | -------- Garden (S)
        |            | -------- Strawberries (S)
        |            |              | ------------ Sensor Cluster (tx)
        |            |
        |            | -------- Carrots (S)
        |            |              | ------------ Sensor Cluster (tx)
        |            |
        |            | ----- Environmental Control (Temperature/ Humidity/ Fan) Cluster (tx/rx)
```


In the above example setup, we can issue messages as thus:

```
    tx /home/Back yard/Wood Shed/Alarm activate 15.0
```
sending the command "activate" with the param "15.0" assuming the device will take in number of seconds to run alarm


We could have each asset then be:

    username:/sector/../asset 

Each asset that can receive (rx) messages must have a file that will be utilized as the `tx` command in the
event that the above command was executed.

Sectors can optionally contain "sentinal" scripts that exists as an asset within the sector that acts with its
authority. Able to pass messages/ command other assets within the sector, and to flag parent sentinals about events.

NOTE: Events should have a "test" header flag that can be set to route to specific "test" functionality (LATER)



The server itself is a sort of sentinal and should have an "Avatar" service that runs the primary sentinal script for
the current user. Configuring behavior of the system, such as triggering text messages. emails, and actions to assets,
will be done with the user's avatar, theryby bypassing clunky configs and logical discontinuity between what the user uses
and what the software is designed around.



ideaing
```

home:
    garage:
        light:
            - asset








```




## Assets

Need to figure out how we want to model general "assets" that can be later refined by
layers/ tags to describe the thing. This needs to be done before datastore AssetStore can
be completed.

Assets are the implementation, or instantiation of, a possibly templated set of information
that describes:
 - physical/ virtual
 - none, rx/ tx/ rxtx
 - if rx/ rxtx, endpoints for submission listed, with name of template that outlines message structure

Each asset template should be stored in a table, and then an actual asset will link to that descrioption


## Datastore 

 - Server Store
    > Complete
 - User store
    > UpdatePassword
    > DeleteUser
 - Asset Store
    > 
 - Event store
