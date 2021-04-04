## Prometheus Query

```shell
curl -g 'http://localhost:9090/api/v1/query?query=sum_over_time(node_network_receive_bytes_total{kubernetes_node="node53",device="ens192"}[5m])' | jq
```

```shell
# Query to retrieve highest metric which will be used as the 100% value
topk(1,sum_over_time(node_network_receive_bytes_total{device="ens192"}[5m]))
```