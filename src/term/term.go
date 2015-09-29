package term

import (
    "fmt"
)

type Term struct {
    key string
    frequency int
    //start int
    //end int
}

type PairTerm struct {
    Term
    first string
    second string
    mi float32
    le float32
    re float32
    score float32
}

type TripleTerm struct {
    PairTerm
    third string
}

type PairTermSlice []PairTerm
type PairTermPtrSlice []*PairTerm

func (t *Term) GetKey() string {
    return t.key
}

func (t *Term) Length() int {
    return len([]rune(t.key))
}

//set the value as newValue and return old value
func (t *Term) SetFrequency(newValue int) int {
    oldValue := t.frequency
    t.frequency = newValue
    return oldValue
}

//func (t *Term) GetStart() int {
//    return t.start
//}

//func (t *Term) GetEnd() int {
//    return t.end
//}

func (t *Term) GetFrequency() int {
    return t.frequency
}

func (t *Term) Increase() {
    t.frequency++
}

func (pt *PairTerm) First() string {
    return pt.first
}

func (pt *PairTerm) Second() string {
    return pt.second
}

func (pt *PairTerm) GetMI() float32 {
    return pt.mi
}

func (pt *PairTerm) SetMI(mi float32) {
    pt.mi = mi
}

func (pt *PairTerm) GetLE() float32 {
    return pt.le
}

func (pt *PairTerm) SetLE(le float32) {
    pt.le = le
}

func (pt *PairTerm) GetRE() float32 {
    return pt.re
}

func (pt *PairTerm) SetRE(re float32) {
    pt.re = re
}

func (pt *PairTerm) GetScore() float32 {
    return pt.score
}

func (pt *PairTerm) SetScore(score float32)  {
    pt.score = score
}

func (pts PairTermSlice) Len() int {
    return len(pts)
}

func (pts PairTermSlice) Swap(i, j int) {
    pts[i], pts[j] = pts[j], pts[i]
}

func (pts PairTermSlice) Less(i, j int) bool {
    //return pts[j].GetMI() < pts[i].GetMI()
    
    return pts[j].GetScore() < pts[i].GetScore()
}


func (ptps PairTermPtrSlice) Len() int {
    return len(ptps)
}

func (ptps PairTermPtrSlice) Swap(i, j int) {
    ptps[i], ptps[j] = ptps[j], ptps[i]
}

func (ptps PairTermPtrSlice) Less(i, j int) bool {
    //return pts[j].GetMI() < pts[i].GetMI()
    
    return ptps[j].GetScore() < ptps[i].GetScore()
}

func (tp *TripleTerm) Third() string {
    return tp.third
}

func NewTerm(key string) *Term {
    return &Term {
        key: key,
        frequency: 1,
    }
}

func NewPairTerm(key string, first, second string) *PairTerm {
    t := NewTerm(key)
    return &PairTerm {
        Term: *t,
        first: first,
        second: second,
        mi: 0.0,
        le: 0.0,
        re: 0.0,
        score: 0.0,
    }
}

func NewTripleTerm(key string, first, second, third string) *TripleTerm {
    pt := NewPairTerm(key, first, second)
    return &TripleTerm {
        PairTerm: *pt,
        third: third,
    }
}

func GetPairTermStr(pairTerms []*PairTerm) string {
    format := "Key:%s\t First:%s\t Second: %s\t Freq: %d\t Len: %d\t MI: %f\t LE:%f\t RE: %f\t Score: %f\n"
    str := ""
    for _, pt := range pairTerms {
        s := fmt.Sprintf(format, pt.key, pt.first, pt.second, pt.frequency, pt.Length(), pt.mi, pt.le, pt.re, pt.score)

        str += s
        //fmt.Print(s)
    }

    return str
}
