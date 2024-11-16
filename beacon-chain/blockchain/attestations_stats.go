package blockchain

import (
	"sync"

	"github.com/sirupsen/logrus"
)

// AttestationsStats represents the stat for every epoch of successful and failed attestation counts
// It records on receiving Attestation in this package file: receive_attestation.go - function: processAttestations
// It prints out the stat for every epoch and stat is reset.
type AttestationsStats struct {
	mu                 sync.Mutex
	successfulCount    int
	failedCount        int
	failedAttestations []FailedAttestation
}

// FailedAttestation represents the failed reason for the attestation
// It is used to debug the attestation of this node operation
type FailedAttestation struct {
	Fields logrus.Fields
	Error  error
}

// NewAttestationsStats instantiates a new AttestationsStats
func NewAttestationsStats() *AttestationsStats {
	return &AttestationsStats{
		mu:                 sync.Mutex{},
		successfulCount:    0,
		failedCount:        0,
		failedAttestations: []FailedAttestation{},
	}
}

// AddSuccess add successfulCount to AttestationsStats
func (a *AttestationsStats) AddSuccess() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.successfulCount += 1
}

// AddFailure add failedCount and append FailedAttestation to AttestationsStats
func (a *AttestationsStats) AddFailure(f FailedAttestation) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.failedCount += 1
	a.failedAttestations = append(a.failedAttestations, f)
}

// ReportEpochTransition generate report of AttestationsStats to log
// It also resets the stats
func (a *AttestationsStats) ReportEpochTransition(currentEpoch uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	log.Infof("Attestation Summary for Epoch %d:", currentEpoch)
	log.Infof("  Successful Attestations Count: %d", a.successfulCount)
	log.Infof("  Failed Attestations Count: %d", a.failedCount)
	for _, v := range a.failedAttestations {
		log.WithFields(v.Fields).Infof("  error: %s", v.Error.Error())
	}

	// reset stats
	a.successfulCount = 0
	a.failedCount = 0
	a.failedAttestations = []FailedAttestation{}
}
