
kxsh - simple shell using go bindings for KNX BAOS IP Service
=============================================================

    > kxsh -h
    Usage of kxsh:
    -dps string
            datapoint file (default "dps.json")
    -group string
            choose datapoint group (default "default")
    -i	interactive mode
    -list
            show available groups
    -url string
            specify URL of BAOS device

Needs the knxbaosip go bindings installed:

    go get github.com/sstark/knxbaosip

Example
=======

    > kxsh -dps dps.json -url http://localhost:8888/baos/ -group n4015
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


    > kxsh -i
    >>> list
        n4016: 24 27 28
        default: 102 700 701 711 712 720 721 722
            alle: 1 2 3 4 5 6 7 8 ...
        praesenz: 748 749 750 751 752 753 754 755 ...
        n4015: 101 102 103 104 105 106 107 108 ...
    >>> group n4015
    [n4015] read 107
    107  DPT1 |Bel. N4.015 Gr.2 E/A Status   | true
