package service_test

import (
	"postcommentservice/graph/model"
	constant "postcommentservice/internal/consts"
	"postcommentservice/internal/service"
	"testing"
)

func TestCommentsSubscriptions_CreateSubscription(t *testing.T) {
	cs := service.NewCommentSubscription()

	postId := 1
	_, ch, err := cs.CreateSubscription(postId)
	if err != nil {
		t.Fatalf("CreateSubscription() error = %v, wantErr %v", err, false)
	}

	if ch == nil {
		t.Errorf("CreateSubscription() channel is nil, want non-nil channel")
	}

	csChans := cs.GetSubscriptions(postId)

	if len(csChans) != 1 {
		t.Errorf("CreateSubscription() len = %d, want %d", len(csChans), 1)
	}

	if csChans[0].Ch != ch {
		t.Errorf("CreateSubscription() created channel is not same as stored channel")
	}
}

func TestCommentsSubscriptions_DeleteSubscription(t *testing.T) {
	cs := service.NewCommentSubscription()

	postId := 1
	chanId, ch, err := cs.CreateSubscription(postId)
	if err != nil {
		t.Fatalf("CreateSubscription() error = %v, wantErr %v", err, false)
	}

	err = cs.DeleteSubscription(postId, chanId)
	if err != nil {
		t.Fatalf("DeleteSubscription() error = %v, wantErr %v", err, false)
	}

	csChans := cs.GetSubscriptions(postId)

	if len(csChans) != 0 {
		t.Errorf("DeleteSubscription() len = %d, want %d", len(csChans), 0)
	}

	select {
	case _, ok := <-ch:
		if ok {
			t.Errorf("DeleteSubscription() channel is not closed")
		}
	default:
	}
}

func TestCommentsSubscriptions_NotifySubscription(t *testing.T) {
	cs := service.NewCommentSubscription()

	postId := 1
	_, ch, err := cs.CreateSubscription(postId)
	if err != nil {
		t.Fatalf("CreateSubscription() error = %v, wantErr %v", err, false)
	}

	comment := model.Comment{ID: 1, Author: "Author", Content: "Content", Post: postId}
	go func() {
		err := cs.NotifySubscription(postId, comment)
		if err != nil {
			t.Errorf("NotifySubscription() error = %v, wantErr %v", err, false)
		}
	}()

	receivedComment := <-ch
	if receivedComment.ID != comment.ID || receivedComment.Author != comment.Author || receivedComment.Content != comment.Content || receivedComment.Post != comment.Post {
		t.Errorf("NotifySubscription() got = %v, want %v", receivedComment, comment)
	}
}

func TestCommentsSubscriptions_NotifySubscription_NoSubscription(t *testing.T) {
	cs := service.NewCommentSubscription()

	postId := 1
	comment := model.Comment{ID: 1, Author: "Author", Content: "Content", Post: postId}
	err := cs.NotifySubscription(postId, comment)
	if err == nil || err.Error() != constant.ThereIsNoSubscriptionError {
		t.Errorf("NotifySubscription() error = %v, want %v", err, constant.ThereIsNoSubscriptionError)
	}
}
