const DashboardViews = Object.freeze({
  SECTORS: "sectors",
  SECTORS_ADD: "sectorsAdd",
  ASSETS:  "asset",
  ACTIONS: "action",
  SIGNALS: "signal",
})

class PageDashboard {
  constructor(alerts) {
    this.alerts = alerts
    this.selected = false
    this.selectedSector = ""
    this.topo = null
  }

  setIdle(contentHook) {
    this.selected = false
    $(contentHook).remove("#emrs-generated-dashboard")
  }

  setSelected(contentHook) {
    console.log("setSelected")
  
    $(contentHook).html('<div id="emrs-generated-dashboard"></div>')

    $("#emrs-generated-dashboard").html(`<div id="dashboard-view"></div>`)
    
    this.contentHook = contentHook
    this.updateTopo()
    this.changeViews(DashboardViews.SECTORS)
    this.selected = true
  }

  updateTopo() {
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
  }

  setTable(table) {
  }

  changeViews(view) {
    $("#dashboard-view").html("")

    let headers = []
    let content = []
    switch (view) {
      case DashboardViews.SECTORS:
        this.drawSectors()
        break
      case DashboardViews.SECTORS_ADD:
        this.drawSectorsForm()
        break
      case DashboardViews.ASSETS:
        headers = ["Asset Name", "Description", ""]
        break;
      case DashboardViews.ACTIONS:
        headers = ["Action Name", "Assigned", "Description", ""]
        break;
      case DashboardViews.SIGNALS:
        headers = ["Signal Name", "In-Use", "Description", ""]
        break;
      default:
        this.alerts.error("Internal error: Invalid view name")
        return
    }
    this.view = view
  }

  drawSectors() {
    $("#dashboard-view").html("")
    
    $("#dashboard-view").append(`
      <button class="edit-button" onclick="app.getPageDashboard().drawSectorsForm();">+ Sector</button>
      `)

    let headers = ["Sector Name", "Assets", "Description", ""]
    let content = []
    for (let i = 0; i < this.topo.Sectors.length; i++) {
      content.push([
        this.topo.Sectors[i].Header.Name,
        this.topo.Sectors[i].Assets.length,
        this.topo.Sectors[i].Header.Description,
        `<button class="edit-button" onclick="app.getPageDashboard().deleteSector('`+ this.topo.Sectors[i].Header.Name +`');">Delete</button>`,
      ])
    }

    $("#dashboard-view").append(new Table(headers, content).value())
  }

  beginInput(path) {
    $("#dashboard-view").append(`
          <div class="row">
            <div class="column column-50" id="emrs_input_col">
              <input placeholder="Name..." type="text" required id="emrs_name_input">
              <textarea placeholder="Description..." id="emrs_description_input"></textarea>
            </div>
            <div class="column column-10">
              <button class="edit-button" onclick="` + path + `">ADD</button>
            </div>
      `)
  }

  endInput() {
    $("#dashboard-view").append(`
          </div>
      `)
  }

  drawSectorsForm() {
    $("#dashboard-view").html("")
    this.beginInput("app.getPageDashboard().addSector();")
    this.endInput()
  }

  deleteSector(name) {
    let msg = new ApiMsg(
      ApiOp.DEL,
      ApiSubject.SECTOR,
      name)

    $.ajax({
      type: "POST",
      url: "api/update" + getApiKeyUrlParam(),
      dataType: 'json',
      data: JSON.stringify(msg),
      async: false,
      error: ((function(obj){
        return function(){ 
          app.alerts.error("Failed to delete sector")
          obj.drawSectors()
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          console.log("complete", data)
          obj.updateTopo()
          obj.drawSectors()
        }
      })(this))
    })
  }

  addSector() {

    console.log("add sector")
    
    let name = $("#emrs_name_input").val()
    $("#emrs_name_input").val('')

    let description = $("#emrs_description_input").val()
    $("#emrs_description_input").val('')

    let sector = new EmrsSector(
      new EmrsHeader(name, description), 
      [])

    let msg = new ApiMsg(
      ApiOp.ADD,
      ApiSubject.SECTOR,
      JSON.stringify(sector))

    console.log(JSON.stringify(msg))
    $.ajax({
      type: "POST",
      url: "api/update" + getApiKeyUrlParam(),
      dataType: 'json',
      data: JSON.stringify(msg),
      async: false,
      error: ((function(obj){
        return function(){ 
          app.alerts.error("Failed to add sector")
          obj.drawSectors()
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          console.log("complete", data)
          obj.updateTopo()
          obj.drawSectors()
        }
      })(this))
    })
  }


}
