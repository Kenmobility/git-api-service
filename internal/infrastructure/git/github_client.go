package git

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kenmobility/github-api-service/common/client"
	"github.com/kenmobility/github-api-service/internal/domains/models"
	"github.com/kenmobility/github-api-service/internal/domains/services"
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
		fmt.Printf("marshal error, [%v]", err)
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

func (g *GitHubClient) FetchCommits(ctx context.Context, repoName string, since time.Time, until time.Time) ([]models.Commit, error) {
	var result []models.Commit

	endpoint := fmt.Sprintf("%s/repos/%s/commits?since=%s&until=%s", g.baseURL, repoName, since.Format(time.RFC3339), until.Format(time.RFC3339))
	for endpoint != "" {

		commitRes, nextURL, err := g.fetchCommitsPage(endpoint)
		if err != nil {
			return nil, err
		}

		for _, c := range commitRes {
			result = append(result, models.Commit{
				CommitID:       c.SHA,
				Message:        c.Commit.Message,
				Author:         c.Commit.Author.Name,
				Date:           c.Commit.Author.Date,
				URL:            c.HtmlURL,
				RepositoryName: repoName,
			})
		}
		/*
			for _, commit := range result {
				commit.RepositoryName = repo.Name

				_, err := g.commitRepository.SaveCommit(ctx, commit)
				if err != nil {
					log.Printf("Error saving commit id-%s: %v\n", commit.CommitID, err)
				}
			}
		*/
		endpoint = nextURL
	}

	return result, nil
}

/*
func (g GitHubClient) StartTracking(ctx context.Context, fetchInterval time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			g.runRepositoryTracker(ctx)
			time.Sleep(fetchInterval)
		}
	}
}


func (g GitHubClient) runRepositoryTracker(ctx context.Context) {

	trackedRepo, err := g.repositoryRepository.TrackedRepository(ctx)
	if err != nil {
		log.Printf("Error fetching tracked repository: %v", err)
		return
	}

	if trackedRepo == nil {
		log.Println("no repository set to track")
		return
	}
	fmt.Printf("********Github repository tracking started for repo %s************\n",
		trackedRepo.Name)
	g.FetchAndSaveCommits(ctx, *trackedRepo, trackedRepo.StartDate, trackedRepo.EndDate)
}
*/

func (g *GitHubClient) fetchCommitsPage(url string) ([]GithubCommitResponse, string, error) {

	response, err := g.client.Get(url, map[string]string{}, g.getHeaders())
	if err != nil {
		log.Println("error fetching commits: ", err)
		return nil, "", err
	}

	if response.StatusCode == http.StatusForbidden {
		return nil, "", fmt.Errorf("*************rate limit exceeded*************")
	}

	g.updateRateLimitHeaders(response)

	if g.rateLimitFields.rateLimitRemaining == 0 {
		waitTime := time.Until(time.Unix(g.rateLimitFields.rateLimitReset, 0))
		log.Printf("Rate limit exceeded. Waiting for %v until reset...", waitTime)
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
	}

	reset := resp.Headers["X-Ratelimit-Reset"]
	if len(reset) > 0 {
		api.rateLimitFields.rateLimitReset, _ = strconv.ParseInt(reset[0], 10, 64)
	}

	used := resp.Headers["X-Ratelimit-Used"]
	if len(used) > 0 {
		usedInt, _ := strconv.Atoi(used[0])
		log.Printf("Rate limit used: %d/%d", usedInt, api.rateLimitFields.rateLimitLimit)
	}
}
