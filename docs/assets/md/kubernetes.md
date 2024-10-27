### Deploying in Kubernetes

1. **Create the deployment file** (nanodns.yaml):
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nanodns-config
data:
  # DNS Server Configuration
  DNS_PORT: "53"
  DNS_RELAY_SERVERS: 8.8.8.8:53,1.1.1.1:53

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
