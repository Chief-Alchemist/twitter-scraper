package twitterscraper_test

import (
	"testing"

	twitterscraper "github.com/Chief-Alchemist/twitter-scraper"
)

func TestFetchFollowing(t *testing.T) {
	if skipAuthTest {
		t.Skip("Skipping test due to environment variable")
	}
	users, _, err := testScraper.FetchFollowing("Support", 20, "")
	if err != nil {
		t.Error(err)
	}
	if len(users) < 1 || users[len(users)-1].Username == "" {
		t.Error("error FetchFollowing() No users found")
	}
}

func TestFetchFollowers(t *testing.T) {
	if skipAuthTest {
		t.Skip("Skipping test due to environment variable")
	}
	users, _, err := testScraper.FetchFollowers("Support", 20, "")
	if err != nil {
		t.Error(err)
	}
	if len(users) < 1 || users[len(users)-1].Username == "" {
		t.Error("error FetchFollowing() No users found")
	}
}

func TestFollowUser(t *testing.T) {
	if skipAuthTest {
		t.Skip("Skipping test due to environment variable")
	}
	err := testScraper.FollowUser("Support", twitterscraper.Follow)
	if err != nil {
		t.Error(err)
	}
}

func TestUnfollowUser(t *testing.T) {
	if skipAuthTest {
		t.Skip("Skipping test due to environment variable")
	}
	err := testScraper.FollowUser("Support", twitterscraper.Unfollow)
	if err != nil {
		t.Error(err)
	}
}
