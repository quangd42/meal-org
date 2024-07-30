package services

import (
	"encoding/gob"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/quangd42/meal-planner/internal/database"
)

func NewSessionManager(store *database.Store) *scs.SessionManager {
	sm := scs.New()
	sm.Lifetime = 24 * time.Hour
	sm.Store = pgxstore.New(store.DB)

	// https://gist.github.com/alexedwards/d6eca7136f98ec12ad606e774d3abad3
	gob.Register(uuid.UUID{})

	return sm
}
