package repository

import (
	"errors"
	"strings"

	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	logger *zap.Logger
	db     *gorm.DB
	mapper mapper
}

func (r Repository) isPostgresqkAlreadyExsitsError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

func (r Repository) isMysqlAlreadyExsitsError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "1062")
}

func (r Repository) SaveEvent(event crawler.Event) error {
	e := r.mapper.mapEventToModelEvent(event)
	err := r.db.Create(&e).Error
	if r.isMysqlAlreadyExsitsError(err) || r.isPostgresqkAlreadyExsitsError(err) {
		return crawler.AlreadyExistsError{}
	} else if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetChannel(userName string) (crawler.Channel, error) {
	var user Channel
	err := r.db.First(&user, "name = ?", userName).Error
	if errors.As(err, &gorm.ErrRecordNotFound) {
		return crawler.Channel{}, nil
	} else if err != nil {
		return crawler.Channel{}, err
	}

	return r.mapper.mapModelChannelToChannel(user), nil
}

func (r Repository) SaveChannels(users []crawler.Channel) error {
	modelUsers := r.mapper.mapChannelsToModelChannels(users)
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}}, // key colume
	}).Create(&modelUsers).Error
}

func NewRepository(logger *zap.Logger, db *gorm.DB) *Repository {
	return &Repository{
		logger: logger,
		db:     db,
	}
}
