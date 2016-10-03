package main

import (
	"os"
	"strconv"
	"time"

	"golang.org/x/net/http2"

	"github.com/empirefox/gotool/paas"
	"github.com/empirefox/wsh2s"
	"github.com/uber-go/zap"
)

func main() {
	http2.VerboseLogs, _ = strconv.ParseBool(os.Getenv("HTTP2_LOG"))
	pingSecond, _ := strconv.ParseUint(os.Getenv("PING_SECOND"), 10, 64)
	wsh2s.Log = zap.New(zap.NewJSONEncoder(zap.NoTime()), zap.AddCaller())
	h2SleepToRunSecond, _ := strconv.ParseUint(os.Getenv("H2_SLEEP_SECOND"), 10, 64)

	s := &wsh2s.Server{
		AcmeDomain:         paas.GetEnv("ACME_DOMAIN", paas.Info.WsDomain),
		DropboxAccessToken: os.Getenv("DROPBOX_AK"),
		DropboxDomainKey:   os.Getenv("DROPBOX_DK"),
		H2SleepToRunSecond: time.Duration(h2SleepToRunSecond),
		PingSecond:         uint(pingSecond),
	}

	err := s.Serve()
	wsh2s.Log.Fatal("Server failed", zap.Error(err))
}
