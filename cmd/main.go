package main

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/spf13/viper"
	"github.com/xiaou/cpCloudMusic/pkg/cloudmusic"
	"github.com/xiaou/cpCloudMusic/pkg/util"
)

var (
	MaxCopyRoutineNumber = 100
	homeDir              string
)

func init() {
	usr, err := user.Current()
	if err != nil {
		glog.Fatalf("err when user.Current(): %v", err)
	}
	homeDir = usr.HomeDir
}

func copySongs(songs []*cloudmusic.SongDetail, toPath string) {
	if len(songs) == 0 {
		return
	}

	var waitgroup sync.WaitGroup
	chBuf := make(chan int, MaxCopyRoutineNumber)

	for _, song := range songs {
		src := song.FileNamePath
		dst := filepath.Join(toPath, song.FileName)

		waitgroup.Add(1)
		chBuf <- 1

		go func(src, dst string, c <-chan int) {
			defer waitgroup.Done()
			defer func() { <-c }()

			err := util.CopyFile(src, dst)
			if err != nil {
				glog.Errorf("err when cp to %v: %v", dst, err)
			} else {
				fmt.Printf(".")
			}
		}(src, dst, chBuf)
	}
	//fmt.Println("")

	waitgroup.Wait()
}

func main() {
	BeginMain()
	defer EndMain()

	playListName := viper.GetString("name")
	outPath := viper.GetString("out")
	if strings.HasPrefix(outPath, "~") {
		outPath = strings.Replace(outPath, "~", homeDir, 1)
	}

	pids, err := cloudmusic.GetPIDsWithName(playListName)
	if err != nil {
		glog.Fatal(err)
	}

	songIDs := make([]cloudmusic.SongID, 0)
	for _, pid := range pids {
		sids, err := cloudmusic.GetSIDsWithPID(pid)
		if err != nil {
			glog.Error(err)
			continue
		}
		songIDs = append(songIDs, sids...)
	}

	glog.Infof("see SongID[] of [%s] = %v\n", playListName, songIDs)
	glog.Infof("see len of SongID[] of [%s] = %d", playListName, len(songIDs))

	// 歌曲详情：
	downloadedSongs := make([]*cloudmusic.SongDetail, 0, len(songIDs))
	for _, sid := range songIDs {
		song, err := cloudmusic.GetDownloadedSongDetailWithSongID(sid)
		if err != nil {
			glog.Error(err)
			continue
		}
		if song != nil {
			downloadedSongs = append(downloadedSongs, song)
		}
	}

	glog.Infof("see total number of songs in [%s] those have been downloaded: %d", playListName, len(downloadedSongs))

	//
	if len(downloadedSongs) > 0 {
		glog.Infof("see detail:")
		for _, song := range downloadedSongs {
			glog.Infof("Name=%v, FileName=%v", song.Name, song.FileName)
		}

		// 拷贝:
		glog.Info("")
		glog.Infof("\tnow cp all above....")
		copySongs(downloadedSongs, outPath)
		glog.Info("")
		glog.Infof("\tDone.\there: %s", outPath)
	}
}
