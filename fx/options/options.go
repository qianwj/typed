package options

const (
	AddressKey         = "address"
	HealthCheckKey     = "healthCheck"
	MetricsPortKey     = "metricsPort"
	defaultMetricsPort = 9090
)

type Options interface {
	GetValue(key string) (any, bool)
}

type options struct {
	data map[string]any
}

func (o *options) GetValue(key string) (any, bool) {
	data, ok := o.data[key]
	return data, ok
}

type metricsOptions struct {
	port int
}

func (o *metricsOptions) GetValue(key string) (any, bool) {
	switch key {
	case MetricsPortKey:
		return o.port, o.port > 0
	}
	return o, false
}

func Address(addr string) Options {
	return &options{data: map[string]any{AddressKey: addr}}
}

func DefaultMetrics() Options {
	return &metricsOptions{
		port: defaultMetricsPort,
	}
}

func Metrics(port int) Options {
	if port < 1 || port > 65535 {
		port = defaultMetricsPort
	}
	return &metricsOptions{port: defaultMetricsPort}
}

func EnableHealthCheck() Options {
	return &options{data: map[string]any{HealthCheckKey: true}}
}
