package occurrence

import (
    //"fmt"
    "segment"
    "term"
    "dict"
    "sort"
	"log"
    //"unicode"
    "regexp"
)

const (
    CNSpace = string(rune(12288))
)

func IsPairTerm(key string, pairTerms []*term.PairTerm) bool {
    for _, term := range pairTerms {
        if term.GetKey() == key {
            return true
        }
    }

    return false
}

func IsSpace(key string) bool {
    if key == CNSpace {
        return true
    } else if isSpace, _ := regexp.MatchString("[\\s]+", key); isSpace {
        return true
    }

    return false
}

func FilterSegment(segments []*segment.Segment, stopdict *dict.Sign) []*segment.Segment {
    newSegs := make([]*segment.Segment, 0)
    for _, seg := range segments {
		key := seg.Text()
        
        if IsSpace(key) {
            continue
        }

        match, _ := regexp.MatchString("[\\d]+", key)
        if match {
            continue
        }

        if !stopdict.IsContain(key) {
            newSegs = append(newSegs, seg)
        }
    }

    return newSegs
}

func MergeSegment(segments []*segment.Segment, pairTerms []*term.PairTerm) []*segment.Segment {
    newSegments := make([]*segment.Segment, 0)

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

func MergeOnce(segments []*segment.Segment, minFreq int, minScore float32, times int) []*term.PairTerm {
    occur := NewOccurrence()
    occur.AddSegments(segments, minFreq)
    occur.Compute(times)
    pairTerms := occur.GetPairTerms(minScore)

    return pairTerms
}

func Merge(segments []*segment.Segment, minFreq int, minScore float32) []*term.PairTerm {
	log.Printf("开始合并词语...")
    pairTerms := make([]*term.PairTerm, 0)
    times := 1
    for {
         terms := MergeOnce(segments, minFreq, minScore, times)
         //fmt.Println("Each merge total terms: ", len(terms))
         pairTerms = append(pairTerms, terms ...)
         minFreq--
         minScore = 0.8 * minScore
         times *= 2   
         segments = MergeSegment(segments, pairTerms)
		 
		 if minFreq < 2 || len(pairTerms) == 0 {
            break
         }
    }
    
    sort.Sort(term.PairTermPtrSlice(pairTerms))
	
	log.Printf("完成合并词语，总数: %d", len(pairTerms))
	
	/*format := "%s,%d,%f\n"
    str := "短语,频率,分数\n"
    for _, pt := range pairTerms {
        //if pt.GetFrequency() >= 1{
            str += fmt.Sprintf(format, pt.GetKey(), pt.GetFrequency(), pt.GetScore())
        //}
    }
	
	log.Println(str)
	*/
	
    return pairTerms
}

