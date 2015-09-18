package dict

import (
    "bytes"
)

type Text []byte

type Token struct {
    text []Text
    frequency int
    nature string
}

func (t *Token) Text() string {
    
    var output string
    for _, w := range t.text {
        output += string(w)
    }

    return output
}

func (t *Token) Frequency() int {
    return t.frequency
}

func (t *Token) Nature() string {
    return t.nature
}

////=================Dictionary================

type Dictionary struct {
    root            node
    maxTokenLen     int
    numTokens       int
    tokens          []*Token
    totalFrequency  int64
}

//trie node
type node struct {
    word        Text
    token       *Token
    children    []*node
}

func (d *Dictionary) MaxTokenLength() int {
    return d.maxTokenLen
}

func (d *Dictionary) NumTokens() int {
    return d.numTokens
}

func (d *Dictionary) TotalFrequency() int64 {
    return d.totalFrequency
}

func (d *Dictionary) addToken(token *Token) {
    current := &d.root

    for _, word := range token.text {
        current = upsert(&current.children, word)
    }
    
    //add the token if there is no added
    if current.token == nil {
        current.token = token
        if len(token.text) > d.maxTokenLen {
            d.maxTokenLen = len(token.text)
        }
        d.numTokens++
        d.tokens = append(d.tokens, token)
        d.totalFrequency += int64(token.frequency)
    }
}

//find all token can match word
func (d *Dictionary) lookupTokens(words []Text, tokens []*Token) int {
	if len(words) == 0 {
		return 0
	}

	current := &d.root
	numTokens := 0
	for _, word := range words {
		// break it go to leaf
		if len(current.children) == 0 {
			break
		}

		// lookup the next token in the children
		index, found := binarySearch(current.children, word)
		if !found {
			break
		}

		// go into the children if matching
		current = current.children[index]
		if current.token != nil {
			tokens[numTokens] = current.token
			numTokens++
		}
	}
	return numTokens
}

// 二分法查找字元在子节点中的位置
// 如果查找成功，第一个返回参数为找到的位置，第二个返回参数为true
// 如果查找失败，第一个返回参数为应当插入的位置，第二个返回参数false
func binarySearch(nodes []*node, word Text) (int, bool) {
	start := 0
	end := len(nodes) - 1

	// 特例：
	if len(nodes) == 0 {
		// 当slice为空时，插入第一位置
		return 0, false
	}
	compareWithFirstWord := bytes.Compare(word, nodes[0].word)
	if compareWithFirstWord < 0 {
		// 当要查找的元素小于首元素时，插入第一位置
		return 0, false
	} else if compareWithFirstWord == 0 {
		// 当首元素等于node时
		return 0, true
	}
	compareWithLastWord := bytes.Compare(word, nodes[end].word)
	if compareWithLastWord == 0 {
		// 当尾元素等于node时
		return end, true
	} else if compareWithLastWord > 0 {
		// 当尾元素小于node时
		return end + 1, false
	}

	// 二分
	current := end / 2
	for end-start > 1 {
		compareWithCurrentWord := bytes.Compare(word, nodes[current].word)
		if compareWithCurrentWord == 0 {
			return current, true
		} else if compareWithCurrentWord < 0 {
			end = current
			current = (start + current) / 2
		} else {
			start = current
			current = (current + end) / 2
		}
	}
	return end, false
}

// 将字元加入节点数组中，并返回插入的节点指针
// 如果字元已经存在则返回存在的节点指针
func upsert(nodes *[]*node, word Text) *node {
	index, found := binarySearch(*nodes, word)
	if found {
		return (*nodes)[index]
	}
	*nodes = append(*nodes, nil)
	copy((*nodes)[index+1:], (*nodes)[index:])
	(*nodes)[index] = &node{word: word}
	return (*nodes)[index]
}

