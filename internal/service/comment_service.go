// mockgen -source=E:\VSCodeFile/testovoeOzon/PostCommentService/internal/service/comment_service.go -destination=E:\VSCodeFile/testovoeOzon/PostCommentService/internal/service/mock/postgetter_mock.go -package=mock service PostGetter
package service

import (
	"database/sql"
	"errors"
	"postcommentservice/graph/model"
	constant "postcommentservice/internal/consts"
	"postcommentservice/internal/gateway"
)

type CommentService struct {
	repo       gateway.Comment
	PostGetter PostGetter
}

type PostGetter interface {
	GetPostById(id int) (model.Post, error)
}

func NewCommentService(repo gateway.Comment, getter PostGetter) *CommentService {
	return &CommentService{repo: repo, PostGetter: getter}
}

func (c CommentService) CreateComment(comment model.Comment) (model.Comment, error) {
	if len(comment.Author) == 0 {
		return model.Comment{}, errors.New(constant.EmptyAuthorError)
	}

	if len(comment.Content) >= constant.MaxContentLength {
		return model.Comment{}, errors.New(constant.LongContentError)
	}

	if comment.Post < 0 {
		return model.Comment{}, errors.New(constant.WrongIdError)
	}

	post, err := c.PostGetter.GetPostById(comment.Post)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Comment{}, errors.New(constant.PostNotFoundError)
		}
	}

	if !post.CommentsAllowed {
		return model.Comment{}, errors.New(constant.CommentsNotAllowedError)
	}

	newComment, err := c.repo.CreateComment(comment)

	if err != nil {
		return model.Comment{}, errors.New(constant.CreatingCommentError)
	}

	return newComment, nil
}

func (c CommentService) GetCommentsByPost(postId int, page *int, pageSize *int) ([]*model.Comment, error) {
	if postId < 0 {
		return nil, errors.New(constant.WrongIdError)
	}

	if page != nil && *page < 0 {
		return nil, errors.New(constant.WrongPageError)
	}

	if pageSize != nil && *pageSize < 0 {
		return nil, errors.New(constant.WrongPageSizeError)
	}

	var offset, limit int
	if page != nil && *page <= 0 {
		page = nil
	}

	if pageSize != nil && *pageSize < 0 {
		pageSize = nil
	}

	if page == nil || pageSize == nil {
		limit = -1
		offset = 0
	} else {
		offset = (*page - 1) * *pageSize
		limit = *pageSize
	}

	comments, err := c.repo.GetCommentsByPost(postId, limit, offset)

	if err != nil {
		return nil, errors.New(constant.GettingCommentError)
	}

	return comments, nil
}

func (c CommentService) GetRepliesOfComment(commentId int) ([]*model.Comment, error) {
	if commentId < 0 {
		return nil, errors.New(constant.WrongIdError)
	}

	comments, err := c.repo.GetRepliesOfComment(commentId)
	if err != nil {
		return nil, errors.New(constant.GettingCommentError)
	}

	return comments, nil
}
