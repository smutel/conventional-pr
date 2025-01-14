package service

import (
	"context"

	"github.com/Namchee/conventional-pr/internal"
	"github.com/Namchee/conventional-pr/internal/entity"
	"github.com/Namchee/conventional-pr/internal/formatter"
	"github.com/google/go-github/v32/github"
)

// GithubService is a service that simplifies GitHub API interaction
type GithubService struct {
	client internal.GithubClient
	config *entity.Config
	meta   *entity.Meta
}

// NewGithubService creates a new GitHub service that simplify API interaction with functions which is actually needed
func NewGithubService(
	client internal.GithubClient,
	config *entity.Config,
	meta *entity.Meta,
) *GithubService {
	return &GithubService{
		client: client,
		config: config,
		meta:   meta,
	}
}

// WriteReport creates a new comment that contains conventional-pr workflow report in markdown format
func (s *GithubService) WriteReport(
	pullRequest *github.PullRequest,
	whitelistResults []*entity.WhitelistResult,
	validationResults []*entity.ValidationResult,
) error {
	report := formatter.FormatResultToTables(whitelistResults, validationResults)

	ctx := context.Background()

	return s.client.Comment(
		ctx,
		s.meta.Owner,
		s.meta.Name,
		pullRequest.GetNumber(),
		&github.IssueComment{
			Body: &report,
		},
	)
}

// WriteTemplate creates a new comment that contains user-desired message
func (s *GithubService) WriteTemplate(
	pullRequest *github.PullRequest,
) error {
	if s.config.Template == "" {
		return nil
	}

	ctx := context.Background()

	return s.client.Comment(
		ctx,
		s.meta.Owner,
		s.meta.Name,
		pullRequest.GetNumber(),
		&github.IssueComment{
			Body: &s.config.Template,
		},
	)
}

// AttachLabel attachs label to invalid pull request
func (s *GithubService) AttachLabel(
	pullRequest *github.PullRequest,
) error {
	if s.config.Label == "" {
		return nil
	}

	ctx := context.Background()

	return s.client.Label(
		ctx,
		s.meta.Owner,
		s.meta.Name,
		pullRequest.GetNumber(),
		s.config.Label,
	)
}

// ClosePullRequest closes invalid pull request
func (s *GithubService) ClosePullRequest(
	pullRequest *github.PullRequest,
) error {
	if !s.config.Close {
		return nil
	}

	ctx := context.Background()

	return s.client.Close(
		ctx,
		s.meta.Owner,
		s.meta.Name,
		pullRequest.GetNumber(),
	)
}
