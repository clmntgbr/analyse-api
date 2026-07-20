package video

import "testing"

func TestSampleTimestampsCountAndBounds(t *testing.T) {
	timestamps := sampleTimestamps(60, 10)
	if len(timestamps) != 10 {
		t.Fatalf("expected 10 timestamps, got %d", len(timestamps))
	}

	for i, ts := range timestamps {
		if ts < 0 || ts >= 60 {
			t.Fatalf("timestamp %d out of bounds: %f", i, ts)
		}
		if i > 0 && ts <= timestamps[i-1] {
			t.Fatalf("timestamps must be increasing")
		}
	}
}

func TestSampleTimestampsShortVideo(t *testing.T) {
	timestamps := sampleTimestamps(2, 10)
	if len(timestamps) != 10 {
		t.Fatalf("expected 10 timestamps, got %d", len(timestamps))
	}
}
