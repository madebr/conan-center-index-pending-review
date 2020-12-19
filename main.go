package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Model
type Package struct {
	FullName        string
	Description     string
	StarsCount      int
	ForksCount      int
	LastUpdatedBy   string
	OpenIssuesCount int
}

func main() {
	context := context.Background()

	var httpClient *http.Client
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		fmt.Printf("There is no token specifiec!")
		os.Exit(1)
	}

	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient = oauth2.NewClient(context, tokenService)
	client := github.NewClient(httpClient)

	repo, _, err := client.Repositories.Get(context, "conan-io", "conan-center-index")

	if err != nil {
		fmt.Printf("Problem in getting repository information %v\n", err)
		os.Exit(1)
	}

	pack := &Package{
		FullName:        *repo.FullName,
		Description:     *repo.Description,
		ForksCount:      *repo.ForksCount,
		StarsCount:      *repo.StargazersCount,
		OpenIssuesCount: *repo.OpenIssuesCount,
	}

	fmt.Printf("%+v\n", pack)

	pulls, _, err := client.PullRequests.List(context, "conan-io", "conan-center-index", &github.PullRequestListOptions{})
	for _, pr := range pulls {
		fmt.Printf("pulls/%d - Reviews: %d", *pr.Number, *pr.Comments)
	}

	commitInfo, _, err := client.Repositories.ListCommits(context, "Golang-Coach", "Lessons", nil)

	if err != nil {
		fmt.Printf("Problem in commit information %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", commitInfo[0]) // Last commit information

	// Get Rate limit information
	rateLimit, _, err := client.RateLimits(context)
	if err != nil {
		fmt.Printf("Problem in getting rate limit information %v\n", err)
		return
	}

	fmt.Printf("Limit: %d \nRemaining %d \n", rateLimit.Core.Limit, rateLimit.Core.Remaining) // Last commit information
}
