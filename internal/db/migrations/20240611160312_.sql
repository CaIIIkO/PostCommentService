-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts(
    id serial primary key,
    createdAt timestamp default now(),
    title varchar(100),
    content varchar(2000),
    author varchar(100),
    commentsAllowed boolean
);

CREATE TABLE comments(
    id serial primary key,
    createdAt timestamp default now(),
    content varchar(2000),
    author varchar(100),
    answers int DEFAULT 0,
    post int not null,
    FOREIGN KEY (post) REFERENCES Posts(id) ON DELETE CASCADE ON UPDATE CASCADE ,
    replyTo int,
    FOREIGN KEY (replyTo) REFERENCES Comments(id) ON DELETE SET NULL ON UPDATE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
DROP TABLE comments;
-- +goose StatementEnd
