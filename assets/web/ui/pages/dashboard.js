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
    this.selected = true
  }

  setupDashboardPage() {
    $("#emrs-generated-dashboard").html(`
    <div class="dashboard-display-control">
      <div class="row">
        <div class="column column-20">
            <button
              class="button button-black button-clear"
              id="emrs_dashboard_add_item_button"
              type="submit"
              onclick="dashboardAddItem()">
            </button>
        </div>
        <div class="column column-20 column-offset-20">
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

  changeViews(view) {
    $("#dashboard-view").html("")

    let buttonText = "[INVALID]"
    switch (view) {
      case DashboardViews.ASSETS:
        buttonText = "+ asset"
        this.makeTable("name", "last seen", "")
        break;
      case DashboardViews.ACTIONS:
        buttonText = "+ action"
        this.makeTable("name", "status", "")
        break;
      case DashboardViews.SIGNALS:
        buttonText = "+ signal"
        this.makeTable("name", "in-use", "")
        break;
      default:
        this.alerts.error("Internal error: Invalid view name")
        return
    }

    $("#emrs_dashboard_add_item_button").html(buttonText)

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

          obj.data_cache = data

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

  setupAssetInput(item, postUrl, use_active) {
    let body = 
    `<div class="row">
      <div class="column column-50">
        <form method="POST" action="` + postUrl + `">
          <fieldset>
            `

    if (use_active) {
      let description = "unknown"
      let items = this.data_cache[this.view]
      for (let i = 0; i < items.length; i++) {
console.log(items[i]["Col1"], " | " , item)
        if (items[i]["Col1"] === item) {
          description = items[i]["Col3"]
          i = items.length
          console.log("found description:", description)
        }
      }
      body += `
        <input type="hidden" value="` + item + `" name="original_name" >
        <p>
        <label for="name">Name</label>
        <input type="text" value="` + item + `" name="name" required>
        <p>
        <label for="description">Description</label>
        <textarea value="` + description + `" name="description" required>` + description + `</textarea>
        <p>`

    } else {
      body += `
        <input type="text" placeholder="/garden/tomatoes/soil" name="name" required>
        <p>
        <label for="description">Description</label>
        <textarea placeholder="Soil moisture sensor" name="description" required></textarea>
        <p>`
    }
    body += `
            <button
              class="button button-black "
              type="submit"> ok
            </button>
          </fieldset>
        </form>            <button
              class="button button-black "
              type="submit"
              onclick='dashboardReset()'> back
            </button>
      </div>
    </div>`
    $("#dashboard-view").html(body)
  }

  setupActionInput(item, postUrl, use_active) {
    console.log("setupActionInput()")
  }

  setupSignalInput(item, postUrl, use_active) {
    console.log("setupSignalInput()")
  }

  routeModificationRequest(item, action_type) {

   let use_active = (action_type == "edit")
   $("#dashboard-view").html("")
    switch (this.view) {
      case DashboardViews.ASSETS:
        this.setupAssetInput(item, "/app/" + action_type + "/asset", use_active)
        break;
      case DashboardViews.ACTIONS:
        this.setupActionInput(item, "/app/" + action_type + "/action", use_active)
        break;
      case DashboardViews.SIGNALS:
        this.setupSignalInput(item, "/app/" + action_type + "/signal", use_active)
        break;
      default:
        console.log("something should be on fire...")
        this.alerts.error("Internal error: Invalid view name for current view!")
        return
    }
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
          obj.alerts.error("error deleting:" + item)
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
