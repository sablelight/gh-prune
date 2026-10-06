package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v68/github"
)

type Client struct {
	client *github.Client
}

func NewClient(token string) *Client {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	return &Client{
		client: github.NewClient(httpClient).WithAuthToken(token),
	}
}

type Branch struct {
	Name     string
	SHA      string
	MergedAt *time.Time
}

func (c *Client) GetDefaultBranch(ctx context.Context, owner, repo string) (string, error) {
	repoInfo, _, err := c.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return "", err
	}
	return repoInfo.GetDefaultBranch(), nil
}

func (c *Client) GetMergedBranches(ctx context.Context, owner, repo, defaultBranch string) ([]Branch, error) {
	var allBranches []Branch
	opts := &github.BranchListOptions{ListOptions: github.ListOptions{PerPage: 100}}

	for {
		branches, resp, err := c.client.Repositories.ListBranches(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("list branches: %w", err)
		}

		for _, b := range branches {
			if b.GetName() == defaultBranch {
				continue
			}

			// Check if merged via PR
			prs, _, err := c.client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
				Head: fmt.Sprintf("%s:%s", owner, b.GetName()),
				State: "closed",
				ListOptions: github.ListOptions{PerPage: 10},
			})
			if err != nil {
				continue
			}

			var mergedAt *time.Time
			for _, pr := range prs {
				if pr.GetMerged() && pr.MergedAt != nil {
					mergedAt = pr.MergedAt
					break
				}
			}

			if mergedAt != nil {
				allBranches = append(allBranches, Branch{
					Name:     b.GetName(),
					SHA:      b.GetCommit().GetSHA(),
					MergedAt: mergedAt,
				})
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allBranches, nil
}

func (c *Client) DeleteBranch(ctx context.Context, owner, repo, branch string) error {
	ref := fmt.Sprintf("refs/heads/%s", branch)
	_, err := c.client.Git.DeleteRef(ctx, owner, repo, ref)
	return err
}