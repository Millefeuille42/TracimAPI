package session

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type TLMHandler func(s *Session, message *TracimLiveMessage)

type TLMEvent struct {
	EventId   int         `json:"event_id"`
	EventType string      `json:"event_type"`
	Read      interface{} `json:"read"`
	Created   time.Time   `json:"created"`
	Fields    interface{} `json:"fields"`
}

type TracimLiveMessage struct {
	Event      string
	Data       string
	DataParsed TLMEvent
}

const (
	TLMConnected = "stream-open"
	TLMMessage   = "message"
	TLMError     = "error"
)

func (a *Session) hookToTLM(userId string) error {
	req, err := a.GenerateRequest("GET", fmt.Sprintf("/users/%s/live_messages", userId), []byte(""))
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

		a.eventChannel <- event
	}
}

func (a *Session) SendError(err error) {
	a.errorChannel <- err
}

func (a *Session) ListenEvents(userId string) {
	go func() {
		err := a.hookToTLM(userId)
		if err != nil {
			a.SendError(err)
		}
	}()

	for {
		select {
		case TLM := <-a.eventChannel:
			switch TLM.Event {
			case TLMConnected:
				if _, ok := a.eventHandler[TLMConnected]; ok {
					a.eventHandler[TLMConnected](a, &TLM)
				}
			case TLMMessage:
				err := json.Unmarshal([]byte(TLM.Data), &TLM.DataParsed)
				if err != nil {
					a.SendError(err)
					continue
				}
				if _, ok := a.eventHandler[TLMMessage]; ok {
					a.eventHandler[TLMMessage](a, &TLM)
				}
				if _, ok := a.eventHandler[TLM.DataParsed.EventType]; ok {
					a.eventHandler[TLM.DataParsed.EventType](a, &TLM)
				}
			}
		case err := <-a.errorChannel:
			if _, ok := a.eventHandler[TLMError]; ok {
				a.eventHandler[TLMError](a, &TracimLiveMessage{
					Event: TLMError,
					Data:  err.Error(),
				})
			}
		}
	}
}

func (a *Session) TLMSubscribe(event string, handler TLMHandler) {
	a.eventHandler[event] = handler
}
