package occurrence

import (
    "term"
    "stats"
    "fmt"
    "util"
    "bytes"
    "sort"
	"segment"
    //"dict"
    //"github.com/huichen/sego"
)

const (
    FrequencyDoor = 4
)

type Occurrence struct {
   pairMap map[string]*term.PairTerm
   wordCountMap map[string] int
   totalTerm int
   totalPair int
}

func (o *Occurrence) statistics(segments []*segment.Segment) {
    for _, seg := range segments {
        //token := seg.Token()
        //key := string(token.Text())
		key := seg.Text()
        count, ok := o.wordCountMap[key]
        if ok {
            o.wordCountMap[key] = count+1
        }else{
            o.wordCountMap[key] = 1
        }

        o.totalTerm++
    }
}

func (o *Occurrence) filterSegment(segments []*segment.Segment, door int) []*segment.Segment {
    newSegs := make([]*segment.Segment, 0)
    for _, seg := range segments {
        //token := seg.Token()
        //key := string(token.Text())
		key := seg.Text()
        //if stopdict.IsContain(key) {
        //    continue
        //}

        count, ok := o.wordCountMap[key]
        if ok && count > door{
            newSegs = append(newSegs, seg)
        }
    }

    return newSegs
}

func (o *Occurrence) addPair(segments []*segment.Segment) {
    
    var first, second *segment.Segment
    for i, count := 1, len(segments); i < count; i++ {
        fmt.Println(i)
        first = segments[i - 1]
        second = segments[i]
        
        //fmt.Println("First: ", first, "Second: ", second)
        if (first.End() != second.Start()){
            continue
        }
         
        //firstKey := first.Token().Text()
        //secondKey := second.Token().Text()
		firstKey := first.Text()
        secondKey := second.Text()
        key := firstKey+secondKey
        
        if v, ok := o.pairMap[key]; ok {
            v.Increase()
            o.pairMap[key] = v 
        }else{
            pt := term.NewPairTerm(key, firstKey, secondKey)
            o.pairMap[key] = pt
        }

        o.totalPair++
    }
}

func (o *Occurrence) sort() []term.PairTerm {
    pairTerms := make([]term.PairTerm, len(o.pairMap))
    i := 0
    for _, pt := range o.pairMap {
        pairTerms[i] = *pt
        i++
    }

    sort.Sort(term.PairTermSlice(pairTerms))
    return pairTerms
}

func (o *Occurrence) AddSegments(segments []*segment.Segment, minFreq int) {
    o.statistics(segments)
    newsegs := o.filterSegment(segments, minFreq)
    //fmt.Println(newsegs)
    //o.OutputSegments(newsegs)
    //fmt.Println(newsegs[0])
    o.addPair(newsegs)
}

func (o *Occurrence) Compute() {
    
    var keyTotal, firstTotal, secondTotal int
    var ok bool

    fmt.Println("Compute: ", len(o.pairMap))
    for key, pt := range o.pairMap {
        keyTotal = pt.GetFrequency()
        if firstTotal, ok = o.wordCountMap[pt.First()]; !ok {
            firstTotal = 0
        }
        if secondTotal, ok = o.wordCountMap[pt.Second()]; !ok {
            secondTotal = 0
        }
        
        keyP := stats.Probability(keyTotal, o.totalPair)
        keyFirstP := stats.Probability(firstTotal, o.totalTerm)
        keySecondP := stats.Probability(secondTotal, o.totalTerm)
        mi := stats.CalcMI(keyP, keyFirstP, keySecondP)
        score := mi * float32(pt.GetFrequency())
        pt.SetMI(mi)
        pt.SetScore(score)
        
        fmt.Println(key, mi, score)
        o.pairMap[key] = pt
    }
}

func (o *Occurrence) GetPairTerms(score float32) []*term.PairTerm {
    pairTerms := make([]*term.PairTerm, 0)
    sortTerms := o.sort()

    for i, size := 0, len(sortTerms); i < size; i++ {
        if sortTerms[i].GetScore() > score {
            
            pairTerms = append(pairTerms, &sortTerms[i])
        }
    }
/*
    for _, t := range sortTerms {
        if t.GetScore() > score {
            fmt.Println(t)
            pairTerms = append(pairTerms, t)
        }
    }
  */  
    return pairTerms
}

func (o *Occurrence) OutputSegments(segments []*segment.Segment) {
    for i, seg := range segments {
        fmt.Print("#", i, "|")
        //if seg != nil {
            fmt.Print(seg.Start(), seg.End())
            //if seg.Token() != nil {
                fmt.Print("Text: ", seg.Text())
                fmt.Println()
            //}
        //}
        break
    }
}

func (o *Occurrence) Output() {

    outBuf := bytes.NewBufferString("Output words: \n")
    format := "Key: %v\t\tFirst: %v\t\t Second: %v\t\t Freq: %v mi: %v score: %v\n"

    terms := o.sort()
    for _, t := range terms {
        key := t.First() + t.Second()
        str := fmt.Sprintf(format, key, t.First(), t.Second(), t.GetFrequency(), t.GetMI(), t.GetScore())
        outBuf.WriteString(str)
    }

    for key, pt := range o.pairMap {
        str := fmt.Sprintf(format, key, pt.First(), pt.Second(), pt.GetMI())
        outBuf.WriteString(str)
    }

    fmt.Println(outBuf.String())
    util.WriteFile("../data/data.log", outBuf.String())
}

func NewOccurrence() *Occurrence{
    o := new(Occurrence)
    o.pairMap = make(map[string]*term.PairTerm)
    o.wordCountMap = make(map[string]int)
    o.totalTerm = 0
    o.totalPair = 0

    return o
}
