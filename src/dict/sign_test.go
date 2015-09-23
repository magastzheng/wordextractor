package dict_test

import (
    "dict"
    "testing"
    "fmt"
)

func Test_Sign(t *testing.T) {
    s := dict.NewSign("../data/dictionary/sign.txt")
    end := "。" 
    fmt.Println(len([]rune(end)))
    ret := s.IsContain(end)
    fmt.Println(ret)

    ret = s.IsContain(".")
    fmt.Println(ret)
}

func Test_Stop(t *testing.T) {
    s := dict.NewSign("../data/dictionary/stopwords.txt")
    word := "out of"
    ret := s.IsContain(word)

    fmt.Println(ret)
}
