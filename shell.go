package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sstark/knxbaosip"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	var cg map[string][]int
	var err error

	b, err := ioutil.ReadFile("dps.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &cg)
	if err != nil {
		log.Fatal(err)
	}

	var fGroup string
	flag.StringVar(&fGroup, "group", "default", "choose datapoint group")
	flag.Parse()

	knx := knxbaosip.NewClient("")
	err, si := knx.GetServerItem()
	if err != nil {
		log.Fatal(err)
	}
	sn := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(si.SerialNumber)), "."), "[]")
	fmt.Printf("%s fw:%d sn:%v\n", knx.Url, si.FirmwareVersion, sn)

	datapoints := cg[fGroup]
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
		fmt.Printf("%5d %5s \"%-32s\": %s\n", d.Datapoint, d.Format, desc, string(d.Value))
	}
}
