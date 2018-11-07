package chat

import (
	"fmt"
)

// Chatter is a chat room participant.  It can
// send messages to a room, and set callback functions
// which fire when a new message is received by someone
// else in the room.
type Chatter struct {
	Name    string
	handler func(Msg)
	room    *Room
	receive chan Msg
	quit    chan struct{}
}

// Room is a chatroom that contains one or more Chatters.
// The main functionality of Room is to add / remove Chatters.
type Room struct {
	Name     string
	chatters []*Chatter
}

// Msg is a chatroom message.  It implements fmt.Stringer,
// so it can be passed to the Print family of functions.
type Msg struct {
	from     string
	contents string
}

// String implements fmt.Stringer.
func (m Msg) String() string {
	return fmt.Sprintf("%s > %s", m.from, m.contents)
}

// NewChatter creates a new Chatter with specified name.
func NewChatter(name string) *Chatter {
	c := &Chatter{
		Name: name,
		// The default chat message handler does nothing.
		handler: func(Msg) { return },
		receive: make(chan Msg),
		quit:    make(chan struct{}),
	}

	go func() {
		defer close(c.receive)

		for {
			select {
			case <-c.quit:
				return
			case m := <-c.receive:
				c.handler(m)
			}
		}
	}()

	return c
}

// Join adds Chatter c to Room r.
func (c *Chatter) Join(r *Room) {
	r.chatters = append(r.chatters, c)
	c.room = r

	r.broadcast("admin", fmt.Sprintf("%s has entered the room", c.Name))
}

// Leave removes Chatter c from Room r.
func (c *Chatter) Leave(r *Room) {
	for i, chatter := range r.chatters {
		if chatter.Name == c.Name {
			r.chatters = append(r.chatters[:i], r.chatters[i+1:]...)
			c.room = nil
			close(c.quit)

			r.broadcast("admin", fmt.Sprintf("%s has left the room", c.Name))
			break
		}
	}
}

// Send sends a mesrage msg on behald of Chatter c.
// The message is sent to the ONE most recent room
// that c was added to.
func (c *Chatter) Send(m string) {
	c.room.broadcast(c.Name, m)
}

// OnMsgReceive executes method f (which takes a single msg parameter)
// every time c receives a message from its current room.
func (c *Chatter) OnMsgReceive(f func(Msg)) {
	c.handler = f
}

// broadcast broadcasts a message msg to all chatters
// in the room r.
func (r *Room) broadcast(from string, m string) {
	for _, c := range r.chatters {
		if c.room != nil {
			c.receive <- Msg{from: from, contents: m}
		}
	}
}

// TODO: room FLUSH
// TODO: room LOGS
