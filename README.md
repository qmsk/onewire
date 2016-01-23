## Usage

        $ go build -o bin/server server.go 
        $ ./bin/server -config-path config/test.toml -influxdb-server=localhost
        2016/01/23 13:41:14 server.LoadConfig config/test.toml
        2016/01/23 13:41:14 server.InfluxWriter {localhost onewire}
        2016/01/23 13:41:14 hidraw.List...
        2016/01/23 13:41:14 hidraw.Monitor...
        2016/01/23 13:41:14 http.ListenAndServe :8283...
        2016/01/23 13:41:21 AddHidrawDevice hidraw.DeviceInfo{Device:libudev.Device{DevPath:"/devices/pci0000:00/0000:00:14.0/usb1/1-2/1-2:1.0/0003:16C0:0480.006D/hidraw/hidraw0", Subsystem:"hidraw", DevType:"", SysPath:"/sys/devices/pci0000:00/0000:00:14.0/usb1/1-2/1-2:1.0/0003:16C0:0480.006D/hidraw/hidraw0", SysName:"hidraw0", SysNum:"0", DevNode:"/dev/hidraw0"}, DeviceConfig:hidraw.DeviceConfig{VendorID:0x16c0, ProductID:0x480}}: &avrtemp.Device{hidrawDevice:(*hidraw.Device)(0xc20803c080), status:avrtemp.Status{Device:"/dev/hidraw0", Time:time.Time{sec:0, nsec:0x0, loc:(*time.Location)(nil)}, SensorCount:0x0}}
        2016/01/23 13:41:21 server.Server: Start avrtemp device hidraw0
        2016/01/23 13:41:23 server.Server: Stat 28ffc256921503ae: 245
        2016/01/23 13:41:24 server.Server: Stat 28ff1032921503ed: 244
        2016/01/23 13:41:25 server.Server: Stat 28ffc256921503ae: 246
        2016/01/23 13:41:26 server.Server: Stat 28ff1032921503ed: 244
        2016/01/23 13:41:27 server.Server: Stat 28ffc256921503ae: 246
        2016/01/23 13:41:28 server.Server: Stat 28ff1032921503ed: 244
        2016/01/23 13:41:28 server.Device hidraw0: avrtemp.Device /dev/hidraw0: Read: read /dev/hidraw0: input/output error
        2016/01/23 13:41:28 RemoveHidrawDevice hidraw0...
        2016/01/23 13:41:28 server.Server: Stop device hidraw0
        2016/01/23 13:41:45 AddHidrawDevice hidraw.DeviceInfo{Device:libudev.Device{DevPath:"/devices/pci0000:00/0000:00:14.0/usb1/1-2/1-2:1.0/0003:16C0:0480.006E/hidraw/hidraw0", Subsystem:"hidraw", DevType:"", SysPath:"/sys/devices/pci0000:00/0000:00:14.0/usb1/1-2/1-2:1.0/0003:16C0:0480.006E/hidraw/hidraw0", SysName:"hidraw0", SysNum:"0", DevNode:"/dev/hidraw0"}, DeviceConfig:hidraw.DeviceConfig{VendorID:0x16c0, ProductID:0x480}}: &avrtemp.Device{hidrawDevice:(*hidraw.Device)(0xc20803c150), status:avrtemp.Status{Device:"/dev/hidraw0", Time:time.Time{sec:0, nsec:0x0, loc:(*time.Location)(nil)}, SensorCount:0x0}}
        2016/01/23 13:41:45 server.Server: Start avrtemp device hidraw0
        2016/01/23 13:41:47 server.Server: Stat 28ffc256921503ae: 245
        2016/01/23 13:41:48 server.Server: Stat 28ff1032921503ed: 243
        2016/01/23 13:41:49 server.Server: Stat 28ffc256921503ae: 245
        2016/01/23 13:41:50 server.Server: Stat 28ff1032921503ed: 243


## InfluxDB

Supports writing stats to `server -influxdb-server=... -influxdb-database=...`.

![Grafana Screenshot](/docs/grafana.png?raw=true "Grafana")

Grafana Query:

    SELECT mean("temperature") FROM "onewire" WHERE "family" = 'ds18b20' AND $timeFilter GROUP BY time($interval), "id", "name" fill(null)

## Config
The configuration file can be used to name connected sensors.

    [sensors.test1]
    ID  = "28ff1032921503ed"

    [sensors.test2]
    ID  = "28ffc256921503ae"

## REST API

    $ curl -s http://localhost:8283/api/config |json_pp
    {
       "Sensors" : {
          "test1" : {
             "ID" : "28ff1032921503ed"
          },
          "test2" : {
             "ID" : "28ffc256921503ae"
          }
       }
    }
    $ curl -s http://localhost:8283/api/ |json_pp
    [
       {
          "hidraw_device" : {
             "VendorID" : 5824,
             "DevNode" : "/dev/hidraw0",
             "SysName" : "hidraw0",
             "ProductID" : 1152,
             "DevPath" : "/devices/pci0000:00/0000:00:14.0/usb1/1-2/1-2:1.0/0003:16C0:0480.0068/hidraw/hidraw0",
             "SysPath" : "/sys/devices/pci0000:00/0000:00:14.0/usb1/1-2/1-2:1.0/0003:16C0:0480.0068/hidraw/hidraw0",
             "DevType" : "",
             "SysNum" : "0",
             "Subsystem" : "hidraw"
          },
          "avrtemp_device" : {
             "sensor_count" : 2,
             "time" : "2016-01-23T12:40:58.593460789+02:00",
             "device" : "/dev/hidraw0"
          },
          "stats" : {
             "28ffc256921503ae" : "test2",
             "28ff1032921503ed" : "test1"
          },
          "name" : "hidraw0"
       }
    ]
    $ curl -s http://localhost:8283/api/stats |json_pp
    [
       {
          "sensor_name" : "test2",
          "family" : "ds18b20",
          "temperature" : 29.9,
          "id" : "28ffc256921503ae",
          "time" : "2016-01-23T12:41:05.703672347+02:00"
       },
       {
          "temperature" : 25.4,
          "sensor_name" : "test1",
          "family" : "ds18b20",
          "time" : "2016-01-23T12:41:04.687577945+02:00",
          "id" : "28ff1032921503ed"
       }
    ]

## Devices

The server uses `libudev` to enumerate and monitor connected USB devices matching the configured `-device-vendor=` `-device-product=`. USB devices can be disconnected and reconnected without needing to restart the server.

### Supported hardware

* Diamex GmbH / AVR Temp Sensor (USB `16c0:0480`)

### Dependencies

* libudev-dev

### Configuration

#### `/etc/udev/rules.d/90-hidraw.rules`

    KERNEL=="hidraw*", \
        GROUP="plugdev", MODE=0660

