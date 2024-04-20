##### required modules 

```
go get cloud.google.com/go/bigquery
go get google.golang.org/api/iterator
go get google.golang.org/api/option
go get github.com/go-sql-driver/mysql
go get github.com/lib/pq
go get github.com/snowflakedb/gosnowflake
```


## Getting Started 

```go
package main 

import (
    "github.com/adarsh-jaiss/library"
)

func main() {
    config := library.Config{
        // document your config here 
    }
    client := library.NewClient(config, library.MySql)

    b, err := client.Tables(config.DatabaseName)
    // Handle error
    // convert byte into []string
    // run a loop over the list of table name 
    for _,v := range data {
        by, err := library.Schema(v)
        // Handle error
        // convert 
    }
}
```

### Running docker-compose to test mysql

- `docker-compose up -d`
- `docker exec -it 86c221092f15 mysql -uroot -p`

### Running docker-compose to test postgresql

- `docker-compose up postgres -d`
- `docker exec -it <container id> psql -U postgres -d postgres`