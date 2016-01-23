package libudev

import (
    "fmt"
)

type Error struct {
    Func        string
}

func (e Error) Error() string {
    return fmt.Sprintf("%s", e.Func)
}
