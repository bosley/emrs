class PageAction {
  constructor(alerts, getTopo) {
    this.alerts = alerts
    this.selected = false
    this.getTopo = getTopo

/*
    I think instead of handing data to the thing directly,
    we should hand an emrs context. This way we can have access from one
    object, all of the api and the data

*/
    this.default_action_body = `

func Exec(data string) {
  println("Action called upon with: ", data)
}
    `
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

    this.drawActionTable()
  }

  updateTopo() {
    this.topo = this.getTopo()
  }

  drawActionTable() {
    this.updateTopo()
    

    $("#action-view").html("")
    
    $("#action-view").append(`
      <button class="edit-button" onclick="app.getPageActions().drawActionForm();">+ Action</button>
      `)

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

  drawActionForm() {
    console.log("draw action form")

    $("#action-view").html(`
          <div class="row">
            <div class="column column-50" id="emrs_input_col">
                <input placeholder="Name..." type="text" required id="emrs_name_input" required>
                <textarea placeholder="Description..." id="emrs_description_input" required></textarea>
            </div>
            <div class="column column-10" id="emrs_input_button_col">
              <button class="edit-button" onclick="app.getPageActions().addAction()">ADD</button>
              <button class="edit-button" onclick="app.getPageActions().drawActionTable();"> << </button>
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
          obj.drawActionTable()
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          obj.updateTopo()
          obj.drawActionTable()
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
          obj.drawActionTable()
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          console.log("complete", data)
          obj.updateTopo()
          obj.drawActionTable()
          app.alerts.info("Action Deleted")
        }
      })(this))
    })
  }

}

