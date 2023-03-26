package slackclient

import (
	"regexp"
	"strings"

	"github.com/slack-go/slack"
)

type ArchiveFilter interface {
	Positive() bool
	Passed(messages []slack.Message) bool
}

// positiveFilter accepts the Passed() value if true, or call the next filter if false
type positiveFilter struct{}

func (f positiveFilter) Positive() bool {
	return true
}

// positiveFilter denies the Passed() value if false, or call the next filter if true
type negativeFilter struct{}

func (f negativeFilter) Positive() bool {
	return false
}

// positive filters

type IsUserMessageFilter struct {
	positiveFilter
}

func (f IsUserMessageFilter) Passed(messages []slack.Message) bool {
	return messages[0].BotID == ""
}

type IsUserReactedFilter struct {
	positiveFilter
}

func (f IsUserReactedFilter) Passed(messages []slack.Message) bool {
	return len(messages[0].Reactions) > 0
}

type IsUserThreadedFilter struct {
	positiveFilter
}

// reply중에 user가 단게 있을 때
func (f IsUserThreadedFilter) Passed(messages []slack.Message) bool {
	if len(messages) == 1 {
		return false
	}

	for _, m := range messages[1:] {
		if m.BotID == "" {
			return true
		}
	}

	return false
}

// negative filters

// subtype이 들어있으면 channel_joined, thread_broadcast 등 특수목적 메시지임
type MessageSubTypeExistsFilter struct {
	negativeFilter
}

func (f MessageSubTypeExistsFilter) Passed(messages []slack.Message) bool {
	return messages[0].SubType == ""
}

// message에 link가 없으면 제외
type NoLinkFilter struct {
	negativeFilter
}

func (f NoLinkFilter) Passed(messages []slack.Message) bool {
	p := regexp.MustCompile("<.+>")
	return len(messages[0].Attachments) > 0 || // attach (thumbnail형태)로 link 존재
		len(p.FindString(messages[0].Text)) > 0 // thumbnail이 없는 link
}

// 특정 이모지가 있으면 제외
type ExcludeEmojiFilter struct {
	negativeFilter
	emojiName string
}

func (f ExcludeEmojiFilter) Passed(messages []slack.Message) bool {
	for _, reaction := range messages[0].Reactions {
		if strings.ToLower(reaction.Name) == f.emojiName {
			return false
		}
	}

	return true
}

func NewExcludeEmojiFilter(emojiName string) *ExcludeEmojiFilter {
	return &ExcludeEmojiFilter{
		emojiName: emojiName,
	}
}
