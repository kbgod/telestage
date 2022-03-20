package telestage

type Middleware func(EventFn) EventFn

func applyMiddleware(ef EventFn, middleware ...Middleware) EventFn {
	for i := len(middleware) - 1; i >= 0; i-- {
		ef = middleware[i](ef)
	}
	return ef
}
