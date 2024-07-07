

Stubbed out are multiple types of signals.
 - onEvent
 - onTimeout
 - onBumpTimeout
 - onSchedule

Only the core emits signals around. The signals should be done so 
via internal "assets" that are most likely yaegi "workers" or "actions"
that act as assets internally
