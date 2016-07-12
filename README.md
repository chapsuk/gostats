# GOSTAT

Statsd client wrapper.

## Example

Run
```
docker-compose up -d
```

Open `<DOCKER_HOST>:13000` (`admin:admin`), set `Data Sources` address
as `http://stastd`, test connection. Create Dashboard and see metrics:

```
stats.timers.Gostats.Example.Timing.time.*
stats.Gostats.Example.Counter.count
```
 