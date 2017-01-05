package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
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
	tcp, _ := strconv.ParseUint(os.Getenv("WSH_TCP"), 10, 64)

	s := &wsh2s.Server{
		AcmeDomain:         paas.GetEnv("ACME_DOMAIN", paas.Info.WsDomain),
		DropboxAccessToken: os.Getenv("DROPBOX_AK"),
		DropboxDomainKey:   os.Getenv("DROPBOX_DK"),
		H2SleepToRunSecond: time.Duration(h2SleepToRunSecond),
		PingSecond:         uint(pingSecond),
		TCP:                tcp,
		ServerCrt:          permFromEnv("SERVER_CRT"),
		ServerKey:          permFromEnv("SERVER_KEY"),
		ChainPerm:          permFromEnv("CHAIN_PERM"),
	}

	err := s.Serve()
	wsh2s.Log.Fatal("Server failed", zap.Error(err))
}

func permFromEnv(e string) []byte {
	perm := formatPerm(os.Getenv(e))
	return []byte(perm)
}

var permOneLineRe = regexp.MustCompile(`\s?\-+\s?`)

func permOneLineReplace(o string) string {
	return strings.Replace(o, " ", "\n", -1)
}

func formatPerm(s string) string {
	return permOneLineRe.ReplaceAllStringFunc(s, permOneLineReplace)
}
