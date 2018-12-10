package main

import (
	"fmt"
	"github.com/GeoNet/kit/metrics"
	"github.com/GeoNet/kit/mseed"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func save(inbound chan []byte) {
	msr := mseed.NewMSRecord()
	defer mseed.FreeMSRecord(msr)

	var file string
	var err error

	for {
		select {
		case b := <-inbound:

			// TODO add timing in
			//t := metrics.Start()

			err = msr.Unpack(b, 512, 0, 0)
			if err != nil {
				metrics.MsgErr()
				log.Printf("unpacking miniSEED record: %s", err.Error())
				continue
			}

			// TODO path join
			file = fmt.Sprintf("%s/%s/%s/%s/%s/%s_%s", dir, msr.Network(), msr.Station(), msr.Location(), msr.Channel(), msr.Starttime().Format(time.RFC3339Nano), msr.Endtime().Format(time.RFC3339Nano))

			err := ioutil.WriteFile(file, b, 0644)
			if err == nil {
				metrics.MsgProc()
				continue
			}

			// Ignore any errors from building the directory structure - there are multiple consumers
			// so there is a race here.
			_ = os.Mkdir(fmt.Sprintf("%s/%s", dir, msr.Network()), 0755)
			_ = os.Mkdir(fmt.Sprintf("%s/%s/%s", dir, msr.Network(), msr.Station()), 0755)
			_ = os.Mkdir(fmt.Sprintf("%s/%s/%s/%s", dir, msr.Network(), msr.Station(), msr.Location()), 0755)
			_ = os.Mkdir(fmt.Sprintf("%s/%s/%s/%s/%s", dir, msr.Network(), msr.Station(), msr.Location(), msr.Channel()), 0755)

			err = ioutil.WriteFile(file, b, 0644)
			if err == nil {
				metrics.MsgProc()
				continue
			}

			metrics.MsgErr()
			log.Println(err)
		}
	}
}

// expire removes old data from the DB.  The archive runs 7 days between real time.  Keep
// 8 days to allow some overlap.
//func (a *app) expire() {
//	ticker := time.NewTicker(time.Minute).C
//	var err error
//	for {
//		select {
//		case <-ticker:
//			_, err = a.db.Exec(`DELETE FROM fdsn.record WHERE start_time < now() - interval '8 days'`)
//			if err != nil {
//				log.Printf("deleting old records: %s", err.Error())
//			}
//		}
//	}
//}
