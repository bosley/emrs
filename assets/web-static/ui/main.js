/*
  UI Element HTML hooks
*/
const ui_elements = new Map()
ui_elements.set("alerts", '#emrs_app_alerts')
ui_elements.set("content", '#emrs_app_content')

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
  app.alerts.error("not yet implemented")
}


