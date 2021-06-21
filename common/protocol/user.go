package protocol

import (
	"fmt"
	rateLimit "github.com/juju/ratelimit"
	"sync"
	"time"
)

func (u *User) GetTypedAccount() (Account, error) {
	if u.GetAccount() == nil {
		return nil, newError("Account missing").AtWarning()
	}

	rawAccount, err := u.Account.GetInstance()
	if err != nil {
		return nil, err
	}
	if asAccount, ok := rawAccount.(AsAccount); ok {
		return asAccount.AsAccount()
	}
	if account, ok := rawAccount.(Account); ok {
		return account, nil
	}
	return nil, newError("Unknown account type: ", u.Account.Type)
}

func (u *User) ToMemoryUser() (*MemoryUser, error) {
	account, err := u.GetTypedAccount()
	if err != nil {
		return nil, err
	}

	return &MemoryUser{
		Account: account,
		Email:   u.Email,
		Level:   u.Level,
	}, nil
}

// MemoryUser is a parsed form of User, to reduce number of parsing of Account proto.
type MemoryUser struct {
	// Account is the parsed account of the protocol.
	Account Account
	Email   string
	Level   uint32
}

var buckets sync.Map

func GetBucket(email string) *rateLimit.Bucket {
	b, ok := buckets.Load(email)
	if ok {
		return b.(*rateLimit.Bucket)
	}
	return nil
}
func SetBucket(u *User) {
	var bucket *rateLimit.Bucket
	b, ok := buckets.Load(u.Email)
	if ok {
		if u.SpeedLimiter != nil && u.SpeedLimiter.Speed > 0 {
			bucket = rateLimit.NewBucketWithQuantum(time.Second, u.SpeedLimiter.Speed, u.SpeedLimiter.Speed/2)
			bucket := b.(*rateLimit.Bucket)
			if bucket.Capacity() != u.SpeedLimiter.Speed {
				newError(fmt.Sprintf("user:%s update speed limit to :%d", u.Email, u.SpeedLimiter.Speed)).WriteToLog()
				bucket = rateLimit.NewBucketWithQuantum(time.Second, u.SpeedLimiter.Speed, u.SpeedLimiter.Speed/2)
				buckets.Store(u.Email, bucket)
			}
		}

		if u.SpeedLimiter != nil && u.SpeedLimiter.Speed < 0 {
			buckets.Delete(u.Email)
		}
	} else {
		if u.SpeedLimiter != nil && u.SpeedLimiter.Speed > 0 {
			bucket = rateLimit.NewBucketWithQuantum(time.Second, u.SpeedLimiter.Speed, u.SpeedLimiter.Speed/2)
			buckets.Store(u.Email, bucket)
			newError(fmt.Sprintf("user:%s speed limit:%d", u.Email, u.SpeedLimiter.Speed)).WriteToLog()
		} else {
			newError(fmt.Sprintf("user:%s no speed limit", u.Email)).WriteToLog()
		}
	}
}
