package utils

import (
	"github.com/mufe/golang-base/camp/xlog"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Get3Code 3位随机数字
func Get3Code() string {
	var letters = []rune("1234567890")
	b := make([]rune, 3)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Get4Code 4位随机数字
func Get4Code() string {
	var letters = []rune("1234567890")
	b := make([]rune, 4)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Get6Code 6位随机数字
func Get6Code() string {
	var letters = []rune("1234567890")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 获取随机字符串
//    length：字符串长度
func GetRandomString(length int) string {
	str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	var (
		result []byte
		b      []byte
		r      *rand.Rand
	)
	b = []byte(str)
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, b[r.Intn(len(b))])
	}
	return string(result)
}

func GetRandomInt(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func WeightedRandomIndex(weight []int64)int64{
	xlog.Info(weight)
	if len(weight)==1{
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	sum:=int64(0)
	for _,v:=range weight{
		sum+=v
	}
	ran := rand.Int63n(sum)
	t:=int64(0)
	xlog.Info(ran)
	for k,v:=range weight{
		t+=v
		if t>ran{
			xlog.Info(k)
			return int64(k)
		}
	}
	return int64(len(weight)-1)
}
