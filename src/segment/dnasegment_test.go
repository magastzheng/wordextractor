package segment_test

import (
    "fmt"
    "segment"
    "testing"
    "util"
)

func Test_DNASegment(t *testing.T){
    s := segment.NewDNASegment()
    //s.InitDict("../data/dictionary/dictionary.txt")
    s.InitDict("../data/dictionary/sogoudictionary.txt")
    
    text := util.ReadFile("../data/testdata/125-1.txt")
    fmt.Println(len(text))
    segs := s.MPSeg(text)
    fmt.Println(len(segs))
    fmt.Println(segs)

    fmt.Println(306234192.0/301869396788.0)
}


