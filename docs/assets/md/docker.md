
### Using Docker Run

```bash
docker run -d \
  --name nanodns \
  -p 10053:10053/udp \
  -e DNS_PORT=10053 \
  -e DNS_RELAY_SERVERS=8.8.8.8:53,1.1.1.1:53 \
  -e "A_REC1=app.example.com|10.10.0.1|300" \
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
      - DNS_PORT=10053  # Optional, defaults to 53
      - DNS_RELAY_SERVERS=8.8.8.8:53,1.1.1.1:53
      # A Records
      - A_REC1=app.example.com|service:webapp
      - A_REC2=api.example.com|10.10.0.10|300
      # TXT Records
      - TXT_REC1=example.com|v=spf1 include:_spf.example.com ~all
    ports:
      - "${DNS_PORT:-10053}:${DNS_PORT:-10053}/udp"  # Uses DNS_PORT if set, otherwise 10053
    networks:
      - app_network

networks:
  app_network:
    driver: bridge
```