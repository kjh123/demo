package logs

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type ClickHouseRepository struct {
	*ClickHouseClient
}

func NewClickHouseRepository(client *ClickHouseClient) *ClickHouseRepository {
	return &ClickHouseRepository{ClickHouseClient: client}
}

func (r *ClickHouseRepository) Writer(ctx context.Context, log BehaviorLog) error {
	err := r.Conn.AsyncInsert(ctx,
		fmt.Sprintf(
			"INSERT INTO %s.%s (uid, ip, tags, ua) VALUES (?, ?, ?, ?)",
			r.DB, r.Table,
		),
		false,
		log.UID, log.IP, log.Tags, log.UA,
	)
	if err != nil {
		return errors.Wrap(err, "ClickHouseRepository insert logs")
	}

	return nil
}
