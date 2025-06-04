package storage

import (
	"github.com/smart-software-engineering/rate-calculator/internal/model"
)

type ScheduleStorage interface {
	GetSchedules() ([]*model.WorkScheduleTemplate, error)
	
	GetScheduleByID(id string) (*model.WorkScheduleTemplate, error)
}