package in_memory

import (
	"postcommentservice/graph/model"
	"sync"
	"time"
)

type CommentsMemory struct {
	mu       sync.RWMutex
	comments map[int]model.Comment
	nextID   int
}

func NewCommentsMemory() *CommentsMemory {
	return &CommentsMemory{
		comments: make(map[int]model.Comment),
		nextID:   1,
	}
}

func (c *CommentsMemory) CreateComment(comment model.Comment) (model.Comment, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	comment.ID = c.nextID
	c.nextID++
	comment.CreatedAt = time.Now()

	c.comments[comment.ID] = comment

	if comment.ReplyTo != nil {
		parentComment := c.comments[*comment.ReplyTo]
		parentComment.Answers++
		c.comments[*comment.ReplyTo] = parentComment
	}

	return comment, nil
}

func (c *CommentsMemory) GetCommentsByPost(postId, limit, offset int) ([]*model.Comment, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var comments []*model.Comment
	count := 0

	for _, comment := range c.comments {
		if comment.Post == postId && comment.ReplyTo == nil {
			if count >= offset {
				comments = append(comments, &comment)
			}
			count++
			if limit > 0 && len(comments) >= limit {
				break
			}
		}
	}

	return comments, nil
}

func (c *CommentsMemory) GetRepliesOfComment(commentId int) ([]*model.Comment, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var comments []*model.Comment

	for _, comment := range c.comments {
		if comment.ReplyTo != nil && *comment.ReplyTo == commentId {
			comments = append(comments, &comment)
		}
	}

	return comments, nil
}
