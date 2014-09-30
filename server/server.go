package main

import ( 
    "github.com/go-martini/martini"
    "os/exec"
    "os"
    "encoding/json"
    "log"
    "net/http"
    "strconv"
)

func main() {

    m := martini.Classic()

    m.Get("/media/:path.js", func(params martini.Params) string {

        cmd := exec.Command("ffprobe", "-loglevel", "error", "-show_format", "media/" + params["path"],"-print_format", "json")
        out, err := cmd.Output()

        if err != nil {
            log.Fatal(err)
        }

        type VideoInfo struct {
            Format struct {
                Duration string `json:"duration"`
            } `json:"format"`
        }

        var vinfo VideoInfo
        json.Unmarshal(out,&vinfo)
        floatDuration, _ := strconv.ParseFloat(vinfo.Format.Duration, 64)
        floatDuration += 1
        strDuration := strconv.FormatUint(uint64(floatDuration), 10)
        vi := map[string]string{"duration":strDuration}
        duration, _ := json.Marshal(vi)

        return(string(duration))

    })
    m.Get("/media/:path.ogv", func(params martini.Params, w http.ResponseWriter, req *http.Request) {

        path := params["path"]

        start := req.URL.Query().Get("start")
        if start == "" {
            start = "0"
        }
        cmd := exec.Command("ffmpeg", "-loglevel", "panic", "-i", "media/"+path, "-ss", start, "-f", "ogg", "-acodec", "libvorbis", "-qscale:v", "10", "pipe:1")

        cmd.Stdout = w
        cmd.Stderr = os.Stderr

        if err := cmd.Run(); err != nil {
            log.Fatal(err)
        }
    })

    m.Run()
}
