# Disaster Recovery System Design

- **Primary Cluster**: production apps
- **Secondary Cluster**: standby DRC
- **MySQL** replication (async / semi-sync)
- **Vault** for secret sync
- **Monitoring** with Grafana + Loki
- **Failover**: DNS switch via Ingress
