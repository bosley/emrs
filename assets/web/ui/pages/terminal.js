class PageTerminal {
  constructor(alerts) {
    this.alerts = alerts
    this.selected = false
  }

  setIdle() {
    this.selected = false
  }

  setSelected(contentHook) {
    this.selected = true
    $(contentHook).html("term")
  }
}
