package main 

import (
    "github.com/adarsh-jaiss/xray"
    "github.com/adarsh-jaiss/xray/types"
)

func main() {
    config := xray.Config{
        Host: "127.0.0.1",
        DatabaseName: "employees",
        Username: "root",
        Port: "5432",
        SSL: "disable"
    }
    client := xray.NewClient(config, types.Postgres)
    data, err := client.Tables(config.DatabaseName)
    if err != nil {
        panic(err)
    } 
    var response = []types.Table
    for _,v := range data {
        table, err := xray.Schema(v)
        if err != nil {
            panic(err)
        } 
        response = response.append(table)
    }
    fmt.Println(response)
}