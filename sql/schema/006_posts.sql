-- +goose Up
CREATE TABLE posts (
	id uuid PRIMARY KEY,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	title text NOT NULL,
	url text UNIQUE NOT NULL,
	description text NOT NULL,
	published_at text NOT NULL,
	feed_id uuid references feeds(id) ON DELETE CASCADE NOT NULL
);

-- +goose Down
DROP TABLE posts;
