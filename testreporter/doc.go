// Package testreporter contains a Reporter for collecting detailed testutil reports.
//
// While you could probably get at all of the exposed information in reports
// by examining the output of "go testutil",
// the generated report is intended to collect the information in one machine-readable file.
//
// Proper use of the reporter requires some slight changes to how you normally write tests.
// If you do any of these "the normal way", the tests will still operate fine;
// you will just miss some detail in the external report.
//
// First, the reporter instance must be initialized and Closed.
// The cmd/ibctest package does this in a MainTest function, similar to this:
//
//	func TestMain(m *testing.M) {
//	  f, _ := os.Create("/tmp/report.json")
//	  reporter := testreporter.NewReporter(f)
//	  code := m.Run()
//	  _ = reporter.Close()
//	  os.Exit(code)
//	}
//
// Next, every testutil that needs to be tracked must call TrackTest.
// If you omit the call to TrackTest, then the testutil's start and end time,
// and skip/fail status, will not be reported.
//
//	var reporter *testreporter.Reporter // Initialized somehow.
//
//	func TestFoo(t *testing.T) {
//	  reporter.TrackTest(t)
//	  // Normal testutil usage continues...
//	}
//
// Calling TrackTest tracks the testutil's start and finish time,
// including whether the testutil was skipped or failed.
//
// Parallel tests should not call t.Parallel directly,
// but instead should use TrackParallel.
// This will track the time the testutil paused waiting for parallel execution
// and when parallel execution resumes.
// If you omit the call to TrackParallel, then at worst you have a misleading testutil duration.
//
//	func TestFooParallel(t *testing.T) {
//	  reporter.TrackTest(t)
//	  reporter.TrackParallel(t)
//	  // Normal testutil usage continues...
//	}
//
// If a testutil needs to be skipped, the TrackSkip method will track the skip reason.
// Like the other Track methods, calling t.Skip directly will still cause the testutil to be skipped,
// and the reporter will note that the testutil was skipped,
// but the reporter would not track the specific skip reason.
//
//	func TestFooSkip(t *testing.T) {
//	  if someReason() {
//	    reporter.TrackSkip(t, "skipping due to %s", whySkipped())
//	  }
//	}
//
// Lastly, and perhaps most importantly, the reporter is designed to integrate
// with testify's require and assert packages.
// Plain "go testutil" runs simply have a stream of log lines and a failure/skip state.
// But if you connect the reporter with a require or assert instance,
// any failed assertions are stored as error messages in the report.
//
//	func TestBar(t *testing.T) {
//	  reporter.TrackTest(t)
//	  req := require.New(reporter.TestifyT(t))
//	  t.Log("About to testutil Bar()") // Goes to "go testutil" output, but not included in report.
//
//	  // If this fails, the report includes a "TestErrorMessage" entry in the report.
//	  req.NoError(Bar(), "failure executing Bar()")
//	}
//
// If you use a plain require.NoError(t, err) call,
// the report will note that the testutil failed, but the report will not include the error line.
package testreporter
