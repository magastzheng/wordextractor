package stats_test

import (
    "stats"
    "fmt"
    "testing"
    "math"
)

func Test_Calcmi(t *testing.T) {
    res := stats.Calcmi(4.0, 1.0, 1.0)
    fmt.Println(res)
    fmt.Println(math.Log2(4.0))
}

func Test_CalcEntropy(t *testing.T) {
    res := stats.CalcEntropy(125, 15435)
    fmt.Println(res)
}
