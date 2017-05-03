package awg

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"
)

var count int64

func slowFunc() error {
	for i := 0; i < 10000000; i++ {
		count *= int64(i)
	}

	return nil
}

func fastFunc() error {
	// do nothing
	return nil
}

func errorFunc() error {
	for i := 0; i < 10000; i++ {
		count *= int64(i)
	}

	return errors.New("Test error")
}

func panicFunc() error {
	panic("Test panic")
}

// TestAdvancedWorkGroupTimeout test for timeout
func Test_AdvancedWorkGroupTimeout(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)

	wg.SetTimeout(time.Nanosecond * 10).Start()

	if wg.Status() != StatusTimeout {
		t.Error("AWG should stops by timeout!")
	}

	err := wg.GetLastError()
	if err == nil {
		t.Error("AWG should stops with error by timeout!")
	}

	if _, ok := err.(ErrorTimeout); !ok {
		t.Errorf("Wrong error type. Got %[1]T: %[1]q", err)
	}
}

// Test_AdvancedWorkGroupTimeout_Context test for timeout
func Test_AdvancedWorkGroupTimeout_Context(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Nanosecond))
	defer cancel()

	wg.WithContext(ctx).Start()
	if wg.Status() != StatusTimeout {
		t.Error("AWG should stops by timeout!")
	}

	err := wg.GetLastError()
	if err == nil {
		t.Error("AWG should stops with error by timeout!")
	}

	if _, ok := err.(ErrorTimeout); !ok {
		t.Errorf("Wrong error type. Got %[1]T: %[1]q", err)
	}
}

// TestAdvancedWorkGroupError test for error
func Test_AdvancedWorkGroupError(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(errorFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)

	wg.SetStopOnError(true).Start()

	if wg.Status() != StatusError {
		t.Error("AWG should stops by error!")
	}

}

// TestAdvancedWorkGroupSuccess test for success case
func Test_AdvancedWorkGroupSuccess(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)

	if errs := wg.SetStopOnError(true).Start().GetAllErrors(); len(errs) != 0 {
		t.Errorf("AWG result should be 'success'! But got errors %v", errs)
	}

	if wg.Status() != StatusSuccess {
		t.Error("AWG result should be 'success'!")
	}
}

// Test_AdvancedWorkGroupSuccess_WithCapacity test for success case with capacity
func Test_AdvancedWorkGroupSuccess_WithCapacity(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)
	wg.SetCapacity(2)

	if errs := wg.SetStopOnError(true).Start().GetAllErrors(); len(errs) != 0 {
		t.Errorf("AWG result should be 'success'! But got errors %v", errs)
	}

	if wg.Status() != StatusSuccess {
		t.Error("AWG result should be 'success'!")
	}
}

// Test_AdvancedWorkGroupCancel_Success test for cancel case
func Test_AdvancedWorkGroupCancel_Success(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(5 * time.Microsecond)
		cancel()
	}()

	if errs := wg.WithContext(ctx).SetStopOnError(true).Start().GetAllErrors(); len(errs) != 0 {
		t.Errorf("AWG result should be 'success'! But got errors %v", errs)
	}

	if wg.Status() != StatusSuccess {
		t.Error("AWG result should be 'success'!")
	}
}

// Test_AdvancedWorkGroupCancelWithCapacity_Success test for cancel case with capacity
func Test_AdvancedWorkGroupCancelWithCapacity_Success(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(fastFunc)
	wg.Add(slowFunc)
	wg.Add(slowFunc)
	wg.SetCapacity(2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(5 * time.Microsecond)
		cancel()
	}()

	if errs := wg.WithContext(ctx).SetStopOnError(true).Start().GetAllErrors(); len(errs) != 0 {
		t.Errorf("AWG result should be 'success'! But got errors %v", errs)
	}

	if wg.Status() != StatusSuccess {
		t.Error("AWG result should be 'success'!")
	}
}

