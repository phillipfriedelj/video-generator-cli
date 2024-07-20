package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type Video struct {
	name string
	extension string
	fullPath string
}

func generateInputStreamString(videoCount int) string {
	inputStream := ""
	for i:= 0; i < videoCount; i++ {
		inputStream += fmt.Sprintf("[%d:v]", i)
	}

	return inputStream
}


func main() {
	shell := "sh"
	commandFlag := "-c"
	if runtime.GOOS == "windows" {
		shell = "cmd.exe"
		commandFlag = "/C"
	}


	//1. get path to folder
	videoFolderPath := flag.String("path", ".", "The folder where the resources are. If none is provided, the program will look in the current folder")
	flag.Parse()
	fmt.Printf("PATH: %s\n", *videoFolderPath)
	if fi, err := os.Stat(*videoFolderPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			panic("provided path does not exist")
		} else if !fi.Mode().IsDir(){
			panic("provided path is not a folder")
		} else {
			panic("path error: " + err.Error())
		}
	}

	//2. get videos in folder
	files, err := os.ReadDir(*videoFolderPath)
	if err != nil {
		panic("error reading files in folder: " + err.Error())
	}

	var videos []Video
	var soundtracks []string	

	for _, file := range files {
		println(file.Name())	
		if !file.IsDir() {
			extension := strings.ToLower(filepath.Ext(file.Name()))
			switch extension{
			case ".mp4":
				videos = append(videos, Video{name: file.Name(), fullPath: path.Join(*videoFolderPath, file.Name()), extension: extension})
				case ".mp3":
				case ".aac": 
				soundtracks = append(soundtracks, file.Name())
			}
		}
	}	

	for _, video := range videos {
		fmt.Printf("V: %s\n", video.name)
	}

	for _, song := range soundtracks {
		fmt.Printf("S: %s\n", song)
	}
	//3. combine them
	combineCmd := "ffmpeg -y"
	for _, clip := range videos {
		combineCmd += fmt.Sprintf(" -i %s", clip.fullPath)
	}
	combineCmd += fmt.Sprintf(" -filter_complex \"%sconcat=n=%d:v=1:a=0[v]\" -map \"[v]\" %s", generateInputStreamString(len(videos)), len(videos), path.Join(*videoFolderPath, "combined.mp4"))

	fmt.Println("------- ", combineCmd)

	if err := exec.Command(shell, commandFlag, combineCmd).Run(); err != nil {
		fmt.Println("Error combining video clips:", err)
		return
	}

	// Add soundtrack to the main video
	// addSoundtrackCmd := fmt.Sprintf("ffmpeg -y -i %s -i %s -c:v copy -c:a aac -strict experimental temp_with_sound.mp4", videos[0], soundtracks[0])
	// if err := exec.Command(shell, commandFlag, addSoundtrackCmd).Run(); err != nil {
	// 	fmt.Println("Error adding soundtrack:", err.Error())
	// 	return
	// }
}