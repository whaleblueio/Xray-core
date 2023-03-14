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

	ipCounter := GetIPCounter(u.Email)
	if ipCounter == nil {
		ipCounter = &IpCounter{
			IpTable: make(map[string]*ConnIP),
			Email:   u.Email,
		}
		AddIp(u.Email, ipCounter)
	}
	return &MemoryUser{
		Account:   account,
		Email:     u.Email,
		Level:     u.Level,
		IpCounter: ipCounter,
	}, nil
}

// MemoryUser is a parsed form of User, to reduce number of parsing of Account proto.
type MemoryUser struct {
	// Account is the parsed account of the protocol.
	Account   Account
	Email     string
	Level     uint32
	IpCounter *IpCounter
}

var buckets sync.Map

var connections sync.Map

func AddIp(email string, ipCounter *IpCounter) {

	connections.Store(email, ipCounter)
	newError("AddIp() email:", email, " do not have counter pointer:", &ipCounter, " created one").AtDebug().WriteToLog()
}
func GetIPCounter(email string) *IpCounter {

	if c, ok := connections.Load(email); ok {
		return c.(*IpCounter)
	} else {
		return nil
	}
}

func GetIPs(email string) []string {
	var ips []string
	if connection, ok := connections.Load(email); ok {
		c := connection.(*IpCounter)
		//newError("GetIPs() email:", email, " have ", len(c.IpTable), " connections, connected ips:", c.IpTable, " counter pointer:", &c).WriteToLog()
		for k, ip := range c.IpTable {
			interval := time.Now().Unix() - ip.Time
			//over 1 minutes not update ,will delete
			if interval > 1*60 {
				newError("GetIPs() email:", email, " IP:", ip.IP, " over 30 seconds not updated,delete.").AtDebug().WriteToLog()
				c.Del(ip.IP)
			} else {
				ips = append(ips, k)
			}
		}
	} else {
		newError("GetIPs() email:", email, " do not have ip connected", email).WriteToLog()
	}
	return ips
}

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
			bucket = rateLimit.NewBucketWithQuantum(time.Second, u.SpeedLimiter.Speed, u.SpeedLimiter.Speed)
			bucket := b.(*rateLimit.Bucket)
			if bucket.Capacity() != u.SpeedLimiter.Speed {
				newError(fmt.Sprintf("user:%s update speed limit to :%d", u.Email, u.SpeedLimiter.Speed)).WriteToLog()
				bucket = rateLimit.NewBucketWithQuantum(time.Second, u.SpeedLimiter.Speed, u.SpeedLimiter.Speed)
				buckets.Store(u.Email, bucket)
			}
		}

		if u.SpeedLimiter != nil && u.SpeedLimiter.Speed < 0 {
			buckets.Delete(u.Email)
		}
	} else {
		if u.SpeedLimiter != nil && u.SpeedLimiter.Speed > 0 {
			bucket = rateLimit.NewBucketWithQuantum(time.Second, u.SpeedLimiter.Speed, u.SpeedLimiter.Speed)
			buckets.Store(u.Email, bucket)
			newError(fmt.Sprintf("user:%s speed limit:%d", u.Email, u.SpeedLimiter.Speed)).WriteToLog()
		}
	}
}
