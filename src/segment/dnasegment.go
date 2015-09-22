package segment
//http://m.blog.csdn.net/blog/wangliang_f/17532633
//Maximum probility segment - a dynamic programming method

import (
    "fmt"
    "math"
    "sort"
    "strings"
    "os"
    "log"
    "bufio"
    "strconv"
)

type NodeState struct {
    preNode     int //pre-node
    probSum     float64 //sum of current probability
}

type NodeStateSlice []NodeState

func (n NodeStateSlice) Len() int {
    return len(n)
}

func (n NodeStateSlice) Swap(i, j int) {
    n[i], n[j] = n[j], n[i]
}

func (n NodeStateSlice) Less(i, j int) bool {
    return n[i].probSum > n[j].probSum
}


type DNASegment struct {
    word1_dict          map[string]float64 //record the probility, 1-gram
    word1_dict_count    map[string]int //record word frequency, 1-gram
    //word2_dict          map[string]float64 //2-gram
    //word2_dict_count    map[string]int //2-gram
    gmax_word_length    int
    all_freq            int //total word frequency. 1-gram
}

//evaluate the unknown word frequency, according to 'beautiful data' algorithm
func (s *DNASegment) getUnknownWordProb(word string) float64 {
    m := float64(s.all_freq) * math.Pow10(len(word))
    m = 10.0 / m
    //fmt.Println(m)
    //return math.Log10(m)
    return m
}

//get the segment probility
func (s *DNASegment) getWordProb(word string) (prob float64) {
    if v, ok := s.word1_dict[word]; ok {
        prob = v
    } else {
        prob = s.getUnknownWordProb(word)
    }

    return
}

//get the two word transfer probility
//func (s *DNASegment) getWordTransProb(first_word, second_word string) (prob float64) {
//    trans_word := first_word + " " + second_word
//    fmt.Println(trans_word)
    
    //Why???
//    if v, ok := s.word2_dict_count[trans_word]; ok {
//        result := float64(v / s.word1_dict_count[first_word])
//        prob = float64(math.Log10(result))
//    } else {
//        prob = s.getWordProb(second_word)
//    }

//    return
//}

//Find the best pre-node of node
//Method: find all the probable pre-segment
func (s *DNASegment) getBestPreNode(sequence string, node int, node_state_list []NodeState) NodeState {
    max_seg_length := s.gmax_word_length
    if node < max_seg_length {
        max_seg_length = node
    }
    
    fmt.Println("max_seg_length: ", max_seg_length, " node: ", node)

    pre_node_list := make([]NodeState, 0)
     
    var prob float64
    //get all pre-segment and record its sum of probility
    for segment_length := 1; segment_length < max_seg_length + 1; segment_length++ {
        //fmt.Println("Current: ", segment_length, " Total: ", max_seg_length)
        segment_start_node := node - segment_length

        //get the segment
        segment := string(sequence[segment_start_node: node])

        //get the segment and store its pre-node
        pre_node := segment_start_node

        //if pre_node == 0 {
            //if the pre-node is the sequence beginning
            //the probility is <S> transfer to current word
        //    prob = s.getWordTransProb("<S>", segment)
        //} else {
        //    pre_pre_node := node_state_list[pre_node].preNode
        //    pre_pre_word := string(sequence[pre_pre_node: pre_node])

        //    prob = s.getWordTransProb(pre_pre_word, segment)
        //}
        
        if v, ok := s.word1_dict[segment]; ok {
            prob = v
        } else {
            prob = s.getUnknownWordProb(segment)
        }
        
        pre_node_prob_sum := node_state_list[pre_node].probSum
        
        
        //the cumulative probility
        candidate_prob_sum := pre_node_prob_sum + prob
        current_node_state := NodeState{
            preNode: pre_node,
            probSum: candidate_prob_sum,
        }
        
        pre_node_list = append(pre_node_list, current_node_state)
    }
    
    sort.Sort(NodeStateSlice(pre_node_list))
    fmt.Println(pre_node_list) 
    return pre_node_list[0]
}

func (s *DNASegment) reverse(arr []int) []int {
    newarr := make([]int, len(arr))

    for i, count := 0, len(arr); i < count; i++ {
        newarr[i] = arr[count - 1 - i]
    }

    return newarr
}

func (s *DNASegment) MPSeg(sequence string) []string {
    sequence = strings.Trim(sequence, " ")

    node_state_list := make([]NodeState, 0)

    ini_state := NodeState {
        preNode: -1,
        probSum: 0,
    }

    node_state_list = append(node_state_list, ini_state)
    for node, count := 1, len(sequence); node < count; node++ {
        //fmt.Println("MPSeg: ", node, " Total: ", count)
        bestNode := s.getBestPreNode(sequence, node, node_state_list)

        cur_node := NodeState {
            preNode: bestNode.preNode,
            probSum: bestNode.probSum,
        }

        node_state_list = append(node_state_list, cur_node)
    }
    
    fmt.Println(node_state_list)
    //get the best path, from end to start
    best_path := make([]int, 0)
    node := len(sequence) - 1
    best_path = append(best_path, node)
    for {
        pre_node := node_state_list[node].preNode
        if pre_node == -1 {
            break
        }

        node = pre_node
        best_path = append(best_path, node)
    }

    new_best_path := s.reverse(best_path)

    //build the split
    word_list := make([]string, 0)
    for i, count := 0, len(new_best_path); i < count - 1; i++ {
        left := new_best_path[i]
        right := new_best_path[i+1]
        word := string(sequence[left:right])
        word_list = append(word_list, word)
    }

    return word_list
}

func (s *DNASegment) InitDict(filename string) {
    file, err := os.Open(filename)
    defer file.Close()
    if err != nil {
        fmt.Println("Fail to load dictionary file: %s", filename)
        log.Fatalf("Fail to load dictionary file: %s", filename)
    }
    
    reader := bufio.NewReader(file)
    var text string
    var freqText string
    var frequency int
    var nature string

    for {
        size, _ := fmt.Fscanln(reader, &text, &freqText, &nature)

        if size == 0 {
            //end of file
            break
        } else if size < 2 {
            //invalid line
            continue
        } else if size == 2 {
            //no word nature, set it empty
            nature = ""
        }

        frequency, err = strconv.Atoi(freqText)
        if err != nil {
            continue
        }

        //if frequency < minTokenFrequency {
        //    continue
        //}

        //words := splitTextToWords([]byte(text))
        //token := Token{text: words, frequency: frequency, nature: nature}
        //dict.addToken(&token)
        s.word1_dict_count[text] = frequency
        s.all_freq += frequency
        if s.gmax_word_length < len(text) {
            s.gmax_word_length = len(text)
        }
    }
    
    for k, v := range s.word1_dict_count {
        //fmt.Println(v, s.all_freq)
        prob := float64(v) / float64(s.all_freq)
        //get Log2
        //prob = float64(math.Log(float64(prob)))
        //fmt.Printf("%f: %f\n", prob, math.Log(float64(prob)))
        if _, ok := s.word1_dict[k]; !ok {
            s.word1_dict[k] = prob
        }
    }
    
    //fmt.Println(s.word1_dict_count)
    //fmt.Println(s.word1_dict)
    //fmt.Println(s.gmax_word_length)
    //fmt.Println(s.all_freq)
    log.Println("Finish to load dict")
}

func NewDNASegment() *DNASegment {
    s := new(DNASegment)
    s.word1_dict = make(map[string]float64)
    s.word1_dict_count = make(map[string]int)
    s.gmax_word_length = 0
    s.all_freq = 0

    return s
}
