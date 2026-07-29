# Strom API Reference

Strom is a Go Minecraft protocol library for building servers, bots, and proxies.

## Packages

### `mc/api` — Accounts & Authentication

```go
// Account holds a player's name, UUID, and access token.
type Account struct{ Name string; UUID uuid.UUID; AccessToken string }

// NewAccount creates a profile from name, UUID, and access token.
func NewAccount(name string, uuid uuid.UUID, ygg string) *Account

// NewAccountFromYGG parses an YGGDRASIL-style access token into an Account.
func NewAccountFromYGG(ygg string) *Account

// NewOfflineAccount creates an offline (cracked) account from a display name.
func NewOfflineAccount(name string) *Account

// HasJoined validates a session with Mojang's session server.
func HasJoined(name, server, ip string) (ProfileResp, error)

// Profile fetches a player profile by UUID (with or without signatures).
func Profile(id string, unsigned bool) (ProfileResp, error)

// NameToId / NameToIdBulk look up UUIDs from player names.
func NameToId(name string) (IdName, error)
func NameToIdBulk(names []string) ([]IdName, error)

// BlockedServers returns Mojang's blocklist.
func BlockedServers() ([]string, error)

// MojangPublicKeys returns Mojang's current signing keys.
func MojangPublicKeys() (MojangPublicKeysResp, error)
```

### `mc/client` — Client Connections

```go
// Connect dials a Minecraft server with SRV record lookup.
func Connect(ctx context.Context, addr string, version string) (*proto.Conn, error)

// ConnectVersionLess dials a server without version negotiation (used for status pings).
func ConnectVersionLess(ctx context.Context, addr string) (*proto.Conn, error)

// Login performs the full login handshake (online mode + encryption).
// Functional options: WithoutDns, WithContext, WithVersion.
func Login(connectTo string, account *api.Account, settings ...func(*loginSettings)) (*proto.Conn, error)

// LoginRaw logs in on an already-connected proto.Conn.
func LoginRaw(c *proto.Conn, account *api.Account) error

// IgnoreConfig discards all Configuration-state packets (use when bundling login → play).
func IgnoreConfig(c *proto.Conn) error

// MakeHandshakePacket builds a Handshake packet for the given next state.
func MakeHandshakePacket(c *proto.Conn, nextState proto_base.State) (*v1_21_8.HandshakingToServerPacketSetProtocol, error)
func MakeHandshakePacketAddr(s *proto.Conn, nextState proto_base.State, addr string) (*v1_21_8.HandshakingToServerPacketSetProtocol, error)

// Status pings a server and returns the parsed response.
func Status(ctx context.Context, addr string) (server.StatusResponse, error)

// StatusNoDns is like Status but skips SRV record resolution.
func StatusNoDns(addr string) (server.StatusResponse, error)

// StatusRaw returns a low-level StatusClient for custom status workflows.
func StatusRaw(ctx context.Context, addr string) (*StatusClient, error)
```

### `mc/client.StatusClient` — Programmatic Status Pings

```go
type StatusClient struct {
    *proto.Conn
    Status           string // raw JSON status response
    PingSendTime     time.Time
    PingReceiveTime  time.Time
    DoPingRoundTripTime bool
}

func (s *StatusClient) OnStatus(p *v1_21_8.StatusToClientPacketServerInfo) error
func (s *StatusClient) OnPong(p *v1_21_8.StatusToClientPacketPing) error
```

### `mc/server` — Server / Listen Side

