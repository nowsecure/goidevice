package installation

// #cgo pkg-config: libimobiledevice-1.0
// #include <stdlib.h>
// #include <libimobiledevice/installation_proxy.h>
import "C"
import (
	"errors"
)

var (
	// ErrInvalidArgs .
	ErrInvalidArgs = errors.New("invalid args")
	// ErrUnknown .
	ErrUnknown = errors.New("unknown")
	// ErrNoDevice .
	ErrNoDevice = errors.New("no device")
	// ErrNotEnoughData .
	ErrNotEnoughData = errors.New("not enough data")
	// ErrBadHeader .
	ErrBadHeader = errors.New("bad header")
	// ErrSslError .
	ErrSslError = errors.New("ssl error")
	// ErrReceiveTimeout .
	ErrReceiveTimeout = errors.New("receive timeout")
	// ErrMuxError .
	ErrMuxError = errors.New("mux error")
	// ErrNoRunningSession .
	ErrNoRunningSession = errors.New("no running session")
	// ErrInvalidResponse .
	ErrInvalidResponse = errors.New("invalid response")
	// ErrMissingKey .
	ErrMissingKey = errors.New("missing key")
	// ErrMissingValue .
	ErrMissingValue = errors.New("missing value")
	// ErrGetProhibited .
	ErrGetProhibited = errors.New("get prohibited")
	// ErrSetProhibited .
	ErrSetProhibited = errors.New("set prohibited")
	// ErrRemoveProhibited .
	ErrRemoveProhibited = errors.New("remove prohibited")
	// ErrImmutableValue .
	ErrImmutableValue = errors.New("immutable value")
	// ErrPasswordProtected .
	ErrPasswordProtected = errors.New("password protected")
	// ErrUserDeniedPairing .
	ErrUserDeniedPairing = errors.New("user denied pairing")
	// ErrPairingDialogResponsePending .
	ErrPairingDialogResponsePending = errors.New("pairing dialog response pending")
	// ErrMissingHostID .
	ErrMissingHostID = errors.New("missing host id")
	// ErrInvalidHostID .
	ErrInvalidHostID = errors.New("invalid host id")
	// ErrSessionActive .
	ErrSessionActive = errors.New("session active")
	// ErrSessionInactive .
	ErrSessionInactive = errors.New("session inactive")
	// ErrMissingSessionID .
	ErrMissingSessionID = errors.New("missing session id")
	// ErrInvalidSessionID .
	ErrInvalidSessionID = errors.New("invalid session id")
	// ErrMissingService .
	ErrMissingService = errors.New("missing service")
	// ErrInvalidService .
	ErrInvalidService = errors.New("invalid service")
	// ErrServiceLimit .
	ErrServiceLimit = errors.New("service limit")
	// ErrMissingPairRecord .
	ErrMissingPairRecord = errors.New("missing pair record")
	// ErrSavePairRecordFailed .
	ErrSavePairRecordFailed = errors.New("save pair record failed")
	// ErrInvalidPairRecord .
	ErrInvalidPairRecord = errors.New("invalid pair record")
	// ErrInvalidActivationPeriod .
	ErrInvalidActivationPeriod = errors.New("invalid activation period")
	// ErrMissingActivationPeriod .
	ErrMissingActivationPeriod = errors.New("missing activation period")
	// ErrServiceProhibited .
	ErrServiceProhibited = errors.New("service prohibited")
	// ErrEscrowLocked .
	ErrEscrowLocked = errors.New("escrow locked")
	// ErrPairingProhibitedOverThisConnection .
	ErrPairingProhibitedOverThisConnection = errors.New("pairing prohibited over this connection")
	// ErrFMIPProtected .
	ErrFMIPProtected = errors.New("FMIP protected")
	// ErrMCProtected .
	ErrMCProtected = errors.New("MC Protected")
	// ErrMCChallengeRequired .
	ErrMCChallengeRequired = errors.New("mc challenge required")
)

func resultToError(result C.instproxy_error_t) error {
	switch result {
	case 0:
		return nil
	case -1:
		return ErrInvalidArgs
	case -2:
		return ErrUnknown
	case -3:
		return ErrNoDevice
	case -4:
		return ErrNotEnoughData
	case -5:
		return ErrBadHeader
	case -6:
		return ErrSslError
	case -7:
		return ErrReceiveTimeout
	case -8:
		return ErrMuxError
	case -9:
		return ErrNoRunningSession
	case -10:
		return ErrInvalidResponse
	case -11:
		return ErrMissingKey
	case -12:
		return ErrMissingValue
	case -13:
		return ErrGetProhibited
	case -14:
		return ErrSetProhibited
	case -15:
		return ErrRemoveProhibited
	case -16:
		return ErrImmutableValue
	case -17:
		return ErrPasswordProtected
	case -18:
		return ErrUserDeniedPairing
	case -19:
		return ErrPairingDialogResponsePending
	case -20:
		return ErrMissingHostID
	case -21:
		return ErrInvalidHostID
	case -22:
		return ErrSessionActive
	case -23:
		return ErrSessionInactive
	case -24:
		return ErrMissingSessionID
	case -25:
		return ErrInvalidSessionID
	case -26:
		return ErrMissingService
	case -27:
		return ErrInvalidService
	case -28:
		return ErrServiceLimit
	case -29:
		return ErrMissingPairRecord
	case -30:
		return ErrSavePairRecordFailed
	case -31:
		return ErrInvalidPairRecord
	case -32:
		return ErrInvalidActivationPeriod
	case -33:
		return ErrMissingActivationPeriod
	case -34:
		return ErrServiceProhibited
	case -35:
		return ErrEscrowLocked
	case -36:
		return ErrPairingProhibitedOverThisConnection
	case -37:
		return ErrFMIPProtected
	case -38:
		return ErrMCProtected
	case -39:
		return ErrMCChallengeRequired
	default:
		return ErrUnknown
	}
}
