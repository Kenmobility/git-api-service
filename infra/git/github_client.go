package git

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kenmobility/git-api-service/common/client"
	"github.com/kenmobility/git-api-service/internal/domains"
	"github.com/rs/zerolog/log"
)

type GitHubClient struct {
	baseURL         string
	token           string
	fetchInterval   time.Duration
	client          *client.RestClient
	rateLimitFields rateLimitFields
}

type rateLimitFields struct {
	rateLimitLimit     int
	rateLimitRemaining int
	rateLimitReset     int
}

func (g *GitHubClient) getHeaders() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", g.token),
	}
}

func NewGitHubClient(baseUrl string, token string, fetchInterval time.Duration) GitManagerClient {
	client := client.NewRestClient()

	gc := GitHubClient{
		baseURL:       baseUrl,
		token:         token,
		fetchInterval: fetchInterval,
		client:        client,
	}
	ts := GitManagerClient(&gc)
	return ts
}

func (g *GitHubClient) FetchRepoMetadata(ctx context.Context, repositoryName string) (*domains.RepoMetadata, error) {
	endpoint := fmt.Sprintf("%s/repos/%s", g.baseURL, repositoryName)

	resp, err := g.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("repo status not fetched")
	}

	var gitHubRepoResponse GitHubRepoMetadataResponse

	if err := json.Unmarshal([]byte(resp.Body), &gitHubRepoResponse); err != nil {
		log.Error().Msgf("marshal error, [%v]", err)
		return nil, errors.New("could not unmarshal repo metadata response")
	}

	repoMetadata := &domains.RepoMetadata{
		Name:            gitHubRepoResponse.FullName,
		Description:     gitHubRepoResponse.Description,
		URL:             gitHubRepoResponse.Url,
		Language:        gitHubRepoResponse.Language,
		ForksCount:      gitHubRepoResponse.ForksCount,
		StarsCount:      gitHubRepoResponse.StargazersCount,
		OpenIssuesCount: gitHubRepoResponse.OpenIssues,
		WatchersCount:   gitHubRepoResponse.WatchersCount,
	}

	return repoMetadata, nil
}

func (g *GitHubClient) FetchCommits(ctx context.Context, repo domains.RepoMetadata, since time.Time, until time.Time, lastFetchedCommit string, page, perPage int) ([]domains.Commit, bool, error) {
	var endpoint string

	if lastFetchedCommit != "" {
		endpoint = fmt.Sprintf("%s/repos/%s/commits?sha=%s&per_page=%d&page=%d", g.baseURL, repo.Name, lastFetchedCommit, perPage, page)
	} else {
		endpoint = fmt.Sprintf("%s/repos/%s/commits?since=%s&until=%s&per_page=%d&page=%d", g.baseURL, repo.Name, since.Format(time.RFC3339), until.Format(time.RFC3339), perPage, page)
	}

	response, err := g.client.Get(endpoint, map[string]string{}, g.getHeaders())
	if err != nil {
		log.Error().Msgf("error fetching commits: %v", err)

		return nil, false, err
	}

	if response.StatusCode == http.StatusForbidden {
		return nil, false, fmt.Errorf("rate limit exceeded")
	}

	g.updateRateLimitHeaders(response)

	if g.rateLimitFields.rateLimitRemaining == 0 {
		waitTime := time.Until(time.Unix(int64(g.rateLimitFields.rateLimitReset), 0))
		log.Info().Msgf("Rate limit exceeded. Waiting for %v until reset...", waitTime)
	}

	if response.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("failed to fetch commits; status code: %v", response.StatusCode)
	}

	var commitRes []GithubCommitResponse

	if err := json.Unmarshal([]byte(response.Body), &commitRes); err != nil {
		fmt.Printf("marshal error, [%v]", err)
		return nil, false, errors.New("could not unmarshal commits response")
	}

	var cc []domains.Commit
	for _, cr := range commitRes {
		commit := domains.Commit{
			CommitID:       cr.SHA,
			Message:        cr.Commit.Message,
			Author:         cr.Commit.Author.Name,
			Date:           cr.Commit.Author.Date,
			URL:            cr.HtmlURL,
			RepositoryName: repo.Name,
		}

		cc = append(cc, commit)
	}

	// Determine if more pages exist
	morePages := false
	linkHeader := response.Headers["Link"]
	if len(linkHeader) > 0 {
		morePages = g.hasNextPage(linkHeader[0])
	}

	return cc, morePages, nil
}

// hasNextPage checks if there is a 'next' link in the Link header
func (g *GitHubClient) hasNextPage(linkHeader string) bool {
	links := g.parseLinkHeader(linkHeader)
	_, hasNext := links["next"]
	return hasNext
}

// parseLinkHeader parses the Link header into a map
func (g *GitHubClient) parseLinkHeader(header string) map[string]string {
	links := make(map[string]string)

	for _, part := range strings.Split(header, ",") {
		sections := strings.Split(part, ";")
		if len(sections) < 2 {
			continue
		}
		url := strings.Trim(sections[0], " <>")
		rel := strings.Trim(sections[1], " rel=\"")
		links[rel] = url
	}
	return links
}

// updateRateLimitHeaders updates the link header rate limit values to the fields
func (api *GitHubClient) updateRateLimitHeaders(resp *client.Response) {
	limit := resp.Headers["X-Ratelimit-Limit"]
	if len(limit) > 0 {
		api.rateLimitFields.rateLimitLimit, _ = strconv.Atoi(limit[0])
	}

	remaining := resp.Headers["X-Ratelimit-Remaining"]
	if len(remaining) > 0 {
		api.rateLimitFields.rateLimitRemaining, _ = strconv.Atoi(remaining[0])
		log.Info().Msgf("Rate limit remaining: %d", api.rateLimitFields.rateLimitRemaining)
	}

	reset := resp.Headers["X-Ratelimit-Reset"]
	if len(reset) > 0 {
		api.rateLimitFields.rateLimitReset, _ = strconv.Atoi(reset[0])
	}

	used := resp.Headers["X-Ratelimit-Used"]
	if len(used) > 0 {
		usedInt, _ := strconv.Atoi(used[0])
		log.Info().Msgf("Rate limit used: %d/%d", usedInt, api.rateLimitFields.rateLimitLimit)
	}
}
