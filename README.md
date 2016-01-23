## Supported devices:

* Diamex GmbH / AVR Temp Sensor (USB `16c0:0480`)

## InfluxDB

Supports writing stats to `server -influxdb-server=... -influxdb-database=...`.

![Grafana Screenshot](/doc/grafana.png?raw=true "Grafana")

Grafana Query:

    SELECT mean("temperature") FROM "onewire" WHERE "family" = 'ds18b20' AND $timeFilter GROUP BY time($interval), "id", "name" fill(null)

## Dependencies

* libudev-dev

## Configuration

### `/etc/udev/rules.d/90-hidraw.rules`

    KERNEL=="hidraw*", \
        GROUP="plugdev", MODE=0660

