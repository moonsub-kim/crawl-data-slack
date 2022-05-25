package spinnaker

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger     *zap.Logger
	httpClient *http.Client

	channel   string
	authToken string
	host      string
}

const GET_APPS_URL string = "%s/applications"                          // c.host
const GET_CONFIGS_URL string = "%s/applications/%s/pipelineConfigs"    // c.host, application
const GOTO_URL string = "%s/#/applications/%s/executions/configure/%s" // c.host, application, pipeline_id

var IGNORE_BRANCHES map[string]struct{} = map[string]struct{}{
	"":       {},
	"master": {},
	"${trigger['payload']['deployment']['sha']}": {},
}

func (c Crawler) GetCrawlerName() string { return "spinnaker" }
func (c Crawler) GetJobName() string     { return "values-branch" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var apps []App
	err := c.requestGet(fmt.Sprintf(GET_APPS_URL, c.host), &apps)
	if err != nil {
		return nil, err
	}

	var events []crawler.Event
	for _, app := range apps {
		var pipelines []Pipeline
		err := c.requestGet(fmt.Sprintf(GET_CONFIGS_URL, c.host, app.Name), &pipelines)
		if err != nil {
			c.logger.Error(
				"get pipeline error",
				zap.Error(err),
			)
			return nil, err
		}

		pipelineEvents, err := c.validatePipeline(pipelines)
		if err != nil {
			return nil, err
		}
		events = append(events, pipelineEvents...)
	}

	return events, nil
}

func (c Crawler) validatePipeline(pipelines []Pipeline) ([]crawler.Event, error) {
	getGithubArtifactFromStage := func(stage Stage) (Artifact, error) {
		for _, a := range stage.InputArtifacts {
			if a.Account == "github-artifact" {
				return a.Artifact, nil
			}
		}
		return Artifact{}, fmt.Errorf("no github artifact for the stage")
	}

	getValuesBranch := func(pipeline Pipeline) (string, error) {
		// custom templates
		for _, s := range pipeline.Stages {
			a, err := getGithubArtifactFromStage(s)
			if err != nil {
				continue
			}

			return a.Version, nil
		}

		// chart template v2
		branch, ok := pipeline.Variables["valueArtifactBranch"]
		if ok {
			return branch, nil
		}

		// chart template v1 always has master branch
		if pipeline.Type == "templatedPipeline" {
			return "master", nil
		}

		return "", fmt.Errorf("no github artifacte for the pipeline %s", pipeline.Name)
	}

	var events []crawler.Event
	for _, p := range pipelines {
		if p.Disabled {
			continue
		}

		branch, err := getValuesBranch(p)
		if err != nil {
			c.logger.Warn(
				"pipeline has no github artifact",
				zap.Any("pipeline", p),
			)
			continue
		}

		if _, ok := IGNORE_BRANCHES[branch]; !ok {
			events = append(
				events,
				crawler.Event{
					Crawler:  c.GetCrawlerName(),
					Job:      c.GetJobName(),
					UserName: c.channel,
					UID:      p.Name,
					Name:     time.Now().String(),
					Message: fmt.Sprintf(
						"[Spinnaker] '%s' of the Application '%s' values has a user branch `%s`.\n <%s|go to the pipeline>",
						p.Name,
						p.Application,
						branch,
						fmt.Sprintf(GOTO_URL, c.host, p.Application, p.ID),
					),
				},
			)
		}
	}

	return events, nil
}

func (c Crawler) requestGet(u string, res interface{}) error {
	url, err := url.Parse(u)
	if err != nil {
		return err
	}

	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	// header.Set("Cookie", c.session)

	req := &http.Request{
		Method: http.MethodGet,
		URL:    url,
		Header: header,
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &res)
}

func NewCrawler(logger *zap.Logger, channel string, host string, authToken string) (*Crawler, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		logger: logger,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // To avoid deniedUntrustedSSL
			},
			Jar: jar, // keep the session cookie
		},

		channel:   channel,
		authToken: authToken,
		host:      host,
	}, nil
}
