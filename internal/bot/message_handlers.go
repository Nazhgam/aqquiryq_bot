package bot

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

func (b *Bot) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		go b.handleUpdate(ctx, update)
	}
}

type Handler func(ctx context.Context, update tgbotapi.Update)

func (b *Bot) handleStart(ctx context.Context, update tgbotapi.Update) {
	userID := update.Message.From.ID

	allowed, err := b.userService.IsAllowed(ctx, userID)
	if err != nil || !allowed {
		msg := tgbotapi.NewMessage(userID, "⛔ У вас нет доступа")
		b.api.Send(msg)
		return
	}

	classes, err := b.contentService.GetClasses(ctx)
	if err != nil {
		return
	}

	if len(classes) == 0 {
		msg := tgbotapi.NewMessage(userID, "Контент пока отсутствует")
		b.api.Send(msg)
		return
	}

	var buttons [][]tgbotapi.InlineKeyboardButton

	for _, class := range classes {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("Класс %d", class),
			fmt.Sprintf("class:%d", class),
		)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(btn))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(userID, "📚 Выберите класс:")
	msg.ReplyMarkup = keyboard

	b.api.Send(msg)
}

func (b *Bot) handleAdmin(ctx context.Context, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	btn1 := tgbotapi.NewInlineKeyboardButtonData(
		"➕ Добавить контент",
		"admin:add_content",
	)

	btn2 := tgbotapi.NewInlineKeyboardButtonData(
		"❌ Удалить контент",
		"admin:delete_content_start",
	)

	btn3 := tgbotapi.NewInlineKeyboardButtonData(
		"👤 Добавить пользователя",
		"admin:add_user",
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btn1),
		tgbotapi.NewInlineKeyboardRow(btn2),
		tgbotapi.NewInlineKeyboardRow(btn3),
	)

	msg := tgbotapi.NewMessage(chatID, "🛠 Админ-панель")
	msg.ReplyMarkup = keyboard

	b.api.Send(msg)
}

func (b *Bot) handleAdminFSM(ctx context.Context, update tgbotapi.Update, session *AdminSession) {

	userID := update.Message.From.ID
	text := update.Message.Text

	switch session.State {

	case StateWaitingTitle:
		session.Title = text
		session.State = StateWaitingURL
		b.api.Send(tgbotapi.NewMessage(userID, "Введите Canva URL:"))

	case StateWaitingURL:
		session.URL = text
		session.State = StateWaitingClass
		msg := tgbotapi.NewMessage(userID, "Выберите класс:")
		msg.ReplyMarkup = classKeyboard()
		b.api.Send(msg)

	case StateWaitingLessonNumber:
		lessonNumber, err := strconv.Atoi(text)
		if err != nil || lessonNumber <= 0 {
			b.api.Send(tgbotapi.NewMessage(userID, "Введите корректный номер урока"))
			return
		}

		session.LessonNumber = lessonNumber

		_, err = b.contentService.AddContent(
			ctx,
			session.Title,
			session.URL,
			session.Class,
			session.Quarter,
			session.LessonNumber,
		)

		if err != nil {
			b.api.Send(tgbotapi.NewMessage(userID, "Ошибка при сохранении"))
			return
		}

		b.deleteSession(userID)
		b.api.Send(tgbotapi.NewMessage(userID, "✅ Контент добавлен"))

	case StateWaitInsertUserID:

		newUserID, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			b.api.Send(tgbotapi.NewMessage(userID, "wrong tg id"))
		}
		err = b.userService.AddUser(ctx, newUserID, uuid.NewString())
		if err != nil {
			b.api.Send(tgbotapi.NewMessage(userID, "can't insert new id. try again"))
			b.deleteSession(userID)
		}

		b.api.Send(tgbotapi.NewMessage(userID, "add new user success"))
	}
}

func (b *Bot) handleMainMenu(ctx context.Context, update tgbotapi.Update) {
	userID := update.Message.From.ID

	_, ok := b.getSession(update.Message.From.ID)
	if ok {
		newUserID, err := strconv.ParseInt(update.Message.Text, 10, 64)
		if err != nil {
			b.api.Send(tgbotapi.NewMessage(userID, "wrong tg id"))
		}
		err = b.userService.AddUser(ctx, newUserID, uuid.NewString())
		if err != nil {
			b.api.Send(tgbotapi.NewMessage(userID, "can't insert new id. try again"))
			b.deleteSession(userID)
		}
		return
	}

	isAdmin, _ := b.userService.IsAdmin(ctx, userID)

	msg := tgbotapi.NewMessage(userID, "👋 Добро пожаловать!")
	msg.ReplyMarkup = mainMenuKeyboard(isAdmin)

	b.api.Send(msg)
}

func mainMenuKeyboard(isAdmin bool) tgbotapi.ReplyKeyboardMarkup {

	btnBrowse := tgbotapi.NewKeyboardButton("📚 Презентации")

	rows := [][]tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButtonRow(btnBrowse),
	}

	if isAdmin {
		btnAdmin := tgbotapi.NewKeyboardButton("🛠 Админ")
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(btnAdmin))
	}

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = false

	return keyboard
}
