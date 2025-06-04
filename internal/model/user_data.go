package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// UserData represents the complete state of a user's session
type UserData struct {
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Schedule data
	ActiveScheduleID string        `json:"active_schedule_id"`
	ScheduleTemplate string        `json:"schedule_template"` // source template name
	Schedule         *WorkSchedule `json:"schedule"`

	// Expense data
	ExpenseModel    *ExpenseModel     `json:"expense_model"`
	CategoryOrder   []string          `json:"category_order"`   // tracks custom ordering by category ID
	CategorySources map[string]string `json:"category_sources"` // tracks which template each category came from by category ID

	// Calculated values (cached for performance)
	LastCalculation *CalculationResult `json:"last_calculation"`
}

// CalculationResult stores the calculated values
type CalculationResult struct {
	TotalYearlyExpenses    int                      `json:"total_yearly_expenses"`
	YearlyTotalWithPercent int                      `json:"yearly_total_with_percent"`
	HourlyRate             float64                  `json:"hourly_rate"`
	CategoryTotals         map[string]CategoryTotal `json:"category_totals"`
	CalculatedAt           time.Time                `json:"calculated_at"`
}

// CategoryTotal stores calculated totals for a category
type CategoryTotal struct {
	PercentageTotal int `json:"percentage_total"`
	YearlyTotal     int `json:"yearly_total"`
}

// NewUserData creates a new UserData instance with a generated UUID
func NewUserData() *UserData {
	now := time.Now()
	return &UserData{
		UserID:          uuid.New().String(),
		CreatedAt:       now,
		UpdatedAt:       now,
		Schedule:        NewWorkSchedule(),
		ExpenseModel:    NewExpenseModel(),
		CategoryOrder:   []string{},
		CategorySources: make(map[string]string),
	}
}

// NewUserDataWithDefaults creates a new UserData with sample data
func NewUserDataWithDefaults() *UserData {
	userData := NewUserData()
	userData.ExpenseModel = CreateSampleExpenseModel()

	// Initialize category order and sources
	for _, category := range userData.ExpenseModel.Categories {
		userData.CategoryOrder = append(userData.CategoryOrder, category.ID)
		userData.CategorySources[category.ID] = "default"
	}

	return userData
}

// UpdateTimestamp updates the UpdatedAt field
func (ud *UserData) UpdateTimestamp() {
	ud.UpdatedAt = time.Now()
}

// SetSchedule updates the schedule and metadata
func (ud *UserData) SetSchedule(schedule *WorkSchedule, templateID string, templateName string) {
	ud.Schedule = schedule
	ud.ActiveScheduleID = templateID
	ud.ScheduleTemplate = templateName
	ud.UpdateTimestamp()
}

// SetExpenses updates the expense model
func (ud *UserData) SetExpenses(expenses *ExpenseModel) {
	ud.ExpenseModel = expenses
	ud.UpdateTimestamp()
}

// SetCalculationResult caches the calculation result
func (ud *UserData) SetCalculationResult(result *CalculationResult) {
	ud.LastCalculation = result
	ud.UpdateTimestamp()
}

// AddCategory adds a new category and updates metadata
func (ud *UserData) AddCategory(category ExpenseCategory, source string) {
	ud.ExpenseModel.AddCategory(category)
	ud.CategoryOrder = append(ud.CategoryOrder, category.ID)
	ud.CategorySources[category.ID] = source
	ud.UpdateTimestamp()
}

// AddItemToCategory adds an expense item to an existing category
func (ud *UserData) AddItemToCategory(categoryID string, item ExpenseItem) bool {
	for i, category := range ud.ExpenseModel.Categories {
		if category.ID == categoryID {
			ud.ExpenseModel.Categories[i].Items = append(ud.ExpenseModel.Categories[i].Items, item)
			ud.UpdateTimestamp()
			return true
		}
	}
	return false
}

// FindExpenseItem finds an expense item by ID across all categories
func (ud *UserData) FindExpenseItem(itemID string) (*ExpenseItem, int, int) {
	for categoryIdx, category := range ud.ExpenseModel.Categories {
		for itemIdx, item := range category.Items {
			if item.ID == itemID {
				return &item, categoryIdx, itemIdx
			}
		}
	}
	return nil, -1, -1
}

// FindExpenseItemByLabel finds an expense item by label across all categories (for backward compatibility)
func (ud *UserData) FindExpenseItemByLabel(itemLabel string) (*ExpenseItem, int, int) {
	for categoryIdx, category := range ud.ExpenseModel.Categories {
		for itemIdx, item := range category.Items {
			if item.Label == itemLabel {
				return &item, categoryIdx, itemIdx
			}
		}
	}
	return nil, -1, -1
}

