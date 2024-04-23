# XRay

## Install 

```bash
go get github.com/thesaas-company/xray@latest
```

## Example

### MySQL

```bash
$ docker run -d \                                                                                                         
  --name mysql-employees \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=college \
  -v $PWD/data:/var/lib/mysql \
  genschsa/mysql-employees

$ export DB_PASSWORD=college
```

```go
package main 

import (
    "github.com/thesaas-company/xray"
)

func main() {
    config := library.Config{
        Host: "127.0.0.1",
        DatabaseName: "employees",
        Username: "root",
        Port: "3306",
        SSL: "false"
    }
    client := library.NewClient(config, xray.MySql)
    data, err := client.Tables(config.DatabaseName)
    if err != nil {
        panic(err)
    } 
    var response = []library.Table
    for _,v := range data {
        table, err := library.Schema(v)
        if err != nil {
            panic(err)
        } 
        response = response.append(table)
    }
    fmt.Println(response)
}
```

### Postgres

```bash
$ TBD

$ export DB_PASSWORD=college
```

```go
package main 

import (
    "github.com/thesaas-company/xray"
)

func main() {
    config := library.Config{
        Host: "127.0.0.1",
        DatabaseName: "employees",
        Username: "root",
        Port: "3306",
        SSL: "disable"
    }
    client := library.NewClient(config, xray.Postgres)
    data, err := client.Tables(config.DatabaseName)
    if err != nil {
        panic(err)
    } 
    var response = []library.Table
    for _,v := range data {
        table, err := library.Schema(v)
        if err != nil {
            panic(err)
        } 
        response = response.append(table)
    }
    fmt.Println(response)
}
```