package cli

import (
	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	"github.com/urfave/cli/v2"
)

type Service interface {
	Run([]string) error
}

type service struct {
	App *cli.App
}

func New(a adding.Service, d deleting.Service, g getting.Service) *service {
	app := &cli.App{
		Name:  "Flash",
		Usage: "a cli flashcard app",
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add a flashcard",
				Action: func(ctx *cli.Context) error {
					return addCard(a)
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "Delete a flashcard",
				Action: func(ctx *cli.Context) error {
					return deleteCard(d)
				},
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "Get all flashcards",
				Action: func(ctx *cli.Context) error {
					return getCards(g)
				},
			},
		},
	}
	return &service{App: app}
}

func addCard(a adding.Service) error {
	return nil
}

func deleteCard(a deleting.Service) error {
	return nil
}

func getCards(a getting.Service) error {
	return nil
}
