package dict_test

import (
    "fmt"
    "dict"
    "testing"
)

func Test_Dictionary(t *testing.T) {
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")
    freq := d.FindWord("大学生")
    fmt.Println(freq)
}
