package server

import (
    "fmt"
    "synch"
)

type Log struct {
    mu  sync.Mutex
    records []Record
}

func NewLog() *Log {
    return &Log{}
}

// Appending a record just tacks data on the end of the slice.
func (c *Log) Append(record Record) (unit64, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
    record.Offset = unit64(len(c.records))
    return record.Offset, nil
}

func (c *Log) Read(offset unit64) (Record, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
    if offset >= unit64(len(c.records)) {
        return Record{}, ErrOffsetNotFound
    }
    return c.records[offset], nil
}

type Record struct {
    Value []byte `json:"value"`
    Offset unit64 `json:"offset"`
} 

var ErrorOffsetNotFound = fmt.Errorf("offset not found")

