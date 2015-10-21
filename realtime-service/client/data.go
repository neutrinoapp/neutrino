package neutrinoclient

type NeutrinoEvent struct {
	Code string
	Data string
}

type NeutrinoData struct {
	*NeutrinoClient
	Name string
	Event chan(NeutrinoEvent)
}

const (
	EVENT_UPDATE string = "update"
	EVENT_CREATE string = "create"
	EVENT_DELETE string = "delete"
)

func newUpdateEvent(data string) *NeutrinoEvent {
	return &NeutrinoEvent{EVENT_UPDATE, data}
}

func newCreateEvent(data string) *NeutrinoEvent {
	return &NeutrinoEvent{EVENT_CREATE, data}
}

func newDeleteEvent(data string) *NeutrinoEvent {
	return &NeutrinoEvent{EVENT_DELETE, data}
}

func (d *NeutrinoData) GetUrl() string {
	return d.NeutrinoClient.Addr + "app"
}

func (c *NeutrinoClient) Data(name string) *NeutrinoData {
	return &NeutrinoData{
		NeutrinoClient: c,
		Name: name,
		Event: make(chan NeutrinoEvent),
	}
}

func (d *NeutrinoData) getDataUrl() string {
	return d.NeutrinoClient.getAppUrl() + "/data"
}