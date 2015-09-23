package segment_test

import (
    "fmt"
    "segment"
    "testing"
    "util"
    "dict"
)

func Test_NewDocument(t *testing.T) {
    doc := segment.NewDocument("../data/testdata/125-1.txt")
    fmt.Println(doc.Filename())
    fmt.Println(string(doc.Buffer()))
}

func Test_SplitSentence(t *testing.T) {
    article := util.ReadFile("../data/testdata/125-1.txt")
    fmt.Println(article)
    article = segment.DeleteSpaceChar(article)
    d := dict.NewSign("../data/dictionary/sign.txt")
    sentences := segment.SplitSentence([]rune(article), d)
    fmt.Println(len(sentences))
    //for _, s := range sentences {
    //    fmt.Println(s.ToString())
    //}
}
