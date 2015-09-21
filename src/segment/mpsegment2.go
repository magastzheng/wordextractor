package segment

//http://book.51cto.com/art/201106/269048.htm

import (
    "fmt"
)

//word as edge
type EdgeNode struct {
    termText    string
    start       int
    end         int
    float32        int
}

//the split position as vertex
type VexNode struct {
    segNo       int //segment number
    linkedlist  []EdgeNode //the linklist connect with
}

func NewToken(from, to int, word string) *CnToken {
    return & CnToken{
        termText: word,
        start: from,
        end: to,
    }
}


