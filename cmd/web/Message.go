package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type MessageType int

const (
	Broadcast MessageType = iota
	Personal
	Control
)

func (mt *MessageType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		s = strings.TrimSpace(strings.ToLower(s))
		switch s {
		case "broadcast":
			*mt = Broadcast
			return nil
		case "personal":
			*mt = Personal
			return nil
		case "control":
			*mt = Control
			return nil
		default:
			return fmt.Errorf("invalid message type: %s", s)
		}
	}

	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		switch MessageType(i) {
		case Broadcast, Personal, Control:
			*mt = MessageType(i)
			return nil
		}
	}

	return fmt.Errorf("invalid message type value")
}

func (mt MessageType) MarshalJSON() ([]byte, error) {
	var s string
	switch mt {
	case Broadcast:
		s = "broadcast"
	case Personal:
		s = "personal"
	case Control:
		s = "control"
	default:
		return nil, fmt.Errorf("invalid message type: %d", mt)
	}
	return json.Marshal(s)
}

type ControlCommand int

const (
	Unknown ControlCommand = iota
	Play
	Pause
	Next
	Previous
	Jump
)

func (cc *ControlCommand) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		s = strings.TrimSpace(strings.ToLower(s))
		switch s {
		case "play":
			*cc = Play
			return nil
		case "pause":
			*cc = Pause
			return nil
		case "next":
			*cc = Next
			return nil
		case "previous":
			*cc = Previous
			return nil
		case "jump":
			*cc = Jump
			return nil
		default:
			return fmt.Errorf("invalid control command: %s", s)
		}
	}

	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		switch ControlCommand(i) {
		case Play, Pause, Next, Previous, Jump:
			*cc = ControlCommand(i)
			return nil
		}
	}

	return fmt.Errorf("invalid control command value")
}

func (cc ControlCommand) MarshalJSON() ([]byte, error) {
	var s string
	switch cc {
	case Play:
		s = "play"
	case Pause:
		s = "pause"
	case Next:
		s = "next"
	case Previous:
		s = "previous"
	case Jump:
		s = "jump"
	default:
		return nil, fmt.Errorf("invalid control command: %d", cc)
	}
	return json.Marshal(s)
}

type Payload struct {
	ChatMessage string         `json:"chatMessage,omitempty"`
	Command     ControlCommand `json:"command,omitempty"`
}

type Message struct {
	Sender   string      `json:"sender"`
	Receiver string      `json:"receiver,omitempty"`
	Type     MessageType `json:"type"`
	Payload  Payload     `json:"payload"`
}

func validatePayload(payload Payload, msgType MessageType) bool {
	if msgType == Control {
		return payload.Command != Unknown
	}

	// Broadcast and Personal require a chat message
	return strings.TrimSpace(payload.ChatMessage) != ""
}

func validateMessage(msg Message) bool {
	if strings.TrimSpace(msg.Sender) == "" {
		return false
	}

	if msg.Type != Broadcast && msg.Type != Personal && msg.Type != Control {
		return false
	}

	switch msg.Type {
	case Broadcast:
		if strings.TrimSpace(msg.Receiver) != "*" {
			return false
		}
	case Personal:
		if strings.TrimSpace(msg.Receiver) == "" || strings.TrimSpace(msg.Receiver) == "*" {
			return false
		}
	case Control:
		if strings.TrimSpace(msg.Receiver) == "" {
			return false
		}
	}

	if !validatePayload(msg.Payload, msg.Type) {
		return false
	}

	return true
}
