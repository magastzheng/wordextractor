package occurrence

import (
    "fmt"
    "segment"
    "term"
    "dict"
)

func IsPairTerm(key string, pairTerms []*term.PairTerm) bool {
    for _, term := range pairTerms {
        if term.GetKey() == key {
            return true
        }
    }

    return false
}

func FilterSegment(segments []*segment.Segment, stopdict *dict.Sign) []*segment.Segment {
    newSegs := make([]*segment.Segment, 0)
    for _, seg := range segments {
		key := seg.Text()
        if !stopdict.IsContain(key) {
            newSegs = append(newSegs, seg)
        }
    }

    return newSegs
}

func MergeSegment(segments []*segment.Segment, pairTerms []*term.PairTerm) []*segment.Segment {
    newSegments := make([]*segment.Segment, 0)
    fmt.Println("Merger start=====", len(segments))
    var first, second *segment.Segment
    i, count := 1, len(segments);
    for {
        if i >= count {
            break
        }

        first = segments[i - 1]
        second = segments[i]
        
        //fmt.Println("First: ", first, "Second: ", second)
        if (first.End() != second.Start()){
            i = i + 1
            continue
        }
         
		firstKey := first.Text()
        secondKey := second.Text()
        key := firstKey+secondKey
      
        if IsPairTerm(key, pairTerms) {
            seg := segment.NewSegment(key, first.Start(), second.End())
            newSegments = append(newSegments, seg)
            i = i + 2
        }else{
            newSegments = append(newSegments, first)
            i = i + 1
        }
    }

    return newSegments
}

func MergeOnce(segments []*segment.Segment, minFreq int, minScore float32) []*term.PairTerm {
    occur := NewOccurrence()
    occur.AddSegments(segments, minFreq)
    occur.Compute()
    pairTerms := occur.GetPairTerms(minScore)

    return pairTerms
}

func Merge(segments []*segment.Segment, minFreq int, minScore float32) []*term.PairTerm {
    pairTerms := make([]*term.PairTerm, 0)
    for {
         if minFreq < 3 {
            break
         }

         terms := MergeOnce(segments, minFreq, minScore)
         pairTerms = append(pairTerms, terms ...)
         minFreq--
         minScore = 0.8 * minScore
         segments = MergeSegment(segments, pairTerms)
    }

    return pairTerms
}
