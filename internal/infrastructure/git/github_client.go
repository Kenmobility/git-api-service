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

	"github.com/kenmobility/github-api-service/common/client"
	"github.com/kenmobility/github-api-service/internal/domains/models"
	"github.com/kenmobility/github-api-service/internal/domains/services"
	"github.com/rs/zerolog/log"
)

type GitHubClient struct {
	baseURL                   string
	token                     string
	fetchInterval             time.Duration
	commitRepository          services.CommitRepository
	gitRepoMetadataRepository services.RepoMetadataRepository
	client                    *client.RestClient
	rateLimitFields           rateLimitFields
}

type rateLimitFields struct {
	rateLimitLimit     int
	rateLimitRemaining int
	rateLimitReset     int64
}

func (g *GitHubClient) getHeaders() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", g.token),
	}
}

func NewGitHubClient(baseUrl string, token string, fetchInterval time.Duration,
	commitRepository services.CommitRepository, gitRepoMetadataRepository services.RepoMetadataRepository) GitManagerClient {
	client := client.NewRestClient()

	gc := GitHubClient{
		baseURL:                   baseUrl,
		token:                     token,
		fetchInterval:             fetchInterval,
		commitRepository:          commitRepository,
		gitRepoMetadataRepository: gitRepoMetadataRepository,
		client:                    client,
	}
	ts := GitManagerClient(&gc)
	return ts
}

func (g *GitHubClient) FetchRepoMetadata(ctx context.Context, repositoryName string) (*models.RepoMetadata, error) {
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

	repoMetadata := &models.RepoMetadata{
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

func (g *GitHubClient) FetchCommits(ctx context.Context, repo models.RepoMetadata, since time.Time, until time.Time, lastFetchedCommit string) ([]models.Commit, error) {
	var cc []models.Commit
	var endpoint string

	if lastFetchedCommit != "" {
		endpoint = fmt.Sprintf("%s/repos/%s/commits?sha=%s", g.baseURL, repo.Name, lastFetchedCommit)
	} else {
		endpoint = fmt.Sprintf("%s/repos/%s/commits?since=%s&until=%s", g.baseURL, repo.Name, since.Format(time.RFC3339), until.Format(time.RFC3339))
	}
	for endpoint != "" {
		commitRes, nextURL, err := g.fetchCommitsPage(endpoint)
		if err != nil {
			return nil, err
		}

		for _, cr := range commitRes {
			commit := models.Commit{
				CommitID:       cr.SHA,
				Message:        cr.Commit.Message,
				Author:         cr.Commit.Author.Name,
				Date:           cr.Commit.Author.Date,
				URL:            cr.HtmlURL,
				RepositoryName: repo.Name,
			}

			cc = append(cc, commit)

			sc, err := g.commitRepository.SaveCommit(ctx, commit)
			if err != nil {
				log.Error().Msgf("Error saving commitId-%s: %v\n", commit.CommitID, err)
				continue
			}
			lastFetchedCommit = sc.CommitID
		}

		repo.LastFetchedCommit = lastFetchedCommit
		_, err = g.gitRepoMetadataRepository.UpdateRepoMetadata(ctx, repo)
		if err != nil {
			log.Error().Msgf("Error updating repository %s: %v", repo.Name, err)
			continue
		}
		endpoint = nextURL
	}

	return cc, nil
}

func (g *GitHubClient) fetchCommitsPage(url string) ([]GithubCommitResponse, string, error) {

	response, err := g.client.Get(url, map[string]string{}, g.getHeaders())
	if err != nil {
		log.Error().Msgf("error fetching commits: %v", err)

		return nil, "", err
	}

	if response.StatusCode == http.StatusForbidden {
		return nil, "", fmt.Errorf("rate limit exceeded")
	}

	g.updateRateLimitHeaders(response)

	if g.rateLimitFields.rateLimitRemaining == 0 {
		waitTime := time.Until(time.Unix(g.rateLimitFields.rateLimitReset, 0))
		log.Info().Msgf("Rate limit exceeded. Waiting for %v until reset...", waitTime)
	}

	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to fetch commits; status code: %v", response.StatusCode)
	}

	var commitRes []GithubCommitResponse

	if err := json.Unmarshal([]byte(response.Body), &commitRes); err != nil {
		fmt.Printf("marshal error, [%v]", err)
		return nil, "", errors.New("could not unmarshal commits response")
	}

	nextURL := g.parseNextURL(response.Headers["Link"])

	return commitRes, nextURL, nil
}

func (api *GitHubClient) parseNextURL(linkHeader []string) string {
	if len(linkHeader) == 0 {
		return ""
	}

	links := strings.Split(linkHeader[0], ",")
	for _, link := range links {
		parts := strings.Split(strings.TrimSpace(link), ";")
		if len(parts) < 2 {
			continue
		}

		urlPart := strings.Trim(parts[0], "<>")
		relPart := strings.TrimSpace(parts[1])

		if relPart == `rel="next"` {
			return urlPart
		}
	}

	return ""
}

func (api *GitHubClient) updateRateLimitHeaders(resp *client.Response) {
	limit := resp.Headers["X-Ratelimit-Limit"]

	if len(limit) > 0 {
		api.rateLimitFields.rateLimitReset, _ = strconv.ParseInt(limit[0], 10, 64)
	}

	remaining := resp.Headers["X-Ratelimit-Remaining"]

	if len(remaining) > 0 {
		api.rateLimitFields.rateLimitRemaining, _ = strconv.Atoi(remaining[0])
		log.Info().Msgf("Rate limit remaining: %d", api.rateLimitFields.rateLimitRemaining)
	}

	reset := resp.Headers["X-Ratelimit-Reset"]
	if len(reset) > 0 {
		api.rateLimitFields.rateLimitReset, _ = strconv.ParseInt(reset[0], 10, 64)
	}

	used := resp.Headers["X-Ratelimit-Used"]
	if len(used) > 0 {
		usedInt, _ := strconv.Atoi(used[0])
		log.Info().Msgf("Rate limit used: %d/%d", usedInt, api.rateLimitFields.rateLimitLimit)
	}
}
