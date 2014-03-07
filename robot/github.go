package robot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// The JSON object returned from a GET to the github API for repositories
// curl -i https://api.github.com/repos/kkochis/adoptmyapp
// TODO Are owners always Users?
type Repository struct {
	URL         string    `json:"url"` // Normalized URL
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Owner       *User     `json:"owner"`
	HtmlURL     string    `json:"html_url"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"private"`
	IsFork      bool      `json:"fork"`
	Created     time.Time `json:"created_at"`
	Updated     time.Time `json:"updated_at"`
	Subscribers int64     `json:"subscribers"`
	// TODO Other fields
}

func ParseRepository(data []byte) (*Repository, error) {
	var repo Repository
	if err := json.Unmarshal(data, &repo); err != nil {
		return nil, err
	}
	return &repo, nil
}

// The JSON object contained in the owner field of the repository API object
type User struct {
	Login     string `json:"login"`
	Id        int64  `json:"id"`
	HtmlURL   string `json:"html_url"`
	AvatarURL string `json:"avatar_url"`
	// TODO Other fields
}

// TODO Get additional collaborators, commits, etc...

// Normalize GitHub URL
// * Convert the Scheme to HTTPS
// * Normalize the Host to github.com
// * Remove any extraneous information at the end of the url
// This normalized URL will be used to check uniqueness in the database
func NormalizeGitHubURL(rawurl string) (*url.URL, error) {
	// SSH links are valid
	if strings.HasPrefix(rawurl, `git@github.com:`) {
		return NormalizeGitHubSSH(rawurl)
	}

	// Is it even a URL?
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	// Either the first part of the path or the host must be github
	if !(parsed.Host == `www.github.com` || parsed.Host == `github.com` || strings.HasPrefix(parsed.Path, `www.github.com`) || strings.HasPrefix(parsed.Path, `github.com`)) {
		return nil, fmt.Errorf("%q does not appear to be a GitHub URL", rawurl)
	}
	// Normalize the host
	parsed.Host = `github.com`

	// Normalize the scheme to HTTPS
	parsed.Scheme = `https`

	// Split the Path apart
	// There should be a blank part, a user, and a repo
	parts := strings.Split(parsed.Path, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("%q must point to a repository", rawurl)
	}

	// Remove any .git suffix from the repository part
	var repo = parts[2]
	if strings.HasSuffix(repo, `.git`) {
		repo = repo[:len(repo)-4]
	}

	parsed.Path = strings.Join([]string{"", parts[1], repo}, "/")

	// Clear any extraneous query or fragments from the URL
	parsed.RawQuery = ""
	parsed.Fragment = ""

	return parsed, nil
}

// Convert a GitHub SSH link into a URL
func NormalizeGitHubSSH(rawssh string) (*url.URL, error) {
	parts := strings.Split(rawssh[15:], "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("%q does not appear to be a repository", rawssh)
	}

	var repo = parts[1]
	if strings.HasSuffix(repo, `.git`) {
		repo = repo[:len(repo)-4]
	}

	// Recreate a URL from this information
	parsed := &url.URL{
		Scheme: `https`,
		Host:   `github.com`,
		Path:   strings.Join([]string{"", parts[0], repo}, "/"),
	}
	return parsed, nil
}

// Convert the given repository URL into an API URL
// https://github.com/kkochis/adoptmyapp
// to
// https://api.github.com/repos/kkochis/adoptmyapp
// TODO Error checking, what if a github.com URL was not given
func ConvertRepoURL(rawurl string) (string, error) {
	parsed, err := NormalizeGitHubURL(rawurl)
	if err != nil {
		return "", err
	}

	// Replace the host with the GitHub API Host
	parsed.Host = `api.github.com`

	// Add repos to the path
	parts := strings.Split(parsed.Path, "/")
	parsed.Path = strings.Join([]string{"", "repos", parts[1], parts[2]}, "/")

	return parsed.String(), nil
}

func GetRepoInfo(repositoryURL string) (*Repository, error) {
	// Convert the repository URL into an API URL
	apiURL, err := ConvertRepoURL(repositoryURL)
	if err != nil {
		return nil, err
	}
	_, err = http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
