package main

import (
	"fmt"

	"github.com/justinawrey/gochat/pkg/chat"
)

func main() {
	r := &chat.Room{Name: "test"}
	c1 := chat.NewChatter("justin")
	c2 := chat.NewChatter("carl")
	c3 := chat.NewChatter("connor")

	for _, c := range []*chat.Chatter{c1, c2, c3} {
		c.OnMsgReceive(func(m chat.Msg) {
			fmt.Println(m)
		})
		c.Join(r)
	}

	c1.Send("testing from justin")
	c2.Send("testing from carl")
	c3.Send("testing from connor")
}
