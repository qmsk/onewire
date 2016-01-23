package server

import (
    "github.com/qmsk/onewire/avrtemp"
    "fmt"
    "github.com/qmsk/onewire/hidraw"
    "log"
    "time"
)

type Stat struct {
    Device          *avrtemp.Device
    ID              avrtemp.ID
    SensorConfig    *SensorConfig

    Time            time.Time
    Temperature     avrtemp.Temperature
}

func (stat Stat) String() string {
    return fmt.Sprintf("%v", stat.ID)
}

type Server struct {
    config          Config
    sensorConfig    map[string]*SensorConfig

    hidrawDevices   map[string]hidraw.DeviceInfo
    avrtempDevices  map[string]*avrtemp.Device
    stats           map[string]Stat

    statChan            chan Stat
}

func New() (*Server, error) {
    server := &Server{
        sensorConfig:   make(map[string]*SensorConfig),

        hidrawDevices:  make(map[string]hidraw.DeviceInfo),
        avrtempDevices: make(map[string]*avrtemp.Device),
        stats:          make(map[string]Stat),

        statChan:       make(chan Stat),
    }

    go server.run()

    return server, nil
}

func (s *Server) run() {
    for {
        select {
        case stat := <-s.statChan:
            log.Printf("server.Server: Stat %v: %v\n", stat, stat.Temperature)

            s.stats[stat.String()] = stat
        }
    }
}

func (s *Server) reader(avrtempDevice *avrtemp.Device) {
    for {
        if report, err := avrtempDevice.Read(); err != nil {
            log.Printf("server.Server: reader: avrtemp.Device %v: Read: %v\n", avrtempDevice, err)
            break
        } else {
            log.Printf("server.Server: reader: avrtemp.Device %v: Read: %v\n", avrtempDevice, report)

            s.statChan <- Stat{
                Device:         avrtempDevice,
                ID:             report.ID,
                SensorConfig:   s.sensorConfig[report.ID.String()],
                Time:           time.Now(),
                Temperature:    report.Temp,
            }
        }
    }
}

func (s *Server) AddHidrawDevice(deviceInfo hidraw.DeviceInfo) {
    s.hidrawDevices[deviceInfo.String()] = deviceInfo

    if hidrawDevice, err := hidraw.Open(deviceInfo); err != nil {
        log.Printf("AddHidrawDevice %#v: hidraw.Open: %v\n", deviceInfo, err)
    } else if avrtempDevice, err := avrtemp.Open(hidrawDevice); err != nil {
        log.Printf("AddHidrawDevice %#v: avrtemp.Open: %v\n", deviceInfo, err)
    } else {
        log.Printf("AddHidrawDevice %#v: %#v\n", deviceInfo, avrtempDevice)

        s.avrtempDevices[deviceInfo.String()] = avrtempDevice

        go s.reader(avrtempDevice)
    }
}

func (s *Server) RemoveHidrawDevice(deviceInfo hidraw.DeviceInfo) {
    log.Printf("RemoveHidrawDevice %v...\n", deviceInfo)

    if avrtempDevice := s.avrtempDevices[deviceInfo.String()]; avrtempDevice != nil {
        avrtempDevice.Close()
    }

    delete(s.avrtempDevices, deviceInfo.String())
    delete(s.hidrawDevices, deviceInfo.String())
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
