class AppSession {

  constructor() {

    this.id = "" 
    this.user = ""
    this.version = ""
    this.valid = false
    this.first = true
  }

  quit() {
    this.id = ""
    this.user = ""
    this.version = ""
    this.valid = false
  }

  validate() {

    this.valid = false;

    if (!this.first && this.id == "") {
      console.log("attempt to validate a completed session")
      return
    }

    $.ajax({
      type: "GET",
      url: "/app/session",
      dataType: 'json',
      async: false,
      error: ((function(obj){
        return function(){ 
          console.log("session not valid")
        }
      })(this)),
      success: ((function(obj){
        return function(data){
          obj.runSessionChecks(data)
        }
      })(this))
    })
  }

  isValid() {
    return this.valid
  }

  runSessionChecks(json) {

    if (this.first) {
      this.id = json["session"]
      this.user = json["user"]
      this.version = json["version"]
      this.first = false
      this.valid = true;
      return
    }

    if (this.version != json["version"]) {
      console.log(
        "session: unexpected version change [",
        json["version"],
        "] expected [",
        this.version,
        "]")
      return
    }

    if (this.user != json["user"]) {
      console.log(
        "session: unexpected user [",
        json["user"],
        "] expected [",
        this.user,
        "]")
      return
    }

    if (this.id != json["session"]) {
      console.log(
        "session: unexpected session id [",
        json["session"],
        "] expected [",
        this.id,
        "]")
      return
    }

    this.valid = true;
  }
}
