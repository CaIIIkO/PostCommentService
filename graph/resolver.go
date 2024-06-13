package graph

import "postcommentservice/internal/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostsService        service.Post
	CommentsService     service.Comment
	SubscriptionService service.Subscriptions
}
