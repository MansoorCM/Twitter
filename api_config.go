package main

import (
	"sync/atomic"

	"github.com/MansoorCM/Twitter/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}
