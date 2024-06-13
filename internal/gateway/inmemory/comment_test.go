package in_memory_test

import (
	"postcommentservice/graph/model"
	in_memory "postcommentservice/internal/gateway/inmemory"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	commentsMemory := in_memory.NewCommentsMemory()

	comment := model.Comment{
		Post:    1,
		Content: "Test Comment",
	}

	createdComment, err := commentsMemory.CreateComment(comment)
	assert.NoError(t, err)
	assert.NotNil(t, createdComment.ID)
	assert.Equal(t, comment.Post, createdComment.Post)
	assert.Equal(t, comment.Content, createdComment.Content)
	assert.WithinDuration(t, time.Now(), createdComment.CreatedAt, 1*time.Second)
}

func TestGetCommentsByPost(t *testing.T) {
	commentsMemory := in_memory.NewCommentsMemory()

	comments := []*model.Comment{
		{ID: 1, Post: 1, Content: "Comment 1", CreatedAt: time.Now()},
		{ID: 2, Post: 1, Content: "Comment 2", CreatedAt: time.Now()},
		{ID: 3, Post: 2, Content: "Comment 3", CreatedAt: time.Now()},
	}

	for _, comment := range comments {
		_, err := commentsMemory.CreateComment(*comment)
		assert.NoError(t, err)
	}

	result, err := commentsMemory.GetCommentsByPost(1, 0, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetRepliesOfComment(t *testing.T) {
	commentsMemory := in_memory.NewCommentsMemory()

	comments := []*model.Comment{
		{ID: 1, Post: 1, Content: "Parent Comment", CreatedAt: time.Now()},
		{ID: 2, Post: 1, Content: "Reply to Comment 1", ReplyTo: ptrInt(1), CreatedAt: time.Now()},
		{ID: 3, Post: 1, Content: "Reply to Comment 1", ReplyTo: ptrInt(1), CreatedAt: time.Now()},
		{ID: 4, Post: 1, Content: "Reply to Comment 2", ReplyTo: ptrInt(2), CreatedAt: time.Now()},
	}

	for _, comment := range comments {
		_, err := commentsMemory.CreateComment(*comment)
		assert.NoError(t, err)
	}

	result, err := commentsMemory.GetRepliesOfComment(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2) // Ожидаем два ответа на комментарий с ID=1
}

func ptrInt(i int) *int {
	return &i
}
