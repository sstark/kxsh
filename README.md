
kxsh - simple shell using go bindings for KNX BAOS IP Service
=============================================================

    > ./kxsh -h
    Usage of ./kxsh:
    -dps string
            datapoint file (default "dps.json")
    -group string
            choose datapoint group (default "default")
    -url string
            specify URL of BAOS device

Needs the knxbaosip go bindings installed:

    go get github.com/sstark/knxbaosip

Example
=======

    > ./kxsh -dps dps.json -url http://localhost:8888/baos/ -group n4015
    http://localhost:8888/baos/ fw:51 sn:0.197.1.1.116.14
    102  DPT1 |Bel. N4.015 Gr.1 E/A Status   | false
    107  DPT1 |Bel. N4.015 Gr.2 E/A Status   | false
    276  DPT1 |Magnetkontakt N4.015 Fenster 1| true
    486  DPT5 |Jalo. N4.015 Behang Status. Ab| 0
    608  DPT1 |Heizung N4.016 Stellgr. Heizen| false
    708 DPT14 |RFT N4.015 Absolute Feuchte g/| 0.000000
    709 DPT14 |RFT N4.015 Absolute Feuchte g/| 0.000000
    710  DPT1 |RFT N4.015 Raumklima Behaglich| false
    711  DPT9 |OFT N4.015 Interner Temperatur| 22.4333
    713  DPT9 |RFT N4.016 Interner Feuchtemes| 0.000000
