package service_test

import (
	"database/sql"
	"errors"
	"testing"

	"postcommentservice/graph/model"
	constant "postcommentservice/internal/consts"
	mock_gateway "postcommentservice/internal/gateway/mock"
	"postcommentservice/internal/service"
	mock_service "postcommentservice/internal/service/mock"

	"go.uber.org/mock/gomock"
)

func IsEqualComment(a, b model.Comment) bool {
	return a.ID == b.ID && a.Author == b.Author && a.Content == b.Content && a.Post == b.Post
}

func IsEqualComments(a, b []*model.Comment) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !IsEqualComment(*a[i], *b[i]) {
			return false
		}
	}
	return true
}

func TestCommentService_CreateComment(t *testing.T) {
	defaultComment := model.Comment{
		ID:      1,
		Author:  "Author",
		Content: "Content",
		Post:    1,
	}
	defaultPost := model.Post{
		ID:              1,
		Author:          "Author",
		Content:         "Content",
		CommentsAllowed: true,
		Title:           "Title",
	}

	tests := []struct {
		Title    string
		comment  model.Comment
		want     model.Comment
		wantErr  bool
		repoRes  model.Comment
		repoErr  error
		repoSkip bool
		postRes  model.Post
		postErr  error
		postSkip bool
	}{
		{
			Title:    "Positive",
			comment:  defaultComment,
			want:     defaultComment,
			wantErr:  false,
			repoErr:  nil,
			repoRes:  defaultComment,
			repoSkip: false,
			postRes:  defaultPost,
			postErr:  nil,
			postSkip: false,
		},
		{
			Title:    "Error from repo",
			comment:  defaultComment,
			want:     model.Comment{},
			wantErr:  true,
			repoErr:  errors.New("some error"),
			repoRes:  model.Comment{},
			repoSkip: false,
			postRes:  defaultPost,
			postErr:  nil,
			postSkip: false,
		},
		{
			Title:    "Wrong author",
			comment:  model.Comment{Author: ""},
			want:     model.Comment{},
			wantErr:  true,
			repoErr:  nil,
			repoRes:  model.Comment{},
			repoSkip: true,
			postRes:  defaultPost,
			postErr:  nil,
			postSkip: true,
		},
		{
			Title:    "Wrong content",
			comment:  model.Comment{Author: "Au1", Content: string(make([]byte, constant.MaxContentLength+1))},
			want:     model.Comment{},
			wantErr:  true,
			repoErr:  nil,
			repoRes:  model.Comment{},
			repoSkip: true,
			postRes:  defaultPost,
			postErr:  nil,
			postSkip: true,
		},
		{
			Title:    "Post not found",
			comment:  defaultComment,
			want:     model.Comment{},
			wantErr:  true,
			repoErr:  nil,
			repoRes:  model.Comment{},
			repoSkip: true,
			postRes:  model.Post{},
			postErr:  sql.ErrNoRows,
			postSkip: false,
		},
		{
			Title:    "Comments not allowed",
			comment:  defaultComment,
			want:     model.Comment{},
			wantErr:  true,
			repoErr:  nil,
			repoRes:  model.Comment{},
			repoSkip: true,
			postRes:  model.Post{ID: 1, CommentsAllowed: false},
			postErr:  nil,
			postSkip: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Title, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_gateway.NewMockComment(ctl)
			getter := mock_service.NewMockPostGetter(ctl)

			if !tt.repoSkip {
				repo.EXPECT().CreateComment(tt.comment).Return(tt.repoRes, tt.repoErr)
			}

			if !tt.postSkip {
				getter.EXPECT().GetPostById(tt.comment.Post).Return(tt.postRes, tt.postErr)
			}

			c := service.NewCommentService(repo, getter)

			got, err := c.CreateComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !IsEqualComment(got, tt.want) {
				t.Errorf("CreateComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentService_GetCommentsByPost(t *testing.T) {
	defaultComment := model.Comment{ID: 1, Author: "Author", Content: "Content", Post: 1}

	tests := []struct {
		Title    string
		postId   int
		page     *int
		pageSize *int
		want     []*model.Comment
		wantErr  bool
		repoRes  []*model.Comment
		repoErr  error
		repoSkip bool
	}{
		{
			Title:    "Positive",
			postId:   1,
			page:     nil,
			pageSize: nil,
			want:     []*model.Comment{&defaultComment},
			wantErr:  false,
			repoErr:  nil,
			repoRes:  []*model.Comment{&defaultComment},
			repoSkip: false,
		},
		{
			Title:    "Error from repo",
			postId:   1,
			page:     nil,
			pageSize: nil,
			want:     nil,
			wantErr:  true,
			repoErr:  errors.New("some error"),
			repoRes:  nil,
			repoSkip: false,
		},
		{
			Title:    "Wrong post ID",
			postId:   -1,
			page:     nil,
			pageSize: nil,
			want:     nil,
			wantErr:  true,
			repoErr:  nil,
			repoRes:  nil,
			repoSkip: true,
		},
		{
			Title:    "Wrong page",
			postId:   1,
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
			postId:   1,
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

			repo := mock_gateway.NewMockComment(ctl)

			if !tt.repoSkip {
				repo.EXPECT().GetCommentsByPost(tt.postId, gomock.Any(), gomock.Any()).Return(tt.repoRes, tt.repoErr)
			}

			c := service.NewCommentService(repo, nil)

			got, err := c.GetCommentsByPost(tt.postId, tt.page, tt.pageSize)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommentsByPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !IsEqualComments(got, tt.want) {
				t.Errorf("GetCommentsByPost() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentService_GetRepliesOfComment(t *testing.T) {
	defaultComment := model.Comment{ID: 1, Author: "Author", Content: "Content", Post: 1}

	tests := []struct {
		Title     string
		commentId int
		want      []*model.Comment
		wantErr   bool
		repoRes   []*model.Comment
		repoErr   error
		repoSkip  bool
	}{
		{
			Title:     "Positive",
			commentId: 1,
			want:      []*model.Comment{&defaultComment},
			wantErr:   false,
			repoErr:   nil,
			repoRes:   []*model.Comment{&defaultComment},
			repoSkip:  false,
		},
		{
			Title:     "Error from repo",
			commentId: 1,
			want:      nil,
			wantErr:   true,
			repoErr:   errors.New("some error"),
			repoRes:   nil,
			repoSkip:  false,
		},
		{
			Title:     "Wrong comment ID",
			commentId: -1,
			want:      nil,
			wantErr:   true,
			repoErr:   nil,
			repoRes:   nil,
			repoSkip:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Title, func(t *testing.T) {

			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_gateway.NewMockComment(ctl)

			if !tt.repoSkip {
				repo.EXPECT().GetRepliesOfComment(tt.commentId).Return(tt.repoRes, tt.repoErr)
			}

			c := service.NewCommentService(repo, nil)

			got, err := c.GetRepliesOfComment(tt.commentId)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRepliesOfComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !IsEqualComments(got, tt.want) {
				t.Errorf("GetRepliesOfComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
