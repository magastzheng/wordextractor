package main

import(
    "fmt"
    "dict"
    "segment"
    "util"
    "occurrence"
    "term"
)

func main() {
    var freqDoor int
    var scoreDoor float32
    freqDoor = 4
    scoreDoor = 0.010

    sign := dict.NewSign("../data/dictionary/sign.txt")
    stop := dict.NewSign("../data/dictionary/stopwords.txt")
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")

    article := util.ReadFile("../data/testdata/125.txt")
    allsegs := segment.SegmentDoc(article, sign, d)
    
    fmt.Println(len(allsegs))
    str := segment.GetSegmentStr(allsegs)

    util.WriteFile("../data/test-125.log", str)
    allsegs = occurrence.FilterSegment(allsegs, stop)

    //occur := occurrence.NewOccurrence()
    //occur.AddSegments(allsegs, stop)
    //occur.Compute()
    //occur.Output()

    pairTerms := occurrence.Merge(allsegs, freqDoor, scoreDoor)
    str = term.GetPairTermStr(pairTerms)
    util.WriteFile("../data/main-test-125-merge-merge.log", str)
}
