// Exercise is from
// https://go.dev/tour/concurrency/8

package main

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *Tree, ch chan<- int) {
	if t == nil {
		return
	}

	// these have to be blocking, otherwise, order will be incorrect!
	Walk(t.Left, ch)
	ch <- t.Value
	Walk(t.Right, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	startWalking := func(t *Tree, ch chan<- int) {
		Walk(t, ch)
		close(ch) // closing channels here is simpler than using WaitGroups
	}

	go startWalking(t1, ch1)
	go startWalking(t2, ch2)

	for {
		a, ok1 := <-ch1
		b, ok2 := <-ch2

		if !ok1 || !ok2 {
			// either channel is closed
			break
		}

		if a != b {
			return false
		}
	}

	return true
}
