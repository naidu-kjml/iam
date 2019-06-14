package cfg

import (
	"time"
)

// Expirations contains expiration times for anything that needs to be cached
var Expirations = struct {
	User             time.Duration
	GroupMemberships time.Duration
	GroupsLastSync   time.Duration
	Teams            time.Duration
}{
	User:             time.Hour * 24,
	GroupMemberships: time.Hour * 24,
	GroupsLastSync:   0,
	Teams:            0,
}
