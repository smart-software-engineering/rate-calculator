package service

import (
	"net/http"
	"time"

	"github.com/smart-software-engineering/rate-calculator/internal/calculator"
	"github.com/smart-software-engineering/rate-calculator/internal/model"
	"github.com/smart-software-engineering/rate-calculator/internal/session"
	"github.com/smart-software-engineering/rate-calculator/internal/storage"
)

type UserDataService struct {
	sessionStore    session.Store
	scheduleStorage storage.ScheduleStorage
}

func NewUserDataService(sessionStore session.Store, scheduleStorage storage.ScheduleStorage) *UserDataService {
	return &UserDataService{
		sessionStore:    sessionStore,
		scheduleStorage: scheduleStorage,
	}
}

func (s *UserDataService) GetOrCreateUserData(w http.ResponseWriter, r *http.Request) (*model.UserData, error) {
	userData, err := s.sessionStore.GetOrCreateUserData(w, r)
	if err != nil {
		return nil, err
	}

	// Initialize with default schedule if none set
	if userData.ActiveScheduleID == "" {
		err = s.initializeDefaultSchedule(w, r, userData)
		if err != nil {
			return nil, err
		}
	}

	// Ensure data consistency on every load
	userData.EnsureExpenseConsistency()

	// Save any consistency fixes
	err = s.sessionStore.SetUserData(w, r, userData)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (s *UserDataService) UpdateScheduleFromTemplate(w http.ResponseWriter, r *http.Request, templateID string) (*model.UserData, error) {
	userData, err := s.GetOrCreateUserData(w, r)
	if err != nil {
		return nil, err
	}

	scheduleTemplate, err := s.scheduleStorage.GetScheduleByID(templateID)
	if err != nil {
		return nil, err
	}

	scheduleTemplate.WorkSchedule.CalculateWorkingTime()
	userData.SetSchedule(scheduleTemplate.WorkSchedule, templateID, scheduleTemplate.Label)

	return s.recalculateAndSave(w, r, userData)
}

func (s *UserDataService) UpdateScheduleValues(w http.ResponseWriter, r *http.Request, schedule *model.WorkSchedule) (*model.UserData, error) {
	userData, err := s.GetOrCreateUserData(w, r)
	if err != nil {
		return nil, err
	}

	schedule.CalculateWorkingTime()
	userData.Schedule = schedule
	userData.UpdateTimestamp()

	return s.recalculateAndSave(w, r, userData)
}

func (s *UserDataService) UpdateExpenseValues(w http.ResponseWriter, r *http.Request, expenses *model.ExpenseModel) (*model.UserData, error) {
	userData, err := s.GetOrCreateUserData(w, r)
	if err != nil {
		return nil, err
	}

	userData.SetExpenses(expenses)

	// Ensure data consistency after updating expenses
	userData.EnsureExpenseConsistency()

	return s.recalculateAndSave(w, r, userData)
}

func (s *UserDataService) recalculateAndSave(w http.ResponseWriter, r *http.Request, userData *model.UserData) (*model.UserData, error) {
	calcResults := calculator.CalculateRates(userData.ExpenseModel, userData.Schedule)

	categoryTotals := s.calculateCategoryTotals(userData.ExpenseModel, calcResults.TotalYearlyExpenses)

	result := &model.CalculationResult{
		TotalYearlyExpenses:    calcResults.TotalYearlyExpenses,
		YearlyTotalWithPercent: calcResults.YearlyTotalWithPercent,
		HourlyRate:             calcResults.HourlyRate,
		CategoryTotals:         categoryTotals,
		CalculatedAt:           time.Now(),
	}

	userData.SetCalculationResult(result)

	err := s.sessionStore.SetUserData(w, r, userData)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (s *UserDataService) calculateCategoryTotals(expenses *model.ExpenseModel, totalYearlyExpenses int) map[string]model.CategoryTotal {
	categoryTotals := make(map[string]model.CategoryTotal)

	baseExpenses := 0
	for _, category := range expenses.Categories {
		for _, item := range category.Items {
			switch item.Type {
			case model.Monthly:
				baseExpenses += item.Amount * 12
			case model.Yearly:
				baseExpenses += item.Amount
			}
		}
	}

	for _, category := range expenses.Categories {
		var percentageTotal int
		var yearlyTotal int

		for _, item := range category.Items {
			switch item.Type {
			case model.Monthly:
				yearlyTotal += item.Amount * 12
			case model.Yearly:
				yearlyTotal += item.Amount
			case model.Percentage:
				percentageTotal += item.Amount
				if baseExpenses > 0 {
					yearlyTotal += (baseExpenses * item.Amount) / 100
				}
			}
		}

		categoryTotals[category.Label] = model.CategoryTotal{
			PercentageTotal: percentageTotal,
			YearlyTotal:     yearlyTotal,
		}
	}

	return categoryTotals
}

func (s *UserDataService) initializeDefaultSchedule(w http.ResponseWriter, r *http.Request, userData *model.UserData) error {
	allSchedules, err := s.scheduleStorage.GetSchedules()
	if err != nil {
		return err
	}

	for _, schedule := range allSchedules {
		if schedule.Public {
			schedule.WorkSchedule.CalculateWorkingTime()
			userData.SetSchedule(schedule.WorkSchedule, schedule.ID, schedule.Label)

			return s.sessionStore.SetUserData(w, r, userData)
		}
	}

	return nil
}
