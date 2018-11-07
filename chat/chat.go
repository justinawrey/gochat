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
	return fmt.Sprintf("%s > %s\n", m.from, m.contents)
}

// NewChatter creates a new Chatter with specified name.
func NewChatter(name string) *Chatter {
	return &Chatter{
		Name:    name,
		receive: make(chan Msg),
		quit:    make(chan struct{}),
	}
}

// Send sends a message msg on behald of Chatter c.
// The message is sent to the ONE most recent room
// that c was added to.
func (c *Chatter) Send(m string) {
	c.room.broadcast(c.Name, m)
}

// OnMsgReceive executes method f (which takes a single msg parameter)
// every time c receives a message from its current room.
func (c *Chatter) OnMsgReceive(f func(Msg)) {
	go func() {
		for {
			select {
			case <-c.quit:
				return
			case m := <-c.receive:
				f(m)
			}
		}
	}()
}

// Add adds a Chatter c to the room r.
func (r *Room) Add(c *Chatter) {
	r.chatters = append(r.chatters, c)
	c.room = r

	r.broadcast("admin", fmt.Sprintf("%s has entered the room\n", c.Name))
}

// Remove removes a Chatter c from the room r.
func (r *Room) Remove(c *Chatter) {
	for i, chatter := range r.chatters {
		if chatter.Name == c.Name {
			r.chatters = append(r.chatters[:i], r.chatters[i+1:]...)
			c.quit <- struct{}{}

			r.broadcast("admin", fmt.Sprintf("%s has left the room\n", c.Name))
			break
		}
	}
}

// INTERNAL broadcast broadcasts a message msg to all chatters
// in the room r.
func (r *Room) broadcast(from string, m string) {
	for _, chatter := range r.chatters {
		chatter.receive <- Msg{from: from, contents: m}
	}
}
