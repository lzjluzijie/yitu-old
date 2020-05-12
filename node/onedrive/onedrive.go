package onedrive

import (
	"fmt"
	"log"
	"time"

	"github.com/lzjluzijie/yitu/node"
)

type File struct {
	ID        uint64 `gorm:"primary_key"`
	TuID      uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	FolderID    string
	FileID      string
	RedirectURL string
}

type Node struct {
	*Config
}

func NewNode(c *Config) (n *Node) {
	n = &Node{Config: c}

	go func() {
		for {
			err := n.Refresh()
			if err != nil {
				log.Println(err.Error())
				time.Sleep(5 * time.Second)
			} else {
				time.Sleep(50 * time.Minute)
			}
		}
	}()
	return
}

func (n *Node) Handle(o *node.Object) (f *File) {
	p := fmt.Sprintf(`/yitu/%s/%s/%s`, time.Now().Format("20060102"), o.SHA256, o.Name)
	id, parent, err := n.Upload(p, o.Data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	url, err := n.Share(id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return &File{FolderID: parent, FileID: id, RedirectURL: url}
}
