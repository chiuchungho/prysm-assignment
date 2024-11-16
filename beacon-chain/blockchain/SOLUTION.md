# Task
Prysm is not verifying attestations correctly. Add logging for the attestations.

- Counts successfully verified attestations.
- Counts failed attestations and records the reason for each failure.
- Outputs a summary of the collected data at the end of each epoch.

## Solution
My idea is to record the attestion in the `beacon-chain/blockchain`. This directory has the service that handles the internal logic of managing the full PoS beacon chain. It will start logging the attestation count of the incoming new slot right after the deployment.

1. Added `beacon-chain/blockchain/attestations_stats.go`
- Create new struct AttestationsStats
- In-memory storage to collect the successful count, failed count and failed resason.
- It can also generate log report with function - ReportEpochTransition
2. Modified `beacon-chain/blockchain/receive_attestation.go`
- add `attestationsStats.AddFailure()` and `attestationsStats.AddSuccess()` in func `processAttestations`
- It is recording the attestation count on receive attestation
3. Modified `beacon-chain/blockchain/receive_block.go`
- added `attestationsStats.ReportEpochTransition()` in `updateCheckpoints`
- it is generating a report for the attestion count for every epoch
4. Modified `beacon-chain/blockchain/service.go`
- added attestationsStats to Service struct
- add initialation of attestationsStats to `NewService`

## Example Log Output
```
time="2024-11-16T02:21:53+01:00" level=info msg="Attestation Summary for Epoch 1:" prefix=blockchain
time="2024-11-16T02:21:53+01:00" level=info msg="  Successful Attestations Count: 2" prefix=blockchain
time="2024-11-16T02:21:53+01:00" level=info msg="  Failed Attestations Count: 2" prefix=blockchain
time="2024-11-16T02:21:53+01:00" level=info msg="  error: failed to attestation" aggregatedCount=1 beaconBlockRoot=1 committeeCount=1 committeeIndices=1 prefix=blockchain slot=1 targetRoot=1
time="2024-11-16T02:21:53+01:00" level=info msg="  error: failed to attestation" aggregatedCount=2 beaconBlockRoot=2 committeeCount=2 committeeIndices=2 prefix=blockchain slot=2 targetRoot=2
time="2024-11-16T02:21:53+01:00" level=info msg="Attestation Summary for Epoch 2:" prefix=blockchain
time="2024-11-16T02:21:53+01:00" level=info msg="  Successful Attestations Count: 6" prefix=blockchain
time="2024-11-16T02:21:53+01:00" level=info msg="  Failed Attestations Count: 0" prefix=blockchain
```

## Test
Unit test: `beacon-chain/blockchain/attestations_stats_test.go`

Todo:
- add unit test in beacon-chain/blockchain/service_test.go to fully test the logging for attestation
- add integration to test if the node can log attestion properly

## Improvement
- make the log report in a prettier format
- it could be tough to gather the attestation from the log, it could also be collected in the local directry of the node as a separated file to append the counts.
- the logging could be also added to `beacon-chain/sync`, if we want to have the attestion logging from a new beacon node from genesis block.