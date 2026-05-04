package eddlogger

import (
	"sync"
	"testing"
)

func TestRouteCollector_AddSave(t *testing.T) {
	logger := NewLogger("test-service")
	collector := logger.NewSDKTrnCollector("req-123")

	collector.Add(&SDKTrnOptions{
		SKU:     "SKU001",
		Quantity: 2,
		CP:      "01000",
		Route:   "RUTA-A",
		StoreID: "STORE-1",
		Score:   0.95,
		IsRecalculated: true,
	})
	collector.Add(&SDKTrnOptions{
		SKU:     "SKU002",
		Quantity: 1,
		CP:      "01000",
		Route:   "RUTA-B",
		StoreID: "STORE-2",
		Score:   0.87,
		IsRecalculated: false,
	})
	collector.Add(&SDKTrnOptions{
		SKU:     "SKU003",
		Quantity: 3,
		CP:      "02000",
		Route:   "RUTA-C",
		StoreID: "STORE-1",
		Score:   0.92,
		IsRecalculated: true,
	})

	if collector.Len() != 3 {
		t.Fatalf("expected 3 routes, got %d", collector.Len())
	}

	n, err := collector.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3 saved, got %d", n)
	}

	n2, _ := collector.Save()
	if n2 != 0 {
		t.Fatalf("second Save should return 0, got %d", n2)
	}
}

func TestRouteCollector_NilOpts(t *testing.T) {
	logger := NewLogger("test")
	collector := logger.NewSDKTrnCollector("req-456")
	collector.Add(nil)
	if collector.Len() != 0 {
		t.Fatalf("expected 0 routes with nil opts, got %d", collector.Len())
	}
}

func TestRouteCollector_Concurrent(t *testing.T) {
	logger := NewLogger("test-concurrent")
	collector := logger.NewSDKTrnCollector("req-concurrent")

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			collector.Add(&SDKTrnOptions{
				SKU:      "SKU",
				Quantity: idx,
				CP:       "01000",
			})
		}(i)
	}
	wg.Wait()

	if collector.Len() != 100 {
		t.Fatalf("expected 100 routes, got %d", collector.Len())
	}

	n, _ := collector.Save()
	if n != 100 {
		t.Fatalf("expected 100 saved, got %d", n)
	}
}

func TestRouteCollector_WinnerAlwaysTrue(t *testing.T) {
	logger := NewLogger("test")
	collector := logger.NewSDKTrnCollector("req-winner")

	collector.Add(&SDKTrnOptions{
		SKU:      "SKU001",
		Quantity: 1,
		CP:       "01000",
	})

	n, err := collector.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 saved, got %d", n)
	}
}

func TestRouteCollector_DefaultServiceName(t *testing.T) {
	logger := NewLogger("my-service")
	collector := logger.NewSDKTrnCollector("req-svc")

	collector.Add(&SDKTrnOptions{
		SKU:      "SKU001",
		Quantity: 1,
		CP:       "01000",
	})

	n, err := collector.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 saved, got %d", n)
	}
}

func TestRouteCollector_CustomServiceName(t *testing.T) {
	logger := NewLogger("default-svc")
	collector := logger.NewSDKTrnCollector("req-custom")

	collector.Add(&SDKTrnOptions{
		SKU:         "SKU001",
		Quantity:    1,
		CP:          "01000",
		ServiceName: "custom-override",
	})

	n, err := collector.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 saved, got %d", n)
	}
}
