package segment

import (
    "fmt"
    "dict"
    "strings"
    "unicode"
    "util"
)

type Document struct {
    filename    string
    buf         []rune
}

type Sentence struct {
    text    []rune
    start   int
    end     int
}

func (s *Sentence) Buffer() []rune {
    return s.text
}

func (s *Sentence) Text() string {
    return string(s.text)
}

func (s *Sentence) Start() int {
    return s.start
}

func (s *Sentence) End() int {
    return s.end
}

func (s *Sentence) ToString() string {
    format := "Start: %d, end: %d, text: %s"
    return fmt.Sprintf(format, s.start, s.end, s.Text())
}

func (d *Document) load(filename string) {
    article := util.ReadFile(filename)
    article = DeleteSpaceChar(article)

    d.filename = filename
    d.buf = []rune(article)
}

func (d *Document) Buffer() []rune {
    return d.buf
}

func (d *Document) Filename() string {
    return d.filename
}

func NewDocument(filename string) *Document {
    d := &Document{}
    d.load(filename)

    return d
}

func IsSpace(ch rune) bool {
    //return ch == ' '
    return unicode.IsSpace(ch)
}

//Delete \n, \r, \t
func DeleteSpaceChar(article string) string {
    article = strings.Replace(article, "\n", "", -1)
    article = strings.Replace(article, "\r", "", -1)
    article = strings.Replace(article, "\t", "", -1)
    return article
}

//Split the arctile into some sentences
func SplitSentence(buf []rune, d *dict.Sign) []*Sentence {
    sentences := make([]*Sentence, 0)
    start := 0
    for i, count := 0, len(buf); i < count; i++ {
        current := string(buf[i]) 
        if IsSpace(buf[i]) {
            fmt.Println(i, start)
            if i == start {
                start++
            }
        }

        if i == count - 1 || d.IsContain(current) {
            if i > start {
                sentence := &Sentence{
                    text: buf[start: i],
                    start: start,
                    end: i,
                }

                sentences = append(sentences, sentence)
            }

            start = i + 1
        }
    }

    return sentences
}
