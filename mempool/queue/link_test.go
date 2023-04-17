package queue

import (
	"fmt"
	"testing"
)

func TestLinkNode_AddAtTail(t *testing.T) {

	link := NewLink(0)
	link.AddAtTail(2)
	link.AddAtTail(3)
	link.AddAtTail(1)
	link.AddAtTail(44)
	link.AddAtTail("2")
	link.Print()
	dataNode := link.RemoveFirst()
	fmt.Println("dataNode ===> ", dataNode.Data)
	link.Print()
}
