# NanoDNS Kubernetes Deployment Guide

This guide explains how to deploy NanoDNS in a Kubernetes environment.

## Table of Contents
- [Quick Start](#quick-start)
- [Configuration Examples](#configuration-examples)
- [Testing DNS Records](#testing-dns-records)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

## Quick Start

Create a file named `nanodns.yaml`:

```yaml
# nanodns.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nanodns-config
data:
  # DNS Server Configuration
  DNS_PORT: "53"

  # A Records
  A_REC1: "app.example.com|service:frontend.default.svc.cluster.local"
  A_REC2: "api.example.com|service:backend.default.svc.cluster.local"
  A_REC3: "static.example.com|192.168.1.10|300"

  # CNAME Records
  CNAME_REC1: "www.example.com|app.example.com|3600"
  CNAME_REC2: "docs.example.com|documentation.default.svc.cluster.local"
  CNAME_REC3: "portal.example.com|app.example.com"

  # MX Records
  MX_REC1: "example.com|10|mail1.example.com|3600"
  MX_REC2: "example.com|20|mail2.example.com|3600"
  MX_REC3: "example.com|30|mailbackup.example.com"

  # TXT Records
  TXT_REC1: "example.com|v=spf1 include:_spf.google.com ~all|3600"
  TXT_REC2: "_dmarc.example.com|v=DMARC1; p=reject; rua=mailto:dmarc@example.com"
  TXT_REC3: "verification.example.com|google-site-verification=example123"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nanodns
  labels:
    app: nanodns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nanodns
  template:
    metadata:
      labels:
        app: nanodns
    spec:
      containers:
      - name: nanodns
        image: ghcr.io/mguptahub/nanodns:latest
        ports:
        - containerPort: 53
          protocol: UDP
        envFrom:
        - configMapRef:
            name: nanodns-config
        resources:
          limits:
            memory: "128Mi"
            cpu: "100m"
          requests:
            memory: "64Mi"
            cpu: "50m"
        securityContext:
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          runAsUser: 1000
        livenessProbe:
          exec:
            command:
            - dig
            - "@127.0.0.1"
            - "app.example.com"
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          exec:
            command:
            - dig
            - "@127.0.0.1"
            - "app.example.com"
          initialDelaySeconds: 5
          periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: nanodns
spec:
  selector:
    app: nanodns
  ports:
  - port: 53
    protocol: UDP
    targetPort: 53
  type: ClusterIP
```

Deploy:
```bash
kubectl apply -f nanodns.yaml
```

## Configuration Examples

### Service Resolution (A Records)
```yaml
# Internal service
A_REC1: "app.example.com|service:frontend.default.svc.cluster.local"

# External IP
A_REC2: "api.example.com|192.168.1.10|300"
```

### Domain Aliases (CNAME Records)
```yaml
# Simple alias
CNAME_REC1: "www.example.com|app.example.com"

# Service alias
CNAME_REC2: "docs.example.com|documentation.default.svc.cluster.local"
```

### Mail Configuration (MX Records)
```yaml
# Primary mail server
MX_REC1: "example.com|10|mail1.example.com|3600"

# Backup mail server
MX_REC2: "example.com|20|mail2.example.com|3600"
```

### Domain Verification (TXT Records)
```yaml
# SPF Record
TXT_REC1: "example.com|v=spf1 include:_spf.google.com ~all|3600"

# DMARC Record
TXT_REC2: "_dmarc.example.com|v=DMARC1; p=reject; rua=mailto:dmarc@example.com"

# Domain Verification
TXT_REC3: "verification.example.com|google-site-verification=example123"
```

## Testing DNS Records

1. **A Records**
```bash
# Test service resolution
kubectl run -it --rm debug --image=alpine/bind-tools -- dig @nanodns.default.svc.cluster.local app.example.com A

# Test static IP
kubectl run -it --rm debug --image=alpine/bind-tools -- dig @nanodns.default.svc.cluster.local static.example.com A
```

2. **CNAME Records**
```bash
# Test CNAME resolution
kubectl run -it --rm debug --image=alpine/bind-tools -- dig @nanodns.default.svc.cluster.local www.example.com CNAME
```

3. **MX Records**
```bash
# Test MX record resolution
kubectl run -it --rm debug --image=alpine/bind-tools -- dig @nanodns.default.svc.cluster.local example.com MX
```

4. **TXT Records**
```bash
# Test SPF record
kubectl run -it --rm debug --image=alpine/bind-tools -- dig @nanodns.default.svc.cluster.local example.com TXT
```

## Troubleshooting

### Check Pod Status
```bash
# Get pod status
kubectl get pods -l app=nanodns

# Get detailed pod information
kubectl describe pod -l app=nanodns
```

### View Logs
```bash
# View pod logs
kubectl logs -l app=nanodns -f

# View previous pod logs if pod was restarted
kubectl logs -l app=nanodns -p
```

### Check Service
```bash
# Get service details
kubectl get svc nanodns

# Describe service
kubectl describe svc nanodns
```

### Common Issues and Solutions

1. **Pod won't start**
   - Check pod events: `kubectl describe pod -l app=nanodns`
   - Check logs: `kubectl logs -l app=nanodns`

2. **DNS resolution not working**
   - Verify ConfigMap: `kubectl get configmap nanodns-config -o yaml`
   - Check service endpoints: `kubectl get endpoints nanodns`

3. **Port conflicts**
   - Check port usage: `kubectl get svc --all-namespaces | grep 53`

## Best Practices

1. **DNS Resolution within Cluster**
   - Use full service DNS names: `service.namespace.svc.cluster.local`
   - Configure appropriate TTL values
   - Consider namespace isolation

2. **Resource Management**
   - Adjust resource limits based on usage
   - Monitor memory consumption
   - Set appropriate replicas

3. **Security**
   - Keep image updated to latest version
   - Use non-root user
   - Implement network policies if needed

4. **Monitoring**
   - Watch pod logs for errors
   - Monitor DNS query latency
   - Check resource utilization