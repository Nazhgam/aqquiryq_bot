package bot

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleClass(ctx context.Context, update tgbotapi.Update, parts []string) {
	class, _ := strconv.Atoi(parts[1])

	quarters, err := b.contentService.GetQuarters(ctx, class)
	if err != nil || len(quarters) == 0 {
		return
	}

	var buttons [][]tgbotapi.InlineKeyboardButton

	for _, q := range quarters {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("Четверть %d", q),
			fmt.Sprintf("quarter:%d:%d", class, q),
		)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(btn))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "📅 Выберите четверть:")
	msg.ReplyMarkup = keyboard

	b.api.Send(msg)
}

func (b *Bot) handleQuarter(ctx context.Context, update tgbotapi.Update, parts []string) {
	class, _ := strconv.Atoi(parts[1])
	quarter, _ := strconv.Atoi(parts[2])

	contents, err := b.contentService.GetContents(ctx, class, quarter)
	if err != nil || len(contents) == 0 {
		return
	}

	var buttons [][]tgbotapi.InlineKeyboardButton

	for _, c := range contents {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			c.Title,
			fmt.Sprintf("content:%d", c.ID),
		)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(btn))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "📊 Выберите презентацию:")
	msg.ReplyMarkup = keyboard

	b.api.Send(msg)
}

func (b *Bot) handleContent(ctx context.Context, update tgbotapi.Update, parts []string) {
	id, _ := strconv.ParseInt(parts[1], 10, 64)

	content, err := b.contentService.GetContent(ctx, id)
	if err != nil || content == nil {
		return
	}

	msgText := fmt.Sprintf("📚 <b>%s</b>\n\n👇 Көру үшін басыңыз:\n<a href=\"%s\">%s</a>", content.Title, content.CanvaURL, content.Title)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msgText)
	msg.ParseMode = "HTML"

	b.api.Send(msg)
}

func (b *Bot) handleAdminAction(ctx context.Context, update tgbotapi.Update, parts []string) {

	userID := update.CallbackQuery.From.ID

	isAdmin, err := b.userService.IsAdmin(ctx, userID)
	if err != nil || !isAdmin {
		b.api.Send(tgbotapi.NewMessage(userID, "⛔ Только для администраторов"))
		return
	}

	action := parts[1]

	switch action {

	case "add_content":
		b.handleAdminAddContent(ctx, update)

	case "delete_content_start":
		b.handleAdminDeleteContent(ctx, update)

	case "delete_content":
		b.handleDeleteContentSelect(ctx, update, parts)

	case "confirm_delete":
		b.handleConfirmDelete(ctx, update, parts)

	case "class":
		b.handleAdminSelectClass(ctx, update, parts)

	case "quarter":
		b.handleAdminSelectQuarter(ctx, update, parts)

	case "add_user":
		b.handleAdminAddUser(ctx, update)

	case "cancel":
		b.deleteSession(userID)
	default:
		b.handleAdminAddUser(ctx, update)
		//b.handleAdmin(ctx, update)
	}
}

func (b *Bot) handleAdminAddContent(ctx context.Context, update tgbotapi.Update) {
	userID := update.CallbackQuery.From.ID

	b.setSession(userID, &AdminSession{
		State: StateWaitingTitle,
	})

	msg := tgbotapi.NewMessage(userID, "Введите название контента:")
	msg.ReplyMarkup = cancelKeyboard()
	b.api.Send(msg)
}

func (b *Bot) handleAdminAddUser(ctx context.Context, update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID

	b.setSession(chatID, &AdminSession{
		State: StateWaitInsertUserID,
	})

	msg := tgbotapi.NewMessage(chatID, "Введите Telegram ID пользователя:")
	msg.ReplyMarkup = cancelKeyboard()
	b.api.Send(msg)

}

func cancelKeyboard() tgbotapi.InlineKeyboardMarkup {
	btn := tgbotapi.NewInlineKeyboardButtonData(
		"❌ Отмена",
		"admin:cancel",
	)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btn),
	)
}

func (b *Bot) handleAdminSelectClass(ctx context.Context, update tgbotapi.Update, parts []string) {

	userID := update.CallbackQuery.From.ID
	session, ok := b.getSession(userID)
	if !ok {
		return
	}

	class, _ := strconv.Atoi(parts[2])
	session.Class = class

	switch session.State {

	case StateWaitingClass:
		session.State = StateWaitingQuarter

		msg := tgbotapi.NewMessage(userID, "Выберите четверть:")
		msg.ReplyMarkup = quarterKeyboard()
		b.api.Send(msg)

	case StateDeleteSelectClass:
		session.State = StateDeleteSelectQuarter

		msg := tgbotapi.NewMessage(userID, "Выберите четверть для удаления:")
		msg.ReplyMarkup = quarterKeyboard()
		b.api.Send(msg)
	}
}

