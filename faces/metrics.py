from prometheus_client import Counter, Histogram

request_counter = Counter(
    'request_count', 'times an endpoint was called', labelnames=['endpoint'])

request_latency_histogram = Histogram('request_latency_seconds', 'the time it takes for an endpoint to respond',
                                      labelnames=['endpoint'], buckets=(.01, .05, .1, .5, 1.0, 2.0, 4.0, 8.0, 10.0))
