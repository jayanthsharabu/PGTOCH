package poller

import (
	"context"
	"fmt"
	"pgtoch/internal/etl"
	"pgtoch/internal/log"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type PollConfig struct {
	Table     string
	DeltaCol  string
	Interval  time.Duration
	Limit     *int
	StartFrom string
	OnData    func(data *etl.TableData) error
}

type Poller struct {
	conn   *pgx.Conn
	config PollConfig
}

func NewPoller(conn *pgx.Conn, config PollConfig) *Poller {
	return &Poller{
		conn:   conn,
		config: config,
	}
}

func (p *Poller) Start(ctx context.Context) error {

	lastSeen := p.config.StartFrom

	ticker := time.NewTicker(p.config.Interval)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Logger.Info("Stopping ctx cancelled")
			return ctx.Err()
		case <-ticker.C:
			log.Logger.Info("Polling for new data",
				zap.String("table", p.config.Table),
				zap.String("last_seen", lastSeen),
			)
			data, err := etl.ExtractTableDataSince(ctx, p.conn, p.config.Table, p.config.DeltaCol, lastSeen, p.config.Limit)
			if err != nil {
				log.Logger.Error("Error extracting table data",
					zap.Error(err),
					zap.String("table", p.config.Table),
				)
				continue
			}

			if len(data.Rows) == 0 {
				log.Logger.Info("No new data found in this cycle",
					zap.String("table", p.config.Table),
					zap.String("last_seen", lastSeen),
				)
				continue
			}

			lastRows := data.Rows[len(data.Rows)-1]

			for i, col := range data.Columns {
				if col.Name == p.config.DeltaCol {
					switch v := lastRows[i].(type) {
					case string:
						lastSeen = v
					case int, int32, int64, int8, int16:
						lastSeen = fmt.Sprintf("%d", v)
					case float64, float32:
						lastSeen = fmt.Sprintf("%f", v)
					default:
						lastSeen = fmt.Sprintf("%v", v)
					}
					break
				}
			}

			log.Logger.Info("New data extracted",
				zap.Int("rows", len(data.Rows)),
				zap.String("table", p.config.Table),
				zap.String("last_seen", lastSeen),
			)

			if err := p.config.OnData(data); err != nil {
				log.Logger.Error("Failed to process extracted data",
					zap.Error(err),
					zap.String("table", p.config.Table),
				)
				continue
			}
		}
	}
}
