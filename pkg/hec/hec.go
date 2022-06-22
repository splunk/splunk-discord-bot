package hec

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/splunkhecexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/otel/metric/nonrecording"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

type HecClient interface {
	SendLogs(logs plog.Logs) error
	SendData(timestamp time.Time, data []byte) error
	Stop() error
}

type HecClientImpl struct {
	exporter component.LogsExporter
}

func (h *HecClientImpl) SendData(timestamp time.Time, data []byte) error {
	l := plog.NewLogs()
	logRecord := l.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
	logRecord.Body().SetStringVal(string(data))
	logRecord.SetTimestamp(pcommon.NewTimestampFromTime(timestamp))
	return h.exporter.ConsumeLogs(context.Background(), l)
}

func (h *HecClientImpl) SendLogs(logs plog.Logs) error {
	return h.exporter.ConsumeLogs(context.Background(), logs)
}

func (h *HecClientImpl) Stop() error {
	return h.exporter.Shutdown(context.Background())
}

func CreateClient(endpoint string, token string, insecureSkipVerify bool, index string, logger *zap.Logger) (HecClient, error) {
	factory := splunkhecexporter.NewFactory()
	hecConfig := factory.CreateDefaultConfig().(*splunkhecexporter.Config)
	hecConfig.LogDataEnabled = true
	hecConfig.Endpoint = endpoint
	hecConfig.Token = token
	hecConfig.TLSSetting.InsecureSkipVerify = insecureSkipVerify
	hecConfig.Index = index
	exporter, err := factory.CreateLogsExporter(context.Background(), component.ExporterCreateSettings{
		TelemetrySettings: component.TelemetrySettings{
			Logger:         logger,
			TracerProvider: trace.NewNoopTracerProvider(),
			MeterProvider:  nonrecording.NewNoopMeterProvider(),
			MetricsLevel:   configtelemetry.LevelNone,
		},
		BuildInfo: component.NewDefaultBuildInfo(),
	}, hecConfig)
	if err != nil {
		return nil, err
	}
	if err = exporter.Start(context.Background(), componenttest.NewNopHost()); err != nil {
		return nil, err
	}

	batchProcessorFactory := batchprocessor.NewFactory()
	processor, err := batchProcessorFactory.CreateLogsProcessor(context.Background(), componenttest.NewNopProcessorCreateSettings(), batchProcessorFactory.CreateDefaultConfig(), exporter)
	if err != nil {
		return nil, err
	}
	if err = processor.Start(context.Background(), componenttest.NewNopHost()); err != nil {
		return nil, err
	}

	return &HecClientImpl{exporter: processor}, nil
}
