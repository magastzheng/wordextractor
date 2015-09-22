package segment_test

import (
    "fmt"
    "segment"
    "testing"
)

func Test_Dictionary(t *testing.T) {
    d := segment.NewDictionary("../data/dictionary/sogoudictionary.txt")
    freq := d.FindWord("大学生")
    fmt.Println(freq)
}
