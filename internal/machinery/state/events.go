/*
Copyright 2020 The routerd Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package state

// Event holds the old and new state of an object.
// An event is emitted for every persisting state change.
type Event struct {
	Old Object `json:"old"`
	New Object `json:"new"`
}

func (e Event) Type() EventType {
	if e.New == nil {
		return Deleted
	}
	if e.Old == nil {
		return Added
	}
	return Modified
}

type EventType string

const (
	Added    EventType = "Added"
	Modified EventType = "Modified"
	Deleted  EventType = "Deleted "
)

type eventHub struct {
	broadcast chan Event

	clients  map[*eventClient]struct{}
	register chan *eventClient
}

func newEventHub() *eventHub {
	return &eventHub{
		broadcast: make(chan Event),

		clients:  map[*eventClient]struct{}{},
		register: make(chan *eventClient),
	}
}

// Broadcast emits an event to all connected clients.
func (h *eventHub) Broadcast(old, new Object) {
	e := Event{
		Old: old, New: new,
	}

	h.broadcast <- e
}

// Register returns a new client connected to the event hub.
func (h *eventHub) Register() *eventClient {
	c := newEventClient()
	h.register <- c
	return c
}

func (h *eventHub) Run(stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			// close all clients
			for c := range h.clients {
				c.Close()
				delete(h.clients, c)
			}
			return

		case c := <-h.register:
			h.clients[c] = struct{}{}

		case message := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.events() <- message:
				default:
					// can't send -> Close()
					c.Close()
					delete(h.clients, c)
				}
			}
		}
	}
}

type eventClient struct {
	recv chan Event
}

func newEventClient() *eventClient {
	return &eventClient{
		recv: make(chan Event, 100),
	}
}

func (c *eventClient) Close() error {
	close(c.recv)
	return nil
}

func (c *eventClient) events() chan<- Event {
	return c.recv
}

func (c *eventClient) ResultChan() <-chan Event {
	return c.recv
}
