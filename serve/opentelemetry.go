package serve

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/rs/zerolog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
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

func parseOtelHeaders(headers []string) map[string]string {
	headerMap := make(map[string]string, len(headers))
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			continue
		}
		headerMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return headerMap
}

type otelConfig struct {
	endpoint string
	insecure bool
	headers  []string
	urlPath  string
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

	if len(opts.headers) > 0 {
		headers := parseOtelHeaders(opts.headers)
		traceOptions = append(traceOptions, otlptracehttp.WithHeaders(headers))
	}

	if opts.urlPath != "" {
		traceOptions = append(traceOptions, otlptracehttp.WithURLPath(opts.urlPath))
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

	if len(opts.headers) > 0 {
		headers := parseOtelHeaders(opts.headers)
		metricOptions = append(metricOptions, otlpmetrichttp.WithHeaders(headers))
	}

	if opts.urlPath != "" {
		metricOptions = append(metricOptions, otlpmetrichttp.WithURLPath(opts.urlPath))
	}

	metricExporter, err := otlpmetrichttp.New(ctx, metricOptions...)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP metric exporter: %w", err)
	}

	reader := metric.NewPeriodicReader(metricExporter, metric.WithInterval(15*time.Second))
	return reader, nil
}

func setupOtel(ctx context.Context, logger zerolog.Logger, p *plugin.Plugin, otelEndpoint string, otelEndpointInsecure bool, otelEndpointHeaders []string, otelEndpointURLPath string) (shutdown func(), err error) {
	if otelEndpoint == "" {
		return func() {}, nil
	}
	opts := otelConfig{
		endpoint: otelEndpoint,
		insecure: otelEndpointInsecure,
		headers:  otelEndpointHeaders,
		urlPath:  otelEndpointURLPath,
	}
	traceExporter, err := getTraceExporter(ctx, opts)
	if err != nil {
		return nil, err
	}

	metricReader, err := getMetricReader(ctx, opts)
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
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Debug().Err(err).Msg("otel error")
	}))
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mt)

	shutdown = func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Error().Err(err).Msg("failed to shutdown OTLP trace exporter")
		}
		if err := mt.Shutdown(context.Background()); err != nil {
			logger.Error().Err(err).Msg("failed to shutdown OTLP metric exporter")
		}
	}

	return shutdown, nil
}
