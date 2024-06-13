// mockgen -source=E:\VSCodeFile/testovoeOzon/PostCommentService/internal/gateway/gateway.go -destination=E:\VSCodeFile/testovoeOzon/PostCommentService/internal/gateway/mock/gateway_mock.go -package=mock
package gateway

import "postcommentservice/graph/model"

type Gateway struct {
	Post
	Comment
}

func NewGateway(post Post, comment Comment) *Gateway {
	return &Gateway{
		Post:    post,
		Comment: comment,
	}
}

type Post interface {
	CreatePost(post model.Post) (model.Post, error)
	GetPostById(id int) (model.Post, error)
	GetAllPosts(limit, offset int) ([]*model.Post, error)
}

type Comment interface {
	CreateComment(comment model.Comment) (model.Comment, error)
	GetCommentsByPost(postId, limit, offset int) ([]*model.Comment, error)
	GetRepliesOfComment(commentId int) ([]*model.Comment, error)
}
