package mysql

import (
	"github.com/quanxiang-cloud/faas/internal/models"
	"gorm.io/gorm"
)

type eventRepo struct {
}

func NewEventRepo() models.EventRepo {
	return &eventRepo{}
}

func (e *eventRepo) Insert(db *gorm.DB, event *models.Event) error {
	return db.Create(event).Error
}

func (e *eventRepo) Update(db *gorm.DB, event *models.Event) error {
	return db.Updates(event).Error
}

func (e *eventRepo) Query(db *gorm.DB, id string) (*models.Event, error) {
	ret := &models.Event{}
	err := db.Where("id = ?", id).Find(ret).Error
	return ret, err
}

func (e *eventRepo) QueryByName(db *gorm.DB, name string) (*models.Event, error) {
	ret := &models.Event{}
	err := db.Where("name = ?", name).Find(ret).Error
	return ret, err
}

func (e *eventRepo) Delete(db *gorm.DB, id string) error {
	return db.Where("id = ?", id).Delete(&models.Event{}).Error
}