```go
// StartServerWithOnConn listens on addr, accepts TCP connections, and calls onConn per connection.
func StartServerWithOnConn(listenAddr string, onConn func(c *proto.Conn) error) error

// Servee wraps a raw net.Conn into a proto.Conn with Servee actor.
func Servee(cNet net.Conn) *proto.Conn

// Kick sends a disconnect packet to the client.
func Kick(c *proto.Conn, reason *text.Component, status string) error

// ServeLogin handles the full login sequence (handshake → encryption → login → config → play).
// Functional options: WithOtherAccount, WithRawStatus, WithStatus,
//   WithCompatibleVersions, WithCompatibleVersionsRange, WithoutOnlineMode, WithoutEncryption.
func ServeLogin(c *proto.Conn, settings ...LoginServerSetting) (*LoginServer, error)

// ServeStatus handles a status-only connection (server list ping).
func ServeStatus(c *proto.Conn, s StatusResponse) error

// ServeConfig drives the configuration phase.
func ServeConfig(c *proto.Conn) error
```

### `mc/proto` — Connection & Packet Helpers

```go
// NewConn / NewConnCtx create a new proto.Conn with an embedded event loop.
func NewConn() *Conn
func NewConnCtx(ctx context.Context) *Conn

// SimplePacketToBytes encodes a typed packet to raw wire bytes (no length prefix).
func SimplePacketToBytes(packet proto_base.EncodeDecodeAble) ([]byte, error)

// SimplePacketToBytesLenPrefix encodes a typed packet with varint length prefix.
func SimplePacketToBytesLenPrefix(packet proto_base.EncodeDecodeAble) ([]byte, error)

// SimpleBytesToPacket decodes raw wire bytes into a typed packet.
func SimpleBytesToPacket(packetBytes []byte, version int32, direction proto_base.Direction, state proto_base.State) (proto_base.EncodeDecodeAble, error)

// Packet lookup by type/name.
func LookUpTypeByPacketInfo(direction proto_base.Direction, state proto_base.State, pid, version int32) (proto_base.PacketInfo, bool)
func LookupPacketInfoByType(packet proto_base.EncodeDecodeAble) (proto_base.PacketInfo, bool)
func LookupPacketInfoByNameProtocolVersionAndState(name string, version int32, state proto_base.State) (proto_base.PacketInfo, bool)
```

#### `proto.Conn` Methods

```go
// SetVersion loads the packet definitions for the given Minecraft version.
func (c *Conn) SetVersion(version string) error

// SetState transitions the connection to a new protocol state (Handshake, Status, Login, Config, Play).
func (c *Conn) SetState(s proto_base.State)

// State returns the current protocol state.
func (c *Conn) State() proto_base.State

// Send encodes and writes a typed packet for the current state.
func (c *Conn) Send(packet proto_base.EncodeDecodeAble) error

// Receive reads, decodes, and returns the next packet.
func (c *Conn) Receive() (proto_base.EncodeDecodeAble, error)

// StartConn begins the event loop — blocks until the connection closes.
func (c *Conn) StartConn() error

// Fire dispatches an event through the handler chain.
func (c *Conn) Fire(event any) error

// RegisterCritical registers a handler that closes the loop if it returns an error.
func (c *Conn) RegisterCritical(hs ...any)

// RegisterCriticalUntilLatest registers handlers for all versions between the packet version and the latest supported version.
// Auto-unregisters when SetVersion bumps the version.
func (c *Conn) RegisterCriticalUntilLatest(hs ...any)

// RegisterCriticalUntil registers handlers up to the specified version.
func (c *Conn) RegisterCriticalUntil(until string, hs ...any)

// HandlerFunctions maps packet types to registered handler reflect.Values.
// Useful for direct handler invocation in tests.
func (c *Conn) HandlerFunctions map[reflect.Type][]reflect.Value
```

### `mc/event` — Event Loop

```go
// NewLoop / NewLoopCtx create an event loop.
func NewLoop() *Loop
func NewLoopCtx(ctx context.Context) *Loop

// RegisterCritical registers a handler that closes the loop on error.
func (l *Loop) RegisterCritical(hs ...any)

// FireFound dispatches an event to matching handlers.
func (l *Loop) FireFound(event any) (found bool, err error)

// ValidateHandler checks that a function matches the handler signature (func(packet T) error).
func ValidateHandler(h any) (eventType reflect.Type, hv reflect.Value)
```

