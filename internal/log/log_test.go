package log

import (
	api "github.com/joostvdg/proglog/api/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T, log *Log){
		"append and read a record succeeds": testAppendRead,
		"offset out of range error":         testOutOfRangeErr,
		"init with existing segments":       testInitExisting,
		"reader":                            testReader,
		"truncate":                          testTruncate,
	} {
		t.Run(scenario, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "store-test")
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			c := Config{}
			c.Segment.MaxStoreBytes = 32
			log, err := NewLog(dir, c)
			require.NoError(t, err)

			fn(t, log)
		})
	}
}

// testAppendRead tests that we can successfully append to and from the log.
// When we append a record to the log, it returns the offset of the appended record.
// We should thus be able to read that same record from the log at that offset value.
func testAppendRead(t *testing.T, log *Log) {
	valueToAppend := &api.Record{
		Value: []byte("hello world"),
	}
	off, err := log.Append(valueToAppend)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	read, err := log.Read(off)
	require.NoError(t, err)
	require.Equal(t, valueToAppend.Value, read.Value)
}

// testOutOfRangeErr tests that we get an error when we attempt to read a record at an offset that does not exist.
// Assuming we have created no records, if we read at offset one (e.g., the first record) it should be out of range.
func testOutOfRangeErr(t *testing.T, log *Log) {
	read, err := log.Read(1)
	require.Nil(t, read)

	apiErr := err.(api.ErrOffsetOutOfRange)
	require.Equal(t, uint64(1), apiErr.Offset)
}

// testInitExisting tests that when we create a log, the log bootstraps itself with data from previous instances.
// As we append three records to the initial log, we should have the appropriate offsets in a "re-create" log.
// Namely, lowestOffset is 0 (as we have not tranced) and highest is 2 (0 index based).
func testInitExisting(t *testing.T, log *Log) {
	valueToAppend := &api.Record{
		Value: []byte("hello world"),
	}
	// Append three records, so lowest Off = 0, highest = 2 (0, 1, 2)
	for i := 0; i < 3; i++ {
		_, err := log.Append(valueToAppend)
		require.NoError(t, err)
	}
	require.NoError(t, log.Close())

	off, err := log.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)
	off, err = log.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), off)

	newLog, err := NewLog(log.Dir, log.Config)
	require.NoError(t, err)

	off, err = newLog.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)
	off, err = newLog.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), off)
}

func testReader(t *testing.T, log *Log) {
	valueToAppend := &api.Record{
		Value: []byte("hello world"),
	}
	off, err := log.Append(valueToAppend)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	reader := log.Reader()
	b, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	read := &api.Record{}
	err = proto.Unmarshal(b[lenWidth:], read)
	require.NoError(t, err)
	require.Equal(t, valueToAppend.Value, read.Value)
}

func testTruncate(t *testing.T, log *Log) {
	valueToAppend := &api.Record{
		Value: []byte("hello world"),
	}
	for i := 0; i < 3; i++ {
		_, err := log.Append(valueToAppend)
		require.NoError(t, err)
	}
	// we now truncate to log, essentially removing all logs < 1
	err := log.Truncate(1)
	require.NoError(t, err)
	// which means that there is no log at offset 0
	_, err = log.Read(0)
	require.Error(t, err)
}
