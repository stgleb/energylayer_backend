package influxdb_storage

import (
	"github.com/influxdata/influxdb/client/v2"
	"time"
)

type MockClient struct {
	Called map[string]int

	FakePing  func(time.Duration) (time.Duration, string, error)
	FakeWrite func(client.BatchPoints) error
	FakeQuery func(client.Query) (*client.Response, error)
	FakeClose func() error
}

func (mock MockClient) Ping(timeout time.Duration) (time.Duration, string, error) {
	return mock.FakePing(timeout)
}

func (mock MockClient) Write(bp client.BatchPoints) error {
	return mock.FakeWrite(bp)
}

func (mock MockClient) Query(q client.Query) (*client.Response, error) {
	println(mock.Called)
	return mock.FakeQuery(q)
}

func (mock MockClient) Close() error {
	return mock.FakeClose()
}

func (mock *MockClient) Refresh() {
	mock.Called = make(map[string]int)
}
