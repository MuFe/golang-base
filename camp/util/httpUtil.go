package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/mufe/golang-base/camp/xlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"net"
	"strconv"
)

func IsMiniProgram(c *gin.Context) bool {
	if c.GetHeader("ProgramTag") == "miniProgram" {
		return true
	} else {
		return false
	}
}

func GetHeaderFromKey(c *gin.Context, key string) string {
	return c.GetHeader(key)
}

func GetInt32ValueFromReq(c *gin.Context, key string) int32 {
	return FormatStrToInt32(c.Query(key))
}

func GetInt64ValueFromReq(c *gin.Context, key string) int64 {
	return FormatStrToInt64(c.Query(key))
}

func GetFloat64ValueFromReq(c *gin.Context, key string) float64 {
	value, err := strconv.ParseFloat(c.Query(key), 64)
	if err != nil {
		value = 0
	}
	return value
}

func FormatStrToInt32(str string) int32 {
	intStr, err := strconv.Atoi(str)
	if err != nil {
		intStr = 0
	}
	return int32(intStr)
}

func FormatStrToInt64(str string) int64 {
	intStr, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		intStr = 0
	}
	return intStr
}

type BaseResult struct {
	Total   int64         `json:"total"`
	Current int64         `json:"current"`
	List    []interface{} `json:"list"`
}

func CreateListResultReturn(total int64, list []interface{}) BaseResult {
	return BaseResult{
		Total: total,
		List:  list,
	}
}
func CreateListCurrentResultReturn(total, current int64, list []interface{}) BaseResult {
	return BaseResult{
		Total:   total,
		Current: current,
		List:    list,
	}
}

// 获取rpc服务(服务发现)
func GetRPCServiceBase(consulIp string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(consulIp, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, err
}

// 获取rpc服务(服务发现)
func GetRPCService(name string, tag string, consulIp string) (*grpc.ClientConn, error) {
	url := fmt.Sprintf("%s://%s/%s/%s", "consul", consulIp, name, tag)
	conn, err := grpc.Dial(url, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithBlock())
	if err != nil {
		xlog.ErrorP(err)
		return nil, err
	}
	return conn, err
}

func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
