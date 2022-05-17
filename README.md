# cassandra-builder

**This repo is on development and is incomplete, do not use even for non production projects**
---

This is a simple Query Builder for Apache Cassandra database. It uses `gocql` driver and `goclqx` Query Builder 
internally to create friendly layer to interact with the database trough struct bindings.

## Installing

Simple run next command to download the latest version.

```shell script
go get -u -v github.com/Drafteame/cassandra-builder
```

You can create a client with this:

```go
package main

import "github.com/Drafteame/cassandra-builder/qb"

func main() {
    config := qb.Config{
        Port:          9042,
        KeyspaceName:  "test",
        Username:      "",
        Password:      "",
        ContactPoints: []string{"127.0.0.1"},
    }

    client, err := qb.NewClient(config)
    if err != nil {
        panic(err)
    }

    // Do stuff...    

    client.Close()
}
```


**Note**

This repo is a fork from the github.com/Drafteame/cassandra-builder