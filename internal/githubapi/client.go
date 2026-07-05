package githubapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const apiBaseURL = "https://api.github.com"

var ErrNotFound = errors.New("github resource not found")

type Client struct {
	httpClient *http.Client
	token      string
	userAgent  string
}

type Release struct {
	TagName string         `json:"tag_name"`
	Name    string         `json:"name"`
	Assets  []ReleaseAsset `json:"assets"`
}

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	ContentType        string `json:"content_type"`
	Size               int64  `json:"size"`
}

func New(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		token:      strings.TrimSpace(token),
		userAgent:  "sing-box-rules-starter/1.0",
	}
}

func (c *Client) LatestRelease(ctx context.Context, repository string) (*Release, error) {
	owner, repo, err := splitRepository(repository)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", apiBaseURL, owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var release Release
	if err := c.doJSON(req, &release); err != nil {
		return nil, err
	}

	return &release, nil
}

func (c *Client) Download(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doDownload(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func FindAsset(release *Release, name string) *ReleaseAsset {
	for _, asset := range release.Assets {
		if asset.Name == name {
			assetCopy := asset
			return &assetCopy
		}
	}

	return nil
}

func (c *Client) doJSON(req *http.Request, output any) error {
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(output)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/vnd.github+json")
	return c.doWithHeaders(req)
}

func (c *Client) doDownload(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "*/*")
	return c.doWithHeaders(req)
}

func (c *Client) doWithHeaders(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", c.userAgent)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, ErrNotFound
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		resp.Body.Close()
		return nil, fmt.Errorf("github request failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	return resp, nil
}

func splitRepository(repository string) (string, string, error) {
	parts := strings.Split(strings.TrimSpace(repository), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repository %q", repository)
	}

	return parts[0], parts[1], nil
}
