# Babex

Babex allows you to make chain of microservices on the fly with the help of RabbitMQ.

## Usage

For example we create service which will add number to counter:

First, create service:

```go
package main

import (
	"encoding/json"
	"log"

	"github.com/matroskin13/babex"
)

func main() {
	service, err := babex.NewService(&babex.ServiceConfig{
	    // rabbit addr
		Address:  "amqp://guest:guest@localhost:5672/",
		// name of service (name of queue)
		Name:     "inc-service",
	})
	if err != nil {
		log.Fatal(err)
	}

	err = service.BindToExchange("example", "inc")
	if err != nil {
		log.Fatal(err)
	}
}
```

Consume messages:

```go
msgs, err := service.GetMessages()
if err != nil {
    log.Fatal(err)
}

errChan := service.GetErrors()

for {
    select {
    case msg := <-msgs:
        data := struct {
            Count int `json:"count"`
        }{}

        if err := json.Unmarshal(msg.Data, &data); err != nil {
            log.Println(err)
            break
        }

        cfg := struct {
            IncStep int `json:"incStep"`
        }{}

        if err := json.Unmarshal(msg.Config, &cfg); err != nil {
            log.Println(err)
            break
        }

        data.Count += cfg.IncStep

        log.Printf("count = %v, incStep = %v \r\n", data.Count, cfg.IncStep)

        service.Next(msg, data, nil)
    case err := <-errChan:
        log.Fatal("err", err)
    }
}
```

And publish message to example/inc:

```json
{
  "data": {
    "count": 0
  },
  "config": {
    "incStep": 2
  },
  "chain": [
    {
      "exchange": "example",
      "key": "inc"
    }
  ]
}
```

Check logs:

```bash
$ 2018/06/13 13:51:35 count = 2, incStep = 2
```

Excellent! Let's change the message:

```json
{
  "data": {
    "count": 0
  },
  "config": {
    "incStep": 2
  },
  "chain": [
    {
      "exchange": "example",
      "key": "inc"
    },
    {
      "exchange": "example",
      "key": "inc"
    }
  ]
}
```

Check logs:

```bash
$ 2018/06/13 13:51:35 count = 2, incStep = 2
$ 2018/06/13 13:51:35 count = 4, incStep = 2
```
