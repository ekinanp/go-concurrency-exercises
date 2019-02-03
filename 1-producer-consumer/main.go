//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
	"sync"
)

var streamMux sync.Mutex

func producer(stream *Stream, tweets chan<- *Tweet) {
	for {
		streamMux.Lock()
		tweet, err := stream.Next()
		streamMux.Unlock()

		if err == ErrEOF {
			return
		}

		tweets <- tweet
	}
}

func consumer(tweets <-chan *Tweet) {
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	num_producers := 1
	num_consumers := 1

	// Make the channel
	tweets := make(chan *Tweet, 4)

	// Spawn the producers
	var producers sync.WaitGroup
	for i := 0; i < num_producers; i++ {
		producers.Add(1)
		go func() {
			producer(&stream, tweets)
			producers.Done()
		}()
	}

	// Spawn the consumers
	var consumers sync.WaitGroup
	for i := 0; i < num_consumers; i++ {
		consumers.Add(1)
		go func() {
			consumer(tweets)
			consumers.Done()
		}()
	}

	// Wait for producers to finish, then close channel
	producers.Wait()
	close(tweets)

	// Wait for the consumers
	consumers.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
