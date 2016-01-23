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

type Server struct {
    // only modified at startup
    config          Config
    sensorConfig    map[string]*SensorConfig

    // state private to run()
    devices         map[string]*Device
    stats           map[string]Stat

    deviceChan          chan Device     // add/remove Device
    statChan            chan Stat       // read in from Devices
    influxChan          chan Stat       // write out to influx

    // HTTP API requests, handled by run()
    apiStatusChan       chan chan APIStatus
    apiStatChan         chan chan APIStat
}

func New() (*Server, error) {
    server := &Server{
        sensorConfig:   make(map[string]*SensorConfig),

        devices:        make(map[string]*Device),
        stats:          make(map[string]Stat),

        deviceChan:     make(chan Device),
        statChan:       make(chan Stat),

        apiStatusChan:  make(chan chan APIStatus),
        apiStatChan:    make(chan chan APIStat),
    }

    go server.run()

    return server, nil
}

func (s *Server) apiStatus(statusChan chan APIStatus) {
    defer close(statusChan)

    for name, device := range s.devices {
        status := APIStatus{
            Name:           name,
            HidrawDevice:   device.hidraw,
            Stats:          make(map[string]string),
        }

        if device.avrtemp != nil {
            status.AvrtempDevice = device.avrtemp.Status()
        }

        for statID, stat := range s.stats {
            if stat.Device != device {
                continue
            } else {
                // nil-safe
                status.Stats[statID] = s.sensorConfig[statID].String()
            }
        }

        statusChan <- status
    }
}

func (s *Server) apiStat(statChan chan APIStat) {
    defer close(statChan)

    for id, stat := range s.stats {
        apiStat := APIStat{
            ID:             id,
            Family:         stat.ID.Family(),
            Time:           stat.Time,
            Temperature:    stat.Temperature.Float64(),
        }

        if stat.SensorConfig != nil {
            apiStat.SensorName = stat.SensorConfig.String()
        }

        statChan <- apiStat
    }
}

func (s *Server) run() {
    var stat Stat

    for {
        select {
        case device := <-s.deviceChan:
            if device.avrtemp != nil {
                runningDevice := &device

                log.Printf("server.Server: Start avrtemp device %v\n", runningDevice)

                // run
                s.devices[device.String()] = runningDevice

                go runningDevice.reader(s.statChan)

            } else if runningDevice := s.devices[device.String()]; runningDevice != nil {
                log.Printf("server.Server: Stop device %v\n", runningDevice)

                // shutdown
                if runningDevice.avrtemp != nil {
                    runningDevice.avrtemp.Close()
                }

                delete(s.devices, device.String())
            }

        case stat = <-s.statChan:
            log.Printf("server.Server: Stat %v: %v\n", stat, stat.Temperature)

            stat.SensorConfig = s.sensorConfig[stat.ID.String()]

            s.stats[stat.String()] = stat

            s.influxChan <- stat

        case statusChan := <-s.apiStatusChan:
            s.apiStatus(statusChan)
        case statChan := <-s.apiStatChan:
            s.apiStat(statChan)
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
    device := Device{
        hidraw: deviceInfo,
    }

    if hidrawDevice, err := hidraw.Open(deviceInfo); err != nil {
        log.Printf("AddHidrawDevice %#v: hidraw.Open: %v\n", deviceInfo, err)
    } else if avrtempDevice, err := avrtemp.Open(hidrawDevice); err != nil {
        log.Printf("AddHidrawDevice %#v: avrtemp.Open: %v\n", deviceInfo, err)
    } else {
        log.Printf("AddHidrawDevice %#v: %#v\n", deviceInfo, avrtempDevice)

        device.avrtemp = avrtempDevice
    }

    s.deviceChan <- device
}

func (s *Server) RemoveHidrawDevice(deviceInfo hidraw.DeviceInfo) {
    log.Printf("RemoveHidrawDevice %v...\n", deviceInfo)

    s.deviceChan <- Device{hidraw: deviceInfo}
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
