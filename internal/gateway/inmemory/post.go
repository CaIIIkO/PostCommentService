package in_memory

import (
	"errors"
	"postcommentservice/graph/model"
	"sync"
	"time"
)

type PostMemory struct {
	mu     sync.RWMutex
	posts  map[int]model.Post
	nextID int
}

func NewPostMemory() *PostMemory {
	return &PostMemory{
		posts:  make(map[int]model.Post),
		nextID: 1,
	}
}

func (p *PostMemory) CreatePost(post model.Post) (model.Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	post.ID = p.nextID
	p.nextID++
	post.CreatedAt = time.Now()

	p.posts[post.ID] = post

	return post, nil
}

func (p *PostMemory) GetPostById(id int) (model.Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	post, exists := p.posts[id]
	if !exists {
		return model.Post{}, errors.New("post not found")
	}

	return post, nil
}

func (p *PostMemory) GetAllPosts(limit, offset int) ([]*model.Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var posts []*model.Post
	count := 0

	for _, post := range p.posts {
		if count >= offset {
			posts = append(posts, &post)
		}
		count++
		if limit > 0 && len(posts) >= limit {
			break
		}
	}

	return posts, nil
}
