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
	Handler func(Msg)
	Room    *Room

	receive chan Msg
	quit    chan struct{}
}

// Room is a chatroom that contains one or more Chatters.
// The main functionality of Room is to add / remove Chatters.
type Room struct {
	Name     string
	Chatters []*Chatter
}

// Msg is a chatroom message.  It implements fmt.Stringer,
// so it can be passed to the Print family of functions.
type Msg struct {
	From     string
	Contents string
}

// String implements fmt.Stringer.
func (m Msg) String() string {
	return fmt.Sprintf("%s > %s", m.From, m.Contents)
}

// NewChatter creates a new Chatter with specified name.
func NewChatter(name string) *Chatter {
	c := &Chatter{
		Name: name,
		// The default chat message handler does nothing.
		Handler: func(Msg) { return },
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
				c.Handler(m)
			}
		}
	}()

	return c
}

// Join adds Chatter c to Room r.
func (c *Chatter) Join(r *Room) {
	r.Chatters = append(r.Chatters, c)
	c.Room = r

	r.broadcast("admin", fmt.Sprintf("%s has entered the room", c.Name))
}

// Leave removes Chatter c from Room r.
func (c *Chatter) Leave(r *Room) {
	for i, chatter := range r.Chatters {
		if chatter.Name == c.Name {
			r.Chatters = append(r.Chatters[:i], r.Chatters[i+1:]...)
			c.Room = nil
			close(c.quit)

			r.broadcast("admin", fmt.Sprintf("%s has left the room", c.Name))
			break
		}
	}
}

// Send sends a message m on behald of Chatter c.
// The message is sent to the ONE most recent room
// that c was added to.
func (c *Chatter) Send(m string) {
	c.Room.broadcast(c.Name, m)
}

// OnMsgReceive executes method f (which takes a single msg parameter)
// every time c receives a message from its current room.
func (c *Chatter) OnMsgReceive(f func(Msg)) {
	c.Handler = f
}

// Flush flushes out all lingering messages
// that still need to be received by every chatter in Room r.
// Chatter.OnMsgReceive will be fired, as usual, for all
// flushed messages.
func (r *Room) Flush() {
	//TODO:
	return
}

// Close flushes all lingering messages, and
// then closes room r, removing all Chatters in r in the process.
func (r *Room) Close() {
	//TODO:
	return
}

// broadcast broadcasts a message msg to all chatters
// in the room r.
func (r *Room) broadcast(from string, m string) {
	for _, c := range r.Chatters {
		if c.Room != nil {
			c.receive <- Msg{From: from, Contents: m}
		}
	}
}
