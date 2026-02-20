package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type State string

const (
	StateWaitingTitle        State = "waiting_title"
	StateWaitingURL          State = "waiting_url"
	StateWaitingClass        State = "waiting_class"
	StateWaitingQuarter      State = "waiting_quarter"
	StateWaitingLessonNumber State = "waiting_lesson_number"

	StateDeleteSelectClass   State = "delete_select_class"
	StateDeleteSelectQuarter State = "delete_select_quarter"
	StateDeleteSelectContent State = "delete_select_content"
	StateDeleteConfirm       State = "delete_confirm"

	StateWaitInsertUserID State = "wait_insert_user_id"
)

type AdminSession struct {
	State        State
	Title        string
	URL          string
	Class        int
	Quarter      int
	LessonNumber int
}

func (b *Bot) getSession(userID int64) (*AdminSession, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	session, ok := b.adminSessions[userID]
	return session, ok
}

func (b *Bot) setSession(userID int64, session *AdminSession) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.adminSessions[userID] = session
}

func (b *Bot) deleteSession(userID int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.adminSessions, userID) // ✅ ПРАВИЛЬНО
}

func classKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for i := 1; i <= 11; i++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d класс", i),
			fmt.Sprintf("admin:class:%d", i),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	cancelBtn := tgbotapi.NewInlineKeyboardButtonData(
		"❌ Отмена",
		"admin:cancel",
	)

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelBtn))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func quarterKeyboard() tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton

	for i := 1; i <= 4; i++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d четверть", i),
			fmt.Sprintf("admin:quarter:%d", i),
		)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	cancelBtn := tgbotapi.NewInlineKeyboardButtonData(
		"❌ Отмена",
		"admin:cancel",
	)

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelBtn))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
