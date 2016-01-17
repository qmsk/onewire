package rawhid

import (
    "github.com/qmsk/onewire/udev"
)

type DeviceInfo struct {
    Name    string
    SysPath string
    DevPath string
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

            if devFile, err := deviceNode.DevFile(); err != nil {
                return nil, err
            } else {
                deviceInfo.DevPath = devFile
            }

            devices = append(devices, deviceInfo)
        }

        return devices, nil
    }
}
