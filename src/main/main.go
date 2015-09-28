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
)

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
    sign := dict.NewSign("../data/dictionary/sign.txt")
    stop := dict.NewSign("../data/dictionary/stopwords.txt")
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")

    s := &WordSetting {
        freqDoor: 6,
        scoreDoor: 0.01,
        signDict: sign,
        stopDict: stop,
        wordDict: d,
    }

    return s
}

func main() {
    flag.Parse()
    root := flag.Arg(0)
    //getFileList(root)
    handlePath(root)
}

func mainMerge() {
    var freqDoor int
    var scoreDoor float32
    freqDoor = 6 
    scoreDoor = 0.010

    sign := dict.NewSign("../data/dictionary/sign.txt")
    stop := dict.NewSign("../data/dictionary/stopwords.txt")
    d := dict.NewDictionary("../data/dictionary/sogoudictionary.txt")

    article := util.ReadFile("../data/testdata/125.txt")
    allsegs := segment.SegmentDoc(article, sign, d)
    
    fmt.Println(len(allsegs))
    str := segment.GetSegmentStr(allsegs)

    util.WriteFile("../data/test-125.log", str)
    allsegs = occurrence.FilterSegment(allsegs, stop)

    //occur := occurrence.NewOccurrence()
    //occur.AddSegments(allsegs, stop)
    //occur.Compute()
    //occur.Output()

    pairTerms := occurrence.Merge(allsegs, freqDoor, scoreDoor)
    str = term.GetPairTermStr(pairTerms)
    util.WriteFile("../data/main-test-125-merge-merge.log", str)
}

func handlePath(root string) {
    ws := NewWordSetting()
    files := getFilePath(root)
    for _, f := range files {
        fullfilepath := filepath.Join(f.folder, f.filename)
        fmt.Println("Handle the file: ", fullfilepath)
        content := util.ReadFile(fullfilepath)
        pairTerm := getWords(content, ws)
        writeOutput(f, pairTerm)
    }
}

func getWords(content string, ws *WordSetting) []*term.PairTerm {
    allsegs := segment.SegmentDoc(content, ws.signDict, ws.wordDict)
    
    //fmt.Println(len(allsegs))
    //str := segment.GetSegmentStr(allsegs)

    //util.WriteFile("../data/test-125.log", str)
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
    
    format := "%s,%d,%f\n"
    str := ""
    for _, pt := range pairTerms {
        str += fmt.Sprintf(format, pt.GetKey(), pt.GetFrequency(), pt.GetScore())
    }

    util.WriteFile(outfile, str)
    fmt.Println("Store the word in: ", outfile)
}
