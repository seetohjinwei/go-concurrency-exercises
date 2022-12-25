// Exercise from:
// https://go.dev/tour/concurrency/10

package main

import (
	"fmt"
	"sync"
)

type SafeSet struct {
	mu sync.Mutex
	v  map[string]bool
}

func NewSafeSet() *SafeSet {
	return &SafeSet{
		v: make(map[string]bool),
	}
}

func (c *SafeSet) Add(url string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.v[url] = true
}

func (c *SafeSet) IsKey(url string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, inMap := c.v[url]
	return inMap
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup, set *SafeSet) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:

	defer wg.Done()

	if depth <= 0 {
		return
	}

	if set.IsKey(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	set.Add(url)
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		wg.Add(1)
		go Crawl(u, depth-1, fetcher, wg, set)
	}
	return
}

func main() {
	var wg sync.WaitGroup
	set := NewSafeSet()

	wg.Add(1)
	go Crawl("https://golang.org/", 4, fetcher, &wg, set)
	wg.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
