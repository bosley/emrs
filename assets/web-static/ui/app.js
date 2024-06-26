/*
  The various pages that the user can load in the UI
*/
const ApplicationPage = Object.freeze({
  NONE: 0,
  DASHBOARD: 1,
  TERMINAL: 2
});

/*
  The controller application for the emrs ui application.
  UI Component interaction is tied to calls within the application
  object to change state. 
  The state will be reflected upon the next UI update 
*/
class Application {
  constructor(ui_elements) {
    this.page = ApplicationPage.NONE

    this.contentHook = ui_elements.get("content")
    this.alerts = new Alerts(ui_elements.get("alerts"))

    this.pages = new Map()
    this.pages.set(ApplicationPage.DASHBOARD, new PageDashboard(this.alerts))
    this.pages.set(ApplicationPage.TERMINAL, new PageTerminal(this.alerts))

    this.loadPage(ApplicationPage.DASHBOARD)
  }

  tick() {
    this.alerts.tick()
    this.pages.get(this.page).render(this.contentHook)
  }

  loadPage(pageId) {
    if (this.page === pageId) {
      return
    }

    if (this.page != ApplicationPage.NONE) {
      this.pages.get(this.page).setIdle()
    }

    this.page = pageId

    this.pages.get(this.page).setSelected()
  }
}


