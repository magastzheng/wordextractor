package occurrence_test

import (
    "occurrence"
    "fmt"
    "testing"
    "dict"
    "util"
    "segment"
    "term"
)

func Test_SegmentMerger(t *testing.T) {
    sign := dict.NewSign("../data/dictionary/sign.txt")
    stop := dict.NewSign("../data/dictionary/stopwords.txt")
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")

    article := util.ReadFile("../data/testdata/125-2.txt")
    allsegs := segment.SegmentDoc(article, sign, d)
    allsegs = occurrence.FilterSegment(allsegs, stop) 
    fmt.Println(len(allsegs))
    str := segment.GetSegmentStr(allsegs)

    util.WriteFile("../data/test-125-2.log", str)

    occur := occurrence.NewOccurrence()
    occur.AddSegments(allsegs, 3)
    occur.Compute()
    occur.Output()

    pairTerms := occur.GetPairTerms(10.0)
    str = term.GetPairTermStr(pairTerms)
    util.WriteFile("../data/test-125-2-occur.log", str)

    newSegments := occurrence.MergeSegment(allsegs, pairTerms)
    fmt.Println(len(newSegments))
    str = segment.GetSegmentStr(newSegments)
    //fmt.Println(str)
    util.WriteFile("../data/test-125-2-merge.log", str)

    
    occur1 := occurrence.NewOccurrence()
    occur1.AddSegments(newSegments, 1)
    occur1.Compute()
    occur1.Output()
    pairTerms = occur1.GetPairTerms(10.0)
    str = term.GetPairTermStr(pairTerms)
    util.WriteFile("../data/test-125-2-second-merge.log", str)
}


func Test_Merge(t *testing.T) {
    sign := dict.NewSign("../data/dictionary/sign.txt")
    stop := dict.NewSign("../data/dictionary/stopwords.txt")
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")

    article := util.ReadFile("../data/testdata/125-2.txt")
    allsegs := segment.SegmentDoc(article, sign, d)
    allsegs = occurrence.FilterSegment(allsegs, stop) 
    
    pairTerms := occurrence.Merge(allsegs, 4, 15.0)
    str := term.GetPairTermStr(pairTerms)
    util.WriteFile("../data/test-125-2-merge-merge.log", str)
}

func Test_Merge125(t *testing.T) {
    sign := dict.NewSign("../data/dictionary/sign.txt")
    stop := dict.NewSign("../data/dictionary/stopwords.txt")
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")

    article := util.ReadFile("../data/testdata/125.txt")
    allsegs := segment.SegmentDoc(article, sign, d)
    allsegs = occurrence.FilterSegment(allsegs, stop) 
    
    pairTerms := occurrence.Merge(allsegs, 4, 15.0)
    str := term.GetPairTermStr(pairTerms)
    util.WriteFile("../data/test-125-merge-merge.log", str)
}
