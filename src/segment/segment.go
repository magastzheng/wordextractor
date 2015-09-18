package segment

import (
    "dict"
)

//a segment in text
type Segment struct {
    //start position of segment in text
    start int

    //end position of segment in text - exclude the position
    end int

    //segment information
    token *dict.Token
}

func (s *Segment) Start () int {
    return s.start
}

func (s *Segment) End() int {
    return s.end
}

func (s *Segment) Token() *dict.Token {
    return s.token
}
