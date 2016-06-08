# GOSTAT

Statsd client wrapper.

## Usage

```go
func main() {
    host := "192.168.99.100:8125"
    stat, err := gostats.NewStatsd(host, "global.prefix")
    if err != nil {
        return
    }
    stat.Write("metric.fzz", 1 * time.Millisecond)
    stat.Write("metric.gzz", 2 * time.Millisecond)
    stat.Write("metric.bzz", 3 * time.Millisecond)
}
```

## Example

See [exmaple](/examples/main.go).

1. Start graphit and statsd with docker compose: `docker-comopse up -d`
1. Run example: `go run example/main.go -h <DOCKER_HOST>`.
1. Open grafana (`<DOCKER_HOST:3000`). Set data source as `<DOCKER_HOST>:8080`.
Create Dashboard. See metrics `stats.*`.
1. Stop example
