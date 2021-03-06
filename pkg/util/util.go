package util

import (
	"database/sql"
	"fmt"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

// For testing, bypass HandleCrash.
var ReallyCrash bool

func CallerPathInline(inPkg string) string {
	callers := ""
	for i := 0; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if i == 0 || i == 1 || i == 2 {
			// skip current stack
			continue
		}
		if file == "<autogenerated>" {
			continue
		}
		// only trace stack in current module
		if strings.Contains(file, inPkg) {
			base := path.Base(file)
			callers = callers + fmt.Sprintf("[%v:%v", base, line)
		}
	}
	return callers
}

// HandleCrash simply catches a crash and logs an error. Meant to be called via defer.
func HandleCrash() {
	if ReallyCrash {
		return
	}

	r := recover()
	if r != nil {
		callers := ""
		for i := 0; true; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			callers = callers + fmt.Sprintf("%v:%v\n", file, line)
		}
		glog.Infof("Recovered from panic: %#v (%v)\n%v", r, r, callers)
	}
}

// Forever loops forever running f every period.  Catches any panics, and keeps going.
func Forever(f func(), period time.Duration) {
	Until(f, period, nil)
}

// Until loops until stop channel is closed, running f every period.
// Catches any panics, and keeps going. f may not be invoked if
// stop channel is already closed.
func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		default:
		}
		func() {
			defer HandleCrash()
			f()
		}()
		time.Sleep(period)
	}
}

// WaitForever wait forever.
func WaitForever() {
	select {}
}

// OpenMysql open mysql from url like this:mysql://root:123@localhost:3306/test
func OpenMysql(mysqlURL string) (*sql.DB, error) {
	s := regexp.MustCompile(`(^mysql://)(.+[@])(.+)(/.+)`).ReplaceAllString(mysqlURL, `${2}tcp(${3})${4}?charset=utf8`) // 转换uri为mysql库识别的
	glog.Infof("OpenMysql: %s->%s", mysqlURL, s)
	return sql.Open("mysql", s)
}

// DeepCopy can copy map[string]interface{} or []interface{} in deeply copy.
func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	}

	return value
}

// ToDouble 把interface{}转换成float64
func ToDouble(i interface{}) float64 {
	f, ok := i.(float64)
	if ok {
		return f
	}

	s, ok := i.(string)
	if !ok {
		return 0.0
	}
	res, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return res
}
