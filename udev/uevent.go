package udev

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type UeventInfo map[string]string

func readUevent(ueventPath string) (UeventInfo, error) {
    file, err := os.Open(ueventPath)
    if err != nil {
        return nil, err
    }

    scanner := bufio.NewScanner(file)
    info := make(UeventInfo)

    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.SplitN(line, "=", 2)

        if len(parts) != 2 {
            return nil, fmt.Errorf("Invalid line: %#v", parts)
        }

        info[parts[0]] = parts[1]
    }

    return info, nil
}

