package cli

import (
	"fmt"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	"github.com/jmcveigh55/flash/pkg/core/updating"
	"github.com/urfave/cli/v2"
)

type Service interface {
	Run([]string) error
}

type service struct {
	app *cli.App
}

func New(a adding.Service, d deleting.Service, g getting.Service, u updating.Service) *service {
	return &service{
		app: &cli.App{
			Name:  "Flash",
			Usage: "a cli flashcard app",
			Flags: []cli.Flag{},
			Commands: []*cli.Command{
				addCmd(a), deleteCmd(d), getCmd(g), updateCmd(u),
			},
		},
	}
}

func (s *service) Run(arguments []string) error {
	return s.app.Run(arguments)
}

func addCmd(a adding.Service) *cli.Command {
	return &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "Add a flashcard",
		Action: func(ctx *cli.Context) error {
			return addCard(ctx, a)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "group",
				Aliases:  []string{"g"},
				Usage:    "Flashcard's group",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "title",
				Aliases:  []string{"t"},
				Usage:    "Flashcard's title",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "description",
				Aliases:  []string{"d"},
				Usage:    "Flashcard's Description",
				Required: true,
			},
		},
	}
}

func deleteCmd(d deleting.Service) *cli.Command {
	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "Delete a flashcard",
		Action: func(ctx *cli.Context) error {
			return deleteCard(ctx, d)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "group",
				Aliases:  []string{"g"},
				Usage:    "Flashcard's group",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "title",
				Aliases:  []string{"t"},
				Usage:    "Flashcard's title",
				Required: true,
			},
		},
	}
}

func getCmd(g getting.Service) *cli.Command {
	return &cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "Get all flashcards",
		Action: func(ctx *cli.Context) error {
			return getCards(ctx, g)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "group",
				Aliases:  []string{"g"},
				Usage:    "Flashcard's group",
				Required: true,
			},
		},
	}
}

func updateCmd(u updating.Service) *cli.Command {
	return &cli.Command{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "Update a flashcard's description",
		Action: func(ctx *cli.Context) error {
			return updateCard(ctx, u)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "group",
				Aliases:  []string{"g"},
				Usage:    "Flashcard's group",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "title",
				Aliases:  []string{"t"},
				Usage:    "Flashcard's title",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "description",
				Aliases:  []string{"d"},
				Usage:    "Flashcard's Description",
				Required: true,
			},
		},
	}
}

func addCard(ctx *cli.Context, a adding.Service) error {
	return a.AddCard(
		ctx.String("g"),
		adding.Card{
			Title: ctx.String("t"),
			Desc:  ctx.String("d"),
		},
	)
}

func deleteCard(ctx *cli.Context, d deleting.Service) error {
	return d.DeleteCard(
		ctx.String("g"),
		deleting.Card{
			Title: ctx.String("t"),
		},
	)
}

func getCards(ctx *cli.Context, g getting.Service) error {
	cards, err := g.GetCards(ctx.String("g"))
	for i, c := range cards {
		fmt.Printf("\t%d) %s -> %s\n", i, c.Title, c.Desc)
	}
	return err
}

func updateCard(ctx *cli.Context, u updating.Service) error {
	return u.UpdateCard(
		ctx.String("g"),
		updating.Card{
			Title: ctx.String("t"),
			Desc:  ctx.String("d"),
		},
	)
}
