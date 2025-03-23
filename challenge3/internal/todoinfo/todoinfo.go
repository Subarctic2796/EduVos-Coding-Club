package todoinfo

import "sync"

type Todo struct {
	Name string
	Id   int
}

func NewTodo(name string) *Todo {
	return &Todo{name, CUR_ID}
}

var (
	TODOS  = make(map[int]*Todo)
	CUR_ID = -1
	LOCK   sync.Mutex
)
