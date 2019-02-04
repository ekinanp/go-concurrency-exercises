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

import "os"
import "os/signal"
import "syscall"

func main() {
	// Create a process
	proc := MockProcess{}

	sigCh := make(chan os.Signal)
	signal.Notify(
		sigCh,
		syscall.SIGINT,
	)

	procStarted := make(chan bool)

	// This routine will poll for the SIGINT signal
	go func() {
		// This avoids erroneous SIGINT signals
		<-procStarted

		// Wait on the SIGINT signal
		<-sigCh
		signal.Stop(sigCh)

		proc.Stop()
	}()

	// Run the process (blocking)
	procStarted <- true
	proc.Run()
}
