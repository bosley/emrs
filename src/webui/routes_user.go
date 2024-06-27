package webui

import (
	"emrs/badger"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	sessionKeyUserId = "user-id"
	loginAttemptKey  = "login-failure"
	existingUserKey  = "user-exists"
)

/*
Retrieves the user data from the session and deletes it.
Once this occurs we send the user to the friendly "message"
page where we same something factual and offer a button to
log back in
*/
func (wc *controller) routeLogout(c *gin.Context) {

	session := sessions.Default(c)
	token := session.Get(sessionKeyUserId)
	if token == nil {
		slog.Debug("invalid session token for logout")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}

	session.Delete(sessionKeyUserId)
	if err := session.Save(); err != nil {
		slog.Debug("Unable to save session")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.HTML(200, "window.html", gin.H{
    "Prompting": false,
    "Message": "You have been logged out",
    "ShowLogin": true,
	})
}

/*
	Clearly, this is the route that is POSTed to to authenticate someone.
	When someone posts here we use badger to check for a match.
	Badger is kind of an authentication bag with "Vouchers" (basically JWTs)
	and key signing/verification stuff.

	As of 24/June/24 badger executes the following for generating authentication
	hashes, and to verify that something matches a stored authentication hash

```

	func RawIsHashMatch(raw []byte, hashed []byte) error {
		defer zeroArr(raw)
		return bcrypt.CompareHashAndPassword(hashed, raw)
	}

	func Hash(raw []byte) ([]byte, error) {
		defer zeroArr(raw)
		return bcrypt.GenerateFromPassword(raw, bcrypt.DefaultCost)
	}

```
*/
func (wc *controller) routeAuth(c *gin.Context) {

	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	db := wc.appCore.GetUserStore()

	storedHash := db.GetAuthHash(username)

	if storedHash == nil ||
		badger.RawIsHashMatch([]byte(password), []byte(*storedHash)) != nil {
		slog.Warn("auth failure", "user", username)
		c.Set(loginAttemptKey, true)
		wc.routeIndex(c)
		return
	}
	slog.Warn("auth success", "user", username)

  userbadge, err := badger.New(badger.Config{
    Nickname: username,
  })

  if err != nil {
	  slog.Warn("badger failure", "user", username, "error", err.Error())
	  c.JSON(500, gin.H{
      "code": 500,
      "message": err.Error(),
	  })
    return 
  }

	session.Set(sessionKeyUserId, userbadge.EncodeIdentityString())

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	slog.Info("user logged in", "user", username)

	c.Redirect(http.StatusMovedPermanently, "/emrs/")
}

/*
   The routes /create/user and /new/user are utilized to make the
   very first user account on the system. The only time they are available
   is when the EMRS instance is brand new, meaning that the identity of
   the server had not yet been established in the database given to
   EMRS at start time.

   We _need_ to add other checks to ensure that the startup conditons are actually
   present, though the core server identity is the biggest tell. We could add checks
   to see if there are 0 users AND 0 assets etc...

   Once the user is created, the webui indicates to appCore that the initial
   setup is complete via `IndicateSetupComplete.` Once this happens both
   http routes will be disabled, and the index `/` page will no longer display
   the "Create User" interface.
*/

func (wc *controller) routeCreateUser(c *gin.Context) {

	if !wc.appCore.RequiresSetup() {
	  c.HTML(http.StatusUnauthorized, "window.html", gin.H{
      "Prompting": false,
      "Message": "Creating users is an unauthorized action after initial setup",
      "ShowLogin": false,
	  })
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	hash, err := badger.Hash([]byte(password))

	if err != nil {
		slog.Error("Error creating user", "err", err.Error())
		panic("Error with badger.Hash() - panik because this shouldn't happen")
	}

	db := wc.appCore.GetUserStore()

	if err := db.AddUser(username, string(hash)); err != nil {
		slog.Error("Error creating user", "err", err.Error())
		c.Set(existingUserKey, true)
		wc.routeNewUser(c)
		return
	}

	//  Here we disable access to /new/user and /create/user
	//
	wc.appCore.IndicateSetupComplete()

	c.Redirect(http.StatusMovedPermanently, "/")
}

/*
New user (as described above) is not made available after the very first
user is created on the system, which only happens on the very first launch
of EMRS. This /new/user route just displays the html form for username/pass
and then posts to /create/user (above)
*/
func (wc *controller) routeNewUser(c *gin.Context) {

  slog.Debug("Creating new user")

	if !wc.appCore.RequiresSetup() {
	  c.HTML(http.StatusUnauthorized, "window.html", gin.H{
      "Prompting": false,
      "Message": "Creating users is an unauthorized action after initial setup",
      "ShowLogin": false,
	  })
		return
	}

	_, exists := c.Get(existingUserKey)
	c.HTML(200, "window.html", gin.H{
    "Topic": "Create Account",
    "Prompting": true,
		"PrevAttempt": false,
    "PostTo": "/create/user",
    "Prompt": "EMRS New User",
		"UserExists": exists,
	})
}

func (wc *controller) routeSessionInfo(c *gin.Context) {
  userInfo := getLoggedInUser(c)
  if userInfo == nil {
	  c.JSON(500, gin.H{
      "code": 500,
      "message": "Unable to retrieve user information",
	  })
    return
  }

  badge, err := badger.DecodeIdentityString(userInfo.(string))
  if err != nil {
	  c.JSON(500, gin.H{
      "code": 500,
      "message": err.Error(),
	  })
    return
  }

	c.JSON(200, gin.H{
    "session": badge.Id(),
    "user": badge.Nickname(),
    "version": wc.appCore.GetVersion(),
	})
}
