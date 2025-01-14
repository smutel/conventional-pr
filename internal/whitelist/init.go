package whitelist

import (
	"sort"
	"sync"

	"github.com/Namchee/conventional-pr/internal"
	"github.com/Namchee/conventional-pr/internal/entity"
	"github.com/google/go-github/v32/github"
)

var (
	whitelists = []func(internal.GithubClient, *entity.Config, *entity.Meta) internal.Whitelist{
		NewBotWhitelist,
		NewDraftWhitelist,
		NewPermissionWhitelist,
		NewUsernameWhitelist,
	}
)

// WhitelistGroup is a collection of whitelisting process, integrated in one single function call
type WhitelistGroup struct {
	client internal.GithubClient
	config *entity.Config
	meta   *entity.Meta
	wg     *sync.WaitGroup
}

// NewWhitelistGroup creates a new WhitelistGroup
func NewWhitelistGroup(
	client internal.GithubClient,
	config *entity.Config,
	meta *entity.Meta,
	wg *sync.WaitGroup,
) *WhitelistGroup {
	return &WhitelistGroup{
		client: client,
		config: config,
		meta:   meta,
		wg:     wg,
	}
}

func (w *WhitelistGroup) processWhitelist(
	whitelist internal.Whitelist,
	pullRequest *github.PullRequest,
	pool chan *entity.WhitelistResult,
) {
	defer w.wg.Done()
	result := whitelist.IsWhitelisted(pullRequest)
	pool <- result
}

func (w *WhitelistGroup) cleanup(
	channel chan *entity.WhitelistResult,
) {
	w.wg.Wait()
	close(channel)
}

// Process the pull request with all available whitelists
func (w *WhitelistGroup) Process(
	pullRequest *github.PullRequest,
) []*entity.WhitelistResult {
	channel := make(chan *entity.WhitelistResult, len(whitelists))
	w.wg.Add(len(whitelists))

	for _, wv := range whitelists {
		wl := wv(w.client, w.config, w.meta)

		go w.processWhitelist(wl, pullRequest, channel)
	}

	go w.cleanup(channel)

	var results []*entity.WhitelistResult

	for result := range channel {
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Result == results[j].Result {
			return results[i].Name < results[j].Name
		}

		if results[i].Result {
			return true
		}

		return false
	})

	return results
}

// IsWhitelisted checks if a pull request is whitelisted or not from whitelist results
func IsWhitelisted(result []*entity.WhitelistResult) bool {
	for _, r := range result {
		if r.Active && r.Result {
			return true
		}
	}

	return false
}
