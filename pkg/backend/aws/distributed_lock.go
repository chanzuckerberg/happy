package aws

import (
	"context"
	"time"

	"cirello.io/dynamolock/v2"
	"github.com/pkg/errors"
)

// Adapted from:
// https://github.com/chanzuckerberg/developer-portal/blob/main/engportal/services/async/distributed_lock.go

// Dynamolock is an interface for a dynamodb distributed lock
type Dynamolock interface {
	CloseWithContext(context.Context) error
	AcquireLockWithContext(
		ctx context.Context,
		key string,
		opts ...dynamolock.AcquireLockOption) (*dynamolock.Lock, error)
	ReleaseLockWithContext(
		ctx context.Context,
		lock *dynamolock.Lock,
		opts ...dynamolock.ReleaseLockOption) (bool, error)
}

// DistributedLock allows for locking and coordination between multiple participants
type DistributedLock struct {
	dynamolock Dynamolock
}

// DistributedLockConfig configures a distributedlock
type DistributedLockConfig struct {
	DynamodbTableName string
	LeaseDuration     time.Duration
	HeartbeatPeriod   time.Duration
}

// NewDistributedLock creates a new distributed lock client
func NewDistributedLock(c *DistributedLockConfig, dynamoClient dynamolock.DynamoDBClient) (*DistributedLock, error) {
	leaseDuration := c.LeaseDuration
	if leaseDuration == 0 {
		leaseDuration = 10 * time.Second
	}

	heartbeatPeriod := c.HeartbeatPeriod
	if heartbeatPeriod == 0 {
		heartbeatPeriod = 2 * time.Second
	}

	client, err := dynamolock.New(
		dynamoClient,
		c.DynamodbTableName,
		dynamolock.WithLeaseDuration(leaseDuration),
		dynamolock.WithHeartbeatPeriod(heartbeatPeriod),
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate locking client")
	}

	return &DistributedLock{
		dynamolock: client,
	}, nil
}

// Close closes our distributed lock service cleanly
func (dl *DistributedLock) Close(ctx context.Context) error {
	return dl.dynamolock.CloseWithContext(ctx)
}

// AcquireLock will attempt to acquire a lock for a given key
func (dl *DistributedLock) AcquireLock(ctx context.Context, key string) (*dynamolock.Lock, error) {
	lock, err := dl.dynamolock.AcquireLockWithContext(ctx, key,
		dynamolock.WithDeleteLockOnRelease(),
		// Wait for at least a minute before giving up
		dynamolock.WithAdditionalTimeToWaitForLock(time.Minute),
	)
	var timeoutError *dynamolock.TimeoutError
	if errors.As(err, &timeoutError) {
		return nil, errors.Wrapf(err, "timed out waiting for lock for %s", key)
	}
	return lock, errors.Wrapf(err, "could not acquire lock for %s", key)
}

// ReleaseLock will release a lock for a given key
func (dl *DistributedLock) ReleaseLock(ctx context.Context, lock *dynamolock.Lock) (bool, error) {
	success, err := dl.dynamolock.ReleaseLockWithContext(ctx, lock)
	return success, errors.Wrap(err, "could not release lock")
}
