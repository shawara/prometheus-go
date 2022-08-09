# Prometheus Go Demo
This example app serves as an example of how one can easily instrument HTTP handlers with Prometheus metrics. It uses the Prometheus go client to create a new Prometheus registry.

A [Docker image](https://hub.docker.com/r/shawara/prometheus-go) is available at: `shawara/prometheus-go`

Prometheus Metrics Route =>  `/metrics`

Delay Route to simulate real apis  =>  `/delay/{millis}`

### Using it inside kubernetes
1. install Prometheus Operator 
    ```shell
     helm install prometheus bitnami/kube-prometheus
     ```
2. Create deployment for the app
    ```shell
    kubectl create -f ./k8s/deployment.yaml
    ```
3. User CRD `PodMonitor` to all prometheus to fetch metrics from app
    ```shell
    kubectl create -f ./k8s/pod-monitor.yaml
    ```
   
### Metrics being tracked
   ```cmd
   # HELP http_requests_total Number of get requests.
   # TYPE http_requests_total counter
   http_requests_total{code="200",path="/"} 2
   # HELP http_response_time_seconds Duration of HTTP requests.
   # TYPE http_response_time_seconds histogram
   http_response_time_seconds_bucket{code="200",path="/",le="0.005"} 2
   http_response_time_seconds_bucket{code="200",path="/",le="0.01"} 2
   http_response_time_seconds_bucket{code="200",path="/",le="0.025"} 2
   http_response_time_seconds_bucket{code="200",path="/",le="0.05"} 2
   http_response_time_seconds_bucket{code="200",path="/",le="0.1"} 2
   http_response_time_seconds_bucket{code="200",path="/",le="0.25"} 2
   http_response_time_seconds_bucket{code="200",path="/",le="0.5"} 2
   http_response_time_seconds_bucket{code="200",path="/",le="1"} 2
   http_response_time_seconds_bucket{code="200",path="/",le="2.5"} 2
   http_response_time_seconds_bucket{code="200",path="/",le="5"} 2
   http_response_time_seconds_bucket{code="200",path="/",le="10"} 2
   http_response_time_seconds_bucket{code="200",path="/",le="+Inf"} 2
   http_response_time_seconds_sum{code="200",path="/"} 4.324e-06
   http_response_time_seconds_count{code="200",path="/"} 2
   ```