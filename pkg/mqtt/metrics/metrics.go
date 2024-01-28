package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Metrics for HandleTCPConnections
	BasicAuthAcceptedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_bridge_basic_auth_accepted_total",
		Help: "The total number of basic auth connections accepted",
	})
	BasicAuthFailedToAcceptTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_bridge_basic_auth_failed_to_accept_total",
		Help: "The total number of basic auth connections failed to accept",
	})
	BasicAuthDeviceAuthFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_bridge_basic_auth_device_auth_failed_total",
		Help: "The total number of basic auth connections closed because the device is not registered, disabled or otherwise not allowed to connect",
	})

	// Metrics for HandleTLSConnections
	TlsAcceptedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_bridge_tls_accepted_total",
		Help: "The total number of TLS connections accepted",
	})
	TlsFailedToAcceptTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_bridge_tls_failed_to_accept_total",
		Help: "The total number of TLS connections failed to accept",
	})
	TlsDeviceAuthFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_bridge_tls_bad_certificate_total",
		Help: "The total number of TLS connections closed because the provided certificate is absent, bad or device is not registered, disabled or otherwise not allowed to connect",
	})

	// Metrics for HandleConn (both BasicAuth and TLS)
	ActiveConnectionsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mqtt_bridge_active_connections_total",
		Help: "The total number of active connections",
	})

	// Other common metrics
	ConnNotAnMqttPacketTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_bridge_conn_not_an_mqtt_packet_total",
		Help: "The total number of connections closed because the packet is not an MQTT packet or malformed otherwise",
	})
)
