package session

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// TLMHandler is the function type for the handlers for TLMs, to be used with TLMSubscribe
type TLMHandler func(s *Session, message *TracimLiveMessage)

// TLMEvent is a parsed TLM with the Fields element as an interface for generic usage (see Tracim doc about TLMs)
type TLMEvent struct {
	EventId   int         `json:"event_id"`
	EventType string      `json:"event_type"`
	Read      interface{} `json:"read"`
	Created   time.Time   `json:"created"`
	Fields    interface{} `json:"fields"`
}

// TracimLiveMessage is a wrapper for a TLM, it also contains the EventStream event type and data
type TracimLiveMessage struct {
	// Event EventStream event type
	Event string
	// Data raw EventStream data
	Data string
	// DataParsed Parsed as TLMEvent EventStream data
	DataParsed TLMEvent
}

const (
	// TLMConnected event type used when connected to the EventStream
	TLMConnected = "stream-open"
	// TLMMessage event type used when receiving a TLM from the EventStream
	TLMMessage = "message"
	// TLMError event type used when for the error handler
	TLMError = "error"
)

func (s *Session) hookToTLM() error {
	req, err := s.GenerateRequest("GET", fmt.Sprintf("/users/%s/live_messages", s.UserID), []byte(""))
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	bufferSize := 1024
	for {
		fullData := make([]byte, 0)
		for {
			data := make([]byte, 1024)
			bytesRead, err := resp.Body.Read(data)
			if err != nil {
				log.Fatal(err)
				return err
			}
			fullData = append(fullData, data[:bytesRead]...)
			if bytesRead < bufferSize {
				break
			}
		}
		dataString := string(fullData[:])
		eventIndex := strings.Index(dataString, "event:")
		if eventIndex < 0 {
			log.Println("No event descriptor")
			continue
		}
		splitMessage := strings.SplitN(dataString[eventIndex+len("event:"):], "\n", 2)
		if splitMessage == nil || len(splitMessage) < 2 {
			log.Println("Unable to split data")
			continue
		}

		event := TracimLiveMessage{
			Event: strings.TrimSpace(splitMessage[0]),
			Data:  strings.TrimSpace(splitMessage[1][len("data:"):]),
		}

		s.eventChannel <- event
	}
}

// SendError used to send error to custom error handler
func (s *Session) SendError(err error) {
	s.errorChannel <- err
}

// ListenEvents Start listening to events, it is preferred to register handlers beforehand
func (s *Session) ListenEvents() {
	go func() {
		err := s.hookToTLM()
		if err != nil {
			s.SendError(err)
		}
	}()

	for {
		select {
		case TLM := <-s.eventChannel:
			switch TLM.Event {
			case TLMConnected:
				if _, ok := s.eventHandler[TLMConnected]; ok {
					s.eventHandler[TLMConnected](s, &TLM)
				}
			case TLMMessage:
				err := json.Unmarshal([]byte(TLM.Data), &TLM.DataParsed)
				if err != nil {
					s.SendError(err)
					continue
				}
				if _, ok := s.eventHandler[TLMMessage]; ok {
					s.eventHandler[TLMMessage](s, &TLM)
				}
				if _, ok := s.eventHandler[TLM.DataParsed.EventType]; ok {
					s.eventHandler[TLM.DataParsed.EventType](s, &TLM)
				}
			}
		case err := <-s.errorChannel:
			if _, ok := s.eventHandler[TLMError]; ok {
				s.eventHandler[TLMError](s, &TracimLiveMessage{
					Event: TLMError,
					Data:  err.Error(),
				})
			}
		}
	}
}

// TLMSubscribe register a handler for a specific event
// use any of the TLM* constants for EventStream events or any of the TLM event types (see Tracim doc about TLMs)
func (s *Session) TLMSubscribe(event string, handler TLMHandler) {
	s.eventHandler[event] = handler
}
