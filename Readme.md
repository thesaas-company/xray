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

TDOD: 
- Please rename the project, You can choose any name that you like
- Below is a example for users 

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
    data, err := client.Tables(config.DatabaseName)
    // Handle error 
    var response = []library.Table
    for _,v := range data {
        table, err := library.Schema(v)
        // Handle error
        response = response.append(table)
    }
    fmt.Println(response)
}
```

## Testing

### Running docker-compose to test mysql

- `docker-compose up -d`
- `docker exec -it 86c221092f15 mysql -uroot -p`

### Running docker-compose to test postgresql

- `docker-compose up postgres -d`
- `docker exec -it <container id> psql -U postgres -d postgres`