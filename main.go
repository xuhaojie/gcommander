package main

import (
	"flag"
	"fmt"
	"gcommander/config"
	"gcommander/waker"
	"github.com/kirsle/configdir"
	"path/filepath"
)

func main() {

	configFile := filepath.Join(configdir.LocalConfig(), "gcommander.cfg")

	var g = flag.Bool("g", false, "generate sample config")
	var l = flag.Bool("l", false, "list targets")
	var m = flag.String("m", "", "mac address of target to wake")
	var n = flag.String("n", "", "name of target to wake")

	flag.Parse()

	if *g {
		defaultCfg := config.Default()
		err := defaultCfg.Save(configFile)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("sample config file", configFile, "generated.")
		}
		return
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Println("can't load config, ", err)
		return
	}

	if *l {
		for _, t := range cfg.Targets {
			fmt.Printf("%10s : %s\n", t.Name, t.Mac)
		}
		return
	}

	var target string
	var mac string
	if len(*m) > 0 {
		mac = *m
	} else if len(*n) > 0 {
		mac, err = cfg.FindTarget(*n)
		if err != nil {
			fmt.Println("target not found!")
			return
		}
	} else {
		fmt.Println("please specific target by mac or name.")
		flag.Usage()
		return
	}

	w, err := waker.New(cfg.Url, cfg.User, cfg.Password)
	if err != nil {
		fmt.Println("create waker failed, ", err)
		return
	}

	if len(target) > 0 {
		fmt.Printf("Wake %s with mac %s ...\n", target, mac)
	} else {
		fmt.Printf("Wake %s ...\n", mac)
	}

	err = w.Wake(mac)
	if err != nil {
		fmt.Println("failed! ", err)
		return
	} else {
		fmt.Println("done.")
	}
}
