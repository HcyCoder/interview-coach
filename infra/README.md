# Local Infrastructure

Start the local backing services with:

```bash
docker compose -f infra/docker-compose.yml up -d
```

The backend MySQL DSN in `backend/.env.example` matches the MySQL container defaults.
