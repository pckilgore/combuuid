package combuuid

import (
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func asHashBucketNumber(id uuid.UUID) uint16 {
	a, b := id[0], id[1]
	return (uint16(a) << 8) | uint16(b)
}

func Test_newUuid(t *testing.T) {
	dayIWroteThisIsh := int64(1660513)
	timeVal := time.Unix(dayIWroteThisIsh, 0)

	// Allocate UUIDS in every possible bin, and two in the first to check
	// wrapping works.
	var uuids [BLOCK_SIZE]uuid.UUID

	for i := range uuids {
		uuids[i] = newUuid(timeVal)
		timeVal = timeVal.Add(60 * time.Second)
	}

	for j := 1; j < BLOCK_SIZE; j++ {
		prev := asHashBucketNumber(uuids[j-1])
		curr := asHashBucketNumber(uuids[j])

		if curr != 0 {
			// Hashes should always increment every 60 seconds...
			assert.Greater(t, curr, prev)
		} else {
			// And wrap over at size_uint16
			assert.Equal(t, uint16(BLOCK_SIZE-1), prev)
		}
	}
}

func Fuzz_NewUuid(f *testing.F) {
	period := time.Duration(DEFAULT_PERIOD_SEC) * time.Second
	testTimePeriodInSeconds := int64(BLOCK_SIZE * DEFAULT_PERIOD_SEC)

	f.Add(int64(1))              // Guarantee we hit "wrap-check" assertion
	f.Add(int64(1660513))        // Day test written
	f.Add(int64(math.MaxInt32))  // 2038 problem
	f.Add(int64(math.MaxUint32)) // 2106 problem

	// See you in several hundred billion years.
	// Note, we do push time forward in the test below, so we leave ourselves room
	// so that the test itself doesn't overflow Int64 as it runs. This is
	// billions of years after the Earth itself will be garbage collected by the
	// sun (black hole or supernova, dealers choice), so whatever.
	postEarth := math.MaxInt64 - testTimePeriodInSeconds
	f.Add(postEarth)
	f.Add(postEarth + 1)

	f.Fuzz(func(t *testing.T, sec int64) {
		// The purpose of this function is forward looking. There's no reasonable
		// way for it to ever generate a time earlier than the day the code was
		// written, so we just don't *deal* with negative unix times, because we
		// just don't need to represent unix times times prior to 1970.
		if sec <= 0 || sec > postEarth {
			return
		}

		start := time.Unix(sec, 0)
		end := start.Add(BLOCK_SIZE * period)

		first := newUuid(start)
		next := newUuid(start.Add(period))
		prev := newUuid(end.Add(-1 * period))
		last := newUuid(end)

		// We should wrap to the first hash of possible values.
		assert.Equal(t, first[:2], last[:2])

		pairs := [2]([2]uuid.UUID){{first, next}, {prev, last}}

		for _, pair := range pairs {
			older := asHashBucketNumber(pair[0])
			newer := asHashBucketNumber(pair[1])

			if newer != 0 {
				// Hashes should always increment every PERIOD seconds...
				assert.Greater(t, newer, older)
			} else {
				// And wrap over at size_uint16
				assert.Equal(t, uint16(BLOCK_SIZE-1), older)
			}
		}
	})
}

func Benchmark_timeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}

func Benchmark_NewUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUuid()
	}
}

func Benchmark_googleUuid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		uuid.New()
	}
}
