package neutrinoclient

import (
	"github.com/go-neutrino/neutrino/models"
	"github.com/go-neutrino/neutrino/log"
)

type NeutrinoEvent struct {
	Code string
	Data string
}

type NeutrinoData struct {
	*NeutrinoClient
	Name string
	Event chan NeutrinoEvent
}

const (
	EVENT_UPDATE string = "update"
	EVENT_CREATE string = "create"
	EVENT_DELETE string = "delete"
)

var dataMap map[string]*NeutrinoData

func init() {
	dataMap = make(map[string]*NeutrinoData)
}

func newUpdateEvent(data string) NeutrinoEvent {
	return NeutrinoEvent{EVENT_UPDATE, data}
}

func newCreateEvent(data string) NeutrinoEvent {
	return NeutrinoEvent{EVENT_CREATE, data}
}

func newDeleteEvent(data string) NeutrinoEvent {
	return NeutrinoEvent{EVENT_DELETE, data}
}

func (c *NeutrinoClient) Data(name string) *NeutrinoData {
	//singleton data
	if (dataMap[name] != nil) {
		return dataMap[name]
	}

	//TODO: find out why multiple data objects do not allow
	//multiple channel receiving
	//workaround: use singleton per type name

	d := &NeutrinoData{
		NeutrinoClient: c,
		Name: name,
		Event: make(chan NeutrinoEvent),
	}

	c.registerDataListener(d)

	dataMap[name] = d

	return d
}

func (d *NeutrinoData) GetUrl() string {
	return d.NeutrinoClient.Addr + "app"
}

func (d *NeutrinoData) onDataMessage(m models.JSON) {
	data, err := m.String()

	if err != nil {
		log.Error(err)
		return
	}

	switch m["op"] {
	case EVENT_CREATE:
		d.dispatchEvent(newCreateEvent(data))
	case EVENT_DELETE:
		d.dispatchEvent(newDeleteEvent(data))
	case EVENT_UPDATE:
		d.dispatchEvent(newUpdateEvent(data))
	}
}

func (d *NeutrinoData) dispatchEvent(ev NeutrinoEvent) {
	log.Info("Dispatching event:", ev.Code, ev.Data)
	d.Event <- ev
}

func (d *NeutrinoData) getDataUrl() string {
	return d.NeutrinoClient.getAppUrl() + "/data"
}