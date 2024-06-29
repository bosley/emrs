const DashboardViews = Object.freeze({
  ASSETS:  "asset",
  ACTIONS: "action",
  SIGNALS: "signal",
})

class PageDashboard {
  constructor(alerts) {
    this.alerts = alerts
    this.selected = false
  }

  setIdle(contentHook) {
    this.selected = false
    $(contentHook).remove("#emrs-generated-dashboard")
  }

  setSelected(contentHook) {
    console.log("setSelected")
    $(contentHook).html('<div id="emrs-generated-dashboard"></div>')
    this.contentHook = contentHook
    this.setupDashboardPage()
    this.changeViews(DashboardViews.ASSETS)
    this.createSubmissionForm()
    this.selected = true
  }

  setupDashboardPage() {
    $("#emrs-generated-dashboard").html(`
    <div class="dashboard-display-control">
      <div class="row">
        <div class="column column-20">
            <button
              class="button button-black button-clear"
              type="submit"
              onclick="dashboardListAssets()">
               assets 
            </button>
        </div>
        <div class="column column-20 column-offset-20">
            <button
              class="button button-black button-clear"
              type="submit"
              onclick="dashboardListActions()">
               actions 
            </button>
        </div>
        <div class="column column-20 column-offset-20">
            <button
              class="button button-black button-clear"
              type="submit"
              onclick="dashboardListSignals()">
               signals
            </button>
        </div>
      </div>
    </div>
    <div style="height:2%;"></div>
    <div id="dashboard-view"></div>`)
  }

  clearTable() {
    $(this.contentHook).remove("#emrs-dashboard-table")
  }

  clearTableBody() {
    $(this.contentHook).remove("#emrs-dashboard-table-body")
  }

  makeTable(c1, c2, c3) {
    $("#dashboard-view").html(`
      <div id="emrs-dashboard-table">
        <table>
          <thead>
            <tr>
              <th>` + c1 + `</th>
              <th>` + c2 + `</th>
              <th>` + c3 + `</th>
            </tr>
          </thead>
          <tbody id="emrs-dashboard-table-body"></tbody>
      </div>
    `)
  }

  changeViews(view) {
    this.clearTable()

    let tableUrl = "/unused/invalid/noexist/null"

    switch (view) {
      case DashboardViews.ASSETS:
        this.makeTable("Asset", "Last Contact", "")
        tableUrl = "/app/get/assets?view=1"
        break;
      case DashboardViews.ACTIONS:
        this.makeTable("Action", "Status", "")
        tableUrl = "/app/get/actions?view=1"
        break;
      case DashboardViews.SIGNALS:
        this.makeTable("Signal", "In-Use", "")
        tableUrl = "/app/get/signals?view=1"
        break;
      default:
        this.alerts.error("Internal error: Invalid view name")
        return
    }

    this.view = view
    this.populateTable(tableUrl)
  }

  populateTable(endpoint) {

    // The provided table url should hand us back json
    // that is a list of 3 items, one for each column

    console.log("need to get json from ", endpoint)

    $("#emrs-dashboard-table-body").append(
      `<tr>
        <td>/home/garden/fan</td>
        <td>2 hours ago</td>
        <td>
          <button
            class="button edit-button"
            type="submit"
            onclick="editItem('UUID-002')">
              edit
          </button>
          <button
            class="button delete-button"
            type="submit"
            onclick="deleteItem('UUID-002')">
              delete
          </button>
        </td>
      </tr>`)

  }

  createSubmissionForm() {
    $("#emrs-generated-dashboard").append(
      `<form onsubmit="event.preventDefault(); dashboardAddItem();">
            <fieldset>
              <div class="row">
                <div class="column column-50">
                  <input type="text" style="width:100%;" required id="emrs_dashboard_input">
                </div>
                <div class="column column-10">
                  <button class="edit-button">ADD</button>
                </div>
              </div>
            </fieldset>
          </form>`)
  }
 
  createItem() {
    let item = $("#emrs_dashboard_input").val()
    $("#emrs_dashboard_input").val('')

    console.log("create", item, "for", this.view)
  }
}
