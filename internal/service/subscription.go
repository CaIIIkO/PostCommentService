package service

import (
	"errors"
	"fmt"
	"postcommentservice/graph/model"
	constant "postcommentservice/internal/consts"
	"sync"
)

type CommentsSubscriptions struct {
	Chans   map[int][]CommentSubscription
	counter int
	mu      sync.Mutex
}

type CommentSubscription struct {
	Ch chan *model.Comment
	id int
}

type Subscriptions interface {
	CreateSubscription(postId int) (int, chan *model.Comment, error)
	DeleteSubscription(postId, chanId int) error
	NotifySubscription(postId int, comment model.Comment) error
}

func NewCommentSubscription() *CommentsSubscriptions {
	return &CommentsSubscriptions{
		Chans:   make(map[int][]CommentSubscription),
		counter: 0,
		mu:      sync.Mutex{},
	}
}

func (c *CommentsSubscriptions) CreateSubscription(postId int) (int, chan *model.Comment, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	Ch := make(chan *model.Comment)
	c.counter++

	c.Chans[postId] = append(c.Chans[postId], CommentSubscription{Ch: Ch, id: c.counter})

	return c.counter, Ch, nil
}

func (c *CommentsSubscriptions) DeleteSubscription(postId, chanId int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	obs := c.Chans[postId]
	for i, sub := range obs {
		if sub.id == chanId {
			close(sub.Ch)
			c.Chans[postId] = append(obs[:i], obs[i+1:]...)
			return nil
		}
	}

	return nil
}

func (c *CommentsSubscriptions) NotifySubscription(postId int, comment model.Comment) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	obs, ok := c.Chans[postId]
	if ok {
		for _, sub := range obs {
			sub.Ch <- &comment
		}
	} else {
		return errors.New(constant.ThereIsNoSubscriptionError)
	}

	fmt.Printf("New Comment, postid = %d\n", postId)

	return nil
}

// Method to access chans field for testing purposes
func (c *CommentsSubscriptions) GetSubscriptions(postId int) []CommentSubscription {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Chans[postId]
}
