package segment

import (
    //"dict"
    "fmt"
)

//a segment in text
type Segment struct {
    //start position of segment in text
    start int

    //end position of segment in text - exclude the position
    end int

    //segment information
    //token *dict.Token
    text string
}

func (s *Segment) Start () int {
    return s.start
}

func (s *Segment) End() int {
    return s.end
}

func (s *Segment) Text() string {
    return s.text
}

func NewSegment(text string, start int, end int) *Segment {
	return &Segment{
		start: start,
		end: end,
		text: text,
	}
}

func (s *Segment) ToString() string {
    format := "%d\t%d\t%s\t%v\n"
    return fmt.Sprintf(format, s.start, s.end, s.text, []rune(s.text))
}

func GetSegmentStr(segments []*Segment) string {
    str := ""
    for _, seg := range segments {
        str += seg.ToString()
    }

    return str
}
