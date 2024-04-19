##### required modules 

```
go get cloud.google.com/go/bigquery
go get google.golang.org/api/iterator
go get google.golang.org/api/option
go get github.com/go-sql-driver/mysql
go get github.com/lib/pq
go get github.com/snowflakedb/gosnowflake
```


### Running docker-compose to test mysql

- `docker-compose up -d`
- `docker exec -it 86c221092f15 mysql -uroot -p`