package mocks

import (
	"context"
	"errors"

	"github.com/Namchee/conventional-pr/internal"
	"github.com/Namchee/conventional-pr/internal/constants"
	"github.com/google/go-github/v32/github"
)

// GitHub's client mock. Used in testing
type githubClientMock struct{}

func (m *githubClientMock) GetPullRequest(
	_ context.Context,
	_ string,
	_ string,
	event int,
) (*github.PullRequest, error) {
	if event == 123 {
		return &github.PullRequest{}, nil
	}

	return nil, errors.New("not found")
}

func (m *githubClientMock) GetUser(_ context.Context, name string) (*github.User, error) {
	bot := constants.BotUser
	user := "user"

	if name == "foo" {
		return &github.User{Type: &bot}, nil
	}

	return &github.User{Type: &user}, nil
}

func (m *githubClientMock) GetIssue(
	ctx context.Context,
	_ string,
	_ string,
	number int,
) (*github.Issue, error) {
	if number == 123 {
		return &github.Issue{}, nil
	}

	return nil, nil
}

func (m *githubClientMock) GetPermissionLevel(
	_ context.Context,
	_ string,
	_ string,
	user string,
) (*github.RepositoryPermissionLevel, error) {
	admin := constants.AdminUser
	writeOnly := "write"

	if user == "foo" {
		return &github.RepositoryPermissionLevel{
			Permission: &admin,
		}, nil
	}

	return &github.RepositoryPermissionLevel{
		Permission: &writeOnly,
	}, nil
}

func (m *githubClientMock) GetCommits(
	ctx context.Context,
	owner string,
	name string,
	event int,
) ([]*github.RepositoryCommit, error) {
	good := "feat(test): test something"
	bad := "this is bad"

	verifiedTrue := true
	verifiedFalse := false

	if event == 123 {
		return []*github.RepositoryCommit{
			{
				Commit: &github.Commit{
					Message: &good,
					Verification: &github.SignatureVerification{
						Verified: &verifiedTrue,
					},
				},
			},
		}, nil
	}

	if event == 69 {
		return []*github.RepositoryCommit{
			{
				SHA: &bad,
				Commit: &github.Commit{
					Message: &bad,
					Verification: &github.SignatureVerification{
						Verified: &verifiedFalse,
					},
				},
			},
		}, nil
	}

	return []*github.RepositoryCommit{}, nil
}

func (m *githubClientMock) Comment(
	_ context.Context,
	_ string,
	_ string,
	event int,
	_ *github.IssueComment,
) error {
	if event == 123 {
		return nil
	}

	return errors.New("Error")
}

func (m *githubClientMock) Label(
	_ context.Context,
	_ string,
	_ string,
	event int,
	_ string,
) error {
	if event == 123 {
		return nil
	}

	return errors.New("Error")
}

func (m *githubClientMock) Close(
	_ context.Context,
	_ string,
	_ string,
	event int,
) error {
	if event == 123 {
		return nil
	}

	return errors.New("Error")
}

// NewGithubClientMock creates a GitHub client mock for testing purposes
func NewGithubClientMock() internal.GithubClient {
	return &githubClientMock{}
}
