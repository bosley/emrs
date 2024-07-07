/*
  The various pages that the user can load in the UI
*/
const ApplicationPage = Object.freeze({
  NONE: 0,      // For uninit case
  DASHBOARD: 1,
  ACTIONS: 2,
  TERMINAL: 3,
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
    this.pages.set(ApplicationPage.DASHBOARD, new PageDashboard(this.alerts, this.getTopo))
    this.pages.set(ApplicationPage.ACTIONS, new PageAction(this.alerts, this.getTopo, this.getActionFiles))
    this.pages.set(ApplicationPage.TERMINAL, new PageTerminal(this.alerts, this.getTopo))

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

  getPageActions() {
    return this.pages.get(
      ApplicationPage.ACTIONS)
  }

  getPageTerminal() {
    return this.pages.get(
      ApplicationPage.TERMINAL)
  }

  getTopo() {
    $.ajax({
      type: "GET",
      url: "api/topo" + getApiKeyUrlParam(),
      dataType: 'json',
      async: false,
      error: ((function(obj){
        return function(){ 
          console.log("failed to retrieve emrs data")
          obj.alerts.error("Failed to retrieve EMRS data")
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          obj.topo = JSON.parse(data.topo)
          console.log("topo representation updated")
          console.log(obj.topo)
        }
      })(this))
    })

    return this.topo
  }

  getActionFiles() {
    $.ajax({
      type: "GET",
      url: "api/actions" + getApiKeyUrlParam(),
      dataType: 'json',
      async: false,
      error: ((function(obj){
        return function(){ 
          console.log("failed to retrieve emrs actions")
          obj.alerts.error("Failed to retrieve EMRS actions")
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          console.log(data)
          obj.actions_list = JSON.parse(data.files)
          console.log("actions list updated")
          console.log(obj.actions_list)
        }
      })(this))
    })

    return this.actions_list
  }
}

