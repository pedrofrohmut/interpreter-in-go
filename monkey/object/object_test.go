// monkey/object/object_test.go

package object

import (
    "testing"
)

func TestStringHashKey(t *testing.T) {
    var hello1 = &String { Value: "Hello World" }
    var hello2 = &String { Value: "Hello World" }
    var diff1 = &String { Value: "My name is johnny" }
    var diff2 = &String { Value: "My name is johnny" }

    if hello1.HashKey() != hello2.HashKey() {
        t.Errorf("Strings with same content have different hash keys")
    }

    if diff1.HashKey() != diff2.HashKey() {
        t.Errorf("Strings with same content have different hash keys")
    }

    if hello1.HashKey() == diff1.HashKey() {
        t.Errorf("Strings with different content have the same hash keys")
    }
}
