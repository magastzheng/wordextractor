package dict

import (
    "fmt"
    "os"
    "bufio"
    "strconv"
    "log"
)

const (
    MaxWordLength = 20
)

type Dictionary struct {
    line        string
    word        string
    word_map    map[string]int
    size        int
    freq_all    int
    arr_1       [MaxWordLength]int
    arr_2       [MaxWordLength]float64
}

func (d *Dictionary) FindWord(word string) int {
    if v, ok := d.word_map[word]; ok {
        return v
    }

    return -1
}

func (d *Dictionary) FreqAll() int {
    return d.freq_all
}

func NewDictionary(filename string) *Dictionary {
    d := new(Dictionary)
	log.Printf("开始加载分词字典: %s", filename)
    d.word_map = make(map[string]int)
    d.freq_all = 0
    d.size = 0
    for i := 0; i < MaxWordLength; i++ {
        d.arr_1[i] = 0
        d.arr_2[i] = 0.0
    }

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
        if err != nil || len(text) == 0 {
            continue
        }
        
        //fmt.Println(text)
        word_len := len([]rune(text))
        d.arr_1[word_len] = d.arr_1[word_len] + frequency
        if _, ok := d.word_map[text]; !ok {
            d.word_map[text] = frequency
        }
        
        d.freq_all += frequency
    }

    d.size = len(d.word_map)
    for i := 0; i < MaxWordLength; i++ {
        d.arr_2[i] = float64(d.arr_1[i]) / float64(d.freq_all)
    }
	
	str := fmt.Sprintf("成功加载分词字典，总数: %d", d.size)
    log.Println(str)
	
    return d
}

