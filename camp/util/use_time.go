package utils

import (
	"fmt"
	"time"
)

type useTime struct {
	Name string
	Time time.Time
}

// UseTime 用于查看操作时间
func UseTime(name string) useTime {
	var u = useTime{
		Name: name,
		Time: time.Now(),
	}
	u.print("start", 0)
	return u
}

func (t *useTime) Mark(s string) {
	var now = time.Now()
	t.print(s, float64(now.Sub(t.Time).Nanoseconds()/1e4)/100.0)
	t.Time = now
}

func (t *useTime) print(s string, u float64) {
	fmt.Println(fmt.Sprintf("----------\t[USE_TIME]\tNAME：%s\tMARK：%s\tUSE：%.2fms\t----------", t.Name, s, u))
}
