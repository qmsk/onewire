package rawhid

import (
    "fmt"
    "github.com/qmsk/onewire/udev"
)

type DeviceInfo struct {
    Name        string
    SysPath     string
    DevPath     string

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
                SysPath:    deviceNode.Path(),
            }

            if devFile := deviceNode.DevFile(); devFile != "" {
                deviceInfo.DevPath = devFile
            }

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
