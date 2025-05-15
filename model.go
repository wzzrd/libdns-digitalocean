package digitalocean

import (
	"strconv"
	"time"

	"github.com/digitalocean/godo"
	"github.com/libdns/libdns"
)

// dns custom struct that implements the libdns.Record interface and keeps the ID field used internally
type dns struct {
	Record libdns.RR
	ID     string
}

func (d dns) RR() libdns.RR {
	return d.Record
}

// fromRecord creates a dns struct from a libdns.RR, with an optional ID
func fromRecord(record libdns.Record, id string) dns {
	rr := record.RR()
	return dns{
		Record: rr,
		ID:     id,
	}
}

// recordToGoDo converts a libdns.RR to the DigitalOcean API format
func recordToGoDo(record libdns.Record) godo.DomainRecordEditRequest {
	rr := record.RR()
	return godo.DomainRecordEditRequest{
		Name: rr.Name,
		Data: rr.Data,
		Type: rr.Type,
		TTL:  int(rr.TTL.Seconds()),
	}
}

// godoToRecord converts a DigitalOcean DNS record to dns type
func godoToRecord(entry godo.DomainRecord) dns {
	rr := libdns.RR{
		Name: entry.Name,
		Data: entry.Data,
		Type: entry.Type,
		TTL:  time.Duration(entry.TTL) * time.Second,
	}

	return dns{
		Record: rr,
		ID:     strconv.Itoa(entry.ID),
	}
}
