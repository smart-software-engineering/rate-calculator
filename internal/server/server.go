package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/smart-software-engineering/rate-calculator/internal/model"
	"github.com/smart-software-engineering/rate-calculator/internal/service"
	"github.com/smart-software-engineering/rate-calculator/internal/session"
	"github.com/smart-software-engineering/rate-calculator/internal/storage"
	tmpl "github.com/smart-software-engineering/rate-calculator/internal/template"
)

//go:embed static
var staticFiles embed.FS

type ServerOptions struct {
	DevMode bool
}

type Server struct {
	Addr            string
	Template        tmpl.Manager
	ScheduleStorage storage.ScheduleStorage
	SessionStore    session.Store
	UserDataService *service.UserDataService
	Options         *ServerOptions
}

func New(addr string, tm tmpl.Manager, ss storage.ScheduleStorage, sessionStore session.Store, options *ServerOptions) *Server {
	if options == nil {
		options = &ServerOptions{
			DevMode: false,
		}
	}

	userDataService := service.NewUserDataService(sessionStore, ss)

	return &Server{
		Addr:            addr,
		Template:        tm,
		ScheduleStorage: ss,
		SessionStore:    sessionStore,
		UserDataService: userDataService,
		Options:         options,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return err
	}

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		userData, err := s.UserDataService.GetOrCreateUserData(w, r)
		if err != nil {
			http.Error(w, "Failed to get user data: "+err.Error(), http.StatusInternalServerError)
			return
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

		expenseCategories := s.prepareExpenseCategoriesFromUserData(userData)

		var calcResults *model.CalculationResult
		if userData.LastCalculation != nil {
			calcResults = userData.LastCalculation
		} else {
			userData, err = s.UserDataService.UpdateExpenseValues(w, r, userData.ExpenseModel)
			if err != nil {
				http.Error(w, "Failed to calculate rates: "+err.Error(), http.StatusInternalServerError)
				return
			}
			calcResults = userData.LastCalculation
		}

		data := tmpl.TemplateData{
			"Expenses":               userData.ExpenseModel,
			"ExpenseCategories":      expenseCategories,
			"Schedule":               userData.Schedule,
			"ScheduleID":             userData.ActiveScheduleID,
			"ScheduleLabel":          userData.ScheduleTemplate,
			"Schedules":              schedules,
			"TotalYearlyExpenses":    calcResults.TotalYearlyExpenses,
			"YearlyTotalWithPercent": calcResults.YearlyTotalWithPercent,
			"HourlyRate":             calcResults.HourlyRate,
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

		userData, err := s.UserDataService.UpdateScheduleFromTemplate(w, r, id)
		if err != nil {
			http.Error(w, "Failed to update schedule: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := tmpl.TemplateData{
			"Schedule":               userData.Schedule,
			"ScheduleLabel":          userData.ScheduleTemplate,
			"Expenses":               userData.ExpenseModel,
			"TotalYearlyExpenses":    userData.LastCalculation.TotalYearlyExpenses,
			"YearlyTotalWithPercent": userData.LastCalculation.YearlyTotalWithPercent,
			"HourlyRate":             userData.LastCalculation.HourlyRate,
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

		userData, err := s.UserDataService.GetOrCreateUserData(w, r)
		if err != nil {
			http.Error(w, "Failed to get user data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		schedule := userData.Schedule
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

		userData, err = s.UserDataService.UpdateScheduleValues(w, r, schedule)
		if err != nil {
			http.Error(w, "Failed to update schedule: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := tmpl.TemplateData{
			"Schedule":               userData.Schedule,
			"ScheduleLabel":          userData.ScheduleTemplate,
			"ScheduleID":             userData.ActiveScheduleID,
			"Expenses":               userData.ExpenseModel,
			"TotalYearlyExpenses":    userData.LastCalculation.TotalYearlyExpenses,
			"YearlyTotalWithPercent": userData.LastCalculation.YearlyTotalWithPercent,
			"HourlyRate":             userData.LastCalculation.HourlyRate,
		}

		if err := s.Template.Render(w, "partials/schedule_form", data); err != nil {
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

		userData, err := s.UserDataService.GetOrCreateUserData(w, r)
		if err != nil {
			http.Error(w, "Failed to get user data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		updatedExpenses := updateExpensesFromForm(userData.ExpenseModel, r.Form)

		userData, err = s.UserDataService.UpdateExpenseValues(w, r, updatedExpenses)
		if err != nil {
			http.Error(w, "Failed to calculate rates: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := tmpl.TemplateData{
			"Expenses":               userData.ExpenseModel,
			"Schedule":               userData.Schedule,
			"TotalYearlyExpenses":    userData.LastCalculation.TotalYearlyExpenses,
			"YearlyTotalWithPercent": userData.LastCalculation.YearlyTotalWithPercent,
			"HourlyRate":             userData.LastCalculation.HourlyRate,
		}

		if err := s.Template.Render(w, "partials/calculation_result", data); err != nil {
			http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("POST /update-category", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		userData, err := s.UserDataService.GetOrCreateUserData(w, r)
		if err != nil {
			http.Error(w, "Failed to get user data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		updatedExpenses := updateExpensesFromForm(userData.ExpenseModel, r.Form)

		userData, err = s.UserDataService.UpdateExpenseValues(w, r, updatedExpenses)
		if err != nil {
			http.Error(w, "Failed to update expenses: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		data := tmpl.TemplateData{
			"Expenses":               userData.ExpenseModel,
			"Schedule":               userData.Schedule,
			"TotalYearlyExpenses":    userData.LastCalculation.TotalYearlyExpenses,
			"YearlyTotalWithPercent": userData.LastCalculation.YearlyTotalWithPercent,
			"HourlyRate":             userData.LastCalculation.HourlyRate,
		}

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

	var handler http.Handler = mux

	server := &http.Server{
		Addr:    s.Addr,
		Handler: handler,
	}

	log.Printf("Visit http://localhost%s to view the rate calculator\n", s.Addr)

	return server.ListenAndServe()
}

func updateExpensesFromForm(expenses *model.ExpenseModel, form map[string][]string) *model.ExpenseModel {
	updatedExpenses := &model.ExpenseModel{
		Categories: make([]model.ExpenseCategory, len(expenses.Categories)),
	}

	for i, category := range expenses.Categories {
		updatedExpenses.Categories[i] = model.ExpenseCategory{
			ID:    category.ID,
			Label: category.Label,
			Items: make([]model.ExpenseItem, len(category.Items)),
		}
		copy(updatedExpenses.Categories[i].Items, category.Items)
	}

	// Track items that need to be created
	itemsToCreate := make(map[string]map[string]interface{}) // itemID -> {amount, type, label}

	for key, values := range form {
		if !strings.HasPrefix(key, "expense[") || !strings.HasSuffix(key, "]") || len(values) == 0 {
			continue
		}

		parts := strings.Split(key, "][")
		if len(parts) != 2 {
			continue
		}

		expenseID := parts[0][len("expense["):]
		paramType := parts[1][:len(parts[1])-1]
		value := values[0]

		found := false
		for categoryIdx, category := range updatedExpenses.Categories {
			for itemIdx, item := range category.Items {
				if item.ID == expenseID {
					if paramType == "amount" {
						if amount, err := strconv.Atoi(value); err == nil && amount >= 0 {
							updatedExpenses.Categories[categoryIdx].Items[itemIdx].Amount = amount
							found = true
						} else {
							log.Printf("Invalid amount value for item ID %s: %s (error: %v)", expenseID, value, err)
						}
					} else if paramType == "type" {
						switch value {
						case "monthly":
							updatedExpenses.Categories[categoryIdx].Items[itemIdx].Type = model.Monthly
							found = true
						case "yearly":
							updatedExpenses.Categories[categoryIdx].Items[itemIdx].Type = model.Yearly
							found = true
						case "percentage":
							updatedExpenses.Categories[categoryIdx].Items[itemIdx].Type = model.Percentage
							found = true
						}
					} else if paramType == "label" {
						updatedExpenses.Categories[categoryIdx].Items[itemIdx].Label = value
						found = true
					}
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			// Item doesn't exist, prepare to create it
			if itemsToCreate[expenseID] == nil {
				itemsToCreate[expenseID] = make(map[string]interface{})
			}

			if paramType == "amount" {
				if amount, err := strconv.Atoi(value); err == nil && amount >= 0 {
					itemsToCreate[expenseID]["amount"] = amount
				}
			} else if paramType == "type" {
				switch value {
				case "monthly":
					itemsToCreate[expenseID]["type"] = model.Monthly
				case "yearly":
					itemsToCreate[expenseID]["type"] = model.Yearly
				case "percentage":
					itemsToCreate[expenseID]["type"] = model.Percentage
				}
			} else if paramType == "label" {
				itemsToCreate[expenseID]["label"] = value
			}

			log.Printf("Expense item not found, will attempt to create with ID: %s", expenseID)
		}
	}

	// Create missing items in appropriate categories
	for itemID, params := range itemsToCreate {
		// Only create if we have both amount and type
		amount, hasAmount := params["amount"].(int)
		expenseType, hasType := params["type"].(model.ExpenseType)
		label, hasLabel := params["label"].(string)

		// Use itemID as label if no label provided
		if !hasLabel || label == "" {
			label = "New Item " + itemID[:8] // Use first 8 chars of ID as fallback
		}

		if hasAmount && hasType {
			// Determine which category to add it to based on item label patterns
			categoryIndex := determineCategoryForItem(label)

			if categoryIndex < len(updatedExpenses.Categories) {
				newItem := model.ExpenseItem{
					ID:     itemID,
					Label:  label,
					Amount: amount,
					Type:   expenseType,
				}

				updatedExpenses.Categories[categoryIndex].Items = append(
					updatedExpenses.Categories[categoryIndex].Items,
					newItem,
				)

				log.Printf("Created new expense item: %s (ID: %s) in category: %s",
					label, itemID, updatedExpenses.Categories[categoryIndex].Label)
			} else {
				log.Printf("Could not determine category for new item: %s (ID: %s)", label, itemID)
			}
		}
	}

	return updatedExpenses
}

// determineCategoryForItem suggests which category a new expense item should belong to
func determineCategoryForItem(itemName string) int {
	itemLower := strings.ToLower(itemName)

	// Professional Expenses (index 1 in sample model)
	professionalKeywords := []string{
		"insurance", "software", "hardware", "office", "communication",
		"training", "education", "mobility", "ads", "advertising",
		"fees", "donation", "business", "professional", "equipment",
	}

	// Private Expenses (index 0 in sample model)
	privateKeywords := []string{
		"rent", "household", "utilities", "groceries", "clothing",
		"internet", "phone", "personal", "private", "home",
	}

	// Financial Goals (index 2 in sample model)
	financialKeywords := []string{
		"savings", "retirement", "emergency", "investment", "fund",
	}

	// Taxes (index 3 in sample model)
	taxKeywords := []string{
		"tax", "pension", "health", "social", "vat", "income",
	}

	for _, keyword := range professionalKeywords {
		if strings.Contains(itemLower, keyword) {
			return 1 // Professional Expenses
		}
	}

	for _, keyword := range financialKeywords {
		if strings.Contains(itemLower, keyword) {
			return 2 // Financial Goals
		}
	}

	for _, keyword := range taxKeywords {
		if strings.Contains(itemLower, keyword) {
			return 3 // Taxes
		}
	}

	for _, keyword := range privateKeywords {
		if strings.Contains(itemLower, keyword) {
			return 0 // Private Expenses
		}
	}

	// Default to Professional Expenses if unsure
	return 1
}

func (s *Server) prepareExpenseCategoriesFromUserData(userData *model.UserData) []model.ExpenseCategoryData {
	var categories []model.ExpenseCategoryData

	for _, category := range userData.ExpenseModel.Categories {
		var percentageTotal int
		var yearlyTotal int

		if userData.LastCalculation != nil && userData.LastCalculation.CategoryTotals != nil {
			if totals, exists := userData.LastCalculation.CategoryTotals[category.Label]; exists {
				percentageTotal = totals.PercentageTotal
				yearlyTotal = totals.YearlyTotal
			}
		}

		categories = append(categories, model.ExpenseCategoryData{
			Category:        category,
			PercentageTotal: percentageTotal,
			YearlyTotal:     yearlyTotal,
		})
	}

	return categories
}
