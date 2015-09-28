package stats

import (
    "math"
    //"fmt"
)

//Mutual Information
//MI(X,Y)=log2{P(X,Y)/[P(X) * P(Y)]}
func CalcMI(xy, x, y float64) float64 {
   z := xy / (x * y)
   return math.Log2(z)
}

//
func CalcEntropy(frequency, totalFrequency float64) float64 {
    p := frequency / totalFrequency
    return -1 * p * math.Log2(p)
}

//Make value normalization
func Normalize(value, totalValue float64) float64 {
    return value / totalValue
}

//Calculate the probability
func Probability(freq, totalFreq int) float64 {
    return float64(freq / totalFreq)
}

