package main

import (
	"YandexAlice/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var userStates = make(map[string]*models.UserState)

func main() {
	http.HandleFunc("/post", handleAlice)
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func handleAlice(w http.ResponseWriter, r *http.Request) {
	var req models.AliceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Request from user %s: %s\n", req.Session.UserID, req.Request.OriginalUtterance)

	userID := req.Session.UserID
	if _, ok := userStates[userID]; !ok {
		userStates[userID] = &models.UserState{}
	}
	user := userStates[userID]

	res := models.AliceResponse{
		Version: req.Version,
	}
	res.Session.SessionID = req.Session.SessionID
	res.Session.MessageID = req.Session.MessageID
	res.Session.UserID = req.Session.UserID

	handleDialog(&req, &res, user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func handleDialog(req *models.AliceRequest, res *models.AliceResponse, user *models.UserState) {
	input := strings.ToLower(req.Request.OriginalUtterance)

	// Этап 1: Согласие с Police Agreement
	if !user.AgreedToPA {
		if contains(input, []string{"да", "согласен", "окей", "подтверждаю"}) {
			user.AgreedToPA = true
			res.Response.Text = "Спасибо за согласие с Police Agreement. Перейдём к следующему этапу!"
			return
		}

		res.Response.Text = "Перед началом вы должны согласиться с Police Agreement. В нём:\n\n- Неприкосновенность автора\n- Обязанности по созданию навыка\n- В случае нарушения: 1 мега сырок от РосАгроКомплекса.\n\nСогласны?"
		res.Response.Buttons = []models.Button{
			{Title: "Да", Hide: true},
			{Title: "Согласен", Hide: true},
		}
		return
	}

	// Отказ после согласия на создание навыка
	if user.AgreedToMakeSkill && contains(input, []string{"не хочу", "не буду", "отказ", "передумал"}) {
		res.Response.Text = "Согласно Police Agreement, вы обязаны выполнить создание навыка. В противном случае вас ждёт 1 мега сырок от РосАгроКомплекса."
		res.Response.Buttons = []models.Button{
			{Title: "Ладно, сделаю", Hide: true},
			{Title: "Прошу прощения", Hide: true},
		}
		return
	}

	// Согласие на создание навыка
	if strings.Contains(input, "сделаю к следующему уроку") {
		user.AgreedToMakeSkill = true
		days, hours, minutes, seconds := timeLeftUntilDeadline()
		res.Response.Text =
			"Хорошо, жду. У вас есть: " +
				formatTime(days, hours, minutes, seconds) +
				"\nПрошло: 0 секунд. Осталось: 5 секунд."
		return
	}

	// Запрос времени
	if contains(input, []string{"сколько осталось времени", "сколько времени осталось", "сколько прошло", "сколько прошло времени"}) {
		days, hours, minutes, seconds := timeLeftUntilDeadline()
		// Возвращаем прошедшее и оставшееся время
		res.Response.Text = "Прошло: 0 секунд. Осталось: " + formatTime(days, hours, minutes, seconds)
		return
	}

	// Финальное согласие
	if contains(input, []string{"ладно", "сделаю", "хорошо"}) {
		res.Response.Text = "Жду ваш навык!"
		res.Response.EndSession = true
		return
	}

	// Общий случай — агитация
	res.Response.Text = "Все говорят \"" + req.Request.OriginalUtterance + "\", а ты сделай навык для Алисы!"
	res.Response.Buttons = []models.Button{
		{Title: "Сделаю к следующему уроку", Hide: true},
		{Title: "Ладно", Hide: true},
	}
}

func timeLeftUntilDeadline() (int, int, int, int) {
	now := time.Now()
	var target time.Time

	switch now.Weekday() {
	case time.Monday:
		// до четверга
		target = now.AddDate(0, 0, int(time.Thursday-now.Weekday()))
	case time.Thursday:
		// до понедельника
		target = now.AddDate(0, 0, int(time.Monday+7-now.Weekday()))
	default:
		// просто +3 дня
		target = now.AddDate(0, 0, 3)
	}

	diff := target.Sub(now)
	days := int(diff.Hours()) / 24
	hours := int(diff.Hours()) % 24
	minutes := int(diff.Minutes()) % 60
	seconds := int(diff.Seconds()) % 60
	return days, hours, minutes, seconds
}

func formatTime(days, hours, minutes, seconds int) string {
	return strconv.Itoa(days) + " дн., " +
		strconv.Itoa(hours) + " ч., " +
		strconv.Itoa(minutes) + " мин., " +
		strconv.Itoa(seconds) + " сек."
}

func contains(input string, options []string) bool {
	for _, opt := range options {
		if strings.Contains(input, opt) {
			return true
		}
	}
	return false
}
