package neutrinoclient

import (
	"github.com/neutrinoapp/neutrino/src/common/log"
	"github.com/neutrinoapp/neutrino/src/common/messaging"
	"github.com/neutrinoapp/neutrino/src/common/models"
)

type NeutrinoEvent struct {
	Code string
	Data string
}

type NeutrinoData struct {
	NeutrinoClient *NeutrinoClient
	DataName       string
	Event          chan NeutrinoEvent
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
	if dataMap[name] != nil {
		return dataMap[name]
	}

	d := &NeutrinoData{
		NeutrinoClient: c,
		DataName:       name,
		Event:          make(chan NeutrinoEvent),
	}

	c.registerDataListener(d)

	dataMap[name] = d

	return d
}

func (d *NeutrinoData) GetUrl() string {
	return d.NeutrinoClient.ApiAddr + "app"
}

func (d *NeutrinoData) Create(m models.JSON) {
	messaging.GetMessageBuilder().Build(
		messaging.OP_CREATE,
		messaging.ORIGIN_CLIENT,
		m,
		nil,
		d.DataName,
		d.NeutrinoClient.AppId,
		d.NeutrinoClient.Token,
	).Send(d.NeutrinoClient.WebsocketClient.GetConnection())
}

func (d *NeutrinoData) Update(id string, m models.JSON) {
	m["_id"] = id
	messaging.GetMessageBuilder().Build(
		messaging.OP_UPDATE,
		messaging.ORIGIN_CLIENT,
		m,
		nil,
		d.DataName,
		d.NeutrinoClient.AppId,
		d.NeutrinoClient.Token,
	).Send(d.NeutrinoClient.WebsocketClient.GetConnection())
}

func (d *NeutrinoData) Delete(id string) {
	messaging.GetMessageBuilder().Build(
		messaging.OP_DELETE,
		messaging.ORIGIN_CLIENT,
		models.JSON{"_id": id},
		nil,
		d.DataName,
		d.NeutrinoClient.AppId,
		d.NeutrinoClient.Token,
	).Send(d.NeutrinoClient.WebsocketClient.GetConnection())
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
