package main

import (
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/empirefox/esecend/admin"
	"github.com/empirefox/esecend/captchar"
	"github.com/empirefox/esecend/cdn"
	"github.com/empirefox/esecend/config"
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/esecend/front"
	"github.com/empirefox/esecend/hub"
	"github.com/empirefox/esecend/search"
	"github.com/empirefox/esecend/sec"
	"github.com/empirefox/esecend/server"
	"github.com/empirefox/esecend/sms"
	"github.com/empirefox/esecend/wo2"
	"github.com/empirefox/esecend/wx"
	"github.com/gin-gonic/gin"
)

var (
	isDevMode bool
	isTLS     bool

	log = logrus.New()
)

func init() {
	isDevMode, _ = strconv.ParseBool(os.Getenv("DEV_MODE"))
	isTLS, _ = strconv.ParseBool(os.Getenv("TLS"))
}

// PORT=9999 DEV_MODE=1 CONFIG5=~/soft/workspace/prouction/esecend/esecend-config.json5 ./goprd
func main() {
	if !isDevMode {
		gin.SetMode(gin.ReleaseMode)
	}
	s, err := newServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(s.StartRun())
}

func newServer() (*server.Server, error) {
	configFile := os.Getenv("CONFIG5")
	if configFile == "" {
		configFile = "./config.json5"
	}

	// 1.
	conf, err := config.Load(configFile)
	if err != nil {
		return nil, err
	}

	// 3.
	wxClient, err := wx.NewWxClient(conf)
	if err != nil {
		return nil, err
	}

	// 4.
	dbs, err := dbsrv.NewDbService(conf, wxClient, isDevMode)
	if err != nil {
		return nil, err
	}

	err = dbs.LoadProfile()
	if err != nil {
		return nil, err
	}

	//	paykey, err := models.EncPaykey([]byte("123456"), 10)
	//	if err != nil {
	//		panic(err)
	//	}
	//	err = dbs.GetDB().UpdateColumns(&models.User{
	//		ID:     1,
	//		Paykey: paykey,
	//	}, "Paykey")
	//	if err != nil {
	//		panic(err)
	//	}

	// 5.
	captcha, err := captchar.NewCaptchar("./comic.ttf")
	if err != nil {
		return nil, err
	}

	// 6.
	secHandler := security.NewHandler(conf, dbs)

	// 12.
	newsResource := &search.Resource{
		Conf: conf,
		Dbs:  dbs,
		View: front.NewsItemTable,
	}
	newsResource.SetDefaultFilters()
	newsResource.SearchAttrs(front.NewsItemTable.Fields()...)

	// 13.
	productResource := &search.Resource{
		Conf: conf,
		Dbs:  dbs,
		View: front.ProductTable,
	}
	productResource.SetDefaultFilters()
	productResource.SearchAttrs("Name", "Intro", "Detail")

	// 14.
	orderResource := &search.Resource{
		Conf: conf,
		Dbs:  dbs,
		View: front.OrderTable,
	}
	orderResource.SetDefaultFilters()
	// TODO search from items
	//	orderResource.SearchAttrs(front.OrderItemTable.Fields()...)

	s := &server.Server{
		IsDevMode:  isDevMode,
		IsTLS:      isTLS,
		Config:     conf,
		Cdn:        cdn.NewQiniu(conf), // 2.
		DB:         dbs,
		WxClient:   wxClient,
		Captcha:    captcha,
		SecHandler: secHandler,
		Auther:     wo2.NewAuther(conf, secHandler), // 7.
		Admin:      admin.NewAdmin(conf),            // 8.
		SmsSender:  sms.NewSender(conf),             // 9.
		ProductHub: hub.NewProductHub(dbs),          // 10.
		OrderHub:   hub.NewOrderHub(conf, dbs),      // 11.

		NewsResource:    newsResource,    // 12.
		ProductResource: productResource, // 13.
		OrderResource:   orderResource,   // 14.
	}

	s.BuildEngine()

	return s, nil
}
