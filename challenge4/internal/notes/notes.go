package notes

import (
	"fmt"
	"sync"
)

type Note struct {
	Id       int    `form:"id"`
	Title    string `form:"title"`
	Contents string `form:"contents"`
}

func (n Note) String() string {
	return fmt.Sprintf("%s [%d]\n%s", n.Title, n.Id, n.Contents)
}

func NewNote(title, contents string) *Note {
	CUR_ID++
	return &Note{CUR_ID, title, contents}
}

var (
	NOTES  = make(map[int]*Note)
	CUR_ID = 0
	LOCK   sync.Mutex
)
