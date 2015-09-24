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
    score float32
}

type PairTermSlice []PairTerm

func (t *Term) GetKey() string {
    return t.key
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
        score: 0.0,
    }
}

func GetPairTermStr(pairTerms []*PairTerm) string {
    format := "%s\t%s\t%s\t%d\n"
    str := ""
    for _, pt := range pairTerms {
        s := fmt.Sprintf(format, pt.key, pt.first, pt.second, pt.frequency)
        str += s
        //fmt.Print(s)
    }

    return str
}
