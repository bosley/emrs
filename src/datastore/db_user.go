package datastore

func (c *controller) Validate(username string, password string) bool {
	return false
}

func (c *controller) AddUser(name string, password string, email string) error {
	return nil
}

func (c *controller) UpdatePassword(username string, password string) error {
	return nil
}

func (c *controller) DeleteUser(username string) bool {
	return false
}
