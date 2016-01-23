package server

import (
    "github.com/qmsk/onewire/avrtemp"
    "fmt"
    "github.com/qmsk/onewire/hidraw"
    "log"
    "time"
)

type Device struct {
    hidraw          hidraw.DeviceInfo
    avrtemp         *avrtemp.Device
}

func (d Device) String() string {
    return d.hidraw.String()
}

type Stat struct {
    Device          *Device
    ID              avrtemp.ID
    SensorConfig    *SensorConfig

    Time            time.Time
    Temperature     avrtemp.Temperature
}

func (stat Stat) String() string {
    return fmt.Sprintf("%v", stat.ID)
}

type Server struct {
    // XXX: these maps are all racy
    config          Config
    sensorConfig    map[string]*SensorConfig

    devices         map[string]*Device
    stats           map[string]Stat

    statChan            chan Stat
    influxChan          chan Stat
}

func New() (*Server, error) {
    server := &Server{
        sensorConfig:   make(map[string]*SensorConfig),

        devices:        make(map[string]*Device),
        stats:          make(map[string]Stat),

        statChan:       make(chan Stat),
    }

    go server.run()

    return server, nil
}

func (s *Server) run() {
    var stat Stat
    var influxChan chan Stat

    for {
        select {
        case stat = <-s.statChan:
            log.Printf("server.Server: Stat %v: %v\n", stat, stat.Temperature)

            stat.SensorConfig = s.sensorConfig[stat.ID.String()]

            s.stats[stat.String()] = stat

            influxChan = s.influxChan

        // send once
        case influxChan <- stat:
            influxChan = nil
        }
    }
}

func (device *Device) reader(statChan chan Stat) {
    for {
        if report, err := device.avrtemp.Read(); err != nil {
            log.Printf("server.Device %v: avrtemp.Device %v: Read: %v\n", device, err)
            break
        } else {
            log.Printf("server.Device %v: avrtemp.Device %v: Read: %v\n", device, report)

            statChan <- Stat{
                Device:         device,
                ID:             report.ID,
                Time:           time.Now(),
                Temperature:    report.Temp,
            }
        }
    }
}

func (s *Server) AddHidrawDevice(deviceInfo hidraw.DeviceInfo) {
    device := &Device{
        hidraw: deviceInfo,
    }

    if hidrawDevice, err := hidraw.Open(deviceInfo); err != nil {
        log.Printf("AddHidrawDevice %#v: hidraw.Open: %v\n", deviceInfo, err)
    } else if avrtempDevice, err := avrtemp.Open(hidrawDevice); err != nil {
        log.Printf("AddHidrawDevice %#v: avrtemp.Open: %v\n", deviceInfo, err)
    } else {
        log.Printf("AddHidrawDevice %#v: %#v\n", deviceInfo, avrtempDevice)

        device.avrtemp = avrtempDevice

        go device.reader(s.statChan)
    }

    s.devices[device.String()] = device
}

func (s *Server) RemoveHidrawDevice(deviceInfo hidraw.DeviceInfo) {
    log.Printf("RemoveHidrawDevice %v...\n", deviceInfo)

    device := s.devices[deviceInfo.String()]

    if device != nil && device.avrtemp != nil {
        device.avrtemp.Close()
    }

    delete(s.devices, deviceInfo.String())
}

func (s *Server) MonitorHidraw(monitorChan chan hidraw.MonitorEvent) {
    for monitorEvent := range monitorChan {
        switch monitorEvent.Action {
        case "add":
            s.AddHidrawDevice(monitorEvent.DeviceInfo)
        case "remove":
            s.RemoveHidrawDevice(monitorEvent.DeviceInfo)
        default:
            log.Printf("MonitorHidraw: %v?! %v\n", monitorEvent.Action, monitorEvent.DeviceInfo)
        }
    }

    log.Printf("server.Server %v: MonitorHidraw: exit\n", s)
}
