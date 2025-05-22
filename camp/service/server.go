package service

import (
	"context"
	"fmt"
	"github.com/mufe/golang-base/camp/xlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"os"
	"runtime/debug"
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
		clientIP := GetClientIPFromMetadata(ctx)
		if clientIP != "" {
			ctx = context.WithValue(ctx, "client-ip", clientIP)
		}
		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	r := grpc.NewServer(opts...)
	s = &server{r: r}
}
func GetClientIPFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if ips := md.Get("client-ip"); len(ips) > 0 {
		return ips[0]
	}
	// 可选 fallback
	if ips := md.Get("x-forwarded-for"); len(ips) > 0 {
		return ips[0]
	}
	return ""
}

func GetClientIPFromCtx(ctx context.Context) string {
	ip, _ := ctx.Value("client-ip").(string)
	return ip
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
