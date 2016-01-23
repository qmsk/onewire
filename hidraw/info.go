package hidraw

import (
    "fmt"
    "log"
    "github.com/qmsk/onewire/libudev"
)

type DeviceConfig struct {
    VendorID    uint
    ProductID   uint
}

type DeviceInfo struct {
    libudev.Device
    DeviceConfig
}

type MonitorEvent struct {
    DeviceInfo

    Action      string
}

func (self *DeviceInfo) fromUdevDevice(udevDevice libudev.Device) error {
    // find USB device
    usbDevice, err := udevDevice.ParentWithSubsystemDevType("usb", "usb_device")
    if err != nil {
        return err
    }

    // USB attrs
    sysAttrs := usbDevice.SysAttrs("idVendor", "idProduct")

    if idVendor := sysAttrs["idVendor"]; idVendor == "" {
        return fmt.Errorf("udev.Device %v: SysAttr %v: null", usbDevice, "idVendor")
    } else if _, err := fmt.Sscanf(idVendor, "%x", &self.VendorID); err != nil {
        return fmt.Errorf("udev.Device %v: SysAttr %v: %v", usbDevice, "idVendor", err)
    }

    if idProduct := sysAttrs["idProduct"]; idProduct == "" {
        return fmt.Errorf("udev.Device %v: SysAttr %v: null", usbDevice, "idProduct")
    } else if _, err := fmt.Sscanf(idProduct, "%x", &self.ProductID); err != nil {
        return fmt.Errorf("udev.Device %v: SysAttr %v: %v", usbDevice, "idProduct", err)
    }

    return nil
}

func (self *DeviceInfo) filter(match DeviceConfig) bool {
    if match.VendorID != 0 && self.VendorID != match.VendorID {
        return false
    }
    if match.ProductID != 0 && self.ProductID != match.ProductID {
        return false
    }

    return true
}

func List(match DeviceConfig) (devices []DeviceInfo, err error) {
    udevDevices, err := libudev.Enumerate(libudev.Device{
        Subsystem: "hidraw",
    })
    if err != nil {
        return nil, err
    }

    for _, udevDevice := range udevDevices {
        deviceInfo := DeviceInfo{Device:udevDevice}

        if err := deviceInfo.fromUdevDevice(udevDevice); err != nil {
            return nil, err
        } else if !deviceInfo.filter(match) {

        } else {
            devices = append(devices, deviceInfo)
        }
    }

    return devices, nil
}

func monitor(monitorChan chan MonitorEvent, udevMonitor *libudev.Monitor, match DeviceConfig) {
    defer udevMonitor.Close()
    defer close(monitorChan)

    for {
        udevDevice, err := udevMonitor.Recv()
        if err != nil {
            log.Printf("rawhid.monitor: udev.Monitor: Recv: %v\n", err)
            continue
        }

        // make
        monitorEvent := MonitorEvent{
            DeviceInfo: DeviceInfo{
                Device: udevDevice.Device,
            },
            Action: udevDevice.Action,
        }

        if udevDevice.Action == "remove" {
            // skip USB data
            // XXX: how to filter?

        } else if err := monitorEvent.DeviceInfo.fromUdevDevice(udevDevice.Device); err != nil {
            log.Printf("rawhid.monitor: fromUdevDevice: %v\n", err)
            continue
        } else if !monitorEvent.DeviceInfo.filter(match) {
            continue
        }

        monitorChan <- monitorEvent
    }

    log.Printf("rawhid.monitor: exit\n")
}

// Only action=add events for the given match are passed through.
// All action=remove events are passed through, as they cannot be filtered.
func Monitor(match DeviceConfig) (chan MonitorEvent, error) {
    udevMonitor, err := libudev.MonitorUdev("hidraw")
    if err != nil {
        return nil, err
    }

    // run
    monitorChan := make(chan MonitorEvent)

    go monitor(monitorChan, udevMonitor, match)

    return monitorChan, nil
}

func Find(config DeviceConfig) (*Device, error) {
    if devices, err := List(config); err != nil {
        return nil, err
    } else {
        for _, deviceInfo := range devices {
            return Open(deviceInfo)
        }

        // not found
        return nil, nil
    }
}
