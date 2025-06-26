package rates

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"

	"github.com/google/uuid"
)

func NewRateCalculator(file fs.File) RateCalculator {
	var schedule Schedule

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&schedule); err != nil {
		fmt.Println(err)
		log.Fatal("Could not decode the schedule file")
	}

	return &rateCalc{schedule: schedule}
}

type rateCalc struct {
	schedule Schedule
}

func (r *rateCalc) Schedules() (Schedule, error) {
	return r.schedule, nil
}

type RateCalculator interface {
	Schedules() (Schedule, error)
}

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
