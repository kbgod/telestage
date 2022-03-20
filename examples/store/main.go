package main

import (
	"fmt"
	"log"
	"os"

	"github.com/askoldex/telestage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	stateStore := NewStateStore()
	stg := telestage.NewStage(stateStore.Getter())
	mainScene := telestage.NewScene()
	stg.Add("main", mainScene)

	mainScene.Use(addUserBalance)

	mainScene.OnMessage(func(ctx telestage.Context) {
		account := ctx.Get("account").(*account)
		ctx.Reply(fmt.Sprintf("Your balance: %d", account.Balance))
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

type account struct {
	Balance int
}

func addUserBalance(ef telestage.EventFn) telestage.EventFn {
	return func(ctx telestage.Context) {
		ctx.Set("account", &account{500})
		ef(ctx)
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
