package slackclient

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type Client struct {
	logger *zap.Logger
	api    *slack.Client
	token  string
	mapper mapper
}

var imageExt = map[string]struct{}{"jpg": {}, "jpeg": {}, "png": {}}
var maxIteration int = 100
var filters []ArchiveFilter = []ArchiveFilter{
	messageSubTypeFilter{},                                               // negative filters, positive보다 우선함
	IsUserMessageFilter{}, IsUserReactedFilter{}, IsUserThreadedFilter{}, // positive filters
}

func (c Client) getConversations(channelID string, from time.Time, to time.Time) ([]slack.Message, error) {
	withRetry := func(params *slack.GetConversationHistoryParameters) (res *slack.GetConversationHistoryResponse, err error) {
		err = retry.Do(func() error {
			var err error
			res, err = c.api.GetConversationHistory(params)
			return err
		})
		c.logger.Debug("getConverstaions", zap.Error(err), zap.Any("params", params), zap.Any("results", res))
		return res, err
	}

	var err error
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Oldest:    strconv.FormatInt(from.Unix(), 10),
		Latest:    strconv.FormatInt(to.Unix(), 10),
		Limit:     100,
	}
	res := &slack.GetConversationHistoryResponse{HasMore: true}

	var messages []slack.Message
	for i := 0; i < maxIteration && res.HasMore; i++ {
		res, err = withRetry(params)
		if err != nil {
			return nil, err
		}
		messages = append(messages, res.Messages...)
		params.Cursor = res.ResponseMetaData.NextCursor
	}

	return messages, nil
}

func (c Client) getThreadMessages(channelID string, m slack.Message) ([]slack.Message, error) {
	withRetry := func(params *slack.GetConversationRepliesParameters) (msgs []slack.Message, hasMore bool, nextCursor string, err error) {
		err = retry.Do(func() error {
			var err error
			msgs, hasMore, nextCursor, err = c.api.GetConversationReplies(params)
			c.logger.Debug("getThreadMessages", zap.Error(err))
			return err
		})
		return msgs, hasMore, nextCursor, err
	}

	var err error
	var nextCursor string
	hasMore := true

	params := &slack.GetConversationRepliesParameters{
		ChannelID: channelID,
		Timestamp: m.Timestamp,
	}

	var messages, resMessages []slack.Message
	for i := 0; i < maxIteration && hasMore; i++ {
		resMessages, hasMore, nextCursor, err = withRetry(params)
		if err != nil {
			return nil, err
		}
		messages = append(messages, resMessages...)
		params.Cursor = nextCursor
	}

	return messages, nil
}

func (c Client) archiveFilter(messages []slack.Message) bool {
	for _, f := range filters {
		if !f.Positive() && !f.Passed(messages) {
			return false
		} else if f.Positive() && f.Passed(messages) {
			return true
		}
	}

	return false
}

// 쓰레드 메시지들을 post로 만들때 다시 구조적으로 정리해야됨 ..

func (c Client) buildPost(channel crawler.Channel, messages []slack.Message) (crawler.Post, error) {
	var title string
	if len(messages[0].Attachments) > 0 {
		title = messages[0].Attachments[0].Title
	} else {
		// 링크 없는 케이스
		splitted := strings.Split(messages[0].Text, "\n")
		for _, l := range splitted {
			if l != "" {
				slackURLWithAliasPattern := regexp.MustCompile(`<(.+?)\|(.+?)>`)
				title = slackURLWithAliasPattern.ReplaceAllString(l, "$2")
				break
			}
		}
	}

	body, err := c.messageToBody(messages[0])
	if err != nil {
		return crawler.Post{}, err
	}

	post := crawler.Post{
		Title:  title,
		Labels: []string{channel.Name}, // TODO label 추가
		Bodies: []crawler.Body{body},
	}
	if len(messages) == 1 {
		return post, nil
	}

	for _, m := range messages[1:] {
		body, err := c.messageToBody(m)
		if err != nil {
			return crawler.Post{}, err
		}

		post.Comments = append(post.Comments, crawler.Comment{
			Bodies: []crawler.Body{body},
		})
	}

	return post, nil
}

func (c Client) ArchivePosts(channel crawler.Channel, from time.Time, to time.Time) ([]crawler.Post, error) {
	messageResponses, err := c.getConversations(channel.ID, from, to)
	if err != nil {
		return nil, err
	}

	var posts []crawler.Post
	var messages [][]slack.Message
	for _, m := range messageResponses {
		replies, err := c.getThreadMessages(channel.ID, m)
		if err != nil {
			return nil, err
		}

		if !c.archiveFilter(replies) {
			continue
		}

		post, err := c.buildPost(channel, replies)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)

		replies[0].BotProfile = nil
		replies[0].Blocks.BlockSet = nil

		messages = append(messages, replies)
	}

	c.logger.Info(
		"ArchivePosts",
		zap.Any("threads", messages),
		zap.Any("posts", posts),
	)

	return posts, nil
}

