package logger

import (
	"os"
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"
)

// log先随便写一下

func InitKlog() {
    // 1. 读取 env
    level := os.Getenv("LOG_LEVEL")
    // 2. 若无 env，则解析命令行参数 --log-level=xxx
    if level == "" {
        for _, arg := range os.Args {
            if strings.HasPrefix(arg, "--log-level=") {
                parts := strings.SplitN(arg, "=", 2)
                if len(parts) == 2 {
                    level = parts[1]
                }
                break
            }
        }
    }

    switch strings.ToLower(strings.TrimSpace(level)) {
    case "debug":
        klog.SetLevel(klog.LevelDebug)
    case "info":
        klog.SetLevel(klog.LevelInfo)
    case "warn":
        klog.SetLevel(klog.LevelWarn)
    case "error":
        klog.SetLevel(klog.LevelError)
    case "fatal":
        klog.SetLevel(klog.LevelFatal)
    default:
        klog.SetLevel(klog.LevelInfo)
    }
}

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
