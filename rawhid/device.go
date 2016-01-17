package rawhid

import (
    "fmt"
    "os"
)

type Device struct {
    file        *os.File
    reportSize  int
}

func Open(info DeviceInfo) (*Device, error) {
    device := &Device{}

    if file, err := os.OpenFile(info.DevFile, os.O_RDWR, 0); err != nil {
        return nil, err
    } else {
        device.file = file
    }

    if err := device.init(); err != nil {
        return nil, err
    }

    return device, nil
}

func (self *Device) init() error {
    if reportSize, err := self.ioctlGetReportSize(); err != nil {
        return err
    } else {
        self.reportSize = reportSize
    }

    return nil
}

func (self *Device) DevInfo() (devInfo DevInfo, err error) {
    return self.ioctlGetDevInfo()
}

func (self *Device) Read() ([]byte, error) {
    buf := make([]byte, self.reportSize)

    if n, err := self.file.Read(buf); err != nil {
        return nil, err
    } else if n != self.reportSize {
        return nil, fmt.Errorf("Short read")
    } else {
        return buf, nil
    }
}
