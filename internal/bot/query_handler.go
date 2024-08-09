package bot

import (
	"errors"
	"github.com/alir32a/jupiter/pkg/tg"
)

type QueryHandler func(callbackQueryID string, query Query) error

type QueryCommander struct {
	queries         map[string]QueryHandler
	NotFoundHandler func() error
	ErrHandler      func(err error) error
}

func NewQueryCommander() *QueryCommander {
	return &QueryCommander{
		queries: make(map[string]QueryHandler),
	}
}

func (q *QueryCommander) Register(name string, handler QueryHandler) {
	q.queries[name] = handler
}

func (q *QueryCommander) Handle(callbackQuery tg.CallbackQuery) error {
	query, err := Unmarshal(callbackQuery.Data)
	if err != nil {
		if q.ErrHandler != nil {
			return q.ErrHandler(err)
		}

		return err
	}

	handler, ok := q.queries[query.Resource]
	if !ok {
		if q.NotFoundHandler != nil {
			return q.NotFoundHandler()
		}

		return errors.New("resource not found")
	}

	return handler(callbackQuery.ID, *query)
}
