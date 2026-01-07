# Deployment Guide

## Local Development

### Prerequisites
- Go 1.21+
- PostgreSQL 15+ (optional)
- Make

### Steps

1. Clone and setup:
```bash
git clone https://github.com/lugondev/go-indexer-solana-starter.git
cd go-indexer-solana-starter
cp .env.example .env
```

2. Configure environment:
Edit `.env` file with your settings

3. Run:
```bash
make run
```

## Docker Deployment

### Using Docker Compose (Recommended)

```bash
docker-compose up -d
```

This will start:
- Indexer service on port 8080
- PostgreSQL on port 5432

### Using Docker Only

```bash
docker build -t solana-indexer .
docker run -d \
  --env-file .env \
  -p 8080:8080 \
  --name solana-indexer \
  solana-indexer
```

## Production Deployment

### System Requirements

- CPU: 2+ cores
- RAM: 4GB minimum, 8GB recommended
- Storage: 50GB+ for database
- Network: Stable connection to Solana RPC

### Environment Variables

Ensure these are set:

```bash
SOLANA_RPC_URL=https://api.mainnet-beta.solana.com
SOLANA_WS_URL=wss://api.mainnet-beta.solana.com
DATABASE_URL=postgres://user:pass@host:5432/db
START_SLOT=0
LOG_LEVEL=info
```

### Using systemd

Create `/etc/systemd/system/solana-indexer.service`:

```ini
[Unit]
Description=Solana Indexer
After=network.target postgresql.service

[Service]
Type=simple
User=indexer
WorkingDirectory=/opt/solana-indexer
EnvironmentFile=/opt/solana-indexer/.env
ExecStart=/opt/solana-indexer/indexer
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable solana-indexer
sudo systemctl start solana-indexer
sudo systemctl status solana-indexer
```

### Using Kubernetes

Example deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: solana-indexer
spec:
  replicas: 1
  selector:
    matchLabels:
      app: solana-indexer
  template:
    metadata:
      labels:
        app: solana-indexer
    spec:
      containers:
      - name: indexer
        image: solana-indexer:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: indexer-config
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2000m"
```

## Monitoring

### Health Checks

```bash
curl http://localhost:8080/health
```

### Logs

```bash
journalctl -u solana-indexer -f
```

or with Docker:
```bash
docker logs -f solana-indexer
```

## Backup and Recovery

### Database Backup

```bash
pg_dump -U postgres solana_indexer > backup.sql
```

### Restore

```bash
psql -U postgres solana_indexer < backup.sql
```

## Scaling

### Horizontal Scaling
- Run multiple instances with different START_SLOT ranges
- Use load balancer for API requests

### Vertical Scaling
- Increase MAX_CONCURRENCY
- Increase BATCH_SIZE
- Add more CPU/RAM resources

## Troubleshooting

### High Memory Usage
- Reduce BATCH_SIZE
- Reduce MAX_CONCURRENCY
- Check for memory leaks with pprof

### Slow Indexing
- Increase MAX_CONCURRENCY
- Use faster RPC endpoint
- Optimize database queries
- Add database indexes

### Connection Issues
- Check RPC endpoint status
- Verify network connectivity
- Check firewall rules