// TestAdvancedWorkGroupPanicError test for success case
func Test_AdvancedWorkGroupPanicError(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(slowFunc)
	wg.Add(panicFunc)

	wg.SetStopOnError(true).Start()

	if wg.Status() != StatusError {
		t.Error("AWG result should be 'error'!")
	}

	t.Log(wg.GetLastError())
}

// TestAdvancedWorkGroupPanicSuccess test for success case
func TestAdvancedWorkGroupPanicSuccess(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(slowFunc)
	wg.Add(panicFunc)

	err := wg.SetStopOnError(false).Start().GetLastError()

	if err == nil {
		t.Error("Panic should be an error")
	}
}

// Test_AdvancedWorkGroupStopOnErrorPanic test for success case
func Test_AdvancedWorkGroupStopOnErrorPanic(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(slowFunc)
	wg.Add(panicFunc)

	wg.SetStopOnError(true).
		Start()

	if wg.Status() != StatusError {
		t.Error("AWG result should be 'error'!")
	}
}

// TestAWGStopOnError tests AdvancedWaitGroup with StopOnError set to false
// and with failing task.
func TestAWGStopOnError(t *testing.T) {
	var wg AdvancedWaitGroup
	wg.Add(fastFunc)
	wg.Add(errorFunc)
	wg.Add(fastFunc)
	wg.SetStopOnError(false).
		Start()
}

// Test_AdvancedWorkGroup_NoLeak tests for goroutines leaks
func Test_AdvancedWorkGroup_NoLeak(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.Add(errorFunc)

	wg.SetStopOnError(true).
		Start()

	time.Sleep(2 * time.Second)

	numGoroutines := runtime.NumGoroutine()

	var wg2 AdvancedWaitGroup

	wg2.Add(errorFunc)
	wg2.Add(slowFunc)

	wg2.SetStopOnError(true).
		Start()

	time.Sleep(3 * time.Second)

	numGoroutines2 := runtime.NumGoroutine()

	if numGoroutines != numGoroutines2 {
		t.Fatalf("We leaked %d goroutine(s)", numGoroutines2-numGoroutines)
	}
}

// Test_AdvancedWorkGroupAddSliceTimeout test for timeout
func Test_AdvancedWorkGroupAddSliceTimeout(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.AddSlice([]WaitgroupFunc{fastFunc, fastFunc, fastFunc, slowFunc, slowFunc, slowFunc})
	wg.SetTimeout(time.Nanosecond * 10).SetStopOnError(true).Start()

	if wg.Status() != StatusTimeout {
		t.Error("AWG should stops by timeout!", wg.Status())
	}

	err := wg.GetLastError()
	if err == nil {
		t.Error("AWG should stops with error by timeout!")
	}

	if _, ok := err.(ErrorTimeout); !ok {
		t.Errorf("Wrong error type. Got %[1]T: %[1]q", err)
	}
}

// Test_AdvancedWorkGroupAddSliceTimeoutSuccess test for timeout that if too long
func Test_AdvancedWorkGroupAddSliceTimeoutSuccess(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.AddSlice([]WaitgroupFunc{fastFunc, fastFunc, fastFunc})
	wg.SetTimeout(time.Second * 10).SetStopOnError(true).Start()

	if wg.Status() != StatusSuccess {
		t.Error("AWG shouldn`t stops by timeout!", wg.Status())
	}

	err := wg.GetLastError()
	if err != nil {
		t.Error("AWG shouldn`t stops with error by timeout!")
	}
}

// Test_AdvancedWorkGroupGetLastError test
func Test_AdvancedWorkGroupGetLastErrorSuccess(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.AddSlice([]WaitgroupFunc{fastFunc})
	wg.Start()

	err := wg.GetLastError()
	if err != nil {
		t.Error("Shouldn`t get errors!")
	}
}

// Test_AdvancedWorkGroupGetLastErrorSingleError test
func Test_AdvancedWorkGroupGetLastErrorSingleError(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.AddSlice([]WaitgroupFunc{fastFunc, errorFunc})
	wg.Start()

	err := wg.GetLastError()
	if err == nil {
		t.Error("Should get error!")
	}
}

