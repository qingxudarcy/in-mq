package inmq

import (
	"log"
	"os"
	"time"
)

type ConnConfig struct {
	Heartbeat time.Duration // default: 20s
	TimeOut   time.Duration // default: 30s
	Logger    Logx
}

var defaultRcConfig = &ConnConfig{
	Heartbeat: 20 * time.Second,
	TimeOut:   30 * time.Second,
	Logger:    log.New(os.Stderr, "", log.LstdFlags),
}

type connOption struct {
	apply func(rc *ConnConfig)
}

func WithRcHeartbeat(heartbeat time.Duration) *connOption {
	return &connOption{
		apply: func(rc *ConnConfig) {
			rc.Heartbeat = heartbeat
		},
	}
}

func WithRcLogger(logger Logx) *connOption {
	return &connOption{
		apply: func(rc *ConnConfig) {
			rc.Logger = logger
		},
	}
}

func WithRcTimeOut(timeout time.Duration) *connOption {
	return &connOption{
		apply: func(rc *ConnConfig) {
			rc.TimeOut = timeout
		},
	}
}

func ConnCfgParse(rcOpts ...*connOption) *ConnConfig {
	rc := defaultRcConfig
	for _, opt := range rcOpts {
		opt.apply(rc)
	}
	return rc
}
