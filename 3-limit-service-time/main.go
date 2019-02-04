//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import "time"
import "sync"
//import "fmt"

var maxTime int64
var users chan *User

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool

	TimeUsed  int64 // in seconds
	mux       sync.Mutex // for accessing the shared variable TimeUsed
}

func (u *User) IncrementTimeUsed() int64 {
	u.mux.Lock()
	defer u.mux.Unlock()

	u.TimeUsed++
	return u.TimeUsed
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		// Premium users have unlimited access, so no need
		// to keep track of their time.
		process()
		return true
	}
	// We have a free user.

	// Possible data-race, but not a big deal because
	// TimeUsed is always incremented + the writes are
	// thread-safe
	if u.TimeUsed >= maxTime {
		return false
	}

	timesUp := make(chan bool)
	processDone := make(chan bool)

	go func() {
		processStarted := make(chan bool)
		go func() {
			processStarted <- true
			process()

			// close to broadcast a signal
			// to both the ticking goroutine
			// and the handler's goroutine
			close(processDone)
		}()
		<-processStarted

		for {
			select {
			case <-processDone:
				return
			default:
				time.Sleep(1 * time.Second)
				timeUsed := u.IncrementTimeUsed()
				if timeUsed >= maxTime {
					timesUp <- true
					return
				}
			}
		}
	}()

	for {
		select {
		case <-timesUp:
			return false
		case <-processDone:
			return true
		}
	}
}

func main() {
	maxTime = 16
	RunMockServer()
}
