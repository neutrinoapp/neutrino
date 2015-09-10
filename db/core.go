package db

var config map[string]interface{}
func Initialize(c map[string]interface{}) {
	if config != nil {
		panic("Initialize must be called once")
	}

	config = c
}