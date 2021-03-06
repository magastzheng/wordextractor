package segment

import (
    "fmt"
    "math"
    "dict"
    "unicode"
    //"util"
)

//const (
//    MaxWordLength = 20
//)

type Candidate struct {
    Start       int
    Length      int
    BestPrev    int
    Fee         float32
    SumFee      float32
    Word        string
    Freq        int
}

func (c *Candidate) ToString() string {
    format := "%s:\t\t%d\t\t%d\t\t%d\t\t%f\t\t%f\n"
    return fmt.Sprintf(format, c.Word, c.Start, c.Length, c.BestPrev, c.Fee, c.SumFee)
}

func Output(vec_cd []*Candidate) string {
    out := ""
    for _, c := range vec_cd {
        out += c.ToString()
    }

    return out
}

func IsLetter(ch rune) bool {
    return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func getFee(freq, totalFreq int) float32 {
    normal := float64(freq * 1 + 1)/ float64(totalFreq)
    fee := float32(-math.Log2(float64(normal)))
    return fee
}

func newCandidate(word string, start int, length int, fee float32, freq int) *Candidate {
    cand := &Candidate {
        Start: start,
        Length: length,
        Word: word,
        Fee: fee,
        SumFee: 0.0,
        Freq: freq,
    }

    return cand
}

func getTempWords(runeBuf []rune, d *dict.Dictionary) []*Candidate {
    freq := 0
    //runeBuf := []rune(sequence)
    totalLen := len(runeBuf)
    word := ""
    vec_cd := make([]*Candidate, 0)
    
    totalFreq := d.FreqAll()
    isNumber := false
    numberStart := 0
    isLetter := false
    letterStart := 0
    for i := 0; i < totalLen; i++ {
        ch := runeBuf[i]
        if unicode.IsDigit(ch) {
            //handle digit
            if !isNumber {
                numberStart = i
                isNumber = true
            }
            
            end := i + 1
            isTail := false
            if i + 1 < totalLen {
                if !unicode.IsDigit(runeBuf[end]) {
                    isTail = true
                }
            } else {
                isTail = true
            }

            if isTail {
                isNumber = false
                word = string(runeBuf[numberStart: end])
                fee := getFee(0, totalFreq)
                cand := newCandidate(word, numberStart, end - numberStart, fee, freq)
                vec_cd = append(vec_cd, cand)
            }
        } else if IsLetter(ch) {
            //handle English characters
            if !isLetter {
                letterStart = i
                isLetter = true
            }
            
            end := i + 1
            isTail := false
            if i + 1 < totalLen {
                if !IsLetter(runeBuf[end]) {
                    isTail = true
                }
            } else {
                isTail = true
            }

            if isTail {
                isLetter = false
                word = string(runeBuf[letterStart: end])
                fee := getFee(0, totalFreq)
                cand := newCandidate(word, letterStart, end - letterStart, fee, freq)
                vec_cd = append(vec_cd, cand)
            }
        } else {
            for length := 1; length < dict.MaxWordLength && i + length <= totalLen; length++ {
                //fmt.Println(i, length)
                word = string(runeBuf[i: i+length])
                freq = d.FindWord(word)
            
                if length > 1 && freq == -1 {
                    //More than one character and cannot find in dictionary
                    //don't sign in
                    continue
                }

                if freq == -1 {
                    //single character and cannot find in dictionary
                    freq = 0
                }
            
                fee := getFee(freq, totalFreq)
                cand := newCandidate(word, i, length, fee, freq)
                vec_cd = append(vec_cd, cand)
            }
        }
    }
    
    return vec_cd
}

func getPrev(vec_cd []*Candidate) {
    min_id := -1  //best preview word number
    j := -1
    size := len(vec_cd)
    for i := 0; i < size; i++ {
        if vec_cd[i].Start == 0 {
            //The first character
            vec_cd[i].BestPrev = -1
            vec_cd[i].SumFee = vec_cd[i].Fee
        }else{
            min_id = -1

            for j = i - 1; j >= 0; j-- {
                //Find all the preview word in the left
                if vec_cd[j].Start + vec_cd[j].Length == vec_cd[i].Start {
                    if min_id == -1 || vec_cd[j].SumFee < vec_cd[min_id].SumFee {
                        min_id = j
                    }
                }
            }
            
            //Store the best preview word number
            vec_cd[i].BestPrev = min_id
            //store the minimum cumulative fee
            vec_cd[i].SumFee = vec_cd[i].Fee + vec_cd[min_id].SumFee
        }
    }

    //s := Output(vec_cd)
    //util.WriteFile("../data/origin_all_in_getprev.txt", s)
}

func reverse(source []*Candidate) []*Candidate {
    size := len(source)
    dest := make([]*Candidate, size)
    for i := 0; i < size; i++ {
        dest[i] = source[size - 1 - i]
    }

    return dest
}

func outputSegment(segments []*Candidate) {
    out := ""
    for _, s := range segments {
        out += s.Word + "||"
    }

    fmt.Println(out)
}

func SegmentSentence_MP(sequence string, d *dict.Dictionary) string {
    in := []rune(sequence)
    runeLen := len(in)
    min_id := -1
    
    //get all candidate
    vec_cd := getTempWords(in, d)
    
    //s := Output(vec_cd)
    //util.WriteFile("../data/origin_all_can.txt", s)

    getPrev(vec_cd)
    
    //s = Output(vec_cd)
    //util.WriteFile("../data/origin_all_after_getprev.txt", s)

    size := len(vec_cd)
    for i := 0; i < size; i++ {
        if vec_cd[i].Start + vec_cd[i].Length == runeLen {
            //The current word is the tail of sequence
            if min_id == -1 || vec_cd[i].SumFee < vec_cd[min_id].SumFee {
                
                min_id = i;
            }
        }
    }
    
    source := make([]*Candidate, 0)
    out := ""
    for i := min_id; i >= 0; i = vec_cd[i].BestPrev {
        start := vec_cd[i].Start
        end := start + vec_cd[i].Length
        source = append(source, vec_cd[i])
        out = string(in[start: end]) + "::" + out
    }
    //dest := reverse(source)

    //fmt.Println(dest)
    //outputSegment(dest)
    return out
}

func SegmentSentenceMP(buf []rune, pos int, d *dict.Dictionary) []*Segment {
    runeLen := len(buf)
    min_id := -1
    
    //get all candidate
    vec_cd := getTempWords(buf, d)
    
    //s := Output(vec_cd)
    //util.WriteFile("../data/origin_all_can.log", s)

    getPrev(vec_cd)
    
    //s = Output(vec_cd)
    //util.WriteFile("../data/origin_all_after_getprev.log", s)

    size := len(vec_cd)
    for i := 0; i < size; i++ {
        if vec_cd[i].Start + vec_cd[i].Length == runeLen {
            //The current word is the tail of sequence
            if min_id == -1 || vec_cd[i].SumFee < vec_cd[min_id].SumFee {
                
                min_id = i;
            }
        }
    }
    
    source := make([]*Candidate, 0)
    for i := min_id; i >= 0; i = vec_cd[i].BestPrev {
        source = append(source, vec_cd[i])
    }
    dest := reverse(source)
    
    segments := make([]*Segment, len(dest))
    for i := 0; i < len(dest); i++ {
        start := dest[i].Start + pos
        end := start + dest[i].Length
        word := dest[i].Word

        seg := &Segment{
            start: start,
            end: end,
            text: word,
        }

        segments[i] = seg
    }

    return segments
}
