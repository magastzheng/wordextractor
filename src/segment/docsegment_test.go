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
    //fmt.Println(string(doc.Buffer()))
}

func Test_SplitSentence(t *testing.T) {
    article := util.ReadFile("../data/testdata/125-1.txt")
    //fmt.Println(article)
    article = segment.DeleteSpaceChar(article)
    d := dict.NewSign("../data/dictionary/sign.txt")
    sentences := segment.SplitSentence([]rune(article), d)
    fmt.Println(len(sentences))
    //for _, s := range sentences {
    //    fmt.Println(s.ToString())
    //}
}

func Test_SplitDocument(t *testing.T){
    article := util.ReadFile("../data/testdata/125-1.txt")
    //fmt.Println(article)
    article = segment.DeleteSpaceChar(article)
    sign := dict.NewSign("../data/dictionary/sign.txt")
    sentences := segment.SplitSentence([]rune(article), sign)
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")
    fmt.Println("Start====")

    allsegs := make([]*segment.Segment, 0)
    for _, sentence := range sentences {
        segments := segment.SegmentSentenceMP(sentence.Buffer(), sentence.Start(), d)
        //fmt.Println(len(segments))
        //fmt.Println(segments)
        //str := ""
        //for _, seg := range segments {
        //    str += seg.ToString()
        //}

        //fmt.Println(sentence.Start(), str)

        allsegs = append(allsegs, segments ...)
    }

    fmt.Println(len(allsegs))
}


func Test_Segment(t *testing.T){
    article := util.ReadFile("../data/testdata/125-1.txt")
	sign := dict.NewSign("../data/dictionary/sign.txt")
    //sentences := SplitSentence([]rune(text), sign)
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")
    allsegs := segment.SegmentDoc(article, sign, d)

    fmt.Println(len(allsegs))
    final := ""
    for _, seg := range allsegs {
        final += seg.ToString()
    }
    
    final += fmt.Sprintf(": %d", len(allsegs))
    util.WriteFile("../data/docsegment_test_125-1.log", final)
}
