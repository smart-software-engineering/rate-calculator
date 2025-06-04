package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/smart-software-engineering/rate-calculator/internal/model"
)

type ScheduleStorage struct {
	schedules map[string]*model.WorkScheduleTemplate
	mu        sync.RWMutex
}

func NewScheduleStorage(dataDir string) (*ScheduleStorage, error) {
	storage := &ScheduleStorage{
		schedules: make(map[string]*model.WorkScheduleTemplate),
	}

	err := storage.loadSchedulesFromDir(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load schedules: %w", err)
	}

	return storage, nil
}

func (s *ScheduleStorage) GetSchedules() ([]*model.WorkScheduleTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	schedules := make([]*model.WorkScheduleTemplate, 0, len(s.schedules))
	for _, schedule := range s.schedules {
		schedules = append(schedules, schedule)
	}

	sort.Slice(schedules, func(i, j int) bool {
		if schedules[i].SortOrder == schedules[j].SortOrder {
			return schedules[i].Label < schedules[j].Label
		}
		return schedules[i].SortOrder < schedules[j].SortOrder
	})

	return schedules, nil
}

func (s *ScheduleStorage) GetScheduleByID(id string) (*model.WorkScheduleTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	schedule, ok := s.schedules[id]
	if !ok {
		return nil, fmt.Errorf("schedule with ID %s not found", id)
	}

	clone := &model.WorkScheduleTemplate{
		ID:           schedule.ID,
		Label:        schedule.Label,
		Description:  schedule.Description,
		UserID:       schedule.UserID,
		Public:       schedule.Public,
		SortOrder:    schedule.SortOrder,
		WorkSchedule: schedule.WorkSchedule.Clone(),
	}

	return clone, nil
}

func (s *ScheduleStorage) loadSchedulesFromDir(dir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		schedule, err := s.loadScheduleFromFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to load schedule from %s: %w", filePath, err)
		}

		if schedule.ID == "" {
			schedule.ID = filepath.Base(file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))])
		}

		schedule.WorkSchedule.CalculateWorkingTime()

		s.schedules[schedule.ID] = schedule
	}

	return nil
}

func (s *ScheduleStorage) loadScheduleFromFile(filePath string) (*model.WorkScheduleTemplate, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var schedule model.WorkScheduleTemplate
	err = json.Unmarshal(data, &schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &schedule, nil
}