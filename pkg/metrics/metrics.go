package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	compbasemetrics "k8s.io/component-base/metrics"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/util/metrics"
)

const (
	SchedulerSubsystem = "llm_d_inference_scheduler"
)

var (
	SchedulerPrefillSelectionDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Subsystem: SchedulerSubsystem,
			Name:      "prefill_selection_duration_seconds",
			Help:      metrics.HelpMsgWithStability("Time taken to select a prefill pod", compbasemetrics.ALPHA),
			Buckets:   prometheus.DefBuckets,
		},
	)
	SchedulerDecodeSelectionDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Subsystem: SchedulerSubsystem,
			Name:      "decode_selection_duration_seconds",
			Help:      metrics.HelpMsgWithStability("Time taken to select a decode pod", compbasemetrics.ALPHA),
			Buckets:   prometheus.DefBuckets,
		},
	)
	SchedulerPDDecisionCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: SchedulerSubsystem,
			Name:      "pd_decision_total",
			Help:      metrics.HelpMsgWithStability("Total number of P/D disaggregation decisions made", compbasemetrics.ALPHA),
		},
		[]string{"decision_type"}, // "split" or "combined"
	)
	SchedulerPDThresholdHitsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: SchedulerSubsystem,
			Name:      "pd_threshold_hits_total",
			Help:      metrics.HelpMsgWithStability("Total number of times the P/D token threshold was met, labeled by the configured threshold", compbasemetrics.ALPHA),
		},
		[]string{"threshold"},
	)
	SchedulerRequestDurationByDecision = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: SchedulerSubsystem,
			Name:      "request_duration_by_decision_seconds",
			Help:      metrics.HelpMsgWithStability("Total time taken by the scheduler to process a request, labeled by P/D decision type", compbasemetrics.ALPHA),
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"decision_type"}, // "split" or "combined"
	)
)

// GetCollectors returns all custom collectors for the llm-d-inference-scheduler.
func GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		SchedulerRequestDurationByDecision,
		SchedulerPrefillSelectionDuration,
		SchedulerDecodeSelectionDuration,
		SchedulerPDDecisionCount,
		SchedulerPDThresholdHitsCount,
	}
}

// RecordRequestDurationByDecision records the total time taken by the scheduler to process a request, labeled by P/D decision type.
func RecordRequestDurationByDecision(duration time.Duration, decisionType string) {
	SchedulerRequestDurationByDecision.WithLabelValues(decisionType).Observe(duration.Seconds())
}

// RecordPrefillSelectionDuration records the time taken to select a prefill pod.
func RecordPrefillSelectionDuration(duration time.Duration) {
	SchedulerPrefillSelectionDuration.Observe(duration.Seconds())
}

// RecordDecodeSelectionDuration records the time taken to select a decode pod.
func RecordDecodeSelectionDuration(duration time.Duration) {
	SchedulerDecodeSelectionDuration.Observe(duration.Seconds())
}

// RecordPDDecisionCounter records the type of P/D disaggregation decision made.
func RecordPDDecisionCounter(decisionType string) {
	SchedulerPDDecisionCount.WithLabelValues(decisionType).Inc()
}

// IncPDThresholdHitsCounter increments the counter for P/D token threshold hits.
func IncPDThresholdHitsCounter(threshold int) {
	SchedulerPDThresholdHitsCount.WithLabelValues(strconv.Itoa(threshold)).Inc()
}
