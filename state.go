package main

import (
	"github.com/Romasav/gator/internal/config"
	"github.com/Romasav/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func newState(db *database.Queries, config *config.Config) *state {
	return &state{db, config}
}
