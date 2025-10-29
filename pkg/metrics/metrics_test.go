package metrics

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSchedulerDecodeSelectionDuration(t *testing.T) {
	RecordDecodeSelectionDuration(150 * time.Millisecond)
	expectedSum := fmt.Sprintf("%.2f", 0.15)
	if err := testutil.CollectAndCompare(SchedulerDecodeSelectionDuration, strings.NewReader(fmt.Sprintf(`
		# HELP llm_d_inference_scheduler_decode_selection_duration_seconds [ALPHA] Time taken to select a decode pod
		# TYPE llm_d_inference_scheduler_decode_selection_duration_seconds histogram
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="0.005"} 0
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="0.01"} 0
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="0.025"} 0
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="0.05"} 0
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="0.1"} 0
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="0.25"} 1
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="0.5"} 1
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="1"} 1
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="2.5"} 1
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="5"} 1
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="10"} 1
		llm_d_inference_scheduler_decode_selection_duration_seconds_bucket{le="+Inf"} 1
		llm_d_inference_scheduler_decode_selection_duration_seconds_sum %s
		llm_d_inference_scheduler_decode_selection_duration_seconds_count 1
	`, expectedSum))); err != nil {
		t.Errorf("RecordDecodeSelectionDuration() failed: %v", err)
	}
}

func TestSchedulerPDDecisionCount(t *testing.T) {
	RecordPDDecisionCounter("split")
	RecordPDDecisionCounter("combined")
	RecordPDDecisionCounter("split")
	if err := testutil.CollectAndCompare(SchedulerPDDecisionCount, strings.NewReader(`
		# HELP llm_d_inference_scheduler_pd_decision_total [ALPHA] Total number of P/D disaggregation decisions made
		# TYPE llm_d_inference_scheduler_pd_decision_total counter
		llm_d_inference_scheduler_pd_decision_total{decision_type="combined"} 1
		llm_d_inference_scheduler_pd_decision_total{decision_type="split"} 2
	`), "decision_type"); err != nil {
		t.Errorf("RecordPDDecisionCounter() failed: %v", err)
	}
}

func TestSchedulerPDThresholdHitsCount(t *testing.T) {
	IncPDThresholdHitsCounter(10)
	IncPDThresholdHitsCounter(20)
	IncPDThresholdHitsCounter(10)
	if err := testutil.CollectAndCompare(SchedulerPDThresholdHitsCount, strings.NewReader(`
		# HELP llm_d_inference_scheduler_pd_threshold_hits_total [ALPHA] Total number of times the P/D token threshold was met, labeled by the configured threshold
		# TYPE llm_d_inference_scheduler_pd_threshold_hits_total counter
		llm_d_inference_scheduler_pd_threshold_hits_total{threshold="10"} 2
		llm_d_inference_scheduler_pd_threshold_hits_total{threshold="20"} 1
	`), "threshold"); err != nil {
		t.Errorf("IncPDThresholdHitsCounter() failed: %v", err)
	}
}

func TestSchedulerRequestDurationByDecision(t *testing.T) {
	RecordRequestDurationByDecision(100*time.Millisecond, "combined")
	RecordRequestDurationByDecision(250*time.Millisecond, "split")
	expectedSumCombined := fmt.Sprintf("%.1f", 0.1)
	expectedSumSplit := fmt.Sprintf("%.2f", 0.25)
	if err := testutil.CollectAndCompare(SchedulerRequestDurationByDecision, strings.NewReader(fmt.Sprintf(`
		# HELP llm_d_inference_scheduler_request_duration_by_decision_seconds [ALPHA] Total time taken by the scheduler to process a request, labeled by P/D decision type
		# TYPE llm_d_inference_scheduler_request_duration_by_decision_seconds histogram
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="0.005"} 0
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="0.01"} 0
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="0.025"} 0
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="0.05"} 0
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="0.1"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="0.25"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="0.5"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="1"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="2.5"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="5"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="10"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="combined",le="+Inf"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_sum{decision_type="combined"} %s
		llm_d_inference_scheduler_request_duration_by_decision_seconds_count{decision_type="combined"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="0.005"} 0
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="0.01"} 0
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="0.025"} 0
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="0.05"} 0
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="0.1"} 0
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="0.25"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="0.5"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="1"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="2.5"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="5"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="10"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_bucket{decision_type="split",le="+Inf"} 1
		llm_d_inference_scheduler_request_duration_by_decision_seconds_sum{decision_type="split"} %s
		llm_d_inference_scheduler_request_duration_by_decision_seconds_count{decision_type="split"} 1
	`, expectedSumCombined, expectedSumSplit)), "decision_type"); err != nil {
		t.Errorf("RecordRequestDurationByDecision() failed: %v", err)
	}
}

func TestSchedulerPrefillSelectionDuration(t *testing.T) {
	RecordPrefillSelectionDuration(120 * time.Millisecond)
	expectedSum := fmt.Sprintf("%.2f", 0.12)
	if err := testutil.CollectAndCompare(SchedulerPrefillSelectionDuration, strings.NewReader(fmt.Sprintf(`
		# HELP llm_d_inference_scheduler_prefill_selection_duration_seconds [ALPHA] Time taken to select a prefill pod
		# TYPE llm_d_inference_scheduler_prefill_selection_duration_seconds histogram
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="0.005"} 0
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="0.01"} 0
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="0.025"} 0
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="0.05"} 0
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="0.1"} 0
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="0.25"} 1
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="0.5"} 1
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="1"} 1
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="2.5"} 1
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="5"} 1
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="10"} 1
		llm_d_inference_scheduler_prefill_selection_duration_seconds_bucket{le="+Inf"} 1
		llm_d_inference_scheduler_prefill_selection_duration_seconds_sum %s
		llm_d_inference_scheduler_prefill_selection_duration_seconds_count 1
	`, expectedSum))); err != nil {
		t.Errorf("RecordPrefillSelectionDuration() failed: %v", err)
	}
}