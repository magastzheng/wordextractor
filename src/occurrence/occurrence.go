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
        //fmt.Println(i)
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

    wordProb := float32(wordFreq) / float32(o.totalPair)
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
        
        //fmt.Println("leftWordMap: ", k, tripleFreq)
        //tripleProb := float32(fripleFreq) / float32(o.totalTriple)
        //tripleLE := stats.CalcEntropy(float64(fripleFreq), float64(o.totalTriple))
        tripleProb := stats.Normalize(float32(tripleFreq), float32(o.totalTriple))
        p := float64(tripleProb / wordProb)
        entropy := -1 * p * math.Log2(p)

        le += float32(entropy)
    }

    for k, _ := range rightWordMap {
        tripleWord := word + k
        tripleFreq := 1
        if tt, ok := o.tripleMap[tripleWord]; ok {
            tripleFreq = tt.GetFrequency()
        }
        
        //fmt.Println("rightWordMap: ", k, tripleFreq, o.totalTriple)
        //tripleProb := float32(fripleFreq) / float32(o.totalTriple)
        //tripleLE := stats.CalcEntropy(float64(fripleFreq), float64(o.totalTriple))
        tripleProb := stats.Normalize(float32(tripleFreq), float32(o.totalTriple))
        p := float64(tripleProb / wordProb)
        entropy := -1 * p * math.Log2(p)
        
        //fmt.Println(entropy)
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

func (o *Occurrence) calcScore(pt *term.PairTerm) float32 {
    multipier := float32(len(pt.GetKey())) * 0.8

    //score := pt.GetMI() * float32(pt.GetFrequency()) * float32(multipier)
    score := pt.GetMI() + pt.GetLE() + pt.GetRE()    
    score = score * multipier
    //fmt.Println("m: ", multipier, " mi: ", pt.GetMI(), " score: ", score)

    return score
}

func (o *Occurrence) AddSegments(segments []*segment.Segment, minFreq int) {
    o.statistics(segments)
    newsegs := o.filterSegment(segments, minFreq)
    //fmt.Println(newsegs)
    //o.OutputSegments(newsegs)
    //fmt.Println(newsegs[0])
    o.addPair(newsegs)
    o.addTriple(newsegs)
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
        pt.SetMI(mi)
        
        le, re := o.computeEntropy(key)
        pt.SetLE(le)
        pt.SetRE(re)

        score := o.calcScore(pt)
        pt.SetScore(score)
        
        //fmt.Println(key, mi, le, re, score)
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
