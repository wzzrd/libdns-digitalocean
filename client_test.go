package digitalocean

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/libdns/libdns"
)

func TestClient_getDNSEntries(t *testing.T) {
	// Mock domain records to return
	mockRecords := []godo.DomainRecord{
		{
			ID:   1,
			Type: "A",
			Name: "test",
			Data: "192.168.1.1",
			TTL:  3600,
		},
		{
			ID:   2,
			Type: "CNAME",
			Name: "www",
			Data: "example.com",
			TTL:  1800,
		},
	}

	// Test successful call
	p := setupTest(mockRecords, nil)
	ctx := context.Background()

	records, err := p.getDNSEntries(ctx, "example.com")

	if err != nil {
		t.Errorf("Client.getDNSEntries() error = %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Client.getDNSEntries() returned %d records, want 2", len(records))
	}

	// Verify first record
	if records[0].Type != "A" || records[0].Name != "test" ||
		records[0].Value != "192.168.1.1" || records[0].ID != "1" {
		t.Errorf("Client.getDNSEntries()[0] = %v, want A record", records[0])
	}

	// Verify second record
	if records[1].Type != "CNAME" || records[1].Name != "www" ||
		records[1].Value != "example.com" || records[1].ID != "2" {
		t.Errorf("Client.getDNSEntries()[1] = %v, want CNAME record", records[1])
	}

	// Test error case
	p = setupTest(nil, errors.New("API error"))

	_, err = p.getDNSEntries(ctx, "example.com")
	if err == nil {
		t.Error("Client.getDNSEntries() expected error, got nil")
	}
}

func TestClient_addDNSEntry(t *testing.T) {
	// Test record to add
	testRecord := libdns.Record{
		Type:  "A",
		Name:  "test",
		Value: "192.168.1.1",
		TTL:   3600 * time.Second,
	}

	// Test successful call
	p := setupTest(nil, nil)
	ctx := context.Background()

	resultRecord, err := p.addDNSEntry(ctx, "example.com", testRecord)

	if err != nil {
		t.Errorf("Client.addDNSEntry() error = %v", err)
	}

	// Verify the returned record
	if resultRecord.Type != testRecord.Type ||
		resultRecord.Name != testRecord.Name ||
		resultRecord.Value != testRecord.Value ||
		resultRecord.ID != "12345" {
		t.Errorf("Client.addDNSEntry() record mismatch, got = %v, want Type=%s, Name=%s, Value=%s, ID=12345",
			resultRecord, testRecord.Type, testRecord.Name, testRecord.Value)
	}

	// Test error case
	p = setupTest(nil, errors.New("API error"))

	_, err = p.addDNSEntry(ctx, "example.com", testRecord)
	if err == nil {
		t.Error("Client.addDNSEntry() expected error, got nil")
	}
}

func TestClient_removeDNSEntry(t *testing.T) {
	// Test record to delete
	testRecord := libdns.Record{
		ID:    "1",
		Type:  "A",
		Name:  "test",
		Value: "192.168.1.1",
	}

	// Test successful call
	p := setupTest(nil, nil)
	ctx := context.Background()

	resultRecord, err := p.removeDNSEntry(ctx, "example.com", testRecord)

	if err != nil {
		t.Errorf("Client.removeDNSEntry() error = %v", err)
	}

	// Verify the ID was preserved
	if resultRecord.ID != testRecord.ID {
		t.Errorf("Client.removeDNSEntry() ID mismatch, got = %v, want = %v", resultRecord.ID, testRecord.ID)
	}

	// Test error case - API error
	p = setupTest(nil, errors.New("API error"))

	_, err = p.removeDNSEntry(ctx, "example.com", testRecord)
	if err == nil {
		t.Error("Client.removeDNSEntry() expected error, got nil")
	}

	// Test error case - invalid ID
	p = setupTest(nil, nil)
	invalidIDRecord := libdns.Record{
		ID:    "invalid", // Non-numeric ID
		Type:  "A",
		Name:  "test",
		Value: "192.168.1.1",
	}

	_, err = p.removeDNSEntry(ctx, "example.com", invalidIDRecord)
	if err == nil {
		t.Error("Client.removeDNSEntry() expected error for invalid ID, got nil")
	}
}

func TestClient_updateDNSEntry(t *testing.T) {
	// Test record to update
	testRecord := libdns.Record{
		ID:    "1",
		Type:  "A",
		Name:  "test",
		Value: "192.168.1.2", // Updated IP
		TTL:   7200 * time.Second,
	}

	// Test successful call
	p := setupTest(nil, nil)
	ctx := context.Background()

	resultRecord, err := p.updateDNSEntry(ctx, "example.com", testRecord)

	if err != nil {
		t.Errorf("Client.updateDNSEntry() error = %v", err)
	}

	// Verify the record was preserved
	if resultRecord.ID != testRecord.ID ||
		resultRecord.Type != testRecord.Type ||
		resultRecord.Name != testRecord.Name ||
		resultRecord.Value != testRecord.Value {
		t.Errorf("Client.updateDNSEntry() record mismatch, got = %v, want = %v", resultRecord, testRecord)
	}

	// Test error case - API error
	p = setupTest(nil, errors.New("API error"))

	_, err = p.updateDNSEntry(ctx, "example.com", testRecord)
	if err == nil {
		t.Error("Client.updateDNSEntry() expected error, got nil")
	}

	// Test error case - invalid ID
	p = setupTest(nil, nil)
	invalidIDRecord := libdns.Record{
		ID:    "invalid", // Non-numeric ID
		Type:  "A",
		Name:  "test",
		Value: "192.168.1.1",
	}

	_, err = p.updateDNSEntry(ctx, "example.com", invalidIDRecord)
	if err == nil {
		t.Error("Client.updateDNSEntry() expected error for invalid ID, got nil")
	}
}
