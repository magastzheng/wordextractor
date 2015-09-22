package segment_test

import (
    "fmt"
    "segment"
    "testing"
    "util"
)

func Test_MPSegment(t *testing.T) {
    d := segment.NewDictionary("../data/dictionary/sogoudictionary.txt")
    text := util.ReadFile("../data/testdata/125-1.txt")

    out := segment.SegmentSentence_MP(text, d)
    fmt.Println(out)
}
