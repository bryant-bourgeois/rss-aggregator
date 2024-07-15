-- +goose Up
CREATE TABLE users_feeds (
	id uuid PRIMARY KEY,
	user_id uuid references users(id) ON DELETE CASCADE NOT NULL,
	feed_id uuid references feeds(id) ON DELETE CASCADE NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);

-- +goose Down
DROP TABLE users_feeds;
