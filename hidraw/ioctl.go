package hidraw

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

func (self *Device) ioctlGetReportDescriptor() ([]byte, error) {
    var size C.int
    var hrd C.struct_hidraw_report_descriptor

    if _, err := C.hidraw_ioctl_getrdescsize(self.fd(), &size); err != nil {
        return nil, err
    }
    if _, err := C.hidraw_ioctl_getrdesc(self.fd(), &hrd, size); err != nil {
        return nil, err
    }

    buf := make([]byte, int(hrd.size))

    for i := 0; i < len(buf); i++ {
        buf[i] = byte(hrd.value[i])
    }

    return buf, nil
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
