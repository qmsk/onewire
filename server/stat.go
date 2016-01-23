package server

import (
    "github.com/qmsk/onewire/avrtemp"
    "fmt"
    "log"
    "time"
)

type Stat struct {
    // set by Device.reader()
    Device          *Device
    ID              avrtemp.ID
    Time            time.Time
    Temperature     avrtemp.Temperature

    // set by Server.run()
    SensorConfig    *SensorConfig
}

func (stat Stat) String() string {
    return fmt.Sprintf("%v", stat.ID)
}

func (s *Server) stat(stat Stat) {
    log.Printf("server.Server: Stat %v: %v\n", stat, stat.Temperature)

    stat.SensorConfig = s.sensorConfig[stat.ID.String()]

    s.stats[stat.String()] = stat

    s.influxChan <- stat
}
