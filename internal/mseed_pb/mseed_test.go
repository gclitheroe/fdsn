package mseed_pb

import (
	"github.com/GeoNet/kit/mseed"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"testing"
)

var results = []struct {
	id    string
	file  string
	count int
	start int64 // start time in Unix nanoseconds
	end   int64 // end time in Unix nanoseconds
}{
	{
		id:    id(),
		file:  "etc/NZ.ABAZ.10.EHE.D.2016.079",
		count: 20532,
		start: 1458345601968393000,
		end:   1458432002998391000,
	},
}

func TestIndex(t *testing.T) {
	for _, r := range results {
		f, err := os.Open(r.file)
		if err != nil {
			t.Fatalf("%s %s\n", r.id, err)
		}

		var idx Index

		err = idx.SingleStream(f)
		if err != nil {
			t.Fatalf("%s SingleStream %s\n", r.id, err)
		}

		if r.count != len(idx.Records) {
			t.Errorf("%s expected %d records got %d\n", r.id, r.count, len(idx.Records))
		}

		if r.start != idx.Records[0].Start.UnixNano {
			t.Errorf("%s expected start %d got %d\n", r.id, r.start, idx.Records[0].Start.UnixNano)
		}

		if r.end != idx.Records[len(idx.Records)-1].End.UnixNano {
			t.Errorf("%s expected end %d got %d\n", r.id, r.end, idx.Records[len(idx.Records)-1].Start.UnixNano)
		}

		o, err := f.Seek(10*512, 0)
		if err != nil {
			log.Print(err)
		}

		t.Log(o)

		msr := mseed.NewMSRecord()
		record := make([]byte, 512)

		_, err = io.ReadFull(f, record)
		if err != nil {
			log.Print(err)
		}

		err = msr.Unpack(record, 512, 1, 1)
		if err != nil {
			log.Print(err)
		}

		t.Log(msr.Starttime())
		t.Log(msr.Endtime())
		t.Log(msr.DataSamples())

	}
}

func id() string {
	_, _, l, _ := runtime.Caller(1)
	return "L" + strconv.Itoa(l)
}
