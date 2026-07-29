package data

import "sync"

//   {
//    "id": 2,
//    "name": "granite",
//    "displayName": "Granite",
//    "stackSize": 64
//  },

type Item struct {
	Id          int32
	Name        string
	DisplayName string
	StackSize   int32
}

var ItemsCache = make(map[string][]*Item)
var ItemsCacheMutex = &sync.RWMutex{}

// ItemsForVersion returns the item definitions for the given Minecraft version, caching them after the first load.
func ItemsForVersion(v string) (ret []*Item) {
	var ok bool
	ItemsCacheMutex.RLock()
	if ret, ok = ItemsCache[v]; ok {
		ItemsCacheMutex.RUnlock()
		return
	}
	ItemsCacheMutex.RUnlock()
	var b []*Item
	must(LoadVersionedJson(v, "items", &b))
	ItemsCacheMutex.Lock()
	ItemsCache[v] = b
	ItemsCacheMutex.Unlock()
	return b
}

// LookupItemById returns the item with the given numeric ID in the given version.
func LookupItemById(version string, id int32) (item *Item, ok bool) {
	for _, i := range ItemsForVersion(version) {
		if i.Id == id {
			item = i
			ok = true
			return
		}
	}
	return
}

// LookupItemByName returns the item with the given name in the given version.
func LookupItemByName(version string, name string) (item *Item, ok bool) {
	for _, i := range ItemsForVersion(version) {
		if i.Name == name {
			item = i
			ok = true
		}
	}
	return
}
