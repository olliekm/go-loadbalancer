# go-loadbalancer
simple load balancer in go

## todos:
1. Improve Load Balancing Strategies

    Your current implementation uses Round Robin, but adding more sophisticated strategies would make it more efficient.
    
    Implement Weighted Load Balancing
    Assign weights to backends so stronger servers get more requests.
    Dynamically adjust weights based on health metrics (e.g., response time).

2. Improve Observability & Metrics

   Current issue: No visibility into traffic, latency, or failures.
    Use Prometheus for Metrics

3. Rate Limiting to Prevent Overload
   
   If your load balancer is receiving too much traffic, rate limit requests.
   Use golang.org/x/time/rate