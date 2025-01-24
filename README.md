# How to develop post-board

## 1. start the dev dependencies

```sh
make -f make.post_board docker_run_dev_dependencies
```

The above command will start the following services:
- `postgres`
- `rabbitmq`
- `elastic search`
- `kibana`

## 2. run the database migration

```sh
make -f make.post_board local_run_migration_up
```

## 3. start the post-board app & logging system

```sh
make -f make.post_board local_run_app
```

```sh
make -f make.post_board local_run_logging
```

## 4. test the Login API

```sh
curl --location 'http://localhost:8082/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "admin@google.com",
    "password": "secret123"
}'
```

## 5. configure Kibana

1. (Sidebar) Analytics > Discover
2. Create data view
   - Name: `post-board-api`
   - Index Pattern: `post-board-api-*`
   - Timestamp field: `startTime`