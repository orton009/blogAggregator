package main

import (
	"blogAggregator/internal/database"
	"errors"

	"internal/config"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	CommandMap map[string](func(*state, command) error)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.CommandMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if handler, ok := c.CommandMap[cmd.Name]; ok {
		return handler(s, cmd)
	} else {
		return errors.New("command not supported: " + cmd.Name)
	}
}
