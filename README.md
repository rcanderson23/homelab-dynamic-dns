# homelab-dynamic-dns
Goal is to provide automatic DNS updates in a homelab setting where your external IP is not known
to the cluster and while also providing updates to internal IPs. This is done through a combination
IPLookup and DNSProvider Custom Resources and Annotations on the Ingresses that should have DNS entries
managed.

Current functionality is limited to Ingress objects with Services to be added. 

## Example
`IPLookup` performs a GET on an HTTP endpoint and updates the `status` of the returned IP address
```yaml
apiVersion: networking.thehomelab.tech/v1alpha1
kind: IPLookup
metadata:
  name: aws
spec:
  type: http
  config:
    http:
      url: https://checkip.amazonaws.com
```

`DNSProvider` is your DNS provider that you are using to update records.
```yaml
apiVersion: networking.thehomelab.tech/v1alpha1
kind: DNSProvider
metadata:
  name: thehomelab-tech-aws
spec:
  type: aws
  config:
    aws:
      hostZoneID: AAABBBCCC123
      ttl: 300
```

`Ingress` requires 3 annotations to be watched by this controller:
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: thehomelab-tech
  annotations:
    thehomelab.tech/dnsprovider: aws # name of the DNSProvider CR to use
    thehomelab.tech/ip-address: external # Public IP; use 'internal' for ingress IP
    thehomelab.tech/iplookup: aws # name of IPLookup CR to use
spec:
  rules:
  - host: www.thehomelab.tech
    http:
      paths:
      - backend:
          service:
            name: web
            port:
              number: 8080
        path: /
        pathType: ImplementationSpecific
  tls:
  - hosts:
    - www.thehomelab.tech
    secretName: thehomelab-tls
```