/*
  UI Element HTML hooks
*/
const ui_elements = new Map()
ui_elements.set("alerts", '#emrs_app_alerts')
ui_elements.set("content", '#emrs_app_content')
ui_elements.set("user_display", '#emrs_user_dropdown')
ui_elements.set("version_display", '#emrs_version_indicator')


const app = new Application(ui_elements);

/*
  UI Input Hooks
*/
function loadDashboard() {
  app.loadPage(ApplicationPage.DASHBOARD)
}

function loadTerminal() {
  app.loadPage(ApplicationPage.TERMINAL)
}

function alertNotYetDone () {

  // This isn't required, but technically it should happen
  app.auth()
  app.alerts.error("not yet implemented")
}

function performLogout() {
  app.quit()
  delete app
}
