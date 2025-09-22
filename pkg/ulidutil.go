package pkg

import (
	cryptoRand "crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	entropyMu sync.Mutex
	entropy   = ulid.Monotonic(cryptoRand.Reader, 0)
)

// New generates a new ULID string using current UTC time.
// Example output: "01HZ3FXCK9DF06JHBAAXM1T1H5"
func NewUlid() string {
	return MustNew(time.Now().UTC())
}

// MustNew generates a ULID from the provided time.
// Panics on failure (safe in most applications).
func MustNew(t time.Time) string {
	entropyMu.Lock()
	defer entropyMu.Unlock()

	id, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {
		panic(fmt.Errorf("ulid generation failed: %w", err))
	}
	return id.String()
}

// Parse parses a ULID string into a ULID type.
// Returns error if format is invalid.
func Parse(id string) (ulid.ULID, error) {
	return ulid.Parse(id)
}
