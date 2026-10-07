package world

import (
	"github.com/admin-else/strom/mc/proto_generated/v26_2"
	"github.com/admin-else/strom/mc/util"
)

// Weather constants mirror net.minecraft.world.level.Level. The rain and
// thunder levels clamp to [0, 1] and the server advances them by 0.01 per tick.
const weatherLevelStep = 0.01

// onGameStateChange applies the weather events of ClientboundGameEventPacket.
// 26.4 carries start/stop raining (reason 1/2) and the smooth rain/thunder
// level changes (reason 7/8, the packet's float parameter).
func (m *Module) onGameStateChange(p *v26_2.PlayToClientPacketGameStateChange) (err error) {
	switch p.Reason {
	case "start_raining":
		m.world.SetRainLevel(0)
	case "stop_raining":
		m.world.SetRainLevel(1)
	case "rain_level_change":
		m.world.SetRainLevel(p.GameMode)
	case "thunder_level_change":
		m.world.SetThunderLevel(p.GameMode)
	}
	return nil
}

// SetRainLevel mirrors Level.setRainLevel: the value is clamped to [0, 1] and
// the previous level is set equal so getRainLevel interpolates from the new
// value.
func (w *World) SetRainLevel(rainLevel float32) {
	clampedRainLevel := util.ClampFloat(rainLevel, 0, 1)
	w.mu.Lock()
	defer w.mu.Unlock()
	w.oRainLevel = clampedRainLevel
	w.rainLevel = clampedRainLevel
}

// SetThunderLevel mirrors Level.setThunderLevel: the value is clamped to [0, 1]
// and the previous level is set equal.
func (w *World) SetThunderLevel(thunderLevel float32) {
	clampedThunderLevel := util.ClampFloat(thunderLevel, 0, 1)
	w.mu.Lock()
	defer w.mu.Unlock()
	w.oThunderLevel = clampedThunderLevel
	w.thunderLevel = clampedThunderLevel
}

// GetRainLevel mirrors Level.getRainLevel, interpolating between the previous
// and current level by the partial tick a.
func (w *World) GetRainLevel(a float32) (ret float32) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return util.LerpFloat(a, w.oRainLevel, w.rainLevel)
}

// GetThunderLevel mirrors Level.getThunderLevel: the interpolated thunder level
// scaled by the rain level.
func (w *World) GetThunderLevel(a float32) (ret float32) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return util.LerpFloat(a, w.oThunderLevel, w.thunderLevel) * util.LerpFloat(a, w.oRainLevel, w.rainLevel)
}

// RainLevel returns the current rain level, i.e. getRainLevel(1.0F).
func (w *World) RainLevel() (ret float32) { return w.GetRainLevel(1) }

// ThunderLevel returns the current thunder level, i.e. getThunderLevel(1.0F).
func (w *World) ThunderLevel() (ret float32) { return w.GetThunderLevel(1) }

// IsRaining mirrors Level.isRaining: the rain level is above 0.2.
func (w *World) IsRaining() (ret bool) { return w.RainLevel() > 0.2 }

// IsThundering mirrors Level.isThundering: the thunder level is above 0.9.
func (w *World) IsThundering() (ret bool) { return w.ThunderLevel() > 0.9 }

// AdvanceWeather mirrors the level portion of ServerLevel.advanceWeatherCycle:
// the rain and thunder levels move by 0.01 per tick toward their target and are
// clamped to [0, 1]. The previous levels are saved first so GetRainLevel/
// GetThunderLevel interpolate within the tick.
func (w *World) AdvanceWeather(raining, thundering bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.oThunderLevel = w.thunderLevel
	if thundering {
		w.thunderLevel += weatherLevelStep
	} else {
		w.thunderLevel -= weatherLevelStep
	}
	w.thunderLevel = util.ClampFloat(w.thunderLevel, 0, 1)
	w.oRainLevel = w.rainLevel
	if raining {
		w.rainLevel += weatherLevelStep
	} else {
		w.rainLevel -= weatherLevelStep
	}
	w.rainLevel = util.ClampFloat(w.rainLevel, 0, 1)
}
