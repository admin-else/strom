package world

import (
	"math"
	"testing"

	"github.com/admin-else/strom/mc/proto_generated/v26_2"
)

func almostEqual(a, b float32) bool { return math.Abs(float64(a-b)) < 1e-6 }

// TestGameStateChangeWeather checks the ClientboundGameEventPacket weather
// mapping (start/stop raining and the smooth level changes).
func TestGameStateChangeWeather(t *testing.T) {
	w := NewWorld("26.4-snapshot-2", -64, 384)
	m := &Module{world: w}

	for _, tc := range []struct {
		reason string
		param  float32
		want   float32
	}{
		{"start_raining", 0, 0},
		{"stop_raining", 0, 1},
		{"rain_level_change", 0.42, 0.42},
	} {
		if err := m.onGameStateChange(&v26_2.PlayToClientPacketGameStateChange{Reason: tc.reason, GameMode: tc.param}); err != nil {
			t.Fatalf("%s: %v", tc.reason, err)
		}
		if got := w.RainLevel(); !almostEqual(got, tc.want) {
			t.Errorf("%s rain level = %v, want %v", tc.reason, got, tc.want)
		}
	}

	if err := m.onGameStateChange(&v26_2.PlayToClientPacketGameStateChange{Reason: "thunder_level_change", GameMode: 0.8}); err != nil {
		t.Fatal(err)
	}
	if got := w.ThunderLevel(); !almostEqual(got, 0.8*0.42) {
		t.Errorf("thunder level = %v, want %v (lerp * rain)", got, 0.8*0.42)
	}
}

// TestSetRainLevelClamps checks Level.setRainLevel clamps to [0, 1].
func TestSetRainLevelClamps(t *testing.T) {
	w := NewWorld("26.4-snapshot-2", -64, 384)
	w.SetRainLevel(2)
	if got := w.RainLevel(); got != 1 {
		t.Errorf("clamped rain = %v, want 1", got)
	}
	w.SetRainLevel(-1)
	if got := w.RainLevel(); got != 0 {
		t.Errorf("clamped rain = %v, want 0", got)
	}
}

// TestAdvanceWeatherSteps checks the +/-0.01 per-tick server projection.
func TestAdvanceWeatherSteps(t *testing.T) {
	w := NewWorld("26.4-snapshot-2", -64, 384)
	for range 5 {
		w.AdvanceWeather(true, true)
	}
	if got := w.RainLevel(); !almostEqual(got, 0.05) {
		t.Errorf("rain after 5 ticks = %v, want 0.05", got)
	}
	if got := w.ThunderLevel(); !almostEqual(got, 0.05*0.05) {
		t.Errorf("thunder after 5 ticks = %v, want 0.0025", got)
	}
	for range 100 {
		w.AdvanceWeather(true, true)
	}
	if got := w.RainLevel(); got != 1 {
		t.Errorf("rain clamped = %v, want 1", got)
	}
	for range 100 {
		w.AdvanceWeather(false, false)
	}
	if got := w.RainLevel(); got > 1e-4 {
		t.Errorf("rain clamped = %v, want ~0", got)
	}
}

// TestIsRainingThresholds checks the 0.2/0.9 thresholds.
func TestIsRainingThresholds(t *testing.T) {
	w := NewWorld("26.4-snapshot-2", -64, 384)
	w.SetRainLevel(0.2)
	if w.IsRaining() {
		t.Errorf("rain level exactly 0.2 must not be raining")
	}
	w.SetRainLevel(0.21)
	if !w.IsRaining() {
		t.Errorf("rain level 0.21 must be raining")
	}
	w.SetRainLevel(1)
	w.SetThunderLevel(0.9)
	if w.IsThundering() {
		t.Errorf("thunder level exactly 0.9 must not be thundering")
	}
	w.SetThunderLevel(1)
	if !w.IsThundering() {
		t.Errorf("thunder level 1 must be thundering")
	}
}
