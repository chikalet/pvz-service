package metrics
import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)
var (
	PVZCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pvz_created_total",
		Help: "Total created pickup points",
	})
	IntakesOpened = promauto.NewCounter(prometheus.CounterOpts{
		Name: "intakes_opened_total",
		Help: "Total opened intakes",
	})
	IntakesClosed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "intakes_closed_total",
		Help: "Total closed intakes",
	})
	ItemsAdded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "items_added_total",
		Help: "Total items added to intakes",
	})
	ItemsQuantity = promauto.NewCounter(prometheus.CounterOpts{
		Name: "items_quantity_total",
		Help: "Total quantity of all items",
	})
	HTTPRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
	HTTPDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.3, 1, 3, 5},
		},
		[]string{"path"},
	)
)