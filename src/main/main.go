package main

import(
    "fmt"
    "dict"
    "segment"
    "util"
)

func main() {
    sign := dict.NewSign("../data/dictionary/sign.txt")
    //stop := dict.NewSign("../data/dictionary/stopwords.txt")
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")

    article := util.ReadFile("../data/testdata/125-1.txt")
    allsegs := segment.SegmentDoc(article, sign, d)

    fmt.Println(len(allsegs))
}
