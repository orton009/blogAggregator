package main

import (
	"blogAggregator/internal/database"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func followFeed(s *state, feedUrl string, user database.User) (database.CreateFeedFollowRow, error) {

	feedFollow := database.CreateFeedFollowRow{}
	feed, err := s.db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		return feedFollow, err
	}

	feedFollow, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return feedFollow, err
	}

	return feedFollow, nil
}

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {

	return func(s *state, c command) error {

		if s.config.CurrentUserName == "" {
			return fmt.Errorf("you must be logged in to perform this action")
		}

		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, c, user)
	}
}
