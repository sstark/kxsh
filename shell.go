package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sstark/knxbaosip"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type GroupMap map[string][]int

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
		l := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v[:ml])), " "), "[]")
		fmt.Printf("%12s: %s%s\n", k, l, ell)
	}
}

func main() {
	var cg GroupMap
	var err error

	var fGroup, fUrl, fDps string
	var fLsGroups bool
	flag.StringVar(&fGroup, "group", "default", "choose datapoint group")
	flag.StringVar(&fUrl, "url", "", "specify URL of BAOS device")
	flag.StringVar(&fDps, "dps", "dps.json", "datapoint file")
	flag.BoolVar(&fLsGroups, "lsgroups", false, "show available groups")
	flag.Parse()

	b, err := ioutil.ReadFile(fDps)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &cg)
	if err != nil {
		log.Fatal(err)
	}

	if fLsGroups {
		listGroups(cg)
		os.Exit(0)
	}

	knx := knxbaosip.NewClient(fUrl)
	err, si := knx.GetServerItem()
	if err != nil {
		log.Fatal(err)
	}
	sn := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(si.SerialNumber)), "."), "[]")
	fmt.Printf("%s fw:%d sn:%v\n", knx.Url, si.FirmwareVersion, sn)

	datapoints := cg[fGroup]
	if len(datapoints) == 0 {
		log.Fatalf("no datapoints in group \"%s\"", fGroup)
	}
	err, ds := knx.GetDescriptionString(datapoints)
	if err != nil {
		log.Fatal(err)
	}
	err, dpv := knx.GetDatapointValue(datapoints)
	if err != nil {
		log.Fatal(err)
	}
	for i, d := range dpv {
		desc := ds[i].Description
		fmt.Printf("%5d %5s |%-30s| %s\n", d.Datapoint, d.Format, desc, string(d.Value))
	}
}
