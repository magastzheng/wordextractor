package segment_test

import (
    "fmt"
    "segment"
    "testing"
    "util"
	"dict"
)

func Test_MPSegment(t *testing.T) {
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")
    text := util.ReadFile("../data/testdata/125-1.txt")

    out := segment.SegmentSentence_MP(text, d)
    fmt.Println(out)
}

func Test_SegmentSentenceMP(t *testing.T) {
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")
    text := util.ReadFile("../data/testdata/125-1.txt")
    segments := segment.SegmentSentenceMP([]rune(text), 0, d)
    //str := segment.Output(segments)
    fmt.Println(len(segments))
}
