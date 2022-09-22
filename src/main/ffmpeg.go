package main

import (
    "avalon-ffmpeg/src/ffmpeg"
    "bytes"
    "context"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "os"
    "os/exec"
)

type FFMPEGServiceWorker struct {
    ffmpeg.UnimplementedFFMPEGServer
}

func (F FFMPEGServiceWorker) MergeVideo(ctx context.Context, param *ffmpeg.Param) (*ffmpeg.Info, error) {
    raw := []string{`-loglevel`, `error`, `-i`, param.GetInputVideo(), `-i`, param.GetInputAudio()}
    srtCache := []string{}
    var args []string
    for i := range param.GetSubtitles() {
        if param.GetSubtitles()[i] != nil {
            input, a, err := GetSubtitleCommand(i, param.GetSubtitles()[i])
            if err != nil {
                return nil, status.Error(codes.Internal, err.Error())
            }
            raw = append(raw, `-i`, input)
            srtCache = append(srtCache, input)
            args = append(args, a...)
        }
    }
    
    raw = append(raw, args...)
    raw = append(raw, `-map`, `0:v`, `-map`, `1:a`, `-c:v`, `copy`, `-c:a`, `copy`, `-c:s`, `srt`, `-y`, param.GetOutputVideo())
    
    cmd := exec.Command(`ffmpeg`, raw...)
    
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    
    err := cmd.Run()
    errStr := string(stderr.Bytes())
    
    if err != nil || errStr != `` {
        if errStr != `` {
            return nil, status.Error(codes.Internal, errStr)
        }
        return nil, status.Error(codes.Internal, err.Error())
    }
    
    for i := range srtCache {
        os.Remove(srtCache[i])
    }
    
    return &ffmpeg.Info{
        Code: 0,
        Msg:  `OK`,
    }, nil
}
