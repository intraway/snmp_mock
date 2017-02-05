package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
	"gopkg.in/errgo.v1"

	"github.com/posteo/go-agentx"
	"github.com/posteo/go-agentx/value"
)

var log = logging.MustGetLogger("iostat_monitor")
var format = logging.MustStringFormatter(
	`%{time:2006-02-01 15:04:05.000} %{level:.4s} %{message}`,
)

func main() {
	logging.SetFormatter(format)
	log_out_backend := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetBackend(log_out_backend)

	log.Info("Running snmp mock")
	oid_files := os.Args[1:]

	config, err := LoadConfig("/app/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	client := &agentx.Client{
		Net:               "tcp",
		Address:           "localhost:705",
		Timeout:           1 * time.Minute,
		ReconnectInterval: 1 * time.Second,
	}

	if err := client.Open(); err != nil {
		log.Fatalf(errgo.Details(err))
	}

	session, err := client.Session()
	if err != nil {
		log.Fatalf(errgo.Details(err))
	}

	snmp_handler := &SNMPHandler{}
	err = LoadOids(snmp_handler, oid_files...)
	if err != nil {
		log.Fatal(err)
	}

	session.Handler = snmp_handler

	if err := session.Register(127, value.MustParseOID(config.BaseOid)); err != nil {
		log.Fatalf(errgo.Details(err))
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