### `mc/proto_base` — Wire Types & Primitives

```go
// Packet types and metadata lookup.
type State int32       // Handshake, Status, Login, Config, Play
type Direction int     // ToServer, ToClient
type Actor int         // Server, Client
type PacketInfo struct { Direction Direction; State State; Version int32; PacketId int32 }

// Varint encoding.
func EncodeVarInt(w io.Writer, v int32) error
func DecodeVarInt(r io.Reader) (int32, error)
func EncodeVarLong(w io.Writer, v int64) error
func DecodeVarLong(r io.Reader) (int64, error)

// String encoding (varint-prefixed UTF-8).
func EncodeString(w io.Writer, s string) error
func DecodeString(r io.Reader) (string, error)

// Packet direction helpers.
func (a Actor) SendDirection() Direction
func (a Actor) ReceiveDirection() Direction
func (d Direction) Opposite() Direction
func (d Direction) String() string
```

### `mc/data` — Minecraft Version Data

```go
// Version lookup.
func LookUpVersionByName(name string) (VersionInfo, error)
func LookUpVersionByProtocolVersion(version int32) (VersionInfo, error)
func LookUpVersionByDataVersion(version int32) (VersionInfo, error)
func MustLookupProtocolVersion(version string) int32

// Block data.
func BlocksForVersion(v string) []*Block
func LookupBlockByStateId(version string, stateId int32) (*Block, bool)
func LookupBlockByName(version string, name string) (*Block, bool)
func StateIdFromBlocKAndStateMap(version string, name string, stateMap map[string]string) (int32, error)
func FromBlockState(version string, stateId int32) (*Block, map[string]string, error)

// Biome data.
func BiomesForVersion(v string) []*Biome
func LookupBiomeById(version string, id int32) (*Biome, bool)
func LookupBiomeByName(version string, name string) (*Biome, bool)

// Item data.
func ItemsForVersion(v string) []*Item
func LookupItemById(version string, id int32) (*Item, bool)
func LookupItemByName(version string, name string) (*Item, bool)

// Noise / generation data.
func LoadMultiNoise(version string) (*MultiNoiseData, error)

// Registry data.
func LoadRegistry(name string) (*nbt2.Tag, error)

// Generic JSON loading (relative to minecraft-data directory).
func LoadJson(path string, data any) error
func LoadVersionedJson(version, dataName string, data any) error
```

### `mc/crypto` — Encryption Helpers

```go
func AuthDigest(elems ...[]byte) string
func NewCFB8Decrypt(c cipher.Block, iv []byte) *CFB8
func NewCFB8Encrypt(c cipher.Block, iv []byte) *CFB8
func FromOfflinePlayer(displayName string) uuid.UUID
```

### `mc/text` — Text Components

```go
// Pretty wraps a plain string in a Component.
func Pretty(s string) *Component

// PrettyF wraps a formatted string in a Component.
func PrettyF(format string, args ...any) *Component

// PrettyPlaceholders replaces placeholders in a translation string.
func PrettyPlaceholders(format string, placeholders map[string]any) *Component
```

### `mc/nbt` — Named Binary Tag I/O

```go
func ReadFile(r io.Reader, data any) error
func ReadUnstructuredFile(r io.Reader) (*Tag, error)
func ReadUnstructuredFilePath(path string) (*Tag, error)
func WriteUnstructuredFile(w io.Writer, n *Tag) error
func WriteUnstructuredFilePath(path string, n *Tag) error
func PrintSNBTAny(t any, w io.Writer) error
```

### `mc/level` — World & Chunk Data

