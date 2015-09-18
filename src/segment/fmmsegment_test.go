package segment_test

import (
    "testing"
    "fmt"
    "segment"
    "util"
)

func Test_FMMSegement(t *testing.T) {
    seg := segment.NewFMMSegment("../data/dictionary/CoreNatureDictionary.mini.txt")
    content := util.ReadFile("../data/testdata/125-1.txt")
    words := seg.Segment(content)
    for _, w := range words {
        fmt.Println(w)
    }
}
