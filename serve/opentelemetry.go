package serve

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/rs/zerolog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otellog "go.opentelemetry.io/otel/log"
	logglobal "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// newResource returns a resource describing this application.
func newResource(p *plugin.Plugin) *resource.Resource {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("cloudquery-"+p.Name()),
			semconv.ServiceVersion(p.Version()),
		),
	)
	if err != nil {
		panic(err)
	}
	return r
}

type otelConfig struct {
	endpoint string
	insecure bool
}

func getTraceExporter(ctx context.Context, opts otelConfig) (*otlptrace.Exporter, error) {
	if opts.endpoint == "" {
		return nil, nil
	}

	traceOptions := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(opts.endpoint),
	}

	if opts.insecure {
		traceOptions = append(traceOptions, otlptracehttp.WithInsecure())
	}

	traceClient := otlptracehttp.NewClient(traceOptions...)
	traceExporter, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	return traceExporter, nil
}

func getMetricReader(ctx context.Context, opts otelConfig) (*metric.PeriodicReader, error) {
	if opts.endpoint == "" {
		return nil, nil
	}

	metricOptions := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(opts.endpoint),
	}

	if opts.insecure {
		metricOptions = append(metricOptions, otlpmetrichttp.WithInsecure())
	}

	metricExporter, err := otlpmetrichttp.New(ctx, metricOptions...)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP metric exporter: %w", err)
	}

	reader := metric.NewPeriodicReader(metricExporter, metric.WithInterval(15*time.Second))
	return reader, nil
}

func getLogsProcessor(ctx context.Context, opts otelConfig) (*log.BatchProcessor, error) {
	if opts.endpoint == "" {
		return nil, nil
	}

	logOptions := []otlploghttp.Option{
		otlploghttp.WithEndpoint(opts.endpoint),
		otlploghttp.WithCompression(otlploghttp.GzipCompression),
	}

	if opts.insecure {
		logOptions = append(logOptions, otlploghttp.WithInsecure())
	}

	exporter, err := otlploghttp.New(ctx, logOptions...)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP log exporter: %w", err)
	}

	processor := log.NewBatchProcessor(exporter)
	return processor, nil
}

func setupOtel(ctx context.Context, logger zerolog.Logger, p *plugin.Plugin, otelEndpoint string, otelEndpointInsecure bool) (shutdown func(), err error) {
	if otelEndpoint == "" {
		return nil, nil
	}
	opts := otelConfig{
		endpoint: otelEndpoint,
		insecure: otelEndpointInsecure,
	}
	traceExporter, err := getTraceExporter(ctx, opts)
	if err != nil {
		return nil, err
	}

	metricReader, err := getMetricReader(ctx, opts)
	if err != nil {
		return nil, err
	}

	logsProcessor, err := getLogsProcessor(ctx, opts)
	if err != nil {
		return nil, err
	}

	pluginResource := newResource(p)
	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(pluginResource),
	)

	mt := metric.NewMeterProvider(
		metric.WithReader(metricReader),
		metric.WithResource(pluginResource),
	)

	lp := log.NewLoggerProvider(
		log.WithProcessor(logsProcessor),
		log.WithResource(pluginResource),
	)

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Warn().Err(err).Msg("otel error")
	}))
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mt)
	logglobal.SetLoggerProvider(lp)

	shutdown = func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Error().Err(err).Msg("failed to shutdown OTLP trace provider")
		}
		if err := mt.Shutdown(context.Background()); err != nil {
			logger.Error().Err(err).Msg("failed to shutdown OTLP metric provider")
		}
		if err := lp.Shutdown(context.Background()); err != nil {
			logger.Error().Err(err).Msg("failed to shutdown OTLP logger provider")
		}
	}

	return shutdown, nil
}

// Similar to https://github.com/AkhigbeEromo/opentelemetry-go-contrib/blob/dedcf91a55a36a5a8589c56f2e43c188eb42f4f2/bridges/otelzerolog/hook.go
// but with `TraceLevel` and attributes support
type otelLoggerHook struct {
	otellog.Logger
	ctx context.Context
}

func (h *otelLoggerHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	record := otellog.Record{}
	record.SetTimestamp(time.Now().UTC())
	record.SetSeverity(otellogSeverity(level))
	record.SetBody(otellog.StringValue(message))
	// See https://github.com/rs/zerolog/issues/493, this is ugly but it works
	// At the moment there's no way to get the log fields from the event, so we use reflection to get the buffer and parse it
	// TODO: Remove this if https://github.com/rs/zerolog/pull/682 is merged
	logData := make(map[string]any)
	eventBuffer := fmt.Sprintf("%s}", reflect.ValueOf(e).Elem().FieldByName("buf"))
	err := json.Unmarshal([]byte(eventBuffer), &logData)
	if err == nil {
		recordAttributes := make([]otellog.KeyValue, 0, len(logData))
		for k, v := range logData {
			if k == "level" {
				continue
			}
			if k == "time" {
				eventTimestamp, ok := v.(string)
				if !ok {
					continue
				}
				t, err := time.Parse(time.RFC3339Nano, eventTimestamp)
				if err == nil {
					record.SetTimestamp(t)
					continue
				}
			}
			var attributeValue otellog.Value
			switch v := v.(type) {
			case string:
				attributeValue = otellog.StringValue(v)
			case int:
				attributeValue = otellog.IntValue(v)
			case int64:
				attributeValue = otellog.Int64Value(v)
			case float64:
				attributeValue = otellog.Float64Value(v)
			case bool:
				attributeValue = otellog.BoolValue(v)
			case []byte:
				attributeValue = otellog.BytesValue(v)
			default:
				attributeValue = otellog.StringValue(fmt.Sprintf("%v", v))
			}
			recordAttributes = append(recordAttributes, otellog.KeyValue{
				Key:   k,
				Value: attributeValue,
			})
		}
		record.AddAttributes(recordAttributes...)
	}

	h.Emit(h.ctx, record)
}

func otellogSeverity(level zerolog.Level) otellog.Severity {
	switch level {
	case zerolog.DebugLevel:
		return otellog.SeverityDebug
	case zerolog.InfoLevel:
		return otellog.SeverityInfo
	case zerolog.WarnLevel:
		return otellog.SeverityWarn
	case zerolog.ErrorLevel:
		return otellog.SeverityError
	case zerolog.FatalLevel:
		return otellog.SeverityFatal2
	case zerolog.PanicLevel:
		return otellog.SeverityFatal1
	case zerolog.TraceLevel:
		return otellog.SeverityTrace
	default:
		return otellog.SeverityUndefined
	}
}

func newOTELLoggerHook() zerolog.Hook {
	return &otelLoggerHook{logglobal.GetLoggerProvider().Logger("cloudquery"), context.Background()}
}
