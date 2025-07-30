// monkey/utils/utils.go

package utils

import (
    "reflect"
)

func IsNill(tmp any) bool {
    if tmp == nil || reflect.ValueOf(tmp).IsNil() {
        return true
    }
    return false
}
