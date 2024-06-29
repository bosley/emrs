const AlertLevel = Object.freeze({
  INFO: "info",
  ERROR: "error",
  WARN: "warning",
  SUCCESS: "success",
})

// An alert to show the user. This alert will exist for "N" "cycles" determined
// by the Alerts object. Once the alert has been shown for "N" cycles, the alert
// will be removed from the screen
class Alert {
  constructor(level, message, max, img) {
    this.shown = 0
    this.max = max
    this.value = $('<div class="alert-tile emrs_alert alert-' + level + '"><img class="panel-tile-img" src="' + img + '"><br>' + message + '</div><p>')
    console.log("new alert:", level)
  }
  shownInd() {
    this.shown += 1
  }
  isIncomplete() {
    return this.shown < this.max
  }
}

// When an alert is asked to be displayed, the alerts object
// may or may not kick off a timeout to continue displaying
// whatever alerts it has stored. Once it runs out of alerts
// the timeout callback cycle will conclude until the next
// alert is given
class Alerts {
  constructor(target) {
    this.active = []
    this.shown = []
    this.target = target
    this.cycling = false
    this.interval = 1000
  }

  info(message) {
    this.active.push(
      new Alert(AlertLevel.INFO, message, 8, "img/icons8-info-100.png"));
    this.kickoff()
  }

  warning(message) {
    this.active.push(
      new Alert(AlertLevel.WARN, message, 5, "img/icons8-important-100.png"));
    this.kickoff()
  }

  error(message) {
    this.active.push(
      new Alert(AlertLevel.ERROR, message, 5, "img/icons8-error-100.png"));
    this.kickoff()
  }

  success(message) {
    this.active.push(
      new Alert(AlertLevel.SUCCESS, message, 10, "img/icons8-done-100.png"));
    this.kickoff()
  }

  kickoff() {
    if (this.cycling) {
      return
    }
    this.execute()
  }

  initAlertCycle() {
    setTimeout(
      (function(obj){
        return function(){ obj.execute(); }
      })(this), this.interval)
  }

  execute() {
    // Remove whatever alerts currently exist
    $("div").remove(".emrs_alert")

    if (this.active.length == 0) {
      this.cycling = false
      return
    }
    this.cycling = true

    // Show the alert, indicate to the alert that we have shown it
    for (let i = this.active.length-1; i >= 0; i--) {
      $(this.target).append(this.active[i].value)
      this.active[i].shownInd()
    }

    // Remove all alerts that are completed
    this.active = this.active.filter((entry) => entry.isIncomplete())

    // Kick off another timeout to check again
    this.initAlertCycle()
  }
}
