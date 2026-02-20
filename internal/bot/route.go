package bot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		b.routeMessage(ctx, update)
	case update.CallbackQuery != nil:
		b.routeCallback(ctx, update)
	}
}

func (b *Bot) routeMessage(ctx context.Context, update tgbotapi.Update) {
	userID := update.Message.From.ID

	// если активна FSM
	if session, ok := b.getSession(userID); ok {
		b.handleAdminFSM(ctx, update, session)
		return
	}

	switch update.Message.Text {

	case "📚 Презентации":
		b.handleStart(ctx, update)

	case "🛠 Админ":
		b.handleAdmin(ctx, update)

	case "/start":
		b.handleMainMenu(ctx, update)

	default:
		b.handleMainMenu(ctx, update)
	}
}

func (b *Bot) routeCallback(ctx context.Context, update tgbotapi.Update) {
	parts := strings.Split(update.CallbackQuery.Data, ":")

	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "menu":
		b.handleMenuAction(ctx, update, parts)
	case "class":
		b.handleClass(ctx, update, parts)
	case "quarter":
		b.handleQuarter(ctx, update, parts)
	case "content":
		b.handleContent(ctx, update, parts)
	case "admin":
		b.handleAdminAction(ctx, update, parts)
	}
}