// FindCategoryByID finds a category by its ID
func (ud *UserData) FindCategoryByID(categoryID string) (*ExpenseCategory, int) {
	for idx, category := range ud.ExpenseModel.Categories {
		if category.ID == categoryID {
			return &category, idx
		}
	}
	return nil, -1
}

// FindCategoryByLabel finds a category by its label (for backward compatibility)
func (ud *UserData) FindCategoryByLabel(categoryLabel string) (*ExpenseCategory, int) {
	for idx, category := range ud.ExpenseModel.Categories {
		if category.Label == categoryLabel {
			return &category, idx
		}
	}
	return nil, -1
}

// UpdateExpenseItem updates an existing expense item by ID
func (ud *UserData) UpdateExpenseItem(itemID string, amount int, expenseType ExpenseType) bool {
	// Try to find existing item by ID
	_, categoryIdx, itemIdx := ud.FindExpenseItem(itemID)

	if categoryIdx >= 0 && itemIdx >= 0 {
		// Update existing item
		ud.ExpenseModel.Categories[categoryIdx].Items[itemIdx].Amount = amount
		ud.ExpenseModel.Categories[categoryIdx].Items[itemIdx].Type = expenseType
		ud.UpdateTimestamp()
		return true
	}

	return false // Item not found
}

// UpdateExpenseItemByLabel updates an existing expense item by label (for backward compatibility)
func (ud *UserData) UpdateExpenseItemByLabel(itemLabel string, amount int, expenseType ExpenseType) bool {
	// Try to find existing item by label
	_, categoryIdx, itemIdx := ud.FindExpenseItemByLabel(itemLabel)

	if categoryIdx >= 0 && itemIdx >= 0 {
		// Update existing item
		ud.ExpenseModel.Categories[categoryIdx].Items[itemIdx].Amount = amount
		ud.ExpenseModel.Categories[categoryIdx].Items[itemIdx].Type = expenseType
		ud.UpdateTimestamp()
		return true
	}

	return false // Item not found
}

// GetOrCreateExpenseItem gets an existing item by ID or creates a new one in the appropriate category
func (ud *UserData) GetOrCreateExpenseItem(itemID string, label string, amount int, expenseType ExpenseType) *ExpenseItem {
	// Try to find existing item by ID
	item, categoryIdx, itemIdx := ud.FindExpenseItem(itemID)

	if item != nil {
		// Update existing item
		ud.ExpenseModel.Categories[categoryIdx].Items[itemIdx].Amount = amount
		ud.ExpenseModel.Categories[categoryIdx].Items[itemIdx].Type = expenseType
		// Update label if provided
		if label != "" {
			ud.ExpenseModel.Categories[categoryIdx].Items[itemIdx].Label = label
		}
		ud.UpdateTimestamp()
		return &ud.ExpenseModel.Categories[categoryIdx].Items[itemIdx]
	}

	// Create new item in appropriate category
	newItem := NewExpenseItem(label, amount, expenseType)
	// Use provided ID if given, otherwise NewExpenseItem already generated one
	if itemID != "" {
		newItem.ID = itemID
	}

	// Determine best category for this item
	categoryIndex := ud.determineBestCategoryForItem(label)

	if categoryIndex < len(ud.ExpenseModel.Categories) {
		ud.ExpenseModel.Categories[categoryIndex].Items = append(
			ud.ExpenseModel.Categories[categoryIndex].Items,
			newItem,
		)
		ud.UpdateTimestamp()
		return &ud.ExpenseModel.Categories[categoryIndex].Items[len(ud.ExpenseModel.Categories[categoryIndex].Items)-1]
	}

	return nil
}

// GetOrCreateExpenseItemByLabel gets an existing item by label or creates a new one (for backward compatibility)
func (ud *UserData) GetOrCreateExpenseItemByLabel(itemLabel string, amount int, expenseType ExpenseType) *ExpenseItem {
	// Try to find existing item by label
	item, categoryIdx, itemIdx := ud.FindExpenseItemByLabel(itemLabel)

	if item != nil {
		// Update existing item
		ud.ExpenseModel.Categories[categoryIdx].Items[itemIdx].Amount = amount
		ud.ExpenseModel.Categories[categoryIdx].Items[itemIdx].Type = expenseType
		ud.UpdateTimestamp()
		return &ud.ExpenseModel.Categories[categoryIdx].Items[itemIdx]
	}

	// Create new item in appropriate category
	newItem := NewExpenseItem(itemLabel, amount, expenseType)

	// Determine best category for this item
	categoryIndex := ud.determineBestCategoryForItem(itemLabel)

	if categoryIndex < len(ud.ExpenseModel.Categories) {
		ud.ExpenseModel.Categories[categoryIndex].Items = append(
			ud.ExpenseModel.Categories[categoryIndex].Items,
			newItem,
		)
		ud.UpdateTimestamp()
		return &ud.ExpenseModel.Categories[categoryIndex].Items[len(ud.ExpenseModel.Categories[categoryIndex].Items)-1]
	}

	return nil
}

