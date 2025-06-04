package session

import (
	"net/http"

	"github.com/smart-software-engineering/rate-calculator/internal/model"
)

// SameSiteMode represents the SameSite cookie attribute
type SameSiteMode int

const (
	SameSiteDefaultMode SameSiteMode = iota
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)

// CookieOptions represents configuration options for cookies
type CookieOptions struct {
	Path     string
	MaxAge   int
	HttpOnly bool
	Secure   bool
	SameSite SameSiteMode
}

type Store interface {
	GetSchedule(r *http.Request) (*model.WorkSchedule, error)
	SetSchedule(w http.ResponseWriter, r *http.Request, schedule *model.WorkSchedule) error
	GetScheduleID(r *http.Request) (string, error)
	SetScheduleID(w http.ResponseWriter, r *http.Request, id string) error
	GetExpenses(r *http.Request) (*model.ExpenseModel, error)
	SetExpenses(w http.ResponseWriter, r *http.Request, expenses *model.ExpenseModel) error

	GetUserData(r *http.Request) (*model.UserData, error)
	SetUserData(w http.ResponseWriter, r *http.Request, userData *model.UserData) error
	GetOrCreateUserData(w http.ResponseWriter, r *http.Request) (*model.UserData, error)

	Save(w http.ResponseWriter, r *http.Request) error
	Clear(w http.ResponseWriter, r *http.Request) error

	// GetAuthKey returns the session authentication key used for CSRF protection
	GetAuthKey() string
}
