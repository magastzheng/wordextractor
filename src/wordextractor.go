package main

import(
    "fmt"
    "dict"
    "segment"
    "util"
    "occurrence"
    "term"
    "path/filepath"
    "os"
    "strings"
    //"os/exec"
    "flag"
    "log"
	"encoding/json"
	"io/ioutil"
    mahonia "github.com/axgle/mahonia"
)

const (
    SignDictionary = "data/dictionary/sign.txt"
    StopDictionary = "data/dictionary/stopwords.txt"
    NormalDictionary = "data/dictionary/sogoudictionary.txt"
	ConfigFile = "config/config.json"
    Encoding = "gb18030"
)

type Config struct {
	FreqDoor 	int 	`json:"freqDoor"`
	ScoreDoor 	float32 `json:"scoreDoor"`
}

type FilePath struct {
    folder string
    filename string
}

type WordSetting struct {
    freqDoor int
    scoreDoor float32
    signDict *dict.Sign
    stopDict *dict.Sign
    wordDict *dict.Dictionary
}

func NewWordSetting() *WordSetting {
    sign := dict.NewSign(SignDictionary)
    stop := dict.NewSign(StopDictionary)
    d := dict.NewDictionary(NormalDictionary)
	config := LoadConfig()
	
    s := &WordSetting {
        freqDoor: config.FreqDoor,
        scoreDoor: config.ScoreDoor,
        signDict: sign,
        stopDict: stop,
        wordDict: d,
    }
	
    return s
}

func LoadConfig() *Config {
	c := &Config{
		FreqDoor: 3,
		ScoreDoor: 0.01,
	}
	chunks, err := ioutil.ReadFile(ConfigFile)
    if err != nil {
		log.Fatalf("Fail to load config file: %s", ConfigFile)
		return c
	}
	
    err = json.Unmarshal(chunks, c)
	if err != nil {
		log.Fatal("Fail to unmarshel config file", err)
		return c
	}
	
	return c
}

func main() {
    flag.Parse()
    root := flag.Arg(0)
    handlePath(root)
}

func handlePath(root string) {
    ws := NewWordSetting()
    decoder := mahonia.NewDecoder(Encoding)
    files := getFilePath(root)
    for _, f := range files {
        fullfilepath := filepath.Join(f.folder, f.filename)
        log.Printf("正在处理文件: %s", fullfilepath)
        content := util.ReadFile(fullfilepath)
        //if ret, ok := decoder.ConvertStringOK(content); ok {
        //    content = ret
        //}
        content = decoder.ConvertString(content)
        pairTerm := getWords(content, ws)
        writeOutput(f, pairTerm)
    }
}

func getWords(content string, ws *WordSetting) []*term.PairTerm {
    allsegs := segment.SegmentDoc(content, ws.signDict, ws.wordDict)
    
    //fmt.Println(len(allsegs))
    //str := segment.GetSegmentStr(allsegs)

    //util.WriteFile("../data/test-segment-125.log", str)
    allsegs = occurrence.FilterSegment(allsegs, ws.stopDict)
    pairTerms := occurrence.Merge(allsegs, ws.freqDoor, ws.scoreDoor)
    //str = term.GetPairTermStr(pairTerms)
    return pairTerms
}

func getFilePath(root string) []*FilePath {
    files := make([]*FilePath, 0)
    
    err := filepath.Walk(root, func(root string, f os.FileInfo, err error) error {
        if f == nil {
            return err
        }
        if f.IsDir() {
            return err
        }

        dir, filename := filepath.Split(root)
        //ext := filepath.Ext(filename)
        //fmt.Println("Ext ****", ext)
        if !strings.Contains(filename, ".csv") {
            fp := &FilePath{
                folder: dir,
                filename: filename,
            }
            files = append(files, fp)
        }
        return err
    })
    
    if err != nil {
        fmt.Println(err)
    }

    return files
}

func writeOutput(file *FilePath, pairTerms []*term.PairTerm) {
    ext := filepath.Ext(file.filename)
    pos := strings.Index(file.filename, ext)
    base := string(file.filename[:pos])
    outfile := base + ".csv"
    //fmt.Println("Outfile name:", outfile)
    outfile = filepath.Join(file.folder, outfile)
    
    //fmt.Println("Total words: ", len(pairTerms))
    format := "%s,%d,%f\n"
    str := "短语,频率,分数\n"
    for _, pt := range pairTerms {
        //if pt.GetFrequency() > 1{
            str += fmt.Sprintf(format, pt.GetKey(), pt.GetFrequency(), pt.GetScore())
        //}
    }
	
	//fmt.Println(str)
    encoder := mahonia.NewEncoder(Encoding)
    str = encoder.ConvertString(str)
    util.WriteFile(outfile, str)
    log.Printf("短语存储到文件: %s", outfile)
}
