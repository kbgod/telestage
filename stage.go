package telestage

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrSceneNotFound = errors.New("scene not found")
)

type StateGetter func(Context) string

type Stage struct {
	scenes      map[string]*Scene
	stateGetter StateGetter
}

func NewStage(stateGetter StateGetter) *Stage {
	return &Stage{
		scenes:      map[string]*Scene{},
		stateGetter: stateGetter,
	}
}

func (s *Stage) Add(state string, scene *Scene) {
	s.scenes[state] = scene
}

func (s *Stage) Run(bot *tgbotapi.BotAPI, upd tgbotapi.Update) error {
	ctx := &NativeContext{
		bot: bot,
		upd: &upd,
	}

	state := s.stateGetter(ctx)
	scene, ok := s.scenes[state]
	if !ok {
		return fmt.Errorf("%w with name %s", ErrSceneNotFound, state)
	}

	events := scene.GetEvents()
	for _, e := range events {
		if e(ctx) {
			return nil
		}
	}

	return nil
}
