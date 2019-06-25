package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/OpenPeeDeeP/xdg"
	"github.com/sstark/knxbaosip"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	confFileName = "kxsh.conf"
)

type GroupMap map[string][]int
type ConfigMap map[string]string

// JoinInts returns a string with the int elements of l joined together using c
func JoinInts(l []int, c string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(l)), c), "[]")
}

func listGroups(groups GroupMap) {
	var ell string
	for k, v := range groups {
		ml := len(v)
		if ml > 8 {
			ell = " ..."
			ml = 8
		} else {
			ell = ""
		}
		l := JoinInts(v[:ml], " ")
		fmt.Printf("%12s: %s%s\n", k, l, ell)
	}
}

func readDatapoints(knx *knxbaosip.Client, datapoints []int, tsv bool) {
	err, ds := knx.GetDescriptionString(datapoints)
	if err != nil {
		log.Fatal(err)
	}
	dsm := make(map[int]string)
	for _, n := range ds {
		dsm[n.Datapoint] = n.Description
	}
	err, dpv := knx.GetDatapointValue(datapoints)
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range dpv {
		if tsv {
			fmt.Printf("%5d\t%5s\t%-30s\t%s\n", d.Datapoint, d.Format, dsm[d.Datapoint], string(d.Value))
		} else {
			fmt.Printf("%5d %5s |%-30s| %s\n", d.Datapoint, d.Format, dsm[d.Datapoint], string(d.Value))
		}
	}
}

func main() {
	var cg GroupMap
	var cf ConfigMap
	var err error

	conf := xdg.New("", "kxsh")
	confFile := conf.QueryConfig(confFileName)
	if confFile != "" {
		fmt.Printf("using config from %s\n", confFile)
		b, err := ioutil.ReadFile(confFile)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(b, &cf)
		if err != nil {
			log.Fatal(err)
		}
	}

	var fGroup, fUrl, fDps string
	var fList, fInteractive, fTsv bool
	flag.StringVar(&fGroup, "group", "default", "choose datapoint group")
	flag.StringVar(&fUrl, "url", cf["url"], "specify URL of BAOS device")
	flag.StringVar(&fDps, "dps", "dps.json", "datapoint file")
	flag.BoolVar(&fList, "list", false, "show available groups")
	flag.BoolVar(&fInteractive, "i", false, "interactive mode")
	flag.BoolVar(&fTsv, "tsv", false, "tab separated output")
	flag.Parse()

	b, err := ioutil.ReadFile(fDps)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &cg)
	if err != nil {
		log.Fatal(err)
	}

	if fList {
		listGroups(cg)
		os.Exit(0)
	}

	if cf["user"] != "" && cf["pass"] != "" {
		newUrl, err := url.Parse(fUrl)
		if err != nil {
			log.Fatalf("could not parse url: %v", err)
		}
		newUrl.User = url.UserPassword(cf["user"], cf["pass"])
		fUrl = newUrl.String()
	}

	knx := knxbaosip.NewClient(fUrl)
	err, si := knx.GetServerItem()
	if err != nil {
		if err == knxbaosip.AuthError {
			log.Fatalf("%v: Please supply user and pass in the config file.", err)
		}
		log.Fatal(err)
	}

	sn := JoinInts(si.SerialNumber, ".")
	fmt.Printf("%s fw:%d sn:%v\n", knx.Url, si.FirmwareVersion, sn)

	if fInteractive {
		err = prompt(knx, cg)
		if err == nil {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	datapoints := cg[fGroup]
	if len(datapoints) == 0 {
		log.Fatalf("no datapoints in group \"%s\"", fGroup)
	}

	readDatapoints(knx, datapoints, fTsv)
}
