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

    this.session_valid = false;
    this.page = ApplicationPage.NONE

    this.auth()

    this.contentHook = ui_elements.get("content")
    this.alerts = new Alerts(ui_elements.get("alerts"))

    this.pages = new Map()
    this.pages.set(ApplicationPage.DASHBOARD, new PageDashboard(this.alerts))
    this.pages.set(ApplicationPage.TERMINAL, new PageTerminal(this.alerts))

    this.loadPage(ApplicationPage.DASHBOARD)
  }

  validateSession() {
    // We can determine if the user's session is still valid by attempting
    // to access a protected region of the server /emrs. If we get anything
    // other than success, the session is no longer valid
    $.ajax({
      type: "GET",
      url: "/emrs",
      async: false,
      error: ((function(obj){
        return function(){ 
          obj.session_valid = false;
          console.log("session no longer valid")
        }
      })(this)),
      success: ((function(obj){
        return function(){
          obj.session_valid = true;
        }
      })(this))
    })
  }

  // Ensure that the user is still authorized, or redirect to main EMRS page
  auth() {
    this.validateSession()
    if (this.session_valid) {
      return
    }
    this.unauthorized()
  }

  // Shutdown the app and logout the user
  quit() {
    this.alerts.warning("logging out")

    if (this.session_valid) {
      this.session_valid = false
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

    this.auth()

    if (this.page === pageId) {
      return
    }

    if (this.page != ApplicationPage.NONE) {
      this.pages.get(this.page).setIdle()
    }

    this.page = pageId

    this.pages.get(this.page).setSelected(this.contentHook)
  }
}

