package chat

import (
	"fmt"
)

// Msg is a chatroom message.  It implements fmt.Stringer,
// so it can be passed to the Print family of functions.
type Msg struct {
	// From is the sender of this message.
	From string

	// Room is the room within which this message was sent.
	Room *Room

	// Contents is the contents of the message.
	Contents string
}

// String implements fmt.Stringer.
func (m Msg) String() string {
	return fmt.Sprintf("%s > %s", m.From, m.Contents)
}

// Chatter is a chat room participant.  It can
// send messages to a room, and set callback functions
// which fire when a new message is received by someone
// else in the room.
type Chatter struct {
	// Name is the name of this chatter.
	Name string

	// Room is the most recent room that this chatter belongs to.
	Room *Room

	handler func(Msg)
	receive chan Msg
	quit    chan struct{}
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
	c.Room = r

	r.broadcast("admin", fmt.Sprintf("%s has entered the room", c.Name))
}

// Leave removes Chatter c from Room r.
func (c *Chatter) Leave(r *Room) {
	for i, chatter := range r.chatters {
		if chatter.Name == c.Name {
			r.chatters = append(r.chatters[:i], r.chatters[i+1:]...)
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
	c.handler = f
}

// Room is a chatroom that contains one or more Chatters.
// The main functionality of Room is to add / remove Chatters.
type Room struct {
	// Name is the name of the room.
	Name     string
	chatters []*Chatter
}

// NewRoom returns a new room with name name.
func NewRoom(name string) *Room {
	return &Room{
		Name: name,
	}
}

// Chatters returns the names of all Chatters in r.
func (r *Room) Chatters() []string {
	var res []string
	for _, c := range r.chatters {
		res = append(res, c.Name)
	}
	return res
}

// Add adds Chatter c to Room r.
// This is a convenience function and is semantically the same
// as c.Join(r).
func (r *Room) Add(c *Chatter) {
	c.Join(r)
}

// Remove removes Chatter c from Room r.
// This is a convenience function and is semantically the same
// as c.Leave(r).
func (r *Room) Remove(c *Chatter) {
	c.Leave(r)
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
	for _, c := range r.chatters {
		if c.Room != nil {
			c.receive <- Msg{From: from, Contents: m}
		}
	}
}
