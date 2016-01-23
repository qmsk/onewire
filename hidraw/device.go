package hidraw

import (
    "os"
)

type Device struct {
    file        *os.File
}

func Open(info DeviceInfo) (*Device, error) {
    device := &Device{}

    if file, err := os.OpenFile(info.DevNode, os.O_RDWR, 0); err != nil {
        return nil, err
    } else {
        device.file = file
    }

    return device, nil
}

func (self *Device) String() string {
    return self.file.Name()
}

// USB HID Device Info
//
// Bus-level data about the device
func (self *Device) DevInfo() (devInfo DevInfo, err error) {
    return self.ioctlGetDevInfo()
}

// USB HID Report Descriptor
//
// This is a magic thing that specifies the format of the reports
// http://stackoverflow.com/questions/21606991/custom-hid-device-hid-report-descriptor
func (self *Device) ReportDescriptor() ([]byte, error) {
    return self.ioctlGetReportDescriptor()
}

func (self *Device) Read(buf []byte) (int, error) {
    return self.file.Read(buf)
}

func (self *Device) Close() {
    self.file.Close()
}
