# Modern Golang Telegram bot fremework

This fremework based on [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) and inspired by [telegraf.js](https://telegraf.js.org/)

## Telestage event driven framework using fsm.

### Installation

`go get github.com/askoldex/telestage`

### Quick Start

```go

package main

import (
	"github.com/askoldex/telestage"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	stateStore := NewStateStore()
	stg := telestage.NewStage(stateStore.Getter())
	mainScene := telestage.NewScene()
	messageScene := telestage.NewScene()
	stg.Add("main", mainScene)
	stg.Add("message", messageScene)

	messageScene.OnStart(func(ctx telestage.Context) {
		ctx.Reply("Hello from message scene, go back: /leave")
	})

	messageScene.OnCommand("leave", func(ctx telestage.Context) {
		stateStore.Set(ctx.Sender().ID, "main")
		ctx.Reply("Welcome in main scene, send: /start")
	})

	messageScene.OnMessage(func(ctx telestage.Context) {
		ctx.Reply("You in message scene")
	})

	mainScene.OnStart(func(ctx telestage.Context) {
		ctx.Reply("Hello world. send: /enter")
	})

	mainScene.OnCommand("enter", func(ctx telestage.Context) {
		stateStore.Set(ctx.Sender().ID, "message")
		ctx.Reply("Now send: /start")
	})

	mainScene.OnMessage(func(ctx telestage.Context) {
		ctx.Reply("Incorrect input")
	})

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	upds := bot.GetUpdatesChan(u)

	for upd := range upds {
		stg.Run(bot, upd)
	}
}

type stateStore struct {
	states map[int64]string
}

func NewStateStore() *stateStore {
	return &stateStore{
		states: map[int64]string{},
	}
}

func (ss *stateStore) Getter() telestage.StateGetter {
	return func(ctx telestage.Context) string {
		return ss.Get(ctx.Sender().ID)
	}
}

func (ss *stateStore) Get(userID int64) string {
	state, ok := ss.states[userID]
	if !ok {
		return "main"
	}

	return state
}

func (ss *stateStore) Set(userID int64, state string) {
	ss.states[userID] = state
}

```

### Scene middlewares

```go
mainScene.Use(func(ef telestage.EventFn) telestage.EventFn {
		return func(ctx telestage.Context) {
			if ctx.Message().Sticker == nil { // ignore if message is sticker
				ef(ctx)
			}
		}
	})

mainScene.OnMessage(func(ctx telestage.Context) {
    ctx.Reply("Hello") // answer on any message
})
```

### Event group middlewares

```go
mainScene.UseGroup(func(s *telestage.Scene) {
    // s is mainScene
    s.OnCommand("ban", func(ctx telestage.Context) {...})
    s.OnCommand("kick", func(ctx telestage.Context) {...})
}, func(ef telestage.EventFn) telestage.EventFn {
    return func(ctx telestage.Context) {
        if isAdmin(ctx.Sender().ID) {
            ef(ctx)
        }
    }
})
```

### Event middlewares

```go
mainScene.OnCommand("ping", func(ctx telestage.Context) {
    ctx.Reply("pong")
}, func(ef telestage.EventFn) telestage.EventFn {
    return func(ctx telestage.Context) {
        if ctx.Upd().FromChat().IsPrivate() {
            ef(ctx)
        } else {
            ctx.Reply("This command available only in private chat")
        }
    }
})
```

More examples see in examples folder.

