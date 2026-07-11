package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/alghiffari10/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerRegister(s *State, cmd Command) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

func HandlerLogin(s *State, cmd Command) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	_, err := s.Db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.Cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func HandlerReset(s *State, cmd Command) error {

	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: database")
	}

	err := s.Db.ResetUser(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Database reset successfully")

	return nil
}
func HandlerUsers(s *State, cmd Command) error {

	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: users")
	}

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %v (current)", user.Name)
		} else {
			fmt.Printf("* %v", user.Name)
		}
	}

	return nil
}

func printUser(user database.User) {

	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
