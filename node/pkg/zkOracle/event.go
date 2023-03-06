package zkOracle

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/event"
)

func WatchEvent[K any](
	ctx context.Context,
	subscribeLog func(*bind.WatchOpts, chan<- K) (event.Subscription, error),
	handleEvent func(context.Context, K) error,
) error {
	events := make(chan K)
	defer close(events)

	sub, err := subscribeLog(
		&bind.WatchOpts{
			Context: ctx,
		},
		events,
	)
	if err != nil {
		return fmt.Errorf("subscribe log: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case e := <-events:
			if handleEvent != nil {
				if err := handleEvent(ctx, e); err != nil {
					return fmt.Errorf("handle event: %w", err)
				}
			}
		case err := <-sub.Err():
			return fmt.Errorf("subscription: %w", err)
		case <-ctx.Done():
			return fmt.Errorf("context: %w", ctx.Err())
		}
	}
}

func WaitEvent[K any](
	ctx context.Context,
	subscribeLog func(*bind.WatchOpts, chan<- K) (event.Subscription, error),
) (K, error) {
	events := make(chan K)
	defer close(events)

	var receivedEvent K

	sub, err := subscribeLog(
		&bind.WatchOpts{
			Context: ctx,
		},
		events,
	)
	if err != nil {
		return receivedEvent, fmt.Errorf("subscribe log: %w", err)
	}
	defer sub.Unsubscribe()

	select {
	case e := <-events:
		receivedEvent = e
		return receivedEvent, nil
	case err := <-sub.Err():
		return receivedEvent, fmt.Errorf("subscription: %w", err)
	case <-ctx.Done():
		return receivedEvent, fmt.Errorf("context: %w", ctx.Err())
	}

}
