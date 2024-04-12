package logs

import (
	"context"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxRepository struct {
	*InfluxClient
}

func NewInfluxRepository(influxClient *InfluxClient) *InfluxRepository {
	return &InfluxRepository{InfluxClient: influxClient}
}

func (i *InfluxRepository) Writer(ctx context.Context, log BehaviorLog) error {
	tags := make(map[string]string)
	for _, tag := range log.Tags {
		tags[tag] = tag
	}

	point := write.NewPoint(i.Measurement, tags, map[string]interface{}{
		"uid": log.UID,
		"ip":  log.IP,
		"ua":  log.UA,
	}, time.Now().In(time.Local))

	if err := i.WriteAPI.WritePoint(ctx, point); err != nil {
		return err
	}

	return nil
}
