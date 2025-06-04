package model

import "github.com/google/uuid"

type ExpenseType string

const (
	Monthly    ExpenseType = "monthly"
	Yearly     ExpenseType = "yearly"
	Percentage ExpenseType = "percentage"
)

type ExpenseItem struct {
	ID     string      `json:"id"`
	Label  string      `json:"label"`
	Amount int         `json:"amount"`
	Type   ExpenseType `json:"type"`
}

type ExpenseCategory struct {
	ID    string        `json:"id"`
	Label string        `json:"label"`
	Items []ExpenseItem `json:"items"`
}

type ExpenseModel struct {
	Categories []ExpenseCategory `json:"categories"`
}

type ExpenseCategoryData struct {
	Category        ExpenseCategory
	PercentageTotal int
	YearlyTotal     int
	CSRFField       interface{}
}

func NewExpenseModel() *ExpenseModel {
	return &ExpenseModel{
		Categories: []ExpenseCategory{},
	}
}

func (m *ExpenseModel) AddCategory(category ExpenseCategory) {
	m.Categories = append(m.Categories, category)
}

func (m *ExpenseModel) AddItemToCategory(categoryIndex int, item ExpenseItem) {
	if categoryIndex >= 0 && categoryIndex < len(m.Categories) {
		m.Categories[categoryIndex].Items = append(m.Categories[categoryIndex].Items, item)
	}
}

// NewExpenseItem creates a new expense item with generated ID
func NewExpenseItem(label string, amount int, expenseType ExpenseType) ExpenseItem {
	return ExpenseItem{
		ID:     uuid.New().String(),
		Label:  label,
		Amount: amount,
		Type:   expenseType,
	}
}

// NewExpenseCategory creates a new expense category with generated ID
func NewExpenseCategory(label string) ExpenseCategory {
	return ExpenseCategory{
		ID:    uuid.New().String(),
		Label: label,
		Items: []ExpenseItem{},
	}
}

func CreateSampleExpenseModel() *ExpenseModel {
	expenses := NewExpenseModel()

	privateExpenses := NewExpenseCategory("Private Expenses")
	privateExpenses.Items = []ExpenseItem{
		NewExpenseItem("Rent + Household", 3000, Monthly),
		NewExpenseItem("Utilities", 200, Monthly),
		NewExpenseItem("Groceries", 400, Monthly),
		NewExpenseItem("Clothing", 250, Monthly),
		NewExpenseItem("Internet + Phone", 70, Monthly),
		NewExpenseItem("Mobility", 70, Monthly),
		NewExpenseItem("Diverse", 50, Monthly),
	}
	expenses.AddCategory(privateExpenses)

	professionalExpenses := NewExpenseCategory("Professional Expenses")
	professionalExpenses.Items = []ExpenseItem{
		NewExpenseItem("Office", 200, Monthly),
		NewExpenseItem("Software + Hardware", 100, Monthly),
		NewExpenseItem("Communication", 100, Yearly),
		NewExpenseItem("Insurances", 100, Monthly),
		NewExpenseItem("Fees, Donations", 50, Yearly),
		NewExpenseItem("Training & Education", 200, Yearly),
		NewExpenseItem("Mobility", 1000, Yearly),
		NewExpenseItem("Ads", 250, Yearly),
	}
	expenses.AddCategory(professionalExpenses)

	financialGoals := NewExpenseCategory("Financial Goals")
	financialGoals.Items = []ExpenseItem{
		NewExpenseItem("Savings", 20, Percentage),
		NewExpenseItem("Retirement Fund", 15, Percentage),
		NewExpenseItem("Emergency Fund", 10, Percentage),
	}
	expenses.AddCategory(financialGoals)

	taxes := NewExpenseCategory("Taxes")
	taxes.Items = []ExpenseItem{
		NewExpenseItem("Pension", 25, Percentage),
		NewExpenseItem("Health Insurance", 10, Percentage),
		NewExpenseItem("Income Tax", 10, Percentage),
	}
	expenses.AddCategory(taxes)

	return expenses
}
