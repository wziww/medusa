package stream

import (
	"encoding/json"
	"github/wziww/medusa/config"
	"github/wziww/medusa/log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

// APIServerInit api server start
func APIServerInit() {
	go func() {
		var enable bool
		var port int
		switch config.C.Base.Client {
		case true:
			enable = config.C.Client.API.Enable
			port = config.C.Client.API.Port
		case false:
			enable = config.C.Server.API.Enable
			port = config.C.Server.API.Port
		}
		log.FMTLog(log.LOGINFO, "api server enable:", enable)
		if !enable {
			return
		}
		server := &http.Server{
			Addr: func() string {
				addr := "0.0.0.0:"
				if port == 0 {
					addr += "8083"
				} else {
					addr += strconv.Itoa(port)
				}
				return addr
			}(),
			WriteTimeout: 60 * time.Second,
			ReadTimeout:  60 * time.Second,
		}
		http.HandleFunc("/v1/date", func(w http.ResponseWriter, r *http.Request) {
			v1Data(w, r)
			return
		})
		log.FMTLog(log.LOGINFO, "api service start listen at", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			os.Exit(1)
		}
	}()
}

func v1Data(w http.ResponseWriter, r *http.Request) {
	type res struct {
		FlowIn  uint64
		FlowOut uint64
		Counter interface{} `json:"counter"`
	}
	d := &res{}
	d.FlowIn = atomic.LoadUint64(FlowIn)
	d.FlowOut = atomic.LoadUint64(FlowOut)
	d.Counter = Counter.GetAll()
	w.WriteHeader(200)
	j, _ := json.Marshal(d)
	w.Write(j)
}
