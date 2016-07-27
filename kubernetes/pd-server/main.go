// Copyright 2016 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/pingcap/pd/server"
	"github.com/pingcap/pd/server/api"
)

var (
	config                 = flag.String("c", "", "config file")
	addr                   = flag.String("addr", "127.0.0.1:1234", "server listening address")
	advertiseAddr          = flag.String("advertise-addr", "", "server advertise listening address [127.0.0.1:1234] for client communication")
	etcdAddrs              = flag.String("etcd", "127.0.0.1:2379", "Etcd endpoints, separated by comma")
	httpAddr               = flag.String("http-addr", ":9090", "http server listening address")
	pprofAddr              = flag.String("pprof", ":6060", "pprof HTTP listening address")
	rootPath               = flag.String("root", "/pd", "pd root path in etcd")
	leaderLease            = flag.Int64("lease", 3, "leader lease time (second)")
	logLevel               = flag.String("L", "debug", "log level: info, debug, warn, error, fatal")
	tsoSaveInterval        = flag.Int64("tso-save-interval", 2000, "the interval time (ms) to save timestamp")
	clusterID              = flag.Uint64("cluster-id", 0, "cluster ID")
	maxPeerCount           = flag.Uint64("max-peer-count", 3, "max peer count for the region")
	metricAddr             = flag.String("metric-addr", "", "metric address")
	minCapUsedRatio        = flag.Float64("min-capacity-used-ratio", 0.4, "min capacity used ratio for choosing store in balance")
	maxCapUsedRatio        = flag.Float64("max-capacity-used-ratio", 0.9, "max capacity used ratio for choosing store in balance")
	maxSendSnapCount       = flag.Uint64("max-sending-snap-count", 3, "max sending snapshot count for choosing store in balance")
	maxRecvSnapCount       = flag.Uint64("max-receiving-snap-count", 3, "max receiving snapshot count for choosing store in balance")
	maxDiffScoreFrac       = flag.Float64("max-diff-score-fraction", 0.1, "max diff score fraction for choosing store in balance")
	balanceInterval        = flag.Uint64("balance-interval", 30, "the interval time (s) to do balance")
	maxBalanceCount        = flag.Uint64("max-balance-count", 16, "the max region count to balance at the same time")
	maxBalanceRetryPerLoop = flag.Uint64("max-balance-retry-per-loop", 10, "the max retry count to balance in a balance schedule")
	maxBalanceCountPerLoop = flag.Uint64("max-balance-count-per-loop", 3, "the max region count to balance in a balance schedule")
)

func setCmdArgs(cfg *server.Config) {
	flag.Visit(func(flag *flag.Flag) {
		flagArgs[flag.Name] = true
	})

	setStringFlagConfig(&cfg.Addr, "addr", *addr)
	setStringFlagConfig(&cfg.AdvertiseAddr, "advertise-addr", *advertiseAddr)
	setStringSliceFlagConfig(&cfg.EtcdAddrs, "etcd", *etcdAddrs)
	setStringFlagConfig(&cfg.HTTPAddr, "http-addr", *httpAddr)
	setStringFlagConfig(&cfg.PprofAddr, "pprof", *pprofAddr)
	setStringFlagConfig(&cfg.RootPath, "root", *rootPath)
	setIntFlagConfig(&cfg.LeaderLease, "lease", *leaderLease)
	setStringFlagConfig(&cfg.LogLevel, "L", *logLevel)
	setIntFlagConfig(&cfg.TsoSaveInterval, "tso-save-interval", *tsoSaveInterval)
	setUintFlagConfig(&cfg.ClusterID, "cluster-id", *clusterID)
	setUintFlagConfig(&cfg.MaxPeerCount, "max-peer-count", *maxPeerCount)
	setStringFlagConfig(&cfg.MetricAddr, "metric-addr", *metricAddr)
	setFloatFlagConfig(&cfg.BalanceCfg.MinCapacityUsedRatio, "min-capacity-used-ratio", *minCapUsedRatio)
	setFloatFlagConfig(&cfg.BalanceCfg.MaxCapacityUsedRatio, "max-capacity-used-ratio", *maxCapUsedRatio)
	setUintFlagConfig(&cfg.BalanceCfg.MaxSendingSnapCount, "max-sending-snap-count", *maxSendSnapCount)
	setUintFlagConfig(&cfg.BalanceCfg.MaxReceivingSnapCount, "max-receiving-snap-count", *maxRecvSnapCount)
	setFloatFlagConfig(&cfg.BalanceCfg.MaxDiffScoreFraction, "max-diff-score-fraction", *maxDiffScoreFrac)
	setUintFlagConfig(&cfg.BalanceCfg.BalanceInterval, "balance-interval", *balanceInterval)
	setUintFlagConfig(&cfg.BalanceCfg.MaxBalanceCount, "max-balance-count", *maxBalanceCount)
	setUintFlagConfig(&cfg.BalanceCfg.MaxBalanceRetryPerLoop, "max-balance-retry-per-loop", *maxBalanceRetryPerLoop)
	setUintFlagConfig(&cfg.BalanceCfg.MaxBalanceCountPerLoop, "max-balance-count-per-loop", *maxBalanceCountPerLoop)
}

func setEnvVar(cfg *server.Config) {
        if val := os.Getenv("${Q8S_ETCD_ENVVAR}"); !(flagArgs["etcd"] || !useConfigFile) && val != "" {
                arr := strings.Split(val, ",")
                for i := 0; i < len(arr); i++ {
                        if v := os.Getenv(arr[i]); v != "" {
                            arr[i] = v
                        }
                }
                cfg.EtcdAddrs = arr    
        }
}

func main() {
	flag.Parse()

	cfg := server.NewConfig()

	if *config != "" {
		if err := cfg.LoadFromFile(*config); err != nil {
			log.Fatalf("load config failed - %s", err)
		}

		useConfigFile = true
		log.Infof("PD init config - %v", cfg)
	}

	setCmdArgs(cfg)

        setEnvVar(cfg)

	log.SetLevelByString(cfg.LogLevel)
	log.SetHighlighting(false)

	log.Infof("PD config - %v", cfg)

	go func() {
		http.ListenAndServe(cfg.PprofAddr, nil)
	}()

	svr, err := server.NewServer(cfg)
	if err != nil {
		log.Errorf("create pd server err %s\n", err)
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-sc
		log.Infof("Got signal [%d] to exit.", sig)
		svr.Close()
		os.Exit(0)
	}()

	go func() {
		err = api.ServeHTTP(cfg.HTTPAddr, svr)
		if err != nil {
			log.Fatalf("serve http failed - %v", errors.Trace(err))
		}
	}()

	err = svr.Run()
	if err != nil {
		log.Fatalf("server run failed - %v", errors.Trace(err))
	}
}

var (
	flagArgs      = map[string]bool{}
	useConfigFile = false
)

func setStringFlagConfig(dest *string, name string, value string) {
	if flagArgs[name] || !useConfigFile {
		*dest = value
	}
}

func setStringSliceFlagConfig(dest *[]string, name string, value string) {
	if flagArgs[name] || !useConfigFile {
		*dest = append([]string{}, strings.Split(value, ",")...)
	}
}

func setIntFlagConfig(dest *int64, name string, value int64) {
	if flagArgs[name] || !useConfigFile {
		*dest = value
	}
}

func setUintFlagConfig(dest *uint64, name string, value uint64) {
	if flagArgs[name] || !useConfigFile {
		*dest = value
	}
}

func setFloatFlagConfig(dest *float64, name string, value float64) {
	if flagArgs[name] || !useConfigFile {
		*dest = value
	}
}
