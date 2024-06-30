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

    switch (view) {
      case DashboardViews.ASSETS:
        this.setupContentAsset()
        break;
      case DashboardViews.ACTIONS:
        this.setupContentAction()
        break;
      case DashboardViews.SIGNALS:
        this.setupContentSignal()
        break;
      default:
        this.alerts.error("Internal error: Invalid view name")
        return
    }

    this.view = view
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

  setupContentAsset() {
    this.makeTable("name", "last seen", "")
    $("#dashboard-view").append(
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
    this.populateTable()
  }

  setupContentAction() {
    this.makeTable("name", "status", "")
    $("#dashboard-view").append(
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
    this.populateTable()
  }

  setupContentSignal() {
    this.makeTable("name", "assigned", "")
    $("#dashboard-view").append(
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
    this.populateTable()
  }

  setupEditAsset(name) {


    /*
        TODO:

        Right now when we hit "edit" the table is wiped and a form will display

        We need a map of maps in the dashboard view -> asset -> data mapping.

        This way when we edit, we can just pull from memory of the currently displayed

        table for autofilled fields in editing.


        Keep in mind I want to use the same screen for adding a new item rather than having
        some button at the bottom of the table with a single input field.

        Instead, an add button at the top should bring this same screen up, but unpopulated.

        Then when they submit, it will go add and the page will come back to /app to
        display the new table

      */

    console.log("edit", name)


    // TODO: Modify this to post to the update area.
    //
    // then modify the views to build their forms on "ADD" rather than have the form under the table
    //
    // depending on whertr we come from we will wan to post to different areas /create vs update/ etc
    // so we will want to function-this-out to change post destinations
    //
    // then we need to change backend to get post data like it shouldve
    //
    // then get all 6 variations working
    //
    // THEN WE CAN START REAL PROGRAMMING
    //
    //
    //
    $("#dashboard-view").append(
`
    <div class="row">
      <div class="column column-50">
        <form>
          <fieldset>
            <label for="assetName">Name</label>
            <input type="text" placeholder="..." id="assetName">
            <p>
            <label for="assetDesc">Description</label>
            <textarea placeholder=". . ." id="assetDesc"></textarea>
            <p>
            <button
              class="button button-black "
              type="submit"> ok
            </button>
          </fieldset>
        </form>
      </div>
    </div>
`
    )


  }

  setupEditAction(name) {
  }

  setupEditSignal(name) {
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

    console.log("edit item")

    $("#dashboard-view").html("")

    switch (this.view) {
      case DashboardViews.ASSETS:
        console.log("its an asset")
        this.setupEditAsset(item)
        break;
      case DashboardViews.ACTIONS:
        this.setupEditAction(item)
        break;
      case DashboardViews.SIGNALS:
        this.setupEditSignal(item)
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
