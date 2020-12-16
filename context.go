package goo

import (
	"context"
	"os"
)

var (
	sig         = make(chan os.Signal)
	ctx, cancel = context.WithCancel(context.Background())
	Context     = ctx
)

// func init() {
// 	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
//
// 	AsyncFunc(func() {
// 		for si := range sig {
// 			switch si {
// 			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
// 				cancel()
// 			case syscall.SIGUSR1:
// 			case syscall.SIGUSR2:
// 			default:
// 			}
// 		}
// 	})
// }
//
// func SetBasePath(basePath string) {
// 	ctx = context.WithValue(ctx, "basePath", basePath)
// }
//
// func BasePath() string {
// 	v := ctx.Value("basePath")
// 	if v == nil {
// 		return ""
// 	}
// 	return v.(string)
// }
