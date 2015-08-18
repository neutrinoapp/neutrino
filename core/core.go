package realbase

type Config interface {
	GetConnectionString() string
}

type config struct {
	connectionString string
}

var c *config = nil
func Initialize(connectionString string) {
	if c != nil {
		panic("Initialize must be called once")
	}

	c = &config{connectionString}
}

func GetConfig() Config {
	if c == nil {
		c = &config{"localhost:27017"}
	}

	return c
}

func (c *config) GetConnectionString() string {
	return c.connectionString
}