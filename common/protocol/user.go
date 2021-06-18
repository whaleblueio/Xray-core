package protocol

import (
	"fmt"
	rateLimit "github.com/juju/ratelimit"
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
	var bucket *rateLimit.Bucket
	if u.SpeedLimiter != nil && u.SpeedLimiter.Speed > 0 {
		bucket = rateLimit.NewBucketWithQuantum(time.Second, u.SpeedLimiter.Speed, u.SpeedLimiter.Speed)
		newError(fmt.Sprintf("user:%s speed limit:%d", u.Email, u.SpeedLimiter.Speed)).WriteToLog()
	}
	return &MemoryUser{
		Account: account,
		Bucket:  bucket,
		Email:   u.Email,
		Level:   u.Level,
	}, nil
}

// MemoryUser is a parsed form of User, to reduce number of parsing of Account proto.
type MemoryUser struct {
	// Account is the parsed account of the protocol.
	Account Account
	Bucket  *rateLimit.Bucket
	Email   string
	Level   uint32
}
