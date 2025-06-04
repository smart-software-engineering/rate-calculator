package calculator

import (
	"math"

	"github.com/smart-software-engineering/rate-calculator/internal/model"
)

type CalculationResult struct {
	TotalYearlyExpenses    int
	YearlyTotalWithPercent int
	HourlyRate             float64
}

func CalculateRates(expenses *model.ExpenseModel, schedule *model.WorkSchedule) *CalculationResult {
	var yearlyTotal int
	var totalPercentage float64

	for _, category := range expenses.Categories {
		for _, item := range category.Items {
			switch item.Type {
			case model.Monthly:
				yearlyTotal += item.Amount * 12
			case model.Yearly:
				yearlyTotal += item.Amount
			case model.Percentage:
				totalPercentage += float64(item.Amount)
			}
		}
	}

	percentageFactor := totalPercentage / 100
	yearlyTotalWithPercent := int(math.Round(float64(yearlyTotal) * (1 + percentageFactor)))

	workingHours := schedule.TotalWorkingHours

	var hourlyRate float64
	if workingHours > 0 {
		hourlyRate = roundToTwoDecimals(float64(yearlyTotalWithPercent) / workingHours)
	}

	return &CalculationResult{
		TotalYearlyExpenses:    yearlyTotal,
		YearlyTotalWithPercent: yearlyTotalWithPercent,
		HourlyRate:             hourlyRate,
	}
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
