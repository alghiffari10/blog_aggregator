package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/alghiffari10/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerAggregator(s *State, cmd Command) error {
	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", duration)

	ticker := time.NewTicker(duration)

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Println(err)
		}
	}
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	feedFollow, err := s.Db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")

	return nil
}

func HandlerFeeds(s *State, cmd Command) error {

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		fmt.Printf("Name: %v\n", feed.Name)
		fmt.Printf("URL: %v\n", feed.Url)
		fmt.Printf("Created by: %v\n\n", feed.UserName)
	}

	return nil
}

func HandlerFollows(s *State, cmd Command, user database.User) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: follow <url> ")
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	follow, err := s.Db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		return err
	}

	printFeedFollow(follow.UserName, follow.FeedName)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {

	follows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, follow := range follows {
		fmt.Printf("- %v\n", follow.FeedName)
	}
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: <unfollow> <url>")
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	err = s.Db.DeleteFeedFollow(context.Background(),
		database.DeleteFeedFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		})
	if err != nil {
		return err
	}

	fmt.Printf("Unfollowed %v", feed.Name)

	return nil
}

func scrapeFeeds(s *State) error {

	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.Db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("Fetching... %v\n\n", feed.Name)

	for _, item := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err = s.Db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
	fmt.Printf("* LastFetchedAt: %v\n", feed.LastFetchedAt.Time)
}
func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
