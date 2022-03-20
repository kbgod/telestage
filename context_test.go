package telestage

import (
	"reflect"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestNativeContext_Bot(t *testing.T) {
	b := &tgbotapi.BotAPI{}
	nc := &NativeContext{bot: b}
	assert.Equal(t, b, nc.Bot(), "get bot from context")
}

func TestNativeContext_Upd(t *testing.T) {
	u := &tgbotapi.Update{}
	nc := &NativeContext{upd: u}
	assert.Equal(t, u, nc.Upd(), "get update from context")
}

func TestNativeContext_Set(t *testing.T) {
	type user struct{}
	u := &user{}
	k := "user"
	nc := &NativeContext{}
	nc.Set(k, u)
	assert.Equal(t, nc.store[k], u, "set context store value")
}

func TestNativeContext_Get(t *testing.T) {
	type user struct{}
	u := &user{}
	k := "user"
	nc := &NativeContext{}
	nc.store = map[string]interface{}{
		k: u,
	}
	assert.Equal(t, nc.Get(k), u, "get context store value")
}

func TestNativeContext_Message(t *testing.T) {
	m := &tgbotapi.Message{}
	pinnedMessage := &tgbotapi.Message{
		PinnedMessage: m,
	}
	tests := []struct {
		name string
		upd  *tgbotapi.Update
		want *tgbotapi.Message
	}{
		{
			name: "message in update",
			upd: &tgbotapi.Update{
				Message: m,
			},
			want: m,
		},
		{
			name: "message in callbackQuery",
			upd: &tgbotapi.Update{
				CallbackQuery: &tgbotapi.CallbackQuery{
					Message: m,
				},
			},
			want: m,
		},
		{
			name: "message in editedMessage",
			upd: &tgbotapi.Update{
				EditedMessage: m,
			},
			want: m,
		},
		{
			name: "message in channelPost",
			upd: &tgbotapi.Update{
				ChannelPost: m,
			},
			want: m,
		},
		{
			name: "message in channelPost (pinned)",
			upd: &tgbotapi.Update{
				ChannelPost: pinnedMessage,
			},
			want: m,
		},
		{
			name: "message in editedChannelPost",
			upd: &tgbotapi.Update{
				EditedChannelPost: m,
			},
			want: m,
		},
		{
			name: "empty message",
			upd:  &tgbotapi.Update{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc := &NativeContext{
				upd: tt.upd,
			}
			if got := nc.Message(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NativeContext.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNativeContext_Sender(t *testing.T) {
	u := &tgbotapi.User{}
	tests := []struct {
		name string
		upd  *tgbotapi.Update
		want *tgbotapi.User
	}{
		{
			name: "sender in message",
			upd: &tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: u,
				},
			},
			want: u,
		},
		{
			name: "sender in callbackQuery",
			upd: &tgbotapi.Update{
				CallbackQuery: &tgbotapi.CallbackQuery{
					From: u,
				},
			},
			want: u,
		},
		{
			name: "sender in inlineQuery",
			upd: &tgbotapi.Update{
				InlineQuery: &tgbotapi.InlineQuery{
					From: u,
				},
			},
			want: u,
		},
		{
			name: "sender in shippingQuery",
			upd: &tgbotapi.Update{
				ShippingQuery: &tgbotapi.ShippingQuery{
					From: u,
				},
			},
			want: u,
		},
		{
			name: "sender in precheckoutQuery",
			upd: &tgbotapi.Update{
				PreCheckoutQuery: &tgbotapi.PreCheckoutQuery{
					From: u,
				},
			},
			want: u,
		},
		{
			name: "sender in pollAnswer",
			upd: &tgbotapi.Update{
				PollAnswer: &tgbotapi.PollAnswer{
					User: *u,
				},
			},
			want: u,
		},
		{
			name: "sender in myChatMember",
			upd: &tgbotapi.Update{
				MyChatMember: &tgbotapi.ChatMemberUpdated{
					From: *u,
				},
			},
			want: u,
		},
		{
			name: "sender in chatMember",
			upd: &tgbotapi.Update{
				ChatMember: &tgbotapi.ChatMemberUpdated{
					From: *u,
				},
			},
			want: u,
		},
		{
			name: "sender in chatJoinRequest",
			upd: &tgbotapi.Update{
				ChatJoinRequest: &tgbotapi.ChatJoinRequest{
					From: *u,
				},
			},
			want: u,
		},
		{
			name: "empty user",
			upd:  &tgbotapi.Update{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc := &NativeContext{
				upd: tt.upd,
			}
			if got := nc.Sender(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NativeContext.Sender() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNativeContext_Chat(t *testing.T) {
	c := &tgbotapi.Chat{}
	tests := []struct {
		name string
		upd  *tgbotapi.Update
		want *tgbotapi.Chat
	}{
		{
			name: "chat in message",
			upd: &tgbotapi.Update{
				Message: &tgbotapi.Message{
					Chat: c,
				},
			},
			want: c,
		},
		{
			name: "chat in callbackQuery",
			upd: &tgbotapi.Update{
				CallbackQuery: &tgbotapi.CallbackQuery{
					Message: &tgbotapi.Message{
						Chat: c,
					},
				},
			},
			want: c,
		},
		{
			name: "chat in myChatMember",
			upd: &tgbotapi.Update{
				MyChatMember: &tgbotapi.ChatMemberUpdated{
					Chat: *c,
				},
			},
			want: c,
		},
		{
			name: "chat in chatMember",
			upd: &tgbotapi.Update{
				ChatMember: &tgbotapi.ChatMemberUpdated{
					Chat: *c,
				},
			},
			want: c,
		},
		{
			name: "chat in chatMember",
			upd: &tgbotapi.Update{
				ChatJoinRequest: &tgbotapi.ChatJoinRequest{
					Chat: *c,
				},
			},
			want: c,
		},
		{
			name: "empty chat",
			upd:  &tgbotapi.Update{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc := &NativeContext{
				upd: tt.upd,
			}
			if got := nc.Chat(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NativeContext.Chat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNativeContext_Text(t *testing.T) {
	txt := "message"
	tests := []struct {
		name string
		upd  *tgbotapi.Update
		want string
	}{
		{
			name: "text in message.text",
			upd: &tgbotapi.Update{
				Message: &tgbotapi.Message{
					Text: txt,
				},
			},
			want: txt,
		},
		{
			name: "text in message.caption",
			upd: &tgbotapi.Update{
				Message: &tgbotapi.Message{
					Caption: txt,
				},
			},
			want: txt,
		},
		{
			name: "text in empty message",
			upd:  &tgbotapi.Update{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc := &NativeContext{
				upd: tt.upd,
			}
			if got := nc.Text(); got != tt.want {
				t.Errorf("NativeContext.Text() = %v, want %v", got, tt.want)
			}
		})
	}
}
