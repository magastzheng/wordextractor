package segment

import (
    "fmt"
    "math"
)

type MPSegment struct {
    wordProbMap     map[string]float32
    wordCountMap    map[string]int
    maxWordLength   int
    totalFrequency  int
}

type NodeState struct {
    preNode     int //pre-node
    probSum     int //sum of current probability
}

func (s *MPSegment) getUnkownWordProb(word string) float32 {
    m := s.totalFrequency * math.Pow10(len(word))
    return math.Log(10.0 / m)
}

func (s *MPSegment) Max(x, y int) int {
    if x < y {
        return y
    }

    return x
}

func (s *MPSegment) GetBestPreNode(sequence string, node int, nodeStatList []NodeState) {
    
}  
