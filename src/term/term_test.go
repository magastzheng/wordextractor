package term_test

import (
    "term"
    "fmt"
    "sort"
    "testing"
)

func Test_NewTerm(t *testing.T) {
    t1 := term.NewTerm("test1")
    fmt.Println(t1)
}

func Test_NewNewPairTerm(t *testing.T) {
    t3 := term.NewPairTerm("key3", "This", "is")
    fmt.Println(t3)
}

func Test_SortPairTerm(t *testing.T) {
    t1 := term.NewPairTerm("key1", "F1", "S1")
    t1.SetMI(6.21)
    t2 := term.NewPairTerm("key2", "F2", "S2")
    t2.SetMI(1.62)
    t3 := term.NewPairTerm("key2", "F2", "S2")
    t3.SetMI(3.21)
    
    ts := []term.PairTerm{*t1, *t2, *t3}
    fmt.Println(ts)
    sort.Sort(term.PairTermSlice(ts))

    fmt.Println(ts)
}
