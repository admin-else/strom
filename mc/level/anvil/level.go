package anvil

import (
	"io"

	"github.com/admin-else/strom/mc/nbt"
)

type Level struct {
	ClearWeatherTime           int32 `nbt:"clearWeatherTime"`
	DataPacks                  DataPacks
	DayTime                    int64
	WanderingTraderSpawnChance int32
	GameType                   int32
	WorldGenSettings           WorldGenSettings
	ScheduledEvents            []interface{}
	WasModded                  bool
	ThunderTime                int32 `nbt:"thunderTime"`
	Difficulty                 int8
	DragonFight                DragonFight
	Initialized                bool
	AllowCommands              bool `nbt:"allowCommands"`
	Thundering                 bool
	LastPlayed                 int64
	GameRules                  map[string]interface{}
	RainTime                   int32 `nbt:"rainTime"`
	DataVersion                int32
	Version                    Version
	WanderingTraderSpawnDelay  int32
	Hardcore                   bool
	Spawn                      Spawn
	Raining                    bool
	Player                     Player
	DifficultyLocked           bool
	LevelName                  string
	Time                       int64
	ServerBrands               []string
	CustomBossEvents           map[string]interface{}
}

type DataPacks struct {
	Enabled  []string
	Disabled []string
}

type WorldGenSettings struct {
	BonusChest       bool
	Seed             int64
	GenerateFeatures bool
	Dimensions       map[string]Dimension
}

type Dimensions struct {
	Overworld Dimension
	Nether    Dimension
	End       Dimension
}

type Dimension struct {
	Generator Generator
	Type      string
}

type Generator struct {
	Type        string
	Settings    string
	BiomeSource BiomeSource
}

type BiomeSource struct {
	Preset string `nbt:"preset,omitempty"`
	Type   string
}

type DragonFight struct {
	NeedsStateScanning bool
	Gateways           []int32
	DragonKilled       bool
	PreviouslyKilled   bool
}

type Version struct {
	Snapshot bool
	Series   string
	Id       int32
	Name     string
}

type Spawn struct {
	Dimension string
	Yaw       float32
	Pos       [3]int32
	Pitch     float32
}

type Player struct {
	XpLevel                              int32
	HurtTime                             int16
	DeathTime                            int16
	XpTotal                              int32
	FoodExhaustionLevel                  float32
	RecipeBook                           RecipeBook
	Fire                                 int16
	Brain                                Brain
	XpSeed                               int32
	XpP                                  float32
	EnderItems                           []interface{}
	Invulnerable                         bool
	Pos                                  [3]float64
	Inventory                            []InventoryItem
	SleepTimer                           int16
	Abilities                            Abilities
	PlayerGameType                       int32
	SeenCredits                          bool
	FoodSaturationLevel                  float32
	FallDistance                         float64
	OnGround                             bool
	Dimension                            string
	PortalCooldown                       int32
	Motion                               [3]float64
	IgnoreFallDamageFromCurrentExplosion bool
	Rotation                             [2]float32
	CurrentImpulseContextResetGraceTime  int32
	WardenSpawnTracker                   WardenSpawnTracker
	PreviousPlayerGameType               int32
	Attributes                           []Attribute
	HurtByTimestamp                      int32
	UUID                                 [4]int32
	Score                                int32
	DataVersion                          int32
	FoodLevel                            int32
	SpawnExtraParticlesOnFall            bool
	SelectedItemSlot                     int32
	FoodTickTimer                        int32
	FallFlying                           bool
	AbsorptionAmount                     float32
	Health                               float32
	Air                                  int16
}

type RecipeBook struct {
	Recipes       []string
	ToBeDisplayed []string
}

type Brain struct {
	Memories map[string]interface{}
}

type InventoryItem struct {
	Count int32
	Slot  int8
	Id    string
}

type Abilities struct {
	Invulnerable bool
	Mayfly       bool
	Instabuild   bool
	WalkSpeed    float32
	MayBuild     bool
	Flying       bool
	FlySpeed     float32
}

type WardenSpawnTracker struct {
	WarningLevel          int32
	TicksSinceLastWarning int32
	CooldownTicks         int32
}

type Attribute struct {
	Id   string
	Base float64
}

func ReadLevelData(r io.Reader) (l Level, err error) {
	var lDataWrapped struct{ Data Level }
	err = nbt.ReadFile(r, &lDataWrapped)
	if err != nil {
		return
	}
	l = lDataWrapped.Data
	return
}
