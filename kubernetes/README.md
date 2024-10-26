# NanoDNS Kubernetes Deployment Guide

A guide for deploying and managing NanoDNS in Kubernetes environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Configuration Management](#configuration-management)
- [Record Type Examples](#record-type-examples)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

## Prerequisites

- Kubernetes cluster running version 1.19+
- kubectl CLI tool installed and configured
- Access to pull images from GitHub Container Registry (ghcr.io)

## Quick Start

1. **Create the deployment file** (nanodns.yaml):

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nanodns-config
data:
  # DNS Server Configuration
  DNS_PORT: "53"
  # Relay Configuration - for unmatched queries
  DNS_RELAY_SERVERS: "8.8.8.8:53,1.1.1.1:53"  # Comma-separated upstream DNS servers

  # A Records
  A_REC1: "app.example.com|service:frontend.default.svc.cluster.local"
  A_REC2: "api.example.com|service:backend.default.svc.cluster.local"
  A_REC3: "static.example.com|192.168.1.10|300"

  # CNAME Records
  CNAME_REC1: "www.example.com|app.example.com|3600"
  CNAME_REC2: "docs.example.com|documentation.default.svc.cluster.local"

  # MX Records
  MX_REC1: "example.com|10|mail1.example.com|3600"
  MX_REC2: "example.com|20|mail2.example.com"

  # TXT Records
  TXT_REC1: "example.com|v=spf1 include:_spf.google.com ~all|3600"
  TXT_REC2: "_dmarc.example.com|v=DMARC1; p=reject; rua=mailto:dmarc@example.com"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nanodns
  labels:
    app: nanodns
  annotations:
    reloader.stakater.com/auto: "true"  # Auto reload on ConfigMap changes
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

2. **Deploy to Kubernetes**:

```bash
kubectl apply -f nanodns.yaml
```

3. **Verify deployment**:

```bash
kubectl get pods -l app=nanodns
kubectl get svc nanodns
```

## Configuration Management

### DNS Resolution Strategy

NanoDNS follows this resolution strategy:

1. Check local records defined in ConfigMap
2. If no local record found, forward to configured relay servers
3. Return first successful response from relay servers

### Relay Configuration

Configure upstream DNS servers for unmatched queries:

```yaml
data:
  # Single upstream server
  DNS_RELAY_SERVERS: "8.8.8.8:53"

  # Multiple upstream servers (failover)
  DNS_RELAY_SERVERS: "8.8.8.8:53,1.1.1.1:53"

  # Custom port example
  DNS_RELAY_SERVERS: "custom.dns.server:5353,8.8.8.8:53"
```

When using multiple servers:

- Servers are tried in order
- First successful response is used
- 5-second timeout per server

### Updating DNS Records

1. **Edit ConfigMap Directly**

```bash
# Open ConfigMap in editor
kubectl edit configmap nanodns-config

# Or update specific values
kubectl patch configmap nanodns-config --type merge -p '
{
  "data": {
    "A_REC1": "app.example.com|service:frontend.default.svc.cluster.local",
    "A_REC2": "api.example.com|192.168.1.10|300"
  }
}'
```

2. **Apply Changes**

```bash
# Force rollout to pick up changes
kubectl rollout restart deployment/nanodns

# Monitor the rollout
kubectl rollout status deployment/nanodns
```

## Record Type Examples

### A Records

```yaml
# Internal Kubernetes service
A_REC1: "app.example.com|service:frontend.default.svc.cluster.local"

# External IP with TTL
A_REC2: "api.example.com|192.168.1.10|300"

# Simple internal IP
A_REC3: "internal.example.com|10.0.0.50"
```

### CNAME Records

```yaml
# Simple alias
CNAME_REC1: "www.example.com|app.example.com"

# Service alias with TTL
CNAME_REC2: "docs.example.com|documentation.default.svc.cluster.local|3600"
```

### MX Records

```yaml
# Primary mail server
MX_REC1: "example.com|10|mail1.example.com|3600"

# Backup mail server
MX_REC2: "example.com|20|mail2.example.com"
```

### TXT Records

```yaml
# SPF Record
TXT_REC1: "example.com|v=spf1 include:_spf.google.com ~all|3600"

# DMARC Record
TXT_REC2: "_dmarc.example.com|v=DMARC1; p=reject; rua=mailto:dmarc@example.com"

# Verification Record
TXT_REC3: "verification.example.com|verify-domain=example123"
```

## Testing

### Basic DNS Resolution

```bash
# Create a debug pod
kubectl run -it --rm debug --image=alpine/bind-tools -- sh

# Test different record types
dig @nanodns.default.svc.cluster.local app.example.com A
dig @nanodns.default.svc.cluster.local www.example.com CNAME
dig @nanodns.default.svc.cluster.local example.com MX
dig @nanodns.default.svc.cluster.local example.com TXT

# Test relay to upstream (for non-local domain)
kubectl run -it --rm debug --image=alpine/bind-tools -- \
  dig @nanodns.default.svc.cluster.local google.com
```

### Service Resolution

```bash
# Test internal service resolution
kubectl run -it --rm debug --image=alpine/bind-tools -- dig @nanodns.default.svc.cluster.local app.example.com

# Verify CNAME resolution
kubectl run -it --rm debug --image=alpine/bind-tools -- dig @nanodns.default.svc.cluster.local www.example.com CNAME +short
```

## Troubleshooting


### DNS Relay Issues

1. **Relay Not Working**

```bash
# Check relay configuration
kubectl describe configmap nanodns-config | grep DNS_RELAY

# Test upstream connectivity
kubectl run -it --rm debug --image=alpine/bind-tools -- \
  dig @8.8.8.8 google.com

# Check DNS server logs
kubectl logs -l app=nanodns -f
```

2. **Slow Resolution**

```bash
# Test resolution time
kubectl run -it --rm debug --image=alpine/bind-tools -- \
  dig @nanodns.default.svc.cluster.local google.com +stats

# If slow, check network policies
kubectl get networkpolicies
```

### Common Issues and Solutions

1. **Pod Won't Start**

```bash
# Check pod status
kubectl get pods -l app=nanodns

# Check pod events
kubectl describe pod -l app=nanodns

# Check logs
kubectl logs -l app=nanodns
```

2. **DNS Resolution Not Working**

```bash
# Verify ConfigMap
kubectl get configmap nanodns-config -o yaml

# Check service endpoints
kubectl get endpoints nanodns

# Test from debug pod
kubectl run -it --rm debug --image=busybox -- nslookup app.example.com nanodns.default.svc.cluster.local
```

3. **Configuration Updates Not Applied**

```bash
# Check ConfigMap changes
kubectl describe configmap nanodns-config

# Force pod restart
kubectl rollout restart deployment/nanodns

# Monitor rollout
kubectl rollout status deployment/nanodns
```

## Best Practices

### DNS Configuration

- Configure reliable upstream DNS servers
- Use multiple relay servers for redundancy
- Consider geographic proximity for relay servers
- Monitor relay server response times
- Document relay server selection criteria

### Architecture Considerations

- **Local Records**:
  - Use for internal services
  - Kubernetes service discovery
  - Custom DNS mappings
  
- **Relay DNS**:
  - External domain resolution
  - Internet access required
  - Fallback resolution
  - Upstream server selection

### Resource Management

- Configure appropriate resource requests and limits
- Monitor resource usage
- Scale replicas based on load

### Security

- Keep the image updated
- Run as non-root user
- Use read-only root filesystem
- Implement network policies if needed

### High Availability

- Use multiple replicas in production
- Configure proper health checks
- Implement proper monitoring

### Monitoring

- Watch pod logs for errors
- Monitor DNS query latency
- Track resource utilization
- Set up alerts for failures

### Configuration

- Regularly backup ConfigMap
- Document all DNS records
- Use meaningful TTL values
- Keep records organized by type

Remember to replace example domains and IPs with your actual values when deploying.