package libudev

// #cgo LDFLAGS: -ludev
// #include <libudev.h>
// #include <stdlib.h>
import "C"
import "unsafe"

type Device struct {
    DevPath     string
    Subsystem   string
    DevType     string
    SysPath     string
    SysName     string
    SysNum      string
    DevNode     string
}

func (self Device) String() string {
    return self.SysName
}

func (self *Device) fromUdev(ptr *C.struct_udev_device) error {
    self.DevPath = C.GoString(C.udev_device_get_devpath(ptr))
    self.Subsystem = C.GoString(C.udev_device_get_subsystem(ptr))
    self.DevType = C.GoString(C.udev_device_get_devtype(ptr))
    self.SysPath = C.GoString(C.udev_device_get_syspath(ptr))
    self.SysName = C.GoString(C.udev_device_get_sysname(ptr))
    self.SysNum = C.GoString(C.udev_device_get_sysnum(ptr))
    self.DevNode = C.GoString(C.udev_device_get_devnode(ptr))

    return nil
}

func (self Device) ParentWithSubsystemDevType(subsystem string, devtype string) (device Device, err error) {
    udev := C.udev_new()
    defer C.udev_unref(udev)

    cSysPath := C.CString(self.SysPath)
    defer C.free(unsafe.Pointer(cSysPath))

    udev_device := C.udev_device_new_from_syspath(udev, cSysPath)
    if udev_device == nil {
        return device, Error{"udev_device_new_from_syspath"}
    }
    defer C.udev_device_unref(udev_device)

    // parent
    cSubsystem := C.CString(subsystem)
    defer C.free(unsafe.Pointer(cSubsystem))

    var cDevType *C.char
    if devtype != "" {
        cDevType = C.CString(devtype)
        defer C.free(unsafe.Pointer(cDevType))
    }

    udev_parent := C.udev_device_get_parent_with_subsystem_devtype(udev_device, cSubsystem, cDevType)
    if udev_parent == nil {
        return device, Error{"udev_device_get_parent"}
    }
    // lifetime tied to udev_device

    // return
    if err := device.fromUdev(udev_parent); err != nil {
        return device, err
    }

    return device, nil
}

// Returns nil on device error, missing entry on sysattr error
func (self Device) SysAttrs(sysAttrs ...string) map[string]string {
    udev := C.udev_new()
    defer C.udev_unref(udev)

    cSysPath := C.CString(self.SysPath)
    defer C.free(unsafe.Pointer(cSysPath))

    udev_device := C.udev_device_new_from_syspath(udev, cSysPath)
    if udev_device == nil {
        return nil
    }
    defer C.udev_device_unref(udev_device)

    // get
    out := make(map[string]string)

    for _, sysAttr := range sysAttrs {
        cSysAttr := C.CString(sysAttr)
        defer C.free(unsafe.Pointer(cSysAttr))

        cValue := C.udev_device_get_sysattr_value(udev_device, cSysAttr)
        if cValue == nil {
            continue
        }

        out[sysAttr] = C.GoString(cValue)
    }

    return out
}


