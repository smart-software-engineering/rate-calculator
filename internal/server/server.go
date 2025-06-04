package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/smart-software-engineering/rate-calculator/internal/calculator"
	"github.com/smart-software-engineering/rate-calculator/internal/model"
	"github.com/smart-software-engineering/rate-calculator/internal/session"
	"github.com/smart-software-engineering/rate-calculator/internal/storage"
	"github.com/smart-software-engineering/rate-calculator/internal/template"
)

//go:embed static
var staticFiles embed.FS

type ServerOptions struct {
	DevMode bool
}

type Server struct {
	Addr            string
	Template        template.Manager
	ScheduleStorage storage.ScheduleStorage
	SessionStore    session.Store
	CSRFKey         []byte
	Options         *ServerOptions
}

func New(addr string, tm template.Manager, ss storage.ScheduleStorage, sessionStore session.Store, options *ServerOptions) *Server {
	// Use the same key as the session for simplicity
	csrfKey := []byte(sessionStore.GetAuthKey())

	// Default options if none provided
	if options == nil {
		options = &ServerOptions{
			DevMode: false,
		}
	}

	return &Server{
		Addr:            addr,
		Template:        tm,
		ScheduleStorage: ss,
		SessionStore:    sessionStore,
		CSRFKey:         csrfKey,
		Options:         options,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Setup CSRF protection with environment-appropriate settings
	csrfOptions := []csrf.Option{
		csrf.Path("/"),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.FieldName("gorilla.csrf.Token"),
	}

	// Configure security options based on environment
	if s.Options.DevMode {
		// Development mode - less secure but more convenient
		csrfOptions = append(csrfOptions,
			csrf.Secure(false),
			csrf.SameSite(csrf.SameSiteLaxMode),
		)
		log.Println("Warning: Running in development mode with reduced security settings")
	} else {
		// Production mode - more secure
		csrfOptions = append(csrfOptions,
			csrf.Secure(true),
			csrf.SameSite(csrf.SameSiteStrictMode),
		)
	}

	// CSRF := csrf.Protect(s.CSRFKey, csrfOptions...)

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return err
	}

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		var expenses *model.ExpenseModel

		sessionExpenses, err := s.SessionStore.GetExpenses(r)
		if err != nil || sessionExpenses == nil {
			expenses = model.CreateSampleExpenseModel()
			s.SessionStore.SetExpenses(w, r, expenses)
		} else {
			expenses = sessionExpenses
		}

		allSchedules, err := s.ScheduleStorage.GetSchedules()
		if err != nil {
			http.Error(w, "Failed to load schedule templates: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var schedules []*model.WorkScheduleTemplate
		for _, schedule := range allSchedules {
			if schedule.Public {
				schedules = append(schedules, schedule)
			}
		}

		var activeSchedule *model.WorkScheduleTemplate
		var schedule *model.WorkSchedule

		sessionScheduleID, err := s.SessionStore.GetScheduleID(r)
		if err == nil && sessionScheduleID != "" {
			for _, s := range schedules {
				if s.ID == sessionScheduleID {
					activeSchedule = s
					break
				}
			}
		}

		if activeSchedule == nil && len(schedules) > 0 {
			activeSchedule = schedules[0]
		}

		sessionSchedule, err := s.SessionStore.GetSchedule(r)
		if err != nil || sessionSchedule == nil {
			schedule = activeSchedule.WorkSchedule
			s.SessionStore.SetSchedule(w, r, schedule)
		} else {
			schedule = sessionSchedule
		}

		if sessionScheduleID == "" || sessionScheduleID != activeSchedule.ID {
			s.SessionStore.SetScheduleID(w, r, activeSchedule.ID)
		}

		// Calculate initial rates
		results := calculator.CalculateRates(expenses, schedule)

		data := template.TemplateData{
			"Expenses":               expenses,
			"Schedule":               schedule,
			"ScheduleID":             activeSchedule.ID,
			"ScheduleLabel":          activeSchedule.Label,
			"Schedules":              schedules,
			"TotalYearlyExpenses":    results.TotalYearlyExpenses,
			"YearlyTotalWithPercent": results.YearlyTotalWithPercent,
			"HourlyRate":             results.HourlyRate,
			"CSRFToken":              csrf.Token(r),
			"CSRFField":              csrf.TemplateField(r),
		}

		if err := s.Template.Render(w, "index.html", data); err != nil {
			http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("POST /schedule", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		id := r.FormValue("scheduleID")
		if id == "" {
			http.Error(w, "Schedule ID is required", http.StatusBadRequest)
			return
		}

		log.Printf("Schedule selection: ID=%s, CSRF Token present: %v", id, r.Header.Get("X-CSRF-Token") != "")

		schedule, err := s.ScheduleStorage.GetScheduleByID(id)
		if err != nil {
			http.Error(w, "Schedule not found: "+err.Error(), http.StatusNotFound)
			return
		}

		if !schedule.Public && schedule.UserID != "system" {
			http.Error(w, "Schedule not available", http.StatusForbidden)
			return
		}

		schedule.WorkSchedule.CalculateWorkingTime()

		s.SessionStore.SetSchedule(w, r, schedule.WorkSchedule)
		s.SessionStore.SetScheduleID(w, r, schedule.ID)

		// Get current expenses for rate calculation
		expenses, _ := s.SessionStore.GetExpenses(r)
		if expenses == nil {
			expenses = model.CreateSampleExpenseModel()
		}

		// Calculate rates with the selected schedule
		results := calculator.CalculateRates(expenses, schedule.WorkSchedule)

		data := template.TemplateData{
			"Schedule":               schedule.WorkSchedule,
			"ScheduleLabel":          schedule.Label,
			"Expenses":               expenses,
			"TotalYearlyExpenses":    results.TotalYearlyExpenses,
			"YearlyTotalWithPercent": results.YearlyTotalWithPercent,
			"HourlyRate":             results.HourlyRate,
			"CSRFToken":              csrf.Token(r),
			"CSRFField":              csrf.TemplateField(r),
		}

		if err := s.Template.Render(w, "partials/schedule_form", data); err != nil {
			http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("POST /calculate-working-hours", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		schedule, err := s.SessionStore.GetSchedule(r)
		if err != nil || schedule == nil {
			schedule = model.NewWorkSchedule()
		}

		for key, values := range r.Form {
			if len(values) == 0 {
				continue
			}

			if strings.HasPrefix(key, "schedule[") && strings.HasSuffix(key, "]") {
				paramName := key[len("schedule[") : len(key)-1]
				value := values[0]

				switch paramName {
				case "hoursPerWeek":
					if val, err := strconv.ParseFloat(value, 64); err == nil {
						schedule.HoursPerWeek = val
					}
				case "vacationDays":
					if val, err := strconv.Atoi(value); err == nil {
						schedule.VacationDays = val
					}
				case "publicHolidays":
					if val, err := strconv.Atoi(value); err == nil {
						schedule.PublicHolidays = val
					}
				case "educationDays":
					if val, err := strconv.Atoi(value); err == nil {
						schedule.EducationDays = val
					}
				case "sickDays":
					if val, err := strconv.Atoi(value); err == nil {
						schedule.SickDays = val
					}
				}
			}
		}

		schedule.CalculateWorkingTime()
		s.SessionStore.SetSchedule(w, r, schedule)

		currentScheduleID, _ := s.SessionStore.GetScheduleID(r)

		// Get current expenses for rate calculation
		expenses, _ := s.SessionStore.GetExpenses(r)
		if expenses == nil {
			expenses = model.CreateSampleExpenseModel()
		}

		// Calculate rates with the updated schedule
		results := calculator.CalculateRates(expenses, schedule)

		data := template.TemplateData{
			"Schedule":               schedule,
			"ScheduleLabel":          "Custom",
			"ScheduleID":             currentScheduleID,
			"Expenses":               expenses,
			"TotalYearlyExpenses":    results.TotalYearlyExpenses,
			"YearlyTotalWithPercent": results.YearlyTotalWithPercent,
			"HourlyRate":             results.HourlyRate,
			"CSRFToken":              csrf.Token(r),
			"CSRFField":              csrf.TemplateField(r),
		}

		if err := s.Template.Render(w, "partials/schedule_summary", data); err != nil {
			http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("POST /calculate", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		expenses, err := s.SessionStore.GetExpenses(r)
		if err != nil || expenses == nil {
			expenses = model.CreateSampleExpenseModel()
		}

		updatedExpenses := updateExpensesFromForm(expenses, r.Form)
		s.SessionStore.SetExpenses(w, r, updatedExpenses)

		// Get the current work schedule
		schedule, err := s.SessionStore.GetSchedule(r)
		if err != nil || schedule == nil {
			// Get the default schedule if none exists
			scheduleID, _ := s.SessionStore.GetScheduleID(r)
			defaultSchedule, err := s.ScheduleStorage.GetScheduleByID(scheduleID)
			if err != nil || defaultSchedule == nil {
				// Use a basic default if nothing else is available
				schedule = model.NewWorkSchedule()
			} else {
				schedule = defaultSchedule.WorkSchedule
			}
		}

		// Calculate the results
		results := calculator.CalculateRates(updatedExpenses, schedule)

		// Prepare data for the template
		data := template.TemplateData{
			"Expenses":               updatedExpenses,
			"Schedule":               schedule,
			"TotalYearlyExpenses":    results.TotalYearlyExpenses,
			"YearlyTotalWithPercent": results.YearlyTotalWithPercent,
			"HourlyRate":             results.HourlyRate,
			"CSRFToken":              csrf.Token(r),
			"CSRFField":              csrf.TemplateField(r),
		}

		// Render just the calculation result template
		if err := s.Template.Render(w, "partials/calculation_result", data); err != nil {
			http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /clear-session", func(w http.ResponseWriter, r *http.Request) {
		s.SessionStore.Clear(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	fileServer := http.FileServer(http.FS(staticFS))
	mux.Handle("GET /app/", http.StripPrefix("/app/", fileServer))

	// Create a middleware chain
	var handler http.Handler = mux

	// Add CSRF protection
	// handler = CSRF(handler)

	// Create server with the middleware chain
	server := &http.Server{
		Addr:    s.Addr,
		Handler: handler,
	}

	log.Printf("Visit http://localhost%s to view the rate calculator\n", s.Addr)

	return server.ListenAndServe()
}

func updateExpensesFromForm(expenses *model.ExpenseModel, form map[string][]string) *model.ExpenseModel {
	for key, values := range form {
		if !strings.HasPrefix(key, "expense[") || !strings.HasSuffix(key, "]") || len(values) == 0 {
			continue
		}

		parts := strings.Split(key, "][")
		if len(parts) != 2 {
			continue
		}

		expenseName := parts[0][len("expense["):]
		paramType := parts[1][:len(parts[1])-1]
		value := values[0]

		for categoryIdx, category := range expenses.Categories {
			for itemIdx, item := range category.Items {
				if item.Label == expenseName {
					if paramType == "amount" {
						if amount, err := strconv.Atoi(value); err == nil {
							expenses.Categories[categoryIdx].Items[itemIdx].Amount = amount
						}
					} else if paramType == "type" {
						switch value {
						case "monthly":
							expenses.Categories[categoryIdx].Items[itemIdx].Type = model.Monthly
						case "yearly":
							expenses.Categories[categoryIdx].Items[itemIdx].Type = model.Yearly
						case "percentage":
							expenses.Categories[categoryIdx].Items[itemIdx].Type = model.Percentage
						}
					}
				}
			}
		}
	}

	return expenses
}
