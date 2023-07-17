package githubclient

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/avast/retry-go"
	"github.com/google/go-github/v49/github"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Client struct {
	logger *zap.Logger
	client *github.Client
	owner  string
	repo   string
}

func (c Client) createIssue(title string, body string, labels []string) (int, error) {
	withRetry := func(req *github.IssueRequest) (issue *github.Issue, res *github.Response, err error) {
		err = retry.Do(func() error {
			issue, res, err = c.client.Issues.Create(
				context.Background(),
				c.owner,
				c.repo,
				req,
			)
			c.logger.Debug("createIssue", zap.Error(err))
			return err
		})
		return issue, res, err
	}

	if labels == nil {
		labels = []string{}
	}

	issue, _, err := withRetry(
		&github.IssueRequest{
			Title:  &title,
			Body:   &body,
			Labels: &labels,
		},
	)
	if err != nil {
		return 0, err
	}

	return *issue.Number, err
}

// idempotent
func (c Client) createRelease(releaseTag string, title string) (int64, error) {
	withRetry := func(req *github.RepositoryRelease) (release *github.RepositoryRelease, res *github.Response, err error) {
		err = retry.Do(func() error {
			release, res, err = c.client.Repositories.GetReleaseByTag(
				context.Background(),
				c.owner,
				c.repo,
				releaseTag,
			)
			if err == nil {
				return nil
			}

			release, res, err = c.client.Repositories.CreateRelease(
				context.Background(),
				c.owner,
				c.repo,
				req,
			)
			c.logger.Debug("createRelease", zap.Error(err), zap.Any("req", req))
			return err
		})
		return release, res, err
	}

	release, _, err := withRetry(&github.RepositoryRelease{
		TagName: &releaseTag,
		Name:    &title,
	})
	if err != nil {
		return 0, err
	}

	return *release.ID, nil
}

func (c Client) uploadAsset(releaseID int64, f crawler.File) (string, error) {
	withRetry := func() (asset *github.ReleaseAsset, res *github.Response, err error) {
		err = retry.Do(func() error {
			// os.File은 readcloser 라서 retry시 file closed 에러가 생김
			f, err := os.Open(f.Path)
			if err != nil {
				return err
			}

			asset, res, err = c.client.Repositories.UploadReleaseAsset(
				context.Background(),
				c.owner,
				c.repo,
				releaseID,
				&github.UploadOptions{Name: f.Name()},
				f,
			)
			c.logger.Debug("uploadAsset", zap.Error(err))
			return err
		})
		return asset, res, err
	}

	asset, _, err := withRetry()
	if err != nil {
		return "", err
	}

	return *asset.BrowserDownloadURL, nil
}

func (c Client) createBody(releaseTag string, title string, bodies []crawler.Body) (string, error) {
	var releaseID int64
	var err error
	bodyText := ""

	for _, b := range bodies {
		bodyText += fmt.Sprintf("%s\n", b.Text)

		for _, file := range b.Files {
			if releaseID == 0 {
				releaseID, err = c.createRelease(releaseTag, title)
				if err != nil {
					return "", err
				}
			}

			fileURL, err := c.uploadAsset(releaseID, file)
			if err != nil {
				return "", err
			}

			if file.IsImage {
				bodyText += fmt.Sprintf("<img alt=\"%s\" src=\"%s\">\n", file.Name(), fileURL)
			} else {
				bodyText += fmt.Sprintf("<a href=\"%s\">%s</a>", fileURL, file.Name())
			}
		}
	}

	return bodyText, nil
}

func (c Client) createIssueComment(issueNumber int, body string) error {
	withRetry := func(req *github.IssueComment) (comment *github.IssueComment, res *github.Response, err error) {
		err = retry.Do(func() error {
			comment, res, err = c.client.Issues.CreateComment(
				context.Background(),
				c.owner,
				c.repo,
				issueNumber,
				req,
			)
			c.logger.Debug("createIssueComment", zap.Error(err), zap.Any("req", req))
			return err
		})
		return comment, res, err
	}

	_, _, err := withRetry(&github.IssueComment{Body: &body})
	return err
}

func (c Client) CreatePost(post crawler.Post) error {
	// file attachment 가 있을경우 사용됨
	releaseTag := time.Now().Format("2006-01-02-15-04-05")
	// 다음 포스트에 영향을 주지 않도록 하기 위함
	time.Sleep(time.Second)

	body, err := c.createBody(releaseTag, post.Title, post.Bodies)
	if err != nil {
		return err
	}

	issueNumber, err := c.createIssue(post.Title, body, post.Labels)
	if err != nil {
		return err
	}

	for _, comment := range post.Comments {
		body, err := c.createBody(releaseTag, post.Title, comment.Bodies)
		if err != nil {
			return err
		}

		err = c.createIssueComment(issueNumber, body)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c Client) CreatePosts(posts []crawler.Post) error {
	for _, p := range posts {
		err := c.CreatePost(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c Client) CreateLabel(name string) error {
	withRetry := func(req *github.Label) (label *github.Label, res *github.Response, err error) {
		err = retry.Do(func() error {
			label, res, err = c.client.Issues.CreateLabel(
				context.Background(),
				c.owner,
				c.repo,
				req,
			)
			c.logger.Warn("CreateLabel", zap.Error(err))
			return err
		})
		return label, res, err
	}

	rand.Seed(time.Now().Unix())
	color := fmt.Sprintf("%x", rand.Int63n(0xFFFFFF)) // random color
	_, _, err := withRetry(&github.Label{
		Name:  &name,
		Color: &color,
	})
	return err
}

func (c Client) ListLabels() (map[string]struct{}, error) {
	withRetry := func() (labels []*github.Label, res *github.Response, err error) {
		err = retry.Do(func() error {
			labels, res, err = c.client.Issues.ListLabels(
				context.Background(),
				c.owner,
				c.repo,
				&github.ListOptions{PerPage: 100},
			)
			c.logger.Warn("ListLabels", zap.Error(err))
			return err
		})
		return labels, res, err
	}

	labels, _, err := withRetry()
	if err != nil {
		return nil, err
	}

	names := map[string]struct{}{}
	for _, l := range labels {
		names[*l.Name] = struct{}{}
	}

	return names, nil
}

func (c Client) SyncLabels(labels []string) error {
	currentLabels, err := c.ListLabels()
	if err != nil {
		return err
	}

	for _, l := range labels {
		if _, ok := currentLabels[l]; !ok {
			c.logger.Info("new_label", zap.String("new_label", l))
			err := c.CreateLabel(l)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewClient(logger *zap.Logger, client *github.Client, owner string, repo string) *Client {
	return &Client{
		logger: logger,
		client: client,
		owner:  owner,
		repo:   repo,
	}
}