func (c Client) Notify(n crawler.Notification) error {
	withRetry := func(channelID string, opts ...slack.MsgOption) error {
		return retry.Do(func() error {
			time.Sleep(time.Second)
			_, _, err := c.api.PostMessage(channelID, opts...)
			c.logger.Debug("PostMessage", zap.Error(err))
			return err
		})
	}

	// short message
	if len(n.Event.Message) < 10000 {
		err := withRetry(n.User.ID, slack.MsgOptionText(n.Event.Message, false))
		c.logger.Info(
			"notify",
			zap.Any("notification", n),
			zap.Any("err", err),
		)
		return err
	}

	lines := strings.Split(n.Event.Message, "\n")
	for from := 0; from < len(lines); from += 6 {
		to := from + 6
		if to > len(lines) {
			to = len(lines)
		}

		text := strings.Join(lines[from:to], "\n")
		err := withRetry(n.User.ID, slack.MsgOptionText(text, false))
		c.logger.Info(
			"notify",
			zap.Any("notification", n),
			zap.Any("err", err),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c Client) GetChannels() ([]crawler.Channel, error) {
	slackUsers, err := c.api.GetUsers()
	if err != nil {
		c.logger.Error("getUsers", zap.Error(err))
		return nil, err
	}

	var activeUsers []slack.User
	for _, u := range slackUsers {
		if u.Deleted || u.IsBot || u.IsRestricted {
			continue
		}

		activeUsers = append(activeUsers, u)
	}

	users := c.mapper.mapSlackUsersToUsers(activeUsers)

	nextCursor := ""
	for {
		var slackChannels []slack.Channel
		param := slack.GetConversationsParameters{Cursor: nextCursor, ExcludeArchived: true}
		slackChannels, nextCursor, err = c.api.GetConversations(&param)
		if err != nil {
			c.logger.Error("getConversations", zap.Error(err))
			return nil, err
		}
		channels := c.mapper.mapSlackChannelsToUsers(slackChannels)
		users = append(users, channels...)
		c.logger.Info(
			"GetConversations",
			zap.Any("channels", channels),
			zap.Any("nextCursor", nextCursor),
		)
		if nextCursor == "" {
			break
		}
		time.Sleep(time.Second * 3)
	}

	return users, nil
}

func (c Client) messageToBody(message slack.Message) (crawler.Body, error) {
	slackURLWithAliasPattern := regexp.MustCompile(`<(.+?)\|(.+?)>`)
	slackURLPattern := regexp.MustCompile(`<(.+?)>`)

	// 슬랙 링크 패턴을 마크다운 링크 패턴으로 교체
	s := message.User + "\n" + slackURLWithAliasPattern.ReplaceAllString(message.Text, "[$2]($1)")
	s = slackURLPattern.ReplaceAllString(s, "[$1]($1)")

	for _, a := range message.Attachments {
		s = strings.ReplaceAll(
			s,
			fmt.Sprintf("[%s]", a.OriginalURL),
			fmt.Sprintf("[%s]", a.Title),
		) // link name에도 링크가 붙어있으면 title로 교체

		imageURL := a.ThumbURL
		if imageURL == "" {
			imageURL = a.ImageURL
		}
		s = strings.ReplaceAll(
			s,
			fmt.Sprintf("(%s)", a.OriginalURL),
			fmt.Sprintf("(%s)\n<image alt=\"%s\" src=\"%s\">\n", a.OriginalURL, a.Title, imageURL),
		) // link 밑에 이미지 넣기
	}

	return crawler.Body{Text: s}, nil

	// 여기서 파일 저장하는게 맘에안듦
	// fileDownloadTextPattern := regexp.MustCompile(`a href="(.+?)"`)
	// var files []crawler.File
	// for _, file := range message.Files {
	// 	req, err := http.NewRequest(http.MethodGet, file.URLPrivateDownload, nil)
	// 	if err != nil {
	// 		return crawler.Body{}, err
	// 	}
	// 	req.Header.Add("Authorization", "Bearer "+c.token)

	// 	res, err := http.DefaultClient.Do(req)
	// 	if err != nil {
	// 		return crawler.Body{}, err
	// 	}
	// 	defer res.Body.Close()

	// 	b, err := io.ReadAll(res.Body)
	// 	if err != nil {
	// 		return crawler.Body{}, err
	// 	}

	// 	c.logger.Info(
	// 		"req",
	// 		zap.String("url", file.URLPrivateDownload),
	// 		zap.String("res", string(b)),
	// 	)

	// 	submatches := fileDownloadTextPattern.FindStringSubmatch(string(b))
	// 	if len(submatches) != 2 {
	// 		return crawler.Body{}, errors.New("len(submatches) != 2")
	// 	}

	// 	// submatches[1]
	// 	req, err = http.NewRequest(http.MethodGet, submatches[1], nil)
	// 	if err != nil {
	// 		return crawler.Body{}, err
	// 	}
	// 	req.Header.Add("Authorization", "Bearer "+c.token)

	// 	res, err = http.DefaultClient.Do(req)
	// 	if err != nil {
	// 		return crawler.Body{}, err
	// 	}
	// 	defer res.Body.Close()

	// 	c.logger.Info(
	// 		"req2",
	// 		zap.String("url", submatches[0]),
	// 	)

	// 	splitted := strings.Split(submatches[1], "%2f")
	// 	name := splitted[len(splitted)-1]
	// 	err = writeFile(name, res.Body)
	// 	if err != nil {
	// 		return crawler.Body{}, err
	// 	}

	// 	c.logger.Info("sleep")
	// 	time.Sleep(time.Hour)

	// 	splitted = strings.Split(name, ".")
	// 	ext := splitted[len(splitted)-1]
	// 	_, isImage := imageExt[ext]

	// 	files = append(files, crawler.File{
	// 		Path:    name,
	// 		IsImage: isImage,
	// 	})
	// }

	// return crawler.Body{
	// 	Text:  s,
	// 	Files: files,
	// }, nil
}

func NewClient(logger *zap.Logger, client *slack.Client, token string) *Client {
	return &Client{
		logger: logger,
		api:    client,
		token:  token,
	}
}

func writeFile(name string, body io.Reader) error {
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, body)
	if err != nil {
		return err
	}

	return nil
}
