package repository

import (
	"errors"
	"time"

	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/logger"
	"gorm.io/gorm"
)

type Repository struct {
	logger logger.Logger
	db     *gorm.DB
	mapper mapper
}

func (r Repository) GetEvents(from time.Time) ([]crawler.Event, error) {
	var events []Event
	err := r.db.Where("created_at >= ?", from).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return r.mapper.mapModelEventsToEvents(events), nil
}

func (r Repository) SaveEvents(events []crawler.Event) error {
	if len(events) == 0 {
		return nil
	}

	modelEvents := r.mapper.mapEventsToModelEvents(events)
	return r.db.Create(&modelEvents).Error
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
	err := r.db.First(&user, "user_name = ?", userName).Error
	if errors.As(err, &gorm.ErrRecordNotFound) {
		return crawler.User{}, nil
	} else if err != nil {
		return crawler.User{}, err
	}

	return r.mapper.mapModelUserToUser(user), nil
}

func (r Repository) SaveUsers(users []crawler.User) error {
	modelUsers := r.mapper.mapUsersToModelUsers(users)
	return r.db.Create(&modelUsers).Error
}

func NewRepository(logger logger.Logger, db *gorm.DB) *Repository {
	return &Repository{
		logger: logger,
		db:     db,
	}
}
