package postgres_test

import (
	"regexp"
	"testing"
	"time"

	"postcommentservice/graph/model"
	"postcommentservice/internal/gateway/postgres"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := postgres.NewPostPostgres(sqlxDB)

	post := model.Post{
		Title:           "Test Title",
		Content:         "Test Content",
		Author:          "Test Author",
		CommentsAllowed: true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO Posts (title, content, author, commentsAllowed) 
				VALUES ($1, $2, $3, $4)
				RETURNING id, createdAt`)).
		WithArgs(post.Title, post.Content, post.Author, post.CommentsAllowed).
		WillReturnRows(sqlmock.NewRows([]string{"id", "createdAt"}).AddRow(1, time.Now()))

	mock.ExpectCommit()

	createdPost, err := repo.CreatePost(post)
	assert.NoError(t, err)
	assert.NotZero(t, createdPost.ID)
	assert.NotZero(t, createdPost.CreatedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}
