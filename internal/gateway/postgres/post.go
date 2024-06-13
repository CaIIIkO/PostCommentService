package postgres

import (
	"postcommentservice/graph/model"

	"github.com/jmoiron/sqlx"
)

type PostPostgres struct {
	db *sqlx.DB
}

func NewPostPostgres(db *sqlx.DB) *PostPostgres {
	return &PostPostgres{db: db}
}

func (p PostPostgres) CreatePost(post model.Post) (model.Post, error) {
	query := `INSERT INTO Posts (title, content, author, commentsAllowed) 
				VALUES ($1, $2, $3, $4)
				RETURNING id, createdAt`

	tx, err := p.db.Begin()
	if err != nil {
		return model.Post{}, err
	}

	row := tx.QueryRow(query, post.Title, post.Content, post.Author, post.CommentsAllowed)
	if err := row.Scan(&post.ID, &post.CreatedAt); err != nil {
		tx.Rollback()
		return model.Post{}, err
	}

	return post, tx.Commit()
}

func (p PostPostgres) GetPostById(id int) (model.Post, error) {

	query := `SELECT * FROM posts WHERE id = $1`

	var post model.Post

	if err := p.db.Get(&post, query, id); err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (p PostPostgres) GetAllPosts(limit, offset int) ([]*model.Post, error) {

	query := "SELECT * FROM posts ORDER BY createdAt OFFSET $1"
	args := []interface{}{offset}

	if limit > 0 {
		query += " LIMIT $2"
		args = append(args, limit)
	}

	var posts []*model.Post

	if err := p.db.Select(&posts, query, args...); err != nil {
		return nil, err
	}

	return posts, nil
}
