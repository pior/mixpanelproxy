# MixpanelProxy

Mixpanel Proxy library in Go


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

MixpanelProxy use github.com/cenkalti/log for logging
