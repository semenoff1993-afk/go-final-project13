package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// nextDateHandler обрабатывает GET-запрос /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}

// NextDate вычисляет следующую дату на основе правила
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат начальной даты")
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	parts := strings.Fields(repeat)
	rule := parts[0]

	switch rule {
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат для y")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format(DateFormat), nil
			}
		}

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат для d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("недопустимое количество дней")
		}
		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				return date.Format(DateFormat), nil
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат для w")
		}
		daysStr := strings.Split(parts[1], ",")
		validDays := make(map[int]bool)
		for _, ds := range daysStr {
			d, err := strconv.Atoi(ds)
			if err != nil || d < 1 || d > 7 {
				return "", fmt.Errorf("недопустимый день недели: %s", ds)
			}
			// В Go Sunday = 0, а в задаче = 7
			if d == 7 {
				d = 0
			}
			validDays[d] = true
		}

		// Если стартовая дата в прошлом по отношению к now, начинаем поиск с now
		startSearch := date
		if date.Before(now) || date.Equal(now) {
			startSearch = now
		}
		// Начинаем проверять дни, начиная со следующего за startSearch
		searchDate := startSearch.AddDate(0, 0, 1)

		for i := 0; i < 7; i++ {
			if validDays[int(searchDate.Weekday())] {
				return searchDate.Format(DateFormat), nil
			}
			searchDate = searchDate.AddDate(0, 0, 1)
		}
		return "", fmt.Errorf("ошибка вычисления даты для правила w")

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("неверный формат для m")
		}

		daysStr := strings.Split(parts[1], ",")
		validDays := make([]int, 0)
		for _, ds := range daysStr {
			d, err := strconv.Atoi(ds)
			if err != nil || d < -2 || d == 0 || d > 31 {
				return "", fmt.Errorf("недопустимый день месяца: %s", ds)
			}
			validDays = append(validDays, d)
		}

		validMonths := make(map[int]bool)
		if len(parts) == 3 {
			monthsStr := strings.Split(parts[2], ",")
			for _, ms := range monthsStr {
				m, err := strconv.Atoi(ms)
				if err != nil || m < 1 || m > 12 {
					return "", fmt.Errorf("недопустимый месяц: %s", ms)
				}
				validMonths[m] = true
			}
		}

		// Начинаем поиск либо от стартовой даты, либо от now (что больше)
		startSearch := date
		if date.Before(now) {
			startSearch = now
		}

		// Ищем совпадения в течение максимум 5 лет
		for limit := 0; limit < 5*365; limit++ {
			searchDate := startSearch.AddDate(0, 0, limit)
			// Мы должны вернуть дату строго больше now
			if !searchDate.After(now) && !searchDate.After(date) {
				continue
			}

			if len(validMonths) > 0 && !validMonths[int(searchDate.Month())] {
				continue
			}

			for _, reqDay := range validDays {
				var targetDay int
				if reqDay > 0 {
					targetDay = reqDay
				} else {
					// Вычисление последних дней месяца (-1, -2)
					nextMonth := time.Date(searchDate.Year(), searchDate.Month()+1, 1, 0, 0, 0, 0, searchDate.Location())
					lastDayOfMonth := nextMonth.AddDate(0, 0, -1).Day()
					targetDay = lastDayOfMonth + reqDay + 1
				}

				if searchDate.Day() == targetDay {
					return searchDate.Format(DateFormat), nil
				}
			}
		}
		return "", fmt.Errorf("подходящая дата не найдена")

	default:
		return "", fmt.Errorf("неизвестное правило: %s", rule)
	}
}