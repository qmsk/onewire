package udev

import (
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "strings"
)

type DeviceNode string // path

func (self DeviceNode) IsNil() bool {
    return self == ""
}

func (self DeviceNode) Name() string {
    return path.Base(string(self))
}

func (self DeviceNode) Path() string {
    return string(self)
}

func (self DeviceNode) path(parts ...string) string {
    return path.Join(append([]string{string(self)}, parts...)...)
}

func (self DeviceNode) Subsystem() string {
    if subsystemNode, err := openDevice(self.path("subsystem")); err != nil {
        return ""
    } else {
        return subsystemNode.Name()
    }
}

func (self DeviceNode) ReadString(name string) (string, error) {
    if data, err := ioutil.ReadFile(self.path(name)); err != nil {
        return "", err
    } else {
        return strings.TrimSpace(string(data)), nil
    }
}

func (self DeviceNode) ReadHex(name string) (value uint, err error) {
    if data, err := ioutil.ReadFile(self.path(name)); err != nil {
        return 0, err
    } else if _, err := fmt.Sscanf(string(data), "%x", &value); err != nil {
        return 0, err
    } else {
        return value, nil
    }

}

func (self DeviceNode) ReadUevent() (UeventInfo, error) {
    return readUevent(path.Join(string(self), "uevent"))
}

func (self DeviceNode) DevFile() string {
    if ueventInfo, err := self.ReadUevent(); err != nil {
        return ""
    } else if devName := ueventInfo["DEVNAME"]; devName == "" {
        return ""
    } else if devName[0] == '/' {
        return devName
    } else {
        return path.Join("/dev/", devName)
    }
}

func (self DeviceNode) DevType() string {
    if ueventInfo, err := self.ReadUevent(); err != nil {
        return ""
    } else if devType := ueventInfo["DEVTYPE"]; devType == "" {
        return ""
    } else {
        return devType
    }
}

func (self DeviceNode) Parent() DeviceNode {
    return DeviceNode(path.Dir(self.Path()))
}

func (self DeviceNode) ParentSubsystemDevType(subsystem string, devType string) DeviceNode {
    node := self

    for ; node.Name() != ""; node = node.Parent() {
        if subsystem != "" && node.Subsystem() != subsystem {
            continue
        }

        if devType != "" && node.DevType() != devType {
            continue
        }

        break
    }

    return node
}

func openDevice(filePath string) (DeviceNode, error) {
    dirPath := path.Dir(filePath)

    if linkPath, err := os.Readlink(filePath); err != nil {
        return "", err
    } else {
        return DeviceNode(path.Clean(path.Join(dirPath, linkPath))), nil
    }
}

func ListClass(class string) (nodes []DeviceNode, err error) {
    dirPath := path.Join("/sys/class", class)

    readdir, err := ioutil.ReadDir(dirPath)
    if err != nil {
        return nil, err
    }

    for _, fileInfo := range readdir {
        if deviceNode, err := openDevice(path.Join(dirPath, fileInfo.Name())); err != nil {
            continue
        } else {
            nodes = append(nodes, deviceNode)
        }
    }

    return nodes, nil
}
