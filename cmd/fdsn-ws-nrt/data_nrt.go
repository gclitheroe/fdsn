package main

import (
	"errors"
	"fmt"
	"github.com/GeoNet/fdsn/internal/fdsn"
	"github.com/GeoNet/kit/metrics"
	"github.com/golang/groupcache"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var errNoData = errors.New("no data")

// holdingsSearchNrt searches for near real time records matching the query.
// start and end should be set for all queries.
func holdingsSearchNrt(d fdsn.DataSelect) ([]string, error) {
	timer := metrics.Start()
	defer timer.Track("holdingsSearchNrt")

	// TODO - list the dir etc etc

	log.Printf("%+v\n", d)

	return keyList(strings.Join(d.Network, ""), strings.Join(d.Station, ""), strings.Join(d.Location, ""), strings.Join(d.Channel, ""), d.StartTime.Time, d.EndTime.Time)
}

func keyList(network, station, location, channel string, start, end time.Time) ([]string, error) {
	dir := fmt.Sprintf("/work/fdsn-nrt/%s/%s/%s/%s", network, station, location, channel)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}, nil
	}

	var f []string

	for _, v := range files {
		s := strings.Split(v.Name(), "_")

		if len(s) != 2 {
			return []string{}, err
		}

		startFile, err := time.Parse(time.RFC3339Nano, s[0])
		if err != nil {
			return []string{}, err
		}

		endFile, err := time.Parse(time.RFC3339Nano, s[1])
		if err != nil {
			return []string{}, err
		}

		if startFile.After(start) && endFile.Before(end) {
			f = append(f, dir+"/"+v.Name())
		}
	}

	return f, nil
}

// recordGetter implements groupcache.Getter for fetching miniSEED records from the cache.
// key is like "NZ_AWRB_HNN_23_2017-04-22T22:38:50.115Z"
// network_station_channel_location_time.RFC3339Nano
func recordGetter(ctx groupcache.Context, key string, dest groupcache.Sink) error {
	b, err := ioutil.ReadFile(key)
	if err != nil {
		return err
	}

	dest.SetBytes(b)
	return nil
}
