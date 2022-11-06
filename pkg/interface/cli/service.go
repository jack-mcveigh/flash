package cli

import (
	"fmt"
	"strings"

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
			Name:  "flash",
			Usage: "a cli flashcard app",
			Flags: []cli.Flag{},
			Commands: []*cli.Command{
				addCmd(a), deleteCmd(d), getCmd(g), getAllCmd(g), updateCmd(u),
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
		ArgsUsage: "[group]",
		Flags: []cli.Flag{
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
		ArgsUsage: "[group]",
		Flags: []cli.Flag{
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
		Usage:   "Get flashcards in the group",
		Action: func(ctx *cli.Context) error {
			return getCards(ctx, g)
		},
		ArgsUsage: "[group]",
	}
}

func getAllCmd(g getting.Service) *cli.Command {
	return &cli.Command{
		Name:    "getall",
		Aliases: []string{"g"},
		Usage:   "Get all flashcards under the group",
		Action: func(ctx *cli.Context) error {
			return getAllCards(ctx, g)
		},
		ArgsUsage: "[group]",
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
		ArgsUsage: "[group]",
		Flags: []cli.Flag{
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
	group := groupFromArgs(ctx.Args())
	items := strings.Split(ctx.String("t"), ".")
	title := items[len(items)-1]
	if len(items) > 1 {
		if group != "" {
			group += "."
		}
		group += strings.Join(items[:1], ".")
	}

	return a.AddCard(
		group,
		adding.Card{
			Title: title,
			Desc:  ctx.String("d"),
		},
	)
}

func deleteCard(ctx *cli.Context, d deleting.Service) error {
	group := groupFromArgs(ctx.Args())
	items := strings.Split(ctx.String("t"), ".")
	title := items[len(items)-1]
	if len(items) > 1 {
		if group != "" {
			group += "."
		}
		group += strings.Join(items[:1], ".")
	}

	return d.DeleteCard(
		group,
		deleting.Card{
			Title: title,
		},
	)
}

func getCards(ctx *cli.Context, g getting.Service) error {
	group := groupFromArgs(ctx.Args())
	cards, err := g.GetCards(group)
	for i, c := range cards {
		fmt.Printf("\t%d) %s -> %s\n", i, c.Title, c.Desc)
	}
	return err
}

func getAllCards(ctx *cli.Context, g getting.Service) error {
	group := groupFromArgs(ctx.Args())
	cards, err := g.GetAllCards(group)
	for i, c := range cards {
		fmt.Printf("\t%d) %s -> %s\n", i, c.Title, c.Desc)
	}
	return err
}

func updateCard(ctx *cli.Context, u updating.Service) error {
	group := groupFromArgs(ctx.Args())
	items := strings.Split(ctx.String("t"), ".")
	title := items[len(items)-1]
	if len(items) > 1 {
		if group != "" {
			group += "."
		}
		group += strings.Join(items[:1], ".")
	}

	return u.UpdateCard(
		group,
		updating.Card{
			Title: title,
			Desc:  ctx.String("d"),
		},
	)
}

func groupFromArgs(a cli.Args) string {
	var group string
	if a.Len() == 1 {
		group = a.First()
	}
	return group
}
