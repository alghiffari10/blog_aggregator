package commands

import (
	"context"
	"errors"

	"github.com/alghiffari10/blog_aggregator/internal/config"
	"github.com/alghiffari10/blog_aggregator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	RegisteredCommands map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.RegisteredCommands[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {

	f, ok := c.RegisteredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	return f(s, cmd)
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(
			context.Background(),
			s.Cfg.CurrentUserName,
		)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}

}
