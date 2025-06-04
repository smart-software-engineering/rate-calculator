package model

type ExpenseType string

const (
	Monthly    ExpenseType = "monthly"
	Yearly     ExpenseType = "yearly"
	Percentage ExpenseType = "percentage"
)

type ExpenseItem struct {
	Label  string
	Amount int
	Type   ExpenseType
}

type ExpenseCategory struct {
	Label string
	Items []ExpenseItem
}

type ExpenseModel struct {
	Categories []ExpenseCategory
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

func CreateSampleExpenseModel() *ExpenseModel {
	expenses := NewExpenseModel()

	privateExpenses := ExpenseCategory{
		Label: "Private Expenses",
		Items: []ExpenseItem{
			{Label: "Rent + Household ", Amount: 3000, Type: Monthly},
			{Label: "Utilities", Amount: 200, Type: Monthly},
			{Label: "Groceries", Amount: 400, Type: Monthly},
			{Label: "Clothing", Amount: 250, Type: Monthly},
			{Label: "Internet + Phone", Amount: 70, Type: Monthly},
			{Label: "Mobility", Amount: 70, Type: Monthly},
			{Label: "Diverse", Amount: 50, Type: Monthly},
		},
	}
	expenses.AddCategory(privateExpenses)

	professionalExpenses := ExpenseCategory{
		Label: "Professional Expenses",
		Items: []ExpenseItem{
			{Label: "Office", Amount: 200, Type: Monthly},
			{Label: "Software + Hardware", Amount: 100, Type: Monthly},
			{Label: "Communication", Amount: 100, Type: Yearly},
			{Label: "Insurances", Amount: 100, Type: Monthly},
			{Label: "Fees, Donations", Amount: 50, Type: Yearly},
			{Label: "Training & Education", Amount: 200, Type: Yearly},
			{Label: "Mobility", Amount: 1000, Type: Yearly},
			{Label: "Ads", Amount: 250, Type: Yearly},
		},
	}
	expenses.AddCategory(professionalExpenses)

	financialGoals := ExpenseCategory{
		Label: "Financial Goals",
		Items: []ExpenseItem{
			{Label: "Savings", Amount: 20, Type: Percentage},
			{Label: "Retirement Fund", Amount: 15, Type: Percentage},
			{Label: "Emergency Fund", Amount: 10, Type: Percentage},
		},
	}
	expenses.AddCategory(financialGoals)

	taxes := ExpenseCategory{
		Label: "Taxes",
		Items: []ExpenseItem{
			{Label: "Pension", Amount: 25, Type: Percentage},
			{Label: "Health Insurance", Amount: 10, Type: Percentage},
			{Label: "Income Tax", Amount: 10, Type: Percentage},
		},
	}
	expenses.AddCategory(taxes)

	return expenses
}
