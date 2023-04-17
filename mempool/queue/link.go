package queue

import "fmt"

type LinkNode struct {
	Data interface{} //data field ,use the interface ,can storage any type
	Next *LinkNode   //pointer
}

func NewLink(data interface{}) *LinkNode {
	return &LinkNode{Data: data}
}

func (head *LinkNode) AddAtTail(data interface{}) {
	if head == nil {
		return
	}
	temp := head
	for temp.Next != nil {
		temp = temp.Next
	}
	if temp.Next == nil {
		newNode := &LinkNode{Data: data}
		temp.Next = newNode
	}
}

func (head *LinkNode) FillData(dataList []interface{}) {
	if head == nil {
		return
	}
	for _, item := range dataList {
		head.AddAtTail(item)
	}
}

func (head *LinkNode) RemoveFirst() *LinkNode {
	if head == nil {
		return nil
	}
	data := head.Next
	temp := head.Next
	*head = *temp
	return &LinkNode{Data: data}
}

// Print the link the linkTable
func (head *LinkNode) Print() {
	if head == nil {
		return
	}
	temp := head
	for temp.Next != nil {
		fmt.Println(temp.Data)
		temp = temp.Next
	}
	fmt.Println(temp.Data, "\n----")
}
