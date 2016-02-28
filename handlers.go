package calendar

import (
	"errors"
	"net/http"

	_ "github.com/itsabot/gcal"
	"github.com/labstack/echo"
)

// HandlerAPICalendarEventCreate creates an event on a user's calendaer. Right
// now support is built in for Google Calendar, but additional support will be
// added for Office/Exchange and others.
func HandlerAPICalendarEventCreate(c *echo.Context) error {
	var req struct {
		Title          string
		StartTime      uint64
		DurationInMins int
		AllDay         bool
		Recurring      bool
		// RecurringFreq  driver.RecurringFreq
		UserID uint64
	}
	if err := c.Bind(&req); err != nil {
		return jsonError(err)
	}
	if len(req.Title) == 0 {
		return jsonError(errors.New("Title cannot be blank"))
	}
	if req.DurationInMins <= 0 {
		return jsonError(errors.New("DurationInMins must be > 0"))
	}
	/*
		if req.RecurringFreq > driver.RecurringFreqYearly {
			return jsonError(errors.New("RecurringFreq is too high"))
		}
	*/
	if req.UserID <= 0 {
		return jsonError(errors.New("UserID must be > 0"))
	}
	/*
		ev := &driver.Event{}
		ev.Title = req.Title
		var tmp time.Time
		tmp = time.Unix(int64(req.StartTime), 0)
		ev.StartTime = &tmp
		ev.DurationInMins = req.DurationInMins
		ev.AllDay = req.AllDay
		ev.Recurring = req.Recurring
		ev.RecurringFreq = req.RecurringFreq
		ev.UserID = req.UserID
		client, err := cal.Client(db, ev.UserID)
		if err != nil {
			return jsonError(err)
		}
		if err = ev.Save(client); err != nil {
			return jsonError(err)
		}
	*/
	return nil
}

// HandlerOAuthConnectGoogleCalendar requests access for a Google Calendar.
func HandlerOAuthConnectGoogleCalendar(c *echo.Context) error {
	var req struct {
		UserID uint64
		Code   string
	}
	if err := c.Bind(&req); err != nil {
		return jsonError(err)
	}
	accessToken, idToken, err := cal.Exchange(req.Code)
	if err != nil {
		return jsonError(err)
	}
	gID, err := cal.DecodeIdToken(idToken)
	if err != nil {
		return jsonError(err)
	}
	q := `INSERT INTO sessions (userid, label, token) VALUES ($1, $2, $3)
	      ON CONFLICT (userid, label) DO UPDATE SET token=$3`
	_, err = db.Exec(q, req.UserID, "gcal_token", accessToken)
	if err != nil {
		return jsonError(err)
	}
	_, err = db.Exec(q, req.UserID, "gcal_id", gID)
	if err != nil {
		return jsonError(err)
	}
	// Ensure we can connect to the client
	_, err = cal.Client(db, req.UserID)
	if err != nil {
		return jsonError(err)
	}
	if err := c.JSON(http.StatusOK, nil); err != nil {
		return jsonError(err)
	}
	return nil
}

func HandlerOAuthDisconnectGoogleCalendar(c *echo.Context) error {
	var req struct {
		UserID uint64
	}
	if err := c.Bind(&req); err != nil {
		return jsonError(err)
	}
	var token string
	q := `SELECT token FROM sessions WHERE userid=$1 AND label='gcal_token'`
	if err := db.Get(&token, q, req.UserID); err != nil {
		return jsonError(err)
	}
	// Execute HTTP GET request to revoke current token
	url := "https://accounts.google.com/o/oauth2/revoke?token=" + token
	resp, err := http.Get(url)
	if err != nil {
		return jsonError(err)
	}
	defer resp.Body.Close()
	q = `DELETE FROM sessions WHERE userid=$1 AND label='gcal_token'`
	if _, err = db.Exec(q, req.UserID); err != nil {
		return jsonError(err)
	}
	if err := c.JSON(http.StatusOK, nil); err != nil {
		return jsonError(err)
	}
	return nil
}
