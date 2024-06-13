package postgres

import (
	"postcommentservice/graph/model"

	"github.com/jmoiron/sqlx"
)

type CommentsPostgres struct {
	db *sqlx.DB
}

func NewCommentsPostgres(db *sqlx.DB) *CommentsPostgres {
	return &CommentsPostgres{db: db}
}

func (c CommentsPostgres) CreateComment(comment model.Comment) (model.Comment, error) {

	tx, err := c.db.Begin()
	if err != nil {
		return model.Comment{}, err
	}

	query := `INSERT INTO comments (content, author, post, replyTo) 
				VALUES ($1, $2, $3, $4) RETURNING id, createdAt`

	row := tx.QueryRow(query, comment.Content, comment.Author, comment.Post, comment.ReplyTo)
	if err := row.Scan(&comment.ID, &comment.CreatedAt); err != nil {
		tx.Rollback()
		return model.Comment{}, err
	}

	if comment.ReplyTo != nil {
		_, err := tx.Exec(`UPDATE comments SET answers = answers + 1 WHERE id = $1`, *comment.ReplyTo)
		if err != nil {
			tx.Rollback()
			return model.Comment{}, err
		}
	}

	return comment, tx.Commit()
}

func (c CommentsPostgres) GetCommentsByPost(postId, limit, offset int) ([]*model.Comment, error) {

	query := `SELECT * FROM comments
	     WHERE post = $1 AND replyTo IS NULL
	     ORDER BY createdAt
	     OFFSET $2`

	args := []interface{}{postId, offset}

	if limit >= 0 {
		query += " LIMIT $3"
		args = append(args, limit)
	}

	var comments []*model.Comment

	if err := c.db.Select(&comments, query, args...); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c CommentsPostgres) GetRepliesOfComment(commentId int) ([]*model.Comment, error) {

	query := `SELECT * FROM comments WHERE replyTo = $1`

	var comments []*model.Comment

	if err := c.db.Select(&comments, query, commentId); err != nil {
		return nil, err
	}

	return comments, nil
}
