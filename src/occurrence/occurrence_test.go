package occurrence_test

import (
    "occurrence"
    "testing"
    "io/ioutil"
    "util"
    "fmt"
    "bytes"
    "segment"
	"dict"
)

func Test_Occurrence_Compute(t *testing.T){
    //var segmenter sego.Segmenter
    //segmenter.LoadDictionary("C:/Go/thirdpartlib/src/github.com/huichen/sego/data/dictionary.txt")
    filename := "../data/testdata/125.txt" 

    buf, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
	
	sign := dict.NewSign("../data/dictionary/sign.txt")
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")
	//stop := dict.NewSign("../data/dictionary/stopwords.txt.txt")
    //segments := segmenter.Segment(buf)
	segments := segment.SegmentDoc(string(buf), sign, d)
    fmt.Println(len(segments))
    LogSegments(segments)
    occur := occurrence.NewOccurrence()
    occur.AddSegments(segments, 3)
    occur.Compute(1)
    occur.Output()
}


func LogSegments(segments []*segment.Segment) {
    format := "%d %d %s\t"
    outBuf := bytes.NewBufferString("Output words: \n")
    for _, v := range segments {
        //if len(v) == 12 {
        //    fmt.Println(v)
        //}
        //token := v.Token()
        //binary.Write(outBuf, binary.BigEndian, v.Start())
        //binary.Write(outBuf, binary.BigEndian, v.End())
        //fmt.Println(v.Start())
        //binary.BigEndian.PutUint32(outBuf, uint32(v.Start()))
        //binary.BigEndian.PutUint32(outBuf, uint32(v.End()))
        
        str := fmt.Sprintf(format, v.Start(), v.End(), v.Text())
        outBuf.WriteString(str)
        //outBuf.Write(Int32ToBytes(int32(v.Start())))
        //outBuf.WriteByte('\t')
        //outBuf.Write(Int32ToBytes(int32(v.End())))
        //outBuf.WriteByte('\t')
        //outBuf.Write(v.Start())
        //outBuf.Write(v.End())
        //outBuf.WriteString(token.Text())
        //outBuf.Write(token.Frequency())
        //binary.Write(outBuf, binary.BigEndian, token.Frequency())
        //outBuf.WriteByte('\t') 
       // outBuf.Write(Int32ToBytes(int32(token.Frequency())))
        //outBuf.WriteByte('\t')
        //outBuf.WriteString(token.Pos())
        //outBuf.WriteByte('\n')
    }

    util.WriteFile("../data/segment.log", outBuf.String())
}
