/**
 * author: NoChopFoundation@gmail.com
 */
package topshot

import (
	"errors"
	"testing"
)

type TestLoggerVerify struct {
	initialDetectionCount  int
	warningDetectionCount  int
	seriousDetectionCount  int
	alertDetectionCount    int
	allClearDetectionCount int
}

func TestConsoleLogger(t *testing.T) {
	control := Backoff_Create(Backoff_AlgorithmNoBlast(Backoff_ConsoleLogger()))

	Backoff_HandleReturn(control, nil)
	for i := 0; i < 200; i++ {
		Backoff_HandleReturn(control, errors.New("Hello"))
	}
	Backoff_HandleReturn(control, nil)
}

func TestBasicSetAndClear(t *testing.T) {
	actuals := TestLoggerVerify{}
	testLogger := BackoffLogger{
		initialDetection:  func(bc *BackoffControl, promoted bool, e error) { actuals.initialDetectionCount++ },
		warningDetection:  func(bc *BackoffControl, promoted bool, e error) { actuals.warningDetectionCount++ },
		seriousDetection:  func(bc *BackoffControl, promoted bool, e error) { actuals.seriousDetectionCount++ },
		alertDetection:    func(bc *BackoffControl, promoted bool, e error) { actuals.alertDetectionCount++ },
		allClearDetection: func() { actuals.allClearDetectionCount++ },
	}

	control := Backoff_Create(Backoff_AlgorithmNoBlast(&testLogger))

	Backoff_HandleReturn(control, nil)
	assetLoggerCounts(t, &actuals, 0, 0, 0, 0, 0)

	Backoff_HandleReturn(control, errors.New("Hello"))
	assetLoggerCounts(t, &actuals, 1, 0, 0, 0, 0)
	assetInt(t, control.currentErrorCount, 1)

	Backoff_HandleReturn(control, errors.New("Hello"))
	assetLoggerCounts(t, &actuals, 1, 0, 0, 0, 0)
	assetInt(t, control.currentErrorCount, 2)

	Backoff_HandleReturn(control, nil)
	assetLoggerCounts(t, &actuals, 1, 0, 0, 0, 1)

	assetInt(t, control.currentErrorCount, 0)
	assetInt(t, control.totalErrors, 1)

	//TODO more test as we develop how we want it to work
}

func assetInt(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("expect %d actual %d", expected, actual)
	}
}

func assetLoggerCounts(t *testing.T, actual *TestLoggerVerify, initialDetectionCount int, warningDetectionCount int,
	seriousDetectionCount int, alertDetectionCount int, allClearDetectionCount int) {
	if actual.initialDetectionCount != initialDetectionCount || actual.warningDetectionCount != warningDetectionCount ||
		actual.seriousDetectionCount != seriousDetectionCount || actual.alertDetectionCount != alertDetectionCount ||
		actual.allClearDetectionCount != allClearDetectionCount {
		t.Error("assertLoggerCounts")
	}
}
