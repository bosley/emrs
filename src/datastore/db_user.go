package datastore

type user struct {
	id   int
	name string
	auth string
}

func (c *controller) Validate(username string, password string) bool {
	return false
}

func (c *controller) AddUser(name string, password string) error {
	return nil
}

func (c *controller) UpdatePassword(username string, password string) error {
	return nil
}

func (c *controller) DeleteUser(username string) bool {
	return false
}

const users_create = `insert into users (id, username, authhash) values (NULL, ?, ?)`
const users_get = `select username, id from users where username = ?`
const users_update = `update users set username = ?, authhash = ? where username = ?`
const users_delete = `delete from users where username = ?`

func (c *controller) retrieveUser(username string) *user {

	return nil
}
