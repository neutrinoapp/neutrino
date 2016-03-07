package models

import "gopkg.in/jcelliott/turnpike.v2"

type SubscribeOptions struct {
	Filter    JSON        `json:"filter"`
	BaseTopic string      `json:"baseTopic"`
	Topic     string      `json:"topic"`
	ClientId  turnpike.ID `json:"clientId"`
	TopicId   string      `json:"topicId"`
	Type      string      `json:"type"`
	AppId     string      `json:"appId"`
}

func (opts SubscribeOptions) IsSpecial() bool {
	return opts.Filter != nil
}
