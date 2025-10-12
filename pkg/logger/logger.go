package logger

import (
	"github.com/cloudwego/kitex/pkg/klog"
)

// log先随便写一下

func Debugf(template string, args ...interface{}) {
	klog.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	klog.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	klog.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	klog.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	klog.Fatalf(template, args...)
}
