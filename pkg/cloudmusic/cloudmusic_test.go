package cloudmusic

import (
	"log"
	"testing"
)

func Test1(t *testing.T) {
	playListName := "我喜欢的音乐"
	pids, err := GetPIDsWithName(playListName)
	if err != nil {
		log.Fatal(err)
	}

	songIDs := make([]SongID, 0)
	for _, pid := range pids {
		sids, err := GetSIDsWithPID(pid)
		if err != nil {
			log.Fatal(err)
		}
		songIDs = append(songIDs, sids...)
	}

	log.Printf("see SongID[] of [%s] = %v\n", playListName, songIDs)
	log.Printf("see len of SongID[] of [%s] = %d", playListName, len(songIDs))

	// 歌曲详情：
	songs := make([]*SongDetail, 0, len(songIDs))
	for _, sid := range songIDs {
		det, err := GetDownloadedSongDetailWithSongID(sid)
		if err != nil {
			log.Fatal(err)
		}
		if det != nil {
			songs = append(songs, det)
		}
	}

	for _, song := range songs {
		log.Printf("see detail of song which have been downloaded: name=%v, FileName=%v", song.Name, song.FileName)
	}
	log.Printf("see total number of songs in [%s] those have been downloaded: %d", playListName, len(songs))
}
