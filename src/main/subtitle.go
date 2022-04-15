package main

import (
    "avalon-ffmpeg/src/ffmpeg"
    "bufio"
    "errors"
    jsoniter "github.com/json-iterator/go"
    "golang.org/x/text/language"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path"
    "strconv"
    "strings"
    "time"
)

type Subtitle struct {
    Body []struct {
        Content  string  `json:"content"`
        From     float64 `json:"from"`
        Location int     `json:"location"`
        To       float64 `json:"to"`
    } `json:"body"`
}

func GetSubtitleCommand(index int, subtitle *ffmpeg.Subtitle) (string, []string, error) {
    fileName := StringBuilder(subtitle.GetLocale(), RandString(8))
    filePath := path.Join(Conf.GetString(`subtitle.dir`), fileName+`.srt`)
    
    err := SubtitleProcess(subtitle.GetSubtitleUrl(), fileName)
    if err != nil {
        return ``, nil, err
    }
    
    locale := strings.ReplaceAll(subtitle.GetLocaleText(), " ", "_")
    
    args := []string{`-map`, strconv.Itoa(index + 2), StringBuilder(`-metadata:s:s:`, strconv.Itoa(index)), StringBuilder(`handler=`, locale)}
    tag, err := language.Parse(subtitle.GetLocale())
    if err == nil {
        b, _ := tag.Base()
        args = append(args, StringBuilder(`-metadata:s:s:`, strconv.Itoa(index)), StringBuilder(`language=`, b.ISO3()))
    }
    
    return filePath, args, nil
}

func SubtitleProcess(url string, name string) error {
    var sub Subtitle
    
    err := getSubtitle(url, &sub)
    if err != nil {
        return err
    }
    
    err = writeSubtitle(&sub, name)
    if err != nil {
        return err
    }
    
    return nil
}

func getSubtitle(url string, bind interface{}) error {
    req, err := http.NewRequest(`GET`, url, nil)
    if err != nil {
        return err
    }
    
    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    
    defer resp.Body.Close()
    
    respBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    err = jsoniter.Unmarshal(respBytes, bind)
    if err != nil {
        return err
    }
    
    return nil
}

func writeSubtitle(sub *Subtitle, name string) error {
    if sub == nil {
        return errors.New(`empty subtitle`)
    }
    
    dstFile, err := os.OpenFile(path.Join(Conf.GetString(`subtitle.dir`), name+`.srt`), os.O_CREATE|os.O_WRONLY, os.ModePerm)
    if err != nil {
        log.Fatalf("open file failed, err:%v", err)
    }
    bufWriter := bufio.NewWriter(dstFile)
    defer func() {
        bufWriter.Flush()
        dstFile.Close()
    }()
    
    for i := range sub.Body {
        _, err = bufWriter.WriteString(StringBuilder(strconv.Itoa(i+1)+"\n", GetSubtitleTime(sub.Body[i].From, sub.Body[i].To), sub.Body[i].Content, "\n"))
        if err != nil {
            return err
        }
        i += 1
    }
    
    return nil
}

/*func GetFilled2Int(i float64) string {
    return strconv.FormatFloat(i, 'f', 2, 64)[2:]
}*/

func GetSubtitleTime(from float64, stop float64) string {
    return StringBuilder(TimeParsing(from), ` --> `, TimeParsing(stop), "\n")
}

func TimeParsing(stamp float64) string {
    str := strconv.FormatFloat(stamp, 'f', -1, 64)
    
    s := strings.Split(str, `.`)
    fl := `00`
    if len(s) == 2 {
        switch len(s[1]) {
        case 1:
            fl = s[1] + `0`
        case 2:
            fl = s[1]
        }
    }
    
    in, _ := strconv.Atoi(s[0])
    
    t := time.Unix(int64(in), 0).UTC().Format(`15:04:05`)
    t += StringBuilder(`,`, fl)
    return t
}
