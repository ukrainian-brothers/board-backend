## Setting up tests
Most of the unit tests can be ran without any external configuration, anyway the integration tests will require database setted up from `sql/create_tables.sql` and config which is not included in source files.

The config should be located in `config/configuration.test.local.json` and look like this:
```json
{
    "postgres_config": {
        "host": "localhost",
        "port": 5438,
        "user": "",
        "password": "",
        "db_name": "test_ukrainian_brothers"
    }
}
```
