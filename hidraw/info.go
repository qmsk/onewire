package hidraw

import (
    "fmt"
    "github.com/qmsk/onewire/udev"
)

type DeviceConfig struct {
    Vendor      uint
    Product     uint
}

type DeviceInfo struct {
    Name        string
    SysDevice   string
    DevFile     string

    VendorID    uint16
    ProductID   uint16
}

func List() (devices []DeviceInfo, err error) {
    if deviceNodes, err := udev.ListClass("hidraw"); err != nil {
        return nil, err
    } else {
        for _, deviceNode := range deviceNodes {
            deviceInfo := DeviceInfo{
                Name:       deviceNode.Name(),
                SysDevice:  deviceNode.Path(),
            }

            if devFile := deviceNode.DevFile(); devFile != "" {
                deviceInfo.DevFile = devFile
            }

            // read USB device info
            if usbDevice := deviceNode.ParentSubsystemDevType("usb", "usb_device"); usbDevice.IsNil() {
                return nil, fmt.Errorf("Could not find parent USB device: %v", deviceNode)
            } else {
                if idVendor, err := usbDevice.ReadHex("idVendor"); err != nil {

                } else {
                    deviceInfo.VendorID = uint16(idVendor)
                }

                if idProduct, err := usbDevice.ReadHex("idProduct"); err != nil {

                } else {
                    deviceInfo.ProductID = uint16(idProduct)
                }
            }

            devices = append(devices, deviceInfo)
        }

        return devices, nil
    }
}

func Select(config DeviceConfig) (*Device, error) {
    if devices, err := List(); err != nil {
        return nil, err
    } else {
        for _, deviceInfo := range devices {
            if config.Vendor != 0 && uint(deviceInfo.VendorID) != config.Vendor {
                continue
            }
            if config.Product != 0 && uint(deviceInfo.ProductID) != config.Product {
                continue
            }

            return Open(deviceInfo)
        }

        // not found
        return nil, nil
    }
}
