/*
  UI Element HTML hooks
*/
const ui_elements = new Map()
ui_elements.set("alerts", '#emrs_app_alerts')
ui_elements.set("content", '#emrs_app_content')
ui_elements.set("user_display", '#emrs_user_name')
ui_elements.set("version_display", '#emrs_version_indicator')
ui_elements.set("dashboard_input", '#emrs_dashboard_input')


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

function editItem(item) {
  console.log("Edit item:", item)
  alertNotYetDone()
}

function deleteItem(item) {
  console.log("Delete item:", item)
  alertNotYetDone()
}

function dashboardListAssets() {
  app.getPageDashboard().changeViews(
    DashboardViews.ASSETS)
}

function dashboardListActions() {
  app.getPageDashboard().changeViews(
    DashboardViews.ACTIONS)
}

function dashboardListSignals() {
  app.getPageDashboard().changeViews(
    DashboardViews.SIGNALS)
}

function dashboardAddItem() {
  app.getPageDashboard().createItem()
}

function alertNotYetDone () {
  app.auth()
  app.alerts.error("not yet implemented")
}

function performLogout() {
  app.quit()
  delete app
}


