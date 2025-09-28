package crawler_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	mock_crawler "github.com/moonsub-kim/crawl-data-slack/mocks/internal_/pkg/crawler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type testFixture struct {
	Logger    *zap.Logger
	Repo      *mock_crawler.Repository
	Crawler   *mock_crawler.Crawler
	Messenger *mock_crawler.Messenger
	Archive   *mock_crawler.Archive
	UseCase   *crawler.UseCase
}

// Setup 메서드: 환경 설정 및 객체 생성 로직을 포함
func setup(t *testing.T) *testFixture {
	l := zap.NewNop()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	l = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		level,
	))

	r := new(mock_crawler.Repository)
	c := new(mock_crawler.Crawler)
	m := new(mock_crawler.Messenger)
	a := new(mock_crawler.Archive)

	u := crawler.NewUseCase(l, r, c, m, a)

	return &testFixture{
		Logger:    l,
		Repo:      r,
		Crawler:   c,
		Messenger: m,
		Archive:   a,
		UseCase:   u,
	}
}

func (f *testFixture) assert(t *testing.T) {
	f.Repo.AssertExpectations(t)
	f.Crawler.AssertExpectations(t)
	f.Messenger.AssertExpectations(t)
	f.Archive.AssertExpectations(t)
}

func TestWork(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		f := setup(t)
		defer f.assert(t)

		events := []crawler.Event{{Crawler: "crawler", UserName: "test", EventTime: time.Now()}}

		f.Crawler.On("Crawl").Return(events, nil).Once()
		for _, e := range events {
			user := crawler.Channel{ID: e.UserName, Name: e.UserName}
			notification := crawler.Notification{Event: e, User: user}
			f.Repo.On("SaveEvent", e).Return(nil)
			f.Repo.On("GetChannel", e.UserName).Return(user, nil)
			f.Messenger.On("Notify", notification).Return(nil)
		}

		err := f.UseCase.Work(time.Now().Add(-time.Hour))
		assert.NoError(t, err)
	})
}

func TestArchive(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		f := setup(t)
		defer f.assert(t)

		channel := crawler.Channel{ID: "id", Name: "name"}
		dateFrom := time.Now().Add(-time.Hour)
		dateTo := time.Now()
		posts := []crawler.Post{
			{
				Title:    "a",
				Labels:   []string{"old", "new"},
				Bodies:   []crawler.Body{},
				Comments: []crawler.Comment{},
			},
			{
				Title:    "b",
				Labels:   []string{"old", "new"},
				Bodies:   []crawler.Body{},
				Comments: []crawler.Comment{},
			},
		}
		existLabels := map[string]struct{}{"old": {}}
		newLabels := []string{"new"}

		f.Repo.On("GetChannel", channel.Name).Return(channel, nil)
		f.Messenger.On("ArchivePosts", channel, dateFrom, dateTo).Return(posts, nil)
		f.Archive.On("ListLabels").Return(existLabels, nil)
		for _, l := range newLabels {
			f.Archive.On("CreateLabel", l).Return(nil)
		}
		f.Archive.On("CreatePosts", posts).Return(nil)

		err := f.UseCase.Archive(channel.Name, dateFrom, dateTo)
		assert.NoError(t, err)
	})
}

func TestSyncLabel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		f := setup(t)
		defer f.assert(t)

		labels := []string{"a", "b"}
		f.Messenger.On("GetLabels").Return(labels, nil)
		f.Archive.On("SyncLabels", labels).Return(nil)

		err := f.UseCase.SyncLabel()
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		f := setup(t)
		defer f.assert(t)

		e := errors.New("")
		f.Messenger.On("GetLabels").Return(nil, e)

		err := f.UseCase.SyncLabel()
		assert.Error(t, err)
	})
}
