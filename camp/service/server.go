package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"runtime/debug"
	"github.com/mufe/golang-base/camp/xlog"
)

var s *server

type server struct {
	r *grpc.Server
}

func init() {
	var opts []grpc.ServerOption
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		//必须要先声明defer，否则不能捕获到panic异常
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("%+v", e)
				xlog.Errorf("%+v\n%s", e, string(debug.Stack()))
			}
		}()
		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	r := grpc.NewServer(opts...)
	s = &server{r: r}
}

func GetRegisterRpc() *grpc.Server {
	return s.r
}

//开启服务
func StartService() {
	//  创建server端监听端口
	list, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		fmt.Println(err)
	}
	err = s.r.Serve(list)
	if err != nil {
		fmt.Println(err)
	}
}
