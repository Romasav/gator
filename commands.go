package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Romasav/gator/internal/database"
	"github.com/Romasav/gator/rssFeed"
	"github.com/google/uuid"
)

type commands struct {
	handlers map[string]func(*state, command) error
}

func newCommands() *commands {
	return &commands{
		handlers: map[string]func(*state, command) error{},
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.Name]
	if !exists {
		return fmt.Errorf("command with a name %v dosent exists", cmd.Name)
	}
	err := handler(s, cmd)
	if err != nil {
		return fmt.Errorf("filed to run the command: %w", err)
	}
	return nil
}

type command struct {
	Name      string
	Arguments []string
}

func newCommand(name string, args []string) *command {
	return &command{
		Name:      name,
		Arguments: args,
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("login requires exactly 1 argument (username), found %v arguments", cmd.Arguments)
	}
	username := cmd.Arguments[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		return errors.New("the user dose not exists")
	}
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	err = s.config.SetUpUser(username)
	if err != nil {
		return fmt.Errorf("failed to set up user: %w", err)
	}
	fmt.Println("The user has been set!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("register requires exactly 1 argument (username), found %v arguments", cmd.Arguments)
	}
	username := cmd.Arguments[0]

	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return errors.New("the user already exists")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	user, err := s.db.CreateUser(context.Background(), createUserParams)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = s.config.SetUpUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to set up user: %w", err)
	}

	fmt.Println("The user has been created!")
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("reset dosent require any arguments, found %v arguments", cmd.Arguments)
	}

	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete all users: %w", err)
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("users dosent require any arguments, found %v arguments", cmd.Arguments)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get all users: %w", err)
	}

	for _, user := range users {
		fmt.Print(user.Name)
		if user.Name == s.config.Username {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}

	return nil
}

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("aggregator dosent require any arguments, found %v arguments", cmd.Arguments)
	}

	url := "https://www.wagslane.dev/index.xml"
	rssFeed, err := rssFeed.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed fetch feed: %w", err)
	}

	fmt.Println(rssFeed)

	return nil
}

func handlerCreateFeed(s *state, cmd command) error {
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("create feed requires 2 arguments, found %v arguments", cmd.Arguments)
	}
	nameFeed := cmd.Arguments[0]
	urlFeed := cmd.Arguments[1]

	user, err := s.db.GetUser(context.Background(), s.config.Username)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	createFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      nameFeed,
		Url:       urlFeed,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), createFeedParams)
	if err != nil {
		return fmt.Errorf("failed to create a feed: %w", err)
	}

	fmt.Println("New Feed Record:")
	fmt.Printf("ID:        %s\n", feed.ID.String())
	fmt.Printf("Name:      %s\n", feed.Name)
	fmt.Printf("URL:       %s\n", feed.Url)
	fmt.Printf("User ID:   %s\n", feed.UserID.String())
	fmt.Printf("CreatedAt: %s\n", feed.CreatedAt.Format(time.RFC3339))
	fmt.Printf("UpdatedAt: %s\n", feed.UpdatedAt.Format(time.RFC3339))

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("feeds dosent require any arguments, found %v arguments", cmd.Arguments)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %w", err)
	}

	for index, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user by id: %w", err)
		}

		fmt.Printf("%v Feed Record:\n", index+1)
		fmt.Printf("Name:      %s\n", feed.Name)
		fmt.Printf("URL:       %s\n", feed.Url)
		fmt.Printf("User Name: %s\n", user.Name)
	}

	return nil
}
