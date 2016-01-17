package rawhid

import (
    "os"
    "syscall"
)

type Device struct {
    file    *os.File
}

func Open(info DeviceInfo) (*Device, error) {
    device := &Device{}

    if file, err := os.OpenFile(info.DevFile, os.O_RDWR|syscall.O_NONBLOCK, 0); err != nil {
        return nil, err
    } else {
        device.file = file
    }

    return device, nil
}

func (self *Device) DevInfo() (DevInfo, error) {
    return ioctlGetRawInfo(self.file)
}
