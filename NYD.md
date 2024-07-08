

Stubbed out are multiple types of signals.
 - onEvent                  <---- Current one we use for now
 - onTimeout
 - onBumpTimeout                    The rest will be available later and triggered/scheduled by core
 - onSchedule

TODO:

 - in action.js when we make the action, we need to also send the signal along with the request
    that way the system can automatically add the sigmap for the action

 - setup action running /scheduling. stop runners on update, redeploy, and show the status on timer in UI

 - terminal page. select active runnier, show terminal that dumps logs and setup "dev" env where the script can be fired without consequence in testing (figure out how)

