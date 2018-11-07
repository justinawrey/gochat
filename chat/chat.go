package chat

// Chatter is a chat room participant.  It can
// send messages to a room, and set callback functions
// which fire when a new message is received by someone
// else in the room.
type Chatter struct {
	name    string
	room    *Room
	receive chan string
	quit    chan struct{}
}

// Room is a chatroom that contains one or more Chatters.
// The main functionality of Room is to add / remove Chatters.
type Room struct {
	name     string
	chatters []*Chatter
}

// Send sends a message msg on behald of Chatter c.
// The message is sent to the ONE most recent room
// that c was added to.
func (c *Chatter) Send(msg string) {
	c.room.broadcast(msg)
}

func (c *Chatter) OnMsgReceive(f func(string)) {
	go func() {
		for {
			select {
			case <-c.quit:
				return
			case msg := <-c.receive:
				f(msg)
			}
		}
	}()
}

func (r *Room) Add(c *Chatter) {
	r.chatters = append(r.chatters, c)
	c.room = r
}

func (r *Room) Remove(c *Chatter) {
	for i, chatter := range r.chatters {
		if chatter.name == c.name {
			r.chatters = append(r.chatters[:i], r.chatters[i+1:]...)
			c.quit <- struct{}{}
			break
		}
	}
}

func (r *Room) broadcast(msg string) {
	for _, chatter := range r.chatters {
		chatter.receive <- msg
	}
}
