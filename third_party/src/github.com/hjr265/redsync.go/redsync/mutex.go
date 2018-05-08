// Package redsync provides a Redis-based distributed mutual exclusion lock implementation as described in the blog post http://antirez.com/news/77.
//
// Values containing the types defined in this package should not be copied.
package redsync

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	// DefaultExpiry is used when Mutex Duration is 0
	DefaultExpiry = 8 * time.Second
	// DefaultTries is used when Mutex Duration is 0
	DefaultTries = 16
	// DefaultDelay is used when Mutex Delay is 0
	DefaultDelay = 512 * time.Millisecond
	// DefaultFactor is used when Mutex Factor is 0
	DefaultFactor = 0.01
)

var (
	// ErrFailed is returned when lock cannot be acquired
	ErrFailed = errors.New("failed to acquire lock")
)

// Locker interface with Lock returning an error when lock cannot be aquired
type Locker interface {
	Lock() error
	Unlock()
}

// A Mutex is a mutual exclusion lock.
//
// Fields of a Mutex must not be changed after first use.
type Mutex struct {
	Name   string        // Resouce name
	Expiry time.Duration // Duration for which the lock is valid, DefaultExpiry if 0

	Tries int           // Number of attempts to acquire lock before admitting failure, DefaultTries if 0
	Delay time.Duration // Delay between two attempts to acquire lock, DefaultDelay if 0

	Factor float64 // Drift factor, DefaultFactor if 0

	Quorum int // Quorum for the lock, set to len(addrs)/2+1 by NewMutex()

	value string
	until time.Time

	nodes []*redis.Pool
	nodem sync.Mutex
}

var _ = Locker(&Mutex{})

// NewMutex returns a new Mutex on a named resource connected to the Redis instances at given addresses.
func NewMutex(name string, addrs []net.Addr) (*Mutex, error) {
	if len(addrs) == 0 {
		panic("redsync: addrs is empty")
	}

	nodes := make([]*redis.Pool, len(addrs))
	for i, addr := range addrs {
		dialTo := addr
		connTimeout := time.Duration(5) * time.Second
		readTimeout := time.Duration(30) * time.Second
		writeTimeout := time.Duration(30) * time.Second
		node := &redis.Pool{
			MaxActive: 1,
			Wait:      true,
			Dial: func() (redis.Conn, error) {
				c, err := redis.DialTimeout("tcp", dialTo.String(), connTimeout, readTimeout, writeTimeout)
				if err != nil {
					return nil, err
				}
				return c, err
			},
		}
		nodes[i] = node
	}

	return NewMutexWithPool(name, nodes)
}

// NewMutexWithPool returns a new Mutex on a named resource connected to the Redis instances at given redis Pools.
func NewMutexWithPool(name string, nodes []*redis.Pool) (*Mutex, error) {
	if len(nodes) == 0 {
		panic("redsync: nodes is empty")
	}

	return &Mutex{
		Name:   name,
		Quorum: len(nodes)/2 + 1,
		nodes:  nodes,
	}, nil
}

// Lock locks m.
// In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *Mutex) Lock() error {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	value := base64.StdEncoding.EncodeToString(b)

	expiry := m.Expiry
	if expiry == 0 {
		expiry = DefaultExpiry
	}

	retries := m.Tries
	if retries == 0 {
		retries = DefaultTries
	}

	for i := 0; i < retries; i++ {
		n := 0
		start := time.Now()
		for _, node := range m.nodes {
			if node == nil {
				continue
			}

			conn := node.Get()
			reply, err := redis.String(conn.Do("set", m.Name, value, "nx", "px", int(expiry/time.Millisecond)))
			conn.Close()
			if err != nil {
				continue
			}
			if reply != "OK" {
				continue
			}
			n++
		}

		factor := m.Factor
		if factor == 0 {
			factor = DefaultFactor
		}

		until := time.Now().Add(expiry - time.Now().Sub(start) - time.Duration(int64(float64(expiry)*factor)) + 2*time.Millisecond)
		if n >= m.Quorum && time.Now().Before(until) {
			m.value = value
			m.until = until
			return nil
		}
		for _, node := range m.nodes {
			if node == nil {
				continue
			}

			conn := node.Get()
			_, err := delScript.Do(conn, m.Name, value)
			conn.Close()
			if err != nil {
				continue
			}
		}

		delay := m.Delay
		if delay == 0 {
			delay = DefaultDelay
		}
		time.Sleep(delay)
	}

	return ErrFailed
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
func (m *Mutex) Unlock() {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	value := m.value
	if value == "" {
		panic("redsync: unlock of unlocked mutex")
	}

	m.value = ""
	m.until = time.Unix(0, 0)

	for _, node := range m.nodes {
		if node == nil {
			continue
		}

		conn := node.Get()
		delScript.Do(conn, m.Name, value)
		conn.Close()
	}
}

var delScript = redis.NewScript(1, `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`)