package telestage

import (
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

var emptyStateGetter = func(ctx Context) string {
	return ""
}

func TestGetEvents(t *testing.T) {
	s := NewScene()
	s.OnMessage(func(_ Context) {})
	s.OnStart(func(_ Context) {})

	assert.Equal(t, len(s.GetEvents()), 2, "get events len should be equal with 2")
}

func TestOnCommand(t *testing.T) {
	s := NewScene()
	invoked := false
	cmd := "test"
	s.OnCommand(cmd, func(_ Context) {
		invoked = true
	})

	stage := NewStage(emptyStateGetter)
	stage.Add("", s)

	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/" + cmd,
			Entities: []tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: len(cmd) + 1, // plus slash
				},
			},
		},
	})

	assert.True(t, invoked, "they should be true if message command is "+cmd)
}

func TestSceneMiddleware(t *testing.T) {
	s := NewScene()
	s.Use(func(ef EventFn) EventFn {
		return func(ctx Context) {
			if ctx.Upd().FromChat().IsPrivate() {
				ef(ctx)
			}
		}
	})
	invoked := false
	s.OnMessage(func(_ Context) {
		invoked = true
	})
	stage := NewStage(emptyStateGetter)
	stage.Add("", s)

	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				Type: "group",
			},
		},
	})

	assert.False(t, invoked, "they should be false, because chat type is not private")

	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				Type: "private",
			},
		},
	})

	assert.True(t, invoked, "they should be true, because chat type is private")
}

func TestEventGroupMiddleware(t *testing.T) {
	s := NewScene()

	groupMiddlewareInvoked := false
	s.UseGroup(func(s *Scene) {
		s.OnSticker(func(_ Context) {})
	}, func(ef EventFn) EventFn {
		return func(ctx Context) {
			groupMiddlewareInvoked = true
			ef(ctx)
		}
	})

	stage := NewStage(emptyStateGetter)
	stage.Add("", s)

	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{},
	})

	assert.False(t, groupMiddlewareInvoked, "they should be false, because event OnSticker not invoked")

	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{
			Sticker: &tgbotapi.Sticker{},
		},
	})

	assert.True(t, groupMiddlewareInvoked, "they should be true, because event OnSticker invoked")
}

func TestEventMiddleware(t *testing.T) {
	s := NewScene()

	eventMiddlewareInvoked := false
	s.OnMessage(func(_ Context) {}, func(ef EventFn) EventFn {
		return func(ctx Context) {
			eventMiddlewareInvoked = true
		}
	})

	stage := NewStage(emptyStateGetter)
	stage.Add("", s)

	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{},
	})

	assert.True(t, eventMiddlewareInvoked, "they should be true, because event OnMessage invoked")
}

func TestOnStart(t *testing.T) {
	s := NewScene()
	invoked := false
	s.OnStart(func(_ Context) {
		invoked = true
	})

	stage := NewStage(emptyStateGetter)
	stage.Add("", s)

	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/start",
			Entities: []tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 6, // /start
				},
			},
		},
	})

	assert.True(t, invoked, "they should be true if message command is /start")
}

func TestOnPhoto(t *testing.T) {
	s := NewScene()
	invoked := false
	s.OnPhoto(func(_ Context) {
		invoked = true
	})

	stage := NewStage(emptyStateGetter)
	stage.Add("", s)

	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{
			Photo: []tgbotapi.PhotoSize{
				{
					FileID: "random_file_id",
				},
			},
		},
	})

	assert.True(t, invoked, "they should be true if message photo is not nil")
}

func TestOnSticker(t *testing.T) {
	s := NewScene()
	invoked := false
	s.OnSticker(func(_ Context) {
		invoked = true
	})

	stage := NewStage(emptyStateGetter)
	stage.Add("", s)

	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{
			Sticker: &tgbotapi.Sticker{
				FileID: "random_file_id",
			},
		},
	})

	assert.True(t, invoked, "they should be true if message sticker is not nil")
}

func TestOnMessage(t *testing.T) {
	s := NewScene()
	invoked := false
	s.OnMessage(func(ctx Context) {
		invoked = true
	})

	stage := NewStage(emptyStateGetter)

	stage.Add("", s)
	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{},
	})

	assert.True(t, invoked, "they should be true if message is not nil")
}

func TestOwnEvent(t *testing.T) {
	s := NewScene()
	messageTextContains := func(text string) EventDeterminant {
		return func(ctx Context) bool {
			return strings.Contains(ctx.Text(), text)
		}
	}

	invoked := false
	s.On(messageTextContains("hello"), func(_ Context) {
		invoked = true
	})
	stage := NewStage(emptyStateGetter)
	stage.Add("", s)

	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "hello, my name is John Doe",
		},
	})
	assert.True(t, invoked, "they should be true if message text contains 'hello'")

	invoked = false
	stage.Run(&tgbotapi.BotAPI{}, tgbotapi.Update{
		Message: &tgbotapi.Message{
			Caption: "hello, its my first picture",
		},
	})
	assert.True(t, invoked, "they should be true if message caption contains 'hello'")
}
