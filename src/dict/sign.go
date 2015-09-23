package dict

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "io"
    "log"
)

type Sign struct {
    signs []string
}

func (s *Sign) IsContain(ch string) (ret bool) {
    for _, v := range s.signs {

        if v == ch {
            ret = true
            return
        }
    }

    return
}

func NewSign(filename string) *Sign {
    s := &Sign{
        signs: make([]string, 0),
    }

    log.Printf("Load sign dictionary: %s", filename)
    file, err := os.Open(filename)
    defer file.Close()
    if err != nil {
        log.Fatalf("Fail to load sign dictionary file: %s", filename)
    }

    reader := bufio.NewReader(file)
    
    //read each line
    for {
        line, err := reader.ReadString('\n')
        if err != nil || io.EOF == err {
            break
        }
        
        line = strings.TrimSpace(line)
        if len(line) == 0 {
            continue
        }
        
        s.signs = append(s.signs, line)
    }
    
    //fmt.Println(s.signs)
    str := fmt.Sprintf("Finish to load sign dictionary: %d .", len(s.signs))
    log.Println(str)

    return s
}
