package main

import (
	"log"
	"os"

	"github.com/askoldex/telestage"

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
