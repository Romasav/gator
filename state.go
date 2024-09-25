package main

import "github.com/Romasav/gator/internal/config"

type state struct {
	config *config.Config
}

func newState(config *config.Config) *state {
	return &state{config}
}
