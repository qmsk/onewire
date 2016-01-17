package rawhid

// #include "ioctl.h"
import "C"

import (
    "os"
)

type DevInfo struct {
    BusType     uint
    Vendor      int
    Product     int
}

func ioctlGetRawInfo(dev *os.File) (devInfo DevInfo, err error) {
    var hdi C.struct_hidraw_devinfo

    if _, err := C.hidraw_ioctl_getrawinfo(C.int(dev.Fd()), &hdi); err != nil {
        return devInfo, err
    }

    return DevInfo{
        uint(hdi.bustype),
        int(hdi.vendor),
        int(hdi.product),
    }, nil
}
