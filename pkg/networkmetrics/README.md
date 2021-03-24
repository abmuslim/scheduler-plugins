## Prometheus Query

```shell
curl -g 'http://localhost:9090/api/v1/query?query=sum_over_time(node_network_receive_bytes_total{kubernetes_node="node53",device="ens192"}[5m])' | jq
```