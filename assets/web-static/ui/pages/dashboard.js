/*
  Note:
      All "Page" Objects have to have the "render" method
*/
class PageDashboard {
  constructor(alerts) {
    this.alerts = alerts
  }

  setSelected() {
    console.log("dashboard set to selected")
  }

  setIdle() {
    console.log("dashboard set to idle")
  }

  render(contentTag) {
    console.log("Need to use the given content tag to draw dashboard data: " + contentTag)
  }
}
