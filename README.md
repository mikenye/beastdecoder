# beastdecoder

This is my effort to create a BEAST protocol decoder, that will decode ADS-B/Mode-S data in BEAST protocol frames.

**This is a work in progress, and has bugs! Use at your own risk.**

## Running

```
# go mod tidy
# go run ./... --connect beasthost:30005 -lat -33.33333 -lon 111.11111 --webview 0.0.0.0:8888
```

Replacing:

* `beasthost:30005` to a host/port that provides BEAST data.
* Lat/Long of your receiver.
* IP:Port to listen on for `webviw`

Once started, connect to the webview IP/port to see the vessels being tracked.
