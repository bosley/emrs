const ApiOp = Object.freeze({
  ADD: "opAdd",
  DEL: "opDel",
})

const ApiSubject = Object.freeze({
  SECTOR: "sector",
  ASSET : "asset",
  SIGNAL: "signal",
  ACTION: "actopm",
  MAPPING : "mapping",
  TOPO: "topo",
})

const SignalTrigger = Object.freeze({
	EVENT         : "onEvent",
	TIMEOUT       : "onTimeout",
	BUMPTIMEOUT   : "onBumpTimeout",
	SHUTDOWN      : "onShutdown",
	SCHEDULE      : "onSchedule",
	EMIT          : "onEmit",
})

const ActionType = Object.freeze({
  FILE: "file",
  EMBEDDED: "embedded",
})

class ApiMsg {
  constructor(op, subject, data) {
    this.op = op
    this.subject = subject
    this.data = data
  }
}

class EmrsHeader {
  constructor(name, description, tags) {
    this.name = name
    this.description = description
    this.tags = tags
  }
}

class EmrsSector {
  constructor(header, assets) {
    this.header = header
    this.assets = assets
  }
}

class EmrsAsset {
  constructor(header) {
    this.header = header
  }
}

class EmrsSignal {
  constructor(header, trigger) {
    this.header = header
    this.trigger = trigger
  }
}

class EmrsAction {
  constructor(header, type, data) {
    this.header = header
    this.type = type
    this.data = data
  }
}
