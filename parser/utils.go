package parser

import (
    "strings"
)

func MangleFunction(name string) string {
    ptrize := strings.Replace(name, "*", "_ptr", -1)
    dotize := strings.Replace(ptrize, ".", "_", -1)
    return dotize
}
