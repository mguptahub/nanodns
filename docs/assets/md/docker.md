
### Using Docker Run

```bash
docker run -d \
  --name nanodns \
  -p 15353:15353/udp \
  -e "DNS_PORT=15353" \
  -e "DNS_RELAY_SERVERS=8.8.8.8:53,1.1.1.1:53" \  # Optional relay configuration
  -e "DNS_API_TOKEN=your-secure-token" \
  -e "A_REC1=app.example.com|192.168.1.10|300" \
  -e "TXT_REC1=example.com|v=spf1 include:_spf.example.com ~all" \
  ghcr.io/mguptahub/nanodns:latest
```

### Using Docker Compose

```yaml
name: nanodns
services:
  dns:
    image: ghcr.io/mguptahub/nanodns:latest
    environment:
      - DNS_PORT=15353  # Optional, defaults to 53
      - DNS_RELAY_SERVERS=8.8.8.8:53,1.1.1.1:53
      - DNS_API_TOKEN=your-secure-token
      # A Records
      - A_REC1=app.example.com|service:webapp
      - A_REC2=api.example.com|192.168.1.10|300
      # TXT Records
      - TXT_REC1=example.com|v=spf1 include:_spf.example.com ~all
    ports:
      - "${DNS_PORT:-15353}:${DNS_PORT:-15353}/udp"  # Uses DNS_PORT if set, otherwise 15353
    networks:
      - app_network

networks:
  app_network:
    driver: bridge
```