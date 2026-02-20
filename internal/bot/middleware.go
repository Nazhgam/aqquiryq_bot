package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) whitelistMiddleware(next Handler) Handler {
	return func(ctx context.Context, update tgbotapi.Update) {

		var telegramID int64

		if update.Message != nil {
			telegramID = update.Message.From.ID
		} else if update.CallbackQuery != nil {
			telegramID = update.CallbackQuery.From.ID
		} else {
			return
		}

		allowed, err := b.userService.IsAllowed(ctx, telegramID)
		if err != nil {
			return
		}

		if !allowed {
			msg := tgbotapi.NewMessage(telegramID, "⛔ У вас нет доступа")
			b.api.Send(msg)
			return
		}

		next(ctx, update)
	}
}

func (b *Bot) requireWhitelist(next func(context.Context, tgbotapi.Update)) func(context.Context, tgbotapi.Update) {
	return func(ctx context.Context, update tgbotapi.Update) {

		var telegramID int64

		if update.Message != nil {
			telegramID = update.Message.From.ID
		} else if update.CallbackQuery != nil {
			telegramID = update.CallbackQuery.From.ID
		}

		allowed, err := b.userService.IsAllowed(ctx, telegramID)
		if err != nil || !allowed {
			msg := tgbotapi.NewMessage(telegramID, "⛔ У вас нет доступа")
			b.api.Send(msg)
			return
		}

		next(ctx, update)
	}
}

func (b *Bot) IsUserMember(userID int64, channelID int64) (bool, error) {
	chatConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: channelID,
			UserID: userID,
		},
	}
	member, err := b.api.GetChatMember(chatConfig)
	if err != nil {
		return false, err
	}

	// Status can be creator, administrator, member, restricted, left, or kicked
	if member.Status == "member" || member.Status == "administrator" || member.Status == "creator" {
		return true, nil
	}
	return false, nil
}

func (b *Bot) PostPresentationToChannel(channelID int64, title, presentationURL string) error {
	// Create a formatted message: [Title](URL)
	// Using HTML parse mode for cleaner links: <a href="URL">Title</a>
	msgText := fmt.Sprintf("📚 <b>%s</b>\n\n👇 Көру үшін басыңыз:\n<a href=\"%s\">%s</a>", title, presentationURL, title)

	msg := tgbotapi.NewMessage(channelID, msgText)
	msg.ParseMode = "HTML"

	_, err := b.api.Send(msg)
	return err
}
