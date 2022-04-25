package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	size int
	head *ListItem
	tail *ListItem
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	if l.size == 0 {
		return nil
	}
	return l.head
}

func (l *list) Back() *ListItem {
	if l.size == 0 {
		return nil
	}
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.head == nil {
		l.head = item
		l.tail = item
	} else {
		item.Next = l.Front()
		l.head.Prev = item
		l.head = item
	}
	l.size++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.head == nil {
		l.head = item
		l.tail = item
	} else {
		item.Prev = l.Back()
		l.tail.Next = item
		l.tail = item
	}
	l.size++
	return item
}

func (l *list) Remove(v *ListItem) {
	if v.Next != nil {
		v.Next.Prev = v.Prev
	}
	if v.Prev != nil {
		v.Prev.Next = v.Next
	}
	l.size--
}

func (l *list) MoveToFront(v *ListItem) {
	if l.head == v {
		return
	}
	if v.Next != nil {
		v.Next.Prev = v.Prev
	}
	if v.Prev != nil {
		v.Prev.Next = v.Next
	}
	v.Next = l.head
	l.head = v
}

func NewList() List {
	return new(list)
}
