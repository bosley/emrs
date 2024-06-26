
/*
  Note:
      All "Page" Objects have to have the "render" method
*/
class PageTerminal {
  constructor(alerts) {
    this.alerts = alerts
  }

  setSelected() {
    console.log("terminal set to selected")
  }

  setIdle() {
    console.log("terminal set to idle")
  }

  render(contentTag) {
    console.log("Need to use the given content tag to draw terminal data: " + contentTag)
  }
}
