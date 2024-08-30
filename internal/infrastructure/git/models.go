package git

import "time"

type (
	GithubCommitResponse struct {
		SHA     string `json:"sha"`
		NodeId  string `json:"node_id"`
		Commit  Commit `json:"commit"`
		URL     string `json:"url"`
		HtmlURL string `json:"html_url"`
	}

	Commit struct {
		Author  Author `json:"author"`
		Message string `json:"message"`
		URL     string `json:"url"`
	}

	Author struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	}
)

type (
	GitHubRepoMetadataResponse struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		HtmlUrl     string `json:"html_url"`
		Description string `json:"description"`
		Url         string `json:"url"`
		Owner       struct {
			Url     string `json:"url"`
			HtmlUrl string `json:"html_url"`
		} `json:"owner"`
		StargazersCount int    `json:"stargazers_count"`
		WatchersCount   int    `json:"watchers_count"`
		Language        string `json:"language"`
		ForksCount      int    `json:"forks_count"`
		OpenIssues      int    `json:"open_issues"`
	}
)
