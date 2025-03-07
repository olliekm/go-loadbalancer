![alt](https://64.media.tumblr.com/3c1c16ee3b29662d997f432cac320ef1/tumblr_nj3kbxu32l1u78x0oo1_500.gif)
# PerfectlyBalanced
Building out robust load balancer

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


