package cloudmusic

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/golang/glog"
	_ "github.com/mattn/go-sqlite3"
)

var (
	isWin bool

	homeDir      string
	dbPath       string
	db           *sql.DB
	downloadPath string
)

type PlayListID uint64 // play list id type.
type PlayListDetail struct {
	PID    PlayListID
	Name   string // 歌单名. Note: ID唯一.Name不唯一.
	Detail map[string]interface{}
}

type SongID uint64 // a song id type.
type SongDetail struct {
	SID          SongID
	Name         string // 歌名
	FileName     string // 文件名(不带路径).
	FileNamePath string // 文件路径.
	Detail       map[string]interface{}
}

func init() {
	isWin = runtime.GOOS == "windows"

	usr, err := user.Current()
	if err != nil {
		glog.Fatalf("err when user.Current(): %v", err)
	}
	homeDir = usr.HomeDir

	downloadPath = fmt.Sprintf("%s/Music/网易云音乐/", homeDir)
	dbPath = homeDir + "/Library/Containers/com.netease.163music/Data/Documents/storage/sqlite_storage.sqlite3"

	// get db
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		glog.Fatalf("err when open %v: %v", dbPath, err)
	}
}

// GetAllPlayLists 获取全部歌单
func GetAllPlayLists() (map[PlayListID]*PlayListDetail, error) {
	rows, err := db.Query(`select * from web_playlist`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(map[PlayListID]*PlayListDetail)
	for rows.Next() {
		var pid PlayListID
		var str string
		err = rows.Scan(&pid, &str)
		if err != nil {
			return nil, err
		}
		detail := PlayListDetail{}
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(str), &data); err != nil {
			glog.Fatalf("err when json.Unmarshal(str=%v): err=%v", str, err)
		}

		detail.PID = pid
		detail.Name = data["name"].(string)
		detail.Detail = data
		res[pid] = &detail
	}
	return res, nil
}

// GetPIDsWithName 根据歌单名获取所有与之对应的PlayListID.(可能对应多个.)
func GetPIDsWithName(name string) ([]PlayListID, error) {
	lists, err := GetAllPlayLists()
	if err != nil {
		return nil, err
	}

	res := make([]PlayListID, 0, len(lists))
	for pid, detail := range lists {
		if detail.Name == name {
			res = append(res, pid)
		}
	}
	return res, nil
}

// GetSIDsWithPID 获取PlayListID对应的SongID
func GetSIDsWithPID(pid PlayListID) ([]SongID, error) {
	rows, err := db.Query(fmt.Sprintf(`select tid from web_playlist_track where pid=%v`, pid))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]SongID, 0)
	for rows.Next() {
		var sid SongID
		err = rows.Scan(&sid)
		if err != nil {
			return nil, err
		}
		res = append(res, sid)
	}
	return res, nil
}

// GetDownloadedSongDetailWithSongID 根据歌曲id获得【下载下来了的歌曲】的详情。 Note: 如果歌曲id没有对应的下载下来了的歌曲的详情记录,则返回nil.
func GetDownloadedSongDetailWithSongID(sid SongID) (*SongDetail, error) {
	rows, err := db.Query(fmt.Sprintf(`select track_id,detail,relative_path,track_name from web_offline_track where track_id=%v`, sid))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sid SongID   // track_id
		var str string   // detail
		var rpath string // relative_path
		var name string  // track_name
		err = rows.Scan(&sid, &str, &rpath, &name)
		if err != nil {
			return nil, err
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(str), &data); err != nil {
			glog.Fatalf("err when json.Unmarshal(str=%v): err=%v", str, err)
		}

		filename := filepath.Base(rpath)
		return &SongDetail{
			SID:          sid,
			Name:         name,
			FileNamePath: filepath.Join(downloadPath, filename),
			FileName:     filename,
			Detail:       data,
		}, nil
	}
	return nil, nil
}
