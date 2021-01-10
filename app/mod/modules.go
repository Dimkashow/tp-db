package mod

import (
	"time"
)

type User struct {
	ID       uint64 `json:"-"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
}

type Forum struct {
	ID            uint64 `json:"-"`
	PostsCount    uint64 `json:"posts"`
	Slug          string `json:"slug"`
	ThreadsCount  uint64 `json:"threads"`
	Title         string `json:"title"`
	AdminNickname string `json:"user"`
	AdminID       uint64 `json:"-"`
}

type Thread struct {
	Author       string    `json:"author"`
	AuthorID     uint64    `json:"-"`
	CreationDate time.Time `json:"created"`
	Forum        string    `json:"forum"`
	ForumID      uint64    `json:"-"`
	ID           uint64    `json:"id,omitempty"`
	About        string    `json:"message"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	Votes        int64     `json:"votes"`
}

//easyjson:json
type EzPost []*Post

type Post struct {
	Author       string    `json:"author"`
	CreationDate time.Time `json:"created"`
	Forum        string    `json:"forum"`
	ID           uint64    `json:"id"`
	IsEdited     bool      `json:"isEdited"`
	Message      string    `json:"message"`
	ParentID     uint64    `json:"parent"`
	ThreadID     uint64    `json:"thread"`
	ForumID      uint64    `json:"-"`
	AuthorID     uint64    `json:"-"`
}

type Vote struct {
	ID 		 uint64 `json:"-"`
	Nickname string `json:"nickname"`
	Voice    int64  `json:"voice"`
	ThreadID uint64 `json:"-"`
	UserID   uint64 `json:"-"`
}

type Mes struct {
	Message string `json:"message"`
}

type PostFull struct {
	Author   *User   `json:"author"`
	Forum    *Forum  `json:"forum"`
	PostData *Post   `json:"post"`
	Thread   *Thread `json:"thread"`
}

type Status struct {
	ForumsCount  uint64 `json:"forum"`
	PostsCount   uint64 `json:"post"`
	ThreadsCount uint64 `json:"thread"`
	UsersCount   uint64 `json:"user"`
}
