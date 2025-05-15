package digitalocean

import (
	"context"
	"strconv"
	"sync"

	"github.com/digitalocean/godo"
	"github.com/libdns/libdns"
)

type Client struct {
	client *godo.Client
	mutex  sync.Mutex
}

func (p *Provider) getClient() error {
	if p.client == nil {
		p.client = godo.NewFromToken(p.APIToken)
	}

	return nil
}

func (p *Provider) getDNSEntries(ctx context.Context, zone string) ([]libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.getClient()

	opt := &godo.ListOptions{}
	var records []libdns.Record
	for {
		domains, resp, err := p.client.Domains.Records(ctx, zone, opt)
		if err != nil {
			return records, err
		}

		for _, entry := range domains {
			record := godoToRecord(entry)
			records = append(records, record)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return records, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return records, nil
}

func (p *Provider) addDNSEntry(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.getClient()

	rr := record.RR()
	entry := godo.DomainRecordEditRequest{
		Name: rr.Name,
		Data: rr.Data,
		Type: rr.Type,
		TTL:  int(rr.TTL.Seconds()),
	}

	rec, _, err := p.client.Domains.CreateRecord(ctx, zone, &entry)
	if err != nil {
		return record, err
	}

	return fromRecord(record, strconv.Itoa(rec.ID)), nil
}

func (p *Provider) removeDNSEntry(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.getClient()

	// Get ID from dns record
	var idRaw string
	if dnsRecord, ok := record.(DNS); ok {
		idRaw = dnsRecord.ID
	}

	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return record, err
	}

	_, err = p.client.Domains.DeleteRecord(ctx, zone, id)
	if err != nil {
		return record, err
	}

	return record, nil
}

func (p *Provider) updateDNSEntry(ctx context.Context, zone string, record libdns.Record) (libdns.Record, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.getClient()

	// Get ID from dns record
	var idRaw string
	if dnsRecord, ok := record.(DNS); ok {
		idRaw = dnsRecord.ID
	}

	id, err := strconv.Atoi(idRaw)
	if err != nil {
		return record, err
	}

	rr := record.RR()
	entry := godo.DomainRecordEditRequest{
		Name: rr.Name,
		Data: rr.Data,
		Type: rr.Type,
		TTL:  int(rr.TTL.Seconds()),
	}

	_, _, err = p.client.Domains.EditRecord(ctx, zone, id, &entry)
	if err != nil {
		return record, err
	}

	return record, nil
}
