
### Sample .env file


```bash
# DNS Server Configuration
DNS_PORT=10053

# Relay Configuration
DNS_RELAY_SERVERS=8.8.8.8:53,1.1.1.1:53

# TTL Configuration (in seconds)
DNS_DEFAULT_TTL=60

# LOGGING Configuration
LOG_DIR="/tmp/log/nanodns"
SERVICE_LOG="service.log"
ACTION_LOG="actions.log"
MAX_LOG_SIZE=1048576 # 1MB
MAX_LOG_BACKUPS=5

# A Records
# Format: domain|ip|ttl or domain|service:servicename|ttl
A_REC1=app1.example.com|10.10.0.1|300
A_REC2=app2.example.com|10.10.0.2|300
# A_REC3=static.example.com|10.0.0.50
# A_REC4=*.example.com|192.168.1.100|300

# CNAME Records
# Format: domain|target|ttl
# CNAME_REC1=www.example.com|app.example.com|3600
# CNAME_REC2=docs.example.com|documentation.service.local
# CNAME_REC3=blog.example.com|app.example.com|600

# MX Records
# Format: domain|priority|mailserver|ttl
# MX_REC1=example.com|10|mail1.example.com|3600
# MX_REC2=example.com|20|mail2.example.com|3600
# MX_REC3=example.com|30|mail-backup.example.com|3600

# TXT Records
# Format: domain|"text value"|ttl
# TXT_REC1=example.com|v=spf1 include:_spf.example.com ~all|3600
# TXT_REC2=_dmarc.example.com|v=DMARC1; p=reject; rua=mailto:dmarc@example.com
# TXT_REC3=_acme-challenge.example.com|validation-token-here|60
# TXT_REC4=mail._domainkey.example.com|v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4|3600

# Service Discovery Examples
# A_REC5=db.local|service:postgres.default.svc.cluster.local|60
# A_REC6=redis.local|service:redis.default.svc.cluster.local|60
# A_REC7=kafka.local|service:kafka.kafka.svc.cluster.local|60

# Load Balancing Examples
# A_REC8=api.local|192.168.1.10|60
# A_REC9=api.local|192.168.1.11|60
# A_REC10=api.local|192.168.1.12|60

# Note: All TTL values are optional and will default to DNS_DEFAULT_TTL if not specified
# Note: Service references will be automatically resolved in Docker/Kubernetes environments

```