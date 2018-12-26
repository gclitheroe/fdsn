package mseed_pb

import (
	"github.com/GeoNet/kit/mseed"
	"io"
)

// the record length of the miniSEED records.  Constant for all GNS miniSEED files.
const recordLength int = 512

func (m *Index) SingleStream(r io.Reader) error {
	msr := mseed.NewMSRecord()
	defer mseed.FreeMSRecord(msr)

	record := make([]byte, recordLength)
	var err error
	var i int64

loop:
	for {
		_, err = io.ReadFull(r, record)
		switch {
		case err == io.EOF:
			break loop
		case err != nil:
			return err
		}

		err = msr.Unpack(record, recordLength, 1, 0)
		if err != nil {
			return err
		}

		m.Records = append(m.Records, &Record{
			Number: i,
			Start: &Timestamp{
				UnixNano: msr.Starttime().UnixNano(),
			},
			End: &Timestamp{
				UnixNano: msr.Endtime().UnixNano(),
			},
		})

		i++
	}

	return nil
}
