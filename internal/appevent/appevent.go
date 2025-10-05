package appevent

type AppEventHandler func(data any)

type AppEvent struct {
	handlers []AppEventHandler
}

func (e *AppEvent) Subscribe(handler AppEventHandler) {
	e.handlers = append(e.handlers, handler)
}

func (e *AppEvent) Trigger(data any) {
	for _, handler := range e.handlers {
		handler(data)
	}
}
