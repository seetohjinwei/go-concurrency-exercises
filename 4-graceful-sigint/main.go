//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"os"
	"os/signal"
)

func main() {
	// Create a process
	proc := MockProcess{}

	// Help from: https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-sigint-and-run-a-cleanup-function-i
	sigint := make(chan os.Signal)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		<-sigint       // first SIGINT
		go proc.Stop() // run this in a separate goroutine so that it won't block this goroutine

		<-sigint // second SIGINT
		close(sigint)
		os.Exit(1)
	}()

	// Run the process (blocking)
	proc.Run()
}