// determineBestCategoryForItem uses intelligent categorization for new items
func (ud *UserData) determineBestCategoryForItem(itemLabel string) int {
	// Use the same logic as the server function but adapted for UserData
	itemLower := strings.ToLower(itemLabel)

	// Map keywords to category indices based on common category names
	categoryMappings := map[string][]string{
		"professional": {"insurance", "software", "hardware", "office", "communication", "training", "education", "mobility", "ads", "advertising", "fees", "donation", "business", "professional", "equipment"},
		"private":      {"rent", "household", "utilities", "groceries", "clothing", "internet", "phone", "personal", "private", "home"},
		"financial":    {"savings", "retirement", "emergency", "investment", "fund"},
		"tax":          {"tax", "pension", "health", "social", "vat", "income"},
	}

	// Find category indices by label
	categoryIndices := make(map[string]int)
	for i, category := range ud.ExpenseModel.Categories {
		lowerLabel := strings.ToLower(category.Label)
		if strings.Contains(lowerLabel, "professional") || strings.Contains(lowerLabel, "business") {
			categoryIndices["professional"] = i
		} else if strings.Contains(lowerLabel, "private") || strings.Contains(lowerLabel, "personal") {
			categoryIndices["private"] = i
		} else if strings.Contains(lowerLabel, "financial") || strings.Contains(lowerLabel, "goal") {
			categoryIndices["financial"] = i
		} else if strings.Contains(lowerLabel, "tax") {
			categoryIndices["tax"] = i
		}
	}

	// Try to match keywords to categories
	for categoryType, keywords := range categoryMappings {
		for _, keyword := range keywords {
			if strings.Contains(itemLower, keyword) {
				if index, exists := categoryIndices[categoryType]; exists {
					return index
				}
			}
		}
	}

	// Default to first category (usually Private Expenses)
	if len(ud.ExpenseModel.Categories) > 0 {
		return 0
	}

	return -1
}

// EnsureExpenseConsistency validates and fixes any data inconsistencies
func (ud *UserData) EnsureExpenseConsistency() {
	// Ensure category order matches actual categories
	actualCategories := make(map[string]bool)
	for _, category := range ud.ExpenseModel.Categories {
		actualCategories[category.ID] = true
	}

	// Remove obsolete entries from order tracking
	var newOrder []string
	for _, label := range ud.CategoryOrder {
		if actualCategories[label] {
			newOrder = append(newOrder, label)
		}
	}

	// Add missing categories to order
	for _, category := range ud.ExpenseModel.Categories {
		found := false
		for _, label := range newOrder {
			if label == category.ID {
				found = true
				break
			}
		}
		if !found {
			newOrder = append(newOrder, category.ID)
		}
	}

	ud.CategoryOrder = newOrder

	// Clean up category sources
	for label := range ud.CategorySources {
		if !actualCategories[label] {
			delete(ud.CategorySources, label)
		}
	}

	// Add missing sources
	for _, category := range ud.ExpenseModel.Categories {
		if _, exists := ud.CategorySources[category.ID]; !exists {
			ud.CategorySources[category.ID] = "user-added"
		}
	}

	ud.UpdateTimestamp()
}

// RemoveCategory removes a category and updates metadata
func (ud *UserData) RemoveCategory(categoryID string) {
	// Remove from expense model
	for i, category := range ud.ExpenseModel.Categories {
		if category.ID == categoryID {
			ud.ExpenseModel.Categories = append(
				ud.ExpenseModel.Categories[:i],
				ud.ExpenseModel.Categories[i+1:]...,
			)
			break
		}
	}

	// Remove from order tracking
	for i, label := range ud.CategoryOrder {
		if label == categoryID {
			ud.CategoryOrder = append(ud.CategoryOrder[:i], ud.CategoryOrder[i+1:]...)
			break
		}
	}

	// Remove from sources tracking
	delete(ud.CategorySources, categoryID)

	ud.UpdateTimestamp()
}
