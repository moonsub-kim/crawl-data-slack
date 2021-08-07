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
	modelEvents := r.mapper.mapEventsToModelEvents(events)
	return r.db.Create(&modelEvents).Error
}

func (r Repository) GetRestriction(c string, j string) (crawler.Restriction, error) {
	var restriction Restriction
	err := r.db.Where("crawler = ? AND job = ?", c, j).Order("created_at DESC").First(&restriction).Error
	if errors.As(err, &gorm.ErrRecordNotFound) {
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

func NewRepository(logger logger.Logger, db *gorm.DB) *Repository {
	return &Repository{
		logger: logger,
		db:     db,
	}
}
