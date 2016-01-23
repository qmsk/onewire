package server

import (
    "github.com/qmsk/onewire/avrtemp"
    "fmt"
    "encoding/json"
    "github.com/qmsk/onewire/hidraw"
    "net/http"
    "log"
)

type APIStatus struct {
    Name            string              `json:"name"`
    HidrawDevice    hidraw.DeviceInfo   `json:"hidraw_device"`
    AvrtempDevice   string              `json:"avrtemp_device"`
}

type Server struct {
    hidrawDevices   map[string]hidraw.DeviceInfo
    avrtempDevices  map[string]*avrtemp.Device
}

func New() (*Server, error) {
    server := &Server{
        hidrawDevices:  make(map[string]hidraw.DeviceInfo),
        avrtempDevices: make(map[string]*avrtemp.Device),
    }

    return server, nil
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

func (s *Server) GetStatus() (interface{}, error) {
    var statusList []APIStatus

    for name, hidrawDevice := range s.hidrawDevices {
        status := APIStatus{
            Name:           name,
            HidrawDevice:   hidrawDevice,
        }

        if avrtempDevice := s.avrtempDevices[name]; avrtempDevice != nil {
            status.AvrtempDevice = avrtempDevice.String()
        }

        statusList = append(statusList, status)
    }

    return statusList, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if jsonData, err := s.GetStatus(); err != nil {
        w.WriteHeader(500)

        fmt.Fprintf(w, "%v\n", err)
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)

        if err := json.NewEncoder(w).Encode(jsonData); err != nil {
            panic(err)
        }
    }
}
