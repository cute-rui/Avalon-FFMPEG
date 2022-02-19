package main

import (
    "avalon-ffmpeg/src/ffmpeg"
    "bytes"
    "context"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "os/exec"
)

type FFMPEGServiceWorker struct {
    ffmpeg.UnimplementedFFMPEGServer
}

func (F FFMPEGServiceWorker) MergeVideo(ctx context.Context, param *ffmpeg.Param) (*ffmpeg.Info, error) {
    cmd := exec.Command(`ffmpeg`, `-loglevel`, `error`, `-i`, param.InputVideo, `-i`, param.InputAudio, `-c`, `copy`, `-y`, param.OutputVideo)
    
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
    
    return &ffmpeg.Info{
        Code: 0,
        Msg:  `OK`,
    }, nil
}
