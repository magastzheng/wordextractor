package main

import(
    "fmt"
    "dict"
    "segment"
    "util"
    "occurrence"
)

func main() {
    sign := dict.NewSign("../data/dictionary/sign.txt")
    stop := dict.NewSign("../data/dictionary/stopwords.txt")
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")

    article := util.ReadFile("../data/testdata/125.txt")
    allsegs := segment.SegmentDoc(article, sign, d)
    
    fmt.Println(len(allsegs))
    str := segment.GetSegmentStr(allsegs)

    util.WriteFile("../data/test-125.log", str)

    occur := occurrence.NewOccurrence()
    occur.AddSegments(allsegs, stop)
    occur.Compute()
    occur.Output()
}
