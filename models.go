package main

import (
	"time"

	"github.com/bryant-bourgeois/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	if !feed.LastFetchedAt.Valid {
		return Feed{
			ID:        feed.ID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
			Name:      feed.Name,
			Url:       feed.Url,
			UserID:    feed.UserID,
		}
	}
	return Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: &feed.LastFetchedAt.Time,
	}
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseFeedFollowToFeedFollow(follow database.UsersFeed) FeedFollow {
	return FeedFollow{
		ID:        follow.ID,
		UserID:    follow.UserID,
		FeedID:    follow.FeedID,
		CreatedAt: follow.CreatedAt,
		UpdatedAt: follow.UpdatedAt,
	}
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	FeedID      uuid.UUID `json:"feed_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedAt string    `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Url         string    `json:"url"`
}

func databasePostsToPosts(posts []database.GetPostsByUserIdRow) []Post {
	output := make([]Post, 0)
	for _, val := range posts {
		output = append(output, Post{
			ID:          val.ID.UUID,
			FeedID:      val.FeedID.UUID,
			Title:       val.Title.String,
			Description: val.Description.String,
			PublishedAt: val.PublishedAt.String,
			CreatedAt:   val.CreatedAt.Time,
			UpdatedAt:   val.UpdatedAt.Time,
			Url:         val.Url.String,
		})
	}
	return output
}
