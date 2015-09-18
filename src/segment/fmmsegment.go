package segment

import (
    "fmt"
    "strings"
    "os"
    "bufio"
    //"util"
)

//Forward maximum matching
//http://m.blog.csdn.net/blog/wangliang_f/17527915
type FMMSegment struct {
    //segments    []Segment
    //segments        []string
    dict        map[string]string
    maxWordLen  int
}

func (s *FMMSegment) getSegment(text string) (currernt, rest string) {
    if len(text) > s.maxWordLen {
        currernt = string(text[:s.maxWordLen])
        rest = string(text[s.maxWordLen:])
    }else{
        currernt = text
        rest = ""
    }

    return
}

func (s *FMMSegment) selectMaxLenWord(text string) (word, rest string) {
    for count := len(text); count >= 0; count-- {
        token := text[:count]

        if _, ok := s.dict[token]; ok {
            word = token
            rest = string(text[count:])
            return
        }
    }

    word = text[0:1]
    rest = text[1:]
    
    return
}

func (s *FMMSegment) Segment(text string) []string {
    text = strings.Trim(text, " ")

    //start := 0
    rest := text
    var current string
    wordList := make([]string, 0)
    for {
        current, rest = s.getSegment(rest)

        word, wordright := s.selectMaxLenWord(current)
        wordList = append(wordList, word)

        rest = wordright + rest
        if len(rest) == 0 {
            break
        }
    }

    return wordList
}

func (s *FMMSegment) InitDictionary(filename string) {
    s.dict = make(map[string]string)

    fmt.Printf("Load file: %s\n", filename)
    file, err := os.Open(filename)
    defer file.Close()
	if err != nil {
		fmt.Printf("Fail to load: \"%s\" \n", filename)
	}

	reader := bufio.NewReader(file)
    var text string
    var maxLen int
    for {
        size, _ := fmt.Fscanln(reader, &text)
        if size  == 0 {
            break
        } else {
            if _, ok := s.dict[text]; !ok {
                s.dict[text] = text
                if maxLen < len(text) {
                    maxLen = len(text)
                }
            }
        }

    }

    s.maxWordLen = maxLen
}

func NewFMMSegment(filename string) *FMMSegment {
    s := &FMMSegment{}
    s.InitDictionary(filename)

    return s
}
