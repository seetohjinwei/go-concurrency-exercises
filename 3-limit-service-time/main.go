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

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

var mu sync.Mutex

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	check := make(chan struct{})

	if u.IsPremium {
		process()
		return true
	}

	start := time.Now()

	go func() {
		process()
		close(check)
	}()

	// mutex for modifying u.TimeUsed
	mu.Lock()
	defer mu.Unlock()

	select {
	case <-check:
		end := time.Now()
		elapsed := end.Sub(start)
		u.TimeUsed += int64(elapsed)
		return u.TimeUsed <= int64(10*time.Second)
	case <-time.After(10*time.Second - time.Duration(u.TimeUsed)):
		return false
	}
}

func main() {
	RunMockServer()
}
