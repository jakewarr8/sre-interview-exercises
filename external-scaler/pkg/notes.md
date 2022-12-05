```
kubectl get --raw /apis/external.metrics.k8s.io/v1beta1/namespaces/tenant-a/s0-licenseThreshold | jq
{
  "kind": "ExternalMetricValueList",
  "apiVersion": "external.metrics.k8s.io/v1beta1",
  "metadata": {},
  "items": [
    {
      "metricName": "s0-licenseThreshold",
      "metricLabels": null,
      "timestamp": "2022-12-05T02:30:40Z",
      "value": "2"
    }
  ]
}
```
