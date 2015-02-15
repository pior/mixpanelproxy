# MixpanelProxy

Mixpanel Proxy library in Go (as a http.Handler)

## Features

- Extract the Token and the DistinctId
- Support Events and People endpoints
- Passthrough for other endpoints
- Custom Director to implement the proxy logic
- Fully functional dummy (successful) response

## TODO

- Test coverage

## Usage

```go
package main

import (
    "github.com/cenkalti/log"
    "github.com/pior/mixpanelproxy"
    "net/http"
)

func main() {
    director := func(m *mixpanelproxy.MixpanelRequest) (err error) {
        log.Infof("Received: %s", m)
        return nil
    }

    url := "http://api.mixpanel.com"

    http.Handle("/", mixpanelproxy.NewProxy(&url, director))

    log.Fatal(http.ListenAndServe(":8080", nil))
}

```

## Logging

MixpanelProxy use [github.com/cenkalti/log](https://github.com/cenkalti/log) for logging
