package combuuid

import (
	"time"

	"github.com/google/uuid"
)

const BLOCK_SIZE = 65536
const DEFAULT_PERIOD_SEC = 60

/*
NewUuid wraps (and returns) a "github.com/google/uuid.UUID" encoded using the
COMB optimization for UUIDs as primary keys in relational databases as shown in
https://www.2ndquadrant.com/en/blog/sequential-uuid-generators/

These Uuids sacrifice small amounts of collision resistance for improved write
and read performance for databases. They are almost certainly good enough for
"application"-level uniqueness, if not "universal"-level.
*/
func NewUuid() uuid.UUID {
	return newUuid(time.Now())
}

func newUuid(t time.Time) uuid.UUID {
	id := uuid.New() // id = 0x??, 0x??, ...

	prefix := uint16((t.Unix() / DEFAULT_PERIOD_SEC) % BLOCK_SIZE) // say, 0xAAFF

	id[0] = byte(prefix >> 8) // id = 0xAA 0x??
	id[1] = byte(prefix)      // id = 0xAA 0xFF

	return id
}
