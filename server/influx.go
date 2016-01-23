package server

import (
    influxdb "github.com/influxdata/influxdb/client/v2"
    "log"
    "net"
)

const INFLUX_PORT = "8086"

type InfluxOptions struct {
    Server      string
    Database    string
}

func (opts InfluxOptions) Empty() bool {
    return opts.Server == ""
}

func (opts InfluxOptions) httpConfig() (httpConfig influxdb.HTTPConfig) {
    if host, port, err := net.SplitHostPort(opts.Server); err == nil && host != "" && port != "" {
        httpConfig.Addr = "http://" + net.JoinHostPort(host, port)
    } else if err == nil && host != "" {
        httpConfig.Addr = "http://" + net.JoinHostPort(host, INFLUX_PORT)
    } else {
        httpConfig.Addr = "http://" + net.JoinHostPort(opts.Server, INFLUX_PORT)
    }

    return
}

func (opts InfluxOptions) batchPoints() (influxdb.BatchPoints, error) {
    return influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
        Database:   opts.Database,
    })
}

func (self *Server) influxWriter(influxClient influxdb.Client, options InfluxOptions) {
    defer close(self.influxChan)
    defer influxClient.Close()

    for stat := range self.influxChan {
        tags := map[string]string{
            "id":       stat.ID.String(),
            "family":   stat.ID.Family(),
            "name":     stat.SensorConfig.String(),
        }
        fields := map[string]interface{}{
            "temperature":  stat.Temperature.Float64(),
        }

        // write
        point, err := influxdb.NewPoint("onewire", tags, fields, stat.Time)
        if err != nil {
            log.Printf("server.Server: influxWriter: influxdb.NewPoint: %v\n", err)
        }

        points, err := options.batchPoints()
        if err != nil {
            log.Printf("server.Server: influxWriter: newBatchPoints: %v\n", err)
            continue
        }

        points.AddPoint(point)

        if err := influxClient.Write(points); err != nil {
            log.Printf("server.Server: influxWriter: influxdb.Client %v: Write %v: %v\n", influxClient, points, err)
            continue
        }
    }
}

func (self *Server) InfluxWriter(options InfluxOptions) error {
    influxClient, err := influxdb.NewHTTPClient(options.httpConfig())
    if err != nil {
        return err
    }

    // start goroutine
    self.influxChan = make(chan Stat)

    go self.influxWriter(influxClient, options)

    return nil
}
