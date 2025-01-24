# How to develop post-board

## 1. Start dev dependencies

```sh
make -f make.post_board docker_run_dev_dependencies
```

The above command will start the following services:
- [PostgreSQL](https://www.postgresql.org/)
- [RabbitMQ](https://www.rabbitmq.com/)
- [Elasticsearch](https://www.elastic.co/elasticsearch)
- [Kibana](https://www.elastic.co/kibana)

## 2. Run the database migration

```sh
make -f make.post_board local_run_migration_up
```

## 3. Start post-board

Start post-board app
```sh
make -f make.post_board local_run_app
```

Start logging system
```sh
make -f make.post_board local_run_logging
```

## 4. Test Login API

```sh
curl --location 'http://localhost:8082/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "admin@google.com",
    "password": "secret123"
}'
# Expected Result: {"data":{"token":"this_suppose_to_be_a_long_jwt_token_string"}}
```

## 5. Configure Kibana

1. Go to `http://localhost:5601`
2. (Sidebar) Analytics > Discover
3. Create data view
   - Name: `post-board-api`
   - Index Pattern: `post-board-api-*`
   - Timestamp field: `startTime`
   - Click `Save data view to Kibana`