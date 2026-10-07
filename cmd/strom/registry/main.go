// Command registry connects to a server, captures the configuration-phase
// dynamic registries and prints the id mappings used by neon.
package registry

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/admin-else/strom/cmd/strom/cmd_util"
	"github.com/admin-else/strom/mc/client"
	stromregistry "github.com/admin-else/strom/mc/registry"
)

var (
	cmd         = flag.NewFlagSet("registry", flag.ContinueOnError)
	addrFlag    = cmd.String("addr", "127.0.0.1:25565", "address to connect to")
	accFlag     = cmd.String("acc", "dev", "account name (offline mode)")
	versionFlag = cmd.String("version", "26.4-snapshot-2", "version to use")
	jsonFlag    = cmd.Bool("json", false, "dump every registry as JSON")
	entryFlag   = cmd.String("registry", "", "print one registry's entries in id order")
	noAssert    = cmd.Bool("no-assert", false, "skip the built-in assertion checks")
)

const (
	biomeRegistry    = "minecraft:worldgen/biome"
	entityRegistries = "minecraft:entity_type"
)

// expectedBiomeIDs pins the 26.4 biome ordering the neon client currently
// hardcodes in net/minecraft/client/color/biome/registry_order.json.
var expectedBiomeIDs = map[string]int{
	"minecraft:dappled_forest": 8,
	"minecraft:sulfur_caves":   54,
}

func dumpJSON(store *stromregistry.Store) (err error) {
	out := map[string][]string{}
	for _, name := range store.Names() {
		r, _ := store.Registry(name)
		out[name] = r.Entries()
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func printRegistry(store *stromregistry.Store, name string) (err error) {
	r, ok := store.Registry(name)
	if !ok {
		return fmt.Errorf("registry %q not captured", name)
	}
	fmt.Printf("%s (%d entries)\n", name, r.Len())
	for id, resourceID := range r.Entries() {
		fmt.Printf("  %3d  %s\n", id, resourceID)
	}
	return
}

func Run(args []string) (err error) {
	err = cmd.Parse(args)
	if err != nil {
		return
	}

	acc, err := cmd_util.Account(*accFlag)
	if err != nil {
		return
	}

	store := stromregistry.NewStore()
	c, err := client.Login(*addrFlag, acc,
		client.WithVersion(*versionFlag),
		client.WithRegistries(store),
	)
	if err != nil {
		return
	}
	defer c.Close()

	if *jsonFlag {
		return dumpJSON(store)
	}

	if *entryFlag != "" {
		return printRegistry(store, *entryFlag)
	}

	names := store.Names()
	fmt.Printf("captured %d registries\n", len(names))
	for _, name := range names {
		r, _ := store.Registry(name)
		fmt.Printf("  %-55s %d\n", name, r.Len())
	}

	biomes, ok := store.Registry(biomeRegistry)
	if !ok {
		return fmt.Errorf("biome registry %q not captured", biomeRegistry)
	}
	fmt.Printf("\n%s: %d entries\n", biomeRegistry, biomes.Len())
	for name, want := range expectedBiomeIDs {
		got, found := store.RegistryID(biomeRegistry, name)
		status := "OK"
		if !found || got != want {
			status = "MISMATCH"
		}
		fmt.Printf("  %s = %d (want %d) %s\n", name, got, want, status)
	}

	if entities, ok := store.Registry(entityRegistries); ok {
		fmt.Printf("\n%s: %d entries\n", entityRegistries, entities.Len())
		for _, name := range []string{"minecraft:cow", "minecraft:zombie", "minecraft:player"} {
			id, found := store.RegistryID(entityRegistries, name)
			fmt.Printf("  %s = %d (found %v)\n", name, id, found)
		}
	} else {
		fmt.Printf("\n%s: NOT sent by the server\n", entityRegistries)
	}

	if *noAssert {
		return
	}

	failures := 0
	for name, want := range expectedBiomeIDs {
		got, found := store.RegistryID(biomeRegistry, name)
		if !found || got != want {
			fmt.Fprintf(os.Stderr, "assertion failed: biome %s = %d (found %v), want %d\n", name, got, found, want)
			failures++
		}
	}
	sort.Strings(names)
	if failures > 0 {
		return fmt.Errorf("%d registry assertions failed", failures)
	}
	fmt.Println("\nall registry assertions passed")
	return
}
