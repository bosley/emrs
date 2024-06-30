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

    this.session = new AppSession()

    this.page = ApplicationPage.NONE

    this.auth()

    $(ui_elements.get("user_display")).html(
      this.session.user)

    console.log(this.session.user)

    $(ui_elements.get("version_display")).html(
      this.session.version)

    this.contentHook = ui_elements.get("content")
    this.alerts = new Alerts(ui_elements.get("alerts"))

    this.pages = new Map()
    this.pages.set(ApplicationPage.DASHBOARD, new PageDashboard(this.alerts))
    this.pages.set(ApplicationPage.TERMINAL, new PageTerminal(this.alerts))

    this.content_hook = ui_elements.get("content")

    this.loadPage(ApplicationPage.DASHBOARD)
  }

  // Ensure that the user is still authorized, or redirect to main EMRS page
  auth() {
    this.session.validate()
    if (this.session.isValid()) {
      return
    }
    this.unauthorized()
  }

  // Shutdown the app and logout the user
  quit() {
    this.alerts.warning("logging out")
    if (this.session.isValid()) {
      this.session.quit()
      location.href = "/logout"
    } else {
      location.href = "/"
    }
  }

  // Redirect an unauthorized session to main EMRS page
  unauthorized() {
    console.log("unauthorized session detected")
    location.href = "/"
  }

  

  // Called by UI Hooks set in main.js
  // Checks the user's session to ensure that we are still
  // using an authenticated session
  loadPage(pageId) {

    if (this.page === pageId) {
      return
    }

    this.forceLoadPage(pageId)
  }

  forceLoadPage(pageId) {

    this.auth()

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

