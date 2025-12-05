package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const Layout = "20060102"

func AfterNow(date, now time.Time) bool {
	return date.Year() > now.Year() ||
		(date.Year() == now.Year() && date.Month() > now.Month()) ||
		(date.Year() == now.Year() && date.Month() == now.Month() && date.Day() > now.Day())
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	date, err := time.Parse(Layout, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date format: %w", err)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat format")
	}

	for {
		switch parts[0] {
		case "d":
			if len(parts) < 2 {
				return "", fmt.Errorf("days count is wrong for 'd' rule")
			}
			days, err := strconv.Atoi(parts[1])
			if err != nil {
				return "", fmt.Errorf("invalid days format: %w", err)
			}
			if days <= 0 || days > 400 {
				return "", errors.New("days must be between 1 and 400")
			}
			date = date.AddDate(0, 0, days)

		case "y":
			date = date.AddDate(1, 0, 0)

		default:
			return "", fmt.Errorf("unknown repeat rule: %s", parts[0])
		}

		if AfterNow(date, now) {
			break
		}
	}
	resDate := date.Format(Layout)
	return resDate, nil
}

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "error: Method %s not allowed. Use GET", r.Method)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	query := r.URL.Query()
	date := query.Get("date")
	repeat := query.Get("repeat")
	nowStr := query.Get("now")

	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "error: Parameter 'date' is required")
		return
	}
	if repeat == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "error: Parameter 'repeat' is required")
		return
	}

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		parsedNow, err := time.Parse(Layout, nowStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "error: Invalid 'now' format: %v", err)
			return
		}
		now = parsedNow
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprint(w, nextDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error to return next date")
		return
	}
}
