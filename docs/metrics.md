# Metrics

The `llm-d-inference-scheduler` exposes the following Prometheus metrics to monitor its behavior and performance, particularly concerning Prefill/Decode Disaggregation.

All metrics are in the `llm_d_inference_scheduler` subsystem.

## Scrape and see the metric

Metrics defined in the scheduler plugin are extention of Inference Gateway metrics. For more details of seeing metrics, see the [Instruction](https://github.com/kubernetes-sigs/gateway-api-inference-extension/blob/main/site-src/guides/metrics-and-observability.md).

## Metric Details

### `pd_decision_total`

*   **Type:** Counter
*   **Labels:**
    *   `decision_type`: string ("split" or "combined")
*   **Release Stage:** ALPHA
*   **Description:** Counts the number of requests processed, broken down by the Prefill/Decode disaggregation decision.
    *   `split`: The request was split into separate Prefill and Decode stages.
    *   `combined`: The request used the Decode-only path.
*   **Usage:** Provides a high-level view of how many requests are utilizing the disaggregated path versus the unified path.
*   **Actionability:**
    *   Monitor the ratio of "split" to "combined" to understand the P/D engagement rate.
    *   Sudden changes in this ratio might indicate configuration issues, changes in workload patterns, or problems with the decision logic.

---

### `pd_threshold_hits_total`

*   **Type:** Counter
*   **Labels:**
    *   `threshold`: string (integer value of the configured threshold)
*   **Release Stage:** ALPHA
*   **Description:** Counts how many times the "new token count < threshold" condition was met in the `pd-profile-handler`, triggering the "combined" (Decode-only) path.
*   **Usage:** Helps understand the effectiveness of the `threshold` parameter in the `pd-profile-handler` configuration.
*   **Actionability:**
    *   Analyze the hit counts for different `threshold` values to tune the parameter for optimal performance.
    *   If the count is unexpectedly low or high for a given threshold, it may indicate issues with token counting or the workload.

---

### `request_duration_by_decision_seconds`

*   **Type:** Histogram
*   **Labels:**
    *   `decision_type`: string ("split" or "combined")
*   **Release Stage:** ALPHA
*   **Description:** Measures the total end-to-end request processing duration by the scheduler, from the start of the `ProcessResults` function in the `pd-profile-handler` to its completion, labeled by the P/D decision type.
*   **Usage:** Critical for evaluating the latency impact of P/D disaggregation.
*   **Actionability:**
    *   Compare the latency distributions for "split" and "combined" requests to assess whether P/D is yielding performance benefits.
    *   Track percentiles (e.g., p95, p99) to understand tail latency for both paths.
    *   Increases in latency can signal bottlenecks in either the scheduler or the downstream model servers.

---

### `prefill_selection_duration_seconds`

*   **Type:** Histogram
*   **Release Stage:** ALPHA
*   **Description:** Measures the time taken by the `PdProfileHandler.Pick` function to decide to run the "prefill" profile. This includes logic like token counting and threshold checks.
*   **Usage:** Monitors the overhead introduced by the scheduler's decision logic for the prefill path.
*   **Actionability:**
    *   High values indicate potential performance issues within the `Pick` function's prefill selection logic (e.g., slow state access, expensive computations).
    *   Helps distinguish scheduler-induced latency from model server latency.

---

### `decode_selection_duration_seconds`

*   **Type:** Histogram
*   **Release Stage:** ALPHA
*   **Description:** Measures the time taken by the `PdProfileHandler.Pick` function to decide to run the "decode" profile. This is typically the initial decision point.
*   **Usage:** Monitors the overhead of the most frequent decision path in the scheduler.
*   **Actionability:**
    *   While expected to be very low, any increase could signal a performance regression in the core `Pick` logic.
