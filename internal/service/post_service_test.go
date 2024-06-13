package service_test

import (
	"database/sql"
	"errors"
	"testing"

	"postcommentservice/graph/model"
	constant "postcommentservice/internal/consts"
	mock_gateway "postcommentservice/internal/gateway/mock"
	"postcommentservice/internal/service"

	"go.uber.org/mock/gomock"
)

func IsEqualPost(a, b model.Post) bool {
	return a.ID == b.ID && a.Author == b.Author && a.Content == b.Content &&
		a.CommentsAllowed == b.CommentsAllowed && a.Title == b.Title
}

func IsEqualPosts(a, b []*model.Post) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !IsEqualPost(*a[i], *b[i]) {
			return false
		}
	}
	return true
}

func TestPostsService_CreatePost(t *testing.T) {

	defaultPost := model.Post{
		ID:              1,
		Author:          "Author",
		Content:         "Cnt",
		CommentsAllowed: true,
		Title:           "Title",
	}

	tests := []struct {
		Title    string
		post     model.Post
		want     model.Post
		wantErr  bool
		repoRes  model.Post
		repoErr  error
		repoSkip bool
	}{
		{
			Title:    "Positive",
			post:     defaultPost,
			want:     defaultPost,
			wantErr:  false,
			repoErr:  nil,
			repoRes:  defaultPost,
			repoSkip: false,
		},
		{
			Title:    "Error from repo",
			post:     defaultPost,
			want:     model.Post{},
			wantErr:  true,
			repoErr:  errors.New("some error"),
			repoRes:  model.Post{},
			repoSkip: false,
		},
		{
			Title:    "Wrong author",
			post:     model.Post{Author: ""},
			want:     model.Post{},
			wantErr:  true,
			repoErr:  nil,
			repoRes:  model.Post{},
			repoSkip: true,
		},
		{
			Title:    "Wrong content",
			post:     model.Post{Author: "Au1", Content: string(make([]byte, constant.MaxContentLength+1))},
			want:     model.Post{},
			wantErr:  true,
			repoErr:  nil,
			repoRes:  model.Post{},
			repoSkip: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Title, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_gateway.NewMockPost(ctl)

			if !tt.repoSkip {
				repo.EXPECT().CreatePost(tt.post).Return(tt.repoRes, tt.repoErr)
			}

			p := service.NewPostService(repo)

			got, err := p.CreatePost(tt.post)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !IsEqualPost(got, tt.want) {
				t.Errorf("CreatePost() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostService_GetPostById(t *testing.T) {
	tests := []struct {
		Title    string
		postId   int
		want     model.Post
		wantErr  bool
		repoRes  model.Post
		repoErr  error
		repoSkip bool
	}{
		{
			Title:    "Positive",
			postId:   1,
			want:     model.Post{ID: 1, Title: "Title", Author: "Author", Content: "Content", CommentsAllowed: true},
			wantErr:  false,
			repoErr:  nil,
			repoRes:  model.Post{ID: 1, Title: "Title", Author: "Author", Content: "Content", CommentsAllowed: true},
			repoSkip: false,
		},
		{
			Title:    "Post not found",
			postId:   1,
			want:     model.Post{},
			wantErr:  true,
			repoErr:  sql.ErrNoRows,
			repoRes:  model.Post{},
			repoSkip: false,
		},
		{
			Title:    "Wrong ID",
			postId:   -1,
			want:     model.Post{},
			wantErr:  true,
			repoErr:  nil,
			repoRes:  model.Post{},
			repoSkip: true,
		},
		{
			Title:    "Error from repo",
			postId:   1,
			want:     model.Post{},
			wantErr:  true,
			repoErr:  errors.New("some error"),
			repoRes:  model.Post{},
			repoSkip: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Title, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_gateway.NewMockPost(ctl)

			if !tt.repoSkip {
				repo.EXPECT().GetPostById(tt.postId).Return(tt.repoRes, tt.repoErr)
			}

			p := service.NewPostService(repo)

			got, err := p.GetPostById(tt.postId)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPostById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !IsEqualPost(got, tt.want) {
				t.Errorf("GetPostById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostService_GetAllPosts(t *testing.T) {
	defaultPost := model.Post{ID: 1, Title: "Title", Author: "Author", Content: "Content", CommentsAllowed: true}
	tests := []struct {
		Title    string
		page     *int
		pageSize *int
		want     []*model.Post
		wantErr  bool
		repoRes  []*model.Post
		repoErr  error
		repoSkip bool
	}{
		{
			Title:    "Positive",
			page:     nil,
			pageSize: nil,
			want:     []*model.Post{&defaultPost},
			wantErr:  false,
			repoErr:  nil,
			repoRes:  []*model.Post{&defaultPost},
			repoSkip: false,
		},
		{
			Title:    "Error from repo",
			page:     nil,
			pageSize: nil,
			want:     nil,
			wantErr:  true,
			repoErr:  errors.New("some error"),
			repoRes:  nil,
			repoSkip: false,
		},
		{
			Title:    "Wrong page",
			page:     &[]int{-1}[0],
			pageSize: nil,
			want:     nil,
			wantErr:  true,
			repoErr:  nil,
			repoRes:  nil,
			repoSkip: true,
		},
		{
			Title:    "Wrong pageSize",
			page:     nil,
			pageSize: &[]int{-1}[0],
			want:     nil,
			wantErr:  true,
			repoErr:  nil,
			repoRes:  nil,
			repoSkip: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Title, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_gateway.NewMockPost(ctl)

			if !tt.repoSkip {
				repo.EXPECT().GetAllPosts(gomock.Any(), gomock.Any()).Return(tt.repoRes, tt.repoErr)
			}

			p := service.NewPostService(repo)

			got, err := p.GetAllPosts(tt.page, tt.pageSize)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !IsEqualPosts(got, tt.want) {
				t.Errorf("GetAllPosts() got = %v, want %v", got, tt.want)
			}
		})
	}
}
