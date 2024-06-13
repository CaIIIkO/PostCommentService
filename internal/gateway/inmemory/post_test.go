package in_memory_test

import (
	"postcommentservice/graph/model"
	in_memory "postcommentservice/internal/gateway/inmemory"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	postMemory := in_memory.NewPostMemory()

	post := model.Post{
		Title:           "Test Title",
		Content:         "Test Content",
		Author:          "Test Author",
		CommentsAllowed: true,
	}

	createdPost, err := postMemory.CreatePost(post)
	assert.NoError(t, err)
	assert.NotNil(t, createdPost.ID)
	assert.Equal(t, post.Title, createdPost.Title)
	assert.Equal(t, post.Content, createdPost.Content)
	assert.Equal(t, post.Author, createdPost.Author)
	assert.Equal(t, post.CommentsAllowed, createdPost.CommentsAllowed)
	assert.WithinDuration(t, time.Now(), createdPost.CreatedAt, 1*time.Second)
}

func TestGetPostById(t *testing.T) {
	postMemory := in_memory.NewPostMemory()

	posts := []*model.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1", Author: "Author 1", CreatedAt: time.Now(), CommentsAllowed: true},
		{ID: 2, Title: "Post 2", Content: "Content 2", Author: "Author 2", CreatedAt: time.Now(), CommentsAllowed: false},
		{ID: 3, Title: "Post 3", Content: "Content 3", Author: "Author 3", CreatedAt: time.Now(), CommentsAllowed: true},
	}

	for _, post := range posts {
		_, err := postMemory.CreatePost(*post)
		assert.NoError(t, err)
	}

	result, err := postMemory.GetPostById(2)
	assert.NoError(t, err)
	assert.NotNil(t, result.ID)
	assert.Equal(t, posts[1].Title, result.Title)
	assert.Equal(t, posts[1].Content, result.Content)
	assert.Equal(t, posts[1].Author, result.Author)
	assert.Equal(t, posts[1].CommentsAllowed, result.CommentsAllowed)
	assert.Equal(t, posts[1].CreatedAt.Unix(), result.CreatedAt.Unix())
}

func TestGetAllPosts(t *testing.T) {
	postMemory := in_memory.NewPostMemory()

	posts := []*model.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1", Author: "Author 1", CreatedAt: time.Now(), CommentsAllowed: true},
		{ID: 2, Title: "Post 2", Content: "Content 2", Author: "Author 2", CreatedAt: time.Now(), CommentsAllowed: false},
		{ID: 3, Title: "Post 3", Content: "Content 3", Author: "Author 3", CreatedAt: time.Now(), CommentsAllowed: true},
	}

	for _, post := range posts {
		_, err := postMemory.CreatePost(*post)
		assert.NoError(t, err)
	}

	allPosts, err := postMemory.GetAllPosts(0, 0)
	assert.NoError(t, err)
	assert.Len(t, allPosts, 3)

	limitPosts, err := postMemory.GetAllPosts(2, 0)
	assert.NoError(t, err)
	assert.Len(t, limitPosts, 2)

	offsetPosts, err := postMemory.GetAllPosts(0, 2)
	assert.NoError(t, err)
	assert.Len(t, offsetPosts, 1)
	assert.Equal(t, posts[2].ID, offsetPosts[0].ID)
}
