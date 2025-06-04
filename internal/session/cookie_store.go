package session

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/smart-software-engineering/rate-calculator/internal/model"
)

const (
	sessionName       = "rate-calculator-session"
	scheduleKey       = "current_schedule"
	scheduleIDKey     = "schedule_id"
	expensesKey       = "expenses"
)

func init() {
	gob.Register(&model.WorkSchedule{})
	gob.Register(&model.ExpenseModel{})
	gob.Register(model.ExpenseType(""))
	gob.Register([]model.ExpenseCategory{})
	gob.Register(model.ExpenseCategory{})
	gob.Register([]model.ExpenseItem{})
	gob.Register(model.ExpenseItem{})
}

type CookieStore struct {
	store   *sessions.CookieStore
	authKey string
}

func NewCookieStore(authKey []byte, encryptionKey []byte, options *CookieOptions) *CookieStore {
	var store *sessions.CookieStore
	
	if encryptionKey != nil {
		store = sessions.NewCookieStore(authKey, encryptionKey)
	} else {
		store = sessions.NewCookieStore(authKey)
	}
	
	// Use provided options or fallback to defaults
	if options == nil {
		options = &CookieOptions{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   true,
			SameSite: SameSiteStrictMode,
		}
	}
	
	// Convert our SameSiteMode to http package SameSite
	var sameSite http.SameSite
	switch options.SameSite {
	case SameSiteLaxMode:
		sameSite = http.SameSiteLaxMode
	case SameSiteStrictMode:
		sameSite = http.SameSiteStrictMode
	case SameSiteNoneMode:
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteLaxMode
	}
	
	store.Options = &sessions.Options{
		Path:     options.Path,
		MaxAge:   options.MaxAge,
		HttpOnly: options.HttpOnly,
		Secure:   options.Secure,
		SameSite: sameSite,
	}

	return &CookieStore{
		store:   store,
		authKey: string(authKey),
	}
}

func (c *CookieStore) getSession(r *http.Request) (*sessions.Session, error) {
	return c.store.Get(r, sessionName)
}

func (c *CookieStore) GetSchedule(r *http.Request) (*model.WorkSchedule, error) {
	session, err := c.getSession(r)
	if err != nil {
		return nil, err
	}

	val, ok := session.Values[scheduleKey]
	if !ok {
		return model.NewWorkSchedule(), nil
	}

	schedule, ok := val.(*model.WorkSchedule)
	if !ok {
		return model.NewWorkSchedule(), nil
	}

	return schedule, nil
}

func (c *CookieStore) SetSchedule(w http.ResponseWriter, r *http.Request, schedule *model.WorkSchedule) error {
	session, err := c.getSession(r)
	if err != nil {
		return err
	}

	session.Values[scheduleKey] = schedule
	return session.Save(r, w)
}

func (c *CookieStore) GetScheduleID(r *http.Request) (string, error) {
	session, err := c.getSession(r)
	if err != nil {
		return "", err
	}

	val, ok := session.Values[scheduleIDKey]
	if !ok {
		return "", nil
	}

	id, ok := val.(string)
	if !ok {
		return "", nil
	}

	return id, nil
}

func (c *CookieStore) SetScheduleID(w http.ResponseWriter, r *http.Request, id string) error {
	session, err := c.getSession(r)
	if err != nil {
		return err
	}

	session.Values[scheduleIDKey] = id
	return session.Save(r, w)
}

func (c *CookieStore) GetExpenses(r *http.Request) (*model.ExpenseModel, error) {
	session, err := c.getSession(r)
	if err != nil {
		return nil, err
	}

	val, ok := session.Values[expensesKey]
	if !ok {
		return nil, nil
	}

	expenses, ok := val.(*model.ExpenseModel)
	if !ok {
		return nil, nil
	}

	return expenses, nil
}

func (c *CookieStore) SetExpenses(w http.ResponseWriter, r *http.Request, expenses *model.ExpenseModel) error {
	session, err := c.getSession(r)
	if err != nil {
		return err
	}

	session.Values[expensesKey] = expenses
	return session.Save(r, w)
}

func (c *CookieStore) Save(w http.ResponseWriter, r *http.Request) error {
	session, err := c.getSession(r)
	if err != nil {
		return err
	}

	return session.Save(r, w)
}

func (c *CookieStore) Clear(w http.ResponseWriter, r *http.Request) error {
	session, err := c.getSession(r)
	if err != nil {
		return err
	}

	session.Values = make(map[interface{}]interface{})
	return session.Save(r, w)
}

func (c *CookieStore) GetAuthKey() string {
	return c.authKey
}