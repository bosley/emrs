const DashboardViews = Object.freeze({
  SECTORS: "sectors",
  ASSETS:  "asset",
  ACTIONS: "action",
  SIGNALS: "signal",
})

class EmrsApiAssetCUD {
  constructor(sector, asset) {
    this.sector = sector
    this.asset = asset
  }
}

class PageDashboard {
  constructor(alerts, getTopo) {
    this.alerts = alerts
    this.selected = false
    this.selectedSector = ""
    this.topo = null
    this.getTopo = getTopo
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
    this.topo = this.getTopo()
  }

  changeViews(view) {
    $("#dashboard-view").html("")

    let headers = []
    let content = []
    switch (view) {
      case DashboardViews.SECTORS:
        this.drawSectors()
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
        `<button class="edit-button" onclick="app.getPageDashboard().loadSector('`+ this.topo.Sectors[i].Header.Name +`');">Edit</button>
         <button class="delete-button" onclick="app.getPageDashboard().deleteSector('`+ this.topo.Sectors[i].Header.Name +`');">Delete</button>`,
      ])
    }

    $("#dashboard-view").append(new Table(headers, content).value())
  }

  beginInput(path) {
    $("#dashboard-view").append(`
          <div class="row">
            <div class="column column-50" id="emrs_input_col">
              <input placeholder="Name..." type="text" required id="emrs_name_input">
              <textarea placeholder="Description..." id="emrs_description_input" required></textarea>
    <!--  Not important enough atm        <textarea placeholder="Comma,sep,tags..." id="emrs_tags_input" required></textarea> -->
            </div>
            <div class="column column-10" id="emrs_input_button_col">
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
          app.alerts.info("Sector Deleted")
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

//    let tags = $("#emrs_tags_input").val().split(",")
  //  $("#emrs_tags_input").val('')

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
          app.alerts.info("Sector Added")
        }
      })(this))
    })
  }

  retrieveActiveSector() {
    for (let i = 0; i < this.topo.Sectors.length; i++) {
      if (this.topo.Sectors[i].Header.Name === this.active_sector) {
        return this.topo.Sectors[i]
      }
    }
    return null
  }

  loadActiveSector() {
    if (null != this.active_sector) {
      this.loadSector(this.active_sector)
    } else {
      this.alers.error("no active sector selected")
    }
  }

  loadSector(sectorName) {
    this.active_sector = sectorName
    console.log("load sector", this.active_sector)

    let sector = this.retrieveActiveSector()
    $("#dashboard-view").html("")

    
    $("#dashboard-view").append(`
      <button class="edit-button" onclick="app.getPageDashboard().drawSectors();"> << </button>
      `)
    
    $("#dashboard-view").append(`
      <button class="edit-button" onclick="app.getPageDashboard().loadAssetForm();">+ ASSET</button>
      `)
    let headers = ["Asset Name", "Description", ""]
    let content = []
    for (let i = 0; i < sector.Assets.length; i++) {
      content.push([
        sector.Assets[i].Header.Name,
        sector.Assets[i].Header.Description,
        `
        <!-- <button class="edit-button" onclick="app.getPageDashboard().editAsset('`+ sector.Assets[i].Header.Name +`');">Edit</button> -->


         <button class="delete-button" onclick="app.getPageDashboard().deleteAsset('`+ sector.Assets[i].Header.Name +`');">Delete</button>`,
      ])
    }

    $("#dashboard-view").append(new Table(headers, content).value())

  }

  loadAssetForm() {
    $("#dashboard-view").html("")
    this.beginInput("app.getPageDashboard().addAsset();")
    this.endInput()

    $("#emrs_input_button_col").append(`
      <button class="edit-button" onclick="app.getPageDashboard().loadActiveSector();"> << </button>
      `)
  }

  addAsset() {
    console.log("add sector")
    
    let name = $("#emrs_name_input").val()
    $("#emrs_name_input").val('')

    let description = $("#emrs_description_input").val()
    $("#emrs_description_input").val('')

    let data  = new EmrsApiAssetCUD(
      this.active_sector,
      new EmrsAsset(
        new EmrsHeader(name, description)))

    let msg = new ApiMsg(
      ApiOp.ADD,
      ApiSubject.ASSET,
      JSON.stringify(data))

    console.log(JSON.stringify(msg))
    $.ajax({
      type: "POST",
      url: "api/update" + getApiKeyUrlParam(),
      dataType: 'json',
      data: JSON.stringify(msg),
      async: false,
      error: ((function(obj){
        return function(){ 
          app.alerts.error("Failed to add asset")
          obj.loadActiveSector()
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          console.log("complete", data)
          obj.updateTopo()
          obj.loadActiveSector()
          app.alerts.info("Asset added")
        }
      })(this))
    })
  }

  deleteAsset(name) {
    console.log("delete asset")

    let data  = new EmrsApiAssetCUD(
      this.active_sector,
      new EmrsAsset(
        new EmrsHeader(name, "")))

    let msg = new ApiMsg(
      ApiOp.DEL,
      ApiSubject.ASSET,
      JSON.stringify(data))

    console.log(JSON.stringify(msg))
    $.ajax({
      type: "POST",
      url: "api/update" + getApiKeyUrlParam(),
      dataType: 'json',
      data: JSON.stringify(msg),
      async: false,
      error: ((function(obj){
        return function(){ 
          app.alerts.error("Failed to delete asset")
          obj.loadActiveSector()
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          console.log("complete", data)
          obj.updateTopo()
          obj.loadActiveSector()
          app.alerts.info("Asset Deleted")
        }
      })(this))
    })
  }


}