func (b *Bot) handleAdminSelectQuarter(ctx context.Context, update tgbotapi.Update, parts []string) {

	userID := update.CallbackQuery.From.ID
	session, ok := b.getSession(userID)
	if !ok {
		return
	}

	quarter, _ := strconv.Atoi(parts[2])
	session.Quarter = quarter

	switch session.State {

	case StateWaitingQuarter:
		session.State = StateWaitingLessonNumber

		msg := tgbotapi.NewMessage(userID, "Введите номер урока:")
		msg.ReplyMarkup = cancelKeyboard()
		b.api.Send(msg)

	case StateDeleteSelectQuarter:
		session.State = StateDeleteSelectContent

		contents, err := b.contentService.GetContents(
			ctx,
			session.Class,
			session.Quarter,
		)
		if err != nil || len(contents) == 0 {
			b.api.Send(tgbotapi.NewMessage(userID, "Контент не найден"))
			return
		}

		var rows [][]tgbotapi.InlineKeyboardButton

		for _, c := range contents {
			btn := tgbotapi.NewInlineKeyboardButtonData(
				c.Title,
				fmt.Sprintf("admin:delete_content:%d", c.ID),
			)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "admin:cancel"),
		))

		keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

		msg := tgbotapi.NewMessage(userID, "Выберите контент для удаления:")
		msg.ReplyMarkup = keyboard

		b.api.Send(msg)
	}
}
func (b *Bot) handleAdminDeleteContent(ctx context.Context, update tgbotapi.Update) {

	userID := update.CallbackQuery.From.ID

	b.setSession(userID, &AdminSession{
		State: StateDeleteSelectClass,
	})

	keyboard := classKeyboard()

	msg := tgbotapi.NewMessage(userID, "Выберите класс для удаления:")
	msg.ReplyMarkup = keyboard

	b.api.Send(msg)
}

func (b *Bot) handleDeleteContentSelect(ctx context.Context, update tgbotapi.Update, parts []string) {

	userID := update.CallbackQuery.From.ID
	session, ok := b.getSession(userID)
	if !ok || session.State != StateDeleteSelectContent {
		return
	}

	contentID, _ := strconv.ParseInt(parts[2], 10, 64)

	session.State = StateDeleteConfirm
	session.Title = fmt.Sprintf("%d", contentID)

	btnYes := tgbotapi.NewInlineKeyboardButtonData(
		"✅ Подтвердить",
		fmt.Sprintf("admin:confirm_delete:%d", contentID),
	)

	btnCancel := tgbotapi.NewInlineKeyboardButtonData(
		"❌ Отмена",
		"admin:cancel",
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(btnYes),
		tgbotapi.NewInlineKeyboardRow(btnCancel),
	)

	msg := tgbotapi.NewMessage(userID, "Подтвердить удаление?")
	msg.ReplyMarkup = keyboard

	b.api.Send(msg)
}

func (b *Bot) handleConfirmDelete(ctx context.Context, update tgbotapi.Update, parts []string) {

	userID := update.CallbackQuery.From.ID
	session, ok := b.getSession(userID)
	if !ok || session.State != StateDeleteConfirm {
		return
	}

	contentID, _ := strconv.ParseInt(parts[2], 10, 64)

	err := b.contentService.DeleteContent(ctx, contentID)
	if err != nil {
		b.api.Send(tgbotapi.NewMessage(userID, "Ошибка при удалении"))
		return
	}

	b.deleteSession(userID)
	b.api.Send(tgbotapi.NewMessage(userID, "🗑 Контент удалён"))
}

func (b *Bot) handleMenuAction(ctx context.Context, update tgbotapi.Update, parts []string) {

	userID := update.CallbackQuery.From.ID

	if len(parts) < 2 {
		return
	}

	switch parts[1] {

	case "browse":
		// вызываем старую логику показа классов
		msg := tgbotapi.NewMessage(userID, "📚 Загружаем классы...")
		b.api.Send(msg)

		// имитируем /start логику
		fakeUpdate := tgbotapi.Update{
			Message: &tgbotapi.Message{
				From: &tgbotapi.User{ID: userID},
			},
		}
		b.handleStart(ctx, fakeUpdate)

	case "admin":
		isAdmin, _ := b.userService.IsAdmin(ctx, userID)
		if !isAdmin {
			return
		}
		b.handleAdmin(ctx, update)
	}
}
