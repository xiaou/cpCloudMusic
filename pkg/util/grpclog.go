package util

import (
	"fmt"

	"github.com/golang/glog"
	"google.golang.org/grpc/grpclog"
)

func SetGRPCLogger(level glog.Level) {
	grpclog.SetLogger(&glogger{level})
}

type glogger struct {
	l glog.Level
}

func (g *glogger) Fatal(args ...interface{}) {
	glog.FatalDepth(2, args...)
}

func (g *glogger) Fatalf(format string, args ...interface{}) {
	glog.FatalDepth(2, fmt.Sprintf(format, args...))
}

func (g *glogger) Fatalln(args ...interface{}) {
	glog.FatalDepth(2, fmt.Sprintln(args...))
}

func (g *glogger) Print(args ...interface{}) {
	if glog.V(g.l) {
		glog.InfoDepth(2, args...)
	}
}

func (g *glogger) Printf(format string, args ...interface{}) {
	if glog.V(g.l) {
		glog.InfoDepth(2, fmt.Sprintf(format, args...))
	}
}

func (g *glogger) Println(args ...interface{}) {
	if glog.V(g.l) {
		glog.InfoDepth(2, fmt.Sprintln(args...))
	}
}
