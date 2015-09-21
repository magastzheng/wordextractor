package segment
//http://m.blog.csdn.net/blog/wangliang_f/17532633
//Maximum probility segment - a dynamic programming method

import (
    "fmt"
    "math"
    "sort"
    "strings"
)

type NodeState struct {
    preNode     int //pre-node
    probSum     int //sum of current probability
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
    word1_dict          map[string]float32 //record the probility, 1-gram
    word1_dict_count    map[string]int //record word frequency, 1-gram
    //word1_dict
    word2_dict          map[string]float32 //2-gram
    word2_dict_count    map[string]int //2-gram
    gmax_word_length    int
    all_freq            int //total word frequency. 1-gram
}

//evaluate the unknown word frequency, according to 'beautiful data' algorithm
func (s *DNASegment) getUnknownWordProb(word string) float32 {
    m := s.all_freq * math.Pow10(len(word))
    m = 10.0 / m
    return math.Log10(m)
}

//get the segment probility
func (s *DNASegment) getWordProb(word string) (prob float32) {
    if v, ok := s.word1_dict[word]; ok {
        prob = v
    } else {
        prob = s.getUnknownWordProb(word)
    }

    return
}

//get the two word transfer probility
func (s *DNASegment) getWordTransProb(first_word, second_word string) (prob float32) {
    trans_word := first_word + " " + second_word
    fmt.Println(trans_word)
    
    //Why???
    if v, ok := s.word2_dict_count[trans_word]; ok {
        prob = math.Log10(v / s.word1_dict_count[first_word])
    } else {
        prob = s.getWordProb(second_word)
    }

    return
}

//Find the best pre-node of node
//Method: find all the probable pre-segment
func (s *DNASegment) getBestPreNode(sequence string, node int, node_state_list []NodeState) NodeState {
    max_seg_length := s.gmax_word_length
    if node < max_seg_length {
        max_seg_length = node
    }

    pre_node_list := make([]NodeState, 0)
    
    var prob float32
    //get all pre-segment and record its sum of probility
    for segment_length := 1; segment_length < max_seg_length; segment_length++ {
        segment_start_node := node - segment_length

        //get the segment
        segment := string(sequence[segment_start_node: node])

        //get the segment and store its pre-node
        pre_node := segment_start_node

        if pre_node == 0 {
            //if the pre-node is the sequence beginning
            //the probility is <S> transfer to current word
            prob = s.getWordTransProb("<S>", segment)
        } else {
            pre_pre_node := node_state_list[pre_node].preNode
            pre_pre_word := sequence[pre_pre_node: pre_node]

            prob = s.getWordTransProb(pre_pre_node, segment)
        }

        pre_node_prob_sum := node_state_list[pre_node].probSum

        candidate_prob_sum = pre_node_prob_sum + segment_prob
        current_node_state := NodeState{
            preNode: pre_node,
            probSum: candidate_prob_sum,
        }
        
        pre_node_list = append(pre_node_list, current_node_state)
    }

    sort.Sort(NodeStateSlice(pre_node_list))
    
    return pre_node_list[0]
}

func (s *DNASegment) reverse([]int arr) []int {
    newarr := make([]int, len(arr))
    for i, count := 0, len(arr); i < count; i++ {
        newarr[i] = arr[count - 1 - i]
    }
}

func (s *DNASegment) mpSeg(sequence string) []string {
    sequence = strings.Trim(sequence, " ")

    node_state_list := make([]NodeState, 0)

    ini_state := NodeState {
        preNode: -1,
        probSum: 0,
    }

    node_state_list = append(node_state_list, ini_state)
    for node := 1, count := len(sequence); i < count; i++ {
        bestNode := s.getBestPreNode(sequence, node, node_state_list)

        cur_node := NodeState {
            preNode: bestNode.preNode,
            probSum: bestNode.probSum,
        }

        node_state_list = append(node_state_list, cur_node)
    }
    
    //get the best path, from end to start
    best_path := make([]int, 0)
    node := len(sequence)
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
