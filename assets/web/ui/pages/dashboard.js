const DashboardViews = Object.freeze({
  ASSETS:  "asset",
  ACTIONS: "action",
  SIGNALS: "signal",
})

class PostCreate {
  constructor(classification, name) {
    this.classification = classification
    this.name = name
  }
}

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
        break;
      case DashboardViews.ACTIONS:
        this.makeTable("Action", "Status", "")
        break;
      case DashboardViews.SIGNALS:
        this.makeTable("Signal", "In-Use", "")
        break;
      default:
        this.alerts.error("Internal error: Invalid view name")
        return
    }

    this.view = view
    this.populateTable()
  }

  reload() {
    this.changeViews(this.view)
  }

  populateTable() {

    $.ajax({
      type: "GET",
      url: "/app/dashboard",
      dataType: 'json',
      error: ((function(obj){
        return function(){ 
          console.log("session not valid")
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          function addItem(item) {
            $("#emrs-dashboard-table-body").append(
              `<tr>
                <td>` + item["Col1"] + `</td>
                <td>` + item["Col2"] + `</td>
                <td>
                  <button
                    class="button edit-button"
                    type="submit"
                    onclick="dashboardEditItem('` + item["Col1"] + `')">
                      edit
                  </button>
                  <button
                    class="button delete-button"
                    type="submit"
                    onclick="dashboardDeleteItem('` + item["Col1"] + `')">
                      delete
                  </button>
                </td>
              </tr>`)
          }
          let items = data[obj.view]
          for (let i = 0; i < items.length; i++) {
            addItem(items[i]) 
          }
        }
      })(this))
    })
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

    let msg = JSON.stringify(new PostCreate(this.view, item)) 

    console.log(msg)

    $.ajax({
      type: "POST",
      url: "/app/create",
      dataType: 'json',
      data: msg,
      contentType: 'application/json',
      error: ((function(obj){
        return function(){ 
          obj.alerts.error("error creating:" + item)
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          obj.alerts.info(data["status"])
          obj.reload()
        }
      })(this))
    })
  }

  editItem(item) {

    console.log("Dashboard::editItem(",item,")")
  }

  deleteItem(item) {

    console.log("delete", item, "for", this.view)

    let msg = JSON.stringify(new PostCreate(this.view, item)) 

    console.log(msg)
    $.ajax({
      type: "POST",
      url: "/app/delete",
      dataType: 'json',
      data: msg,
      contentType: 'application/json',
      error: ((function(obj){
        return function(){ 
          obj.alerts.error("error creating:" + item)
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          obj.alerts.info(data["status"])
          obj.reload()
        }
      })(this))
    })
  }
}
