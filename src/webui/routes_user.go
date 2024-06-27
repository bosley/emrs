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
Posts to /auth (routeAuth) - shown below
*/
func (wc *controller) routeLogin(c *gin.Context) {
	_, attempted := c.Get(loginAttemptKey)
	c.HTML(200, "login.html", gin.H{
		"PageHeader":  buildPageHeader("Login"),
		"NavData":     buildNavData(c),
		"PrevAttempt": attempted,
	})
}

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

	c.HTML(200, "message.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Logged Out"),
		"Message":    "You have been logged out",
		"ShowLogin":  true,
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
		wc.routeLogin(c)
		return
	}
	slog.Warn("auth success", "user", username)

	session.Set(sessionKeyUserId, username)

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
		c.HTML(http.StatusUnauthorized, "message.html", gin.H{
			"NavData":    buildNavData(c),
			"PageHeader": buildPageHeader("Not Available"),
			"Message":    "Creating users is an unauthorized action after initial setup",
			"ShowLogin":  false,
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

	c.Redirect(http.StatusMovedPermanently, "/login")
}

/*
New user (as described above) is not made available after the very first
user is created on the system, which only happens on the very first launch
of EMRS. This /new/user route just displays the html form for username/pass
and then posts to /create/user (above)
*/
func (wc *controller) routeNewUser(c *gin.Context) {

	if !wc.appCore.RequiresSetup() {
		c.HTML(http.StatusUnauthorized, "message.html", gin.H{
			"NavData":    buildNavData(c),
			"PageHeader": buildPageHeader("Not Available"),
			"Message":    "Creating users is an unauthorized action after initial setup",
			"ShowLogin":  false,
		})
		return
	}

	_, exists := c.Get(existingUserKey)
	c.HTML(200, "new_user.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Create User"),
		"UserExists": exists,
	})
}
