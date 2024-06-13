// mockgen -source=E:\VSCodeFile/testovoeOzon/PostCommentService/internal/service/service.go -destination=E:\VSCodeFile/testovoeOzon/PostCommentService/internal/service/mock/service_mock.go -package=mock service Comment Post
package service

import (
	"postcommentservice/graph/model"
	"postcommentservice/internal/gateway"
)

type Service struct {
	Post
	Comment
}

func NewService(gateway *gateway.Gateway) *Service {
	return &Service{
		Post:    NewPostService(gateway.Post),
		Comment: NewCommentService(gateway.Comment, gateway.Post),
	}
}

type Post interface {
	CreatePost(post model.Post) (model.Post, error)
	GetPostById(id int) (model.Post, error)
	GetAllPosts(page, pageSize *int) ([]*model.Post, error)
}

type Comment interface {
	CreateComment(comment model.Comment) (model.Comment, error)
	GetCommentsByPost(postId int, page *int, pageSize *int) ([]*model.Comment, error)
	GetRepliesOfComment(commentId int) ([]*model.Comment, error)
}
