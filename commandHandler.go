package main

import (
	"blogAggregator/internal/database"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const FeedURL string = "https://www.wagslane.dev/index.xml"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("expected an argument")
	}

	username := cmd.Args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	err = s.config.SetUser(user.Name)
	if err == nil {

		fmt.Println("User " + username + " has been set.")
		return nil
	}

	fmt.Println("welcome back " + user.Name)
	return err
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("expected username as argument")
	}

	username := cmd.Args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("user already exist %v", user)
	}

	user, createError := s.db.CreateUser(context.Background(), database.CreateUserParams{
		Name:      username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
	})
	if createError != nil {
		return createError
	}

	fmt.Println("User successfully created")
	fmt.Println(user)

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("User successfully Set")

	return nil
}

func handlerReset(s *state, c command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleteing entries in database\n %w ", err)
	}

	fmt.Println("all users removed successfully!")
	return nil
}

func handlerList(s *state, c command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error listing users %w", err)
	}

	currentUser := s.config.CurrentUserName

	for _, u := range users {
		if u.Name == currentUser {

			fmt.Println(u.Name + " (current)")
		} else {

			fmt.Println(u.Name)
		}
	}

	return nil
}

func handlerAgg(s *state, c command) error {
	if len(c.Args) < 1 {
		return fmt.Errorf("expected atleast one argument, time between fetching")
	}

	timeBetween := c.Args[0]

	timeBetweenDuration, err := time.ParseDuration(timeBetween)
	if err != nil {
		return fmt.Errorf("error parsing time between fetching \n %w", err)
	}

	ticker := time.NewTicker(timeBetweenDuration)
	defer ticker.Stop()

	for range ticker.C {

		feed, err := s.db.GetNextFeedToFetch(context.Background())
		if err != nil {
			return fmt.Errorf("error getting next feed to fetch \n %w", err)
		}

		rssFeed, err := fetchFeed(context.Background(), feed.Url)
		if err != nil {
			return fmt.Errorf("error fetching rss feed \n %w %v", err, feed)
		}

		err = s.db.MarkFeedAsFetched(context.Background(), feed.ID)
		if err != nil {
			return fmt.Errorf("error marking feed as fetched \n %w", err)
		}

		fmt.Println("rss feed fetched successfully\n", rssFeed)
	}
	return nil
}

func handlerCreateFeed(s *state, c command, user database.User) error {
	if len(c.Args) < 2 {
		return fmt.Errorf("expected atleast two arguments, name and url of the feed")
	}

	rssFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:      c.Args[0],
		Url:       c.Args[1],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		ID:        uuid.New(),
	})

	if err != nil {
		return fmt.Errorf("error creating feed entry\n %w", err)
	}

	_, err = followFeed(s, rssFeed.Url, user)
	if err != nil {
		return fmt.Errorf("error following feed \n %w", err)
	}

	fmt.Println("created rss feed successfully \n", rssFeed)
	return nil

}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {

		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}

		fmt.Println(feed.Name, "\t", feed.Url, "\t", user.Name)
	}

	return nil
}

func handlerFollowing(s *state, _ command) error {
	feedFollows, err := s.db.GetFeedFollows(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	for _, feedFollow := range feedFollows {

		fmt.Println("- ", feedFollow.FeedName)
	}

	return nil
}

func handlerFollow(s *state, c command, user database.User) error {
	if len(c.Args) < 1 {
		return fmt.Errorf("expected atleast one argument, url of the feed to follow")
	}

	feedUrl := c.Args[0]

	feedFollow, err := followFeed(s, feedUrl, user)

	if err != nil {
		return fmt.Errorf("error creating feed follow \n %w", err)
	}

	fmt.Println("feed follow created \n", feedFollow)
	return nil
}

func handlerUnfollow(s *state, c command, user database.User) error {
	if len(c.Args) < 1 {
		return fmt.Errorf("expected atleast one argument, url of the feed to unfollow")
	}

	feedUrl := c.Args[0]

	feed, err := s.db.GetFeed(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("error deleting feed follow \n %w", err)
	}

	fmt.Println("feed follow deleted successfully")
	return nil
}
