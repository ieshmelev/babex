package babex

import (
	"encoding/json"
)

type RawMessage interface {
	Ack(multiple bool) error
	Nack(multiple bool) error
}

type Message struct {
	Key     string
	Chain   Chain
	Data    json.RawMessage
	Headers map[string]interface{}
	Config  []byte
	Meta    map[string]string

	InitialMessage *InitialMessage
	RawMessage     RawMessage
}

type InitialMessage struct {
	Chain  Chain             `json:"chain"`
	Data   json.RawMessage   `json:"data"`
	Config json.RawMessage   `json:"config"`
	Meta   map[string]string `json:"meta"`
}

func (m Message) Ack(multiple bool) error {
	return m.RawMessage.Ack(multiple)
}

func (m Message) Nack(multiple bool) error {
	return m.RawMessage.Nack(multiple)
}
