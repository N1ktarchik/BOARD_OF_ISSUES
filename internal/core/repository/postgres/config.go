package postgres

import "time"

type DatabaseConfig struct {
	connStr           string
	maxConns          int32
	minConns          int32
	maxConnLifetime   time.Duration
	maxConnIdleTime   time.Duration
	healthCheckPeriod time.Duration
}

func NewDatabaseConfig(connStr string, maxConns, minConns int32, connMaxLifetime, connMaxIdleTime, healthCheckPeriod time.Duration) *DatabaseConfig {
	return &DatabaseConfig{
		connStr:           connStr,
		maxConns:          maxConns,
		minConns:          minConns,
		maxConnLifetime:   connMaxLifetime,
		maxConnIdleTime:   connMaxIdleTime,
		healthCheckPeriod: healthCheckPeriod,
	}
}
