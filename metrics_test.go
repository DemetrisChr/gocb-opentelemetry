package gocbopentelemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestOpenTelemetryMetricsInSeconds(t *testing.T) {
	rdr := metric.NewManualReader()

	provider := metric.NewMeterProvider(
		metric.WithReader(rdr),
	)

	meter := NewOpenTelemetryMeter(provider)
	recorder, err := meter.ValueRecorder("test_recorder", map[string]string{
		"foo":    "bar",
		"__unit": "s",
	})
	require.Nil(t, err)

	recorder.RecordValue(2_000_000)
	recorder.RecordValue(500_000)

	var data metricdata.ResourceMetrics
	err = rdr.Collect(context.Background(), &data)
	require.Nil(t, err)

	require.Len(t, data.ScopeMetrics, 1)
	require.Len(t, data.ScopeMetrics[0].Metrics, 1)
	require.Equal(t, "s", data.ScopeMetrics[0].Metrics[0].Unit)

	histogram, ok := data.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
	require.True(t, ok)

	require.Len(t, histogram.DataPoints, 1)

	dataPoint := histogram.DataPoints[0]
	assert.Equal(t, 2.5, dataPoint.Sum)
	assert.Equal(t, uint64(2), dataPoint.Count)

	assert.Equal(t, attribute.NewSet(attribute.String("foo", "bar")), dataPoint.Attributes)
}
