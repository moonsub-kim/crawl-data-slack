package main

import (
	"errors"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/githubclient"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var (
	githubCreateIssueArgTitle  = "title"
	githubCreateIssueArgText   = "text"
	githubCreateIssueArgImage  = "image-file"
	githubCreateIssueArgFile   = "file"
	githubCreateIssueArgLabels = "labels"

	commandGithubCreateIssue *cli.Command = &cli.Command{
		Name: "create-issue",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: githubCreateIssueArgTitle, Required: true},
			&cli.StringFlag{Name: githubCreateIssueArgText, Required: true},
			&cli.StringFlag{Name: githubCreateIssueArgImage},
			&cli.StringFlag{Name: githubCreateIssueArgFile},
			&cli.StringSliceFlag{Name: githubCreateIssueArgLabels},
		},
		Action: RunGithub(
			func(ctx *cli.Context, logger *zap.Logger, client *githubclient.Client) error {
				title := ctx.String(githubCreateIssueArgTitle)
				text := ctx.String(githubCreateIssueArgText)
				imageFile := ctx.String(githubCreateIssueArgImage)
				fname := ctx.String(githubCreateIssueArgFile)
				labels := ctx.StringSlice(githubCreateIssueArgLabels)
				if labels == nil {
					labels = []string{}
				}

				// github에 존재하는 label 가져옴
				existLabels, err := client.ListLabels()
				if err != nil {
					return err
				}

				// github에 존재하지 않는 label들 추가
				for _, l := range labels {
					if _, ok := existLabels[l]; !ok {
						err = client.CreateLabel(l)
						if err != nil {
							return err
						}
					}
				}

				if fname != "" && imageFile != "" {
					return errors.New("file, image 동시 사용 금지")
				}

				var file crawler.File
				if fname != "" || imageFile != "" {
					if imageFile != "" {
						file = crawler.File{
							Path:    fname,
							IsImage: true,
						}
					} else {
						file = crawler.File{Path: fname}
					}
				}

				return client.CreatePost(crawler.Post{
					Title:  title,
					Labels: labels,
					Bodies: []crawler.Body{
						{
							Text:  text,
							Files: []crawler.File{file},
						},
					},
				})
			},
		),
	}
)
