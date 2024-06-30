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
    this.maxTableDesc = 20 // must be > 3
    this.maxTableItemName = 20 // must be > 3
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

  changeViews(view) {
    $("#dashboard-view").html("")

    let buttonText = "[INVALID]"
    switch (view) {
      case DashboardViews.ASSETS:
        buttonText = "+ asset"
        break;
      case DashboardViews.ACTIONS:
        buttonText = "+ action"
        break;
      case DashboardViews.SIGNALS:
        buttonText = "+ signal"
        break;
      default:
        this.alerts.error("Internal error: Invalid view name")
        return
    }

    this.makeTable("name", "description", "")
    $("#emrs_dashboard_add_item_button").html(buttonText)

    this.view = view
    this.populateTable()
  }

  reload() {
    this.changeViews(this.view)
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

  retrieveCachedItem(item) {
      let items = this.data_cache[this.view]
      for (let i = 0; i < items.length; i++) {
        if (items[i]["Col1"] === item) {
          return items[i]
        }
      }
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

            let name = item["Col1"]
            if (name.length > obj.maxTableItemName - 3) {
              name = name.substring(0, obj.maxTableItemName - 3) + "..."
            }

            let description = item["Col2"]
            if (description.length > obj.maxTableDesc - 3) {
              description = description.substring(0, obj.maxTableDesc - 3) + "..."
            }

            obj.data_cache["short_name"] = name 
            obj.data_cache["short_description"] = description

            $("#emrs-dashboard-table-body").append(
              `<tr>
                <td>` + name + `</td>
                <td>` + description + `</td>
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

  setupAssetInput(item, postUrl, use_active) {
    let body = 
    `<div class="row">
      <div class="column column-50">
        <form method="POST" action="` + postUrl + `">
          <fieldset>
            `

    if (use_active) {
      let description = this.retrieveCachedItem(item)["Col2"]

      body += `
        <label for="name">Name</label>
        <input type="text" value="` + item + `" name="name" required>
        <p>
        <label for="description">Description</label>
        <textarea value="` + description + `" name="description" required>` + description + `</textarea>
        <p>
        <input type="hidden" value="` + item + `" name="original_name" >
        <p>`
    } else {
      body += `
        <label for="name">Name</label>
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
    let body = 
    `<div class="row">
      <div class="column column-50">
        <form method="POST" action="` + postUrl + `">
          <fieldset>
            `

    if (use_active) {
      let cached_obj = this.retrieveCachedItem(item)
      let description = cached_obj["Col2"]
      let execInfo = cached_obj["Col3"]

      body += `
        <label for="name">Name</label>
        <input type="text" value="` + item + `" name="name" required>
        <p>
        <label for="description">Description</label>
        <textarea value="` + description + `" name="description" required>` + description + `</textarea>
        <p>
        <label for="execution_info">Execution Information</label>
        <textarea value="` + execInfo + `" name="execution_info" required>` + execInfo + `</textarea>
        <p>
        <input type="hidden" value="` + item + `" name="original_name" >
        <p>
        `
    } else {
      body += `
        <label for="name">Name</label>
        <input type="text" placeholder="[UNDER CONSTRUCTION]" name="name" required>
        <p>
        <label for="description">Description</label>
        <textarea placeholder="[UNDER CONSTRUCTION]" name="description" required></textarea>
        <p>
        <label for="execution_info">Execution Information</label>
        <textarea placeholder=" [UNDER CONSTRUCTION] " name="execution_info" required></textarea>
        <p>
        `
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

  setupSignalInput(item, postUrl, use_active) {
    let body = 
    `<div class="row">
      <div class="column column-50">
        <form method="POST" action="` + postUrl + `">
          <fieldset>
            `

    if (use_active) {
      let cached_obj = this.retrieveCachedItem(item)
      let description = cached_obj["Col2"]
      let triggers = cached_obj["Col3"]

      body += `
        <label for="name">Name</label>
        <input type="text" value="` + item + `" name="name" required>
        <p>
        <label for="description">Description</label>
        <textarea value="` + description + `" name="description" required>` + description + `</textarea>
        <p>
        <label for="triggers">Triggers</label>
        <textarea value="` + triggers + `" name="triggers" required>` + triggers + `</textarea>
        <p>
        <input type="hidden" value="` + item + `" name="original_name" >
        <p>`
    } else {
      body += `
        <label for="name">Name</label>
        <input type="text" placeholder="[UNDER CONSTRUCTION]" name="name" required>
        <p>
        <label for="description">Description</label>
        <textarea placeholder="[UNDER CONSTRUCTION]" name="description" required></textarea>
        <p>
        <label for="triggers">Triggers</label>
        <textarea placeholder=" [UNDER CONSTRUCTION] " name="triggers" required></textarea>
        <p>
        `
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
}
