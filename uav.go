package UAV_Express

import (
	"errors"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/MeteorKL/koala"
)

type Pos struct {
	Longitude float64
	Latitude  float64
	Addr      []string
}

func pos(Longitude float64, Latitude float64, Addr []string) *Pos {
	return &Pos{
		Longitude: Longitude,
		Latitude:  Latitude,
		Addr:      Addr,
	}
}

type Uav struct {
	ID        string
	From      *Pos
	FromTime  int64
	To        *Pos
	Longitude float64
	Latitude  float64
	Status    string
	Mutex     *sync.Mutex
}

var uavs = make(map[string]*Uav)
var shitang2Pos = pos(120.131658, 30.273568, []string{"杭州市", "浙大路", "浙江大学玉泉校区二食堂"})
var she30Pos = pos(120.132193, 30.275040, []string{"杭州市", "浙大路", "浙江大学玉泉校区30舍"})
var she31Pos = pos(120.131402, 30.275439, []string{"杭州市", "浙大路", "浙江大学玉泉校区31舍"})

func randomPos(Longitude float64, Latitude float64, to *Pos) (float64, float64) {
	t := time.Now().Unix()
	r := float64(rand.Int63n(t)) / float64(t)
	println(r)
	Longitude = r*(to.Longitude-Longitude) + Longitude
	Latitude = r*(to.Latitude-Latitude) + Latitude
	return Longitude, Latitude
}

var uavIDs = []string{
	"0号",
	"Inszva",
}

func addUavs(id []string, initPos *Pos) {
	for _, ID := range id {
		uavs[ID] = &Uav{
			ID:        ID,
			From:      initPos,
			Longitude: initPos.Longitude,
			Latitude:  initPos.Latitude,
			Status:    "ob",
			Mutex:     new(sync.Mutex),
		}
	}
}

func addUav(ID string, initPos *Pos) {
	uavs[ID] = &Uav{
		ID:        ID,
		From:      initPos,
		Longitude: initPos.Longitude,
		Latitude:  initPos.Latitude,
		Status:    "ob",
		Mutex:     new(sync.Mutex),
	}
}

func removeUav(id string) {
	delete(uavs, id)
}

func distance(from *Pos, to *Pos) float64 {
	return (((to.Longitude - from.Longitude) * math.Pi * 12656 * math.Cos(((to.Latitude+from.Latitude)/2)*math.Pi/180) / 180) *
		((to.Longitude - from.Longitude) * math.Pi * 12656 * math.Cos(((to.Latitude+from.Latitude)/2)*math.Pi/180) / 180)) +
		(((to.Latitude - from.Latitude) * math.Pi * 12656 / 180) *
			((to.Latitude - from.Latitude) * math.Pi * 12656 / 180))
}

func getAvailableUav(From *Pos) (string, error) {
	for ID := range uavs {
		if uavs[ID].From == From && uavs[ID].Status == "ob" {
			return ID, nil
		}
	}
	return "", errors.New("没有可用的无人机")
}

func assignUavTasks(From *Pos, To *Pos) string {
	if distance(From, To) > distance(she31Pos, shitang2Pos)*5 {
		println(distance(From, To))
		println(distance(she31Pos, shitang2Pos) * 5)
		return "附近没有无人机"
	}
	ID, err := getAvailableUav(From)
	if err != nil {
		return err.Error()
	}
	uavs[ID].FromTime = time.Now().Unix()
	uavs[ID].To = To
	uavs[ID].Status = "配送中"
	return "请求无人机配送成功"
}

func (u *Uav) randomMove() {
	u.Longitude, u.Latitude = randomPos(u.Longitude, u.Latitude, u.To)
}

func randomMove() {
	for id := range uavs {
		if uavs[id].To != nil {
			uavs[id].randomMove()
		}
	}
}

const v float64 = 0.030

func (u *Uav) regularMove(dstStatus string) {
	r := v / distance(u.From, u.To)
	u.Mutex.Lock()
	u.Longitude = r*(u.To.Longitude-u.From.Longitude) + u.Longitude
	u.Latitude = r*(u.To.Latitude-u.From.Latitude) + u.Latitude
	if distance(pos(u.Longitude, u.Latitude, []string{}), u.From) > distance(u.From, u.To) {
		u.Longitude = u.To.Longitude
		u.Latitude = u.To.Latitude
		u.Status = dstStatus
	}
	u.Mutex.Unlock()
}

func (u *Uav) reachDstUav() {
	u.Status = "正在配对停机坪"
	time.Sleep(time.Second * 5)
	u.Mutex.Lock()
	u.Status = "正在返程"
	u.From, u.To = u.To, u.From
	u.Mutex.Unlock()
}

func regularMove() {
	for {
		for id := range uavs {
			switch uavs[id].Status {
			case "订单已送达":
				go uavs[id].reachDstUav()
			case "正在返程":
				uavs[id].regularMove("任务已完成")
			case "配送中":
				uavs[id].regularMove("订单已送达")
			case "任务已完成":
				uavs[id].Mutex.Lock()
				uavs[id].Status = "ob"
				uavs[id].From, uavs[id].To = uavs[id].To, uavs[id].From
				uavs[id].Mutex.Unlock()
			}
		}
		time.Sleep(time.Second)
	}
}

func uavHandlers() {
	println(distance(she30Pos, shitang2Pos) * 1000)
	println(distance(she30Pos, she31Pos) * 1000)
	println(distance(she31Pos, shitang2Pos) * 1000)
	addUavs(uavIDs, shitang2Pos)
	go regularMove()
	koala.Get("/api/getuavs", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		koala.WriteJSON(w, map[string]interface{}{
			"status":  0,
			"message": "获取无人机信息成功",
			"data":    uavs,
		})
	})

	koala.Get("/api/requestuav", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		var status = 0
		var message string
		var data interface{}
		var Longitude, Latitude float64
		var addr []string
		if arr, ok := p.ParamGet["Longitude"]; ok {
			a := arr[0]
			Longitude, _ = strconv.ParseFloat(a, 64)
		} else {
			status = 10000
			message = "请输入经度"
		}
		if arr, ok := p.ParamGet["Latitude"]; ok {
			a := arr[0]
			Latitude, _ = strconv.ParseFloat(a, 64)
		} else {
			status = 10001
			message = "请输入纬度"
		}
		if arr, ok := p.ParamGet["detail"]; ok {
			addr = arr
		}
		pos := pos(Longitude, Latitude, addr)
		println(pos)
		message = assignUavTasks(shitang2Pos, pos)
		data = pos
		koala.WriteJSON(w, map[string]interface{}{
			"status":  status,
			"message": message,
			"data":    data,
		})
	})

}
