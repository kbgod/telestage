package telestage

type EventFn func(Context)
type Event func(Context) bool
type EventDeterminant func(Context) bool

type Scene struct {
	events      []Event
	middlewares []Middleware
}

func NewScene() *Scene {
	return &Scene{}
}

func (s *Scene) GetEvents() []Event {
	return s.events
}

func (s *Scene) Use(mw ...Middleware) {
	s.middlewares = append(s.middlewares, mw...)
}

func (s *Scene) UseGroup(group func(*Scene), mw ...Middleware) {
	original := s.middlewares
	s.middlewares = append(s.middlewares, mw...)
	group(s)
	s.middlewares = original
}

// OnCommand handle the command specified by first argument
func (s *Scene) OnCommand(cmd string, ef EventFn, mw ...Middleware) {
	ef = applyMiddleware(ef, append(s.middlewares, mw...)...)
	s.events = append(s.events, func(ctx Context) bool {
		if ctx.Upd().Message == nil || ctx.Upd().Message.Command() != cmd {
			return false
		}
		ef(ctx)
		return true
	})
}

// OnMessage handle any message type (photo, text, sticker etc.)
func (s *Scene) OnMessage(ef EventFn, mw ...Middleware) {
	ef = applyMiddleware(ef, append(s.middlewares, mw...)...)
	s.events = append(s.events, func(ctx Context) bool {
		if ctx.Upd().Message == nil {
			return false
		}
		ef(ctx)
		return true
	})
}

// OnPhoto handle sending a photo
func (s *Scene) OnPhoto(ef EventFn, mw ...Middleware) {
	ef = applyMiddleware(ef, append(s.middlewares, mw...)...)
	s.events = append(s.events, func(ctx Context) bool {
		m := ctx.Message()
		if m == nil || len(m.Photo) == 0 {
			return false
		}
		ef(ctx)
		return true
	})
}

// OnSticker handle sending a sticker
func (s *Scene) OnSticker(ef EventFn, mw ...Middleware) {
	ef = applyMiddleware(ef, append(s.middlewares, mw...)...)
	s.events = append(s.events, func(ctx Context) bool {
		m := ctx.Message()
		if m == nil || m.Sticker == nil {
			return false
		}
		ef(ctx)
		return true
	})
}

// OnPhoto handle the "/start" command
func (s *Scene) OnStart(ef EventFn, mw ...Middleware) {
	s.OnCommand("start", ef, mw...)
}

// On handle the your own event determinator
func (s *Scene) On(determinant EventDeterminant, ef EventFn, mw ...Middleware) {
	ef = applyMiddleware(ef, append(s.middlewares, mw...)...)
	s.events = append(s.events, func(ctx Context) bool {
		if !determinant(ctx) {
			return false
		}
		ef(ctx)

		return true
	})
}
