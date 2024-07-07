/*
  The various pages that the user can load in the UI
*/
const ApplicationPage = Object.freeze({
  NONE: 0,      // For uninit case
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

    $(ui_elements.get("user_display")).html(
      "[USER]")

    $(ui_elements.get("version_display")).html(
      "0.0.0.0 (TBD)")

    this.contentHook = ui_elements.get("content")
    this.alerts = new Alerts(ui_elements.get("alerts"))

    this.pages = new Map()
    this.pages.set(ApplicationPage.DASHBOARD, new PageDashboard(this.alerts))
    this.pages.set(ApplicationPage.TERMINAL, new PageTerminal(this.alerts))

    this.content_hook = ui_elements.get("content")

    this.loadPage(ApplicationPage.DASHBOARD)
  }

  loadPage(pageId) {

    if (this.page === pageId) {
      return
    }

    if (this.page != ApplicationPage.NONE) {
      this.pages.get(this.page).setIdle()
    }

    this.page = pageId

    $(this.content_hook).html("<p>.")

    this.pages.get(this.page).setSelected(this.contentHook)
  }

  getPageDashboard() {
    return this.pages.get(
      ApplicationPage.DASHBOARD)
  }

  getPageTerminal() {
    return this.pages.get(
      ApplicationPage.TERMINAL)
  }


}

