package otel

import (
	"time"

	"github.com/segmentio/ksuid"

	"gorm.io/gorm"
)

// Logs are designed to be written via an OTLP exporter.
//
// https://opentelemetry.io/docs/specs/otel/logs/bridge-api/
//
// The clickhouse exporter, is a good reference point for this
// https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/exporter/clickhouseexporter/exporter_logs.go
type OtelLogRecord struct {
	ID        string    `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"notnull"`
	UpdatedAt time.Time `json:"updated_at" gorm:"notnull"`

	// OTEL log message attributes
	Timestamp        time.Time           `json:"timestamp" gorm:"type:DateTime64(9);codec:Delta(8),ZSTD(1)"`
	Exemplars        []map[string]string `json:"-" gorm:"type:Nested(filtered_attributes Map(LowCardinality(String), String), time_unix DateTime64(9), value Float64, span_id String, trace_id String); codec:ZSTD(1);"`
	EventsTimestamp  []time.Time         `json:"-" gorm:"type:DateTime64(9);column:events.timestamp"`
	EventsName       []string            `json:"-" gorm:"type:LowCardinality(String);column:events.name"`
	EventsAttributes []map[string]string `json:"-" gorm:"type:Map(LowCardinality(String), String);column:events.attributes"`
}

func (r *OtelLogRecord) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = ksuid.New().String()
	}
	return nil
}

func (r OtelLogRecord) GetTableOptions() (string, bool) {
	opts := `ENGINE = MergeTree()
	PARTITION BY toDate(timestamp_time)
	PRIMARY KEY (service_name, timestamp_time)
	ORDER BY (service_name, timestamp_time, timestamp)
	TTL toDateTime("timestamp") + toIntervalDay(720)
	SETTINGS index_granularity = 8192, ttl_only_drop_parts = 0;`
	return opts, true
}

func (r OtelLogRecord) MigrateDB(db *gorm.DB) *gorm.DB {
	opts, hasOpts := r.GetTableOptions()
	if !hasOpts {
		return db
	}
	return db.Set("gorm:table_options", opts)
}
