package rawhid

// #include "ioctl.h"
import "C"

type DevInfo struct {
    BusType     uint
    Vendor      int
    Product     int
}

func (self *Device) fd() C.int {
    return C.int(self.file.Fd())
}

func (self *Device) ioctlGetReportSize() (int, error) {
    var size C.int

    if _, err := C.hidraw_ioctl_getrdescsize(self.fd(), &size); err != nil {
        return 0, err
    }
    
    return int(size), nil
}

func (self *Device) ioctlGetDevInfo() (devInfo DevInfo, err error) {
    var hdi C.struct_hidraw_devinfo

    if _, err := C.hidraw_ioctl_getrawinfo(self.fd(), &hdi); err != nil {
        return devInfo, err
    }

    devInfo = DevInfo{
        uint(hdi.bustype),
        int(hdi.vendor),
        int(hdi.product),
    }
    return
}
