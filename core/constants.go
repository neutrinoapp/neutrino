package realbase

type constants struct {
}

var Constants = &constants{}

func (c *constants) DatabaseName() string {
	return "realbase"
}

func (c *constants) UsersCollection() string {
	return "users"
}

func (c *constants) ApplicationsCollection() string {
	return "applications"
}