```go
// Chunk reading.
func ReadChunkFromChunkPacketData(r io.Reader, version string, worldHeight int) (*Chunk, error)

// Section / storage format helpers.
func MakeBiomeFormat(version string) StorageFormat
func MakeBlockFormat(version string) StorageFormat
func ReadSectionStorage(r io.Reader, format StorageFormat) (*Storage, error)
func WriteSectionStorage(w io.Writer, s *Storage) error
func SectionDecodePacket(r io.Reader, version string) (Section, error)
func SectionEncodePacket(w io.Writer, s Section) error
func SectionEquals(a, b Section) bool
func CountStorage(s *Storage, v int32) int
```

#### `mc/level/anvil` — Anvil Region Format

```go
func ReadLevelData(r io.Reader) (Level, error)
func ReadChunkHeader(r *os.File) (*ChunkHeader, error)
```

#### `mc/level/structure` / `mc/level/litematic`

```go
func Load(fileName string) (*Structure, error)              // structure
func LoadFromPath(path string) (*Structure, error)           // litematic
func LoadFromFile(r io.Reader) (*Structure, error)           // litematic
```

### `mc/gen` — World Generation

```go
func XoroshiroFromSeed(seed int64) Xoroshiro
func XoroshiroFromString(s string) Xoroshiro
func PerlinNoiseFromXoroshiro(x Xoroshiro) PerlinNoise
func OctaveFromXoroshiro(x *Xoroshiro, amplitudes []float64, omin, lenn, nmax int) Octave
func DoublePerlinFromXoroshiro(x *Xoroshiro, amplitudes []float64, omin, lenn, nmax int) DoublePerlin
func BiomeNoiseFromXoroshiro(x Xoroshiro, large bool) BiomeNoise
func CreateBiomeDepthSpline() Spline
```

### `mc/mapstructure` — Tag-Based Decoder

```go
// Decode populates a struct from a map using struct tags.
func Decode(data any, v any) error

// NewFormat creates a custom format with options.
func NewFormat(name string, opts ...Option) *Format

// Options.
func WithRequireAll() Option
func WithErrorOnExtra() Option
func WithTrySnakeCase() Option
func WithTryLowCase() Option
func WithTypeCodec(decodeFn, encodeFn any) Option
```

### `mc/util` — Miscellaneous Helpers

```go
func CamelCase(s string) string
func SnakeCase(s string) string
func FirstLetterLower(s string) string
func FloorDiv(a, b int32) int32
func BpeByNum(x float64) uint8
func Must(err error)
```

## Handler Signature

All event handlers must match:
```go
func(typedPacket) error
```
The argument must implement `proto_base.EncodeDecodeAble` (has `Decode`/`Encode` methods). Return `event.HandlerDoneErr{}` to stop the loop cleanly, or any other error to close the connection.

## Quick Recipes

### Server list ping
```go
// Programmatic
resp, err := client.Status(context.Background(), "mc.hypixel.net:25565")

// With raw JSON
s, err := client.StatusRaw(ctx, addr)
_ = json.Unmarshal([]byte(s.Status), &resp)
```

### Start a server
```go
server.StartServerWithOnConn(":25565", func(c *proto.Conn) error {
    c.RegisterCriticalUntilLatest(handleHandshake, handleLogin)
    return c.StartConn()
})
```

### Construct a packet without a connection
```go
// Encode to wire format
b, err := proto.SimplePacketToBytesLenPrefix(&v1_21_8.StatusToServerPacketPing{Time: 1234})

// Decode from wire format
packet, err := proto.SimpleBytesToPacket(bytes, 772, proto_base.ToClient, proto_base.Status)
info, ok := packet.(*v1_21_8.StatusToClientPacketServerInfo)
```

### Look up version data
```go
v, _ := data.LookUpVersionByName("1.21.8")
block, _ := data.LookupBlockByName("1.21.8", "minecraft:stone")
protocol := data.MustLookupProtocolVersion("1.21.8")
```

### Read NBT
```go
tag, _ := nbt.ReadUnstructuredFile(file)
data.Level level
nbt.ReadFile(file, &data)
```
