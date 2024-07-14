-- +goose Up
CREATE TABLE feeds (
	id uuid PRIMARY KEY,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	name varchar(64) NOT NULL,
	url text UNIQUE NOT NULL,
	user_id uuid references users(id) ON DELETE CASCADE NOT NULL
);

-- +goose Down
DROP TABLE feeds;
