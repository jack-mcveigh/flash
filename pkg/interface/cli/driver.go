package cli

import (
	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
)

type Service interface {
}

type service struct {
	adder   adding.Service
	deleter deleting.Service
	getter  getting.Service
}

func New(a adding.Service, d deleting.Service, g getting.Service) *service {
	return &service{a, d, g}
}