// Test_AdvancedWorkGroupGetAllErrorsSuccess test
func Test_AdvancedWorkGroupGetAllErrorsSuccess(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.AddSlice([]WaitgroupFunc{fastFunc})
	wg.Start()

	errs := wg.GetAllErrors()
	if len(errs) != 0 {
		t.Error("Shouldn`t get errors!")
	}
}

// Test_AdvancedWorkGroupGetAllErrorsManyErrors test
func Test_AdvancedWorkGroupGetAllErrorsManyErrors(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.AddSlice([]WaitgroupFunc{fastFunc, errorFunc, errorFunc})
	wg.Start()

	errs := wg.GetAllErrors()
	if len(errs) != 2 {
		t.Error("Should get errors")
	}
}

// Test_AdvancedWorkGroupGetAllErrorsSingleError test
func Test_AdvancedWorkGroupGetAllErrorsSingleError(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.AddSlice([]WaitgroupFunc{fastFunc, errorFunc})
	wg.Start()

	errs := wg.GetAllErrors()
	if len(errs) != 1 {
		t.Error("Should get one error!")
	}
}

// Test_AdvancedWorkGroupReset test
func Test_AdvancedWorkGroupReset(t *testing.T) {
	var wg AdvancedWaitGroup

	wg.AddSlice([]WaitgroupFunc{fastFunc, errorFunc})
	wg.Start()

	wg.Reset()
	if wg.Status() != StatusIdle {
		t.Error("Cleaned wg should have idle status")
	}

	wg.Add(errorFunc, errorFunc)

	if errs := wg.Start().GetAllErrors(); len(errs) != 2 {
		t.Error("Should get two errors on cleaned wg")
	}
}

// Test_AdvancedWorkGroupDoubleStart test
func Test_AdvancedWorkGroupDoubleStart(t *testing.T) {
	var wg1, wg2 AdvancedWaitGroup

	wg1.AddSlice([]WaitgroupFunc{slowFunc, errorFunc, fastFunc, errorFunc})
	chDone := make(chan struct{})
	go func() {
		wg1.Start()
		chDone <- struct{}{}
	}()

	<-chDone
	if errs := wg1.Start().GetAllErrors(); len(errs) != 2 {
		t.Error("Should get two errors on wg")
	}

	chDone2 := make(chan struct{})
	wg2.AddSlice([]WaitgroupFunc{slowFunc, errorFunc, fastFunc, errorFunc})
	go func() {
		<-chDone2
		wg2.Start()
	}()

	wg2.Start()
	chDone2 <- struct{}{}
	if errs := wg2.GetAllErrors(); len(errs) != 2 {
		t.Error("Should get two errors on wg")
	}
}

var results = make(chan bool, 100)

func fastFuncWithResult() error {
	results <- true
	return nil
}

// TestAdvancedWorkGroupTimeout test for timeout
func Test_AdvancedWorkGroupTimeout_Execution(t *testing.T) {
	var wg AdvancedWaitGroup
	maxProcs := 8 * runtime.NumCPU()

	for i := 0; i < maxProcs; i++ {
		wg.Add(errorFunc)
	}

	for i := 0; i < maxProcs; i++ {
		wg.Add(fastFuncWithResult)
	}

	wg.SetTimeout(time.Microsecond).SetStopOnError(true).Start()

	time.Sleep(time.Second)
	close(results)

	err := wg.GetLastError()
	if err == nil {
		t.Error("AWG should stops with error")
	}

	if errs := wg.GetAllErrors(); len(errs) > 0 {
		for _, err := range errs {
			t.Log(err)
		}
	}

	count := 0
	for range results {
		count++
	}

	if count >= maxProcs {
		t.Errorf("Some wg functions should be interrupted")
	}

	//Debug
	t.Logf("Done %v of %v", count, maxProcs)
}
