package neutrino

import (
	"encoding/json"
	"fmt"
	"github.com/gngeorgiev/sockjs-go/sockjs"
	"net/http"
)

type MessageService interface {
	InitSocketHandler() http.Handler
	GetSessions() []sockjs.Session
	Broadcast(message string)
	BroadcastJSON(message map[string]interface{})
}

type messageService struct {
	sessions []sockjs.Session
}

func NewMessageService() MessageService {
	return &messageService{make([]sockjs.Session, 0)}
}

func (m *messageService) InitSocketHandler() http.Handler {
	handler := sockjs.NewHandler("/socket", sockjs.DefaultOptions, func(so sockjs.Session) {
		m.sessions = append(m.sessions, so)
		msg, _ := so.Recv()
		fmt.Println(msg)
		//TODO:
	})

	return handler
}

func (m *messageService) GetSessions() []sockjs.Session {
	return m.sessions
}

func (m *messageService) Broadcast(message string) {
	for _, so := range m.sessions {
		if so.GetSessionState() > 1 {
			so.Send(message)
		} else {
			//m.Sessions = append(m.Sessions[:i], m.Sessions[i+1:])
		}
	}
}

func (m *messageService) BroadcastJSON(obj map[string]interface{}) {
	json, err := json.Marshal(obj)

	if err != nil {
		panic(err)
	}

	jsonString := string(json)
	m.Broadcast(jsonString)
}
