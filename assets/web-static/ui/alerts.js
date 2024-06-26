const AlertLevel = Object.freeze({
  INFO: "info",
  ERROR: "danger",
  WARN: "warning",
  SUCCESS: "success",
})

class Alert {
  constructor(level, message, max) {
    this.shown = 0
    this.max = max
    this.value = $('<div class="emrs_alert alert alert-' + level + '", role="alert", id="' + this.emrsId + '"> ' + message + '</div>')
  }
  shownInd() {
    this.shown += 1
  }
  isIncomplete() {
    return this.shown < this.max
  }
}

class Alerts {
  constructor(target) {
    this.active = []
    this.shown = []
    this.target = target
  }

  info(message) {
    this.active.push(
      new Alert(AlertLevel.INFO, message, 10));
  }
  warning(message) {
    this.active.push(
      new Alert(AlertLevel.WARN, message, 12));
  }
  error(message) {
    this.active.push(
      new Alert(AlertLevel.ERROR, message, 15));
  }
  success(message) {
    this.active.push(
      new Alert(AlertLevel.SUCCESS, message, 20));
  }
  tick() {
    $("div").remove(".emrs_alert")

    for (let i = 0; i < this.active.length; i++) {
      $(this.target).append(this.active[i].value)
      this.active[i].shownInd()
    }
    this.active = this.active.filter((entry) => entry.isIncomplete())
  }
}
