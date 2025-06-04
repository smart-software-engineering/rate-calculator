package model

type WorkSchedule struct {
	HoursPerWeek        float64
	VacationDays        int
	PublicHolidays      int
	EducationDays       int
	SickDays            int
	TotalWorkingDays    int
	TotalWorkingHours   float64
	EffectiveHourlyRate float64
}

type WorkScheduleTemplate struct {
	ID           string
	Label        string
	Description  string
	UserID       string
	Public       bool
	SortOrder    int
	WorkSchedule *WorkSchedule
}

func NewWorkSchedule() *WorkSchedule {
	return &WorkSchedule{
		HoursPerWeek:   40.0,
		VacationDays:   20,
		PublicHolidays: 10,
		EducationDays:  5,
		SickDays:       5,
	}
}

func (ws *WorkSchedule) CalculateWorkingTime() {
	daysPerWeek := 5
	
	weeksPerYear := 52
	totalPossibleWorkDays := daysPerWeek * weeksPerYear
	
	nonWorkingDays := ws.VacationDays + ws.PublicHolidays + ws.EducationDays + ws.SickDays
	
	ws.TotalWorkingDays = totalPossibleWorkDays - nonWorkingDays
	if ws.TotalWorkingDays < 0 {
		ws.TotalWorkingDays = 0
	}
	
	hoursPerDay := ws.HoursPerWeek / 5
	ws.TotalWorkingHours = float64(ws.TotalWorkingDays) * hoursPerDay
}

func (ws *WorkSchedule) Clone() *WorkSchedule {
	return &WorkSchedule{
		HoursPerWeek:        ws.HoursPerWeek,
		VacationDays:        ws.VacationDays,
		PublicHolidays:      ws.PublicHolidays,
		EducationDays:       ws.EducationDays,
		SickDays:            ws.SickDays,
		TotalWorkingDays:    ws.TotalWorkingDays,
		TotalWorkingHours:   ws.TotalWorkingHours,
		EffectiveHourlyRate: ws.EffectiveHourlyRate,
	}
}