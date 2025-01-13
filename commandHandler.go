package main

import (
	"blogAggregator/internal/database"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

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
	fmt.Println(user)
	return err
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("expected username as argument")
	}

	username := cmd.Args[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("user already exist \n %v \n", user)
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
	err := s.db.DeleteAll(context.Background())
	if err != nil {
		fmt.Errorf("error deleteing entries in database\n %w ", err)
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
