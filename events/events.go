package events

import (
	"context"
	"fmt"
	"time"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
)

// TODO: increase once implementation is somewhat complete
const (
	period      = 1 * time.Second
	chanTimeout = period / 2
)

type Listener struct {
	ticker *time.Ticker
	done   chan bool
	types  map[string]bool
	Events chan []flow.Event
	fc     *client.Client
}

type TimeOutError struct {
	error
}

func NewListener(fc *client.Client) *Listener {
	return &Listener{
		nil,
		make(chan bool),
		make(map[string]bool),
		nil,
		fc,
	}
}

func (l *Listener) process(ctx context.Context, start, end uint64) error {
	ee := make([]flow.Event, 0)

	for t := range l.types {
		r, err := l.fc.GetEventsForHeightRange(ctx, client.EventRangeQuery{
			Type:        t,
			StartHeight: start,
			EndHeight:   end,
		})
		if err != nil {
			return fmt.Errorf("error while fetching events: %w", err)
		}
		for _, b := range r {
			ee = append(ee, b.Events...)
		}
	}

	select {
	case l.Events <- ee:
		// Sent
		return nil
	case <-time.After(chanTimeout):
		// Timed out while waiting for channel
		return &TimeOutError{}
	}
}

func (l *Listener) handleError(err error) {
	_, ok := err.(*TimeOutError)
	if ok {
		// Ignore timeout errors
		return
	}
	fmt.Println(err)
}

func (l *Listener) AddType(t string) {
	l.types[t] = true
}

func (l *Listener) RemoveType(t string) {
	delete(l.types, t)
}

func (l *Listener) Start() *Listener {
	if l.ticker != nil {
		return l
	}

	l.ticker = time.NewTicker(period)
	l.Events = make(chan []flow.Event)

	go func() {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var start, end uint64
		cur, _ := l.fc.GetLatestBlock(ctx, true)
		start = cur.Height

		// TODO:

		// psiemens:
		// Rather than starting the listener from the latest block, I recommend
		// storing the "last fetched height" in the database. This way if the server
		// stops or falls behind the chain, it can catch up to the latest block next
		// time it starts up.
		//
		// However, if the listener falls behind, it may need to catch up on many
		// blocks at once (e.g. 1000+). The chain can only handle queries for
		// roughly 200 blocks at a time, so you'll need to batch these requests.
		//
		// You could do this by enforcing a maximum value on the block difference
		// (end - start) below. It probably makes to make this maximum a
		// configurable value.

		for {
			select {
			case <-l.done:
				return
			case <-l.ticker.C:
				cur, _ = l.fc.GetLatestBlock(ctx, true)
				end = cur.Height
				if end > start {
					// start has already been checked, add 1
					if err := l.process(ctx, start+1, end); err != nil {
						l.handleError(err)
					} else {
						start = end
					}
				}
			}
		}
	}()

	return l
}

func (l *Listener) Stop() {
	l.ticker.Stop()
	l.done <- true
	close(l.Events)
	l.ticker = nil
}