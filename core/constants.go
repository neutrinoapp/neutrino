package realbase

type constants struct {
}

var Constants = &constants{}

func (c *constants) DatabaseName() string {
	return "Realbase"
}

func (c *constants) UsersCollection() string {
	return "Users"
}