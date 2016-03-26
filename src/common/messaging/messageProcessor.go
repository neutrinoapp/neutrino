package messaging

import (
	"errors"
	"fmt"

	"time"

	"github.com/neutrinoapp/neutrino/src/common/client"
	"github.com/neutrinoapp/neutrino/src/common/db"
	"github.com/neutrinoapp/neutrino/src/common/log"
	"gopkg.in/redis.v3"
)

type MessageProcessor struct {
	dbService   db.DbService
	redisClient *redis.Client
}

func NewMessageProcessor() MessageProcessor {
	return MessageProcessor{
		dbService:   db.NewDbService(),
		redisClient: client.GetNewRedisClient(),
	}
}

func (p MessageProcessor) shouldProcessMessage(m Message) bool {
	key := m.GetRedisKey()
	timestamp := p.redisClient.Get(key).Val()
	if timestamp == "" {
		return true
	}

	cachedTime, parseTimeError := time.Parse(time.RFC3339, timestamp)
	if parseTimeError == nil {
		log.Error("Error parsing message time:", parseTimeError)
		return true
	}

	messageTime, parseTimeError := time.Parse(time.RFC3339, m.Timestamp)
	if parseTimeError == nil {
		log.Error("Error parsing message time:", parseTimeError)
		return true
	}

	return cachedTime.Before(messageTime)
}

func (p MessageProcessor) Process(m Message) (res interface{}, err error) {
	shouldProcess := p.shouldProcessMessage(m)
	if !shouldProcess {
		log.Info("Skipping old message:", m)
		return nil, nil
	}

	p.redisClient.Set(m.GetRedisKey(), m.Timestamp, 0)

	if m.Operation == OP_CREATE {
		return p.handleCreateMessage(m)
	} else if m.Operation == OP_READ {
		return p.handleReadMessage(m)
	} else if m.Operation == OP_UPDATE {
		err = p.handleUpdateMessage(m)
		return
	} else if m.Operation == OP_DELETE {
		err = p.handleDeleteMessage(m)
		return
	}

	return nil, errors.New(fmt.Sprintf("Invalid message! %v", m))
}

func (p MessageProcessor) handleCreateMessage(m Message) (id string, err error) {
	return p.dbService.CreateItem(m.App, m.Type, m.Payload)
}

func (p MessageProcessor) handleReadMessage(m Message) (data interface{}, err error) {
	if id, ok := m.Payload[db.ID_FIELD].(string); ok {
		return p.dbService.GetItemById(id)
	}

	return p.dbService.GetItems(m.App, m.Type, m.Options.Filter)
}

func (p MessageProcessor) handleUpdateMessage(m Message) (err error) {
	if id, ok := m.Payload[db.ID_FIELD].(string); ok {
		return p.dbService.UpdateItemById(id, m.Payload)
	}

	return errors.New(fmt.Sprintf("Cannot update item without id %v", m))
}

func (p MessageProcessor) handleDeleteMessage(m Message) (err error) {
	if id, ok := m.Payload[db.ID_FIELD].(string); ok {
		return p.dbService.DeleteItemById(id)
	}

	return p.dbService.DeleteAllItems(m.App, m.Type)
}
