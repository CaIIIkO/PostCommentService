package service

import (
	"database/sql"
	"errors"
	"postcommentservice/graph/model"

	constant "postcommentservice/internal/consts"
	"postcommentservice/internal/gateway"
)

type PostService struct {
	repo gateway.Post
}

func NewPostService(repo gateway.Post) *PostService {
	return &PostService{repo: repo}
}

func (p PostService) CreatePost(post model.Post) (model.Post, error) {

	if len(post.Author) == 0 {
		return model.Post{}, errors.New(constant.EmptyAuthorError)
	}

	if len(post.Content) >= constant.MaxContentLength {
		return model.Post{}, errors.New(constant.LongContentError)
	}

	newPost, err := p.repo.CreatePost(post)
	if err != nil {
		return model.Post{}, errors.New(constant.CreatingPostError)
	}

	return newPost, nil
}

func (p PostService) GetPostById(postId int) (model.Post, error) {

	if postId < 0 {
		return model.Post{}, errors.New(constant.WrongIdError)
	}

	post, err := p.repo.GetPostById(postId)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return model.Post{}, errors.New(constant.PostNotFoundError)
		}

		return model.Post{}, errors.New(constant.GettingPostError)
	}

	return post, nil
}

func (p PostService) GetAllPosts(page, pageSize *int) ([]*model.Post, error) {

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

	posts, err := p.repo.GetAllPosts(limit, offset)
	if err != nil {
		return nil, errors.New(constant.GettingPostError)
	}

	return posts, nil
}
