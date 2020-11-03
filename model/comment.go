package model

type Comment struct {
	CommentId int
	UserId int
	PicId int
	Content string
	CreatedAt string
	ReplyOf int
	Likes int
	Deleted bool
}