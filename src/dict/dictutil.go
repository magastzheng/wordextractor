package dict

import (
    "fmt"
    "log"
    "os"
    "bufio"
    "strings"
    "strconv"
    "unicode"
    "unicode/utf8"
)

const (
    minTokenFrequency = 2
)

func LoadDictionary(filenames string) *Dictionary {
    dict := new(Dictionary)
    for _, filename := range strings.Split(filenames, ",") {
        log.Printf("Load dictionary: %s", filename)
        file, err := os.Open(filename)
        defer file.Close()
        if err != nil {
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

            if frequency < minTokenFrequency {
                continue
            }

            words := splitTextToWords([]byte(text))
            token := Token{text: words, frequency: frequency, nature: nature}
            dict.addToken(&token)
        }
    }

    log.Println("Finish to load dict")

    return dict
}

//func LookupTokens(

// 将文本划分成字元
func splitTextToWords(text Text) []Text {
	output := make([]Text, 0, len(text)/8)
	current := 0
	inAlphanumeric := true
	alphanumericStart := 0
	for current < len(text) {
		r, size := utf8.DecodeRune(text[current:])
		if size <= 2 && (unicode.IsLetter(r) || unicode.IsNumber(r)) {
			// 当前是拉丁字母或数字（非中日韩文字）
			if !inAlphanumeric {
				alphanumericStart = current
				inAlphanumeric = true
			}
		} else {
			if inAlphanumeric {
				inAlphanumeric = false
				if current != 0 {
					output = append(output, toLower(text[alphanumericStart:current]))
				}
			}
			output = append(output, text[current:current+size])
		}
		current += size
	}

	// 处理最后一个字元是英文的情况
	if inAlphanumeric {
		if current != 0 {
			output = append(output, toLower(text[alphanumericStart:current]))
		}
	}

	return output
}

// 将英文词转化为小写
func toLower(text []byte) []byte {
	output := make([]byte, len(text))
	for i, t := range text {
		if t >= 'A' && t <= 'Z' {
			output[i] = t - 'A' + 'a'
		} else {
			output[i] = t
		}
	}
	return output
}
