package data

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

func ItemsForVersion(v string) (ret []*Item) {
	var ok bool
	if ret, ok = ItemsCache[v]; ok {
		return
	}
	var b []*Item
	must(LoadVersionedJson(v, "items", &b))
	ItemsCache[v] = b
	return b
}

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

func LookupItemByName(version string, name string) (item *Item, ok bool) {
	for _, i := range ItemsForVersion(version) {
		if i.Name == name {
			item = i
			ok = true
		}
	}
	return
}
