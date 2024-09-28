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

func handlerCreateFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("create feed requires 2 arguments, found %v arguments", cmd.Arguments)
	}
	nameFeed := cmd.Arguments[0]
	urlFeed := cmd.Arguments[1]

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

	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), createFeedFollowParams)
	if err != nil {
		return fmt.Errorf("failed to create a new feed follow: %w", err)
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("follow requires 1 argument, found %v arguments", cmd.Arguments)
	}
	feedUrl := cmd.Arguments[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("failed to find feed by url: %w", err)
	}

	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), createFeedFollowParams)
	if err != nil {
		return fmt.Errorf("failed to create a new feed follow: %w", err)
	}

	fmt.Printf("User %s is now following the feed %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("following dosent require any arguments, found %v arguments", cmd.Arguments)
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows for current user: %w", err)
	}

	fmt.Println("You are following these feeds:")
	for _, follow := range feedFollows {
		fmt.Printf("- %s\n", follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("unfollow requires exactly 1 argument (feed URL), found %v arguments", len(cmd.Arguments))
	}
	feedURL := cmd.Arguments[0]

	unfollowArgs := database.DeleteFeedFollowByUserAndFeedURLParams{
		UserID: user.ID,
		Url:    feedURL,
	}

	err := s.db.DeleteFeedFollowByUserAndFeedURL(context.Background(), unfollowArgs)
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}

	fmt.Printf("Feed '%s' unfollowed successfully by user '%s'.\n", feedURL, user.Name)
	return nil
}
