package server

import (
    "github.com/qmsk/onewire/avrtemp"
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

func (device *Device) reader(statChan chan Stat) {
    for {
        if report, err := device.avrtemp.Read(); err != nil {
            log.Printf("server.Device %v: avrtemp.Device %v: Read: %v\n", device, err)
            break
        } else {
            statChan <- Stat{
                Device:         device,
                ID:             report.ID,
                Time:           time.Now(),
                Temperature:    report.Temp,
            }
        }
    }
}

func (s *Server) addDevice(device *Device) {
    log.Printf("server.Server: Start avrtemp device %v\n", device)

    // run
    s.devices[device.String()] = device

    go device.reader(s.statChan)
}

func (s *Server) removeDevice(device *Device) {
    log.Printf("server.Server: Stop device %v\n", device)

    // shutdown
    if device.avrtemp != nil {
        device.avrtemp.Close()
    }

    delete(s.devices, device.String())
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
