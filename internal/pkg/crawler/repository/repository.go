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

func (r Repository) SaveEvent(event crawler.Event) error {
	e := r.mapper.mapEventToModelEvent(event)
	err := r.db.Create(&e).Error
	if err != nil && strings.Contains(err.Error(), "1062") {
		return crawler.AlreadyExistsError{}
	} else if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetRestriction(c string, j string) (crawler.Restriction, error) {
	var restriction Restriction
	err := r.db.Where("crawler = ? AND job = ?", c, j).Order("created_at DESC").First(&restriction).Error
	if errors.As(err, &gorm.ErrRecordNotFound) {
		r.logger.Info("Ignore empty restriction")
		return crawler.Restriction{}, nil
	} else if err != nil {
		return crawler.Restriction{}, err
	}

	return r.mapper.mapModelRestrictionToRestriction(restriction), nil
}

func (r Repository) SaveRestriction(restriction crawler.Restriction) error {
	model := r.mapper.mapRestrictionToModelRestriction(restriction)
	return r.db.Create(&model).Error
}

func (r Repository) GetUser(userName string) (crawler.User, error) {
	var user User
	err := r.db.First(&user, "name = ?", userName).Error
	if errors.As(err, &gorm.ErrRecordNotFound) {
		return crawler.User{}, nil
	} else if err != nil {
		return crawler.User{}, err
	}

	return r.mapper.mapModelUserToUser(user), nil
}

func (r Repository) SaveUsers(users []crawler.User) error {
	modelUsers := r.mapper.mapUsersToModelUsers(users)
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
