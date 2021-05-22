package topshot

import (
	"fmt"
)

// return is sleep ms
type OnErrorFunc func(*BackoffControl, error) (int, error)

// *BackoffControl ->
// bool -> true if error has just been promoted to the current level
// error -> The error which just occurred
type OnLogFunc func(*BackoffControl, bool, error)

type BackoffLogger struct {
	initialDetection  OnLogFunc
	warningDetection  OnLogFunc
	seriousDetection  OnLogFunc
	alertDetection    OnLogFunc
	allClearDetection func()
}

type BackoffAlgorithm struct {
	onError OnErrorFunc
	logger  *BackoffLogger
}

type BackoffControl struct {
	currentError      error
	currentErrorCount int

	lastError   error
	totalErrors int

	errorHandler *BackoffAlgorithm
}

func Backoff_ConsoleLogger() *BackoffLogger {
	return &BackoffLogger{
		initialDetection: func(bc *BackoffControl, promoted bool, e error) {
			fmt.Println("BackoffLogger initial error detected(", bc.currentErrorCount, ")", e)
		},
		warningDetection: func(bc *BackoffControl, promoted bool, e error) {
			fmt.Println("BackoffLogger WARNING error detected(", bc.currentErrorCount, ")", e)
		},
		seriousDetection: func(bc *BackoffControl, promoted bool, e error) {
			fmt.Println("BackoffLogger SERIOUS error detected(", bc.currentErrorCount, ")", e)
		},
		alertDetection: func(bc *BackoffControl, promoted bool, e error) {
			fmt.Println("BackoffLogger ALERT error detected(", bc.currentErrorCount, ")", e)
		},
		allClearDetection: func() {
			fmt.Println("BackoffLogger working again")
		},
	}
}

func Backoff_AlgorithmNoBlast(logger *BackoffLogger) *BackoffAlgorithm {
	const WARNING_LEVEL = 10
	const SERIOUS_LEVEL = 40
	const ALERT_LEVEL = 70

	return &BackoffAlgorithm{
		logger: logger,
		onError: func(control *BackoffControl, errEvent error) (int, error) {
			if control.currentErrorCount == 1 {
				logger.initialDetection(control, true, errEvent)
				return 10, nil
			} else if control.currentErrorCount < WARNING_LEVEL {
				logger.initialDetection(control, false, errEvent)
				return 1000, nil
			} else if control.currentErrorCount == WARNING_LEVEL {
				logger.warningDetection(control, true, errEvent)
				return 60 * 1000, nil
			} else if control.currentErrorCount < SERIOUS_LEVEL {
				logger.warningDetection(control, false, errEvent)
				return 60 * 3 * 1000, nil
			} else if control.currentErrorCount == SERIOUS_LEVEL {
				logger.seriousDetection(control, true, errEvent)
				return 60 * 3 * 1000, nil
			} else if control.currentErrorCount < ALERT_LEVEL {
				logger.seriousDetection(control, false, errEvent)
				return 60 * 3 * 1000, nil
			} else if control.currentErrorCount == ALERT_LEVEL {
				logger.alertDetection(control, true, errEvent)
				return 15 * 60 * 1000, nil
			} else {
				logger.alertDetection(control, false, errEvent)
				return 15 * 60 * 1000, nil
			}
		},
	}
}

func Backoff_Create(errorHandler *BackoffAlgorithm) *BackoffControl {
	return &BackoffControl{
		currentError:      nil,
		currentErrorCount: 0,
		lastError:         nil,
		totalErrors:       0,
		errorHandler:      errorHandler,
	}
}

func Backoff_HandleReturn(control *BackoffControl, err error) (int, error) {
	if err != nil {
		if control.currentErrorCount == 0 {
			control.lastError = err
			control.totalErrors += 1
		}
		control.currentError = err
		control.currentErrorCount += 1
		return control.errorHandler.onError(control, err)
	} else {
		if control.currentErrorCount > 0 {
			control.currentError = nil
			control.currentErrorCount = 0
			control.errorHandler.logger.allClearDetection()
			// Could report offline time?
		}
		return 0, nil
	}
}
