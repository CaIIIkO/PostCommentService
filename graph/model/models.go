package model

func (p InputPost) FromInput() Post {
	return Post{
		Title:           p.Title,
		Author:          p.Author,
		Content:         p.Content,
		CommentsAllowed: p.CommentsAllowed,
	}
}

func (c InputComment) FromInput() Comment {
	return Comment{
		Author:  c.Author,
		Content: c.Content,
		Post:    c.Post,
		ReplyTo: c.ReplyTo,
	}
}
