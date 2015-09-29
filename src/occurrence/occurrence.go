package occurrence

import (
    "term"
    "stats"
    "fmt"
    "util"
    "bytes"
    "sort"
	"segment"
    "strings"
    "math"
    //"dict"
    //"github.com/huichen/sego"
)

const (
    FrequencyDoor = 4
)

type Occurrence struct {
   pairMap map[string]*term.PairTerm
   tripleMap map[string]*term.TripleTerm
   wordCountMap map[string] int
   totalTerm int
   totalPair int
   totalTriple int
}

func (o *Occurrence) statistics(segments []*segment.Segment) {
    for _, seg := range segments {
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
		key := seg.Text()
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
        first = segments[i - 1]
        second = segments[i]
        
        if (first.End() != second.Start()){
            continue
        }

		firstKey := first.Text()
        secondKey := second.Text()
        key := firstKey+secondKey

        if t, ok := o.pairMap[key]; ok {
            t.Increase()
            o.pairMap[key] = t 
        }else{
            pt := term.NewPairTerm(key, firstKey, secondKey)
            o.pairMap[key] = pt
        }

        o.totalPair++
    }
}

func (o *Occurrence) addTriple(segments []*segment.Segment) {
    var first, second, third *segment.Segment
    for i, count := 2, len(segments); i < count; i++ {
        first = segments[i - 2]
        second = segments[i - 1]
        third = segments[i]

        if first.End() != second.Start() && second.End() != third.Start() {
            continue
        }
        
        firstKey := first.Text()
        secondKey := second.Text()
        thirdKey := third.Text()
        key := firstKey+secondKey+thirdKey
        
        if t, ok := o.tripleMap[key]; ok {
            t.Increase()
            o.tripleMap[key] = t 
        }else{
            t := term.NewTripleTerm(key, firstKey, secondKey, thirdKey)
            o.tripleMap[key] = t
        }

        o.totalTriple++
    }
}

func (o *Occurrence) computeEntropy(word string) (le, re float32) {
    leftWordMap := make(map[string]int)
    rightWordMap := make(map[string]int)
    le = 0.0
    re = 0.0

    wordFreq := 1
    pt, yes := o.pairMap[word]
    if yes {
        wordFreq = pt.GetFrequency()
    }

    wordProb := stats.Probability(wordFreq, o.totalPair)
    for k, t := range o.tripleMap {
        pos := strings.Index(k, word)
        if pos == 0 {
            //left, right word is the third one
            right := t.Third()
            if _, ok := rightWordMap[right]; !ok {
                rightWordMap[right] = 1
            } else {
                rightWordMap[right] = rightWordMap[right] + 1
            }
            
        } else if pos > 0 {
            //right, left word is the first one
            left := t.First()
            if _, ok := leftWordMap[left]; !ok {
                leftWordMap[left] = 1
            } else {
                leftWordMap[left] = leftWordMap[left] + 1
            }

        } else {
            //nothing
        }
    }
    
    
    for k, _ := range leftWordMap {
        tripleWord := k + word
        tripleFreq := 1
        if tt, ok := o.tripleMap[tripleWord]; ok {
            tripleFreq = tt.GetFrequency()
        }
        
        tripleProb := stats.Normalize(float64(tripleFreq), float64(o.totalTriple))
        p := stats.Probability64(tripleProb, wordProb)
        entropy := -1 * p * math.Log2(p)

        le += float32(entropy)
    }

    for k, _ := range rightWordMap {
        tripleWord := word + k
        tripleFreq := 1
        if tt, ok := o.tripleMap[tripleWord]; ok {
            tripleFreq = tt.GetFrequency()
        }
        
        tripleProb := stats.Normalize(float64(tripleFreq), float64(o.totalTriple))
        p := stats.Probability64(tripleProb, wordProb)
        entropy := -1 * p * math.Log2(p)
        
        re += float32(entropy)
    }
    
    //format := "Key: %s\t\t wordProb: %f\t le: %f\t re: %f \n"
    //s := fmt.Sprintf(format, word, wordProb, le, re)
    //fmt.Println(s)

    return
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

func (o *Occurrence) calcScore(times int, pt *term.PairTerm) float32 {
    multipier := 1.5 * float32(pt.Length()) 
    multipier *= float32(pt.GetFrequency())

    //score := pt.GetMI() * float32(pt.GetFrequency()) * float32(multipier)
    score := 2.5 * pt.GetMI() + pt.GetLE() + pt.GetRE()    
    score = score * multipier * float32(times)
    //fmt.Println(pt.GetKey(), "multipier: ", multipier, " mi: ", pt.GetMI(), "LE: ", pt.GetLE(), "RE: ", pt.GetRE(), " score: ", score, " times: ", times)

    return score
}

func (o *Occurrence) AddSegments(segments []*segment.Segment, minFreq int) {
    o.statistics(segments)
    newsegs := o.filterSegment(segments, minFreq)
    o.addPair(newsegs)
    o.addTriple(newsegs)
}

func (o *Occurrence) Compute(times int) {
    
    var keyTotal, firstTotal, secondTotal int
    var totalMI, totalLE, totalRE float32
    var ok bool

    //fmt.Println("Compute: ", len(o.pairMap))
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
        pt.SetMI(float32(mi))
        
        le, re := o.computeEntropy(key)
        pt.SetLE(le)
        pt.SetRE(re)
        
        totalMI += float32(mi)
        totalLE += le
        totalRE += re

        //fmt.Println(key, keyTotal,  o.totalPair, o.totalTerm, mi, le, re)
        o.pairMap[key] = pt
    }
    
    //Normalize
    for key, pt := range o.pairMap {
        mi := pt.GetMI() / totalMI
        le := pt.GetLE() / totalLE
        re := pt.GetRE() / totalRE
        
        pt.SetMI(mi)
        pt.SetLE(le)
        pt.SetRE(re)
        score := o.calcScore(times, pt)
        pt.SetScore(score)

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
    format := "Key: %v\t\tFirst: %v\t\t Second: %v\t\t Freq: %v mi: %f le: %f re: %f score: %v\n"

    terms := o.sort()
    for _, t := range terms {
        key := t.First() + t.Second()
        str := fmt.Sprintf(format, key, t.First(), t.Second(), t.GetFrequency(), t.GetMI(), t.GetLE(), t.GetRE(), t.GetScore())
        outBuf.WriteString(str)
    }
    
    outBuf.WriteString("========================PairMap===================\n")
    for key, pt := range o.pairMap {
        str := fmt.Sprintf(format, key, pt.First(), pt.Second(), pt.GetFrequency(), pt.GetMI(), pt.GetLE(), pt.GetRE(), pt.GetScore())
        outBuf.WriteString(str)
    }

    //fmt.Println(outBuf.String())
    util.WriteFile("../data/data.log", outBuf.String())
}

func NewOccurrence() *Occurrence{
    o := new(Occurrence)
    o.pairMap = make(map[string]*term.PairTerm)
    o.tripleMap = make(map[string]*term.TripleTerm)
    o.wordCountMap = make(map[string]int)

    o.totalTerm = 0
    o.totalPair = 0
    o.totalTriple = 0

    return o
}
