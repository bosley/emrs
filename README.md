# Environmental Monitoring and Response System EMRS
_(suggested pronunciation: immerse)_



This system aims to provide a configurable management for 
monitoring and responding to asset-generated events. EMRS will do
this by providing an easy-to-use watchdog scripting system that integrates a user's
configured assets, a persistent event-environment, and a custom-job execution api.


The idea is:

The user will be able to use a scripting language to create watchdogs that
monitor events from one-or-more assets as they submit events to the system.
These watchdogs can modify environment variables that persist over the lifetime
of the server. Using triggers, or built-in functions within the environment, the
watchdogs can accumulate data for a section of time, and queue it to be received
by yet-another-script that the user provides and configures as an asset that
receives events. If a trigger goes off, it will call a function that attempts to
send an event to an asset. With user-provided scripts as assets they can chain
behavior to process event data doing whatever they want. 

One idea would be to use this as a greenhouse monitor. Events could be published
to the server with watchdogs keeping track of temperature. If a watchdog detects
temp is at a boundary, then it can execute a correcting action by sending an
event to a fan or heat controller. 

Similarly, scheduled asset eventing could be created. This would permit the server
to send out automated ping messsages, but also send wake-up events to
watchdog scripts that may need to monitor last-contact with another asset
(temperature sensor, or remote server) and respond by asking them to phone-home,
or maybe by sending a text message to the user.

Think "IFTTT but for someone who wants a self-hosted offline-possible monitoring
and response system that doesn't want my data"
