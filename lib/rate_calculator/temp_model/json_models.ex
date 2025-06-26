defmodule RateCalculator.TempModel.JsonModels do

"""
  type Schedule struct {
	Id           uuid.UUID    `json:"id"`
	Label        string       `json:"label"`
	Description  string       `json:"description"`
	UserId       uuid.UUID    `json:"userId"`
	IsPublic     bool         `json:"isPublic"`
	WorkSchedule WorkSchedule `json:"workSchedule"`
}

type WorkSchedule struct {
	HoursPerWeek    float64 `json:"hoursPerWeek"`
	PrivateHolidays int     `json:"privateHolidays"`
	PublicHolidays  int     `json:"publicHolidays"`
	EducationDays   int     `json:"educationDays"`
	SickDays        int     `json:"sickDays"`
}
"""

end
