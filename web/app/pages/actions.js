class PageAction {
  constructor(alerts, getTopo, getActions) {
    this.alerts = alerts
    this.selected = false
    this.getTopo = getTopo
    this.getActions = getActions
  }

  setIdle(contentHook) {
    this.selected = false
    $(contentHook).remove("#emrs-generated-viewport")
  }

  setSelected(contentHook) {
    this.selected = true
    $(contentHook).html("actions")

    $(contentHook).html('<div id="emrs-generated-viewport"></div>')
    $("#emrs-generated-viewport").html(`<div id="action-view"></div>`)

    this.loadView()


    /*


      TODO:


          Need to create the page to associate signals to actions.

          we can use this.getActions() to get a list of files that the
          user has on disk.

          adding an action is adding the file path
          to the EmrsAction.info

          Signals are in the topo representation, and sig


      */



  }

  updateTopo() {
    this.topo = this.getTopo()
  }

  updateActionsFileList() {
    this.action_files = this.getActions()
  }

  updateInfo() {
    this.updateTopo()
    this.updateActionsFileList()
  }

  loadView() {
    this.updateInfo()

    $("#action-view").html("")
    
    $("#action-view").append(`
      <button class="edit-button" onclick="app.getPageActions().loadActionEditor(null);">+ Action</button>
      `)

    this.drawActionTable()
  }

  drawActionTable() {

    let headers = ["Action Name", "Assigned", "Description", ""]
    let content = []

    let actions = this.topo.Actions
    for (let i = 0; i < actions.length; i++) {
      content.push([
        actions[i].Header.Name,
        "N/I",
        actions[i].Header.Description,
        `
         <button class="edit-button" onclick="app.getPageActions().editAction('`+ actions[i].Header.Name +`');">Edit</button>
         <button class="delete-button" onclick="app.getPageActions().deleteAction('`+ actions[i].Header.Name +`');">Delete</button>`,
      ])
    }

    $("#action-view").append(new Table(headers, content).value())
  }












  loadActionEditor(action) {

    console.log("draw action form")

    if (null === action) {
      console.log("load action editor with no pre-filled data")
    }

    $("#action-view").html(`
          <div class="row">
            <div class="column column-50" id="emrs_input_col">
                <input placeholder="Name..." type="text" required id="emrs_name_input" required>
                <textarea placeholder="Description..." id="emrs_description_input" required></textarea>
            </div>
            <div class="column column-10" id="emrs_input_button_col">
              <button class="edit-button" onclick="app.getPageActions().addAction()">ADD</button>
              <button class="edit-button" onclick="app.getPageActions().loadView();"> << </button>
            </div>
      `)
  }

  addAction() {

    console.log("add action")
    
    let name = $("#emrs_name_input").val()
    $("#emrs_name_input").val('')

    let description = $("#emrs_description_input").val()
    $("#emrs_description_input").val('')

    let type = $("#emrs_checkbox_one_input").val()
    $("#emrs_checkbox_one_input").val('')

    let file = $("#emrs_selected_file").val()
    $("#emrs_selected_file").val('')

    console.log(file)

    let x = new EmrsAction(
      new EmrsHeader(name, description),
      ActionType.FILE,
      "SomeFillerText")

    console.log(x)
    let msg = new ApiMsg(
      ApiOp.ADD,
      ApiSubject.ACTION,
      JSON.stringify(x))

    console.log(JSON.stringify(msg))
    $.ajax({
      type: "POST",
      url: "api/update" + getApiKeyUrlParam(),
      dataType: 'json',
      data: JSON.stringify(msg),
      async: false,
      error: ((function(obj){
        return function(){ 
          app.alerts.error("Failed to add action")
          obj.loadView()
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          obj.loadView()
          app.alerts.info("Action Added")
        }
      })(this))
    })
  }

  editAction(name) {
    console.log("edit", name)

    
  }

  deleteAction(name) {
    let msg = new ApiMsg(
      ApiOp.DEL,
      ApiSubject.ACTION,
      name)

    $.ajax({
      type: "POST",
      url: "api/update" + getApiKeyUrlParam(),
      dataType: 'json',
      data: JSON.stringify(msg),
      async: false,
      error: ((function(obj){
        return function(){ 
          app.alerts.error("Failed to delete action")
          obj.loadView()
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          console.log("complete", data)
          obj.loadView()
          app.alerts.info("Action Deleted")
        }
      })(this))
    })
  }

}

