# coredns-sql
Simplified version of `https://github.com/wenerme/coredns-pdsql/`

 * Case insensitive searches added
 * Wildcard searches removed

# Usage
plugin.cfg
```
sql:git.dev.a1fred.com/protocloud/zibort-coredns/pdsql
sql_postgres:github.com/jinzhu/gorm/dialects/postgres
```

Corefile
```
sql postgres "host=postgres user=postgres password=postgres dbname=dns port=5432 sslmode=disable" {
    # debug db  # Uncomment to log sql queries
    auto-migrate
}
```
