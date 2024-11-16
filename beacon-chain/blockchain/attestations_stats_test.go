package blockchain

import (
	"errors"
	"testing"

	"github.com/prysmaticlabs/prysm/v5/testing/assert"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
)

func TestAddSuccess(t *testing.T) {
	a := NewAttestationsStats()

	a.AddSuccess()
	a.AddSuccess()

	assert.Equal(t, 2, a.successfulCount)
}

func TestFailure(t *testing.T) {
	a := NewAttestationsStats()

	a.AddFailure(FailedAttestation{})
	a.AddFailure(FailedAttestation{})

	assert.Equal(t, 2, a.failedCount)
	assert.Equal(t, 2, len(a.failedAttestations))
}

func TestReportEpochTransition(t *testing.T) {
	hook := logtest.NewGlobal()

	a := NewAttestationsStats()

	a.AddSuccess()
	a.AddSuccess()

	a.AddFailure(FailedAttestation{
		Fields: logrus.Fields{
			"slot":             "1",
			"committeeCount":   "1",
			"committeeIndices": "1",
			"beaconBlockRoot":  "1",
			"targetRoot":       "1",
			"aggregatedCount":  "1",
		},
		Error: errors.New("failed to attestation"),
	})
	a.AddFailure(FailedAttestation{
		Fields: logrus.Fields{
			"slot":             "2",
			"committeeCount":   "2",
			"committeeIndices": "2",
			"beaconBlockRoot":  "2",
			"targetRoot":       "2",
			"aggregatedCount":  "2",
		},
		Error: errors.New("failed to attestation"),
	})

	a.ReportEpochTransition(1)

	assert.LogsContain(t, hook,
		`"Attestation Summary for Epoch 1:" prefix=blockchain`,
	)
	assert.LogsContain(t, hook,
		`"  Successful Attestations Count: 2" prefix=blockchain`,
	)
	assert.LogsContain(t, hook,
		`"  Failed Attestations Count: 2" prefix=blockchain`,
	)
	assert.LogsContain(t, hook,
		`"  error: failed to attestation" aggregatedCount=1 beaconBlockRoot=1 committeeCount=1 committeeIndices=1 prefix=blockchain slot=1 targetRoot=1`,
	)
	assert.LogsContain(t, hook,
		`"  error: failed to attestation" aggregatedCount=2 beaconBlockRoot=2 committeeCount=2 committeeIndices=2 prefix=blockchain slot=2 targetRoot=2`,
	)

	hook = logtest.NewGlobal()
	a.AddSuccess()
	a.AddSuccess()
	a.AddSuccess()
	a.AddSuccess()
	a.AddSuccess()
	a.AddSuccess()

	a.ReportEpochTransition(2)
	assert.LogsContain(t, hook,
		`"Attestation Summary for Epoch 2:" prefix=blockchain`,
	)
	assert.LogsContain(t, hook,
		`"  Successful Attestations Count: 6" prefix=blockchain`,
	)
	assert.LogsContain(t, hook,
		`"  Failed Attestations Count: 0" prefix=blockchain`,
	)
}
