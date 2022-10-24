package protocol

import (
	"sync"
	"time"
)

type IpCounter struct {
	ipTableLock sync.Mutex
	IpTable     map[string]*ConnIP
}

type ConnIP struct {
	IP   string
	Time int64
}

// Add implements stats.IpCounter.
func (c *IpCounter) Add(ip string) {
	c.ipTableLock.Lock()
	defer c.ipTableLock.Unlock()
	connected, found := c.IpTable[ip]
	if found {
		connected.Time = time.Now().Unix()
		return
	}
	c.IpTable[ip] = &ConnIP{
		IP:   ip,
		Time: time.Now().Unix(),
	}
}

// Del implements stats.IpCounter.
func (c *IpCounter) Del(ip string) {
	c.ipTableLock.Lock()
	defer c.ipTableLock.Unlock()

	_, found := c.IpTable[ip]
	if !found {
		return
	}
	delete(c.IpTable, ip)
}

func (c *IpCounter) getIP(ip string) *ConnIP {
	IPCon, found := c.IpTable[ip]
	if !found {
		return IPCon
	}
	return nil
}
