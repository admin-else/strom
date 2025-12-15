package v1_21_8

import (
	"encoding/binary"
	"io"

	"github.com/admin-else/queser"
	"github.com/admin-else/queser/nbt"
	"github.com/google/uuid"
)

type ArmorTrimMaterial struct {
	AssetBase           string
	OverrideArmorAssets []struct {
		Key   string
		Value string
	}
	Description nbt.Anon
}

func (_ ArmorTrimMaterial) Decode(r io.Reader) (ret ArmorTrimMaterial, err error) {
	ret.AssetBase, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var lArmorTrimMaterialOverrideArmorAssets queser.VarInt
	lArmorTrimMaterialOverrideArmorAssets, err = lArmorTrimMaterialOverrideArmorAssets.Decode(r)
	if err != nil {
		return
	}
	ret.OverrideArmorAssets = []struct {
		Key   string
		Value string
	}{}
	for range lArmorTrimMaterialOverrideArmorAssets {
		var ArmorTrimMaterialOverrideArmorAssetsElement struct {
			Key   string
			Value string
		}
		ArmorTrimMaterialOverrideArmorAssetsElement.Key, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ArmorTrimMaterialOverrideArmorAssetsElement.Value, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.OverrideArmorAssets = append(ret.OverrideArmorAssets, ArmorTrimMaterialOverrideArmorAssetsElement)
	}
	ret.Description, err = ret.Description.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret ArmorTrimMaterial) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.AssetBase)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.OverrideArmorAssets)).Encode(w)
	if err != nil {
		return
	}
	for iArmorTrimMaterialOverrideArmorAssets := range len(ret.OverrideArmorAssets) {
		err = queser.EncodeString(w, ret.OverrideArmorAssets[iArmorTrimMaterialOverrideArmorAssets].Key)
		if err != nil {
			return
		}
		err = queser.EncodeString(w, ret.OverrideArmorAssets[iArmorTrimMaterialOverrideArmorAssets].Value)
		if err != nil {
			return
		}
	}
	err = ret.Description.Encode(w)
	if err != nil {
		return
	}
	return
}

type ArmorTrimPattern struct {
	AssetId     string
	Description nbt.Anon
	Decal       bool
}

func (_ ArmorTrimPattern) Decode(r io.Reader) (ret ArmorTrimPattern, err error) {
	ret.AssetId, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Description, err = ret.Description.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Decal)
	if err != nil {
		return
	}
	return
}
func (ret ArmorTrimPattern) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.AssetId)
	if err != nil {
		return
	}
	err = ret.Description.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Decal)
	if err != nil {
		return
	}
	return
}

type BannerPattern struct {
	AssetId        string
	TranslationKey string
}

func (_ BannerPattern) Decode(r io.Reader) (ret BannerPattern, err error) {
	ret.AssetId, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.TranslationKey, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	return
}
func (ret BannerPattern) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.AssetId)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.TranslationKey)
	if err != nil {
		return
	}
	return
}

type BannerPatternLayer struct {
	Pattern any
	ColorId queser.VarInt
}

func (_ BannerPatternLayer) Decode(r io.Reader) (ret BannerPatternLayer, err error) {
	var BannerPatternLayerPatternId queser.VarInt
	BannerPatternLayerPatternId, err = BannerPatternLayerPatternId.Decode(r)
	if err != nil {
		return
	}
	if BannerPatternLayerPatternId != 0 {
		ret.Pattern = BannerPatternLayerPatternId
		return
	}
	var BannerPatternLayerPatternResult BannerPattern
	BannerPatternLayerPatternResult, err = BannerPatternLayerPatternResult.Decode(r)
	if err != nil {
		return
	}
	ret.Pattern = BannerPatternLayerPatternResult
	ret.ColorId, err = ret.ColorId.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret BannerPatternLayer) Encode(w io.Writer) (err error) {
	switch BannerPatternLayerPatternKnownType := ret.Pattern.(type) {
	case queser.VarInt:
		err = BannerPatternLayerPatternKnownType.Encode(w)
		if err != nil {
			return
		}
	case BannerPattern:
		err = BannerPatternLayerPatternKnownType.Encode(w)
		if err != nil {
			return
		}
	default:
		err = queser.BadTypeError
	}
	err = ret.ColorId.Encode(w)
	if err != nil {
		return
	}
	return
}

type ByteArray struct {
	Val []byte
}

func (_ ByteArray) Decode(r io.Reader) (ret ByteArray, err error) {
	var lByteArray queser.VarInt
	lByteArray, err = lByteArray.Decode(r)
	if err != nil {
		return
	}
	ret.Val, err = io.ReadAll(io.LimitReader(r, int64(lByteArray)))
	if err != nil {
		return
	}
	return
}
func (ret ByteArray) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Val)).Encode(w)
	if err != nil {
		return
	}
	_, err = w.Write(ret.Val)
	if err != nil {
		return
	}
	return
}

type ContainerID struct {
	Val queser.VarInt
}

func (_ ContainerID) Decode(r io.Reader) (ret ContainerID, err error) {
	ret.Val, err = ret.Val.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret ContainerID) Encode(w io.Writer) (err error) {
	err = ret.Val.Encode(w)
	if err != nil {
		return
	}
	return
}

type DataComponentMatchers struct {
	ExactMatchers   ExactComponentMatcher
	PartialMatchers []queser.VarInt
}

func (_ DataComponentMatchers) Decode(r io.Reader) (ret DataComponentMatchers, err error) {
	ret.ExactMatchers, err = ret.ExactMatchers.Decode(r)
	if err != nil {
		return
	}
	var lDataComponentMatchersPartialMatchers queser.VarInt
	lDataComponentMatchersPartialMatchers, err = lDataComponentMatchersPartialMatchers.Decode(r)
	if err != nil {
		return
	}
	ret.PartialMatchers = []queser.VarInt{}
	for range lDataComponentMatchersPartialMatchers {
		var DataComponentMatchersPartialMatchersElement queser.VarInt
		DataComponentMatchersPartialMatchersElement, err = DataComponentMatchersPartialMatchersElement.Decode(r)
		if err != nil {
			return
		}
		ret.PartialMatchers = append(ret.PartialMatchers, DataComponentMatchersPartialMatchersElement)
	}
	return
}
func (ret DataComponentMatchers) Encode(w io.Writer) (err error) {
	err = ret.ExactMatchers.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.PartialMatchers)).Encode(w)
	if err != nil {
		return
	}
	for iDataComponentMatchersPartialMatchers := range len(ret.PartialMatchers) {
		err = ret.PartialMatchers[iDataComponentMatchersPartialMatchers].Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type EntityMetadataPaintingVariant struct {
	Width   int32
	Height  int32
	AssetId string
	Title   *nbt.Anon
	Author  *nbt.Anon
}

func (_ EntityMetadataPaintingVariant) Decode(r io.Reader) (ret EntityMetadataPaintingVariant, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Width)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Height)
	if err != nil {
		return
	}
	ret.AssetId, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var EntityMetadataPaintingVariantTitlePresent bool
	err = binary.Read(r, binary.BigEndian, &EntityMetadataPaintingVariantTitlePresent)
	if err != nil {
		return
	}
	if EntityMetadataPaintingVariantTitlePresent {
		var EntityMetadataPaintingVariantTitlePresentValue nbt.Anon
		EntityMetadataPaintingVariantTitlePresentValue, err = EntityMetadataPaintingVariantTitlePresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.Title = &EntityMetadataPaintingVariantTitlePresentValue
	}
	var EntityMetadataPaintingVariantAuthorPresent bool
	err = binary.Read(r, binary.BigEndian, &EntityMetadataPaintingVariantAuthorPresent)
	if err != nil {
		return
	}
	if EntityMetadataPaintingVariantAuthorPresent {
		var EntityMetadataPaintingVariantAuthorPresentValue nbt.Anon
		EntityMetadataPaintingVariantAuthorPresentValue, err = EntityMetadataPaintingVariantAuthorPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.Author = &EntityMetadataPaintingVariantAuthorPresentValue
	}
	return
}
func (ret EntityMetadataPaintingVariant) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Width)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Height)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.AssetId)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Title != nil)
	if err != nil {
		return
	}
	if ret.Title != nil {
		err = (*ret.Title).Encode(w)
		if err != nil {
			return
		}
	}
	err = binary.Write(w, binary.BigEndian, ret.Author != nil)
	if err != nil {
		return
	}
	if ret.Author != nil {
		err = (*ret.Author).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type ExactComponentMatcher struct {
	Val []SlotComponent
}

func (_ ExactComponentMatcher) Decode(r io.Reader) (ret ExactComponentMatcher, err error) {
	var lExactComponentMatcher queser.VarInt
	lExactComponentMatcher, err = lExactComponentMatcher.Decode(r)
	if err != nil {
		return
	}
	ret.Val = []SlotComponent{}
	for range lExactComponentMatcher {
		var ExactComponentMatcherElement SlotComponent
		ExactComponentMatcherElement, err = ExactComponentMatcherElement.Decode(r)
		if err != nil {
			return
		}
		ret.Val = append(ret.Val, ExactComponentMatcherElement)
	}
	return
}
func (ret ExactComponentMatcher) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Val)).Encode(w)
	if err != nil {
		return
	}
	for iExactComponentMatcher := range len(ret.Val) {
		err = ret.Val[iExactComponentMatcher].Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type HashedSlot struct {
	ItemId     queser.VarInt
	ItemCount  queser.VarInt
	Components []struct {
		Type SlotComponentType
		Hash int32
	}
	RemoveComponents []struct {
		Type SlotComponentType
	}
}

func (_ HashedSlot) Decode(r io.Reader) (ret HashedSlot, err error) {
	ret.ItemId, err = ret.ItemId.Decode(r)
	if err != nil {
		return
	}
	ret.ItemCount, err = ret.ItemCount.Decode(r)
	if err != nil {
		return
	}
	var lHashedSlotComponents queser.VarInt
	lHashedSlotComponents, err = lHashedSlotComponents.Decode(r)
	if err != nil {
		return
	}
	ret.Components = []struct {
		Type SlotComponentType
		Hash int32
	}{}
	for range lHashedSlotComponents {
		var HashedSlotComponentsElement struct {
			Type SlotComponentType
			Hash int32
		}
		HashedSlotComponentsElement.Type, err = HashedSlotComponentsElement.Type.Decode(r)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &HashedSlotComponentsElement.Hash)
		if err != nil {
			return
		}
		ret.Components = append(ret.Components, HashedSlotComponentsElement)
	}
	var lHashedSlotRemoveComponents queser.VarInt
	lHashedSlotRemoveComponents, err = lHashedSlotRemoveComponents.Decode(r)
	if err != nil {
		return
	}
	ret.RemoveComponents = []struct {
		Type SlotComponentType
	}{}
	for range lHashedSlotRemoveComponents {
		var HashedSlotRemoveComponentsElement struct {
			Type SlotComponentType
		}
		HashedSlotRemoveComponentsElement.Type, err = HashedSlotRemoveComponentsElement.Type.Decode(r)
		if err != nil {
			return
		}
		ret.RemoveComponents = append(ret.RemoveComponents, HashedSlotRemoveComponentsElement)
	}
	return
}
func (ret HashedSlot) Encode(w io.Writer) (err error) {
	err = ret.ItemId.Encode(w)
	if err != nil {
		return
	}
	err = ret.ItemCount.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Components)).Encode(w)
	if err != nil {
		return
	}
	for iHashedSlotComponents := range len(ret.Components) {
		err = ret.Components[iHashedSlotComponents].Type.Encode(w)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Components[iHashedSlotComponents].Hash)
		if err != nil {
			return
		}
	}
	err = queser.VarInt(len(ret.RemoveComponents)).Encode(w)
	if err != nil {
		return
	}
	for iHashedSlotRemoveComponents := range len(ret.RemoveComponents) {
		err = ret.RemoveComponents[iHashedSlotRemoveComponents].Type.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type IDSet struct {
	Val any
}

func (_ IDSet) Decode(r io.Reader) (ret IDSet, err error) {
	var lIDSet queser.VarInt
	lIDSet, err = lIDSet.Decode(r)
	if err != nil {
		return
	}
	lIDSet = lIDSet - 1
	if lIDSet == -1 {
		var IDSetResult string
		IDSetResult, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Val = IDSetResult
		return
	}
	var IDSetResult []queser.VarInt
	for range lIDSet {
		var IDSetElement queser.VarInt
		IDSetElement, err = IDSetElement.Decode(r)
		if err != nil {
			return
		}
		IDSetResult = append(IDSetResult, IDSetElement)
	}
	return
}
func (ret IDSet) Encode(w io.Writer) (err error) {
	err = queser.ToDoError
	return
}

type InstrumentData struct {
	SoundEvent  ItemSoundHolder
	UseDuration float32
	Range       float32
	Description nbt.Anon
}

func (_ InstrumentData) Decode(r io.Reader) (ret InstrumentData, err error) {
	ret.SoundEvent, err = ret.SoundEvent.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.UseDuration)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Range)
	if err != nil {
		return
	}
	ret.Description, err = ret.Description.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret InstrumentData) Encode(w io.Writer) (err error) {
	err = ret.SoundEvent.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.UseDuration)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Range)
	if err != nil {
		return
	}
	err = ret.Description.Encode(w)
	if err != nil {
		return
	}
	return
}

type ItemBlockPredicate struct {
	BlockSet   *any
	Properties *[]ItemBlockProperty
	Nbt        nbt.Anon
	Components DataComponentMatchers
}

func (_ ItemBlockPredicate) Decode(r io.Reader) (ret ItemBlockPredicate, err error) {
	var ItemBlockPredicateBlockSetPresent bool
	err = binary.Read(r, binary.BigEndian, &ItemBlockPredicateBlockSetPresent)
	if err != nil {
		return
	}
	if ItemBlockPredicateBlockSetPresent {
		var ItemBlockPredicateBlockSetPresentValue any
		var lItemBlockPredicateBlockSet queser.VarInt
		lItemBlockPredicateBlockSet, err = lItemBlockPredicateBlockSet.Decode(r)
		if err != nil {
			return
		}
		lItemBlockPredicateBlockSet = lItemBlockPredicateBlockSet - 1
		if lItemBlockPredicateBlockSet == -1 {
			var ItemBlockPredicateBlockSetResult string
			ItemBlockPredicateBlockSetResult, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			ItemBlockPredicateBlockSetPresentValue = ItemBlockPredicateBlockSetResult
			return
		}
		var ItemBlockPredicateBlockSetResult []queser.VarInt
		for range lItemBlockPredicateBlockSet {
			var ItemBlockPredicateBlockSetElement queser.VarInt
			ItemBlockPredicateBlockSetElement, err = ItemBlockPredicateBlockSetElement.Decode(r)
			if err != nil {
				return
			}
			ItemBlockPredicateBlockSetResult = append(ItemBlockPredicateBlockSetResult, ItemBlockPredicateBlockSetElement)
		}
		ret.BlockSet = &ItemBlockPredicateBlockSetPresentValue
	}
	var ItemBlockPredicatePropertiesPresent bool
	err = binary.Read(r, binary.BigEndian, &ItemBlockPredicatePropertiesPresent)
	if err != nil {
		return
	}
	if ItemBlockPredicatePropertiesPresent {
		var ItemBlockPredicatePropertiesPresentValue []ItemBlockProperty
		var lItemBlockPredicateProperties queser.VarInt
		lItemBlockPredicateProperties, err = lItemBlockPredicateProperties.Decode(r)
		if err != nil {
			return
		}
		ItemBlockPredicatePropertiesPresentValue = []ItemBlockProperty{}
		for range lItemBlockPredicateProperties {
			var ItemBlockPredicatePropertiesElement ItemBlockProperty
			ItemBlockPredicatePropertiesElement, err = ItemBlockPredicatePropertiesElement.Decode(r)
			if err != nil {
				return
			}
			ItemBlockPredicatePropertiesPresentValue = append(ItemBlockPredicatePropertiesPresentValue, ItemBlockPredicatePropertiesElement)
		}
		ret.Properties = &ItemBlockPredicatePropertiesPresentValue
	}
	ret.Nbt, err = ret.Nbt.Decode(r)
	if err != nil {
		return
	}
	ret.Components, err = ret.Components.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret ItemBlockPredicate) Encode(w io.Writer) (err error) {
	err = queser.ToDoError
	return
}

type ItemBlockProperty struct {
	Name         string
	IsExactMatch bool
	Value        any
}

func (_ ItemBlockProperty) Decode(r io.Reader) (ret ItemBlockProperty, err error) {
	ret.Name, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IsExactMatch)
	if err != nil {
		return
	}
	switch ret.IsExactMatch {
	case false:
		var ItemBlockPropertyValueTmp struct {
			MinValue string
			MaxValue string
		}
		ItemBlockPropertyValueTmp.MinValue, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ItemBlockPropertyValueTmp.MaxValue, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Value = ItemBlockPropertyValueTmp
	case true:
		var ItemBlockPropertyValueTmp struct {
			ExactValue string
		}
		ItemBlockPropertyValueTmp.ExactValue, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Value = ItemBlockPropertyValueTmp
	}
	return
}
func (ret ItemBlockProperty) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Name)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IsExactMatch)
	if err != nil {
		return
	}
	switch ret.IsExactMatch {
	case false:
		ItemBlockPropertyValue, ok := ret.Value.(struct {
			MinValue string
			MaxValue string
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.EncodeString(w, ItemBlockPropertyValue.MinValue)
		if err != nil {
			return
		}
		err = queser.EncodeString(w, ItemBlockPropertyValue.MaxValue)
		if err != nil {
			return
		}
	case true:
		ItemBlockPropertyValue, ok := ret.Value.(struct {
			ExactValue string
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.EncodeString(w, ItemBlockPropertyValue.ExactValue)
		if err != nil {
			return
		}
	}
	return
}

type ItemBookPage struct {
	Content         string
	FilteredContent *string
}

func (_ ItemBookPage) Decode(r io.Reader) (ret ItemBookPage, err error) {
	ret.Content, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var ItemBookPageFilteredContentPresent bool
	err = binary.Read(r, binary.BigEndian, &ItemBookPageFilteredContentPresent)
	if err != nil {
		return
	}
	if ItemBookPageFilteredContentPresent {
		var ItemBookPageFilteredContentPresentValue string
		ItemBookPageFilteredContentPresentValue, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.FilteredContent = &ItemBookPageFilteredContentPresentValue
	}
	return
}
func (ret ItemBookPage) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Content)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.FilteredContent != nil)
	if err != nil {
		return
	}
	if ret.FilteredContent != nil {
		err = queser.EncodeString(w, *ret.FilteredContent)
		if err != nil {
			return
		}
	}
	return
}

type ItemConsumeEffect struct {
	Type string
	Anon any
}

var ItemConsumeEffectTypeMap = map[queser.VarInt]string{0: "apply_effects", 1: "remove_effects", 2: "clear_all_effects", 3: "teleport_randomly", 4: "play_sound"}

func (_ ItemConsumeEffect) Decode(r io.Reader) (ret ItemConsumeEffect, err error) {
	var ItemConsumeEffectTypeKey queser.VarInt
	ItemConsumeEffectTypeKey, err = ItemConsumeEffectTypeKey.Decode(r)
	if err != nil {
		return
	}
	ret.Type, err = queser.ErroringIndex(ItemConsumeEffectTypeMap, ItemConsumeEffectTypeKey)
	if err != nil {
		return
	}
	switch ret.Type {
	case "apply_effects":
		var ItemConsumeEffectAnonTmp struct {
			Effects     []ItemPotionEffect
			Probability float32
		}
		var lItemConsumeEffectAnonEffects queser.VarInt
		lItemConsumeEffectAnonEffects, err = lItemConsumeEffectAnonEffects.Decode(r)
		if err != nil {
			return
		}
		ItemConsumeEffectAnonTmp.Effects = []ItemPotionEffect{}
		for range lItemConsumeEffectAnonEffects {
			var ItemConsumeEffectAnonEffectsElement ItemPotionEffect
			ItemConsumeEffectAnonEffectsElement, err = ItemConsumeEffectAnonEffectsElement.Decode(r)
			if err != nil {
				return
			}
			ItemConsumeEffectAnonTmp.Effects = append(ItemConsumeEffectAnonTmp.Effects, ItemConsumeEffectAnonEffectsElement)
		}
		err = binary.Read(r, binary.BigEndian, &ItemConsumeEffectAnonTmp.Probability)
		if err != nil {
			return
		}
		ret.Anon = ItemConsumeEffectAnonTmp
	case "clear_all_effects":
		var ItemConsumeEffectAnonTmp queser.Void
		ItemConsumeEffectAnonTmp, err = ItemConsumeEffectAnonTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Anon = ItemConsumeEffectAnonTmp
	case "play_sound":
		var ItemConsumeEffectAnonTmp struct {
			Sound ItemSoundHolder
		}
		ItemConsumeEffectAnonTmp.Sound, err = ItemConsumeEffectAnonTmp.Sound.Decode(r)
		if err != nil {
			return
		}
		ret.Anon = ItemConsumeEffectAnonTmp
	case "remove_effects":
		var ItemConsumeEffectAnonTmp struct {
			Effects IDSet
		}
		ItemConsumeEffectAnonTmp.Effects, err = ItemConsumeEffectAnonTmp.Effects.Decode(r)
		if err != nil {
			return
		}
		ret.Anon = ItemConsumeEffectAnonTmp
	case "teleport_randomly":
		var ItemConsumeEffectAnonTmp struct {
			Diameter float32
		}
		err = binary.Read(r, binary.BigEndian, &ItemConsumeEffectAnonTmp.Diameter)
		if err != nil {
			return
		}
		ret.Anon = ItemConsumeEffectAnonTmp
	}
	return
}

var ItemConsumeEffectTypeReverseMap = map[string]queser.VarInt{"apply_effects": 0, "remove_effects": 1, "clear_all_effects": 2, "teleport_randomly": 3, "play_sound": 4}

func (ret ItemConsumeEffect) Encode(w io.Writer) (err error) {
	var vItemConsumeEffectType queser.VarInt
	vItemConsumeEffectType, err = queser.ErroringIndex(ItemConsumeEffectTypeReverseMap, ret.Type)
	if err != nil {
		return
	}
	err = vItemConsumeEffectType.Encode(w)
	if err != nil {
		return
	}
	switch ret.Type {
	case "apply_effects":
		ItemConsumeEffectAnon, ok := ret.Anon.(struct {
			Effects     []ItemPotionEffect
			Probability float32
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.VarInt(len(ItemConsumeEffectAnon.Effects)).Encode(w)
		if err != nil {
			return
		}
		for iItemConsumeEffectAnonEffects := range len(ItemConsumeEffectAnon.Effects) {
			err = ItemConsumeEffectAnon.Effects[iItemConsumeEffectAnonEffects].Encode(w)
			if err != nil {
				return
			}
		}
		err = binary.Write(w, binary.BigEndian, ItemConsumeEffectAnon.Probability)
		if err != nil {
			return
		}
	case "clear_all_effects":
		ItemConsumeEffectAnon, ok := ret.Anon.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ItemConsumeEffectAnon.Encode(w)
		if err != nil {
			return
		}
	case "play_sound":
		ItemConsumeEffectAnon, ok := ret.Anon.(struct {
			Sound ItemSoundHolder
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ItemConsumeEffectAnon.Sound.Encode(w)
		if err != nil {
			return
		}
	case "remove_effects":
		ItemConsumeEffectAnon, ok := ret.Anon.(struct {
			Effects IDSet
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ItemConsumeEffectAnon.Effects.Encode(w)
		if err != nil {
			return
		}
	case "teleport_randomly":
		ItemConsumeEffectAnon, ok := ret.Anon.(struct {
			Diameter float32
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, ItemConsumeEffectAnon.Diameter)
		if err != nil {
			return
		}
	}
	return
}

type ItemEffectDetail struct {
	Amplifier     queser.VarInt
	Duration      queser.VarInt
	Ambient       bool
	ShowParticles bool
	ShowIcon      bool
	HiddenEffect  *ItemEffectDetail
}

func (_ ItemEffectDetail) Decode(r io.Reader) (ret ItemEffectDetail, err error) {
	ret.Amplifier, err = ret.Amplifier.Decode(r)
	if err != nil {
		return
	}
	ret.Duration, err = ret.Duration.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Ambient)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.ShowParticles)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.ShowIcon)
	if err != nil {
		return
	}
	var ItemEffectDetailHiddenEffectPresent bool
	err = binary.Read(r, binary.BigEndian, &ItemEffectDetailHiddenEffectPresent)
	if err != nil {
		return
	}
	if ItemEffectDetailHiddenEffectPresent {
		var ItemEffectDetailHiddenEffectPresentValue ItemEffectDetail
		ItemEffectDetailHiddenEffectPresentValue, err = ItemEffectDetailHiddenEffectPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.HiddenEffect = &ItemEffectDetailHiddenEffectPresentValue
	}
	return
}
func (ret ItemEffectDetail) Encode(w io.Writer) (err error) {
	err = ret.Amplifier.Encode(w)
	if err != nil {
		return
	}
	err = ret.Duration.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Ambient)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.ShowParticles)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.ShowIcon)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.HiddenEffect != nil)
	if err != nil {
		return
	}
	if ret.HiddenEffect != nil {
		err = (*ret.HiddenEffect).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type ItemFireworkExplosion struct {
	Shape      string
	Colors     []int32
	FadeColors []int32
	HasTrail   bool
	HasTwinkle bool
}

var ItemFireworkExplosionShapeMap = map[queser.VarInt]string{0: "small_ball", 1: "large_ball", 2: "star", 3: "creeper", 4: "burst"}

func (_ ItemFireworkExplosion) Decode(r io.Reader) (ret ItemFireworkExplosion, err error) {
	var ItemFireworkExplosionShapeKey queser.VarInt
	ItemFireworkExplosionShapeKey, err = ItemFireworkExplosionShapeKey.Decode(r)
	if err != nil {
		return
	}
	ret.Shape, err = queser.ErroringIndex(ItemFireworkExplosionShapeMap, ItemFireworkExplosionShapeKey)
	if err != nil {
		return
	}
	var lItemFireworkExplosionColors queser.VarInt
	lItemFireworkExplosionColors, err = lItemFireworkExplosionColors.Decode(r)
	if err != nil {
		return
	}
	ret.Colors = []int32{}
	for range lItemFireworkExplosionColors {
		var ItemFireworkExplosionColorsElement int32
		err = binary.Read(r, binary.BigEndian, &ItemFireworkExplosionColorsElement)
		if err != nil {
			return
		}
		ret.Colors = append(ret.Colors, ItemFireworkExplosionColorsElement)
	}
	var lItemFireworkExplosionFadeColors queser.VarInt
	lItemFireworkExplosionFadeColors, err = lItemFireworkExplosionFadeColors.Decode(r)
	if err != nil {
		return
	}
	ret.FadeColors = []int32{}
	for range lItemFireworkExplosionFadeColors {
		var ItemFireworkExplosionFadeColorsElement int32
		err = binary.Read(r, binary.BigEndian, &ItemFireworkExplosionFadeColorsElement)
		if err != nil {
			return
		}
		ret.FadeColors = append(ret.FadeColors, ItemFireworkExplosionFadeColorsElement)
	}
	err = binary.Read(r, binary.BigEndian, &ret.HasTrail)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.HasTwinkle)
	if err != nil {
		return
	}
	return
}

var ItemFireworkExplosionShapeReverseMap = map[string]queser.VarInt{"small_ball": 0, "large_ball": 1, "star": 2, "creeper": 3, "burst": 4}

func (ret ItemFireworkExplosion) Encode(w io.Writer) (err error) {
	var vItemFireworkExplosionShape queser.VarInt
	vItemFireworkExplosionShape, err = queser.ErroringIndex(ItemFireworkExplosionShapeReverseMap, ret.Shape)
	if err != nil {
		return
	}
	err = vItemFireworkExplosionShape.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Colors)).Encode(w)
	if err != nil {
		return
	}
	for iItemFireworkExplosionColors := range len(ret.Colors) {
		err = binary.Write(w, binary.BigEndian, ret.Colors[iItemFireworkExplosionColors])
		if err != nil {
			return
		}
	}
	err = queser.VarInt(len(ret.FadeColors)).Encode(w)
	if err != nil {
		return
	}
	for iItemFireworkExplosionFadeColors := range len(ret.FadeColors) {
		err = binary.Write(w, binary.BigEndian, ret.FadeColors[iItemFireworkExplosionFadeColors])
		if err != nil {
			return
		}
	}
	err = binary.Write(w, binary.BigEndian, ret.HasTrail)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.HasTwinkle)
	if err != nil {
		return
	}
	return
}

type ItemPotionEffect struct {
	Id      queser.VarInt
	Details ItemEffectDetail
}

func (_ ItemPotionEffect) Decode(r io.Reader) (ret ItemPotionEffect, err error) {
	ret.Id, err = ret.Id.Decode(r)
	if err != nil {
		return
	}
	ret.Details, err = ret.Details.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret ItemPotionEffect) Encode(w io.Writer) (err error) {
	err = ret.Id.Encode(w)
	if err != nil {
		return
	}
	err = ret.Details.Encode(w)
	if err != nil {
		return
	}
	return
}

type ItemSoundEvent struct {
	SoundName  string
	FixedRange *float32
}

func (_ ItemSoundEvent) Decode(r io.Reader) (ret ItemSoundEvent, err error) {
	ret.SoundName, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var ItemSoundEventFixedRangePresent bool
	err = binary.Read(r, binary.BigEndian, &ItemSoundEventFixedRangePresent)
	if err != nil {
		return
	}
	if ItemSoundEventFixedRangePresent {
		var ItemSoundEventFixedRangePresentValue float32
		err = binary.Read(r, binary.BigEndian, &ItemSoundEventFixedRangePresentValue)
		if err != nil {
			return
		}
		ret.FixedRange = &ItemSoundEventFixedRangePresentValue
	}
	return
}
func (ret ItemSoundEvent) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.SoundName)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.FixedRange != nil)
	if err != nil {
		return
	}
	if ret.FixedRange != nil {
		err = binary.Write(w, binary.BigEndian, *ret.FixedRange)
		if err != nil {
			return
		}
	}
	return
}

type ItemSoundHolder struct {
	Val any
}

func (_ ItemSoundHolder) Decode(r io.Reader) (ret ItemSoundHolder, err error) {
	var ItemSoundHolderId queser.VarInt
	ItemSoundHolderId, err = ItemSoundHolderId.Decode(r)
	if err != nil {
		return
	}
	if ItemSoundHolderId != 0 {
		ret.Val = ItemSoundHolderId
		return
	}
	var ItemSoundHolderResult ItemSoundEvent
	ItemSoundHolderResult, err = ItemSoundHolderResult.Decode(r)
	if err != nil {
		return
	}
	ret.Val = ItemSoundHolderResult
	return
}
func (ret ItemSoundHolder) Encode(w io.Writer) (err error) {
	switch ItemSoundHolderKnownType := ret.Val.(type) {
	case queser.VarInt:
		err = ItemSoundHolderKnownType.Encode(w)
		if err != nil {
			return
		}
	case ItemSoundEvent:
		err = ItemSoundHolderKnownType.Encode(w)
		if err != nil {
			return
		}
	default:
		err = queser.BadTypeError
	}
	return
}

type ItemWrittenBookPage struct {
	Content         nbt.Anon
	FilteredContent nbt.Anon
}

func (_ ItemWrittenBookPage) Decode(r io.Reader) (ret ItemWrittenBookPage, err error) {
	ret.Content, err = ret.Content.Decode(r)
	if err != nil {
		return
	}
	ret.FilteredContent, err = ret.FilteredContent.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret ItemWrittenBookPage) Encode(w io.Writer) (err error) {
	err = ret.Content.Encode(w)
	if err != nil {
		return
	}
	err = ret.FilteredContent.Encode(w)
	if err != nil {
		return
	}
	return
}

type JukeboxSongData struct {
	SoundEvent       ItemSoundHolder
	Description      nbt.Anon
	LengthInSeconds  float32
	ComparatorOutput queser.VarInt
}

func (_ JukeboxSongData) Decode(r io.Reader) (ret JukeboxSongData, err error) {
	ret.SoundEvent, err = ret.SoundEvent.Decode(r)
	if err != nil {
		return
	}
	ret.Description, err = ret.Description.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.LengthInSeconds)
	if err != nil {
		return
	}
	ret.ComparatorOutput, err = ret.ComparatorOutput.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret JukeboxSongData) Encode(w io.Writer) (err error) {
	err = ret.SoundEvent.Encode(w)
	if err != nil {
		return
	}
	err = ret.Description.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.LengthInSeconds)
	if err != nil {
		return
	}
	err = ret.ComparatorOutput.Encode(w)
	if err != nil {
		return
	}
	return
}

type Particle struct {
	Type string
	Data any
}

var ParticleTypeMap = map[queser.VarInt]string{0: "angry_villager", 1: "block", 10: "landing_lava", 100: "electric_spark", 101: "scrape", 102: "shriek", 103: "egg_crack", 104: "dust_plume", 105: "trial_spawner_detected_player", 106: "trial_spawner_detected_player_ominous", 107: "vault_connection", 108: "dust_pillar", 109: "ominous_spawning", 11: "dripping_water", 110: "raid_omen", 111: "trial_omen", 112: "block_crumble", 113: "firefly", 12: "falling_water", 13: "dust", 14: "dust_color_transition", 15: "effect", 16: "elder_guardian", 17: "enchanted_hit", 18: "enchant", 19: "end_rod", 2: "block_marker", 20: "entity_effect", 21: "explosion_emitter", 22: "explosion", 23: "gust", 24: "small_gust", 25: "gust_emitter_large", 26: "gust_emitter_small", 27: "sonic_boom", 28: "falling_dust", 29: "firework", 3: "bubble", 30: "fishing", 31: "flame", 32: "infested", 33: "cherry_leaves", 34: "pale_oak_leaves", 35: "tinted_leaves", 36: "sculk_soul", 37: "sculk_charge", 38: "sculk_charge_pop", 39: "soul_fire_flame", 4: "cloud", 40: "soul", 41: "flash", 42: "happy_villager", 43: "composter", 44: "heart", 45: "instant_effect", 46: "item", 47: "vibration", 48: "trail", 49: "item_slime", 5: "crit", 50: "item_cobweb", 51: "item_snowball", 52: "large_smoke", 53: "lava", 54: "mycelium", 55: "note", 56: "poof", 57: "portal", 58: "rain", 59: "smoke", 6: "damage_indicator", 60: "white_smoke", 61: "sneeze", 62: "spit", 63: "squid_ink", 64: "sweep_attack", 65: "totem_of_undying", 66: "underwater", 67: "splash", 68: "witch", 69: "bubble_pop", 7: "dragon_breath", 70: "current_down", 71: "bubble_column_up", 72: "nautilus", 73: "dolphin", 74: "campfire_cosy_smoke", 75: "campfire_signal_smoke", 76: "dripping_honey", 77: "falling_honey", 78: "landing_honey", 79: "falling_nectar", 8: "dripping_lava", 80: "falling_spore_blossom", 81: "ash", 82: "crimson_spore", 83: "warped_spore", 84: "spore_blossom_air", 85: "dripping_obsidian_tear", 86: "falling_obsidian_tear", 87: "landing_obsidian_tear", 88: "reverse_portal", 89: "white_ash", 9: "falling_lava", 90: "small_flame", 91: "snowflake", 92: "dripping_dripstone_lava", 93: "falling_dripstone_lava", 94: "dripping_dripstone_water", 95: "falling_dripstone_water", 96: "glow_squid_ink", 97: "glow", 98: "wax_on", 99: "wax_off"}
var ParticleDataPositionTypeMap = map[queser.VarInt]string{0: "block", 1: "entity"}

func (_ Particle) Decode(r io.Reader) (ret Particle, err error) {
	var ParticleTypeKey queser.VarInt
	ParticleTypeKey, err = ParticleTypeKey.Decode(r)
	if err != nil {
		return
	}
	ret.Type, err = queser.ErroringIndex(ParticleTypeMap, ParticleTypeKey)
	if err != nil {
		return
	}
	switch ret.Type {
	case "block":
		var ParticleDataTmp queser.VarInt
		ParticleDataTmp, err = ParticleDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "block_crumble":
		var ParticleDataTmp queser.VarInt
		ParticleDataTmp, err = ParticleDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "block_marker":
		var ParticleDataTmp queser.VarInt
		ParticleDataTmp, err = ParticleDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "dust":
		var ParticleDataTmp struct {
			Red   float32
			Green float32
			Blue  float32
			Scale float32
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.Red)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.Green)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.Blue)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.Scale)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "dust_color_transition":
		var ParticleDataTmp struct {
			FromRed   float32
			FromGreen float32
			FromBlue  float32
			Scale     float32
			ToRed     float32
			ToGreen   float32
			ToBlue    float32
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.FromRed)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.FromGreen)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.FromBlue)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.Scale)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.ToRed)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.ToGreen)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.ToBlue)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "dust_pillar":
		var ParticleDataTmp queser.VarInt
		ParticleDataTmp, err = ParticleDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "entity_effect":
		var ParticleDataTmp int32
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "falling_dust":
		var ParticleDataTmp queser.VarInt
		ParticleDataTmp, err = ParticleDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "firefly":
		var ParticleDataTmp queser.Void
		ParticleDataTmp, err = ParticleDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "item":
		var ParticleDataTmp Slot
		ParticleDataTmp, err = ParticleDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "sculk_charge":
		var ParticleDataTmp float32
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "shriek":
		var ParticleDataTmp queser.VarInt
		ParticleDataTmp, err = ParticleDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "tinted_leaves":
		var ParticleDataTmp int32
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "trail":
		var ParticleDataTmp struct {
			Target Vec3f64
			Color  uint8
		}
		ParticleDataTmp.Target, err = ParticleDataTmp.Target.Decode(r)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &ParticleDataTmp.Color)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	case "vibration":
		var ParticleDataTmp struct {
			PositionType string
			Position     any
			Ticks        queser.VarInt
		}
		var ParticleDataPositionTypeKey queser.VarInt
		ParticleDataPositionTypeKey, err = ParticleDataPositionTypeKey.Decode(r)
		if err != nil {
			return
		}
		ParticleDataTmp.PositionType, err = queser.ErroringIndex(ParticleDataPositionTypeMap, ParticleDataPositionTypeKey)
		if err != nil {
			return
		}
		switch ParticleDataTmp.PositionType {
		case "block":
			var ParticleDataPositionTmp Position
			ParticleDataPositionTmp, err = ParticleDataPositionTmp.Decode(r)
			if err != nil {
				return
			}
			ParticleDataTmp.Position = ParticleDataPositionTmp
		case "entity":
			var ParticleDataPositionTmp struct {
				EntityId        queser.VarInt
				EntityEyeHeight float32
			}
			ParticleDataPositionTmp.EntityId, err = ParticleDataPositionTmp.EntityId.Decode(r)
			if err != nil {
				return
			}
			err = binary.Read(r, binary.BigEndian, &ParticleDataPositionTmp.EntityEyeHeight)
			if err != nil {
				return
			}
			ParticleDataTmp.Position = ParticleDataPositionTmp
		}
		ParticleDataTmp.Ticks, err = ParticleDataTmp.Ticks.Decode(r)
		if err != nil {
			return
		}
		ret.Data = ParticleDataTmp
	}
	return
}

var ParticleTypeReverseMap = map[string]queser.VarInt{"angry_villager": 0, "block": 1, "landing_lava": 10, "electric_spark": 100, "scrape": 101, "shriek": 102, "egg_crack": 103, "dust_plume": 104, "trial_spawner_detected_player": 105, "trial_spawner_detected_player_ominous": 106, "vault_connection": 107, "dust_pillar": 108, "ominous_spawning": 109, "dripping_water": 11, "raid_omen": 110, "trial_omen": 111, "block_crumble": 112, "firefly": 113, "falling_water": 12, "dust": 13, "dust_color_transition": 14, "effect": 15, "elder_guardian": 16, "enchanted_hit": 17, "enchant": 18, "end_rod": 19, "block_marker": 2, "entity_effect": 20, "explosion_emitter": 21, "explosion": 22, "gust": 23, "small_gust": 24, "gust_emitter_large": 25, "gust_emitter_small": 26, "sonic_boom": 27, "falling_dust": 28, "firework": 29, "bubble": 3, "fishing": 30, "flame": 31, "infested": 32, "cherry_leaves": 33, "pale_oak_leaves": 34, "tinted_leaves": 35, "sculk_soul": 36, "sculk_charge": 37, "sculk_charge_pop": 38, "soul_fire_flame": 39, "cloud": 4, "soul": 40, "flash": 41, "happy_villager": 42, "composter": 43, "heart": 44, "instant_effect": 45, "item": 46, "vibration": 47, "trail": 48, "item_slime": 49, "crit": 5, "item_cobweb": 50, "item_snowball": 51, "large_smoke": 52, "lava": 53, "mycelium": 54, "note": 55, "poof": 56, "portal": 57, "rain": 58, "smoke": 59, "damage_indicator": 6, "white_smoke": 60, "sneeze": 61, "spit": 62, "squid_ink": 63, "sweep_attack": 64, "totem_of_undying": 65, "underwater": 66, "splash": 67, "witch": 68, "bubble_pop": 69, "dragon_breath": 7, "current_down": 70, "bubble_column_up": 71, "nautilus": 72, "dolphin": 73, "campfire_cosy_smoke": 74, "campfire_signal_smoke": 75, "dripping_honey": 76, "falling_honey": 77, "landing_honey": 78, "falling_nectar": 79, "dripping_lava": 8, "falling_spore_blossom": 80, "ash": 81, "crimson_spore": 82, "warped_spore": 83, "spore_blossom_air": 84, "dripping_obsidian_tear": 85, "falling_obsidian_tear": 86, "landing_obsidian_tear": 87, "reverse_portal": 88, "white_ash": 89, "falling_lava": 9, "small_flame": 90, "snowflake": 91, "dripping_dripstone_lava": 92, "falling_dripstone_lava": 93, "dripping_dripstone_water": 94, "falling_dripstone_water": 95, "glow_squid_ink": 96, "glow": 97, "wax_on": 98, "wax_off": 99}
var ParticleDataPositionTypeReverseMap = map[string]queser.VarInt{"block": 0, "entity": 1}

func (ret Particle) Encode(w io.Writer) (err error) {
	var vParticleType queser.VarInt
	vParticleType, err = queser.ErroringIndex(ParticleTypeReverseMap, ret.Type)
	if err != nil {
		return
	}
	err = vParticleType.Encode(w)
	if err != nil {
		return
	}
	switch ret.Type {
	case "block":
		ParticleData, ok := ret.Data.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ParticleData.Encode(w)
		if err != nil {
			return
		}
	case "block_crumble":
		ParticleData, ok := ret.Data.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ParticleData.Encode(w)
		if err != nil {
			return
		}
	case "block_marker":
		ParticleData, ok := ret.Data.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ParticleData.Encode(w)
		if err != nil {
			return
		}
	case "dust":
		ParticleData, ok := ret.Data.(struct {
			Red   float32
			Green float32
			Blue  float32
			Scale float32
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.Red)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.Green)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.Blue)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.Scale)
		if err != nil {
			return
		}
	case "dust_color_transition":
		ParticleData, ok := ret.Data.(struct {
			FromRed   float32
			FromGreen float32
			FromBlue  float32
			Scale     float32
			ToRed     float32
			ToGreen   float32
			ToBlue    float32
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.FromRed)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.FromGreen)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.FromBlue)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.Scale)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.ToRed)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.ToGreen)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.ToBlue)
		if err != nil {
			return
		}
	case "dust_pillar":
		ParticleData, ok := ret.Data.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ParticleData.Encode(w)
		if err != nil {
			return
		}
	case "entity_effect":
		ParticleData, ok := ret.Data.(int32)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData)
		if err != nil {
			return
		}
	case "falling_dust":
		ParticleData, ok := ret.Data.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ParticleData.Encode(w)
		if err != nil {
			return
		}
	case "firefly":
		ParticleData, ok := ret.Data.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ParticleData.Encode(w)
		if err != nil {
			return
		}
	case "item":
		ParticleData, ok := ret.Data.(Slot)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ParticleData.Encode(w)
		if err != nil {
			return
		}
	case "sculk_charge":
		ParticleData, ok := ret.Data.(float32)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData)
		if err != nil {
			return
		}
	case "shriek":
		ParticleData, ok := ret.Data.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ParticleData.Encode(w)
		if err != nil {
			return
		}
	case "tinted_leaves":
		ParticleData, ok := ret.Data.(int32)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData)
		if err != nil {
			return
		}
	case "trail":
		ParticleData, ok := ret.Data.(struct {
			Target Vec3f64
			Color  uint8
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ParticleData.Target.Encode(w)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ParticleData.Color)
		if err != nil {
			return
		}
	case "vibration":
		ParticleData, ok := ret.Data.(struct {
			PositionType string
			Position     any
			Ticks        queser.VarInt
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		var vParticleDataPositionType queser.VarInt
		vParticleDataPositionType, err = queser.ErroringIndex(ParticleDataPositionTypeReverseMap, ParticleData.PositionType)
		if err != nil {
			return
		}
		err = vParticleDataPositionType.Encode(w)
		if err != nil {
			return
		}
		switch ParticleData.PositionType {
		case "block":
			ParticleDataPosition, ok := ParticleData.Position.(Position)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = ParticleDataPosition.Encode(w)
			if err != nil {
				return
			}
		case "entity":
			ParticleDataPosition, ok := ParticleData.Position.(struct {
				EntityId        queser.VarInt
				EntityEyeHeight float32
			})
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = ParticleDataPosition.EntityId.Encode(w)
			if err != nil {
				return
			}
			err = binary.Write(w, binary.BigEndian, ParticleDataPosition.EntityEyeHeight)
			if err != nil {
				return
			}
		}
		err = ParticleData.Ticks.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type ServerLinkType struct {
	Val string
}

var ServerLinkTypeMap = map[queser.VarInt]string{0: "bug_report", 1: "community_guidelines", 2: "support", 3: "status", 4: "feedback", 5: "community", 6: "website", 7: "forums", 8: "news", 9: "announcements"}

func (_ ServerLinkType) Decode(r io.Reader) (ret ServerLinkType, err error) {
	var ServerLinkTypeKey queser.VarInt
	ServerLinkTypeKey, err = ServerLinkTypeKey.Decode(r)
	if err != nil {
		return
	}
	ret.Val, err = queser.ErroringIndex(ServerLinkTypeMap, ServerLinkTypeKey)
	if err != nil {
		return
	}
	return
}

var ServerLinkTypeReverseMap = map[string]queser.VarInt{"bug_report": 0, "community_guidelines": 1, "support": 2, "status": 3, "feedback": 4, "community": 5, "website": 6, "forums": 7, "news": 8, "announcements": 9}

func (ret ServerLinkType) Encode(w io.Writer) (err error) {
	var vServerLinkType queser.VarInt
	vServerLinkType, err = queser.ErroringIndex(ServerLinkTypeReverseMap, ret.Val)
	if err != nil {
		return
	}
	err = vServerLinkType.Encode(w)
	if err != nil {
		return
	}
	return
}

type Slot struct {
	ItemCount queser.VarInt
	Anon      any
}

func (_ Slot) Decode(r io.Reader) (ret Slot, err error) {
	err = queser.ToDoError
	return
}
func (ret Slot) Encode(w io.Writer) (err error) {
	err = queser.ToDoError
	return
}

type SlotComponent struct {
	Type SlotComponentType
	Data any
}

func (_ SlotComponent) Decode(r io.Reader) (ret SlotComponent, err error) {
	err = queser.ToDoError
	return
}
func (ret SlotComponent) Encode(w io.Writer) (err error) {
	err = queser.ToDoError
	return
}

type SlotComponentType struct {
	Val string
}

var SlotComponentTypeMap = map[queser.VarInt]string{0: "custom_data", 1: "max_stack_size", 10: "enchantments", 11: "can_place_on", 12: "can_break", 13: "attribute_modifiers", 14: "custom_model_data", 15: "tooltip_display", 16: "repair_cost", 17: "creative_slot_lock", 18: "enchantment_glint_override", 19: "intangible_projectile", 2: "max_damage", 20: "food", 21: "consumable", 22: "use_remainder", 23: "use_cooldown", 24: "damage_resistant", 25: "tool", 26: "weapon", 27: "enchantable", 28: "equippable", 29: "repairable", 3: "damage", 30: "glider", 31: "tooltip_style", 32: "death_protection", 33: "blocks_attacks", 34: "stored_enchantments", 35: "dyed_color", 36: "map_color", 37: "map_id", 38: "map_decorations", 39: "map_post_processing", 4: "unbreakable", 40: "potion_duration_scale", 41: "charged_projectiles", 42: "bundle_contents", 43: "potion_contents", 44: "suspicious_stew_effects", 45: "writable_book_content", 46: "written_book_content", 47: "trim", 48: "debug_stick_state", 49: "entity_data", 5: "custom_name", 50: "bucket_entity_data", 51: "block_entity_data", 52: "instrument", 53: "provides_trim_material", 54: "ominous_bottle_amplifier", 55: "jukebox_playable", 56: "provides_banner_patterns", 57: "recipes", 58: "lodestone_tracker", 59: "firework_explosion", 6: "item_name", 60: "fireworks", 61: "profile", 62: "note_block_sound", 63: "banner_patterns", 64: "base_color", 65: "pot_decorations", 66: "container", 67: "block_state", 68: "bees", 69: "lock", 7: "item_model", 70: "container_loot", 71: "break_sound", 72: "villager/variant", 73: "wolf/variant", 74: "wolf/sound_variant", 75: "wolf/collar", 76: "fox/variant", 77: "salmon/size", 78: "parrot/variant", 79: "tropical_fish/pattern", 8: "lore", 80: "tropical_fish/base_color", 81: "tropical_fish/pattern_color", 82: "mooshroom/variant", 83: "rabbit/variant", 84: "pig/variant", 85: "cow/variant", 86: "chicken/variant", 87: "frog/variant", 88: "horse/variant", 89: "painting/variant", 9: "rarity", 90: "llama/variant", 91: "axolotl/variant", 92: "cat/variant", 93: "cat/collar", 94: "sheep/color", 95: "shulker/color"}

func (_ SlotComponentType) Decode(r io.Reader) (ret SlotComponentType, err error) {
	var SlotComponentTypeKey queser.VarInt
	SlotComponentTypeKey, err = SlotComponentTypeKey.Decode(r)
	if err != nil {
		return
	}
	ret.Val, err = queser.ErroringIndex(SlotComponentTypeMap, SlotComponentTypeKey)
	if err != nil {
		return
	}
	return
}

var SlotComponentTypeReverseMap = map[string]queser.VarInt{"custom_data": 0, "max_stack_size": 1, "enchantments": 10, "can_place_on": 11, "can_break": 12, "attribute_modifiers": 13, "custom_model_data": 14, "tooltip_display": 15, "repair_cost": 16, "creative_slot_lock": 17, "enchantment_glint_override": 18, "intangible_projectile": 19, "max_damage": 2, "food": 20, "consumable": 21, "use_remainder": 22, "use_cooldown": 23, "damage_resistant": 24, "tool": 25, "weapon": 26, "enchantable": 27, "equippable": 28, "repairable": 29, "damage": 3, "glider": 30, "tooltip_style": 31, "death_protection": 32, "blocks_attacks": 33, "stored_enchantments": 34, "dyed_color": 35, "map_color": 36, "map_id": 37, "map_decorations": 38, "map_post_processing": 39, "unbreakable": 4, "potion_duration_scale": 40, "charged_projectiles": 41, "bundle_contents": 42, "potion_contents": 43, "suspicious_stew_effects": 44, "writable_book_content": 45, "written_book_content": 46, "trim": 47, "debug_stick_state": 48, "entity_data": 49, "custom_name": 5, "bucket_entity_data": 50, "block_entity_data": 51, "instrument": 52, "provides_trim_material": 53, "ominous_bottle_amplifier": 54, "jukebox_playable": 55, "provides_banner_patterns": 56, "recipes": 57, "lodestone_tracker": 58, "firework_explosion": 59, "item_name": 6, "fireworks": 60, "profile": 61, "note_block_sound": 62, "banner_patterns": 63, "base_color": 64, "pot_decorations": 65, "container": 66, "block_state": 67, "bees": 68, "lock": 69, "item_model": 7, "container_loot": 70, "break_sound": 71, "villager/variant": 72, "wolf/variant": 73, "wolf/sound_variant": 74, "wolf/collar": 75, "fox/variant": 76, "salmon/size": 77, "parrot/variant": 78, "tropical_fish/pattern": 79, "lore": 8, "tropical_fish/base_color": 80, "tropical_fish/pattern_color": 81, "mooshroom/variant": 82, "rabbit/variant": 83, "pig/variant": 84, "cow/variant": 85, "chicken/variant": 86, "frog/variant": 87, "horse/variant": 88, "painting/variant": 89, "rarity": 9, "llama/variant": 90, "axolotl/variant": 91, "cat/variant": 92, "cat/collar": 93, "sheep/color": 94, "shulker/color": 95}

func (ret SlotComponentType) Encode(w io.Writer) (err error) {
	var vSlotComponentType queser.VarInt
	vSlotComponentType, err = queser.ErroringIndex(SlotComponentTypeReverseMap, ret.Val)
	if err != nil {
		return
	}
	err = vSlotComponentType.Encode(w)
	if err != nil {
		return
	}
	return
}

type UntrustedSlot struct {
	ItemCount queser.VarInt
	Anon      any
}

func (_ UntrustedSlot) Decode(r io.Reader) (ret UntrustedSlot, err error) {
	err = queser.ToDoError
	return
}
func (ret UntrustedSlot) Encode(w io.Writer) (err error) {
	err = queser.ToDoError
	return
}

type UntrustedSlotComponent struct {
	Type SlotComponentType
	Data ByteArray
}

func (_ UntrustedSlotComponent) Decode(r io.Reader) (ret UntrustedSlotComponent, err error) {
	ret.Type, err = ret.Type.Decode(r)
	if err != nil {
		return
	}
	ret.Data, err = ret.Data.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret UntrustedSlotComponent) Encode(w io.Writer) (err error) {
	err = ret.Type.Encode(w)
	if err != nil {
		return
	}
	err = ret.Data.Encode(w)
	if err != nil {
		return
	}
	return
}

type ChatSession struct {
	Val *struct {
		Uuid      uuid.UUID
		PublicKey struct {
			ExpireTime   int64
			KeyBytes     []byte
			KeySignature []byte
		}
	}
}

func (_ ChatSession) Decode(r io.Reader) (ret ChatSession, err error) {
	var ChatSessionPresent bool
	err = binary.Read(r, binary.BigEndian, &ChatSessionPresent)
	if err != nil {
		return
	}
	if ChatSessionPresent {
		var ChatSessionPresentValue struct {
			Uuid      uuid.UUID
			PublicKey struct {
				ExpireTime   int64
				KeyBytes     []byte
				KeySignature []byte
			}
		}
		_, err = io.ReadFull(r, ChatSessionPresentValue.Uuid[:])
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &ChatSessionPresentValue.PublicKey.ExpireTime)
		if err != nil {
			return
		}
		var lChatSessionPublicKeyKeyBytes queser.VarInt
		lChatSessionPublicKeyKeyBytes, err = lChatSessionPublicKeyKeyBytes.Decode(r)
		if err != nil {
			return
		}
		ChatSessionPresentValue.PublicKey.KeyBytes, err = io.ReadAll(io.LimitReader(r, int64(lChatSessionPublicKeyKeyBytes)))
		if err != nil {
			return
		}
		var lChatSessionPublicKeyKeySignature queser.VarInt
		lChatSessionPublicKeyKeySignature, err = lChatSessionPublicKeyKeySignature.Decode(r)
		if err != nil {
			return
		}
		ChatSessionPresentValue.PublicKey.KeySignature, err = io.ReadAll(io.LimitReader(r, int64(lChatSessionPublicKeyKeySignature)))
		if err != nil {
			return
		}
		ret.Val = &ChatSessionPresentValue
	}
	return
}
func (ret ChatSession) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Val != nil)
	if err != nil {
		return
	}
	if ret.Val != nil {
		_, err = w.Write((*ret.Val).Uuid[:])
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, (*ret.Val).PublicKey.ExpireTime)
		if err != nil {
			return
		}
		err = queser.VarInt(len((*ret.Val).PublicKey.KeyBytes)).Encode(w)
		if err != nil {
			return
		}
		_, err = w.Write((*ret.Val).PublicKey.KeyBytes)
		if err != nil {
			return
		}
		err = queser.VarInt(len((*ret.Val).PublicKey.KeySignature)).Encode(w)
		if err != nil {
			return
		}
		_, err = w.Write((*ret.Val).PublicKey.KeySignature)
		if err != nil {
			return
		}
	}
	return
}

type ChunkBlockEntity struct {
	Anon    uint8
	Y       int16
	Type    queser.VarInt
	NbtData nbt.Anon
}

func (_ ChunkBlockEntity) Decode(r io.Reader) (ret ChunkBlockEntity, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Anon)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	ret.Type, err = ret.Type.Decode(r)
	if err != nil {
		return
	}
	ret.NbtData, err = ret.NbtData.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret ChunkBlockEntity) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Anon)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = ret.Type.Encode(w)
	if err != nil {
		return
	}
	err = ret.NbtData.Encode(w)
	if err != nil {
		return
	}
	return
}

type CommandNode struct {
	Flags         uint8
	Children      []queser.VarInt
	RedirectNode  any
	ExtraNodeData any
}

var CommandNodeExtraNodeDataParserMap = map[queser.VarInt]string{0: "brigadier:bool", 1: "brigadier:float", 10: "minecraft:vec3", 11: "minecraft:vec2", 12: "minecraft:block_state", 13: "minecraft:block_predicate", 14: "minecraft:item_stack", 15: "minecraft:item_predicate", 16: "minecraft:color", 17: "minecraft:hex_color", 18: "minecraft:component", 19: "minecraft:style", 2: "brigadier:double", 20: "minecraft:message", 21: "minecraft:nbt", 22: "minecraft:nbt_tag", 23: "minecraft:nbt_path", 24: "minecraft:objective", 25: "minecraft:objective_criteria", 26: "minecraft:operation", 27: "minecraft:particle", 28: "minecraft:angle", 29: "minecraft:rotation", 3: "brigadier:integer", 30: "minecraft:scoreboard_slot", 31: "minecraft:score_holder", 32: "minecraft:swizzle", 33: "minecraft:team", 34: "minecraft:item_slot", 35: "minecraft:item_slots", 36: "minecraft:resource_location", 37: "minecraft:function", 38: "minecraft:entity_anchor", 39: "minecraft:int_range", 4: "brigadier:long", 40: "minecraft:float_range", 41: "minecraft:dimension", 42: "minecraft:gamemode", 43: "minecraft:time", 44: "minecraft:resource_or_tag", 45: "minecraft:resource_or_tag_key", 46: "minecraft:resource", 47: "minecraft:resource_key", 48: "minecraft:resource_selector", 49: "minecraft:template_mirror", 5: "brigadier:string", 50: "minecraft:template_rotation", 51: "minecraft:heightmap", 52: "minecraft:loot_table", 53: "minecraft:loot_predicate", 54: "minecraft:loot_modifier", 55: "minecraft:dialog", 56: "minecraft:uuid", 6: "minecraft:entity", 7: "minecraft:game_profile", 8: "minecraft:block_pos", 9: "minecraft:column_pos"}
var CommandNodeExtraNodeDataPropertiesMap = map[queser.VarInt]string{0: "SINGLE_WORD", 1: "QUOTABLE_PHRASE", 2: "GREEDY_PHRASE"}

func (_ CommandNode) Decode(r io.Reader) (ret CommandNode, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Flags)
	if err != nil {
		return
	}
	var lCommandNodeChildren queser.VarInt
	lCommandNodeChildren, err = lCommandNodeChildren.Decode(r)
	if err != nil {
		return
	}
	ret.Children = []queser.VarInt{}
	for range lCommandNodeChildren {
		var CommandNodeChildrenElement queser.VarInt
		CommandNodeChildrenElement, err = CommandNodeChildrenElement.Decode(r)
		if err != nil {
			return
		}
		ret.Children = append(ret.Children, CommandNodeChildrenElement)
	}
	switch ret.Flags {
	case 1:
		var CommandNodeRedirectNodeTmp queser.VarInt
		CommandNodeRedirectNodeTmp, err = CommandNodeRedirectNodeTmp.Decode(r)
		if err != nil {
			return
		}
		ret.RedirectNode = CommandNodeRedirectNodeTmp
	default:
		var CommandNodeRedirectNodeTmp queser.Void
		CommandNodeRedirectNodeTmp, err = CommandNodeRedirectNodeTmp.Decode(r)
		if err != nil {
			return
		}
		ret.RedirectNode = CommandNodeRedirectNodeTmp
	}
	switch ret.Flags {
	case 0:
		var CommandNodeExtraNodeDataTmp queser.Void
		CommandNodeExtraNodeDataTmp, err = CommandNodeExtraNodeDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.ExtraNodeData = CommandNodeExtraNodeDataTmp
	case 1:
		var CommandNodeExtraNodeDataTmp struct {
			Name string
		}
		CommandNodeExtraNodeDataTmp.Name, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.ExtraNodeData = CommandNodeExtraNodeDataTmp
	case 2:
		var CommandNodeExtraNodeDataTmp struct {
			Name           string
			Parser         string
			Properties     any
			SuggestionType any
		}
		CommandNodeExtraNodeDataTmp.Name, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		var CommandNodeExtraNodeDataParserKey queser.VarInt
		CommandNodeExtraNodeDataParserKey, err = CommandNodeExtraNodeDataParserKey.Decode(r)
		if err != nil {
			return
		}
		CommandNodeExtraNodeDataTmp.Parser, err = queser.ErroringIndex(CommandNodeExtraNodeDataParserMap, CommandNodeExtraNodeDataParserKey)
		if err != nil {
			return
		}
		switch CommandNodeExtraNodeDataTmp.Parser {
		case "brigadier:bool":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "brigadier:double":
			var CommandNodeExtraNodeDataPropertiesTmp struct {
				Flags uint8
				Min   any
				Max   any
			}
			err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesTmp.Flags)
			if err != nil {
				return
			}
			switch CommandNodeExtraNodeDataPropertiesTmp.Flags {
			case 1:
				var CommandNodeExtraNodeDataPropertiesMinTmp float64
				err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesMinTmp)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Min = CommandNodeExtraNodeDataPropertiesMinTmp
			default:
				var CommandNodeExtraNodeDataPropertiesMinTmp queser.Void
				CommandNodeExtraNodeDataPropertiesMinTmp, err = CommandNodeExtraNodeDataPropertiesMinTmp.Decode(r)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Min = CommandNodeExtraNodeDataPropertiesMinTmp
			}
			switch CommandNodeExtraNodeDataPropertiesTmp.Flags {
			case 1:
				var CommandNodeExtraNodeDataPropertiesMaxTmp float64
				err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesMaxTmp)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Max = CommandNodeExtraNodeDataPropertiesMaxTmp
			default:
				var CommandNodeExtraNodeDataPropertiesMaxTmp queser.Void
				CommandNodeExtraNodeDataPropertiesMaxTmp, err = CommandNodeExtraNodeDataPropertiesMaxTmp.Decode(r)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Max = CommandNodeExtraNodeDataPropertiesMaxTmp
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "brigadier:float":
			var CommandNodeExtraNodeDataPropertiesTmp struct {
				Flags uint8
				Min   any
				Max   any
			}
			err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesTmp.Flags)
			if err != nil {
				return
			}
			switch CommandNodeExtraNodeDataPropertiesTmp.Flags {
			case 1:
				var CommandNodeExtraNodeDataPropertiesMinTmp float32
				err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesMinTmp)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Min = CommandNodeExtraNodeDataPropertiesMinTmp
			default:
				var CommandNodeExtraNodeDataPropertiesMinTmp queser.Void
				CommandNodeExtraNodeDataPropertiesMinTmp, err = CommandNodeExtraNodeDataPropertiesMinTmp.Decode(r)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Min = CommandNodeExtraNodeDataPropertiesMinTmp
			}
			switch CommandNodeExtraNodeDataPropertiesTmp.Flags {
			case 1:
				var CommandNodeExtraNodeDataPropertiesMaxTmp float32
				err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesMaxTmp)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Max = CommandNodeExtraNodeDataPropertiesMaxTmp
			default:
				var CommandNodeExtraNodeDataPropertiesMaxTmp queser.Void
				CommandNodeExtraNodeDataPropertiesMaxTmp, err = CommandNodeExtraNodeDataPropertiesMaxTmp.Decode(r)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Max = CommandNodeExtraNodeDataPropertiesMaxTmp
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "brigadier:integer":
			var CommandNodeExtraNodeDataPropertiesTmp struct {
				Flags uint8
				Min   any
				Max   any
			}
			err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesTmp.Flags)
			if err != nil {
				return
			}
			switch CommandNodeExtraNodeDataPropertiesTmp.Flags {
			case 1:
				var CommandNodeExtraNodeDataPropertiesMinTmp int32
				err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesMinTmp)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Min = CommandNodeExtraNodeDataPropertiesMinTmp
			default:
				var CommandNodeExtraNodeDataPropertiesMinTmp queser.Void
				CommandNodeExtraNodeDataPropertiesMinTmp, err = CommandNodeExtraNodeDataPropertiesMinTmp.Decode(r)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Min = CommandNodeExtraNodeDataPropertiesMinTmp
			}
			switch CommandNodeExtraNodeDataPropertiesTmp.Flags {
			case 1:
				var CommandNodeExtraNodeDataPropertiesMaxTmp int32
				err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesMaxTmp)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Max = CommandNodeExtraNodeDataPropertiesMaxTmp
			default:
				var CommandNodeExtraNodeDataPropertiesMaxTmp queser.Void
				CommandNodeExtraNodeDataPropertiesMaxTmp, err = CommandNodeExtraNodeDataPropertiesMaxTmp.Decode(r)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Max = CommandNodeExtraNodeDataPropertiesMaxTmp
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "brigadier:long":
			var CommandNodeExtraNodeDataPropertiesTmp struct {
				Flags uint8
				Min   any
				Max   any
			}
			err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesTmp.Flags)
			if err != nil {
				return
			}
			switch CommandNodeExtraNodeDataPropertiesTmp.Flags {
			case 1:
				var CommandNodeExtraNodeDataPropertiesMinTmp int64
				err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesMinTmp)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Min = CommandNodeExtraNodeDataPropertiesMinTmp
			default:
				var CommandNodeExtraNodeDataPropertiesMinTmp queser.Void
				CommandNodeExtraNodeDataPropertiesMinTmp, err = CommandNodeExtraNodeDataPropertiesMinTmp.Decode(r)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Min = CommandNodeExtraNodeDataPropertiesMinTmp
			}
			switch CommandNodeExtraNodeDataPropertiesTmp.Flags {
			case 1:
				var CommandNodeExtraNodeDataPropertiesMaxTmp int64
				err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesMaxTmp)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Max = CommandNodeExtraNodeDataPropertiesMaxTmp
			default:
				var CommandNodeExtraNodeDataPropertiesMaxTmp queser.Void
				CommandNodeExtraNodeDataPropertiesMaxTmp, err = CommandNodeExtraNodeDataPropertiesMaxTmp.Decode(r)
				if err != nil {
					return
				}
				CommandNodeExtraNodeDataPropertiesTmp.Max = CommandNodeExtraNodeDataPropertiesMaxTmp
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "brigadier:string":
			var CommandNodeExtraNodeDataPropertiesTmp string
			var CommandNodeExtraNodeDataPropertiesKey queser.VarInt
			CommandNodeExtraNodeDataPropertiesKey, err = CommandNodeExtraNodeDataPropertiesKey.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataPropertiesTmp, err = queser.ErroringIndex(CommandNodeExtraNodeDataPropertiesMap, CommandNodeExtraNodeDataPropertiesKey)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:angle":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:block_pos":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:block_predicate":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:block_state":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:color":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:column_pos":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:component":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:dialog":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:dimension":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:entity":
			var CommandNodeExtraNodeDataPropertiesTmp uint8
			err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesTmp)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:entity_anchor":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:float_range":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:function":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:game_profile":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:gamemode":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:heightmap":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:hex_color":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:int_range":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:item_predicate":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:item_slot":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:item_stack":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:message":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:nbt":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:nbt_path":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:objective":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:objective_criteria":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:operation":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:particle":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:resource":
			var CommandNodeExtraNodeDataPropertiesTmp struct {
				Registry string
			}
			CommandNodeExtraNodeDataPropertiesTmp.Registry, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:resource_key":
			var CommandNodeExtraNodeDataPropertiesTmp struct {
				Registry string
			}
			CommandNodeExtraNodeDataPropertiesTmp.Registry, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:resource_location":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:resource_or_tag":
			var CommandNodeExtraNodeDataPropertiesTmp struct {
				Registry string
			}
			CommandNodeExtraNodeDataPropertiesTmp.Registry, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:resource_or_tag_key":
			var CommandNodeExtraNodeDataPropertiesTmp struct {
				Registry string
			}
			CommandNodeExtraNodeDataPropertiesTmp.Registry, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:resource_selector":
			var CommandNodeExtraNodeDataPropertiesTmp struct {
				Registry string
			}
			CommandNodeExtraNodeDataPropertiesTmp.Registry, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:rotation":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:score_holder":
			var CommandNodeExtraNodeDataPropertiesTmp uint8
			err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesTmp)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:scoreboard_slot":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:swizzle":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:team":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:template_mirror":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:template_rotation":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:time":
			var CommandNodeExtraNodeDataPropertiesTmp struct {
				Min int32
			}
			err = binary.Read(r, binary.BigEndian, &CommandNodeExtraNodeDataPropertiesTmp.Min)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:uuid":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:vec2":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		case "minecraft:vec3":
			var CommandNodeExtraNodeDataPropertiesTmp queser.Void
			CommandNodeExtraNodeDataPropertiesTmp, err = CommandNodeExtraNodeDataPropertiesTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.Properties = CommandNodeExtraNodeDataPropertiesTmp
		}
		switch ret.Flags {
		case 1:
			var CommandNodeExtraNodeDataSuggestionTypeTmp string
			CommandNodeExtraNodeDataSuggestionTypeTmp, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.SuggestionType = CommandNodeExtraNodeDataSuggestionTypeTmp
		default:
			var CommandNodeExtraNodeDataSuggestionTypeTmp queser.Void
			CommandNodeExtraNodeDataSuggestionTypeTmp, err = CommandNodeExtraNodeDataSuggestionTypeTmp.Decode(r)
			if err != nil {
				return
			}
			CommandNodeExtraNodeDataTmp.SuggestionType = CommandNodeExtraNodeDataSuggestionTypeTmp
		}
		ret.ExtraNodeData = CommandNodeExtraNodeDataTmp
	}
	return
}

var CommandNodeExtraNodeDataParserReverseMap = map[string]queser.VarInt{"brigadier:bool": 0, "brigadier:float": 1, "minecraft:vec3": 10, "minecraft:vec2": 11, "minecraft:block_state": 12, "minecraft:block_predicate": 13, "minecraft:item_stack": 14, "minecraft:item_predicate": 15, "minecraft:color": 16, "minecraft:hex_color": 17, "minecraft:component": 18, "minecraft:style": 19, "brigadier:double": 2, "minecraft:message": 20, "minecraft:nbt": 21, "minecraft:nbt_tag": 22, "minecraft:nbt_path": 23, "minecraft:objective": 24, "minecraft:objective_criteria": 25, "minecraft:operation": 26, "minecraft:particle": 27, "minecraft:angle": 28, "minecraft:rotation": 29, "brigadier:integer": 3, "minecraft:scoreboard_slot": 30, "minecraft:score_holder": 31, "minecraft:swizzle": 32, "minecraft:team": 33, "minecraft:item_slot": 34, "minecraft:item_slots": 35, "minecraft:resource_location": 36, "minecraft:function": 37, "minecraft:entity_anchor": 38, "minecraft:int_range": 39, "brigadier:long": 4, "minecraft:float_range": 40, "minecraft:dimension": 41, "minecraft:gamemode": 42, "minecraft:time": 43, "minecraft:resource_or_tag": 44, "minecraft:resource_or_tag_key": 45, "minecraft:resource": 46, "minecraft:resource_key": 47, "minecraft:resource_selector": 48, "minecraft:template_mirror": 49, "brigadier:string": 5, "minecraft:template_rotation": 50, "minecraft:heightmap": 51, "minecraft:loot_table": 52, "minecraft:loot_predicate": 53, "minecraft:loot_modifier": 54, "minecraft:dialog": 55, "minecraft:uuid": 56, "minecraft:entity": 6, "minecraft:game_profile": 7, "minecraft:block_pos": 8, "minecraft:column_pos": 9}
var CommandNodeExtraNodeDataPropertiesReverseMap = map[string]queser.VarInt{"SINGLE_WORD": 0, "QUOTABLE_PHRASE": 1, "GREEDY_PHRASE": 2}

func (ret CommandNode) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Flags)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Children)).Encode(w)
	if err != nil {
		return
	}
	for iCommandNodeChildren := range len(ret.Children) {
		err = ret.Children[iCommandNodeChildren].Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Flags {
	case 1:
		CommandNodeRedirectNode, ok := ret.RedirectNode.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = CommandNodeRedirectNode.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.RedirectNode.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.RedirectNode.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Flags {
	case 0:
		CommandNodeExtraNodeData, ok := ret.ExtraNodeData.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = CommandNodeExtraNodeData.Encode(w)
		if err != nil {
			return
		}
	case 1:
		CommandNodeExtraNodeData, ok := ret.ExtraNodeData.(struct {
			Name string
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.EncodeString(w, CommandNodeExtraNodeData.Name)
		if err != nil {
			return
		}
	case 2:
		CommandNodeExtraNodeData, ok := ret.ExtraNodeData.(struct {
			Name           string
			Parser         string
			Properties     any
			SuggestionType any
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.EncodeString(w, CommandNodeExtraNodeData.Name)
		if err != nil {
			return
		}
		var vCommandNodeExtraNodeDataParser queser.VarInt
		vCommandNodeExtraNodeDataParser, err = queser.ErroringIndex(CommandNodeExtraNodeDataParserReverseMap, CommandNodeExtraNodeData.Parser)
		if err != nil {
			return
		}
		err = vCommandNodeExtraNodeDataParser.Encode(w)
		if err != nil {
			return
		}
		switch CommandNodeExtraNodeData.Parser {
		case "brigadier:bool":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "brigadier:double":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(struct {
				Flags uint8
				Min   any
				Max   any
			})
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataProperties.Flags)
			if err != nil {
				return
			}
			switch CommandNodeExtraNodeDataProperties.Flags {
			case 1:
				CommandNodeExtraNodeDataPropertiesMin, ok := CommandNodeExtraNodeDataProperties.Min.(float64)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataPropertiesMin)
				if err != nil {
					return
				}
			default:
				_, ok := CommandNodeExtraNodeDataProperties.Min.(queser.Void)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = CommandNodeExtraNodeDataProperties.Min.(queser.Void).Encode(w)
				if err != nil {
					return
				}
			}
			switch CommandNodeExtraNodeDataProperties.Flags {
			case 1:
				CommandNodeExtraNodeDataPropertiesMax, ok := CommandNodeExtraNodeDataProperties.Max.(float64)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataPropertiesMax)
				if err != nil {
					return
				}
			default:
				_, ok := CommandNodeExtraNodeDataProperties.Max.(queser.Void)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = CommandNodeExtraNodeDataProperties.Max.(queser.Void).Encode(w)
				if err != nil {
					return
				}
			}
		case "brigadier:float":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(struct {
				Flags uint8
				Min   any
				Max   any
			})
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataProperties.Flags)
			if err != nil {
				return
			}
			switch CommandNodeExtraNodeDataProperties.Flags {
			case 1:
				CommandNodeExtraNodeDataPropertiesMin, ok := CommandNodeExtraNodeDataProperties.Min.(float32)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataPropertiesMin)
				if err != nil {
					return
				}
			default:
				_, ok := CommandNodeExtraNodeDataProperties.Min.(queser.Void)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = CommandNodeExtraNodeDataProperties.Min.(queser.Void).Encode(w)
				if err != nil {
					return
				}
			}
			switch CommandNodeExtraNodeDataProperties.Flags {
			case 1:
				CommandNodeExtraNodeDataPropertiesMax, ok := CommandNodeExtraNodeDataProperties.Max.(float32)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataPropertiesMax)
				if err != nil {
					return
				}
			default:
				_, ok := CommandNodeExtraNodeDataProperties.Max.(queser.Void)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = CommandNodeExtraNodeDataProperties.Max.(queser.Void).Encode(w)
				if err != nil {
					return
				}
			}
		case "brigadier:integer":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(struct {
				Flags uint8
				Min   any
				Max   any
			})
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataProperties.Flags)
			if err != nil {
				return
			}
			switch CommandNodeExtraNodeDataProperties.Flags {
			case 1:
				CommandNodeExtraNodeDataPropertiesMin, ok := CommandNodeExtraNodeDataProperties.Min.(int32)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataPropertiesMin)
				if err != nil {
					return
				}
			default:
				_, ok := CommandNodeExtraNodeDataProperties.Min.(queser.Void)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = CommandNodeExtraNodeDataProperties.Min.(queser.Void).Encode(w)
				if err != nil {
					return
				}
			}
			switch CommandNodeExtraNodeDataProperties.Flags {
			case 1:
				CommandNodeExtraNodeDataPropertiesMax, ok := CommandNodeExtraNodeDataProperties.Max.(int32)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataPropertiesMax)
				if err != nil {
					return
				}
			default:
				_, ok := CommandNodeExtraNodeDataProperties.Max.(queser.Void)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = CommandNodeExtraNodeDataProperties.Max.(queser.Void).Encode(w)
				if err != nil {
					return
				}
			}
		case "brigadier:long":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(struct {
				Flags uint8
				Min   any
				Max   any
			})
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataProperties.Flags)
			if err != nil {
				return
			}
			switch CommandNodeExtraNodeDataProperties.Flags {
			case 1:
				CommandNodeExtraNodeDataPropertiesMin, ok := CommandNodeExtraNodeDataProperties.Min.(int64)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataPropertiesMin)
				if err != nil {
					return
				}
			default:
				_, ok := CommandNodeExtraNodeDataProperties.Min.(queser.Void)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = CommandNodeExtraNodeDataProperties.Min.(queser.Void).Encode(w)
				if err != nil {
					return
				}
			}
			switch CommandNodeExtraNodeDataProperties.Flags {
			case 1:
				CommandNodeExtraNodeDataPropertiesMax, ok := CommandNodeExtraNodeDataProperties.Max.(int64)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataPropertiesMax)
				if err != nil {
					return
				}
			default:
				_, ok := CommandNodeExtraNodeDataProperties.Max.(queser.Void)
				if !ok {
					err = queser.BadTypeError
					return
				}
				err = CommandNodeExtraNodeDataProperties.Max.(queser.Void).Encode(w)
				if err != nil {
					return
				}
			}
		case "brigadier:string":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(string)
			if !ok {
				err = queser.BadTypeError
				return
			}
			var vCommandNodeExtraNodeDataProperties queser.VarInt
			vCommandNodeExtraNodeDataProperties, err = queser.ErroringIndex(CommandNodeExtraNodeDataPropertiesReverseMap, CommandNodeExtraNodeDataProperties)
			if err != nil {
				return
			}
			err = vCommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:angle":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:block_pos":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:block_predicate":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:block_state":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:color":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:column_pos":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:component":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:dialog":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:dimension":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:entity":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(uint8)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataProperties)
			if err != nil {
				return
			}
		case "minecraft:entity_anchor":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:float_range":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:function":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:game_profile":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:gamemode":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:heightmap":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:hex_color":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:int_range":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:item_predicate":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:item_slot":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:item_stack":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:message":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:nbt":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:nbt_path":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:objective":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:objective_criteria":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:operation":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:particle":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:resource":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(struct {
				Registry string
			})
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = queser.EncodeString(w, CommandNodeExtraNodeDataProperties.Registry)
			if err != nil {
				return
			}
		case "minecraft:resource_key":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(struct {
				Registry string
			})
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = queser.EncodeString(w, CommandNodeExtraNodeDataProperties.Registry)
			if err != nil {
				return
			}
		case "minecraft:resource_location":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:resource_or_tag":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(struct {
				Registry string
			})
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = queser.EncodeString(w, CommandNodeExtraNodeDataProperties.Registry)
			if err != nil {
				return
			}
		case "minecraft:resource_or_tag_key":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(struct {
				Registry string
			})
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = queser.EncodeString(w, CommandNodeExtraNodeDataProperties.Registry)
			if err != nil {
				return
			}
		case "minecraft:resource_selector":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(struct {
				Registry string
			})
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = queser.EncodeString(w, CommandNodeExtraNodeDataProperties.Registry)
			if err != nil {
				return
			}
		case "minecraft:rotation":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:score_holder":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(uint8)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataProperties)
			if err != nil {
				return
			}
		case "minecraft:scoreboard_slot":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:swizzle":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:team":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:template_mirror":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:template_rotation":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:time":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(struct {
				Min int32
			})
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = binary.Write(w, binary.BigEndian, CommandNodeExtraNodeDataProperties.Min)
			if err != nil {
				return
			}
		case "minecraft:uuid":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:vec2":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		case "minecraft:vec3":
			CommandNodeExtraNodeDataProperties, ok := CommandNodeExtraNodeData.Properties.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeDataProperties.Encode(w)
			if err != nil {
				return
			}
		}
		switch ret.Flags {
		case 1:
			CommandNodeExtraNodeDataSuggestionType, ok := CommandNodeExtraNodeData.SuggestionType.(string)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = queser.EncodeString(w, CommandNodeExtraNodeDataSuggestionType)
			if err != nil {
				return
			}
		default:
			_, ok := CommandNodeExtraNodeData.SuggestionType.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = CommandNodeExtraNodeData.SuggestionType.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
	}
	return
}

type EntityMetadata struct {
	Val queser.ToDo
}

func (_ EntityMetadata) Decode(r io.Reader) (ret EntityMetadata, err error) {
	err = queser.ToDoError
	return
}
func (ret EntityMetadata) Encode(w io.Writer) (err error) {
	err = queser.ToDoError
	return
}

type EntityMetadataEntry struct {
	Key   uint8
	Type  string
	Value any
}

var EntityMetadataEntryTypeMap = map[queser.VarInt]string{0: "byte", 1: "int", 10: "block_pos", 11: "optional_block_pos", 12: "direction", 13: "optional_uuid", 14: "block_state", 15: "optional_block_state", 16: "compound_tag", 17: "particle", 18: "particles", 19: "villager_data", 2: "long", 20: "optional_unsigned_int", 21: "pose", 22: "cat_variant", 23: "cow_variant", 24: "wolf_variant", 25: "wolf_sound_variant", 26: "frog_variant", 27: "pig_variant", 28: "chicken_variant", 29: "optional_global_pos", 3: "float", 30: "painting_variant", 31: "sniffer_state", 32: "armadillo_state", 33: "vector3", 34: "quaternion", 4: "string", 5: "component", 6: "optional_component", 7: "item_stack", 8: "boolean", 9: "rotations"}

func (_ EntityMetadataEntry) Decode(r io.Reader) (ret EntityMetadataEntry, err error) {
	err = queser.ToDoError
	return
}

var EntityMetadataEntryTypeReverseMap = map[string]queser.VarInt{"byte": 0, "int": 1, "block_pos": 10, "optional_block_pos": 11, "direction": 12, "optional_uuid": 13, "block_state": 14, "optional_block_state": 15, "compound_tag": 16, "particle": 17, "particles": 18, "villager_data": 19, "long": 2, "optional_unsigned_int": 20, "pose": 21, "cat_variant": 22, "cow_variant": 23, "wolf_variant": 24, "wolf_sound_variant": 25, "frog_variant": 26, "pig_variant": 27, "chicken_variant": 28, "optional_global_pos": 29, "float": 3, "painting_variant": 30, "sniffer_state": 31, "armadillo_state": 32, "vector3": 33, "quaternion": 34, "string": 4, "component": 5, "optional_component": 6, "item_stack": 7, "boolean": 8, "rotations": 9}

func (ret EntityMetadataEntry) Encode(w io.Writer) (err error) {
	err = queser.ToDoError
	return
}

type GameProfile struct {
	Name       string
	Properties []struct {
		Name      string
		Value     string
		Signature *string
	}
}

func (_ GameProfile) Decode(r io.Reader) (ret GameProfile, err error) {
	ret.Name, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var lGameProfileProperties queser.VarInt
	lGameProfileProperties, err = lGameProfileProperties.Decode(r)
	if err != nil {
		return
	}
	ret.Properties = []struct {
		Name      string
		Value     string
		Signature *string
	}{}
	for range lGameProfileProperties {
		var GameProfilePropertiesElement struct {
			Name      string
			Value     string
			Signature *string
		}
		GameProfilePropertiesElement.Name, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		GameProfilePropertiesElement.Value, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		var GameProfilePropertiesElementSignaturePresent bool
		err = binary.Read(r, binary.BigEndian, &GameProfilePropertiesElementSignaturePresent)
		if err != nil {
			return
		}
		if GameProfilePropertiesElementSignaturePresent {
			var GameProfilePropertiesElementSignaturePresentValue string
			GameProfilePropertiesElementSignaturePresentValue, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			GameProfilePropertiesElement.Signature = &GameProfilePropertiesElementSignaturePresentValue
		}
		ret.Properties = append(ret.Properties, GameProfilePropertiesElement)
	}
	return
}
func (ret GameProfile) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Name)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Properties)).Encode(w)
	if err != nil {
		return
	}
	for iGameProfileProperties := range len(ret.Properties) {
		err = queser.EncodeString(w, ret.Properties[iGameProfileProperties].Name)
		if err != nil {
			return
		}
		err = queser.EncodeString(w, ret.Properties[iGameProfileProperties].Value)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Properties[iGameProfileProperties].Signature != nil)
		if err != nil {
			return
		}
		if ret.Properties[iGameProfileProperties].Signature != nil {
			err = queser.EncodeString(w, *ret.Properties[iGameProfileProperties].Signature)
			if err != nil {
				return
			}
		}
	}
	return
}

type Ingredient struct {
	Val []Slot
}

func (_ Ingredient) Decode(r io.Reader) (ret Ingredient, err error) {
	var lIngredient queser.VarInt
	lIngredient, err = lIngredient.Decode(r)
	if err != nil {
		return
	}
	ret.Val = []Slot{}
	for range lIngredient {
		var IngredientElement Slot
		IngredientElement, err = IngredientElement.Decode(r)
		if err != nil {
			return
		}
		ret.Val = append(ret.Val, IngredientElement)
	}
	return
}
func (ret Ingredient) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Val)).Encode(w)
	if err != nil {
		return
	}
	for iIngredient := range len(ret.Val) {
		err = ret.Val[iIngredient].Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type Optvarint struct {
	Val queser.VarInt
}

func (_ Optvarint) Decode(r io.Reader) (ret Optvarint, err error) {
	ret.Val, err = ret.Val.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret Optvarint) Encode(w io.Writer) (err error) {
	err = ret.Val.Encode(w)
	if err != nil {
		return
	}
	return
}

type PackedChunkPos struct {
	Z int32
	X int32
}

func (_ PackedChunkPos) Decode(r io.Reader) (ret PackedChunkPos, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	return
}
func (ret PackedChunkPos) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	return
}

type PacketCommonAddResourcePack struct {
	Uuid          uuid.UUID
	Url           string
	Hash          string
	Forced        bool
	PromptMessage *nbt.Anon
}

func (_ PacketCommonAddResourcePack) Decode(r io.Reader) (ret PacketCommonAddResourcePack, err error) {
	_, err = io.ReadFull(r, ret.Uuid[:])
	if err != nil {
		return
	}
	ret.Url, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Hash, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Forced)
	if err != nil {
		return
	}
	var PacketCommonAddResourcePackPromptMessagePresent bool
	err = binary.Read(r, binary.BigEndian, &PacketCommonAddResourcePackPromptMessagePresent)
	if err != nil {
		return
	}
	if PacketCommonAddResourcePackPromptMessagePresent {
		var PacketCommonAddResourcePackPromptMessagePresentValue nbt.Anon
		PacketCommonAddResourcePackPromptMessagePresentValue, err = PacketCommonAddResourcePackPromptMessagePresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.PromptMessage = &PacketCommonAddResourcePackPromptMessagePresentValue
	}
	return
}
func (ret PacketCommonAddResourcePack) Encode(w io.Writer) (err error) {
	_, err = w.Write(ret.Uuid[:])
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Url)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Hash)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Forced)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.PromptMessage != nil)
	if err != nil {
		return
	}
	if ret.PromptMessage != nil {
		err = (*ret.PromptMessage).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PacketCommonClearDialog struct {
}

func (_ PacketCommonClearDialog) Decode(r io.Reader) (ret PacketCommonClearDialog, err error) {
	return
}
func (ret PacketCommonClearDialog) Encode(w io.Writer) (err error) {
	return
}

type PacketCommonCookieRequest struct {
	Cookie string
}

func (_ PacketCommonCookieRequest) Decode(r io.Reader) (ret PacketCommonCookieRequest, err error) {
	ret.Cookie, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	return
}
func (ret PacketCommonCookieRequest) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Cookie)
	if err != nil {
		return
	}
	return
}

type PacketCommonCookieResponse struct {
	Key   string
	Value *ByteArray
}

func (_ PacketCommonCookieResponse) Decode(r io.Reader) (ret PacketCommonCookieResponse, err error) {
	ret.Key, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var PacketCommonCookieResponseValuePresent bool
	err = binary.Read(r, binary.BigEndian, &PacketCommonCookieResponseValuePresent)
	if err != nil {
		return
	}
	if PacketCommonCookieResponseValuePresent {
		var PacketCommonCookieResponseValuePresentValue ByteArray
		PacketCommonCookieResponseValuePresentValue, err = PacketCommonCookieResponseValuePresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.Value = &PacketCommonCookieResponseValuePresentValue
	}
	return
}
func (ret PacketCommonCookieResponse) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Key)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Value != nil)
	if err != nil {
		return
	}
	if ret.Value != nil {
		err = (*ret.Value).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PacketCommonCustomClickAction struct {
	Id  string
	Nbt *nbt.Anon
}

func (_ PacketCommonCustomClickAction) Decode(r io.Reader) (ret PacketCommonCustomClickAction, err error) {
	ret.Id, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var PacketCommonCustomClickActionNbtPresent bool
	err = binary.Read(r, binary.BigEndian, &PacketCommonCustomClickActionNbtPresent)
	if err != nil {
		return
	}
	if PacketCommonCustomClickActionNbtPresent {
		var PacketCommonCustomClickActionNbtPresentValue nbt.Anon
		PacketCommonCustomClickActionNbtPresentValue, err = PacketCommonCustomClickActionNbtPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.Nbt = &PacketCommonCustomClickActionNbtPresentValue
	}
	return
}
func (ret PacketCommonCustomClickAction) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Id)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Nbt != nil)
	if err != nil {
		return
	}
	if ret.Nbt != nil {
		err = (*ret.Nbt).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PacketCommonCustomReportDetails struct {
	Details []struct {
		Key   string
		Value string
	}
}

func (_ PacketCommonCustomReportDetails) Decode(r io.Reader) (ret PacketCommonCustomReportDetails, err error) {
	var lPacketCommonCustomReportDetailsDetails queser.VarInt
	lPacketCommonCustomReportDetailsDetails, err = lPacketCommonCustomReportDetailsDetails.Decode(r)
	if err != nil {
		return
	}
	ret.Details = []struct {
		Key   string
		Value string
	}{}
	for range lPacketCommonCustomReportDetailsDetails {
		var PacketCommonCustomReportDetailsDetailsElement struct {
			Key   string
			Value string
		}
		PacketCommonCustomReportDetailsDetailsElement.Key, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		PacketCommonCustomReportDetailsDetailsElement.Value, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Details = append(ret.Details, PacketCommonCustomReportDetailsDetailsElement)
	}
	return
}
func (ret PacketCommonCustomReportDetails) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Details)).Encode(w)
	if err != nil {
		return
	}
	for iPacketCommonCustomReportDetailsDetails := range len(ret.Details) {
		err = queser.EncodeString(w, ret.Details[iPacketCommonCustomReportDetailsDetails].Key)
		if err != nil {
			return
		}
		err = queser.EncodeString(w, ret.Details[iPacketCommonCustomReportDetailsDetails].Value)
		if err != nil {
			return
		}
	}
	return
}

type PacketCommonRemoveResourcePack struct {
	Uuid *uuid.UUID
}

func (_ PacketCommonRemoveResourcePack) Decode(r io.Reader) (ret PacketCommonRemoveResourcePack, err error) {
	var PacketCommonRemoveResourcePackUuidPresent bool
	err = binary.Read(r, binary.BigEndian, &PacketCommonRemoveResourcePackUuidPresent)
	if err != nil {
		return
	}
	if PacketCommonRemoveResourcePackUuidPresent {
		var PacketCommonRemoveResourcePackUuidPresentValue uuid.UUID
		_, err = io.ReadFull(r, PacketCommonRemoveResourcePackUuidPresentValue[:])
		if err != nil {
			return
		}
		ret.Uuid = &PacketCommonRemoveResourcePackUuidPresentValue
	}
	return
}
func (ret PacketCommonRemoveResourcePack) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Uuid != nil)
	if err != nil {
		return
	}
	if ret.Uuid != nil {
		_, err = w.Write((*ret.Uuid)[:])
		if err != nil {
			return
		}
	}
	return
}

type PacketCommonSelectKnownPacks struct {
	Packs []struct {
		Namespace string
		Id        string
		Version   string
	}
}

func (_ PacketCommonSelectKnownPacks) Decode(r io.Reader) (ret PacketCommonSelectKnownPacks, err error) {
	var lPacketCommonSelectKnownPacksPacks queser.VarInt
	lPacketCommonSelectKnownPacksPacks, err = lPacketCommonSelectKnownPacksPacks.Decode(r)
	if err != nil {
		return
	}
	ret.Packs = []struct {
		Namespace string
		Id        string
		Version   string
	}{}
	for range lPacketCommonSelectKnownPacksPacks {
		var PacketCommonSelectKnownPacksPacksElement struct {
			Namespace string
			Id        string
			Version   string
		}
		PacketCommonSelectKnownPacksPacksElement.Namespace, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		PacketCommonSelectKnownPacksPacksElement.Id, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		PacketCommonSelectKnownPacksPacksElement.Version, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Packs = append(ret.Packs, PacketCommonSelectKnownPacksPacksElement)
	}
	return
}
func (ret PacketCommonSelectKnownPacks) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Packs)).Encode(w)
	if err != nil {
		return
	}
	for iPacketCommonSelectKnownPacksPacks := range len(ret.Packs) {
		err = queser.EncodeString(w, ret.Packs[iPacketCommonSelectKnownPacksPacks].Namespace)
		if err != nil {
			return
		}
		err = queser.EncodeString(w, ret.Packs[iPacketCommonSelectKnownPacksPacks].Id)
		if err != nil {
			return
		}
		err = queser.EncodeString(w, ret.Packs[iPacketCommonSelectKnownPacksPacks].Version)
		if err != nil {
			return
		}
	}
	return
}

type PacketCommonServerLinks struct {
	Links []struct {
		HasKnownType bool
		KnownType    any
		UnknownType  any
		Link         string
	}
}

func (_ PacketCommonServerLinks) Decode(r io.Reader) (ret PacketCommonServerLinks, err error) {
	var lPacketCommonServerLinksLinks queser.VarInt
	lPacketCommonServerLinksLinks, err = lPacketCommonServerLinksLinks.Decode(r)
	if err != nil {
		return
	}
	ret.Links = []struct {
		HasKnownType bool
		KnownType    any
		UnknownType  any
		Link         string
	}{}
	for range lPacketCommonServerLinksLinks {
		var PacketCommonServerLinksLinksElement struct {
			HasKnownType bool
			KnownType    any
			UnknownType  any
			Link         string
		}
		err = binary.Read(r, binary.BigEndian, &PacketCommonServerLinksLinksElement.HasKnownType)
		if err != nil {
			return
		}
		switch PacketCommonServerLinksLinksElement.HasKnownType {
		case true:
			var PacketCommonServerLinksLinksElementKnownTypeTmp ServerLinkType
			PacketCommonServerLinksLinksElementKnownTypeTmp, err = PacketCommonServerLinksLinksElementKnownTypeTmp.Decode(r)
			if err != nil {
				return
			}
			PacketCommonServerLinksLinksElement.KnownType = PacketCommonServerLinksLinksElementKnownTypeTmp
		}
		switch PacketCommonServerLinksLinksElement.HasKnownType {
		case false:
			var PacketCommonServerLinksLinksElementUnknownTypeTmp nbt.Anon
			PacketCommonServerLinksLinksElementUnknownTypeTmp, err = PacketCommonServerLinksLinksElementUnknownTypeTmp.Decode(r)
			if err != nil {
				return
			}
			PacketCommonServerLinksLinksElement.UnknownType = PacketCommonServerLinksLinksElementUnknownTypeTmp
		}
		PacketCommonServerLinksLinksElement.Link, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Links = append(ret.Links, PacketCommonServerLinksLinksElement)
	}
	return
}
func (ret PacketCommonServerLinks) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Links)).Encode(w)
	if err != nil {
		return
	}
	for iPacketCommonServerLinksLinks := range len(ret.Links) {
		err = binary.Write(w, binary.BigEndian, ret.Links[iPacketCommonServerLinksLinks].HasKnownType)
		if err != nil {
			return
		}
		switch ret.Links[iPacketCommonServerLinksLinks].HasKnownType {
		case true:
			PacketCommonServerLinksLinksInnerKnownType, ok := ret.Links[iPacketCommonServerLinksLinks].KnownType.(ServerLinkType)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PacketCommonServerLinksLinksInnerKnownType.Encode(w)
			if err != nil {
				return
			}
		}
		switch ret.Links[iPacketCommonServerLinksLinks].HasKnownType {
		case false:
			PacketCommonServerLinksLinksInnerUnknownType, ok := ret.Links[iPacketCommonServerLinksLinks].UnknownType.(nbt.Anon)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PacketCommonServerLinksLinksInnerUnknownType.Encode(w)
			if err != nil {
				return
			}
		}
		err = queser.EncodeString(w, ret.Links[iPacketCommonServerLinksLinks].Link)
		if err != nil {
			return
		}
	}
	return
}

type PacketCommonSettings struct {
	Locale              string
	ViewDistance        int8
	ChatFlags           queser.VarInt
	ChatColors          bool
	SkinParts           uint8
	MainHand            queser.VarInt
	EnableTextFiltering bool
	EnableServerListing bool
	ParticleStatus      string
}

var PacketCommonSettingsParticleStatusMap = map[queser.VarInt]string{0: "all", 1: "decreased", 2: "minimal"}

func (_ PacketCommonSettings) Decode(r io.Reader) (ret PacketCommonSettings, err error) {
	ret.Locale, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.ViewDistance)
	if err != nil {
		return
	}
	ret.ChatFlags, err = ret.ChatFlags.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.ChatColors)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.SkinParts)
	if err != nil {
		return
	}
	ret.MainHand, err = ret.MainHand.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.EnableTextFiltering)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.EnableServerListing)
	if err != nil {
		return
	}
	var PacketCommonSettingsParticleStatusKey queser.VarInt
	PacketCommonSettingsParticleStatusKey, err = PacketCommonSettingsParticleStatusKey.Decode(r)
	if err != nil {
		return
	}
	ret.ParticleStatus, err = queser.ErroringIndex(PacketCommonSettingsParticleStatusMap, PacketCommonSettingsParticleStatusKey)
	if err != nil {
		return
	}
	return
}

var PacketCommonSettingsParticleStatusReverseMap = map[string]queser.VarInt{"all": 0, "decreased": 1, "minimal": 2}

func (ret PacketCommonSettings) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Locale)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.ViewDistance)
	if err != nil {
		return
	}
	err = ret.ChatFlags.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.ChatColors)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.SkinParts)
	if err != nil {
		return
	}
	err = ret.MainHand.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.EnableTextFiltering)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.EnableServerListing)
	if err != nil {
		return
	}
	var vPacketCommonSettingsParticleStatus queser.VarInt
	vPacketCommonSettingsParticleStatus, err = queser.ErroringIndex(PacketCommonSettingsParticleStatusReverseMap, ret.ParticleStatus)
	if err != nil {
		return
	}
	err = vPacketCommonSettingsParticleStatus.Encode(w)
	if err != nil {
		return
	}
	return
}

type PacketCommonStoreCookie struct {
	Key   string
	Value ByteArray
}

func (_ PacketCommonStoreCookie) Decode(r io.Reader) (ret PacketCommonStoreCookie, err error) {
	ret.Key, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Value, err = ret.Value.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PacketCommonStoreCookie) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Key)
	if err != nil {
		return
	}
	err = ret.Value.Encode(w)
	if err != nil {
		return
	}
	return
}

type PacketCommonTransfer struct {
	Host string
	Port queser.VarInt
}

func (_ PacketCommonTransfer) Decode(r io.Reader) (ret PacketCommonTransfer, err error) {
	ret.Host, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Port, err = ret.Port.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PacketCommonTransfer) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Host)
	if err != nil {
		return
	}
	err = ret.Port.Encode(w)
	if err != nil {
		return
	}
	return
}

type Position struct {
	Val uint64
}

func (_ Position) Decode(r io.Reader) (ret Position, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Val)
	if err != nil {
		return
	}
	return
}
func (ret Position) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Val)
	if err != nil {
		return
	}
	return
}

type PreviousMessages struct {
	Val []struct {
		Id        queser.VarInt
		Signature any
	}
}

func (_ PreviousMessages) Decode(r io.Reader) (ret PreviousMessages, err error) {
	var lPreviousMessages queser.VarInt
	lPreviousMessages, err = lPreviousMessages.Decode(r)
	if err != nil {
		return
	}
	ret.Val = []struct {
		Id        queser.VarInt
		Signature any
	}{}
	for range lPreviousMessages {
		var PreviousMessagesElement struct {
			Id        queser.VarInt
			Signature any
		}
		PreviousMessagesElement.Id, err = PreviousMessagesElement.Id.Decode(r)
		if err != nil {
			return
		}
		switch PreviousMessagesElement.Id {
		case 0:
			var PreviousMessagesElementSignatureTmp [256]byte
			_, err = r.Read(PreviousMessagesElementSignatureTmp[:])
			if err != nil {
				return
			}
			PreviousMessagesElement.Signature = PreviousMessagesElementSignatureTmp
		default:
			var PreviousMessagesElementSignatureTmp queser.Void
			PreviousMessagesElementSignatureTmp, err = PreviousMessagesElementSignatureTmp.Decode(r)
			if err != nil {
				return
			}
			PreviousMessagesElement.Signature = PreviousMessagesElementSignatureTmp
		}
		ret.Val = append(ret.Val, PreviousMessagesElement)
	}
	return
}
func (ret PreviousMessages) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Val)).Encode(w)
	if err != nil {
		return
	}
	for iPreviousMessages := range len(ret.Val) {
		err = ret.Val[iPreviousMessages].Id.Encode(w)
		if err != nil {
			return
		}
		switch ret.Val[iPreviousMessages].Id {
		case 0:
			PreviousMessagesInnerSignature, ok := ret.Val[iPreviousMessages].Signature.([256]byte)
			if !ok {
				err = queser.BadTypeError
				return
			}
			arr := PreviousMessagesInnerSignature
			_, err = w.Write(arr[:])
			if err != nil {
				return
			}
		default:
			_, ok := ret.Val[iPreviousMessages].Signature.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = ret.Val[iPreviousMessages].Signature.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
	}
	return
}

type SoundSource struct {
	Val string
}

var SoundSourceMap = map[queser.VarInt]string{0: "master", 1: "music", 10: "ui", 2: "record", 3: "weather", 4: "block", 5: "hostile", 6: "neutral", 7: "player", 8: "ambient", 9: "voice"}

func (_ SoundSource) Decode(r io.Reader) (ret SoundSource, err error) {
	var SoundSourceKey queser.VarInt
	SoundSourceKey, err = SoundSourceKey.Decode(r)
	if err != nil {
		return
	}
	ret.Val, err = queser.ErroringIndex(SoundSourceMap, SoundSourceKey)
	if err != nil {
		return
	}
	return
}

var SoundSourceReverseMap = map[string]queser.VarInt{"master": 0, "music": 1, "ui": 10, "record": 2, "weather": 3, "block": 4, "hostile": 5, "neutral": 6, "player": 7, "ambient": 8, "voice": 9}

func (ret SoundSource) Encode(w io.Writer) (err error) {
	var vSoundSource queser.VarInt
	vSoundSource, err = queser.ErroringIndex(SoundSourceReverseMap, ret.Val)
	if err != nil {
		return
	}
	err = vSoundSource.Encode(w)
	if err != nil {
		return
	}
	return
}

type Tags struct {
	Val []struct {
		TagName string
		Entries []queser.VarInt
	}
}

func (_ Tags) Decode(r io.Reader) (ret Tags, err error) {
	var lTags queser.VarInt
	lTags, err = lTags.Decode(r)
	if err != nil {
		return
	}
	ret.Val = []struct {
		TagName string
		Entries []queser.VarInt
	}{}
	for range lTags {
		var TagsElement struct {
			TagName string
			Entries []queser.VarInt
		}
		TagsElement.TagName, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		var lTagsElementEntries queser.VarInt
		lTagsElementEntries, err = lTagsElementEntries.Decode(r)
		if err != nil {
			return
		}
		TagsElement.Entries = []queser.VarInt{}
		for range lTagsElementEntries {
			var TagsElementEntriesElement queser.VarInt
			TagsElementEntriesElement, err = TagsElementEntriesElement.Decode(r)
			if err != nil {
				return
			}
			TagsElement.Entries = append(TagsElement.Entries, TagsElementEntriesElement)
		}
		ret.Val = append(ret.Val, TagsElement)
	}
	return
}
func (ret Tags) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Val)).Encode(w)
	if err != nil {
		return
	}
	for iTags := range len(ret.Val) {
		err = queser.EncodeString(w, ret.Val[iTags].TagName)
		if err != nil {
			return
		}
		err = queser.VarInt(len(ret.Val[iTags].Entries)).Encode(w)
		if err != nil {
			return
		}
		for iTagsInnerEntries := range len(ret.Val[iTags].Entries) {
			err = ret.Val[iTags].Entries[iTagsInnerEntries].Encode(w)
			if err != nil {
				return
			}
		}
	}
	return
}

type Vec2f struct {
	X float32
	Y float32
}

func (_ Vec2f) Decode(r io.Reader) (ret Vec2f, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	return
}
func (ret Vec2f) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	return
}

type Vec3f struct {
	X float32
	Y float32
	Z float32
}

func (_ Vec3f) Decode(r io.Reader) (ret Vec3f, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	return
}
func (ret Vec3f) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	return
}

type Vec3f64 struct {
	X float64
	Y float64
	Z float64
}

func (_ Vec3f64) Decode(r io.Reader) (ret Vec3f64, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	return
}
func (ret Vec3f64) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	return
}

type Vec3i struct {
	X queser.VarInt
	Y queser.VarInt
	Z queser.VarInt
}

func (_ Vec3i) Decode(r io.Reader) (ret Vec3i, err error) {
	ret.X, err = ret.X.Decode(r)
	if err != nil {
		return
	}
	ret.Y, err = ret.Y.Decode(r)
	if err != nil {
		return
	}
	ret.Z, err = ret.Z.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret Vec3i) Encode(w io.Writer) (err error) {
	err = ret.X.Encode(w)
	if err != nil {
		return
	}
	err = ret.Y.Encode(w)
	if err != nil {
		return
	}
	err = ret.Z.Encode(w)
	if err != nil {
		return
	}
	return
}

type Vec4f struct {
	X float32
	Y float32
	Z float32
	W float32
}

func (_ Vec4f) Decode(r io.Reader) (ret Vec4f, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.W)
	if err != nil {
		return
	}
	return
}
func (ret Vec4f) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.W)
	if err != nil {
		return
	}
	return
}

type HandshakingToServerPacket struct {
	Name   string
	Params any
}

var HandshakingToServerPacketNameMap = map[queser.VarInt]string{0x00: "set_protocol", 0xfe: "legacy_server_list_ping"}

func (_ HandshakingToServerPacket) Decode(r io.Reader) (ret HandshakingToServerPacket, err error) {
	var HandshakingToServerPacketNameKey queser.VarInt
	HandshakingToServerPacketNameKey, err = HandshakingToServerPacketNameKey.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.ErroringIndex(HandshakingToServerPacketNameMap, HandshakingToServerPacketNameKey)
	if err != nil {
		return
	}
	switch ret.Name {
	case "legacy_server_list_ping":
		var HandshakingToServerPacketParamsTmp HandshakingToServerPacketLegacyServerListPing
		HandshakingToServerPacketParamsTmp, err = HandshakingToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = HandshakingToServerPacketParamsTmp
	case "set_protocol":
		var HandshakingToServerPacketParamsTmp HandshakingToServerPacketSetProtocol
		HandshakingToServerPacketParamsTmp, err = HandshakingToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = HandshakingToServerPacketParamsTmp
	}
	return
}

var HandshakingToServerPacketNameReverseMap = map[string]queser.VarInt{"set_protocol": 0x00, "legacy_server_list_ping": 0xfe}

func (ret HandshakingToServerPacket) Encode(w io.Writer) (err error) {
	var vHandshakingToServerPacketName queser.VarInt
	vHandshakingToServerPacketName, err = queser.ErroringIndex(HandshakingToServerPacketNameReverseMap, ret.Name)
	if err != nil {
		return
	}
	err = vHandshakingToServerPacketName.Encode(w)
	if err != nil {
		return
	}
	switch ret.Name {
	case "legacy_server_list_ping":
		HandshakingToServerPacketParams, ok := ret.Params.(HandshakingToServerPacketLegacyServerListPing)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = HandshakingToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_protocol":
		HandshakingToServerPacketParams, ok := ret.Params.(HandshakingToServerPacketSetProtocol)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = HandshakingToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type HandshakingToServerPacketLegacyServerListPing struct {
	Payload uint8
}

func (_ HandshakingToServerPacketLegacyServerListPing) Decode(r io.Reader) (ret HandshakingToServerPacketLegacyServerListPing, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Payload)
	if err != nil {
		return
	}
	return
}
func (ret HandshakingToServerPacketLegacyServerListPing) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Payload)
	if err != nil {
		return
	}
	return
}

type HandshakingToServerPacketSetProtocol struct {
	ProtocolVersion queser.VarInt
	ServerHost      string
	ServerPort      uint16
	NextState       queser.VarInt
}

func (_ HandshakingToServerPacketSetProtocol) Decode(r io.Reader) (ret HandshakingToServerPacketSetProtocol, err error) {
	ret.ProtocolVersion, err = ret.ProtocolVersion.Decode(r)
	if err != nil {
		return
	}
	ret.ServerHost, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.ServerPort)
	if err != nil {
		return
	}
	ret.NextState, err = ret.NextState.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret HandshakingToServerPacketSetProtocol) Encode(w io.Writer) (err error) {
	err = ret.ProtocolVersion.Encode(w)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.ServerHost)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.ServerPort)
	if err != nil {
		return
	}
	err = ret.NextState.Encode(w)
	if err != nil {
		return
	}
	return
}

type HandshakingToClientPacket struct {
	Name   string
	Params any
}

var HandshakingToClientPacketNameMap = map[queser.VarInt]string{}

func (_ HandshakingToClientPacket) Decode(r io.Reader) (ret HandshakingToClientPacket, err error) {
	var HandshakingToClientPacketNameKey queser.VarInt
	HandshakingToClientPacketNameKey, err = HandshakingToClientPacketNameKey.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.ErroringIndex(HandshakingToClientPacketNameMap, HandshakingToClientPacketNameKey)
	if err != nil {
		return
	}
	return
}

var HandshakingToClientPacketNameReverseMap = map[string]queser.VarInt{}

func (ret HandshakingToClientPacket) Encode(w io.Writer) (err error) {
	var vHandshakingToClientPacketName queser.VarInt
	vHandshakingToClientPacketName, err = queser.ErroringIndex(HandshakingToClientPacketNameReverseMap, ret.Name)
	if err != nil {
		return
	}
	err = vHandshakingToClientPacketName.Encode(w)
	if err != nil {
		return
	}
	return
}

type StatusToServerPacket struct {
	Name   string
	Params any
}

var StatusToServerPacketNameMap = map[queser.VarInt]string{0x00: "ping_start", 0x01: "ping"}

func (_ StatusToServerPacket) Decode(r io.Reader) (ret StatusToServerPacket, err error) {
	var StatusToServerPacketNameKey queser.VarInt
	StatusToServerPacketNameKey, err = StatusToServerPacketNameKey.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.ErroringIndex(StatusToServerPacketNameMap, StatusToServerPacketNameKey)
	if err != nil {
		return
	}
	switch ret.Name {
	case "ping":
		var StatusToServerPacketParamsTmp StatusToServerPacketPing
		StatusToServerPacketParamsTmp, err = StatusToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = StatusToServerPacketParamsTmp
	case "ping_start":
		var StatusToServerPacketParamsTmp StatusToServerPacketPingStart
		StatusToServerPacketParamsTmp, err = StatusToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = StatusToServerPacketParamsTmp
	}
	return
}

var StatusToServerPacketNameReverseMap = map[string]queser.VarInt{"ping_start": 0x00, "ping": 0x01}

func (ret StatusToServerPacket) Encode(w io.Writer) (err error) {
	var vStatusToServerPacketName queser.VarInt
	vStatusToServerPacketName, err = queser.ErroringIndex(StatusToServerPacketNameReverseMap, ret.Name)
	if err != nil {
		return
	}
	err = vStatusToServerPacketName.Encode(w)
	if err != nil {
		return
	}
	switch ret.Name {
	case "ping":
		StatusToServerPacketParams, ok := ret.Params.(StatusToServerPacketPing)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = StatusToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "ping_start":
		StatusToServerPacketParams, ok := ret.Params.(StatusToServerPacketPingStart)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = StatusToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type StatusToServerPacketPing struct {
	Time int64
}

func (_ StatusToServerPacketPing) Decode(r io.Reader) (ret StatusToServerPacketPing, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Time)
	if err != nil {
		return
	}
	return
}
func (ret StatusToServerPacketPing) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Time)
	if err != nil {
		return
	}
	return
}

type StatusToServerPacketPingStart struct {
}

func (_ StatusToServerPacketPingStart) Decode(r io.Reader) (ret StatusToServerPacketPingStart, err error) {
	return
}
func (ret StatusToServerPacketPingStart) Encode(w io.Writer) (err error) {
	return
}

type StatusToClientPacket struct {
	Name   string
	Params any
}

var StatusToClientPacketNameMap = map[queser.VarInt]string{0x00: "server_info", 0x01: "ping"}

func (_ StatusToClientPacket) Decode(r io.Reader) (ret StatusToClientPacket, err error) {
	var StatusToClientPacketNameKey queser.VarInt
	StatusToClientPacketNameKey, err = StatusToClientPacketNameKey.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.ErroringIndex(StatusToClientPacketNameMap, StatusToClientPacketNameKey)
	if err != nil {
		return
	}
	switch ret.Name {
	case "ping":
		var StatusToClientPacketParamsTmp StatusToClientPacketPing
		StatusToClientPacketParamsTmp, err = StatusToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = StatusToClientPacketParamsTmp
	case "server_info":
		var StatusToClientPacketParamsTmp StatusToClientPacketServerInfo
		StatusToClientPacketParamsTmp, err = StatusToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = StatusToClientPacketParamsTmp
	}
	return
}

var StatusToClientPacketNameReverseMap = map[string]queser.VarInt{"server_info": 0x00, "ping": 0x01}

func (ret StatusToClientPacket) Encode(w io.Writer) (err error) {
	var vStatusToClientPacketName queser.VarInt
	vStatusToClientPacketName, err = queser.ErroringIndex(StatusToClientPacketNameReverseMap, ret.Name)
	if err != nil {
		return
	}
	err = vStatusToClientPacketName.Encode(w)
	if err != nil {
		return
	}
	switch ret.Name {
	case "ping":
		StatusToClientPacketParams, ok := ret.Params.(StatusToClientPacketPing)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = StatusToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "server_info":
		StatusToClientPacketParams, ok := ret.Params.(StatusToClientPacketServerInfo)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = StatusToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type StatusToClientPacketPing struct {
	Time int64
}

func (_ StatusToClientPacketPing) Decode(r io.Reader) (ret StatusToClientPacketPing, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Time)
	if err != nil {
		return
	}
	return
}
func (ret StatusToClientPacketPing) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Time)
	if err != nil {
		return
	}
	return
}

type StatusToClientPacketServerInfo struct {
	Response string
}

func (_ StatusToClientPacketServerInfo) Decode(r io.Reader) (ret StatusToClientPacketServerInfo, err error) {
	ret.Response, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	return
}
func (ret StatusToClientPacketServerInfo) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Response)
	if err != nil {
		return
	}
	return
}

type LoginToServerPacket struct {
	Name   string
	Params any
}

var LoginToServerPacketNameMap = map[queser.VarInt]string{0x00: "login_start", 0x01: "encryption_begin", 0x02: "login_plugin_response", 0x03: "login_acknowledged", 0x04: "cookie_response"}

func (_ LoginToServerPacket) Decode(r io.Reader) (ret LoginToServerPacket, err error) {
	var LoginToServerPacketNameKey queser.VarInt
	LoginToServerPacketNameKey, err = LoginToServerPacketNameKey.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.ErroringIndex(LoginToServerPacketNameMap, LoginToServerPacketNameKey)
	if err != nil {
		return
	}
	switch ret.Name {
	case "cookie_response":
		var LoginToServerPacketParamsTmp PacketCommonCookieResponse
		LoginToServerPacketParamsTmp, err = LoginToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = LoginToServerPacketParamsTmp
	case "encryption_begin":
		var LoginToServerPacketParamsTmp LoginToServerPacketEncryptionBegin
		LoginToServerPacketParamsTmp, err = LoginToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = LoginToServerPacketParamsTmp
	case "login_acknowledged":
		var LoginToServerPacketParamsTmp LoginToServerPacketLoginAcknowledged
		LoginToServerPacketParamsTmp, err = LoginToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = LoginToServerPacketParamsTmp
	case "login_plugin_response":
		var LoginToServerPacketParamsTmp LoginToServerPacketLoginPluginResponse
		LoginToServerPacketParamsTmp, err = LoginToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = LoginToServerPacketParamsTmp
	case "login_start":
		var LoginToServerPacketParamsTmp LoginToServerPacketLoginStart
		LoginToServerPacketParamsTmp, err = LoginToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = LoginToServerPacketParamsTmp
	}
	return
}

var LoginToServerPacketNameReverseMap = map[string]queser.VarInt{"login_start": 0x00, "encryption_begin": 0x01, "login_plugin_response": 0x02, "login_acknowledged": 0x03, "cookie_response": 0x04}

func (ret LoginToServerPacket) Encode(w io.Writer) (err error) {
	var vLoginToServerPacketName queser.VarInt
	vLoginToServerPacketName, err = queser.ErroringIndex(LoginToServerPacketNameReverseMap, ret.Name)
	if err != nil {
		return
	}
	err = vLoginToServerPacketName.Encode(w)
	if err != nil {
		return
	}
	switch ret.Name {
	case "cookie_response":
		LoginToServerPacketParams, ok := ret.Params.(PacketCommonCookieResponse)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = LoginToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "encryption_begin":
		LoginToServerPacketParams, ok := ret.Params.(LoginToServerPacketEncryptionBegin)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = LoginToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "login_acknowledged":
		LoginToServerPacketParams, ok := ret.Params.(LoginToServerPacketLoginAcknowledged)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = LoginToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "login_plugin_response":
		LoginToServerPacketParams, ok := ret.Params.(LoginToServerPacketLoginPluginResponse)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = LoginToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "login_start":
		LoginToServerPacketParams, ok := ret.Params.(LoginToServerPacketLoginStart)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = LoginToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type LoginToServerPacketEncryptionBegin struct {
	SharedSecret []byte
	VerifyToken  []byte
}

func (_ LoginToServerPacketEncryptionBegin) Decode(r io.Reader) (ret LoginToServerPacketEncryptionBegin, err error) {
	var lLoginToServerPacketEncryptionBeginSharedSecret queser.VarInt
	lLoginToServerPacketEncryptionBeginSharedSecret, err = lLoginToServerPacketEncryptionBeginSharedSecret.Decode(r)
	if err != nil {
		return
	}
	ret.SharedSecret, err = io.ReadAll(io.LimitReader(r, int64(lLoginToServerPacketEncryptionBeginSharedSecret)))
	if err != nil {
		return
	}
	var lLoginToServerPacketEncryptionBeginVerifyToken queser.VarInt
	lLoginToServerPacketEncryptionBeginVerifyToken, err = lLoginToServerPacketEncryptionBeginVerifyToken.Decode(r)
	if err != nil {
		return
	}
	ret.VerifyToken, err = io.ReadAll(io.LimitReader(r, int64(lLoginToServerPacketEncryptionBeginVerifyToken)))
	if err != nil {
		return
	}
	return
}
func (ret LoginToServerPacketEncryptionBegin) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.SharedSecret)).Encode(w)
	if err != nil {
		return
	}
	_, err = w.Write(ret.SharedSecret)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.VerifyToken)).Encode(w)
	if err != nil {
		return
	}
	_, err = w.Write(ret.VerifyToken)
	if err != nil {
		return
	}
	return
}

type LoginToServerPacketLoginAcknowledged struct {
}

func (_ LoginToServerPacketLoginAcknowledged) Decode(r io.Reader) (ret LoginToServerPacketLoginAcknowledged, err error) {
	return
}
func (ret LoginToServerPacketLoginAcknowledged) Encode(w io.Writer) (err error) {
	return
}

type LoginToServerPacketLoginPluginResponse struct {
	MessageId queser.VarInt
	Data      *queser.RestBuffer
}

func (_ LoginToServerPacketLoginPluginResponse) Decode(r io.Reader) (ret LoginToServerPacketLoginPluginResponse, err error) {
	ret.MessageId, err = ret.MessageId.Decode(r)
	if err != nil {
		return
	}
	var LoginToServerPacketLoginPluginResponseDataPresent bool
	err = binary.Read(r, binary.BigEndian, &LoginToServerPacketLoginPluginResponseDataPresent)
	if err != nil {
		return
	}
	if LoginToServerPacketLoginPluginResponseDataPresent {
		var LoginToServerPacketLoginPluginResponseDataPresentValue queser.RestBuffer
		LoginToServerPacketLoginPluginResponseDataPresentValue, err = LoginToServerPacketLoginPluginResponseDataPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.Data = &LoginToServerPacketLoginPluginResponseDataPresentValue
	}
	return
}
func (ret LoginToServerPacketLoginPluginResponse) Encode(w io.Writer) (err error) {
	err = ret.MessageId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Data != nil)
	if err != nil {
		return
	}
	if ret.Data != nil {
		err = (*ret.Data).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type LoginToServerPacketLoginStart struct {
	Username   string
	PlayerUUID uuid.UUID
}

func (_ LoginToServerPacketLoginStart) Decode(r io.Reader) (ret LoginToServerPacketLoginStart, err error) {
	ret.Username, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	_, err = io.ReadFull(r, ret.PlayerUUID[:])
	if err != nil {
		return
	}
	return
}
func (ret LoginToServerPacketLoginStart) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Username)
	if err != nil {
		return
	}
	_, err = w.Write(ret.PlayerUUID[:])
	if err != nil {
		return
	}
	return
}

type LoginToClientPacket struct {
	Name   string
	Params any
}

var LoginToClientPacketNameMap = map[queser.VarInt]string{0x00: "disconnect", 0x01: "encryption_begin", 0x02: "success", 0x03: "compress", 0x04: "login_plugin_request", 0x05: "cookie_request"}

func (_ LoginToClientPacket) Decode(r io.Reader) (ret LoginToClientPacket, err error) {
	var LoginToClientPacketNameKey queser.VarInt
	LoginToClientPacketNameKey, err = LoginToClientPacketNameKey.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.ErroringIndex(LoginToClientPacketNameMap, LoginToClientPacketNameKey)
	if err != nil {
		return
	}
	switch ret.Name {
	case "compress":
		var LoginToClientPacketParamsTmp LoginToClientPacketCompress
		LoginToClientPacketParamsTmp, err = LoginToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = LoginToClientPacketParamsTmp
	case "cookie_request":
		var LoginToClientPacketParamsTmp PacketCommonCookieRequest
		LoginToClientPacketParamsTmp, err = LoginToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = LoginToClientPacketParamsTmp
	case "disconnect":
		var LoginToClientPacketParamsTmp LoginToClientPacketDisconnect
		LoginToClientPacketParamsTmp, err = LoginToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = LoginToClientPacketParamsTmp
	case "encryption_begin":
		var LoginToClientPacketParamsTmp LoginToClientPacketEncryptionBegin
		LoginToClientPacketParamsTmp, err = LoginToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = LoginToClientPacketParamsTmp
	case "login_plugin_request":
		var LoginToClientPacketParamsTmp LoginToClientPacketLoginPluginRequest
		LoginToClientPacketParamsTmp, err = LoginToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = LoginToClientPacketParamsTmp
	case "success":
		var LoginToClientPacketParamsTmp LoginToClientPacketSuccess
		LoginToClientPacketParamsTmp, err = LoginToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = LoginToClientPacketParamsTmp
	}
	return
}

var LoginToClientPacketNameReverseMap = map[string]queser.VarInt{"disconnect": 0x00, "encryption_begin": 0x01, "success": 0x02, "compress": 0x03, "login_plugin_request": 0x04, "cookie_request": 0x05}

func (ret LoginToClientPacket) Encode(w io.Writer) (err error) {
	var vLoginToClientPacketName queser.VarInt
	vLoginToClientPacketName, err = queser.ErroringIndex(LoginToClientPacketNameReverseMap, ret.Name)
	if err != nil {
		return
	}
	err = vLoginToClientPacketName.Encode(w)
	if err != nil {
		return
	}
	switch ret.Name {
	case "compress":
		LoginToClientPacketParams, ok := ret.Params.(LoginToClientPacketCompress)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = LoginToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "cookie_request":
		LoginToClientPacketParams, ok := ret.Params.(PacketCommonCookieRequest)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = LoginToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "disconnect":
		LoginToClientPacketParams, ok := ret.Params.(LoginToClientPacketDisconnect)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = LoginToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "encryption_begin":
		LoginToClientPacketParams, ok := ret.Params.(LoginToClientPacketEncryptionBegin)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = LoginToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "login_plugin_request":
		LoginToClientPacketParams, ok := ret.Params.(LoginToClientPacketLoginPluginRequest)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = LoginToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "success":
		LoginToClientPacketParams, ok := ret.Params.(LoginToClientPacketSuccess)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = LoginToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type LoginToClientPacketCompress struct {
	Threshold queser.VarInt
}

func (_ LoginToClientPacketCompress) Decode(r io.Reader) (ret LoginToClientPacketCompress, err error) {
	ret.Threshold, err = ret.Threshold.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret LoginToClientPacketCompress) Encode(w io.Writer) (err error) {
	err = ret.Threshold.Encode(w)
	if err != nil {
		return
	}
	return
}

type LoginToClientPacketDisconnect struct {
	Reason string
}

func (_ LoginToClientPacketDisconnect) Decode(r io.Reader) (ret LoginToClientPacketDisconnect, err error) {
	ret.Reason, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	return
}
func (ret LoginToClientPacketDisconnect) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Reason)
	if err != nil {
		return
	}
	return
}

type LoginToClientPacketEncryptionBegin struct {
	ServerId           string
	PublicKey          []byte
	VerifyToken        []byte
	ShouldAuthenticate bool
}

func (_ LoginToClientPacketEncryptionBegin) Decode(r io.Reader) (ret LoginToClientPacketEncryptionBegin, err error) {
	ret.ServerId, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var lLoginToClientPacketEncryptionBeginPublicKey queser.VarInt
	lLoginToClientPacketEncryptionBeginPublicKey, err = lLoginToClientPacketEncryptionBeginPublicKey.Decode(r)
	if err != nil {
		return
	}
	ret.PublicKey, err = io.ReadAll(io.LimitReader(r, int64(lLoginToClientPacketEncryptionBeginPublicKey)))
	if err != nil {
		return
	}
	var lLoginToClientPacketEncryptionBeginVerifyToken queser.VarInt
	lLoginToClientPacketEncryptionBeginVerifyToken, err = lLoginToClientPacketEncryptionBeginVerifyToken.Decode(r)
	if err != nil {
		return
	}
	ret.VerifyToken, err = io.ReadAll(io.LimitReader(r, int64(lLoginToClientPacketEncryptionBeginVerifyToken)))
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.ShouldAuthenticate)
	if err != nil {
		return
	}
	return
}
func (ret LoginToClientPacketEncryptionBegin) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.ServerId)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.PublicKey)).Encode(w)
	if err != nil {
		return
	}
	_, err = w.Write(ret.PublicKey)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.VerifyToken)).Encode(w)
	if err != nil {
		return
	}
	_, err = w.Write(ret.VerifyToken)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.ShouldAuthenticate)
	if err != nil {
		return
	}
	return
}

type LoginToClientPacketLoginPluginRequest struct {
	MessageId queser.VarInt
	Channel   string
	Data      queser.RestBuffer
}

func (_ LoginToClientPacketLoginPluginRequest) Decode(r io.Reader) (ret LoginToClientPacketLoginPluginRequest, err error) {
	ret.MessageId, err = ret.MessageId.Decode(r)
	if err != nil {
		return
	}
	ret.Channel, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Data, err = ret.Data.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret LoginToClientPacketLoginPluginRequest) Encode(w io.Writer) (err error) {
	err = ret.MessageId.Encode(w)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Channel)
	if err != nil {
		return
	}
	err = ret.Data.Encode(w)
	if err != nil {
		return
	}
	return
}

type LoginToClientPacketSuccess struct {
	Uuid       uuid.UUID
	Username   string
	Properties []struct {
		Name      string
		Value     string
		Signature *string
	}
}

func (_ LoginToClientPacketSuccess) Decode(r io.Reader) (ret LoginToClientPacketSuccess, err error) {
	_, err = io.ReadFull(r, ret.Uuid[:])
	if err != nil {
		return
	}
	ret.Username, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var lLoginToClientPacketSuccessProperties queser.VarInt
	lLoginToClientPacketSuccessProperties, err = lLoginToClientPacketSuccessProperties.Decode(r)
	if err != nil {
		return
	}
	ret.Properties = []struct {
		Name      string
		Value     string
		Signature *string
	}{}
	for range lLoginToClientPacketSuccessProperties {
		var LoginToClientPacketSuccessPropertiesElement struct {
			Name      string
			Value     string
			Signature *string
		}
		LoginToClientPacketSuccessPropertiesElement.Name, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		LoginToClientPacketSuccessPropertiesElement.Value, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		var LoginToClientPacketSuccessPropertiesElementSignaturePresent bool
		err = binary.Read(r, binary.BigEndian, &LoginToClientPacketSuccessPropertiesElementSignaturePresent)
		if err != nil {
			return
		}
		if LoginToClientPacketSuccessPropertiesElementSignaturePresent {
			var LoginToClientPacketSuccessPropertiesElementSignaturePresentValue string
			LoginToClientPacketSuccessPropertiesElementSignaturePresentValue, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			LoginToClientPacketSuccessPropertiesElement.Signature = &LoginToClientPacketSuccessPropertiesElementSignaturePresentValue
		}
		ret.Properties = append(ret.Properties, LoginToClientPacketSuccessPropertiesElement)
	}
	return
}
func (ret LoginToClientPacketSuccess) Encode(w io.Writer) (err error) {
	_, err = w.Write(ret.Uuid[:])
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Username)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Properties)).Encode(w)
	if err != nil {
		return
	}
	for iLoginToClientPacketSuccessProperties := range len(ret.Properties) {
		err = queser.EncodeString(w, ret.Properties[iLoginToClientPacketSuccessProperties].Name)
		if err != nil {
			return
		}
		err = queser.EncodeString(w, ret.Properties[iLoginToClientPacketSuccessProperties].Value)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Properties[iLoginToClientPacketSuccessProperties].Signature != nil)
		if err != nil {
			return
		}
		if ret.Properties[iLoginToClientPacketSuccessProperties].Signature != nil {
			err = queser.EncodeString(w, *ret.Properties[iLoginToClientPacketSuccessProperties].Signature)
			if err != nil {
				return
			}
		}
	}
	return
}

type ConfigurationToServerPacket struct {
	Name   string
	Params any
}

var ConfigurationToServerPacketNameMap = map[queser.VarInt]string{0x00: "settings", 0x01: "cookie_response", 0x02: "custom_payload", 0x03: "finish_configuration", 0x04: "keep_alive", 0x05: "pong", 0x06: "resource_pack_receive", 0x07: "select_known_packs", 0x08: "custom_click_action"}

func (_ ConfigurationToServerPacket) Decode(r io.Reader) (ret ConfigurationToServerPacket, err error) {
	var ConfigurationToServerPacketNameKey queser.VarInt
	ConfigurationToServerPacketNameKey, err = ConfigurationToServerPacketNameKey.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.ErroringIndex(ConfigurationToServerPacketNameMap, ConfigurationToServerPacketNameKey)
	if err != nil {
		return
	}
	switch ret.Name {
	case "cookie_response":
		var ConfigurationToServerPacketParamsTmp PacketCommonCookieResponse
		ConfigurationToServerPacketParamsTmp, err = ConfigurationToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToServerPacketParamsTmp
	case "custom_click_action":
		var ConfigurationToServerPacketParamsTmp PacketCommonCustomClickAction
		ConfigurationToServerPacketParamsTmp, err = ConfigurationToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToServerPacketParamsTmp
	case "custom_payload":
		var ConfigurationToServerPacketParamsTmp ConfigurationToServerPacketCustomPayload
		ConfigurationToServerPacketParamsTmp, err = ConfigurationToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToServerPacketParamsTmp
	case "finish_configuration":
		var ConfigurationToServerPacketParamsTmp ConfigurationToServerPacketFinishConfiguration
		ConfigurationToServerPacketParamsTmp, err = ConfigurationToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToServerPacketParamsTmp
	case "keep_alive":
		var ConfigurationToServerPacketParamsTmp ConfigurationToServerPacketKeepAlive
		ConfigurationToServerPacketParamsTmp, err = ConfigurationToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToServerPacketParamsTmp
	case "pong":
		var ConfigurationToServerPacketParamsTmp ConfigurationToServerPacketPong
		ConfigurationToServerPacketParamsTmp, err = ConfigurationToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToServerPacketParamsTmp
	case "resource_pack_receive":
		var ConfigurationToServerPacketParamsTmp ConfigurationToServerPacketResourcePackReceive
		ConfigurationToServerPacketParamsTmp, err = ConfigurationToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToServerPacketParamsTmp
	case "select_known_packs":
		var ConfigurationToServerPacketParamsTmp PacketCommonSelectKnownPacks
		ConfigurationToServerPacketParamsTmp, err = ConfigurationToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToServerPacketParamsTmp
	case "settings":
		var ConfigurationToServerPacketParamsTmp PacketCommonSettings
		ConfigurationToServerPacketParamsTmp, err = ConfigurationToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToServerPacketParamsTmp
	}
	return
}

var ConfigurationToServerPacketNameReverseMap = map[string]queser.VarInt{"settings": 0x00, "cookie_response": 0x01, "custom_payload": 0x02, "finish_configuration": 0x03, "keep_alive": 0x04, "pong": 0x05, "resource_pack_receive": 0x06, "select_known_packs": 0x07, "custom_click_action": 0x08}

func (ret ConfigurationToServerPacket) Encode(w io.Writer) (err error) {
	var vConfigurationToServerPacketName queser.VarInt
	vConfigurationToServerPacketName, err = queser.ErroringIndex(ConfigurationToServerPacketNameReverseMap, ret.Name)
	if err != nil {
		return
	}
	err = vConfigurationToServerPacketName.Encode(w)
	if err != nil {
		return
	}
	switch ret.Name {
	case "cookie_response":
		ConfigurationToServerPacketParams, ok := ret.Params.(PacketCommonCookieResponse)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "custom_click_action":
		ConfigurationToServerPacketParams, ok := ret.Params.(PacketCommonCustomClickAction)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "custom_payload":
		ConfigurationToServerPacketParams, ok := ret.Params.(ConfigurationToServerPacketCustomPayload)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "finish_configuration":
		ConfigurationToServerPacketParams, ok := ret.Params.(ConfigurationToServerPacketFinishConfiguration)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "keep_alive":
		ConfigurationToServerPacketParams, ok := ret.Params.(ConfigurationToServerPacketKeepAlive)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "pong":
		ConfigurationToServerPacketParams, ok := ret.Params.(ConfigurationToServerPacketPong)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "resource_pack_receive":
		ConfigurationToServerPacketParams, ok := ret.Params.(ConfigurationToServerPacketResourcePackReceive)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "select_known_packs":
		ConfigurationToServerPacketParams, ok := ret.Params.(PacketCommonSelectKnownPacks)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "settings":
		ConfigurationToServerPacketParams, ok := ret.Params.(PacketCommonSettings)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type ConfigurationToServerPacketCustomPayload struct {
	Channel string
	Data    queser.RestBuffer
}

func (_ ConfigurationToServerPacketCustomPayload) Decode(r io.Reader) (ret ConfigurationToServerPacketCustomPayload, err error) {
	ret.Channel, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Data, err = ret.Data.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret ConfigurationToServerPacketCustomPayload) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Channel)
	if err != nil {
		return
	}
	err = ret.Data.Encode(w)
	if err != nil {
		return
	}
	return
}

type ConfigurationToServerPacketFinishConfiguration struct {
}

func (_ ConfigurationToServerPacketFinishConfiguration) Decode(r io.Reader) (ret ConfigurationToServerPacketFinishConfiguration, err error) {
	return
}
func (ret ConfigurationToServerPacketFinishConfiguration) Encode(w io.Writer) (err error) {
	return
}

type ConfigurationToServerPacketKeepAlive struct {
	KeepAliveId int64
}

func (_ ConfigurationToServerPacketKeepAlive) Decode(r io.Reader) (ret ConfigurationToServerPacketKeepAlive, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.KeepAliveId)
	if err != nil {
		return
	}
	return
}
func (ret ConfigurationToServerPacketKeepAlive) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.KeepAliveId)
	if err != nil {
		return
	}
	return
}

type ConfigurationToServerPacketPong struct {
	Id int32
}

func (_ ConfigurationToServerPacketPong) Decode(r io.Reader) (ret ConfigurationToServerPacketPong, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Id)
	if err != nil {
		return
	}
	return
}
func (ret ConfigurationToServerPacketPong) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Id)
	if err != nil {
		return
	}
	return
}

type ConfigurationToServerPacketResourcePackReceive struct {
	Uuid   uuid.UUID
	Result queser.VarInt
}

func (_ ConfigurationToServerPacketResourcePackReceive) Decode(r io.Reader) (ret ConfigurationToServerPacketResourcePackReceive, err error) {
	_, err = io.ReadFull(r, ret.Uuid[:])
	if err != nil {
		return
	}
	ret.Result, err = ret.Result.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret ConfigurationToServerPacketResourcePackReceive) Encode(w io.Writer) (err error) {
	_, err = w.Write(ret.Uuid[:])
	if err != nil {
		return
	}
	err = ret.Result.Encode(w)
	if err != nil {
		return
	}
	return
}

type ConfigurationToClientPacket struct {
	Name   string
	Params any
}

var ConfigurationToClientPacketNameMap = map[queser.VarInt]string{0x00: "cookie_request", 0x01: "custom_payload", 0x02: "disconnect", 0x03: "finish_configuration", 0x04: "keep_alive", 0x05: "ping", 0x06: "reset_chat", 0x07: "registry_data", 0x08: "remove_resource_pack", 0x09: "add_resource_pack", 0x0a: "store_cookie", 0x0b: "transfer", 0x0c: "feature_flags", 0x0d: "tags", 0x0e: "select_known_packs", 0x0f: "custom_report_details", 0x10: "server_links", 0x11: "clear_dialog", 0x12: "show_dialog"}

func (_ ConfigurationToClientPacket) Decode(r io.Reader) (ret ConfigurationToClientPacket, err error) {
	var ConfigurationToClientPacketNameKey queser.VarInt
	ConfigurationToClientPacketNameKey, err = ConfigurationToClientPacketNameKey.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.ErroringIndex(ConfigurationToClientPacketNameMap, ConfigurationToClientPacketNameKey)
	if err != nil {
		return
	}
	switch ret.Name {
	case "add_resource_pack":
		var ConfigurationToClientPacketParamsTmp PacketCommonAddResourcePack
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "clear_dialog":
		var ConfigurationToClientPacketParamsTmp PacketCommonClearDialog
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "cookie_request":
		var ConfigurationToClientPacketParamsTmp PacketCommonCookieRequest
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "custom_payload":
		var ConfigurationToClientPacketParamsTmp ConfigurationToClientPacketCustomPayload
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "custom_report_details":
		var ConfigurationToClientPacketParamsTmp PacketCommonCustomReportDetails
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "disconnect":
		var ConfigurationToClientPacketParamsTmp ConfigurationToClientPacketDisconnect
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "feature_flags":
		var ConfigurationToClientPacketParamsTmp ConfigurationToClientPacketFeatureFlags
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "finish_configuration":
		var ConfigurationToClientPacketParamsTmp ConfigurationToClientPacketFinishConfiguration
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "keep_alive":
		var ConfigurationToClientPacketParamsTmp ConfigurationToClientPacketKeepAlive
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "ping":
		var ConfigurationToClientPacketParamsTmp ConfigurationToClientPacketPing
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "registry_data":
		var ConfigurationToClientPacketParamsTmp ConfigurationToClientPacketRegistryData
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "remove_resource_pack":
		var ConfigurationToClientPacketParamsTmp PacketCommonRemoveResourcePack
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "reset_chat":
		var ConfigurationToClientPacketParamsTmp ConfigurationToClientPacketResetChat
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "select_known_packs":
		var ConfigurationToClientPacketParamsTmp PacketCommonSelectKnownPacks
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "server_links":
		var ConfigurationToClientPacketParamsTmp PacketCommonServerLinks
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "show_dialog":
		var ConfigurationToClientPacketParamsTmp ConfigurationToClientPacketShowDialog
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "store_cookie":
		var ConfigurationToClientPacketParamsTmp PacketCommonStoreCookie
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "tags":
		var ConfigurationToClientPacketParamsTmp ConfigurationToClientPacketTags
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	case "transfer":
		var ConfigurationToClientPacketParamsTmp PacketCommonTransfer
		ConfigurationToClientPacketParamsTmp, err = ConfigurationToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = ConfigurationToClientPacketParamsTmp
	}
	return
}

var ConfigurationToClientPacketNameReverseMap = map[string]queser.VarInt{"cookie_request": 0x00, "custom_payload": 0x01, "disconnect": 0x02, "finish_configuration": 0x03, "keep_alive": 0x04, "ping": 0x05, "reset_chat": 0x06, "registry_data": 0x07, "remove_resource_pack": 0x08, "add_resource_pack": 0x09, "store_cookie": 0x0a, "transfer": 0x0b, "feature_flags": 0x0c, "tags": 0x0d, "select_known_packs": 0x0e, "custom_report_details": 0x0f, "server_links": 0x10, "clear_dialog": 0x11, "show_dialog": 0x12}

func (ret ConfigurationToClientPacket) Encode(w io.Writer) (err error) {
	var vConfigurationToClientPacketName queser.VarInt
	vConfigurationToClientPacketName, err = queser.ErroringIndex(ConfigurationToClientPacketNameReverseMap, ret.Name)
	if err != nil {
		return
	}
	err = vConfigurationToClientPacketName.Encode(w)
	if err != nil {
		return
	}
	switch ret.Name {
	case "add_resource_pack":
		ConfigurationToClientPacketParams, ok := ret.Params.(PacketCommonAddResourcePack)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "clear_dialog":
		ConfigurationToClientPacketParams, ok := ret.Params.(PacketCommonClearDialog)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "cookie_request":
		ConfigurationToClientPacketParams, ok := ret.Params.(PacketCommonCookieRequest)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "custom_payload":
		ConfigurationToClientPacketParams, ok := ret.Params.(ConfigurationToClientPacketCustomPayload)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "custom_report_details":
		ConfigurationToClientPacketParams, ok := ret.Params.(PacketCommonCustomReportDetails)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "disconnect":
		ConfigurationToClientPacketParams, ok := ret.Params.(ConfigurationToClientPacketDisconnect)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "feature_flags":
		ConfigurationToClientPacketParams, ok := ret.Params.(ConfigurationToClientPacketFeatureFlags)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "finish_configuration":
		ConfigurationToClientPacketParams, ok := ret.Params.(ConfigurationToClientPacketFinishConfiguration)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "keep_alive":
		ConfigurationToClientPacketParams, ok := ret.Params.(ConfigurationToClientPacketKeepAlive)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "ping":
		ConfigurationToClientPacketParams, ok := ret.Params.(ConfigurationToClientPacketPing)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "registry_data":
		ConfigurationToClientPacketParams, ok := ret.Params.(ConfigurationToClientPacketRegistryData)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "remove_resource_pack":
		ConfigurationToClientPacketParams, ok := ret.Params.(PacketCommonRemoveResourcePack)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "reset_chat":
		ConfigurationToClientPacketParams, ok := ret.Params.(ConfigurationToClientPacketResetChat)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "select_known_packs":
		ConfigurationToClientPacketParams, ok := ret.Params.(PacketCommonSelectKnownPacks)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "server_links":
		ConfigurationToClientPacketParams, ok := ret.Params.(PacketCommonServerLinks)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "show_dialog":
		ConfigurationToClientPacketParams, ok := ret.Params.(ConfigurationToClientPacketShowDialog)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "store_cookie":
		ConfigurationToClientPacketParams, ok := ret.Params.(PacketCommonStoreCookie)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "tags":
		ConfigurationToClientPacketParams, ok := ret.Params.(ConfigurationToClientPacketTags)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "transfer":
		ConfigurationToClientPacketParams, ok := ret.Params.(PacketCommonTransfer)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ConfigurationToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type ConfigurationToClientPacketCustomPayload struct {
	Channel string
	Data    queser.RestBuffer
}

func (_ ConfigurationToClientPacketCustomPayload) Decode(r io.Reader) (ret ConfigurationToClientPacketCustomPayload, err error) {
	ret.Channel, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Data, err = ret.Data.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret ConfigurationToClientPacketCustomPayload) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Channel)
	if err != nil {
		return
	}
	err = ret.Data.Encode(w)
	if err != nil {
		return
	}
	return
}

type ConfigurationToClientPacketDisconnect struct {
	Reason nbt.Anon
}

func (_ ConfigurationToClientPacketDisconnect) Decode(r io.Reader) (ret ConfigurationToClientPacketDisconnect, err error) {
	ret.Reason, err = ret.Reason.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret ConfigurationToClientPacketDisconnect) Encode(w io.Writer) (err error) {
	err = ret.Reason.Encode(w)
	if err != nil {
		return
	}
	return
}

type ConfigurationToClientPacketFeatureFlags struct {
	Features []string
}

func (_ ConfigurationToClientPacketFeatureFlags) Decode(r io.Reader) (ret ConfigurationToClientPacketFeatureFlags, err error) {
	var lConfigurationToClientPacketFeatureFlagsFeatures queser.VarInt
	lConfigurationToClientPacketFeatureFlagsFeatures, err = lConfigurationToClientPacketFeatureFlagsFeatures.Decode(r)
	if err != nil {
		return
	}
	ret.Features = []string{}
	for range lConfigurationToClientPacketFeatureFlagsFeatures {
		var ConfigurationToClientPacketFeatureFlagsFeaturesElement string
		ConfigurationToClientPacketFeatureFlagsFeaturesElement, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Features = append(ret.Features, ConfigurationToClientPacketFeatureFlagsFeaturesElement)
	}
	return
}
func (ret ConfigurationToClientPacketFeatureFlags) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Features)).Encode(w)
	if err != nil {
		return
	}
	for iConfigurationToClientPacketFeatureFlagsFeatures := range len(ret.Features) {
		err = queser.EncodeString(w, ret.Features[iConfigurationToClientPacketFeatureFlagsFeatures])
		if err != nil {
			return
		}
	}
	return
}

type ConfigurationToClientPacketFinishConfiguration struct {
}

func (_ ConfigurationToClientPacketFinishConfiguration) Decode(r io.Reader) (ret ConfigurationToClientPacketFinishConfiguration, err error) {
	return
}
func (ret ConfigurationToClientPacketFinishConfiguration) Encode(w io.Writer) (err error) {
	return
}

type ConfigurationToClientPacketKeepAlive struct {
	KeepAliveId int64
}

func (_ ConfigurationToClientPacketKeepAlive) Decode(r io.Reader) (ret ConfigurationToClientPacketKeepAlive, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.KeepAliveId)
	if err != nil {
		return
	}
	return
}
func (ret ConfigurationToClientPacketKeepAlive) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.KeepAliveId)
	if err != nil {
		return
	}
	return
}

type ConfigurationToClientPacketPing struct {
	Id int32
}

func (_ ConfigurationToClientPacketPing) Decode(r io.Reader) (ret ConfigurationToClientPacketPing, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Id)
	if err != nil {
		return
	}
	return
}
func (ret ConfigurationToClientPacketPing) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Id)
	if err != nil {
		return
	}
	return
}

type ConfigurationToClientPacketRegistryData struct {
	Id      string
	Entries []struct {
		Key   string
		Value *nbt.Anon
	}
}

func (_ ConfigurationToClientPacketRegistryData) Decode(r io.Reader) (ret ConfigurationToClientPacketRegistryData, err error) {
	ret.Id, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var lConfigurationToClientPacketRegistryDataEntries queser.VarInt
	lConfigurationToClientPacketRegistryDataEntries, err = lConfigurationToClientPacketRegistryDataEntries.Decode(r)
	if err != nil {
		return
	}
	ret.Entries = []struct {
		Key   string
		Value *nbt.Anon
	}{}
	for range lConfigurationToClientPacketRegistryDataEntries {
		var ConfigurationToClientPacketRegistryDataEntriesElement struct {
			Key   string
			Value *nbt.Anon
		}
		ConfigurationToClientPacketRegistryDataEntriesElement.Key, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		var ConfigurationToClientPacketRegistryDataEntriesElementValuePresent bool
		err = binary.Read(r, binary.BigEndian, &ConfigurationToClientPacketRegistryDataEntriesElementValuePresent)
		if err != nil {
			return
		}
		if ConfigurationToClientPacketRegistryDataEntriesElementValuePresent {
			var ConfigurationToClientPacketRegistryDataEntriesElementValuePresentValue nbt.Anon
			ConfigurationToClientPacketRegistryDataEntriesElementValuePresentValue, err = ConfigurationToClientPacketRegistryDataEntriesElementValuePresentValue.Decode(r)
			if err != nil {
				return
			}
			ConfigurationToClientPacketRegistryDataEntriesElement.Value = &ConfigurationToClientPacketRegistryDataEntriesElementValuePresentValue
		}
		ret.Entries = append(ret.Entries, ConfigurationToClientPacketRegistryDataEntriesElement)
	}
	return
}
func (ret ConfigurationToClientPacketRegistryData) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Id)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Entries)).Encode(w)
	if err != nil {
		return
	}
	for iConfigurationToClientPacketRegistryDataEntries := range len(ret.Entries) {
		err = queser.EncodeString(w, ret.Entries[iConfigurationToClientPacketRegistryDataEntries].Key)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Entries[iConfigurationToClientPacketRegistryDataEntries].Value != nil)
		if err != nil {
			return
		}
		if ret.Entries[iConfigurationToClientPacketRegistryDataEntries].Value != nil {
			err = (*ret.Entries[iConfigurationToClientPacketRegistryDataEntries].Value).Encode(w)
			if err != nil {
				return
			}
		}
	}
	return
}

type ConfigurationToClientPacketResetChat struct {
}

func (_ ConfigurationToClientPacketResetChat) Decode(r io.Reader) (ret ConfigurationToClientPacketResetChat, err error) {
	return
}
func (ret ConfigurationToClientPacketResetChat) Encode(w io.Writer) (err error) {
	return
}

type ConfigurationToClientPacketShowDialog struct {
	Dialog nbt.Anon
}

func (_ ConfigurationToClientPacketShowDialog) Decode(r io.Reader) (ret ConfigurationToClientPacketShowDialog, err error) {
	ret.Dialog, err = ret.Dialog.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret ConfigurationToClientPacketShowDialog) Encode(w io.Writer) (err error) {
	err = ret.Dialog.Encode(w)
	if err != nil {
		return
	}
	return
}

type ConfigurationToClientPacketTags struct {
	Tags []struct {
		TagType string
		Tags    Tags
	}
}

func (_ ConfigurationToClientPacketTags) Decode(r io.Reader) (ret ConfigurationToClientPacketTags, err error) {
	var lConfigurationToClientPacketTagsTags queser.VarInt
	lConfigurationToClientPacketTagsTags, err = lConfigurationToClientPacketTagsTags.Decode(r)
	if err != nil {
		return
	}
	ret.Tags = []struct {
		TagType string
		Tags    Tags
	}{}
	for range lConfigurationToClientPacketTagsTags {
		var ConfigurationToClientPacketTagsTagsElement struct {
			TagType string
			Tags    Tags
		}
		ConfigurationToClientPacketTagsTagsElement.TagType, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ConfigurationToClientPacketTagsTagsElement.Tags, err = ConfigurationToClientPacketTagsTagsElement.Tags.Decode(r)
		if err != nil {
			return
		}
		ret.Tags = append(ret.Tags, ConfigurationToClientPacketTagsTagsElement)
	}
	return
}
func (ret ConfigurationToClientPacketTags) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Tags)).Encode(w)
	if err != nil {
		return
	}
	for iConfigurationToClientPacketTagsTags := range len(ret.Tags) {
		err = queser.EncodeString(w, ret.Tags[iConfigurationToClientPacketTagsTags].TagType)
		if err != nil {
			return
		}
		err = ret.Tags[iConfigurationToClientPacketTagsTags].Tags.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToServerMovementFlags struct {
	Val uint8
}

func (_ PlayToServerMovementFlags) Decode(r io.Reader) (ret PlayToServerMovementFlags, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Val)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerMovementFlags) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Val)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacket struct {
	Name   string
	Params any
}

var PlayToServerPacketNameMap = map[queser.VarInt]string{0x00: "teleport_confirm", 0x01: "query_block_nbt", 0x02: "select_bundle_item", 0x03: "set_difficulty", 0x04: "change_gamemode", 0x05: "message_acknowledgement", 0x06: "chat_command", 0x07: "chat_command_signed", 0x08: "chat_message", 0x09: "chat_session_update", 0x0a: "chunk_batch_received", 0x0b: "client_command", 0x0c: "tick_end", 0x0d: "settings", 0x0e: "tab_complete", 0x0f: "configuration_acknowledged", 0x10: "enchant_item", 0x11: "window_click", 0x12: "close_window", 0x13: "set_slot_state", 0x14: "cookie_response", 0x15: "custom_payload", 0x16: "debug_sample_subscription", 0x17: "edit_book", 0x18: "query_entity_nbt", 0x19: "use_entity", 0x1a: "generate_structure", 0x1b: "keep_alive", 0x1c: "lock_difficulty", 0x1d: "position", 0x1e: "position_look", 0x1f: "look", 0x20: "flying", 0x21: "vehicle_move", 0x22: "steer_boat", 0x23: "pick_item_from_block", 0x24: "pick_item_from_entity", 0x25: "ping_request", 0x26: "craft_recipe_request", 0x27: "abilities", 0x28: "block_dig", 0x29: "entity_action", 0x2a: "player_input", 0x2b: "player_loaded", 0x2c: "pong", 0x2d: "recipe_book", 0x2e: "displayed_recipe", 0x2f: "name_item", 0x30: "resource_pack_receive", 0x31: "advancement_tab", 0x32: "select_trade", 0x33: "set_beacon_effect", 0x34: "held_item_slot", 0x35: "update_command_block", 0x36: "update_command_block_minecart", 0x37: "set_creative_slot", 0x38: "update_jigsaw_block", 0x39: "update_structure_block", 0x3a: "set_test_block", 0x3b: "update_sign", 0x3c: "arm_animation", 0x3d: "spectate", 0x3e: "test_instance_block_action", 0x3f: "block_place", 0x40: "use_item", 0x41: "custom_click_action"}

func (_ PlayToServerPacket) Decode(r io.Reader) (ret PlayToServerPacket, err error) {
	var PlayToServerPacketNameKey queser.VarInt
	PlayToServerPacketNameKey, err = PlayToServerPacketNameKey.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.ErroringIndex(PlayToServerPacketNameMap, PlayToServerPacketNameKey)
	if err != nil {
		return
	}
	switch ret.Name {
	case "abilities":
		var PlayToServerPacketParamsTmp PlayToServerPacketAbilities
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "advancement_tab":
		var PlayToServerPacketParamsTmp PlayToServerPacketAdvancementTab
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "arm_animation":
		var PlayToServerPacketParamsTmp PlayToServerPacketArmAnimation
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "block_dig":
		var PlayToServerPacketParamsTmp PlayToServerPacketBlockDig
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "block_place":
		var PlayToServerPacketParamsTmp PlayToServerPacketBlockPlace
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "change_gamemode":
		var PlayToServerPacketParamsTmp PlayToServerPacketChangeGamemode
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "chat_command":
		var PlayToServerPacketParamsTmp PlayToServerPacketChatCommand
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "chat_command_signed":
		var PlayToServerPacketParamsTmp PlayToServerPacketChatCommandSigned
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "chat_message":
		var PlayToServerPacketParamsTmp PlayToServerPacketChatMessage
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "chat_session_update":
		var PlayToServerPacketParamsTmp PlayToServerPacketChatSessionUpdate
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "chunk_batch_received":
		var PlayToServerPacketParamsTmp PlayToServerPacketChunkBatchReceived
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "client_command":
		var PlayToServerPacketParamsTmp PlayToServerPacketClientCommand
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "close_window":
		var PlayToServerPacketParamsTmp PlayToServerPacketCloseWindow
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "configuration_acknowledged":
		var PlayToServerPacketParamsTmp PlayToServerPacketConfigurationAcknowledged
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "cookie_response":
		var PlayToServerPacketParamsTmp PacketCommonCookieResponse
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "craft_recipe_request":
		var PlayToServerPacketParamsTmp PlayToServerPacketCraftRecipeRequest
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "custom_click_action":
		var PlayToServerPacketParamsTmp PacketCommonCustomClickAction
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "custom_payload":
		var PlayToServerPacketParamsTmp PlayToServerPacketCustomPayload
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "debug_sample_subscription":
		var PlayToServerPacketParamsTmp PlayToServerPacketDebugSampleSubscription
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "displayed_recipe":
		var PlayToServerPacketParamsTmp PlayToServerPacketDisplayedRecipe
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "edit_book":
		var PlayToServerPacketParamsTmp PlayToServerPacketEditBook
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "enchant_item":
		var PlayToServerPacketParamsTmp PlayToServerPacketEnchantItem
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "entity_action":
		var PlayToServerPacketParamsTmp PlayToServerPacketEntityAction
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "flying":
		var PlayToServerPacketParamsTmp PlayToServerPacketFlying
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "generate_structure":
		var PlayToServerPacketParamsTmp PlayToServerPacketGenerateStructure
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "held_item_slot":
		var PlayToServerPacketParamsTmp PlayToServerPacketHeldItemSlot
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "keep_alive":
		var PlayToServerPacketParamsTmp PlayToServerPacketKeepAlive
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "lock_difficulty":
		var PlayToServerPacketParamsTmp PlayToServerPacketLockDifficulty
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "look":
		var PlayToServerPacketParamsTmp PlayToServerPacketLook
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "message_acknowledgement":
		var PlayToServerPacketParamsTmp PlayToServerPacketMessageAcknowledgement
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "name_item":
		var PlayToServerPacketParamsTmp PlayToServerPacketNameItem
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "pick_item_from_block":
		var PlayToServerPacketParamsTmp PlayToServerPacketPickItemFromBlock
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "pick_item_from_entity":
		var PlayToServerPacketParamsTmp PlayToServerPacketPickItemFromEntity
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "ping_request":
		var PlayToServerPacketParamsTmp PlayToServerPacketPingRequest
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "player_input":
		var PlayToServerPacketParamsTmp PlayToServerPacketPlayerInput
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "player_loaded":
		var PlayToServerPacketParamsTmp PlayToServerPacketPlayerLoaded
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "pong":
		var PlayToServerPacketParamsTmp PlayToServerPacketPong
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "position":
		var PlayToServerPacketParamsTmp PlayToServerPacketPosition
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "position_look":
		var PlayToServerPacketParamsTmp PlayToServerPacketPositionLook
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "query_block_nbt":
		var PlayToServerPacketParamsTmp PlayToServerPacketQueryBlockNbt
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "query_entity_nbt":
		var PlayToServerPacketParamsTmp PlayToServerPacketQueryEntityNbt
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "recipe_book":
		var PlayToServerPacketParamsTmp PlayToServerPacketRecipeBook
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "resource_pack_receive":
		var PlayToServerPacketParamsTmp PlayToServerPacketResourcePackReceive
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "select_bundle_item":
		var PlayToServerPacketParamsTmp PlayToServerPacketSelectBundleItem
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "select_trade":
		var PlayToServerPacketParamsTmp PlayToServerPacketSelectTrade
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "set_beacon_effect":
		var PlayToServerPacketParamsTmp PlayToServerPacketSetBeaconEffect
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "set_creative_slot":
		var PlayToServerPacketParamsTmp PlayToServerPacketSetCreativeSlot
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "set_difficulty":
		var PlayToServerPacketParamsTmp PlayToServerPacketSetDifficulty
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "set_slot_state":
		var PlayToServerPacketParamsTmp PlayToServerPacketSetSlotState
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "set_test_block":
		var PlayToServerPacketParamsTmp PlayToServerPacketSetTestBlock
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "settings":
		var PlayToServerPacketParamsTmp PacketCommonSettings
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "spectate":
		var PlayToServerPacketParamsTmp PlayToServerPacketSpectate
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "steer_boat":
		var PlayToServerPacketParamsTmp PlayToServerPacketSteerBoat
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "tab_complete":
		var PlayToServerPacketParamsTmp PlayToServerPacketTabComplete
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "teleport_confirm":
		var PlayToServerPacketParamsTmp PlayToServerPacketTeleportConfirm
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "test_instance_block_action":
		var PlayToServerPacketParamsTmp PlayToServerPacketTestInstanceBlockAction
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "tick_end":
		var PlayToServerPacketParamsTmp PlayToServerPacketTickEnd
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "update_command_block":
		var PlayToServerPacketParamsTmp PlayToServerPacketUpdateCommandBlock
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "update_command_block_minecart":
		var PlayToServerPacketParamsTmp PlayToServerPacketUpdateCommandBlockMinecart
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "update_jigsaw_block":
		var PlayToServerPacketParamsTmp PlayToServerPacketUpdateJigsawBlock
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "update_sign":
		var PlayToServerPacketParamsTmp PlayToServerPacketUpdateSign
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "update_structure_block":
		var PlayToServerPacketParamsTmp PlayToServerPacketUpdateStructureBlock
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "use_entity":
		var PlayToServerPacketParamsTmp PlayToServerPacketUseEntity
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "use_item":
		var PlayToServerPacketParamsTmp PlayToServerPacketUseItem
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "vehicle_move":
		var PlayToServerPacketParamsTmp PlayToServerPacketVehicleMove
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	case "window_click":
		var PlayToServerPacketParamsTmp PlayToServerPacketWindowClick
		PlayToServerPacketParamsTmp, err = PlayToServerPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToServerPacketParamsTmp
	}
	return
}

var PlayToServerPacketNameReverseMap = map[string]queser.VarInt{"teleport_confirm": 0x00, "query_block_nbt": 0x01, "select_bundle_item": 0x02, "set_difficulty": 0x03, "change_gamemode": 0x04, "message_acknowledgement": 0x05, "chat_command": 0x06, "chat_command_signed": 0x07, "chat_message": 0x08, "chat_session_update": 0x09, "chunk_batch_received": 0x0a, "client_command": 0x0b, "tick_end": 0x0c, "settings": 0x0d, "tab_complete": 0x0e, "configuration_acknowledged": 0x0f, "enchant_item": 0x10, "window_click": 0x11, "close_window": 0x12, "set_slot_state": 0x13, "cookie_response": 0x14, "custom_payload": 0x15, "debug_sample_subscription": 0x16, "edit_book": 0x17, "query_entity_nbt": 0x18, "use_entity": 0x19, "generate_structure": 0x1a, "keep_alive": 0x1b, "lock_difficulty": 0x1c, "position": 0x1d, "position_look": 0x1e, "look": 0x1f, "flying": 0x20, "vehicle_move": 0x21, "steer_boat": 0x22, "pick_item_from_block": 0x23, "pick_item_from_entity": 0x24, "ping_request": 0x25, "craft_recipe_request": 0x26, "abilities": 0x27, "block_dig": 0x28, "entity_action": 0x29, "player_input": 0x2a, "player_loaded": 0x2b, "pong": 0x2c, "recipe_book": 0x2d, "displayed_recipe": 0x2e, "name_item": 0x2f, "resource_pack_receive": 0x30, "advancement_tab": 0x31, "select_trade": 0x32, "set_beacon_effect": 0x33, "held_item_slot": 0x34, "update_command_block": 0x35, "update_command_block_minecart": 0x36, "set_creative_slot": 0x37, "update_jigsaw_block": 0x38, "update_structure_block": 0x39, "set_test_block": 0x3a, "update_sign": 0x3b, "arm_animation": 0x3c, "spectate": 0x3d, "test_instance_block_action": 0x3e, "block_place": 0x3f, "use_item": 0x40, "custom_click_action": 0x41}

func (ret PlayToServerPacket) Encode(w io.Writer) (err error) {
	var vPlayToServerPacketName queser.VarInt
	vPlayToServerPacketName, err = queser.ErroringIndex(PlayToServerPacketNameReverseMap, ret.Name)
	if err != nil {
		return
	}
	err = vPlayToServerPacketName.Encode(w)
	if err != nil {
		return
	}
	switch ret.Name {
	case "abilities":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketAbilities)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "advancement_tab":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketAdvancementTab)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "arm_animation":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketArmAnimation)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "block_dig":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketBlockDig)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "block_place":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketBlockPlace)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "change_gamemode":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketChangeGamemode)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "chat_command":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketChatCommand)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "chat_command_signed":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketChatCommandSigned)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "chat_message":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketChatMessage)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "chat_session_update":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketChatSessionUpdate)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "chunk_batch_received":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketChunkBatchReceived)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "client_command":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketClientCommand)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "close_window":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketCloseWindow)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "configuration_acknowledged":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketConfigurationAcknowledged)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "cookie_response":
		PlayToServerPacketParams, ok := ret.Params.(PacketCommonCookieResponse)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "craft_recipe_request":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketCraftRecipeRequest)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "custom_click_action":
		PlayToServerPacketParams, ok := ret.Params.(PacketCommonCustomClickAction)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "custom_payload":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketCustomPayload)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "debug_sample_subscription":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketDebugSampleSubscription)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "displayed_recipe":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketDisplayedRecipe)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "edit_book":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketEditBook)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "enchant_item":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketEnchantItem)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_action":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketEntityAction)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "flying":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketFlying)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "generate_structure":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketGenerateStructure)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "held_item_slot":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketHeldItemSlot)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "keep_alive":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketKeepAlive)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "lock_difficulty":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketLockDifficulty)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "look":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketLook)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "message_acknowledgement":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketMessageAcknowledgement)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "name_item":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketNameItem)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "pick_item_from_block":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketPickItemFromBlock)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "pick_item_from_entity":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketPickItemFromEntity)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "ping_request":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketPingRequest)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "player_input":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketPlayerInput)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "player_loaded":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketPlayerLoaded)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "pong":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketPong)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "position":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketPosition)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "position_look":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketPositionLook)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "query_block_nbt":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketQueryBlockNbt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "query_entity_nbt":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketQueryEntityNbt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "recipe_book":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketRecipeBook)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "resource_pack_receive":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketResourcePackReceive)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "select_bundle_item":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketSelectBundleItem)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "select_trade":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketSelectTrade)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_beacon_effect":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketSetBeaconEffect)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_creative_slot":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketSetCreativeSlot)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_difficulty":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketSetDifficulty)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_slot_state":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketSetSlotState)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_test_block":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketSetTestBlock)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "settings":
		PlayToServerPacketParams, ok := ret.Params.(PacketCommonSettings)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "spectate":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketSpectate)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "steer_boat":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketSteerBoat)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "tab_complete":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketTabComplete)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "teleport_confirm":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketTeleportConfirm)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "test_instance_block_action":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketTestInstanceBlockAction)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "tick_end":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketTickEnd)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "update_command_block":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketUpdateCommandBlock)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "update_command_block_minecart":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketUpdateCommandBlockMinecart)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "update_jigsaw_block":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketUpdateJigsawBlock)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "update_sign":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketUpdateSign)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "update_structure_block":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketUpdateStructureBlock)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "use_entity":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketUseEntity)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "use_item":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketUseItem)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "vehicle_move":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketVehicleMove)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "window_click":
		PlayToServerPacketParams, ok := ret.Params.(PlayToServerPacketWindowClick)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketParams.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToServerPacketAbilities struct {
	Flags int8
}

func (_ PlayToServerPacketAbilities) Decode(r io.Reader) (ret PlayToServerPacketAbilities, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Flags)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketAbilities) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Flags)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketAdvancementTab struct {
	Action queser.VarInt
	TabId  any
}

func (_ PlayToServerPacketAdvancementTab) Decode(r io.Reader) (ret PlayToServerPacketAdvancementTab, err error) {
	ret.Action, err = ret.Action.Decode(r)
	if err != nil {
		return
	}
	switch ret.Action {
	case 0:
		var PlayToServerPacketAdvancementTabTabIdTmp string
		PlayToServerPacketAdvancementTabTabIdTmp, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.TabId = PlayToServerPacketAdvancementTabTabIdTmp
	case 1:
		var PlayToServerPacketAdvancementTabTabIdTmp queser.Void
		PlayToServerPacketAdvancementTabTabIdTmp, err = PlayToServerPacketAdvancementTabTabIdTmp.Decode(r)
		if err != nil {
			return
		}
		ret.TabId = PlayToServerPacketAdvancementTabTabIdTmp
	}
	return
}
func (ret PlayToServerPacketAdvancementTab) Encode(w io.Writer) (err error) {
	err = ret.Action.Encode(w)
	if err != nil {
		return
	}
	switch ret.Action {
	case 0:
		PlayToServerPacketAdvancementTabTabId, ok := ret.TabId.(string)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.EncodeString(w, PlayToServerPacketAdvancementTabTabId)
		if err != nil {
			return
		}
	case 1:
		PlayToServerPacketAdvancementTabTabId, ok := ret.TabId.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketAdvancementTabTabId.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToServerPacketArmAnimation struct {
	Hand queser.VarInt
}

func (_ PlayToServerPacketArmAnimation) Decode(r io.Reader) (ret PlayToServerPacketArmAnimation, err error) {
	ret.Hand, err = ret.Hand.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketArmAnimation) Encode(w io.Writer) (err error) {
	err = ret.Hand.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketBlockDig struct {
	Status   queser.VarInt
	Location Position
	Face     int8
	Sequence queser.VarInt
}

func (_ PlayToServerPacketBlockDig) Decode(r io.Reader) (ret PlayToServerPacketBlockDig, err error) {
	ret.Status, err = ret.Status.Decode(r)
	if err != nil {
		return
	}
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Face)
	if err != nil {
		return
	}
	ret.Sequence, err = ret.Sequence.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketBlockDig) Encode(w io.Writer) (err error) {
	err = ret.Status.Encode(w)
	if err != nil {
		return
	}
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Face)
	if err != nil {
		return
	}
	err = ret.Sequence.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketBlockPlace struct {
	Hand           queser.VarInt
	Location       Position
	Direction      queser.VarInt
	CursorX        float32
	CursorY        float32
	CursorZ        float32
	InsideBlock    bool
	WorldBorderHit bool
	Sequence       queser.VarInt
}

func (_ PlayToServerPacketBlockPlace) Decode(r io.Reader) (ret PlayToServerPacketBlockPlace, err error) {
	ret.Hand, err = ret.Hand.Decode(r)
	if err != nil {
		return
	}
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	ret.Direction, err = ret.Direction.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.CursorX)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.CursorY)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.CursorZ)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.InsideBlock)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.WorldBorderHit)
	if err != nil {
		return
	}
	ret.Sequence, err = ret.Sequence.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketBlockPlace) Encode(w io.Writer) (err error) {
	err = ret.Hand.Encode(w)
	if err != nil {
		return
	}
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = ret.Direction.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.CursorX)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.CursorY)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.CursorZ)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.InsideBlock)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.WorldBorderHit)
	if err != nil {
		return
	}
	err = ret.Sequence.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketChangeGamemode struct {
	Mode string
}

var PlayToServerPacketChangeGamemodeModeMap = map[queser.VarInt]string{0: "survival", 1: "creative", 2: "adventure", 3: "spectator"}

func (_ PlayToServerPacketChangeGamemode) Decode(r io.Reader) (ret PlayToServerPacketChangeGamemode, err error) {
	var PlayToServerPacketChangeGamemodeModeKey queser.VarInt
	PlayToServerPacketChangeGamemodeModeKey, err = PlayToServerPacketChangeGamemodeModeKey.Decode(r)
	if err != nil {
		return
	}
	ret.Mode, err = queser.ErroringIndex(PlayToServerPacketChangeGamemodeModeMap, PlayToServerPacketChangeGamemodeModeKey)
	if err != nil {
		return
	}
	return
}

var PlayToServerPacketChangeGamemodeModeReverseMap = map[string]queser.VarInt{"survival": 0, "creative": 1, "adventure": 2, "spectator": 3}

func (ret PlayToServerPacketChangeGamemode) Encode(w io.Writer) (err error) {
	var vPlayToServerPacketChangeGamemodeMode queser.VarInt
	vPlayToServerPacketChangeGamemodeMode, err = queser.ErroringIndex(PlayToServerPacketChangeGamemodeModeReverseMap, ret.Mode)
	if err != nil {
		return
	}
	err = vPlayToServerPacketChangeGamemodeMode.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketChatCommand struct {
	Command string
}

func (_ PlayToServerPacketChatCommand) Decode(r io.Reader) (ret PlayToServerPacketChatCommand, err error) {
	ret.Command, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketChatCommand) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Command)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketChatCommandSigned struct {
	Command            string
	Timestamp          int64
	Salt               int64
	ArgumentSignatures []struct {
		ArgumentName string
		Signature    [256]byte
	}
	MessageCount queser.VarInt
	Acknowledged [3]byte
	Checksum     int8
}

func (_ PlayToServerPacketChatCommandSigned) Decode(r io.Reader) (ret PlayToServerPacketChatCommandSigned, err error) {
	ret.Command, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Timestamp)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Salt)
	if err != nil {
		return
	}
	var lPlayToServerPacketChatCommandSignedArgumentSignatures queser.VarInt
	lPlayToServerPacketChatCommandSignedArgumentSignatures, err = lPlayToServerPacketChatCommandSignedArgumentSignatures.Decode(r)
	if err != nil {
		return
	}
	ret.ArgumentSignatures = []struct {
		ArgumentName string
		Signature    [256]byte
	}{}
	for range lPlayToServerPacketChatCommandSignedArgumentSignatures {
		var PlayToServerPacketChatCommandSignedArgumentSignaturesElement struct {
			ArgumentName string
			Signature    [256]byte
		}
		PlayToServerPacketChatCommandSignedArgumentSignaturesElement.ArgumentName, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		_, err = r.Read(PlayToServerPacketChatCommandSignedArgumentSignaturesElement.Signature[:])
		if err != nil {
			return
		}
		ret.ArgumentSignatures = append(ret.ArgumentSignatures, PlayToServerPacketChatCommandSignedArgumentSignaturesElement)
	}
	ret.MessageCount, err = ret.MessageCount.Decode(r)
	if err != nil {
		return
	}
	_, err = r.Read(ret.Acknowledged[:])
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Checksum)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketChatCommandSigned) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Command)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Timestamp)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Salt)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.ArgumentSignatures)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToServerPacketChatCommandSignedArgumentSignatures := range len(ret.ArgumentSignatures) {
		err = queser.EncodeString(w, ret.ArgumentSignatures[iPlayToServerPacketChatCommandSignedArgumentSignatures].ArgumentName)
		if err != nil {
			return
		}
		arr := ret.ArgumentSignatures[iPlayToServerPacketChatCommandSignedArgumentSignatures].Signature
		_, err = w.Write(arr[:])
		if err != nil {
			return
		}
	}
	err = ret.MessageCount.Encode(w)
	if err != nil {
		return
	}
	arr := ret.Acknowledged
	_, err = w.Write(arr[:])
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Checksum)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketChatMessage struct {
	Message      string
	Timestamp    int64
	Salt         int64
	Signature    *[256]byte
	Offset       queser.VarInt
	Acknowledged [3]byte
	Checksum     uint8
}

func (_ PlayToServerPacketChatMessage) Decode(r io.Reader) (ret PlayToServerPacketChatMessage, err error) {
	ret.Message, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Timestamp)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Salt)
	if err != nil {
		return
	}
	var PlayToServerPacketChatMessageSignaturePresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToServerPacketChatMessageSignaturePresent)
	if err != nil {
		return
	}
	if PlayToServerPacketChatMessageSignaturePresent {
		var PlayToServerPacketChatMessageSignaturePresentValue [256]byte
		_, err = r.Read(PlayToServerPacketChatMessageSignaturePresentValue[:])
		if err != nil {
			return
		}
		ret.Signature = &PlayToServerPacketChatMessageSignaturePresentValue
	}
	ret.Offset, err = ret.Offset.Decode(r)
	if err != nil {
		return
	}
	_, err = r.Read(ret.Acknowledged[:])
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Checksum)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketChatMessage) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Message)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Timestamp)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Salt)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Signature != nil)
	if err != nil {
		return
	}
	if ret.Signature != nil {
		arr := *ret.Signature
		_, err = w.Write(arr[:])
		if err != nil {
			return
		}
	}
	err = ret.Offset.Encode(w)
	if err != nil {
		return
	}
	arr := ret.Acknowledged
	_, err = w.Write(arr[:])
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Checksum)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketChatSessionUpdate struct {
	SessionUUID uuid.UUID
	ExpireTime  int64
	PublicKey   ByteArray
	Signature   ByteArray
}

func (_ PlayToServerPacketChatSessionUpdate) Decode(r io.Reader) (ret PlayToServerPacketChatSessionUpdate, err error) {
	_, err = io.ReadFull(r, ret.SessionUUID[:])
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.ExpireTime)
	if err != nil {
		return
	}
	ret.PublicKey, err = ret.PublicKey.Decode(r)
	if err != nil {
		return
	}
	ret.Signature, err = ret.Signature.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketChatSessionUpdate) Encode(w io.Writer) (err error) {
	_, err = w.Write(ret.SessionUUID[:])
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.ExpireTime)
	if err != nil {
		return
	}
	err = ret.PublicKey.Encode(w)
	if err != nil {
		return
	}
	err = ret.Signature.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketChunkBatchReceived struct {
	ChunksPerTick float32
}

func (_ PlayToServerPacketChunkBatchReceived) Decode(r io.Reader) (ret PlayToServerPacketChunkBatchReceived, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.ChunksPerTick)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketChunkBatchReceived) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.ChunksPerTick)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketClientCommand struct {
	ActionId queser.VarInt
}

func (_ PlayToServerPacketClientCommand) Decode(r io.Reader) (ret PlayToServerPacketClientCommand, err error) {
	ret.ActionId, err = ret.ActionId.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketClientCommand) Encode(w io.Writer) (err error) {
	err = ret.ActionId.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketCloseWindow struct {
	WindowId ContainerID
}

func (_ PlayToServerPacketCloseWindow) Decode(r io.Reader) (ret PlayToServerPacketCloseWindow, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketCloseWindow) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketConfigurationAcknowledged struct {
}

func (_ PlayToServerPacketConfigurationAcknowledged) Decode(r io.Reader) (ret PlayToServerPacketConfigurationAcknowledged, err error) {
	return
}
func (ret PlayToServerPacketConfigurationAcknowledged) Encode(w io.Writer) (err error) {
	return
}

type PlayToServerPacketCraftRecipeRequest struct {
	WindowId ContainerID
	RecipeId queser.VarInt
	MakeAll  bool
}

func (_ PlayToServerPacketCraftRecipeRequest) Decode(r io.Reader) (ret PlayToServerPacketCraftRecipeRequest, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	ret.RecipeId, err = ret.RecipeId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.MakeAll)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketCraftRecipeRequest) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	err = ret.RecipeId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.MakeAll)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketCustomPayload struct {
	Channel string
	Data    queser.RestBuffer
}

func (_ PlayToServerPacketCustomPayload) Decode(r io.Reader) (ret PlayToServerPacketCustomPayload, err error) {
	ret.Channel, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Data, err = ret.Data.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketCustomPayload) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Channel)
	if err != nil {
		return
	}
	err = ret.Data.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketDebugSampleSubscription struct {
	Type queser.VarInt
}

func (_ PlayToServerPacketDebugSampleSubscription) Decode(r io.Reader) (ret PlayToServerPacketDebugSampleSubscription, err error) {
	ret.Type, err = ret.Type.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketDebugSampleSubscription) Encode(w io.Writer) (err error) {
	err = ret.Type.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketDisplayedRecipe struct {
	RecipeId queser.VarInt
}

func (_ PlayToServerPacketDisplayedRecipe) Decode(r io.Reader) (ret PlayToServerPacketDisplayedRecipe, err error) {
	ret.RecipeId, err = ret.RecipeId.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketDisplayedRecipe) Encode(w io.Writer) (err error) {
	err = ret.RecipeId.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketEditBook struct {
	Hand  queser.VarInt
	Pages []string
	Title *string
}

func (_ PlayToServerPacketEditBook) Decode(r io.Reader) (ret PlayToServerPacketEditBook, err error) {
	ret.Hand, err = ret.Hand.Decode(r)
	if err != nil {
		return
	}
	var lPlayToServerPacketEditBookPages queser.VarInt
	lPlayToServerPacketEditBookPages, err = lPlayToServerPacketEditBookPages.Decode(r)
	if err != nil {
		return
	}
	ret.Pages = []string{}
	for range lPlayToServerPacketEditBookPages {
		var PlayToServerPacketEditBookPagesElement string
		PlayToServerPacketEditBookPagesElement, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Pages = append(ret.Pages, PlayToServerPacketEditBookPagesElement)
	}
	var PlayToServerPacketEditBookTitlePresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToServerPacketEditBookTitlePresent)
	if err != nil {
		return
	}
	if PlayToServerPacketEditBookTitlePresent {
		var PlayToServerPacketEditBookTitlePresentValue string
		PlayToServerPacketEditBookTitlePresentValue, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Title = &PlayToServerPacketEditBookTitlePresentValue
	}
	return
}
func (ret PlayToServerPacketEditBook) Encode(w io.Writer) (err error) {
	err = ret.Hand.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Pages)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToServerPacketEditBookPages := range len(ret.Pages) {
		err = queser.EncodeString(w, ret.Pages[iPlayToServerPacketEditBookPages])
		if err != nil {
			return
		}
	}
	err = binary.Write(w, binary.BigEndian, ret.Title != nil)
	if err != nil {
		return
	}
	if ret.Title != nil {
		err = queser.EncodeString(w, *ret.Title)
		if err != nil {
			return
		}
	}
	return
}

type PlayToServerPacketEnchantItem struct {
	WindowId    ContainerID
	Enchantment int8
}

func (_ PlayToServerPacketEnchantItem) Decode(r io.Reader) (ret PlayToServerPacketEnchantItem, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Enchantment)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketEnchantItem) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Enchantment)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketEntityAction struct {
	EntityId  queser.VarInt
	ActionId  string
	JumpBoost queser.VarInt
}

var PlayToServerPacketEntityActionActionIdMap = map[queser.VarInt]string{0: "leave_bed", 1: "start_sprinting", 2: "stop_sprinting", 3: "start_horse_jump", 4: "stop_horse_jump", 5: "open_vehicle_inventory", 6: "start_elytra_flying"}

func (_ PlayToServerPacketEntityAction) Decode(r io.Reader) (ret PlayToServerPacketEntityAction, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	var PlayToServerPacketEntityActionActionIdKey queser.VarInt
	PlayToServerPacketEntityActionActionIdKey, err = PlayToServerPacketEntityActionActionIdKey.Decode(r)
	if err != nil {
		return
	}
	ret.ActionId, err = queser.ErroringIndex(PlayToServerPacketEntityActionActionIdMap, PlayToServerPacketEntityActionActionIdKey)
	if err != nil {
		return
	}
	ret.JumpBoost, err = ret.JumpBoost.Decode(r)
	if err != nil {
		return
	}
	return
}

var PlayToServerPacketEntityActionActionIdReverseMap = map[string]queser.VarInt{"leave_bed": 0, "start_sprinting": 1, "stop_sprinting": 2, "start_horse_jump": 3, "stop_horse_jump": 4, "open_vehicle_inventory": 5, "start_elytra_flying": 6}

func (ret PlayToServerPacketEntityAction) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	var vPlayToServerPacketEntityActionActionId queser.VarInt
	vPlayToServerPacketEntityActionActionId, err = queser.ErroringIndex(PlayToServerPacketEntityActionActionIdReverseMap, ret.ActionId)
	if err != nil {
		return
	}
	err = vPlayToServerPacketEntityActionActionId.Encode(w)
	if err != nil {
		return
	}
	err = ret.JumpBoost.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketFlying struct {
	Flags PlayToServerMovementFlags
}

func (_ PlayToServerPacketFlying) Decode(r io.Reader) (ret PlayToServerPacketFlying, err error) {
	ret.Flags, err = ret.Flags.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketFlying) Encode(w io.Writer) (err error) {
	err = ret.Flags.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketGenerateStructure struct {
	Location    Position
	Levels      queser.VarInt
	KeepJigsaws bool
}

func (_ PlayToServerPacketGenerateStructure) Decode(r io.Reader) (ret PlayToServerPacketGenerateStructure, err error) {
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	ret.Levels, err = ret.Levels.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.KeepJigsaws)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketGenerateStructure) Encode(w io.Writer) (err error) {
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = ret.Levels.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.KeepJigsaws)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketHeldItemSlot struct {
	SlotId int16
}

func (_ PlayToServerPacketHeldItemSlot) Decode(r io.Reader) (ret PlayToServerPacketHeldItemSlot, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.SlotId)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketHeldItemSlot) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.SlotId)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketKeepAlive struct {
	KeepAliveId int64
}

func (_ PlayToServerPacketKeepAlive) Decode(r io.Reader) (ret PlayToServerPacketKeepAlive, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.KeepAliveId)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketKeepAlive) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.KeepAliveId)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketLockDifficulty struct {
	Locked bool
}

func (_ PlayToServerPacketLockDifficulty) Decode(r io.Reader) (ret PlayToServerPacketLockDifficulty, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Locked)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketLockDifficulty) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Locked)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketLook struct {
	Yaw   float32
	Pitch float32
	Flags PlayToServerMovementFlags
}

func (_ PlayToServerPacketLook) Decode(r io.Reader) (ret PlayToServerPacketLook, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	ret.Flags, err = ret.Flags.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketLook) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	err = ret.Flags.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketMessageAcknowledgement struct {
	Count queser.VarInt
}

func (_ PlayToServerPacketMessageAcknowledgement) Decode(r io.Reader) (ret PlayToServerPacketMessageAcknowledgement, err error) {
	ret.Count, err = ret.Count.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketMessageAcknowledgement) Encode(w io.Writer) (err error) {
	err = ret.Count.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketNameItem struct {
	Name string
}

func (_ PlayToServerPacketNameItem) Decode(r io.Reader) (ret PlayToServerPacketNameItem, err error) {
	ret.Name, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketNameItem) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Name)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketPickItemFromBlock struct {
	Position    Position
	IncludeData bool
}

func (_ PlayToServerPacketPickItemFromBlock) Decode(r io.Reader) (ret PlayToServerPacketPickItemFromBlock, err error) {
	ret.Position, err = ret.Position.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IncludeData)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketPickItemFromBlock) Encode(w io.Writer) (err error) {
	err = ret.Position.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IncludeData)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketPickItemFromEntity struct {
	EntityId    queser.VarInt
	IncludeData bool
}

func (_ PlayToServerPacketPickItemFromEntity) Decode(r io.Reader) (ret PlayToServerPacketPickItemFromEntity, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IncludeData)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketPickItemFromEntity) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IncludeData)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketPingRequest struct {
	Id int64
}

func (_ PlayToServerPacketPingRequest) Decode(r io.Reader) (ret PlayToServerPacketPingRequest, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Id)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketPingRequest) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Id)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketPlayerInput struct {
	Inputs uint8
}

func (_ PlayToServerPacketPlayerInput) Decode(r io.Reader) (ret PlayToServerPacketPlayerInput, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Inputs)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketPlayerInput) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Inputs)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketPlayerLoaded struct {
}

func (_ PlayToServerPacketPlayerLoaded) Decode(r io.Reader) (ret PlayToServerPacketPlayerLoaded, err error) {
	return
}
func (ret PlayToServerPacketPlayerLoaded) Encode(w io.Writer) (err error) {
	return
}

type PlayToServerPacketPong struct {
	Id int32
}

func (_ PlayToServerPacketPong) Decode(r io.Reader) (ret PlayToServerPacketPong, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Id)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketPong) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Id)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketPosition struct {
	X     float64
	Y     float64
	Z     float64
	Flags PlayToServerMovementFlags
}

func (_ PlayToServerPacketPosition) Decode(r io.Reader) (ret PlayToServerPacketPosition, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	ret.Flags, err = ret.Flags.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketPosition) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = ret.Flags.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketPositionLook struct {
	X     float64
	Y     float64
	Z     float64
	Yaw   float32
	Pitch float32
	Flags PlayToServerMovementFlags
}

func (_ PlayToServerPacketPositionLook) Decode(r io.Reader) (ret PlayToServerPacketPositionLook, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	ret.Flags, err = ret.Flags.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketPositionLook) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	err = ret.Flags.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketQueryBlockNbt struct {
	TransactionId queser.VarInt
	Location      Position
}

func (_ PlayToServerPacketQueryBlockNbt) Decode(r io.Reader) (ret PlayToServerPacketQueryBlockNbt, err error) {
	ret.TransactionId, err = ret.TransactionId.Decode(r)
	if err != nil {
		return
	}
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketQueryBlockNbt) Encode(w io.Writer) (err error) {
	err = ret.TransactionId.Encode(w)
	if err != nil {
		return
	}
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketQueryEntityNbt struct {
	TransactionId queser.VarInt
	EntityId      queser.VarInt
}

func (_ PlayToServerPacketQueryEntityNbt) Decode(r io.Reader) (ret PlayToServerPacketQueryEntityNbt, err error) {
	ret.TransactionId, err = ret.TransactionId.Decode(r)
	if err != nil {
		return
	}
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketQueryEntityNbt) Encode(w io.Writer) (err error) {
	err = ret.TransactionId.Encode(w)
	if err != nil {
		return
	}
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketRecipeBook struct {
	BookId       queser.VarInt
	BookOpen     bool
	FilterActive bool
}

func (_ PlayToServerPacketRecipeBook) Decode(r io.Reader) (ret PlayToServerPacketRecipeBook, err error) {
	ret.BookId, err = ret.BookId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.BookOpen)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.FilterActive)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketRecipeBook) Encode(w io.Writer) (err error) {
	err = ret.BookId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.BookOpen)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.FilterActive)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketResourcePackReceive struct {
	Uuid   uuid.UUID
	Result queser.VarInt
}

func (_ PlayToServerPacketResourcePackReceive) Decode(r io.Reader) (ret PlayToServerPacketResourcePackReceive, err error) {
	_, err = io.ReadFull(r, ret.Uuid[:])
	if err != nil {
		return
	}
	ret.Result, err = ret.Result.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketResourcePackReceive) Encode(w io.Writer) (err error) {
	_, err = w.Write(ret.Uuid[:])
	if err != nil {
		return
	}
	err = ret.Result.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketSelectBundleItem struct {
	SlotId            queser.VarInt
	SelectedItemIndex queser.VarInt
}

func (_ PlayToServerPacketSelectBundleItem) Decode(r io.Reader) (ret PlayToServerPacketSelectBundleItem, err error) {
	ret.SlotId, err = ret.SlotId.Decode(r)
	if err != nil {
		return
	}
	ret.SelectedItemIndex, err = ret.SelectedItemIndex.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketSelectBundleItem) Encode(w io.Writer) (err error) {
	err = ret.SlotId.Encode(w)
	if err != nil {
		return
	}
	err = ret.SelectedItemIndex.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketSelectTrade struct {
	Slot queser.VarInt
}

func (_ PlayToServerPacketSelectTrade) Decode(r io.Reader) (ret PlayToServerPacketSelectTrade, err error) {
	ret.Slot, err = ret.Slot.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketSelectTrade) Encode(w io.Writer) (err error) {
	err = ret.Slot.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketSetBeaconEffect struct {
	PrimaryEffect   *queser.VarInt
	SecondaryEffect *queser.VarInt
}

func (_ PlayToServerPacketSetBeaconEffect) Decode(r io.Reader) (ret PlayToServerPacketSetBeaconEffect, err error) {
	var PlayToServerPacketSetBeaconEffectPrimaryEffectPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToServerPacketSetBeaconEffectPrimaryEffectPresent)
	if err != nil {
		return
	}
	if PlayToServerPacketSetBeaconEffectPrimaryEffectPresent {
		var PlayToServerPacketSetBeaconEffectPrimaryEffectPresentValue queser.VarInt
		PlayToServerPacketSetBeaconEffectPrimaryEffectPresentValue, err = PlayToServerPacketSetBeaconEffectPrimaryEffectPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.PrimaryEffect = &PlayToServerPacketSetBeaconEffectPrimaryEffectPresentValue
	}
	var PlayToServerPacketSetBeaconEffectSecondaryEffectPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToServerPacketSetBeaconEffectSecondaryEffectPresent)
	if err != nil {
		return
	}
	if PlayToServerPacketSetBeaconEffectSecondaryEffectPresent {
		var PlayToServerPacketSetBeaconEffectSecondaryEffectPresentValue queser.VarInt
		PlayToServerPacketSetBeaconEffectSecondaryEffectPresentValue, err = PlayToServerPacketSetBeaconEffectSecondaryEffectPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.SecondaryEffect = &PlayToServerPacketSetBeaconEffectSecondaryEffectPresentValue
	}
	return
}
func (ret PlayToServerPacketSetBeaconEffect) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.PrimaryEffect != nil)
	if err != nil {
		return
	}
	if ret.PrimaryEffect != nil {
		err = (*ret.PrimaryEffect).Encode(w)
		if err != nil {
			return
		}
	}
	err = binary.Write(w, binary.BigEndian, ret.SecondaryEffect != nil)
	if err != nil {
		return
	}
	if ret.SecondaryEffect != nil {
		err = (*ret.SecondaryEffect).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToServerPacketSetCreativeSlot struct {
	Slot int16
	Item UntrustedSlot
}

func (_ PlayToServerPacketSetCreativeSlot) Decode(r io.Reader) (ret PlayToServerPacketSetCreativeSlot, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Slot)
	if err != nil {
		return
	}
	ret.Item, err = ret.Item.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketSetCreativeSlot) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Slot)
	if err != nil {
		return
	}
	err = ret.Item.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketSetDifficulty struct {
	NewDifficulty string
}

var PlayToServerPacketSetDifficultyNewDifficultyMap = map[queser.VarInt]string{0: "peaceful", 1: "easy", 2: "normal", 3: "hard"}

func (_ PlayToServerPacketSetDifficulty) Decode(r io.Reader) (ret PlayToServerPacketSetDifficulty, err error) {
	var PlayToServerPacketSetDifficultyNewDifficultyKey queser.VarInt
	PlayToServerPacketSetDifficultyNewDifficultyKey, err = PlayToServerPacketSetDifficultyNewDifficultyKey.Decode(r)
	if err != nil {
		return
	}
	ret.NewDifficulty, err = queser.ErroringIndex(PlayToServerPacketSetDifficultyNewDifficultyMap, PlayToServerPacketSetDifficultyNewDifficultyKey)
	if err != nil {
		return
	}
	return
}

var PlayToServerPacketSetDifficultyNewDifficultyReverseMap = map[string]queser.VarInt{"peaceful": 0, "easy": 1, "normal": 2, "hard": 3}

func (ret PlayToServerPacketSetDifficulty) Encode(w io.Writer) (err error) {
	var vPlayToServerPacketSetDifficultyNewDifficulty queser.VarInt
	vPlayToServerPacketSetDifficultyNewDifficulty, err = queser.ErroringIndex(PlayToServerPacketSetDifficultyNewDifficultyReverseMap, ret.NewDifficulty)
	if err != nil {
		return
	}
	err = vPlayToServerPacketSetDifficultyNewDifficulty.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketSetSlotState struct {
	SlotId   queser.VarInt
	WindowId ContainerID
	State    bool
}

func (_ PlayToServerPacketSetSlotState) Decode(r io.Reader) (ret PlayToServerPacketSetSlotState, err error) {
	ret.SlotId, err = ret.SlotId.Decode(r)
	if err != nil {
		return
	}
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.State)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketSetSlotState) Encode(w io.Writer) (err error) {
	err = ret.SlotId.Encode(w)
	if err != nil {
		return
	}
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.State)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketSetTestBlock struct {
	Position Position
	Mode     queser.VarInt
	Message  string
}

func (_ PlayToServerPacketSetTestBlock) Decode(r io.Reader) (ret PlayToServerPacketSetTestBlock, err error) {
	ret.Position, err = ret.Position.Decode(r)
	if err != nil {
		return
	}
	ret.Mode, err = ret.Mode.Decode(r)
	if err != nil {
		return
	}
	ret.Message, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketSetTestBlock) Encode(w io.Writer) (err error) {
	err = ret.Position.Encode(w)
	if err != nil {
		return
	}
	err = ret.Mode.Encode(w)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Message)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketSpectate struct {
	Target uuid.UUID
}

func (_ PlayToServerPacketSpectate) Decode(r io.Reader) (ret PlayToServerPacketSpectate, err error) {
	_, err = io.ReadFull(r, ret.Target[:])
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketSpectate) Encode(w io.Writer) (err error) {
	_, err = w.Write(ret.Target[:])
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketSteerBoat struct {
	LeftPaddle  bool
	RightPaddle bool
}

func (_ PlayToServerPacketSteerBoat) Decode(r io.Reader) (ret PlayToServerPacketSteerBoat, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.LeftPaddle)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.RightPaddle)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketSteerBoat) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.LeftPaddle)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.RightPaddle)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketTabComplete struct {
	TransactionId queser.VarInt
	Text          string
}

func (_ PlayToServerPacketTabComplete) Decode(r io.Reader) (ret PlayToServerPacketTabComplete, err error) {
	ret.TransactionId, err = ret.TransactionId.Decode(r)
	if err != nil {
		return
	}
	ret.Text, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketTabComplete) Encode(w io.Writer) (err error) {
	err = ret.TransactionId.Encode(w)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Text)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketTeleportConfirm struct {
	TeleportId queser.VarInt
}

func (_ PlayToServerPacketTeleportConfirm) Decode(r io.Reader) (ret PlayToServerPacketTeleportConfirm, err error) {
	ret.TeleportId, err = ret.TeleportId.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketTeleportConfirm) Encode(w io.Writer) (err error) {
	err = ret.TeleportId.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketTestInstanceBlockAction struct {
	Pos    Position
	Action queser.VarInt
	Data   struct {
		Test           *string
		Size           Vec3i
		Rotation       queser.VarInt
		IgnoreEntities bool
		Status         queser.VarInt
		ErrorMessage   *nbt.Anon
	}
}

func (_ PlayToServerPacketTestInstanceBlockAction) Decode(r io.Reader) (ret PlayToServerPacketTestInstanceBlockAction, err error) {
	ret.Pos, err = ret.Pos.Decode(r)
	if err != nil {
		return
	}
	ret.Action, err = ret.Action.Decode(r)
	if err != nil {
		return
	}
	var PlayToServerPacketTestInstanceBlockActionDataTestPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToServerPacketTestInstanceBlockActionDataTestPresent)
	if err != nil {
		return
	}
	if PlayToServerPacketTestInstanceBlockActionDataTestPresent {
		var PlayToServerPacketTestInstanceBlockActionDataTestPresentValue string
		PlayToServerPacketTestInstanceBlockActionDataTestPresentValue, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Data.Test = &PlayToServerPacketTestInstanceBlockActionDataTestPresentValue
	}
	ret.Data.Size, err = ret.Data.Size.Decode(r)
	if err != nil {
		return
	}
	ret.Data.Rotation, err = ret.Data.Rotation.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Data.IgnoreEntities)
	if err != nil {
		return
	}
	ret.Data.Status, err = ret.Data.Status.Decode(r)
	if err != nil {
		return
	}
	var PlayToServerPacketTestInstanceBlockActionDataErrorMessagePresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToServerPacketTestInstanceBlockActionDataErrorMessagePresent)
	if err != nil {
		return
	}
	if PlayToServerPacketTestInstanceBlockActionDataErrorMessagePresent {
		var PlayToServerPacketTestInstanceBlockActionDataErrorMessagePresentValue nbt.Anon
		PlayToServerPacketTestInstanceBlockActionDataErrorMessagePresentValue, err = PlayToServerPacketTestInstanceBlockActionDataErrorMessagePresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.Data.ErrorMessage = &PlayToServerPacketTestInstanceBlockActionDataErrorMessagePresentValue
	}
	return
}
func (ret PlayToServerPacketTestInstanceBlockAction) Encode(w io.Writer) (err error) {
	err = ret.Pos.Encode(w)
	if err != nil {
		return
	}
	err = ret.Action.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Data.Test != nil)
	if err != nil {
		return
	}
	if ret.Data.Test != nil {
		err = queser.EncodeString(w, *ret.Data.Test)
		if err != nil {
			return
		}
	}
	err = ret.Data.Size.Encode(w)
	if err != nil {
		return
	}
	err = ret.Data.Rotation.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Data.IgnoreEntities)
	if err != nil {
		return
	}
	err = ret.Data.Status.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Data.ErrorMessage != nil)
	if err != nil {
		return
	}
	if ret.Data.ErrorMessage != nil {
		err = (*ret.Data.ErrorMessage).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToServerPacketTickEnd struct {
}

func (_ PlayToServerPacketTickEnd) Decode(r io.Reader) (ret PlayToServerPacketTickEnd, err error) {
	return
}
func (ret PlayToServerPacketTickEnd) Encode(w io.Writer) (err error) {
	return
}

type PlayToServerPacketUpdateCommandBlock struct {
	Location Position
	Command  string
	Mode     queser.VarInt
	Flags    uint8
}

func (_ PlayToServerPacketUpdateCommandBlock) Decode(r io.Reader) (ret PlayToServerPacketUpdateCommandBlock, err error) {
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	ret.Command, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Mode, err = ret.Mode.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Flags)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketUpdateCommandBlock) Encode(w io.Writer) (err error) {
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Command)
	if err != nil {
		return
	}
	err = ret.Mode.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Flags)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketUpdateCommandBlockMinecart struct {
	EntityId    queser.VarInt
	Command     string
	TrackOutput bool
}

func (_ PlayToServerPacketUpdateCommandBlockMinecart) Decode(r io.Reader) (ret PlayToServerPacketUpdateCommandBlockMinecart, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	ret.Command, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.TrackOutput)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketUpdateCommandBlockMinecart) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Command)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.TrackOutput)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketUpdateJigsawBlock struct {
	Location          Position
	Name              string
	Target            string
	Pool              string
	FinalState        string
	JointType         string
	SelectionPriority queser.VarInt
	PlacementPriority queser.VarInt
}

func (_ PlayToServerPacketUpdateJigsawBlock) Decode(r io.Reader) (ret PlayToServerPacketUpdateJigsawBlock, err error) {
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Target, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Pool, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.FinalState, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.JointType, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.SelectionPriority, err = ret.SelectionPriority.Decode(r)
	if err != nil {
		return
	}
	ret.PlacementPriority, err = ret.PlacementPriority.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketUpdateJigsawBlock) Encode(w io.Writer) (err error) {
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Name)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Target)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Pool)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.FinalState)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.JointType)
	if err != nil {
		return
	}
	err = ret.SelectionPriority.Encode(w)
	if err != nil {
		return
	}
	err = ret.PlacementPriority.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketUpdateSign struct {
	Location    Position
	IsFrontText bool
	Text1       string
	Text2       string
	Text3       string
	Text4       string
}

func (_ PlayToServerPacketUpdateSign) Decode(r io.Reader) (ret PlayToServerPacketUpdateSign, err error) {
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IsFrontText)
	if err != nil {
		return
	}
	ret.Text1, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Text2, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Text3, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Text4, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketUpdateSign) Encode(w io.Writer) (err error) {
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IsFrontText)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Text1)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Text2)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Text3)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Text4)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketUpdateStructureBlock struct {
	Location  Position
	Action    queser.VarInt
	Mode      queser.VarInt
	Name      string
	OffsetX   int8
	OffsetY   int8
	OffsetZ   int8
	SizeX     int8
	SizeY     int8
	SizeZ     int8
	Mirror    queser.VarInt
	Rotation  queser.VarInt
	Metadata  string
	Integrity float32
	Seed      queser.VarInt
	Flags     string
}

var PlayToServerPacketUpdateStructureBlockFlagsMap = map[uint8]string{0: "ignore_entities", 1: "show_air", 2: "show_bounding_box", 3: "strict"}

func (_ PlayToServerPacketUpdateStructureBlock) Decode(r io.Reader) (ret PlayToServerPacketUpdateStructureBlock, err error) {
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	ret.Action, err = ret.Action.Decode(r)
	if err != nil {
		return
	}
	ret.Mode, err = ret.Mode.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OffsetX)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OffsetY)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OffsetZ)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.SizeX)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.SizeY)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.SizeZ)
	if err != nil {
		return
	}
	ret.Mirror, err = ret.Mirror.Decode(r)
	if err != nil {
		return
	}
	ret.Rotation, err = ret.Rotation.Decode(r)
	if err != nil {
		return
	}
	ret.Metadata, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Integrity)
	if err != nil {
		return
	}
	ret.Seed, err = ret.Seed.Decode(r)
	if err != nil {
		return
	}
	var PlayToServerPacketUpdateStructureBlockFlagsKey uint8
	err = binary.Read(r, binary.BigEndian, &PlayToServerPacketUpdateStructureBlockFlagsKey)
	if err != nil {
		return
	}
	ret.Flags, err = queser.ErroringIndex(PlayToServerPacketUpdateStructureBlockFlagsMap, PlayToServerPacketUpdateStructureBlockFlagsKey)
	if err != nil {
		return
	}
	return
}

var PlayToServerPacketUpdateStructureBlockFlagsReverseMap = map[string]uint8{"ignore_entities": 0, "show_air": 1, "show_bounding_box": 2, "strict": 3}

func (ret PlayToServerPacketUpdateStructureBlock) Encode(w io.Writer) (err error) {
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = ret.Action.Encode(w)
	if err != nil {
		return
	}
	err = ret.Mode.Encode(w)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Name)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OffsetX)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OffsetY)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OffsetZ)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.SizeX)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.SizeY)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.SizeZ)
	if err != nil {
		return
	}
	err = ret.Mirror.Encode(w)
	if err != nil {
		return
	}
	err = ret.Rotation.Encode(w)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Metadata)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Integrity)
	if err != nil {
		return
	}
	err = ret.Seed.Encode(w)
	if err != nil {
		return
	}
	var vPlayToServerPacketUpdateStructureBlockFlags uint8
	vPlayToServerPacketUpdateStructureBlockFlags, err = queser.ErroringIndex(PlayToServerPacketUpdateStructureBlockFlagsReverseMap, ret.Flags)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, vPlayToServerPacketUpdateStructureBlockFlags)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketUseEntity struct {
	Target   queser.VarInt
	Mouse    queser.VarInt
	X        any
	Y        any
	Z        any
	Hand     any
	Sneaking bool
}

func (_ PlayToServerPacketUseEntity) Decode(r io.Reader) (ret PlayToServerPacketUseEntity, err error) {
	ret.Target, err = ret.Target.Decode(r)
	if err != nil {
		return
	}
	ret.Mouse, err = ret.Mouse.Decode(r)
	if err != nil {
		return
	}
	switch ret.Mouse {
	case 2:
		var PlayToServerPacketUseEntityXTmp float32
		err = binary.Read(r, binary.BigEndian, &PlayToServerPacketUseEntityXTmp)
		if err != nil {
			return
		}
		ret.X = PlayToServerPacketUseEntityXTmp
	default:
		var PlayToServerPacketUseEntityXTmp queser.Void
		PlayToServerPacketUseEntityXTmp, err = PlayToServerPacketUseEntityXTmp.Decode(r)
		if err != nil {
			return
		}
		ret.X = PlayToServerPacketUseEntityXTmp
	}
	switch ret.Mouse {
	case 2:
		var PlayToServerPacketUseEntityYTmp float32
		err = binary.Read(r, binary.BigEndian, &PlayToServerPacketUseEntityYTmp)
		if err != nil {
			return
		}
		ret.Y = PlayToServerPacketUseEntityYTmp
	default:
		var PlayToServerPacketUseEntityYTmp queser.Void
		PlayToServerPacketUseEntityYTmp, err = PlayToServerPacketUseEntityYTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Y = PlayToServerPacketUseEntityYTmp
	}
	switch ret.Mouse {
	case 2:
		var PlayToServerPacketUseEntityZTmp float32
		err = binary.Read(r, binary.BigEndian, &PlayToServerPacketUseEntityZTmp)
		if err != nil {
			return
		}
		ret.Z = PlayToServerPacketUseEntityZTmp
	default:
		var PlayToServerPacketUseEntityZTmp queser.Void
		PlayToServerPacketUseEntityZTmp, err = PlayToServerPacketUseEntityZTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Z = PlayToServerPacketUseEntityZTmp
	}
	switch ret.Mouse {
	case 0:
		var PlayToServerPacketUseEntityHandTmp queser.VarInt
		PlayToServerPacketUseEntityHandTmp, err = PlayToServerPacketUseEntityHandTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Hand = PlayToServerPacketUseEntityHandTmp
	case 2:
		var PlayToServerPacketUseEntityHandTmp queser.VarInt
		PlayToServerPacketUseEntityHandTmp, err = PlayToServerPacketUseEntityHandTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Hand = PlayToServerPacketUseEntityHandTmp
	default:
		var PlayToServerPacketUseEntityHandTmp queser.Void
		PlayToServerPacketUseEntityHandTmp, err = PlayToServerPacketUseEntityHandTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Hand = PlayToServerPacketUseEntityHandTmp
	}
	err = binary.Read(r, binary.BigEndian, &ret.Sneaking)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketUseEntity) Encode(w io.Writer) (err error) {
	err = ret.Target.Encode(w)
	if err != nil {
		return
	}
	err = ret.Mouse.Encode(w)
	if err != nil {
		return
	}
	switch ret.Mouse {
	case 2:
		PlayToServerPacketUseEntityX, ok := ret.X.(float32)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToServerPacketUseEntityX)
		if err != nil {
			return
		}
	default:
		_, ok := ret.X.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.X.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Mouse {
	case 2:
		PlayToServerPacketUseEntityY, ok := ret.Y.(float32)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToServerPacketUseEntityY)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Y.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Y.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Mouse {
	case 2:
		PlayToServerPacketUseEntityZ, ok := ret.Z.(float32)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToServerPacketUseEntityZ)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Z.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Z.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Mouse {
	case 0:
		PlayToServerPacketUseEntityHand, ok := ret.Hand.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketUseEntityHand.Encode(w)
		if err != nil {
			return
		}
	case 2:
		PlayToServerPacketUseEntityHand, ok := ret.Hand.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToServerPacketUseEntityHand.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Hand.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Hand.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	err = binary.Write(w, binary.BigEndian, ret.Sneaking)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketUseItem struct {
	Hand     queser.VarInt
	Sequence queser.VarInt
	Rotation Vec2f
}

func (_ PlayToServerPacketUseItem) Decode(r io.Reader) (ret PlayToServerPacketUseItem, err error) {
	ret.Hand, err = ret.Hand.Decode(r)
	if err != nil {
		return
	}
	ret.Sequence, err = ret.Sequence.Decode(r)
	if err != nil {
		return
	}
	ret.Rotation, err = ret.Rotation.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketUseItem) Encode(w io.Writer) (err error) {
	err = ret.Hand.Encode(w)
	if err != nil {
		return
	}
	err = ret.Sequence.Encode(w)
	if err != nil {
		return
	}
	err = ret.Rotation.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketVehicleMove struct {
	X        float64
	Y        float64
	Z        float64
	Yaw      float32
	Pitch    float32
	OnGround bool
}

func (_ PlayToServerPacketVehicleMove) Decode(r io.Reader) (ret PlayToServerPacketVehicleMove, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OnGround)
	if err != nil {
		return
	}
	return
}
func (ret PlayToServerPacketVehicleMove) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OnGround)
	if err != nil {
		return
	}
	return
}

type PlayToServerPacketWindowClick struct {
	WindowId     ContainerID
	StateId      queser.VarInt
	Slot         int16
	MouseButton  int8
	Mode         queser.VarInt
	ChangedSlots []struct {
		Location int16
		Item     *HashedSlot
	}
	CursorItem *HashedSlot
}

func (_ PlayToServerPacketWindowClick) Decode(r io.Reader) (ret PlayToServerPacketWindowClick, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	ret.StateId, err = ret.StateId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Slot)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.MouseButton)
	if err != nil {
		return
	}
	ret.Mode, err = ret.Mode.Decode(r)
	if err != nil {
		return
	}
	var lPlayToServerPacketWindowClickChangedSlots queser.VarInt
	lPlayToServerPacketWindowClickChangedSlots, err = lPlayToServerPacketWindowClickChangedSlots.Decode(r)
	if err != nil {
		return
	}
	ret.ChangedSlots = []struct {
		Location int16
		Item     *HashedSlot
	}{}
	for range lPlayToServerPacketWindowClickChangedSlots {
		var PlayToServerPacketWindowClickChangedSlotsElement struct {
			Location int16
			Item     *HashedSlot
		}
		err = binary.Read(r, binary.BigEndian, &PlayToServerPacketWindowClickChangedSlotsElement.Location)
		if err != nil {
			return
		}
		var PlayToServerPacketWindowClickChangedSlotsElementItemPresent bool
		err = binary.Read(r, binary.BigEndian, &PlayToServerPacketWindowClickChangedSlotsElementItemPresent)
		if err != nil {
			return
		}
		if PlayToServerPacketWindowClickChangedSlotsElementItemPresent {
			var PlayToServerPacketWindowClickChangedSlotsElementItemPresentValue HashedSlot
			PlayToServerPacketWindowClickChangedSlotsElementItemPresentValue, err = PlayToServerPacketWindowClickChangedSlotsElementItemPresentValue.Decode(r)
			if err != nil {
				return
			}
			PlayToServerPacketWindowClickChangedSlotsElement.Item = &PlayToServerPacketWindowClickChangedSlotsElementItemPresentValue
		}
		ret.ChangedSlots = append(ret.ChangedSlots, PlayToServerPacketWindowClickChangedSlotsElement)
	}
	var PlayToServerPacketWindowClickCursorItemPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToServerPacketWindowClickCursorItemPresent)
	if err != nil {
		return
	}
	if PlayToServerPacketWindowClickCursorItemPresent {
		var PlayToServerPacketWindowClickCursorItemPresentValue HashedSlot
		PlayToServerPacketWindowClickCursorItemPresentValue, err = PlayToServerPacketWindowClickCursorItemPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.CursorItem = &PlayToServerPacketWindowClickCursorItemPresentValue
	}
	return
}
func (ret PlayToServerPacketWindowClick) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	err = ret.StateId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Slot)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.MouseButton)
	if err != nil {
		return
	}
	err = ret.Mode.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.ChangedSlots)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToServerPacketWindowClickChangedSlots := range len(ret.ChangedSlots) {
		err = binary.Write(w, binary.BigEndian, ret.ChangedSlots[iPlayToServerPacketWindowClickChangedSlots].Location)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.ChangedSlots[iPlayToServerPacketWindowClickChangedSlots].Item != nil)
		if err != nil {
			return
		}
		if ret.ChangedSlots[iPlayToServerPacketWindowClickChangedSlots].Item != nil {
			err = (*ret.ChangedSlots[iPlayToServerPacketWindowClickChangedSlots].Item).Encode(w)
			if err != nil {
				return
			}
		}
	}
	err = binary.Write(w, binary.BigEndian, ret.CursorItem != nil)
	if err != nil {
		return
	}
	if ret.CursorItem != nil {
		err = (*ret.CursorItem).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientChatType struct {
	TranslationKey string
	Parameters     []PlayToClientChatTypeParameterType
	Style          nbt.Anon
}

func (_ PlayToClientChatType) Decode(r io.Reader) (ret PlayToClientChatType, err error) {
	ret.TranslationKey, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var lPlayToClientChatTypeParameters queser.VarInt
	lPlayToClientChatTypeParameters, err = lPlayToClientChatTypeParameters.Decode(r)
	if err != nil {
		return
	}
	ret.Parameters = []PlayToClientChatTypeParameterType{}
	for range lPlayToClientChatTypeParameters {
		var PlayToClientChatTypeParametersElement PlayToClientChatTypeParameterType
		PlayToClientChatTypeParametersElement, err = PlayToClientChatTypeParametersElement.Decode(r)
		if err != nil {
			return
		}
		ret.Parameters = append(ret.Parameters, PlayToClientChatTypeParametersElement)
	}
	ret.Style, err = ret.Style.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientChatType) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.TranslationKey)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Parameters)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientChatTypeParameters := range len(ret.Parameters) {
		err = ret.Parameters[iPlayToClientChatTypeParameters].Encode(w)
		if err != nil {
			return
		}
	}
	err = ret.Style.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientChatTypeParameterType struct {
	Val string
}

var PlayToClientChatTypeParameterTypeMap = map[queser.VarInt]string{0: "content", 1: "sender", 2: "target"}

func (_ PlayToClientChatTypeParameterType) Decode(r io.Reader) (ret PlayToClientChatTypeParameterType, err error) {
	var PlayToClientChatTypeParameterTypeKey queser.VarInt
	PlayToClientChatTypeParameterTypeKey, err = PlayToClientChatTypeParameterTypeKey.Decode(r)
	if err != nil {
		return
	}
	ret.Val, err = queser.ErroringIndex(PlayToClientChatTypeParameterTypeMap, PlayToClientChatTypeParameterTypeKey)
	if err != nil {
		return
	}
	return
}

var PlayToClientChatTypeParameterTypeReverseMap = map[string]queser.VarInt{"content": 0, "sender": 1, "target": 2}

func (ret PlayToClientChatTypeParameterType) Encode(w io.Writer) (err error) {
	var vPlayToClientChatTypeParameterType queser.VarInt
	vPlayToClientChatTypeParameterType, err = queser.ErroringIndex(PlayToClientChatTypeParameterTypeReverseMap, ret.Val)
	if err != nil {
		return
	}
	err = vPlayToClientChatTypeParameterType.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientChatTypes struct {
	Chat      PlayToClientChatType
	Narration PlayToClientChatType
}

func (_ PlayToClientChatTypes) Decode(r io.Reader) (ret PlayToClientChatTypes, err error) {
	ret.Chat, err = ret.Chat.Decode(r)
	if err != nil {
		return
	}
	ret.Narration, err = ret.Narration.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientChatTypes) Encode(w io.Writer) (err error) {
	err = ret.Chat.Encode(w)
	if err != nil {
		return
	}
	err = ret.Narration.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientChatTypesHolder struct {
	Val any
}

func (_ PlayToClientChatTypesHolder) Decode(r io.Reader) (ret PlayToClientChatTypesHolder, err error) {
	var PlayToClientChatTypesHolderId queser.VarInt
	PlayToClientChatTypesHolderId, err = PlayToClientChatTypesHolderId.Decode(r)
	if err != nil {
		return
	}
	if PlayToClientChatTypesHolderId != 0 {
		ret.Val = PlayToClientChatTypesHolderId
		return
	}
	var PlayToClientChatTypesHolderResult PlayToClientChatTypes
	PlayToClientChatTypesHolderResult, err = PlayToClientChatTypesHolderResult.Decode(r)
	if err != nil {
		return
	}
	ret.Val = PlayToClientChatTypesHolderResult
	return
}
func (ret PlayToClientChatTypesHolder) Encode(w io.Writer) (err error) {
	switch PlayToClientChatTypesHolderKnownType := ret.Val.(type) {
	case queser.VarInt:
		err = PlayToClientChatTypesHolderKnownType.Encode(w)
		if err != nil {
			return
		}
	case PlayToClientChatTypes:
		err = PlayToClientChatTypesHolderKnownType.Encode(w)
		if err != nil {
			return
		}
	default:
		err = queser.BadTypeError
	}
	return
}

type PlayToClientPositionUpdateRelatives struct {
	Val uint32
}

func (_ PlayToClientPositionUpdateRelatives) Decode(r io.Reader) (ret PlayToClientPositionUpdateRelatives, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Val)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPositionUpdateRelatives) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Val)
	if err != nil {
		return
	}
	return
}

type PlayToClientRecipeBookSetting struct {
	Open      bool
	Filtering bool
}

func (_ PlayToClientRecipeBookSetting) Decode(r io.Reader) (ret PlayToClientRecipeBookSetting, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Open)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Filtering)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientRecipeBookSetting) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Open)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Filtering)
	if err != nil {
		return
	}
	return
}

type PlayToClientRecipeDisplay struct {
	Type string
	Data any
}

var PlayToClientRecipeDisplayTypeMap = map[queser.VarInt]string{0: "crafting_shapeless", 1: "crafting_shaped", 2: "furnace", 3: "stonecutter", 4: "smithing"}

func (_ PlayToClientRecipeDisplay) Decode(r io.Reader) (ret PlayToClientRecipeDisplay, err error) {
	var PlayToClientRecipeDisplayTypeKey queser.VarInt
	PlayToClientRecipeDisplayTypeKey, err = PlayToClientRecipeDisplayTypeKey.Decode(r)
	if err != nil {
		return
	}
	ret.Type, err = queser.ErroringIndex(PlayToClientRecipeDisplayTypeMap, PlayToClientRecipeDisplayTypeKey)
	if err != nil {
		return
	}
	switch ret.Type {
	case "crafting_shaped":
		var PlayToClientRecipeDisplayDataTmp struct {
			Width           queser.VarInt
			Height          queser.VarInt
			Ingredients     []PlayToClientSlotDisplay
			Result          PlayToClientSlotDisplay
			CraftingStation PlayToClientSlotDisplay
		}
		PlayToClientRecipeDisplayDataTmp.Width, err = PlayToClientRecipeDisplayDataTmp.Width.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.Height, err = PlayToClientRecipeDisplayDataTmp.Height.Decode(r)
		if err != nil {
			return
		}
		var lPlayToClientRecipeDisplayDataIngredients queser.VarInt
		lPlayToClientRecipeDisplayDataIngredients, err = lPlayToClientRecipeDisplayDataIngredients.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.Ingredients = []PlayToClientSlotDisplay{}
		for range lPlayToClientRecipeDisplayDataIngredients {
			var PlayToClientRecipeDisplayDataIngredientsElement PlayToClientSlotDisplay
			PlayToClientRecipeDisplayDataIngredientsElement, err = PlayToClientRecipeDisplayDataIngredientsElement.Decode(r)
			if err != nil {
				return
			}
			PlayToClientRecipeDisplayDataTmp.Ingredients = append(PlayToClientRecipeDisplayDataTmp.Ingredients, PlayToClientRecipeDisplayDataIngredientsElement)
		}
		PlayToClientRecipeDisplayDataTmp.Result, err = PlayToClientRecipeDisplayDataTmp.Result.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.CraftingStation, err = PlayToClientRecipeDisplayDataTmp.CraftingStation.Decode(r)
		if err != nil {
			return
		}
		ret.Data = PlayToClientRecipeDisplayDataTmp
	case "crafting_shapeless":
		var PlayToClientRecipeDisplayDataTmp struct {
			Ingredients     []PlayToClientSlotDisplay
			Result          PlayToClientSlotDisplay
			CraftingStation PlayToClientSlotDisplay
		}
		var lPlayToClientRecipeDisplayDataIngredients queser.VarInt
		lPlayToClientRecipeDisplayDataIngredients, err = lPlayToClientRecipeDisplayDataIngredients.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.Ingredients = []PlayToClientSlotDisplay{}
		for range lPlayToClientRecipeDisplayDataIngredients {
			var PlayToClientRecipeDisplayDataIngredientsElement PlayToClientSlotDisplay
			PlayToClientRecipeDisplayDataIngredientsElement, err = PlayToClientRecipeDisplayDataIngredientsElement.Decode(r)
			if err != nil {
				return
			}
			PlayToClientRecipeDisplayDataTmp.Ingredients = append(PlayToClientRecipeDisplayDataTmp.Ingredients, PlayToClientRecipeDisplayDataIngredientsElement)
		}
		PlayToClientRecipeDisplayDataTmp.Result, err = PlayToClientRecipeDisplayDataTmp.Result.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.CraftingStation, err = PlayToClientRecipeDisplayDataTmp.CraftingStation.Decode(r)
		if err != nil {
			return
		}
		ret.Data = PlayToClientRecipeDisplayDataTmp
	case "furnace":
		var PlayToClientRecipeDisplayDataTmp struct {
			Ingredient      PlayToClientSlotDisplay
			Fuel            PlayToClientSlotDisplay
			Result          PlayToClientSlotDisplay
			CraftingStation PlayToClientSlotDisplay
			Duration        queser.VarInt
			Experience      float32
		}
		PlayToClientRecipeDisplayDataTmp.Ingredient, err = PlayToClientRecipeDisplayDataTmp.Ingredient.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.Fuel, err = PlayToClientRecipeDisplayDataTmp.Fuel.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.Result, err = PlayToClientRecipeDisplayDataTmp.Result.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.CraftingStation, err = PlayToClientRecipeDisplayDataTmp.CraftingStation.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.Duration, err = PlayToClientRecipeDisplayDataTmp.Duration.Decode(r)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientRecipeDisplayDataTmp.Experience)
		if err != nil {
			return
		}
		ret.Data = PlayToClientRecipeDisplayDataTmp
	case "smithing":
		var PlayToClientRecipeDisplayDataTmp struct {
			Template        PlayToClientSlotDisplay
			Base            PlayToClientSlotDisplay
			Addition        PlayToClientSlotDisplay
			Result          PlayToClientSlotDisplay
			CraftingStation PlayToClientSlotDisplay
		}
		PlayToClientRecipeDisplayDataTmp.Template, err = PlayToClientRecipeDisplayDataTmp.Template.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.Base, err = PlayToClientRecipeDisplayDataTmp.Base.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.Addition, err = PlayToClientRecipeDisplayDataTmp.Addition.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.Result, err = PlayToClientRecipeDisplayDataTmp.Result.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.CraftingStation, err = PlayToClientRecipeDisplayDataTmp.CraftingStation.Decode(r)
		if err != nil {
			return
		}
		ret.Data = PlayToClientRecipeDisplayDataTmp
	case "stonecutter":
		var PlayToClientRecipeDisplayDataTmp struct {
			Ingredient      PlayToClientSlotDisplay
			Result          PlayToClientSlotDisplay
			CraftingStation PlayToClientSlotDisplay
		}
		PlayToClientRecipeDisplayDataTmp.Ingredient, err = PlayToClientRecipeDisplayDataTmp.Ingredient.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.Result, err = PlayToClientRecipeDisplayDataTmp.Result.Decode(r)
		if err != nil {
			return
		}
		PlayToClientRecipeDisplayDataTmp.CraftingStation, err = PlayToClientRecipeDisplayDataTmp.CraftingStation.Decode(r)
		if err != nil {
			return
		}
		ret.Data = PlayToClientRecipeDisplayDataTmp
	}
	return
}

var PlayToClientRecipeDisplayTypeReverseMap = map[string]queser.VarInt{"crafting_shapeless": 0, "crafting_shaped": 1, "furnace": 2, "stonecutter": 3, "smithing": 4}

func (ret PlayToClientRecipeDisplay) Encode(w io.Writer) (err error) {
	var vPlayToClientRecipeDisplayType queser.VarInt
	vPlayToClientRecipeDisplayType, err = queser.ErroringIndex(PlayToClientRecipeDisplayTypeReverseMap, ret.Type)
	if err != nil {
		return
	}
	err = vPlayToClientRecipeDisplayType.Encode(w)
	if err != nil {
		return
	}
	switch ret.Type {
	case "crafting_shaped":
		PlayToClientRecipeDisplayData, ok := ret.Data.(struct {
			Width           queser.VarInt
			Height          queser.VarInt
			Ingredients     []PlayToClientSlotDisplay
			Result          PlayToClientSlotDisplay
			CraftingStation PlayToClientSlotDisplay
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientRecipeDisplayData.Width.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.Height.Encode(w)
		if err != nil {
			return
		}
		err = queser.VarInt(len(PlayToClientRecipeDisplayData.Ingredients)).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientRecipeDisplayDataIngredients := range len(PlayToClientRecipeDisplayData.Ingredients) {
			err = PlayToClientRecipeDisplayData.Ingredients[iPlayToClientRecipeDisplayDataIngredients].Encode(w)
			if err != nil {
				return
			}
		}
		err = PlayToClientRecipeDisplayData.Result.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.CraftingStation.Encode(w)
		if err != nil {
			return
		}
	case "crafting_shapeless":
		PlayToClientRecipeDisplayData, ok := ret.Data.(struct {
			Ingredients     []PlayToClientSlotDisplay
			Result          PlayToClientSlotDisplay
			CraftingStation PlayToClientSlotDisplay
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.VarInt(len(PlayToClientRecipeDisplayData.Ingredients)).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientRecipeDisplayDataIngredients := range len(PlayToClientRecipeDisplayData.Ingredients) {
			err = PlayToClientRecipeDisplayData.Ingredients[iPlayToClientRecipeDisplayDataIngredients].Encode(w)
			if err != nil {
				return
			}
		}
		err = PlayToClientRecipeDisplayData.Result.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.CraftingStation.Encode(w)
		if err != nil {
			return
		}
	case "furnace":
		PlayToClientRecipeDisplayData, ok := ret.Data.(struct {
			Ingredient      PlayToClientSlotDisplay
			Fuel            PlayToClientSlotDisplay
			Result          PlayToClientSlotDisplay
			CraftingStation PlayToClientSlotDisplay
			Duration        queser.VarInt
			Experience      float32
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientRecipeDisplayData.Ingredient.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.Fuel.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.Result.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.CraftingStation.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.Duration.Encode(w)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToClientRecipeDisplayData.Experience)
		if err != nil {
			return
		}
	case "smithing":
		PlayToClientRecipeDisplayData, ok := ret.Data.(struct {
			Template        PlayToClientSlotDisplay
			Base            PlayToClientSlotDisplay
			Addition        PlayToClientSlotDisplay
			Result          PlayToClientSlotDisplay
			CraftingStation PlayToClientSlotDisplay
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientRecipeDisplayData.Template.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.Base.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.Addition.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.Result.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.CraftingStation.Encode(w)
		if err != nil {
			return
		}
	case "stonecutter":
		PlayToClientRecipeDisplayData, ok := ret.Data.(struct {
			Ingredient      PlayToClientSlotDisplay
			Result          PlayToClientSlotDisplay
			CraftingStation PlayToClientSlotDisplay
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientRecipeDisplayData.Ingredient.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.Result.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientRecipeDisplayData.CraftingStation.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientSlotDisplay struct {
	Type string
	Data any
}

var PlayToClientSlotDisplayTypeMap = map[queser.VarInt]string{0: "empty", 1: "any_fuel", 2: "item", 3: "item_stack", 4: "tag", 5: "smithing_trim", 6: "with_remainder", 7: "composite"}

func (_ PlayToClientSlotDisplay) Decode(r io.Reader) (ret PlayToClientSlotDisplay, err error) {
	var PlayToClientSlotDisplayTypeKey queser.VarInt
	PlayToClientSlotDisplayTypeKey, err = PlayToClientSlotDisplayTypeKey.Decode(r)
	if err != nil {
		return
	}
	ret.Type, err = queser.ErroringIndex(PlayToClientSlotDisplayTypeMap, PlayToClientSlotDisplayTypeKey)
	if err != nil {
		return
	}
	switch ret.Type {
	case "any_fuel":
		var PlayToClientSlotDisplayDataTmp queser.Void
		PlayToClientSlotDisplayDataTmp, err = PlayToClientSlotDisplayDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = PlayToClientSlotDisplayDataTmp
	case "composite":
		var PlayToClientSlotDisplayDataTmp []PlayToClientSlotDisplay
		var lPlayToClientSlotDisplayData queser.VarInt
		lPlayToClientSlotDisplayData, err = lPlayToClientSlotDisplayData.Decode(r)
		if err != nil {
			return
		}
		PlayToClientSlotDisplayDataTmp = []PlayToClientSlotDisplay{}
		for range lPlayToClientSlotDisplayData {
			var PlayToClientSlotDisplayDataElement PlayToClientSlotDisplay
			PlayToClientSlotDisplayDataElement, err = PlayToClientSlotDisplayDataElement.Decode(r)
			if err != nil {
				return
			}
			PlayToClientSlotDisplayDataTmp = append(PlayToClientSlotDisplayDataTmp, PlayToClientSlotDisplayDataElement)
		}
		ret.Data = PlayToClientSlotDisplayDataTmp
	case "empty":
		var PlayToClientSlotDisplayDataTmp queser.Void
		PlayToClientSlotDisplayDataTmp, err = PlayToClientSlotDisplayDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = PlayToClientSlotDisplayDataTmp
	case "item":
		var PlayToClientSlotDisplayDataTmp queser.VarInt
		PlayToClientSlotDisplayDataTmp, err = PlayToClientSlotDisplayDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = PlayToClientSlotDisplayDataTmp
	case "item_stack":
		var PlayToClientSlotDisplayDataTmp Slot
		PlayToClientSlotDisplayDataTmp, err = PlayToClientSlotDisplayDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = PlayToClientSlotDisplayDataTmp
	case "smithing_trim":
		var PlayToClientSlotDisplayDataTmp struct {
			Base     PlayToClientSlotDisplay
			Material PlayToClientSlotDisplay
			Pattern  any
		}
		PlayToClientSlotDisplayDataTmp.Base, err = PlayToClientSlotDisplayDataTmp.Base.Decode(r)
		if err != nil {
			return
		}
		PlayToClientSlotDisplayDataTmp.Material, err = PlayToClientSlotDisplayDataTmp.Material.Decode(r)
		if err != nil {
			return
		}
		var PlayToClientSlotDisplayDataPatternId queser.VarInt
		PlayToClientSlotDisplayDataPatternId, err = PlayToClientSlotDisplayDataPatternId.Decode(r)
		if err != nil {
			return
		}
		if PlayToClientSlotDisplayDataPatternId != 0 {
			PlayToClientSlotDisplayDataTmp.Pattern = PlayToClientSlotDisplayDataPatternId
			return
		}
		var PlayToClientSlotDisplayDataPatternResult ArmorTrimPattern
		PlayToClientSlotDisplayDataPatternResult, err = PlayToClientSlotDisplayDataPatternResult.Decode(r)
		if err != nil {
			return
		}
		PlayToClientSlotDisplayDataTmp.Pattern = PlayToClientSlotDisplayDataPatternResult
		ret.Data = PlayToClientSlotDisplayDataTmp
	case "tag":
		var PlayToClientSlotDisplayDataTmp string
		PlayToClientSlotDisplayDataTmp, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Data = PlayToClientSlotDisplayDataTmp
	case "with_remainder":
		var PlayToClientSlotDisplayDataTmp struct {
			Input     PlayToClientSlotDisplay
			Remainder PlayToClientSlotDisplay
		}
		PlayToClientSlotDisplayDataTmp.Input, err = PlayToClientSlotDisplayDataTmp.Input.Decode(r)
		if err != nil {
			return
		}
		PlayToClientSlotDisplayDataTmp.Remainder, err = PlayToClientSlotDisplayDataTmp.Remainder.Decode(r)
		if err != nil {
			return
		}
		ret.Data = PlayToClientSlotDisplayDataTmp
	}
	return
}

var PlayToClientSlotDisplayTypeReverseMap = map[string]queser.VarInt{"empty": 0, "any_fuel": 1, "item": 2, "item_stack": 3, "tag": 4, "smithing_trim": 5, "with_remainder": 6, "composite": 7}

func (ret PlayToClientSlotDisplay) Encode(w io.Writer) (err error) {
	var vPlayToClientSlotDisplayType queser.VarInt
	vPlayToClientSlotDisplayType, err = queser.ErroringIndex(PlayToClientSlotDisplayTypeReverseMap, ret.Type)
	if err != nil {
		return
	}
	err = vPlayToClientSlotDisplayType.Encode(w)
	if err != nil {
		return
	}
	switch ret.Type {
	case "any_fuel":
		PlayToClientSlotDisplayData, ok := ret.Data.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientSlotDisplayData.Encode(w)
		if err != nil {
			return
		}
	case "composite":
		PlayToClientSlotDisplayData, ok := ret.Data.([]PlayToClientSlotDisplay)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.VarInt(len(PlayToClientSlotDisplayData)).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientSlotDisplayData := range len(PlayToClientSlotDisplayData) {
			err = PlayToClientSlotDisplayData[iPlayToClientSlotDisplayData].Encode(w)
			if err != nil {
				return
			}
		}
	case "empty":
		PlayToClientSlotDisplayData, ok := ret.Data.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientSlotDisplayData.Encode(w)
		if err != nil {
			return
		}
	case "item":
		PlayToClientSlotDisplayData, ok := ret.Data.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientSlotDisplayData.Encode(w)
		if err != nil {
			return
		}
	case "item_stack":
		PlayToClientSlotDisplayData, ok := ret.Data.(Slot)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientSlotDisplayData.Encode(w)
		if err != nil {
			return
		}
	case "smithing_trim":
		PlayToClientSlotDisplayData, ok := ret.Data.(struct {
			Base     PlayToClientSlotDisplay
			Material PlayToClientSlotDisplay
			Pattern  any
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientSlotDisplayData.Base.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientSlotDisplayData.Material.Encode(w)
		if err != nil {
			return
		}
		switch PlayToClientSlotDisplayDataPatternKnownType := PlayToClientSlotDisplayData.Pattern.(type) {
		case queser.VarInt:
			err = PlayToClientSlotDisplayDataPatternKnownType.Encode(w)
			if err != nil {
				return
			}
		case ArmorTrimPattern:
			err = PlayToClientSlotDisplayDataPatternKnownType.Encode(w)
			if err != nil {
				return
			}
		default:
			err = queser.BadTypeError
		}
	case "tag":
		PlayToClientSlotDisplayData, ok := ret.Data.(string)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.EncodeString(w, PlayToClientSlotDisplayData)
		if err != nil {
			return
		}
	case "with_remainder":
		PlayToClientSlotDisplayData, ok := ret.Data.(struct {
			Input     PlayToClientSlotDisplay
			Remainder PlayToClientSlotDisplay
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientSlotDisplayData.Input.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientSlotDisplayData.Remainder.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientSpawnInfo struct {
	Dimension        queser.VarInt
	Name             string
	HashedSeed       int64
	Gamemode         string
	PreviousGamemode uint8
	IsDebug          bool
	IsFlat           bool
	Death            *struct {
		DimensionName string
		Location      Position
	}
	PortalCooldown queser.VarInt
	SeaLevel       queser.VarInt
}

var PlayToClientSpawnInfoGamemodeMap = map[int8]string{0: "survival", 1: "creative", 2: "adventure", 3: "spectator"}

func (_ PlayToClientSpawnInfo) Decode(r io.Reader) (ret PlayToClientSpawnInfo, err error) {
	ret.Dimension, err = ret.Dimension.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.HashedSeed)
	if err != nil {
		return
	}
	var PlayToClientSpawnInfoGamemodeKey int8
	err = binary.Read(r, binary.BigEndian, &PlayToClientSpawnInfoGamemodeKey)
	if err != nil {
		return
	}
	ret.Gamemode, err = queser.ErroringIndex(PlayToClientSpawnInfoGamemodeMap, PlayToClientSpawnInfoGamemodeKey)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.PreviousGamemode)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IsDebug)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IsFlat)
	if err != nil {
		return
	}
	var PlayToClientSpawnInfoDeathPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientSpawnInfoDeathPresent)
	if err != nil {
		return
	}
	if PlayToClientSpawnInfoDeathPresent {
		var PlayToClientSpawnInfoDeathPresentValue struct {
			DimensionName string
			Location      Position
		}
		PlayToClientSpawnInfoDeathPresentValue.DimensionName, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		PlayToClientSpawnInfoDeathPresentValue.Location, err = PlayToClientSpawnInfoDeathPresentValue.Location.Decode(r)
		if err != nil {
			return
		}
		ret.Death = &PlayToClientSpawnInfoDeathPresentValue
	}
	ret.PortalCooldown, err = ret.PortalCooldown.Decode(r)
	if err != nil {
		return
	}
	ret.SeaLevel, err = ret.SeaLevel.Decode(r)
	if err != nil {
		return
	}
	return
}

var PlayToClientSpawnInfoGamemodeReverseMap = map[string]int8{"survival": 0, "creative": 1, "adventure": 2, "spectator": 3}

func (ret PlayToClientSpawnInfo) Encode(w io.Writer) (err error) {
	err = ret.Dimension.Encode(w)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Name)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.HashedSeed)
	if err != nil {
		return
	}
	var vPlayToClientSpawnInfoGamemode int8
	vPlayToClientSpawnInfoGamemode, err = queser.ErroringIndex(PlayToClientSpawnInfoGamemodeReverseMap, ret.Gamemode)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, vPlayToClientSpawnInfoGamemode)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.PreviousGamemode)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IsDebug)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IsFlat)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Death != nil)
	if err != nil {
		return
	}
	if ret.Death != nil {
		err = queser.EncodeString(w, (*ret.Death).DimensionName)
		if err != nil {
			return
		}
		err = (*ret.Death).Location.Encode(w)
		if err != nil {
			return
		}
	}
	err = ret.PortalCooldown.Encode(w)
	if err != nil {
		return
	}
	err = ret.SeaLevel.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacket struct {
	Name   string
	Params any
}

var PlayToClientPacketNameMap = map[queser.VarInt]string{0x00: "bundle_delimiter", 0x01: "spawn_entity", 0x02: "animation", 0x03: "statistics", 0x04: "acknowledge_player_digging", 0x05: "block_break_animation", 0x06: "tile_entity_data", 0x07: "block_action", 0x08: "block_change", 0x09: "boss_bar", 0x0a: "difficulty", 0x0b: "chunk_batch_finished", 0x0c: "chunk_batch_start", 0x0d: "chunk_biomes", 0x0e: "clear_titles", 0x0f: "tab_complete", 0x10: "declare_commands", 0x11: "close_window", 0x12: "window_items", 0x13: "craft_progress_bar", 0x14: "set_slot", 0x15: "cookie_request", 0x16: "set_cooldown", 0x17: "chat_suggestions", 0x18: "custom_payload", 0x19: "damage_event", 0x1a: "debug_sample", 0x1b: "hide_message", 0x1c: "kick_disconnect", 0x1d: "profileless_chat", 0x1e: "entity_status", 0x1f: "sync_entity_position", 0x20: "explosion", 0x21: "unload_chunk", 0x22: "game_state_change", 0x23: "open_horse_window", 0x24: "hurt_animation", 0x25: "initialize_world_border", 0x26: "keep_alive", 0x27: "map_chunk", 0x28: "world_event", 0x29: "world_particles", 0x2a: "update_light", 0x2b: "login", 0x2c: "map", 0x2d: "trade_list", 0x2e: "rel_entity_move", 0x2f: "entity_move_look", 0x30: "move_minecart", 0x31: "entity_look", 0x32: "vehicle_move", 0x33: "open_book", 0x34: "open_window", 0x35: "open_sign_entity", 0x36: "ping", 0x37: "ping_response", 0x38: "craft_recipe_response", 0x39: "abilities", 0x3a: "player_chat", 0x3b: "end_combat_event", 0x3c: "enter_combat_event", 0x3d: "death_combat_event", 0x3e: "player_remove", 0x3f: "player_info", 0x40: "face_player", 0x41: "position", 0x42: "player_rotation", 0x43: "recipe_book_add", 0x44: "recipe_book_remove", 0x45: "recipe_book_settings", 0x46: "entity_destroy", 0x47: "remove_entity_effect", 0x48: "reset_score", 0x49: "remove_resource_pack", 0x4a: "add_resource_pack", 0x4b: "respawn", 0x4c: "entity_head_rotation", 0x4d: "multi_block_change", 0x4e: "select_advancement_tab", 0x4f: "server_data", 0x50: "action_bar", 0x51: "world_border_center", 0x52: "world_border_lerp_size", 0x53: "world_border_size", 0x54: "world_border_warning_delay", 0x55: "world_border_warning_reach", 0x56: "camera", 0x57: "update_view_position", 0x58: "update_view_distance", 0x59: "set_cursor_item", 0x5a: "spawn_position", 0x5b: "scoreboard_display_objective", 0x5c: "entity_metadata", 0x5d: "attach_entity", 0x5e: "entity_velocity", 0x5f: "entity_equipment", 0x60: "experience", 0x61: "update_health", 0x62: "held_item_slot", 0x63: "scoreboard_objective", 0x64: "set_passengers", 0x65: "set_player_inventory", 0x66: "teams", 0x67: "scoreboard_score", 0x68: "simulation_distance", 0x69: "set_title_subtitle", 0x6a: "update_time", 0x6b: "set_title_text", 0x6c: "set_title_time", 0x6d: "entity_sound_effect", 0x6e: "sound_effect", 0x6f: "start_configuration", 0x70: "stop_sound", 0x71: "store_cookie", 0x72: "system_chat", 0x73: "playerlist_header", 0x74: "nbt_query_response", 0x75: "collect", 0x76: "entity_teleport", 0x77: "test_instance_block_status", 0x78: "set_ticking_state", 0x79: "step_tick", 0x7a: "transfer", 0x7b: "advancements", 0x7c: "entity_update_attributes", 0x7d: "entity_effect", 0x7e: "declare_recipes", 0x7f: "tags", 0x80: "set_projectile_power", 0x81: "custom_report_details", 0x82: "server_links", 0x83: "tracked_waypoint", 0x84: "clear_dialog", 0x85: "show_dialog"}

func (_ PlayToClientPacket) Decode(r io.Reader) (ret PlayToClientPacket, err error) {
	var PlayToClientPacketNameKey queser.VarInt
	PlayToClientPacketNameKey, err = PlayToClientPacketNameKey.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.ErroringIndex(PlayToClientPacketNameMap, PlayToClientPacketNameKey)
	if err != nil {
		return
	}
	switch ret.Name {
	case "abilities":
		var PlayToClientPacketParamsTmp PlayToClientPacketAbilities
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "acknowledge_player_digging":
		var PlayToClientPacketParamsTmp PlayToClientPacketAcknowledgePlayerDigging
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "action_bar":
		var PlayToClientPacketParamsTmp PlayToClientPacketActionBar
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "add_resource_pack":
		var PlayToClientPacketParamsTmp PacketCommonAddResourcePack
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "advancements":
		var PlayToClientPacketParamsTmp PlayToClientPacketAdvancements
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "animation":
		var PlayToClientPacketParamsTmp PlayToClientPacketAnimation
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "attach_entity":
		var PlayToClientPacketParamsTmp PlayToClientPacketAttachEntity
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "block_action":
		var PlayToClientPacketParamsTmp PlayToClientPacketBlockAction
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "block_break_animation":
		var PlayToClientPacketParamsTmp PlayToClientPacketBlockBreakAnimation
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "block_change":
		var PlayToClientPacketParamsTmp PlayToClientPacketBlockChange
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "boss_bar":
		var PlayToClientPacketParamsTmp PlayToClientPacketBossBar
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "bundle_delimiter":
		var PlayToClientPacketParamsTmp queser.Void
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "camera":
		var PlayToClientPacketParamsTmp PlayToClientPacketCamera
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "chat_suggestions":
		var PlayToClientPacketParamsTmp PlayToClientPacketChatSuggestions
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "chunk_batch_finished":
		var PlayToClientPacketParamsTmp PlayToClientPacketChunkBatchFinished
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "chunk_batch_start":
		var PlayToClientPacketParamsTmp PlayToClientPacketChunkBatchStart
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "chunk_biomes":
		var PlayToClientPacketParamsTmp PlayToClientPacketChunkBiomes
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "clear_dialog":
		var PlayToClientPacketParamsTmp PacketCommonClearDialog
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "clear_titles":
		var PlayToClientPacketParamsTmp PlayToClientPacketClearTitles
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "close_window":
		var PlayToClientPacketParamsTmp PlayToClientPacketCloseWindow
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "collect":
		var PlayToClientPacketParamsTmp PlayToClientPacketCollect
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "cookie_request":
		var PlayToClientPacketParamsTmp PacketCommonCookieRequest
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "craft_progress_bar":
		var PlayToClientPacketParamsTmp PlayToClientPacketCraftProgressBar
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "craft_recipe_response":
		var PlayToClientPacketParamsTmp PlayToClientPacketCraftRecipeResponse
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "custom_payload":
		var PlayToClientPacketParamsTmp PlayToClientPacketCustomPayload
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "custom_report_details":
		var PlayToClientPacketParamsTmp PacketCommonCustomReportDetails
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "damage_event":
		var PlayToClientPacketParamsTmp PlayToClientPacketDamageEvent
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "death_combat_event":
		var PlayToClientPacketParamsTmp PlayToClientPacketDeathCombatEvent
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "debug_sample":
		var PlayToClientPacketParamsTmp PlayToClientPacketDebugSample
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "declare_commands":
		var PlayToClientPacketParamsTmp PlayToClientPacketDeclareCommands
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "declare_recipes":
		var PlayToClientPacketParamsTmp PlayToClientPacketDeclareRecipes
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "difficulty":
		var PlayToClientPacketParamsTmp PlayToClientPacketDifficulty
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "end_combat_event":
		var PlayToClientPacketParamsTmp PlayToClientPacketEndCombatEvent
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "enter_combat_event":
		var PlayToClientPacketParamsTmp PlayToClientPacketEnterCombatEvent
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_destroy":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntityDestroy
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_effect":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntityEffect
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_equipment":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntityEquipment
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_head_rotation":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntityHeadRotation
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_look":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntityLook
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_metadata":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntityMetadata
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_move_look":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntityMoveLook
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_sound_effect":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntitySoundEffect
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_status":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntityStatus
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_teleport":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntityTeleport
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_update_attributes":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntityUpdateAttributes
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "entity_velocity":
		var PlayToClientPacketParamsTmp PlayToClientPacketEntityVelocity
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "experience":
		var PlayToClientPacketParamsTmp PlayToClientPacketExperience
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "explosion":
		var PlayToClientPacketParamsTmp PlayToClientPacketExplosion
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "face_player":
		var PlayToClientPacketParamsTmp PlayToClientPacketFacePlayer
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "game_state_change":
		var PlayToClientPacketParamsTmp PlayToClientPacketGameStateChange
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "held_item_slot":
		var PlayToClientPacketParamsTmp PlayToClientPacketHeldItemSlot
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "hide_message":
		var PlayToClientPacketParamsTmp PlayToClientPacketHideMessage
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "hurt_animation":
		var PlayToClientPacketParamsTmp PlayToClientPacketHurtAnimation
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "initialize_world_border":
		var PlayToClientPacketParamsTmp PlayToClientPacketInitializeWorldBorder
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "keep_alive":
		var PlayToClientPacketParamsTmp PlayToClientPacketKeepAlive
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "kick_disconnect":
		var PlayToClientPacketParamsTmp PlayToClientPacketKickDisconnect
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "login":
		var PlayToClientPacketParamsTmp PlayToClientPacketLogin
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "map":
		var PlayToClientPacketParamsTmp PlayToClientPacketMap
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "map_chunk":
		var PlayToClientPacketParamsTmp PlayToClientPacketMapChunk
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "move_minecart":
		var PlayToClientPacketParamsTmp PlayToClientPacketMoveMinecart
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "multi_block_change":
		var PlayToClientPacketParamsTmp PlayToClientPacketMultiBlockChange
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "nbt_query_response":
		var PlayToClientPacketParamsTmp PlayToClientPacketNbtQueryResponse
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "open_book":
		var PlayToClientPacketParamsTmp PlayToClientPacketOpenBook
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "open_horse_window":
		var PlayToClientPacketParamsTmp PlayToClientPacketOpenHorseWindow
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "open_sign_entity":
		var PlayToClientPacketParamsTmp PlayToClientPacketOpenSignEntity
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "open_window":
		var PlayToClientPacketParamsTmp PlayToClientPacketOpenWindow
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "ping":
		var PlayToClientPacketParamsTmp PlayToClientPacketPing
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "ping_response":
		var PlayToClientPacketParamsTmp PlayToClientPacketPingResponse
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "player_chat":
		var PlayToClientPacketParamsTmp PlayToClientPacketPlayerChat
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "player_info":
		var PlayToClientPacketParamsTmp PlayToClientPacketPlayerInfo
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "player_remove":
		var PlayToClientPacketParamsTmp PlayToClientPacketPlayerRemove
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "player_rotation":
		var PlayToClientPacketParamsTmp PlayToClientPacketPlayerRotation
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "playerlist_header":
		var PlayToClientPacketParamsTmp PlayToClientPacketPlayerlistHeader
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "position":
		var PlayToClientPacketParamsTmp PlayToClientPacketPosition
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "profileless_chat":
		var PlayToClientPacketParamsTmp PlayToClientPacketProfilelessChat
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "recipe_book_add":
		var PlayToClientPacketParamsTmp PlayToClientPacketRecipeBookAdd
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "recipe_book_remove":
		var PlayToClientPacketParamsTmp PlayToClientPacketRecipeBookRemove
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "recipe_book_settings":
		var PlayToClientPacketParamsTmp PlayToClientPacketRecipeBookSettings
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "rel_entity_move":
		var PlayToClientPacketParamsTmp PlayToClientPacketRelEntityMove
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "remove_entity_effect":
		var PlayToClientPacketParamsTmp PlayToClientPacketRemoveEntityEffect
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "remove_resource_pack":
		var PlayToClientPacketParamsTmp PacketCommonRemoveResourcePack
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "reset_score":
		var PlayToClientPacketParamsTmp PlayToClientPacketResetScore
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "respawn":
		var PlayToClientPacketParamsTmp PlayToClientPacketRespawn
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "scoreboard_display_objective":
		var PlayToClientPacketParamsTmp PlayToClientPacketScoreboardDisplayObjective
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "scoreboard_objective":
		var PlayToClientPacketParamsTmp PlayToClientPacketScoreboardObjective
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "scoreboard_score":
		var PlayToClientPacketParamsTmp PlayToClientPacketScoreboardScore
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "select_advancement_tab":
		var PlayToClientPacketParamsTmp PlayToClientPacketSelectAdvancementTab
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "server_data":
		var PlayToClientPacketParamsTmp PlayToClientPacketServerData
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "server_links":
		var PlayToClientPacketParamsTmp PacketCommonServerLinks
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "set_cooldown":
		var PlayToClientPacketParamsTmp PlayToClientPacketSetCooldown
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "set_cursor_item":
		var PlayToClientPacketParamsTmp PlayToClientPacketSetCursorItem
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "set_passengers":
		var PlayToClientPacketParamsTmp PlayToClientPacketSetPassengers
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "set_player_inventory":
		var PlayToClientPacketParamsTmp PlayToClientPacketSetPlayerInventory
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "set_projectile_power":
		var PlayToClientPacketParamsTmp PlayToClientPacketSetProjectilePower
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "set_slot":
		var PlayToClientPacketParamsTmp PlayToClientPacketSetSlot
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "set_ticking_state":
		var PlayToClientPacketParamsTmp PlayToClientPacketSetTickingState
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "set_title_subtitle":
		var PlayToClientPacketParamsTmp PlayToClientPacketSetTitleSubtitle
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "set_title_text":
		var PlayToClientPacketParamsTmp PlayToClientPacketSetTitleText
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "set_title_time":
		var PlayToClientPacketParamsTmp PlayToClientPacketSetTitleTime
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "show_dialog":
		var PlayToClientPacketParamsTmp PlayToClientPacketShowDialog
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "simulation_distance":
		var PlayToClientPacketParamsTmp PlayToClientPacketSimulationDistance
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "sound_effect":
		var PlayToClientPacketParamsTmp PlayToClientPacketSoundEffect
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "spawn_entity":
		var PlayToClientPacketParamsTmp PlayToClientPacketSpawnEntity
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "spawn_position":
		var PlayToClientPacketParamsTmp PlayToClientPacketSpawnPosition
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "start_configuration":
		var PlayToClientPacketParamsTmp PlayToClientPacketStartConfiguration
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "statistics":
		var PlayToClientPacketParamsTmp PlayToClientPacketStatistics
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "step_tick":
		var PlayToClientPacketParamsTmp PlayToClientPacketStepTick
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "stop_sound":
		var PlayToClientPacketParamsTmp PlayToClientPacketStopSound
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "store_cookie":
		var PlayToClientPacketParamsTmp PacketCommonStoreCookie
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "sync_entity_position":
		var PlayToClientPacketParamsTmp PlayToClientPacketSyncEntityPosition
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "system_chat":
		var PlayToClientPacketParamsTmp PlayToClientPacketSystemChat
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "tab_complete":
		var PlayToClientPacketParamsTmp PlayToClientPacketTabComplete
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "tags":
		var PlayToClientPacketParamsTmp PlayToClientPacketTags
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "teams":
		var PlayToClientPacketParamsTmp PlayToClientPacketTeams
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "test_instance_block_status":
		var PlayToClientPacketParamsTmp PlayToClientPacketTestInstanceBlockStatus
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "tile_entity_data":
		var PlayToClientPacketParamsTmp PlayToClientPacketTileEntityData
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "tracked_waypoint":
		var PlayToClientPacketParamsTmp PlayToClientPacketTrackedWaypoint
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "trade_list":
		var PlayToClientPacketParamsTmp PlayToClientPacketTradeList
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "transfer":
		var PlayToClientPacketParamsTmp PacketCommonTransfer
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "unload_chunk":
		var PlayToClientPacketParamsTmp PlayToClientPacketUnloadChunk
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "update_health":
		var PlayToClientPacketParamsTmp PlayToClientPacketUpdateHealth
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "update_light":
		var PlayToClientPacketParamsTmp PlayToClientPacketUpdateLight
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "update_time":
		var PlayToClientPacketParamsTmp PlayToClientPacketUpdateTime
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "update_view_distance":
		var PlayToClientPacketParamsTmp PlayToClientPacketUpdateViewDistance
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "update_view_position":
		var PlayToClientPacketParamsTmp PlayToClientPacketUpdateViewPosition
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "vehicle_move":
		var PlayToClientPacketParamsTmp PlayToClientPacketVehicleMove
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "window_items":
		var PlayToClientPacketParamsTmp PlayToClientPacketWindowItems
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "world_border_center":
		var PlayToClientPacketParamsTmp PlayToClientPacketWorldBorderCenter
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "world_border_lerp_size":
		var PlayToClientPacketParamsTmp PlayToClientPacketWorldBorderLerpSize
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "world_border_size":
		var PlayToClientPacketParamsTmp PlayToClientPacketWorldBorderSize
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "world_border_warning_delay":
		var PlayToClientPacketParamsTmp PlayToClientPacketWorldBorderWarningDelay
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "world_border_warning_reach":
		var PlayToClientPacketParamsTmp PlayToClientPacketWorldBorderWarningReach
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "world_event":
		var PlayToClientPacketParamsTmp PlayToClientPacketWorldEvent
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	case "world_particles":
		var PlayToClientPacketParamsTmp PlayToClientPacketWorldParticles
		PlayToClientPacketParamsTmp, err = PlayToClientPacketParamsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Params = PlayToClientPacketParamsTmp
	}
	return
}

var PlayToClientPacketNameReverseMap = map[string]queser.VarInt{"bundle_delimiter": 0x00, "spawn_entity": 0x01, "animation": 0x02, "statistics": 0x03, "acknowledge_player_digging": 0x04, "block_break_animation": 0x05, "tile_entity_data": 0x06, "block_action": 0x07, "block_change": 0x08, "boss_bar": 0x09, "difficulty": 0x0a, "chunk_batch_finished": 0x0b, "chunk_batch_start": 0x0c, "chunk_biomes": 0x0d, "clear_titles": 0x0e, "tab_complete": 0x0f, "declare_commands": 0x10, "close_window": 0x11, "window_items": 0x12, "craft_progress_bar": 0x13, "set_slot": 0x14, "cookie_request": 0x15, "set_cooldown": 0x16, "chat_suggestions": 0x17, "custom_payload": 0x18, "damage_event": 0x19, "debug_sample": 0x1a, "hide_message": 0x1b, "kick_disconnect": 0x1c, "profileless_chat": 0x1d, "entity_status": 0x1e, "sync_entity_position": 0x1f, "explosion": 0x20, "unload_chunk": 0x21, "game_state_change": 0x22, "open_horse_window": 0x23, "hurt_animation": 0x24, "initialize_world_border": 0x25, "keep_alive": 0x26, "map_chunk": 0x27, "world_event": 0x28, "world_particles": 0x29, "update_light": 0x2a, "login": 0x2b, "map": 0x2c, "trade_list": 0x2d, "rel_entity_move": 0x2e, "entity_move_look": 0x2f, "move_minecart": 0x30, "entity_look": 0x31, "vehicle_move": 0x32, "open_book": 0x33, "open_window": 0x34, "open_sign_entity": 0x35, "ping": 0x36, "ping_response": 0x37, "craft_recipe_response": 0x38, "abilities": 0x39, "player_chat": 0x3a, "end_combat_event": 0x3b, "enter_combat_event": 0x3c, "death_combat_event": 0x3d, "player_remove": 0x3e, "player_info": 0x3f, "face_player": 0x40, "position": 0x41, "player_rotation": 0x42, "recipe_book_add": 0x43, "recipe_book_remove": 0x44, "recipe_book_settings": 0x45, "entity_destroy": 0x46, "remove_entity_effect": 0x47, "reset_score": 0x48, "remove_resource_pack": 0x49, "add_resource_pack": 0x4a, "respawn": 0x4b, "entity_head_rotation": 0x4c, "multi_block_change": 0x4d, "select_advancement_tab": 0x4e, "server_data": 0x4f, "action_bar": 0x50, "world_border_center": 0x51, "world_border_lerp_size": 0x52, "world_border_size": 0x53, "world_border_warning_delay": 0x54, "world_border_warning_reach": 0x55, "camera": 0x56, "update_view_position": 0x57, "update_view_distance": 0x58, "set_cursor_item": 0x59, "spawn_position": 0x5a, "scoreboard_display_objective": 0x5b, "entity_metadata": 0x5c, "attach_entity": 0x5d, "entity_velocity": 0x5e, "entity_equipment": 0x5f, "experience": 0x60, "update_health": 0x61, "held_item_slot": 0x62, "scoreboard_objective": 0x63, "set_passengers": 0x64, "set_player_inventory": 0x65, "teams": 0x66, "scoreboard_score": 0x67, "simulation_distance": 0x68, "set_title_subtitle": 0x69, "update_time": 0x6a, "set_title_text": 0x6b, "set_title_time": 0x6c, "entity_sound_effect": 0x6d, "sound_effect": 0x6e, "start_configuration": 0x6f, "stop_sound": 0x70, "store_cookie": 0x71, "system_chat": 0x72, "playerlist_header": 0x73, "nbt_query_response": 0x74, "collect": 0x75, "entity_teleport": 0x76, "test_instance_block_status": 0x77, "set_ticking_state": 0x78, "step_tick": 0x79, "transfer": 0x7a, "advancements": 0x7b, "entity_update_attributes": 0x7c, "entity_effect": 0x7d, "declare_recipes": 0x7e, "tags": 0x7f, "set_projectile_power": 0x80, "custom_report_details": 0x81, "server_links": 0x82, "tracked_waypoint": 0x83, "clear_dialog": 0x84, "show_dialog": 0x85}

func (ret PlayToClientPacket) Encode(w io.Writer) (err error) {
	var vPlayToClientPacketName queser.VarInt
	vPlayToClientPacketName, err = queser.ErroringIndex(PlayToClientPacketNameReverseMap, ret.Name)
	if err != nil {
		return
	}
	err = vPlayToClientPacketName.Encode(w)
	if err != nil {
		return
	}
	switch ret.Name {
	case "abilities":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketAbilities)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "acknowledge_player_digging":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketAcknowledgePlayerDigging)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "action_bar":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketActionBar)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "add_resource_pack":
		PlayToClientPacketParams, ok := ret.Params.(PacketCommonAddResourcePack)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "advancements":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketAdvancements)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "animation":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketAnimation)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "attach_entity":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketAttachEntity)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "block_action":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketBlockAction)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "block_break_animation":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketBlockBreakAnimation)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "block_change":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketBlockChange)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "boss_bar":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketBossBar)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "bundle_delimiter":
		PlayToClientPacketParams, ok := ret.Params.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "camera":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketCamera)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "chat_suggestions":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketChatSuggestions)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "chunk_batch_finished":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketChunkBatchFinished)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "chunk_batch_start":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketChunkBatchStart)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "chunk_biomes":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketChunkBiomes)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "clear_dialog":
		PlayToClientPacketParams, ok := ret.Params.(PacketCommonClearDialog)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "clear_titles":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketClearTitles)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "close_window":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketCloseWindow)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "collect":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketCollect)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "cookie_request":
		PlayToClientPacketParams, ok := ret.Params.(PacketCommonCookieRequest)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "craft_progress_bar":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketCraftProgressBar)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "craft_recipe_response":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketCraftRecipeResponse)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "custom_payload":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketCustomPayload)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "custom_report_details":
		PlayToClientPacketParams, ok := ret.Params.(PacketCommonCustomReportDetails)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "damage_event":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketDamageEvent)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "death_combat_event":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketDeathCombatEvent)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "debug_sample":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketDebugSample)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "declare_commands":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketDeclareCommands)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "declare_recipes":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketDeclareRecipes)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "difficulty":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketDifficulty)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "end_combat_event":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEndCombatEvent)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "enter_combat_event":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEnterCombatEvent)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_destroy":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntityDestroy)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_effect":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntityEffect)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_equipment":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntityEquipment)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_head_rotation":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntityHeadRotation)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_look":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntityLook)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_metadata":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntityMetadata)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_move_look":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntityMoveLook)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_sound_effect":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntitySoundEffect)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_status":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntityStatus)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_teleport":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntityTeleport)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_update_attributes":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntityUpdateAttributes)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "entity_velocity":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketEntityVelocity)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "experience":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketExperience)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "explosion":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketExplosion)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "face_player":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketFacePlayer)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "game_state_change":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketGameStateChange)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "held_item_slot":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketHeldItemSlot)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "hide_message":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketHideMessage)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "hurt_animation":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketHurtAnimation)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "initialize_world_border":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketInitializeWorldBorder)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "keep_alive":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketKeepAlive)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "kick_disconnect":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketKickDisconnect)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "login":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketLogin)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "map":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketMap)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "map_chunk":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketMapChunk)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "move_minecart":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketMoveMinecart)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "multi_block_change":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketMultiBlockChange)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "nbt_query_response":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketNbtQueryResponse)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "open_book":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketOpenBook)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "open_horse_window":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketOpenHorseWindow)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "open_sign_entity":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketOpenSignEntity)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "open_window":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketOpenWindow)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "ping":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketPing)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "ping_response":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketPingResponse)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "player_chat":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketPlayerChat)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "player_info":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketPlayerInfo)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "player_remove":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketPlayerRemove)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "player_rotation":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketPlayerRotation)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "playerlist_header":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketPlayerlistHeader)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "position":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketPosition)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "profileless_chat":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketProfilelessChat)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "recipe_book_add":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketRecipeBookAdd)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "recipe_book_remove":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketRecipeBookRemove)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "recipe_book_settings":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketRecipeBookSettings)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "rel_entity_move":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketRelEntityMove)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "remove_entity_effect":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketRemoveEntityEffect)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "remove_resource_pack":
		PlayToClientPacketParams, ok := ret.Params.(PacketCommonRemoveResourcePack)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "reset_score":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketResetScore)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "respawn":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketRespawn)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "scoreboard_display_objective":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketScoreboardDisplayObjective)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "scoreboard_objective":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketScoreboardObjective)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "scoreboard_score":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketScoreboardScore)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "select_advancement_tab":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSelectAdvancementTab)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "server_data":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketServerData)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "server_links":
		PlayToClientPacketParams, ok := ret.Params.(PacketCommonServerLinks)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_cooldown":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSetCooldown)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_cursor_item":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSetCursorItem)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_passengers":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSetPassengers)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_player_inventory":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSetPlayerInventory)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_projectile_power":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSetProjectilePower)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_slot":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSetSlot)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_ticking_state":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSetTickingState)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_title_subtitle":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSetTitleSubtitle)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_title_text":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSetTitleText)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "set_title_time":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSetTitleTime)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "show_dialog":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketShowDialog)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "simulation_distance":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSimulationDistance)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "sound_effect":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSoundEffect)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "spawn_entity":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSpawnEntity)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "spawn_position":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSpawnPosition)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "start_configuration":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketStartConfiguration)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "statistics":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketStatistics)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "step_tick":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketStepTick)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "stop_sound":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketStopSound)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "store_cookie":
		PlayToClientPacketParams, ok := ret.Params.(PacketCommonStoreCookie)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "sync_entity_position":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSyncEntityPosition)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "system_chat":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketSystemChat)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "tab_complete":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketTabComplete)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "tags":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketTags)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "teams":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketTeams)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "test_instance_block_status":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketTestInstanceBlockStatus)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "tile_entity_data":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketTileEntityData)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "tracked_waypoint":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketTrackedWaypoint)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "trade_list":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketTradeList)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "transfer":
		PlayToClientPacketParams, ok := ret.Params.(PacketCommonTransfer)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "unload_chunk":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketUnloadChunk)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "update_health":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketUpdateHealth)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "update_light":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketUpdateLight)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "update_time":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketUpdateTime)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "update_view_distance":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketUpdateViewDistance)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "update_view_position":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketUpdateViewPosition)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "vehicle_move":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketVehicleMove)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "window_items":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketWindowItems)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "world_border_center":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketWorldBorderCenter)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "world_border_lerp_size":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketWorldBorderLerpSize)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "world_border_size":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketWorldBorderSize)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "world_border_warning_delay":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketWorldBorderWarningDelay)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "world_border_warning_reach":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketWorldBorderWarningReach)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "world_event":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketWorldEvent)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	case "world_particles":
		PlayToClientPacketParams, ok := ret.Params.(PlayToClientPacketWorldParticles)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketParams.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketAbilities struct {
	Flags        int8
	FlyingSpeed  float32
	WalkingSpeed float32
}

func (_ PlayToClientPacketAbilities) Decode(r io.Reader) (ret PlayToClientPacketAbilities, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Flags)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.FlyingSpeed)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.WalkingSpeed)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketAbilities) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Flags)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.FlyingSpeed)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.WalkingSpeed)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketAcknowledgePlayerDigging struct {
	SequenceId queser.VarInt
}

func (_ PlayToClientPacketAcknowledgePlayerDigging) Decode(r io.Reader) (ret PlayToClientPacketAcknowledgePlayerDigging, err error) {
	ret.SequenceId, err = ret.SequenceId.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketAcknowledgePlayerDigging) Encode(w io.Writer) (err error) {
	err = ret.SequenceId.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketActionBar struct {
	Text nbt.Anon
}

func (_ PlayToClientPacketActionBar) Decode(r io.Reader) (ret PlayToClientPacketActionBar, err error) {
	ret.Text, err = ret.Text.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketActionBar) Encode(w io.Writer) (err error) {
	err = ret.Text.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketAdvancements struct {
	Val queser.ToDo
}

func (_ PlayToClientPacketAdvancements) Decode(r io.Reader) (ret PlayToClientPacketAdvancements, err error) {
	err = queser.ToDoError
	return
}
func (ret PlayToClientPacketAdvancements) Encode(w io.Writer) (err error) {
	err = queser.ToDoError
	return
}

type PlayToClientPacketAnimation struct {
	EntityId  queser.VarInt
	Animation uint8
}

func (_ PlayToClientPacketAnimation) Decode(r io.Reader) (ret PlayToClientPacketAnimation, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Animation)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketAnimation) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Animation)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketAttachEntity struct {
	EntityId  int32
	VehicleId int32
}

func (_ PlayToClientPacketAttachEntity) Decode(r io.Reader) (ret PlayToClientPacketAttachEntity, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.EntityId)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.VehicleId)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketAttachEntity) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.EntityId)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.VehicleId)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketBlockAction struct {
	Location Position
	Byte1    uint8
	Byte2    uint8
	BlockId  queser.VarInt
}

func (_ PlayToClientPacketBlockAction) Decode(r io.Reader) (ret PlayToClientPacketBlockAction, err error) {
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Byte1)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Byte2)
	if err != nil {
		return
	}
	ret.BlockId, err = ret.BlockId.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketBlockAction) Encode(w io.Writer) (err error) {
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Byte1)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Byte2)
	if err != nil {
		return
	}
	err = ret.BlockId.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketBlockBreakAnimation struct {
	EntityId     queser.VarInt
	Location     Position
	DestroyStage int8
}

func (_ PlayToClientPacketBlockBreakAnimation) Decode(r io.Reader) (ret PlayToClientPacketBlockBreakAnimation, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.DestroyStage)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketBlockBreakAnimation) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.DestroyStage)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketBlockChange struct {
	Location Position
	Type     queser.VarInt
}

func (_ PlayToClientPacketBlockChange) Decode(r io.Reader) (ret PlayToClientPacketBlockChange, err error) {
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	ret.Type, err = ret.Type.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketBlockChange) Encode(w io.Writer) (err error) {
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = ret.Type.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketBossBar struct {
	EntityUUID uuid.UUID
	Action     queser.VarInt
	Title      any
	Health     any
	Color      any
	Dividers   any
	Flags      any
}

func (_ PlayToClientPacketBossBar) Decode(r io.Reader) (ret PlayToClientPacketBossBar, err error) {
	_, err = io.ReadFull(r, ret.EntityUUID[:])
	if err != nil {
		return
	}
	ret.Action, err = ret.Action.Decode(r)
	if err != nil {
		return
	}
	switch ret.Action {
	case 0:
		var PlayToClientPacketBossBarTitleTmp nbt.Anon
		PlayToClientPacketBossBarTitleTmp, err = PlayToClientPacketBossBarTitleTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Title = PlayToClientPacketBossBarTitleTmp
	case 3:
		var PlayToClientPacketBossBarTitleTmp nbt.Anon
		PlayToClientPacketBossBarTitleTmp, err = PlayToClientPacketBossBarTitleTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Title = PlayToClientPacketBossBarTitleTmp
	default:
		var PlayToClientPacketBossBarTitleTmp queser.Void
		PlayToClientPacketBossBarTitleTmp, err = PlayToClientPacketBossBarTitleTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Title = PlayToClientPacketBossBarTitleTmp
	}
	switch ret.Action {
	case 0:
		var PlayToClientPacketBossBarHealthTmp float32
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketBossBarHealthTmp)
		if err != nil {
			return
		}
		ret.Health = PlayToClientPacketBossBarHealthTmp
	case 2:
		var PlayToClientPacketBossBarHealthTmp float32
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketBossBarHealthTmp)
		if err != nil {
			return
		}
		ret.Health = PlayToClientPacketBossBarHealthTmp
	default:
		var PlayToClientPacketBossBarHealthTmp queser.Void
		PlayToClientPacketBossBarHealthTmp, err = PlayToClientPacketBossBarHealthTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Health = PlayToClientPacketBossBarHealthTmp
	}
	switch ret.Action {
	case 0:
		var PlayToClientPacketBossBarColorTmp queser.VarInt
		PlayToClientPacketBossBarColorTmp, err = PlayToClientPacketBossBarColorTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Color = PlayToClientPacketBossBarColorTmp
	case 4:
		var PlayToClientPacketBossBarColorTmp queser.VarInt
		PlayToClientPacketBossBarColorTmp, err = PlayToClientPacketBossBarColorTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Color = PlayToClientPacketBossBarColorTmp
	default:
		var PlayToClientPacketBossBarColorTmp queser.Void
		PlayToClientPacketBossBarColorTmp, err = PlayToClientPacketBossBarColorTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Color = PlayToClientPacketBossBarColorTmp
	}
	switch ret.Action {
	case 0:
		var PlayToClientPacketBossBarDividersTmp queser.VarInt
		PlayToClientPacketBossBarDividersTmp, err = PlayToClientPacketBossBarDividersTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Dividers = PlayToClientPacketBossBarDividersTmp
	case 4:
		var PlayToClientPacketBossBarDividersTmp queser.VarInt
		PlayToClientPacketBossBarDividersTmp, err = PlayToClientPacketBossBarDividersTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Dividers = PlayToClientPacketBossBarDividersTmp
	default:
		var PlayToClientPacketBossBarDividersTmp queser.Void
		PlayToClientPacketBossBarDividersTmp, err = PlayToClientPacketBossBarDividersTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Dividers = PlayToClientPacketBossBarDividersTmp
	}
	switch ret.Action {
	case 0:
		var PlayToClientPacketBossBarFlagsTmp uint8
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketBossBarFlagsTmp)
		if err != nil {
			return
		}
		ret.Flags = PlayToClientPacketBossBarFlagsTmp
	case 5:
		var PlayToClientPacketBossBarFlagsTmp uint8
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketBossBarFlagsTmp)
		if err != nil {
			return
		}
		ret.Flags = PlayToClientPacketBossBarFlagsTmp
	default:
		var PlayToClientPacketBossBarFlagsTmp queser.Void
		PlayToClientPacketBossBarFlagsTmp, err = PlayToClientPacketBossBarFlagsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Flags = PlayToClientPacketBossBarFlagsTmp
	}
	return
}
func (ret PlayToClientPacketBossBar) Encode(w io.Writer) (err error) {
	_, err = w.Write(ret.EntityUUID[:])
	if err != nil {
		return
	}
	err = ret.Action.Encode(w)
	if err != nil {
		return
	}
	switch ret.Action {
	case 0:
		PlayToClientPacketBossBarTitle, ok := ret.Title.(nbt.Anon)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketBossBarTitle.Encode(w)
		if err != nil {
			return
		}
	case 3:
		PlayToClientPacketBossBarTitle, ok := ret.Title.(nbt.Anon)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketBossBarTitle.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Title.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Title.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Action {
	case 0:
		PlayToClientPacketBossBarHealth, ok := ret.Health.(float32)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToClientPacketBossBarHealth)
		if err != nil {
			return
		}
	case 2:
		PlayToClientPacketBossBarHealth, ok := ret.Health.(float32)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToClientPacketBossBarHealth)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Health.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Health.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Action {
	case 0:
		PlayToClientPacketBossBarColor, ok := ret.Color.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketBossBarColor.Encode(w)
		if err != nil {
			return
		}
	case 4:
		PlayToClientPacketBossBarColor, ok := ret.Color.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketBossBarColor.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Color.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Color.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Action {
	case 0:
		PlayToClientPacketBossBarDividers, ok := ret.Dividers.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketBossBarDividers.Encode(w)
		if err != nil {
			return
		}
	case 4:
		PlayToClientPacketBossBarDividers, ok := ret.Dividers.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketBossBarDividers.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Dividers.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Dividers.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Action {
	case 0:
		PlayToClientPacketBossBarFlags, ok := ret.Flags.(uint8)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToClientPacketBossBarFlags)
		if err != nil {
			return
		}
	case 5:
		PlayToClientPacketBossBarFlags, ok := ret.Flags.(uint8)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToClientPacketBossBarFlags)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Flags.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Flags.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketCamera struct {
	CameraId queser.VarInt
}

func (_ PlayToClientPacketCamera) Decode(r io.Reader) (ret PlayToClientPacketCamera, err error) {
	ret.CameraId, err = ret.CameraId.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketCamera) Encode(w io.Writer) (err error) {
	err = ret.CameraId.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketChatSuggestions struct {
	Action  queser.VarInt
	Entries []string
}

func (_ PlayToClientPacketChatSuggestions) Decode(r io.Reader) (ret PlayToClientPacketChatSuggestions, err error) {
	ret.Action, err = ret.Action.Decode(r)
	if err != nil {
		return
	}
	var lPlayToClientPacketChatSuggestionsEntries queser.VarInt
	lPlayToClientPacketChatSuggestionsEntries, err = lPlayToClientPacketChatSuggestionsEntries.Decode(r)
	if err != nil {
		return
	}
	ret.Entries = []string{}
	for range lPlayToClientPacketChatSuggestionsEntries {
		var PlayToClientPacketChatSuggestionsEntriesElement string
		PlayToClientPacketChatSuggestionsEntriesElement, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Entries = append(ret.Entries, PlayToClientPacketChatSuggestionsEntriesElement)
	}
	return
}
func (ret PlayToClientPacketChatSuggestions) Encode(w io.Writer) (err error) {
	err = ret.Action.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Entries)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketChatSuggestionsEntries := range len(ret.Entries) {
		err = queser.EncodeString(w, ret.Entries[iPlayToClientPacketChatSuggestionsEntries])
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketChunkBatchFinished struct {
	BatchSize queser.VarInt
}

func (_ PlayToClientPacketChunkBatchFinished) Decode(r io.Reader) (ret PlayToClientPacketChunkBatchFinished, err error) {
	ret.BatchSize, err = ret.BatchSize.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketChunkBatchFinished) Encode(w io.Writer) (err error) {
	err = ret.BatchSize.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketChunkBatchStart struct {
}

func (_ PlayToClientPacketChunkBatchStart) Decode(r io.Reader) (ret PlayToClientPacketChunkBatchStart, err error) {
	return
}
func (ret PlayToClientPacketChunkBatchStart) Encode(w io.Writer) (err error) {
	return
}

type PlayToClientPacketChunkBiomes struct {
	Biomes []struct {
		Position PackedChunkPos
		Data     ByteArray
	}
}

func (_ PlayToClientPacketChunkBiomes) Decode(r io.Reader) (ret PlayToClientPacketChunkBiomes, err error) {
	var lPlayToClientPacketChunkBiomesBiomes queser.VarInt
	lPlayToClientPacketChunkBiomesBiomes, err = lPlayToClientPacketChunkBiomesBiomes.Decode(r)
	if err != nil {
		return
	}
	ret.Biomes = []struct {
		Position PackedChunkPos
		Data     ByteArray
	}{}
	for range lPlayToClientPacketChunkBiomesBiomes {
		var PlayToClientPacketChunkBiomesBiomesElement struct {
			Position PackedChunkPos
			Data     ByteArray
		}
		PlayToClientPacketChunkBiomesBiomesElement.Position, err = PlayToClientPacketChunkBiomesBiomesElement.Position.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketChunkBiomesBiomesElement.Data, err = PlayToClientPacketChunkBiomesBiomesElement.Data.Decode(r)
		if err != nil {
			return
		}
		ret.Biomes = append(ret.Biomes, PlayToClientPacketChunkBiomesBiomesElement)
	}
	return
}
func (ret PlayToClientPacketChunkBiomes) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Biomes)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketChunkBiomesBiomes := range len(ret.Biomes) {
		err = ret.Biomes[iPlayToClientPacketChunkBiomesBiomes].Position.Encode(w)
		if err != nil {
			return
		}
		err = ret.Biomes[iPlayToClientPacketChunkBiomesBiomes].Data.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketClearTitles struct {
	Reset bool
}

func (_ PlayToClientPacketClearTitles) Decode(r io.Reader) (ret PlayToClientPacketClearTitles, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Reset)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketClearTitles) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Reset)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketCloseWindow struct {
	WindowId ContainerID
}

func (_ PlayToClientPacketCloseWindow) Decode(r io.Reader) (ret PlayToClientPacketCloseWindow, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketCloseWindow) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketCollect struct {
	CollectedEntityId queser.VarInt
	CollectorEntityId queser.VarInt
	PickupItemCount   queser.VarInt
}

func (_ PlayToClientPacketCollect) Decode(r io.Reader) (ret PlayToClientPacketCollect, err error) {
	ret.CollectedEntityId, err = ret.CollectedEntityId.Decode(r)
	if err != nil {
		return
	}
	ret.CollectorEntityId, err = ret.CollectorEntityId.Decode(r)
	if err != nil {
		return
	}
	ret.PickupItemCount, err = ret.PickupItemCount.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketCollect) Encode(w io.Writer) (err error) {
	err = ret.CollectedEntityId.Encode(w)
	if err != nil {
		return
	}
	err = ret.CollectorEntityId.Encode(w)
	if err != nil {
		return
	}
	err = ret.PickupItemCount.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketCraftProgressBar struct {
	WindowId ContainerID
	Property int16
	Value    int16
}

func (_ PlayToClientPacketCraftProgressBar) Decode(r io.Reader) (ret PlayToClientPacketCraftProgressBar, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Property)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Value)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketCraftProgressBar) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Property)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Value)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketCraftRecipeResponse struct {
	WindowId      ContainerID
	RecipeDisplay PlayToClientRecipeDisplay
}

func (_ PlayToClientPacketCraftRecipeResponse) Decode(r io.Reader) (ret PlayToClientPacketCraftRecipeResponse, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	ret.RecipeDisplay, err = ret.RecipeDisplay.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketCraftRecipeResponse) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	err = ret.RecipeDisplay.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketCustomPayload struct {
	Channel string
	Data    queser.RestBuffer
}

func (_ PlayToClientPacketCustomPayload) Decode(r io.Reader) (ret PlayToClientPacketCustomPayload, err error) {
	ret.Channel, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.Data, err = ret.Data.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketCustomPayload) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Channel)
	if err != nil {
		return
	}
	err = ret.Data.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketDamageEvent struct {
	EntityId       queser.VarInt
	SourceTypeId   queser.VarInt
	SourceCauseId  queser.VarInt
	SourceDirectId queser.VarInt
	SourcePosition *Vec3f64
}

func (_ PlayToClientPacketDamageEvent) Decode(r io.Reader) (ret PlayToClientPacketDamageEvent, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	ret.SourceTypeId, err = ret.SourceTypeId.Decode(r)
	if err != nil {
		return
	}
	ret.SourceCauseId, err = ret.SourceCauseId.Decode(r)
	if err != nil {
		return
	}
	ret.SourceDirectId, err = ret.SourceDirectId.Decode(r)
	if err != nil {
		return
	}
	var PlayToClientPacketDamageEventSourcePositionPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketDamageEventSourcePositionPresent)
	if err != nil {
		return
	}
	if PlayToClientPacketDamageEventSourcePositionPresent {
		var PlayToClientPacketDamageEventSourcePositionPresentValue Vec3f64
		PlayToClientPacketDamageEventSourcePositionPresentValue, err = PlayToClientPacketDamageEventSourcePositionPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.SourcePosition = &PlayToClientPacketDamageEventSourcePositionPresentValue
	}
	return
}
func (ret PlayToClientPacketDamageEvent) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = ret.SourceTypeId.Encode(w)
	if err != nil {
		return
	}
	err = ret.SourceCauseId.Encode(w)
	if err != nil {
		return
	}
	err = ret.SourceDirectId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.SourcePosition != nil)
	if err != nil {
		return
	}
	if ret.SourcePosition != nil {
		err = (*ret.SourcePosition).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketDeathCombatEvent struct {
	PlayerId queser.VarInt
	Message  nbt.Anon
}

func (_ PlayToClientPacketDeathCombatEvent) Decode(r io.Reader) (ret PlayToClientPacketDeathCombatEvent, err error) {
	ret.PlayerId, err = ret.PlayerId.Decode(r)
	if err != nil {
		return
	}
	ret.Message, err = ret.Message.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketDeathCombatEvent) Encode(w io.Writer) (err error) {
	err = ret.PlayerId.Encode(w)
	if err != nil {
		return
	}
	err = ret.Message.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketDebugSample struct {
	Sample []int64
	Type   queser.VarInt
}

func (_ PlayToClientPacketDebugSample) Decode(r io.Reader) (ret PlayToClientPacketDebugSample, err error) {
	var lPlayToClientPacketDebugSampleSample queser.VarInt
	lPlayToClientPacketDebugSampleSample, err = lPlayToClientPacketDebugSampleSample.Decode(r)
	if err != nil {
		return
	}
	ret.Sample = []int64{}
	for range lPlayToClientPacketDebugSampleSample {
		var PlayToClientPacketDebugSampleSampleElement int64
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketDebugSampleSampleElement)
		if err != nil {
			return
		}
		ret.Sample = append(ret.Sample, PlayToClientPacketDebugSampleSampleElement)
	}
	ret.Type, err = ret.Type.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketDebugSample) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Sample)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketDebugSampleSample := range len(ret.Sample) {
		err = binary.Write(w, binary.BigEndian, ret.Sample[iPlayToClientPacketDebugSampleSample])
		if err != nil {
			return
		}
	}
	err = ret.Type.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketDeclareCommands struct {
	Nodes     []CommandNode
	RootIndex queser.VarInt
}

func (_ PlayToClientPacketDeclareCommands) Decode(r io.Reader) (ret PlayToClientPacketDeclareCommands, err error) {
	var lPlayToClientPacketDeclareCommandsNodes queser.VarInt
	lPlayToClientPacketDeclareCommandsNodes, err = lPlayToClientPacketDeclareCommandsNodes.Decode(r)
	if err != nil {
		return
	}
	ret.Nodes = []CommandNode{}
	for range lPlayToClientPacketDeclareCommandsNodes {
		var PlayToClientPacketDeclareCommandsNodesElement CommandNode
		PlayToClientPacketDeclareCommandsNodesElement, err = PlayToClientPacketDeclareCommandsNodesElement.Decode(r)
		if err != nil {
			return
		}
		ret.Nodes = append(ret.Nodes, PlayToClientPacketDeclareCommandsNodesElement)
	}
	ret.RootIndex, err = ret.RootIndex.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketDeclareCommands) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Nodes)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketDeclareCommandsNodes := range len(ret.Nodes) {
		err = ret.Nodes[iPlayToClientPacketDeclareCommandsNodes].Encode(w)
		if err != nil {
			return
		}
	}
	err = ret.RootIndex.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketDeclareRecipes struct {
	Recipes []struct {
		Name  string
		Items []queser.VarInt
	}
	StoneCutterRecipes []struct {
		Input       IDSet
		SlotDisplay PlayToClientSlotDisplay
	}
}

func (_ PlayToClientPacketDeclareRecipes) Decode(r io.Reader) (ret PlayToClientPacketDeclareRecipes, err error) {
	var lPlayToClientPacketDeclareRecipesRecipes queser.VarInt
	lPlayToClientPacketDeclareRecipesRecipes, err = lPlayToClientPacketDeclareRecipesRecipes.Decode(r)
	if err != nil {
		return
	}
	ret.Recipes = []struct {
		Name  string
		Items []queser.VarInt
	}{}
	for range lPlayToClientPacketDeclareRecipesRecipes {
		var PlayToClientPacketDeclareRecipesRecipesElement struct {
			Name  string
			Items []queser.VarInt
		}
		PlayToClientPacketDeclareRecipesRecipesElement.Name, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		var lPlayToClientPacketDeclareRecipesRecipesElementItems queser.VarInt
		lPlayToClientPacketDeclareRecipesRecipesElementItems, err = lPlayToClientPacketDeclareRecipesRecipesElementItems.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketDeclareRecipesRecipesElement.Items = []queser.VarInt{}
		for range lPlayToClientPacketDeclareRecipesRecipesElementItems {
			var PlayToClientPacketDeclareRecipesRecipesElementItemsElement queser.VarInt
			PlayToClientPacketDeclareRecipesRecipesElementItemsElement, err = PlayToClientPacketDeclareRecipesRecipesElementItemsElement.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketDeclareRecipesRecipesElement.Items = append(PlayToClientPacketDeclareRecipesRecipesElement.Items, PlayToClientPacketDeclareRecipesRecipesElementItemsElement)
		}
		ret.Recipes = append(ret.Recipes, PlayToClientPacketDeclareRecipesRecipesElement)
	}
	var lPlayToClientPacketDeclareRecipesStoneCutterRecipes queser.VarInt
	lPlayToClientPacketDeclareRecipesStoneCutterRecipes, err = lPlayToClientPacketDeclareRecipesStoneCutterRecipes.Decode(r)
	if err != nil {
		return
	}
	ret.StoneCutterRecipes = []struct {
		Input       IDSet
		SlotDisplay PlayToClientSlotDisplay
	}{}
	for range lPlayToClientPacketDeclareRecipesStoneCutterRecipes {
		var PlayToClientPacketDeclareRecipesStoneCutterRecipesElement struct {
			Input       IDSet
			SlotDisplay PlayToClientSlotDisplay
		}
		PlayToClientPacketDeclareRecipesStoneCutterRecipesElement.Input, err = PlayToClientPacketDeclareRecipesStoneCutterRecipesElement.Input.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketDeclareRecipesStoneCutterRecipesElement.SlotDisplay, err = PlayToClientPacketDeclareRecipesStoneCutterRecipesElement.SlotDisplay.Decode(r)
		if err != nil {
			return
		}
		ret.StoneCutterRecipes = append(ret.StoneCutterRecipes, PlayToClientPacketDeclareRecipesStoneCutterRecipesElement)
	}
	return
}
func (ret PlayToClientPacketDeclareRecipes) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Recipes)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketDeclareRecipesRecipes := range len(ret.Recipes) {
		err = queser.EncodeString(w, ret.Recipes[iPlayToClientPacketDeclareRecipesRecipes].Name)
		if err != nil {
			return
		}
		err = queser.VarInt(len(ret.Recipes[iPlayToClientPacketDeclareRecipesRecipes].Items)).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketDeclareRecipesRecipesInnerItems := range len(ret.Recipes[iPlayToClientPacketDeclareRecipesRecipes].Items) {
			err = ret.Recipes[iPlayToClientPacketDeclareRecipesRecipes].Items[iPlayToClientPacketDeclareRecipesRecipesInnerItems].Encode(w)
			if err != nil {
				return
			}
		}
	}
	err = queser.VarInt(len(ret.StoneCutterRecipes)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketDeclareRecipesStoneCutterRecipes := range len(ret.StoneCutterRecipes) {
		err = ret.StoneCutterRecipes[iPlayToClientPacketDeclareRecipesStoneCutterRecipes].Input.Encode(w)
		if err != nil {
			return
		}
		err = ret.StoneCutterRecipes[iPlayToClientPacketDeclareRecipesStoneCutterRecipes].SlotDisplay.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketDifficulty struct {
	Difficulty       string
	DifficultyLocked bool
}

var PlayToClientPacketDifficultyDifficultyMap = map[queser.VarInt]string{0: "peaceful", 1: "easy", 2: "normal", 3: "hard"}

func (_ PlayToClientPacketDifficulty) Decode(r io.Reader) (ret PlayToClientPacketDifficulty, err error) {
	var PlayToClientPacketDifficultyDifficultyKey queser.VarInt
	PlayToClientPacketDifficultyDifficultyKey, err = PlayToClientPacketDifficultyDifficultyKey.Decode(r)
	if err != nil {
		return
	}
	ret.Difficulty, err = queser.ErroringIndex(PlayToClientPacketDifficultyDifficultyMap, PlayToClientPacketDifficultyDifficultyKey)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.DifficultyLocked)
	if err != nil {
		return
	}
	return
}

var PlayToClientPacketDifficultyDifficultyReverseMap = map[string]queser.VarInt{"peaceful": 0, "easy": 1, "normal": 2, "hard": 3}

func (ret PlayToClientPacketDifficulty) Encode(w io.Writer) (err error) {
	var vPlayToClientPacketDifficultyDifficulty queser.VarInt
	vPlayToClientPacketDifficultyDifficulty, err = queser.ErroringIndex(PlayToClientPacketDifficultyDifficultyReverseMap, ret.Difficulty)
	if err != nil {
		return
	}
	err = vPlayToClientPacketDifficultyDifficulty.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.DifficultyLocked)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketEndCombatEvent struct {
	Duration queser.VarInt
}

func (_ PlayToClientPacketEndCombatEvent) Decode(r io.Reader) (ret PlayToClientPacketEndCombatEvent, err error) {
	ret.Duration, err = ret.Duration.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketEndCombatEvent) Encode(w io.Writer) (err error) {
	err = ret.Duration.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketEnterCombatEvent struct {
}

func (_ PlayToClientPacketEnterCombatEvent) Decode(r io.Reader) (ret PlayToClientPacketEnterCombatEvent, err error) {
	return
}
func (ret PlayToClientPacketEnterCombatEvent) Encode(w io.Writer) (err error) {
	return
}

type PlayToClientPacketEntityDestroy struct {
	EntityIds []queser.VarInt
}

func (_ PlayToClientPacketEntityDestroy) Decode(r io.Reader) (ret PlayToClientPacketEntityDestroy, err error) {
	var lPlayToClientPacketEntityDestroyEntityIds queser.VarInt
	lPlayToClientPacketEntityDestroyEntityIds, err = lPlayToClientPacketEntityDestroyEntityIds.Decode(r)
	if err != nil {
		return
	}
	ret.EntityIds = []queser.VarInt{}
	for range lPlayToClientPacketEntityDestroyEntityIds {
		var PlayToClientPacketEntityDestroyEntityIdsElement queser.VarInt
		PlayToClientPacketEntityDestroyEntityIdsElement, err = PlayToClientPacketEntityDestroyEntityIdsElement.Decode(r)
		if err != nil {
			return
		}
		ret.EntityIds = append(ret.EntityIds, PlayToClientPacketEntityDestroyEntityIdsElement)
	}
	return
}
func (ret PlayToClientPacketEntityDestroy) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.EntityIds)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketEntityDestroyEntityIds := range len(ret.EntityIds) {
		err = ret.EntityIds[iPlayToClientPacketEntityDestroyEntityIds].Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketEntityEffect struct {
	EntityId  queser.VarInt
	EffectId  queser.VarInt
	Amplifier queser.VarInt
	Duration  queser.VarInt
	Flags     uint8
}

func (_ PlayToClientPacketEntityEffect) Decode(r io.Reader) (ret PlayToClientPacketEntityEffect, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	ret.EffectId, err = ret.EffectId.Decode(r)
	if err != nil {
		return
	}
	ret.Amplifier, err = ret.Amplifier.Decode(r)
	if err != nil {
		return
	}
	ret.Duration, err = ret.Duration.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Flags)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketEntityEffect) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = ret.EffectId.Encode(w)
	if err != nil {
		return
	}
	err = ret.Amplifier.Encode(w)
	if err != nil {
		return
	}
	err = ret.Duration.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Flags)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketEntityEquipment struct {
	Val queser.ToDo
}

func (_ PlayToClientPacketEntityEquipment) Decode(r io.Reader) (ret PlayToClientPacketEntityEquipment, err error) {
	err = queser.ToDoError
	return
}
func (ret PlayToClientPacketEntityEquipment) Encode(w io.Writer) (err error) {
	err = queser.ToDoError
	return
}

type PlayToClientPacketEntityHeadRotation struct {
	EntityId queser.VarInt
	HeadYaw  int8
}

func (_ PlayToClientPacketEntityHeadRotation) Decode(r io.Reader) (ret PlayToClientPacketEntityHeadRotation, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.HeadYaw)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketEntityHeadRotation) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.HeadYaw)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketEntityLook struct {
	EntityId queser.VarInt
	Yaw      int8
	Pitch    int8
	OnGround bool
}

func (_ PlayToClientPacketEntityLook) Decode(r io.Reader) (ret PlayToClientPacketEntityLook, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OnGround)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketEntityLook) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OnGround)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketEntityMetadata struct {
	EntityId queser.VarInt
	Metadata EntityMetadata
}

func (_ PlayToClientPacketEntityMetadata) Decode(r io.Reader) (ret PlayToClientPacketEntityMetadata, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	ret.Metadata, err = ret.Metadata.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketEntityMetadata) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = ret.Metadata.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketEntityMoveLook struct {
	EntityId queser.VarInt
	DX       int16
	DY       int16
	DZ       int16
	Yaw      int8
	Pitch    int8
	OnGround bool
}

func (_ PlayToClientPacketEntityMoveLook) Decode(r io.Reader) (ret PlayToClientPacketEntityMoveLook, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.DX)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.DY)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.DZ)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OnGround)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketEntityMoveLook) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.DX)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.DY)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.DZ)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OnGround)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketEntitySoundEffect struct {
	Sound         ItemSoundHolder
	SoundCategory SoundSource
	EntityId      queser.VarInt
	Volume        float32
	Pitch         float32
	Seed          int64
}

func (_ PlayToClientPacketEntitySoundEffect) Decode(r io.Reader) (ret PlayToClientPacketEntitySoundEffect, err error) {
	ret.Sound, err = ret.Sound.Decode(r)
	if err != nil {
		return
	}
	ret.SoundCategory, err = ret.SoundCategory.Decode(r)
	if err != nil {
		return
	}
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Volume)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Seed)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketEntitySoundEffect) Encode(w io.Writer) (err error) {
	err = ret.Sound.Encode(w)
	if err != nil {
		return
	}
	err = ret.SoundCategory.Encode(w)
	if err != nil {
		return
	}
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Volume)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Seed)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketEntityStatus struct {
	EntityId     int32
	EntityStatus int8
}

func (_ PlayToClientPacketEntityStatus) Decode(r io.Reader) (ret PlayToClientPacketEntityStatus, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.EntityId)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.EntityStatus)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketEntityStatus) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.EntityId)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.EntityStatus)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketEntityTeleport struct {
	EntityId queser.VarInt
	X        float64
	Y        float64
	Z        float64
	Yaw      int8
	Pitch    int8
	OnGround bool
}

func (_ PlayToClientPacketEntityTeleport) Decode(r io.Reader) (ret PlayToClientPacketEntityTeleport, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OnGround)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketEntityTeleport) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OnGround)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketEntityUpdateAttributes struct {
	EntityId   queser.VarInt
	Properties []struct {
		Key       string
		Value     float64
		Modifiers []struct {
			Uuid      string
			Amount    float64
			Operation int8
		}
	}
}

var PlayToClientPacketEntityUpdateAttributesPropertiesElementKeyMap = map[queser.VarInt]string{0: "generic.armor", 1: "generic.armor_toughness", 10: "player.entity_interaction_range", 11: "generic.fall_damage_multiplier", 12: "generic.flying_speed", 13: "generic.follow_range", 14: "generic.gravity", 15: "generic.jump_strength", 16: "generic.knockback_resistance", 17: "generic.luck", 18: "generic.max_absorption", 19: "generic.max_health", 2: "generic.attack_damage", 20: "generic.movement_speed", 21: "generic.safe_fall_distance", 22: "generic.scale", 23: "zombie.spawn_reinforcements", 24: "generic.step_height", 25: "submerged_mining_speed", 26: "sweeping_damage_ratio", 27: "tempt_range", 28: "water_movement_efficiency", 29: "waypoint_transmit_range", 3: "generic.attack_knockback", 30: "waypoint_receive_range", 4: "generic.attack_speed", 5: "player.block_break_speed", 6: "player.block_interaction_range", 7: "burning_time", 8: "camera_distance", 9: "explosion_knockback_resistance"}

func (_ PlayToClientPacketEntityUpdateAttributes) Decode(r io.Reader) (ret PlayToClientPacketEntityUpdateAttributes, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	var lPlayToClientPacketEntityUpdateAttributesProperties queser.VarInt
	lPlayToClientPacketEntityUpdateAttributesProperties, err = lPlayToClientPacketEntityUpdateAttributesProperties.Decode(r)
	if err != nil {
		return
	}
	ret.Properties = []struct {
		Key       string
		Value     float64
		Modifiers []struct {
			Uuid      string
			Amount    float64
			Operation int8
		}
	}{}
	for range lPlayToClientPacketEntityUpdateAttributesProperties {
		var PlayToClientPacketEntityUpdateAttributesPropertiesElement struct {
			Key       string
			Value     float64
			Modifiers []struct {
				Uuid      string
				Amount    float64
				Operation int8
			}
		}
		var PlayToClientPacketEntityUpdateAttributesPropertiesElementKeyKey queser.VarInt
		PlayToClientPacketEntityUpdateAttributesPropertiesElementKeyKey, err = PlayToClientPacketEntityUpdateAttributesPropertiesElementKeyKey.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketEntityUpdateAttributesPropertiesElement.Key, err = queser.ErroringIndex(PlayToClientPacketEntityUpdateAttributesPropertiesElementKeyMap, PlayToClientPacketEntityUpdateAttributesPropertiesElementKeyKey)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketEntityUpdateAttributesPropertiesElement.Value)
		if err != nil {
			return
		}
		var lPlayToClientPacketEntityUpdateAttributesPropertiesElementModifiers queser.VarInt
		lPlayToClientPacketEntityUpdateAttributesPropertiesElementModifiers, err = lPlayToClientPacketEntityUpdateAttributesPropertiesElementModifiers.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketEntityUpdateAttributesPropertiesElement.Modifiers = []struct {
			Uuid      string
			Amount    float64
			Operation int8
		}{}
		for range lPlayToClientPacketEntityUpdateAttributesPropertiesElementModifiers {
			var PlayToClientPacketEntityUpdateAttributesPropertiesElementModifiersElement struct {
				Uuid      string
				Amount    float64
				Operation int8
			}
			PlayToClientPacketEntityUpdateAttributesPropertiesElementModifiersElement.Uuid, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketEntityUpdateAttributesPropertiesElementModifiersElement.Amount)
			if err != nil {
				return
			}
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketEntityUpdateAttributesPropertiesElementModifiersElement.Operation)
			if err != nil {
				return
			}
			PlayToClientPacketEntityUpdateAttributesPropertiesElement.Modifiers = append(PlayToClientPacketEntityUpdateAttributesPropertiesElement.Modifiers, PlayToClientPacketEntityUpdateAttributesPropertiesElementModifiersElement)
		}
		ret.Properties = append(ret.Properties, PlayToClientPacketEntityUpdateAttributesPropertiesElement)
	}
	return
}

var PlayToClientPacketEntityUpdateAttributesPropertiesKeyReverseMap = map[string]queser.VarInt{"generic.armor": 0, "generic.armor_toughness": 1, "player.entity_interaction_range": 10, "generic.fall_damage_multiplier": 11, "generic.flying_speed": 12, "generic.follow_range": 13, "generic.gravity": 14, "generic.jump_strength": 15, "generic.knockback_resistance": 16, "generic.luck": 17, "generic.max_absorption": 18, "generic.max_health": 19, "generic.attack_damage": 2, "generic.movement_speed": 20, "generic.safe_fall_distance": 21, "generic.scale": 22, "zombie.spawn_reinforcements": 23, "generic.step_height": 24, "submerged_mining_speed": 25, "sweeping_damage_ratio": 26, "tempt_range": 27, "water_movement_efficiency": 28, "waypoint_transmit_range": 29, "generic.attack_knockback": 3, "waypoint_receive_range": 30, "generic.attack_speed": 4, "player.block_break_speed": 5, "player.block_interaction_range": 6, "burning_time": 7, "camera_distance": 8, "explosion_knockback_resistance": 9}
var PlayToClientPacketEntityUpdateAttributesPropertiesInnerKeyReverseMap = map[string]queser.VarInt{"generic.armor": 0, "generic.armor_toughness": 1, "player.entity_interaction_range": 10, "generic.fall_damage_multiplier": 11, "generic.flying_speed": 12, "generic.follow_range": 13, "generic.gravity": 14, "generic.jump_strength": 15, "generic.knockback_resistance": 16, "generic.luck": 17, "generic.max_absorption": 18, "generic.max_health": 19, "generic.attack_damage": 2, "generic.movement_speed": 20, "generic.safe_fall_distance": 21, "generic.scale": 22, "zombie.spawn_reinforcements": 23, "generic.step_height": 24, "submerged_mining_speed": 25, "sweeping_damage_ratio": 26, "tempt_range": 27, "water_movement_efficiency": 28, "waypoint_transmit_range": 29, "generic.attack_knockback": 3, "waypoint_receive_range": 30, "generic.attack_speed": 4, "player.block_break_speed": 5, "player.block_interaction_range": 6, "burning_time": 7, "camera_distance": 8, "explosion_knockback_resistance": 9}

func (ret PlayToClientPacketEntityUpdateAttributes) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Properties)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketEntityUpdateAttributesProperties := range len(ret.Properties) {
		var vPlayToClientPacketEntityUpdateAttributesPropertiesInnerKey queser.VarInt
		vPlayToClientPacketEntityUpdateAttributesPropertiesInnerKey, err = queser.ErroringIndex(PlayToClientPacketEntityUpdateAttributesPropertiesInnerKeyReverseMap, ret.Properties[iPlayToClientPacketEntityUpdateAttributesProperties].Key)
		if err != nil {
			return
		}
		err = vPlayToClientPacketEntityUpdateAttributesPropertiesInnerKey.Encode(w)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Properties[iPlayToClientPacketEntityUpdateAttributesProperties].Value)
		if err != nil {
			return
		}
		err = queser.VarInt(len(ret.Properties[iPlayToClientPacketEntityUpdateAttributesProperties].Modifiers)).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketEntityUpdateAttributesPropertiesInnerModifiers := range len(ret.Properties[iPlayToClientPacketEntityUpdateAttributesProperties].Modifiers) {
			err = queser.EncodeString(w, ret.Properties[iPlayToClientPacketEntityUpdateAttributesProperties].Modifiers[iPlayToClientPacketEntityUpdateAttributesPropertiesInnerModifiers].Uuid)
			if err != nil {
				return
			}
			err = binary.Write(w, binary.BigEndian, ret.Properties[iPlayToClientPacketEntityUpdateAttributesProperties].Modifiers[iPlayToClientPacketEntityUpdateAttributesPropertiesInnerModifiers].Amount)
			if err != nil {
				return
			}
			err = binary.Write(w, binary.BigEndian, ret.Properties[iPlayToClientPacketEntityUpdateAttributesProperties].Modifiers[iPlayToClientPacketEntityUpdateAttributesPropertiesInnerModifiers].Operation)
			if err != nil {
				return
			}
		}
	}
	return
}

type PlayToClientPacketEntityVelocity struct {
	EntityId  queser.VarInt
	VelocityX int16
	VelocityY int16
	VelocityZ int16
}

func (_ PlayToClientPacketEntityVelocity) Decode(r io.Reader) (ret PlayToClientPacketEntityVelocity, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.VelocityX)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.VelocityY)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.VelocityZ)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketEntityVelocity) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.VelocityX)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.VelocityY)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.VelocityZ)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketExperience struct {
	ExperienceBar   float32
	Level           queser.VarInt
	TotalExperience queser.VarInt
}

func (_ PlayToClientPacketExperience) Decode(r io.Reader) (ret PlayToClientPacketExperience, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.ExperienceBar)
	if err != nil {
		return
	}
	ret.Level, err = ret.Level.Decode(r)
	if err != nil {
		return
	}
	ret.TotalExperience, err = ret.TotalExperience.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketExperience) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.ExperienceBar)
	if err != nil {
		return
	}
	err = ret.Level.Encode(w)
	if err != nil {
		return
	}
	err = ret.TotalExperience.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketExplosion struct {
	X                 float64
	Y                 float64
	Z                 float64
	PlayerKnockback   *Vec3f
	ExplosionParticle Particle
	Sound             ItemSoundHolder
}

func (_ PlayToClientPacketExplosion) Decode(r io.Reader) (ret PlayToClientPacketExplosion, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	var PlayToClientPacketExplosionPlayerKnockbackPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketExplosionPlayerKnockbackPresent)
	if err != nil {
		return
	}
	if PlayToClientPacketExplosionPlayerKnockbackPresent {
		var PlayToClientPacketExplosionPlayerKnockbackPresentValue Vec3f
		PlayToClientPacketExplosionPlayerKnockbackPresentValue, err = PlayToClientPacketExplosionPlayerKnockbackPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.PlayerKnockback = &PlayToClientPacketExplosionPlayerKnockbackPresentValue
	}
	ret.ExplosionParticle, err = ret.ExplosionParticle.Decode(r)
	if err != nil {
		return
	}
	ret.Sound, err = ret.Sound.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketExplosion) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.PlayerKnockback != nil)
	if err != nil {
		return
	}
	if ret.PlayerKnockback != nil {
		err = (*ret.PlayerKnockback).Encode(w)
		if err != nil {
			return
		}
	}
	err = ret.ExplosionParticle.Encode(w)
	if err != nil {
		return
	}
	err = ret.Sound.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketFacePlayer struct {
	FeetEyes       queser.VarInt
	X              float64
	Y              float64
	Z              float64
	IsEntity       bool
	EntityId       any
	EntityFeetEyes any
}

func (_ PlayToClientPacketFacePlayer) Decode(r io.Reader) (ret PlayToClientPacketFacePlayer, err error) {
	ret.FeetEyes, err = ret.FeetEyes.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IsEntity)
	if err != nil {
		return
	}
	switch ret.IsEntity {
	case true:
		var PlayToClientPacketFacePlayerEntityIdTmp queser.VarInt
		PlayToClientPacketFacePlayerEntityIdTmp, err = PlayToClientPacketFacePlayerEntityIdTmp.Decode(r)
		if err != nil {
			return
		}
		ret.EntityId = PlayToClientPacketFacePlayerEntityIdTmp
	default:
		var PlayToClientPacketFacePlayerEntityIdTmp queser.Void
		PlayToClientPacketFacePlayerEntityIdTmp, err = PlayToClientPacketFacePlayerEntityIdTmp.Decode(r)
		if err != nil {
			return
		}
		ret.EntityId = PlayToClientPacketFacePlayerEntityIdTmp
	}
	switch ret.IsEntity {
	case true:
		var PlayToClientPacketFacePlayerEntityFeetEyesTmp queser.VarInt
		PlayToClientPacketFacePlayerEntityFeetEyesTmp, err = PlayToClientPacketFacePlayerEntityFeetEyesTmp.Decode(r)
		if err != nil {
			return
		}
		ret.EntityFeetEyes = PlayToClientPacketFacePlayerEntityFeetEyesTmp
	default:
		var PlayToClientPacketFacePlayerEntityFeetEyesTmp queser.Void
		PlayToClientPacketFacePlayerEntityFeetEyesTmp, err = PlayToClientPacketFacePlayerEntityFeetEyesTmp.Decode(r)
		if err != nil {
			return
		}
		ret.EntityFeetEyes = PlayToClientPacketFacePlayerEntityFeetEyesTmp
	}
	return
}
func (ret PlayToClientPacketFacePlayer) Encode(w io.Writer) (err error) {
	err = ret.FeetEyes.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IsEntity)
	if err != nil {
		return
	}
	switch ret.IsEntity {
	case true:
		PlayToClientPacketFacePlayerEntityId, ok := ret.EntityId.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketFacePlayerEntityId.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.EntityId.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.EntityId.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.IsEntity {
	case true:
		PlayToClientPacketFacePlayerEntityFeetEyes, ok := ret.EntityFeetEyes.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketFacePlayerEntityFeetEyes.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.EntityFeetEyes.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.EntityFeetEyes.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketGameStateChange struct {
	Reason   string
	GameMode float32
}

var PlayToClientPacketGameStateChangeReasonMap = map[uint8]string{0: "no_respawn_block_available", 1: "start_raining", 10: "guardian_elder_effect", 11: "immediate_respawn", 12: "limited_crafting", 13: "level_chunks_load_start", 2: "stop_raining", 3: "change_game_mode", 4: "win_game", 5: "demo_event", 6: "play_arrow_hit_sound", 7: "rain_level_change", 8: "thunder_level_change", 9: "puffer_fish_sting"}

func (_ PlayToClientPacketGameStateChange) Decode(r io.Reader) (ret PlayToClientPacketGameStateChange, err error) {
	var PlayToClientPacketGameStateChangeReasonKey uint8
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketGameStateChangeReasonKey)
	if err != nil {
		return
	}
	ret.Reason, err = queser.ErroringIndex(PlayToClientPacketGameStateChangeReasonMap, PlayToClientPacketGameStateChangeReasonKey)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.GameMode)
	if err != nil {
		return
	}
	return
}

var PlayToClientPacketGameStateChangeReasonReverseMap = map[string]uint8{"no_respawn_block_available": 0, "start_raining": 1, "guardian_elder_effect": 10, "immediate_respawn": 11, "limited_crafting": 12, "level_chunks_load_start": 13, "stop_raining": 2, "change_game_mode": 3, "win_game": 4, "demo_event": 5, "play_arrow_hit_sound": 6, "rain_level_change": 7, "thunder_level_change": 8, "puffer_fish_sting": 9}

func (ret PlayToClientPacketGameStateChange) Encode(w io.Writer) (err error) {
	var vPlayToClientPacketGameStateChangeReason uint8
	vPlayToClientPacketGameStateChangeReason, err = queser.ErroringIndex(PlayToClientPacketGameStateChangeReasonReverseMap, ret.Reason)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, vPlayToClientPacketGameStateChangeReason)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.GameMode)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketHeldItemSlot struct {
	Slot queser.VarInt
}

func (_ PlayToClientPacketHeldItemSlot) Decode(r io.Reader) (ret PlayToClientPacketHeldItemSlot, err error) {
	ret.Slot, err = ret.Slot.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketHeldItemSlot) Encode(w io.Writer) (err error) {
	err = ret.Slot.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketHideMessage struct {
	Id        queser.VarInt
	Signature any
}

func (_ PlayToClientPacketHideMessage) Decode(r io.Reader) (ret PlayToClientPacketHideMessage, err error) {
	ret.Id, err = ret.Id.Decode(r)
	if err != nil {
		return
	}
	switch ret.Id {
	case 0:
		var PlayToClientPacketHideMessageSignatureTmp [256]byte
		_, err = r.Read(PlayToClientPacketHideMessageSignatureTmp[:])
		if err != nil {
			return
		}
		ret.Signature = PlayToClientPacketHideMessageSignatureTmp
	default:
		var PlayToClientPacketHideMessageSignatureTmp queser.Void
		PlayToClientPacketHideMessageSignatureTmp, err = PlayToClientPacketHideMessageSignatureTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Signature = PlayToClientPacketHideMessageSignatureTmp
	}
	return
}
func (ret PlayToClientPacketHideMessage) Encode(w io.Writer) (err error) {
	err = ret.Id.Encode(w)
	if err != nil {
		return
	}
	switch ret.Id {
	case 0:
		PlayToClientPacketHideMessageSignature, ok := ret.Signature.([256]byte)
		if !ok {
			err = queser.BadTypeError
			return
		}
		arr := PlayToClientPacketHideMessageSignature
		_, err = w.Write(arr[:])
		if err != nil {
			return
		}
	default:
		_, ok := ret.Signature.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Signature.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketHurtAnimation struct {
	EntityId queser.VarInt
	Yaw      float32
}

func (_ PlayToClientPacketHurtAnimation) Decode(r io.Reader) (ret PlayToClientPacketHurtAnimation, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketHurtAnimation) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketInitializeWorldBorder struct {
	X                      float64
	Z                      float64
	OldDiameter            float64
	NewDiameter            float64
	Speed                  queser.VarInt
	PortalTeleportBoundary queser.VarInt
	WarningBlocks          queser.VarInt
	WarningTime            queser.VarInt
}

func (_ PlayToClientPacketInitializeWorldBorder) Decode(r io.Reader) (ret PlayToClientPacketInitializeWorldBorder, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OldDiameter)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.NewDiameter)
	if err != nil {
		return
	}
	ret.Speed, err = ret.Speed.Decode(r)
	if err != nil {
		return
	}
	ret.PortalTeleportBoundary, err = ret.PortalTeleportBoundary.Decode(r)
	if err != nil {
		return
	}
	ret.WarningBlocks, err = ret.WarningBlocks.Decode(r)
	if err != nil {
		return
	}
	ret.WarningTime, err = ret.WarningTime.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketInitializeWorldBorder) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OldDiameter)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.NewDiameter)
	if err != nil {
		return
	}
	err = ret.Speed.Encode(w)
	if err != nil {
		return
	}
	err = ret.PortalTeleportBoundary.Encode(w)
	if err != nil {
		return
	}
	err = ret.WarningBlocks.Encode(w)
	if err != nil {
		return
	}
	err = ret.WarningTime.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketKeepAlive struct {
	KeepAliveId int64
}

func (_ PlayToClientPacketKeepAlive) Decode(r io.Reader) (ret PlayToClientPacketKeepAlive, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.KeepAliveId)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketKeepAlive) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.KeepAliveId)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketKickDisconnect struct {
	Reason nbt.Anon
}

func (_ PlayToClientPacketKickDisconnect) Decode(r io.Reader) (ret PlayToClientPacketKickDisconnect, err error) {
	ret.Reason, err = ret.Reason.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketKickDisconnect) Encode(w io.Writer) (err error) {
	err = ret.Reason.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketLogin struct {
	EntityId            int32
	IsHardcore          bool
	WorldNames          []string
	MaxPlayers          queser.VarInt
	ViewDistance        queser.VarInt
	SimulationDistance  queser.VarInt
	ReducedDebugInfo    bool
	EnableRespawnScreen bool
	DoLimitedCrafting   bool
	WorldState          PlayToClientSpawnInfo
	EnforcesSecureChat  bool
}

func (_ PlayToClientPacketLogin) Decode(r io.Reader) (ret PlayToClientPacketLogin, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.EntityId)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IsHardcore)
	if err != nil {
		return
	}
	var lPlayToClientPacketLoginWorldNames queser.VarInt
	lPlayToClientPacketLoginWorldNames, err = lPlayToClientPacketLoginWorldNames.Decode(r)
	if err != nil {
		return
	}
	ret.WorldNames = []string{}
	for range lPlayToClientPacketLoginWorldNames {
		var PlayToClientPacketLoginWorldNamesElement string
		PlayToClientPacketLoginWorldNamesElement, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.WorldNames = append(ret.WorldNames, PlayToClientPacketLoginWorldNamesElement)
	}
	ret.MaxPlayers, err = ret.MaxPlayers.Decode(r)
	if err != nil {
		return
	}
	ret.ViewDistance, err = ret.ViewDistance.Decode(r)
	if err != nil {
		return
	}
	ret.SimulationDistance, err = ret.SimulationDistance.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.ReducedDebugInfo)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.EnableRespawnScreen)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.DoLimitedCrafting)
	if err != nil {
		return
	}
	ret.WorldState, err = ret.WorldState.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.EnforcesSecureChat)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketLogin) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.EntityId)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IsHardcore)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.WorldNames)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketLoginWorldNames := range len(ret.WorldNames) {
		err = queser.EncodeString(w, ret.WorldNames[iPlayToClientPacketLoginWorldNames])
		if err != nil {
			return
		}
	}
	err = ret.MaxPlayers.Encode(w)
	if err != nil {
		return
	}
	err = ret.ViewDistance.Encode(w)
	if err != nil {
		return
	}
	err = ret.SimulationDistance.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.ReducedDebugInfo)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.EnableRespawnScreen)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.DoLimitedCrafting)
	if err != nil {
		return
	}
	err = ret.WorldState.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.EnforcesSecureChat)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketMap struct {
	ItemDamage queser.VarInt
	Scale      int8
	Locked     bool
	Icons      *[]struct {
		Type        queser.VarInt
		X           int8
		Z           int8
		Direction   uint8
		DisplayName *nbt.Anon
	}
	Columns uint8
	Rows    any
	X       any
	Y       any
	Data    any
}

func (_ PlayToClientPacketMap) Decode(r io.Reader) (ret PlayToClientPacketMap, err error) {
	ret.ItemDamage, err = ret.ItemDamage.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Scale)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Locked)
	if err != nil {
		return
	}
	var PlayToClientPacketMapIconsPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapIconsPresent)
	if err != nil {
		return
	}
	if PlayToClientPacketMapIconsPresent {
		var PlayToClientPacketMapIconsPresentValue []struct {
			Type        queser.VarInt
			X           int8
			Z           int8
			Direction   uint8
			DisplayName *nbt.Anon
		}
		var lPlayToClientPacketMapIcons queser.VarInt
		lPlayToClientPacketMapIcons, err = lPlayToClientPacketMapIcons.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketMapIconsPresentValue = []struct {
			Type        queser.VarInt
			X           int8
			Z           int8
			Direction   uint8
			DisplayName *nbt.Anon
		}{}
		for range lPlayToClientPacketMapIcons {
			var PlayToClientPacketMapIconsElement struct {
				Type        queser.VarInt
				X           int8
				Z           int8
				Direction   uint8
				DisplayName *nbt.Anon
			}
			PlayToClientPacketMapIconsElement.Type, err = PlayToClientPacketMapIconsElement.Type.Decode(r)
			if err != nil {
				return
			}
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapIconsElement.X)
			if err != nil {
				return
			}
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapIconsElement.Z)
			if err != nil {
				return
			}
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapIconsElement.Direction)
			if err != nil {
				return
			}
			var PlayToClientPacketMapIconsElementDisplayNamePresent bool
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapIconsElementDisplayNamePresent)
			if err != nil {
				return
			}
			if PlayToClientPacketMapIconsElementDisplayNamePresent {
				var PlayToClientPacketMapIconsElementDisplayNamePresentValue nbt.Anon
				PlayToClientPacketMapIconsElementDisplayNamePresentValue, err = PlayToClientPacketMapIconsElementDisplayNamePresentValue.Decode(r)
				if err != nil {
					return
				}
				PlayToClientPacketMapIconsElement.DisplayName = &PlayToClientPacketMapIconsElementDisplayNamePresentValue
			}
			PlayToClientPacketMapIconsPresentValue = append(PlayToClientPacketMapIconsPresentValue, PlayToClientPacketMapIconsElement)
		}
		ret.Icons = &PlayToClientPacketMapIconsPresentValue
	}
	err = binary.Read(r, binary.BigEndian, &ret.Columns)
	if err != nil {
		return
	}
	switch ret.Columns {
	case 0:
		var PlayToClientPacketMapRowsTmp queser.Void
		PlayToClientPacketMapRowsTmp, err = PlayToClientPacketMapRowsTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Rows = PlayToClientPacketMapRowsTmp
	default:
		var PlayToClientPacketMapRowsTmp uint8
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapRowsTmp)
		if err != nil {
			return
		}
		ret.Rows = PlayToClientPacketMapRowsTmp
	}
	switch ret.Columns {
	case 0:
		var PlayToClientPacketMapXTmp queser.Void
		PlayToClientPacketMapXTmp, err = PlayToClientPacketMapXTmp.Decode(r)
		if err != nil {
			return
		}
		ret.X = PlayToClientPacketMapXTmp
	default:
		var PlayToClientPacketMapXTmp uint8
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapXTmp)
		if err != nil {
			return
		}
		ret.X = PlayToClientPacketMapXTmp
	}
	switch ret.Columns {
	case 0:
		var PlayToClientPacketMapYTmp queser.Void
		PlayToClientPacketMapYTmp, err = PlayToClientPacketMapYTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Y = PlayToClientPacketMapYTmp
	default:
		var PlayToClientPacketMapYTmp uint8
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapYTmp)
		if err != nil {
			return
		}
		ret.Y = PlayToClientPacketMapYTmp
	}
	switch ret.Columns {
	case 0:
		var PlayToClientPacketMapDataTmp queser.Void
		PlayToClientPacketMapDataTmp, err = PlayToClientPacketMapDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Data = PlayToClientPacketMapDataTmp
	default:
		var PlayToClientPacketMapDataTmp []byte
		var lPlayToClientPacketMapData queser.VarInt
		lPlayToClientPacketMapData, err = lPlayToClientPacketMapData.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketMapDataTmp, err = io.ReadAll(io.LimitReader(r, int64(lPlayToClientPacketMapData)))
		if err != nil {
			return
		}
		ret.Data = PlayToClientPacketMapDataTmp
	}
	return
}
func (ret PlayToClientPacketMap) Encode(w io.Writer) (err error) {
	err = ret.ItemDamage.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Scale)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Locked)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Icons != nil)
	if err != nil {
		return
	}
	if ret.Icons != nil {
		err = queser.VarInt(len(*ret.Icons)).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketMapIcons := range len(*ret.Icons) {
			err = (*ret.Icons)[iPlayToClientPacketMapIcons].Type.Encode(w)
			if err != nil {
				return
			}
			err = binary.Write(w, binary.BigEndian, (*ret.Icons)[iPlayToClientPacketMapIcons].X)
			if err != nil {
				return
			}
			err = binary.Write(w, binary.BigEndian, (*ret.Icons)[iPlayToClientPacketMapIcons].Z)
			if err != nil {
				return
			}
			err = binary.Write(w, binary.BigEndian, (*ret.Icons)[iPlayToClientPacketMapIcons].Direction)
			if err != nil {
				return
			}
			err = binary.Write(w, binary.BigEndian, (*ret.Icons)[iPlayToClientPacketMapIcons].DisplayName != nil)
			if err != nil {
				return
			}
			if (*ret.Icons)[iPlayToClientPacketMapIcons].DisplayName != nil {
				err = (*(*ret.Icons)[iPlayToClientPacketMapIcons].DisplayName).Encode(w)
				if err != nil {
					return
				}
			}
		}
	}
	err = binary.Write(w, binary.BigEndian, ret.Columns)
	if err != nil {
		return
	}
	switch ret.Columns {
	case 0:
		PlayToClientPacketMapRows, ok := ret.Rows.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketMapRows.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Rows.(uint8)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Rows.(uint8))
		if err != nil {
			return
		}
	}
	switch ret.Columns {
	case 0:
		PlayToClientPacketMapX, ok := ret.X.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketMapX.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.X.(uint8)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.X.(uint8))
		if err != nil {
			return
		}
	}
	switch ret.Columns {
	case 0:
		PlayToClientPacketMapY, ok := ret.Y.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketMapY.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Y.(uint8)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Y.(uint8))
		if err != nil {
			return
		}
	}
	switch ret.Columns {
	case 0:
		PlayToClientPacketMapData, ok := ret.Data.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketMapData.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Data.([]byte)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.VarInt(len(ret.Data.([]byte))).Encode(w)
		if err != nil {
			return
		}
		_, err = w.Write(ret.Data.([]byte))
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketMapChunk struct {
	X          int32
	Z          int32
	Heightmaps []struct {
		Type string
		Data []int64
	}
	ChunkData           ByteArray
	BlockEntities       []ChunkBlockEntity
	SkyLightMask        []int64
	BlockLightMask      []int64
	EmptySkyLightMask   []int64
	EmptyBlockLightMask []int64
	SkyLight            [][]uint8
	BlockLight          [][]uint8
}

var PlayToClientPacketMapChunkHeightmapsElementTypeMap = map[queser.VarInt]string{0: "world_surface_wg", 1: "world_surface", 2: "ocean_floor_wg", 3: "ocean_floor", 4: "motion_blocking", 5: "motion_blocking_no_leaves"}

func (_ PlayToClientPacketMapChunk) Decode(r io.Reader) (ret PlayToClientPacketMapChunk, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	var lPlayToClientPacketMapChunkHeightmaps queser.VarInt
	lPlayToClientPacketMapChunkHeightmaps, err = lPlayToClientPacketMapChunkHeightmaps.Decode(r)
	if err != nil {
		return
	}
	ret.Heightmaps = []struct {
		Type string
		Data []int64
	}{}
	for range lPlayToClientPacketMapChunkHeightmaps {
		var PlayToClientPacketMapChunkHeightmapsElement struct {
			Type string
			Data []int64
		}
		var PlayToClientPacketMapChunkHeightmapsElementTypeKey queser.VarInt
		PlayToClientPacketMapChunkHeightmapsElementTypeKey, err = PlayToClientPacketMapChunkHeightmapsElementTypeKey.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketMapChunkHeightmapsElement.Type, err = queser.ErroringIndex(PlayToClientPacketMapChunkHeightmapsElementTypeMap, PlayToClientPacketMapChunkHeightmapsElementTypeKey)
		if err != nil {
			return
		}
		var lPlayToClientPacketMapChunkHeightmapsElementData queser.VarInt
		lPlayToClientPacketMapChunkHeightmapsElementData, err = lPlayToClientPacketMapChunkHeightmapsElementData.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketMapChunkHeightmapsElement.Data = []int64{}
		for range lPlayToClientPacketMapChunkHeightmapsElementData {
			var PlayToClientPacketMapChunkHeightmapsElementDataElement int64
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapChunkHeightmapsElementDataElement)
			if err != nil {
				return
			}
			PlayToClientPacketMapChunkHeightmapsElement.Data = append(PlayToClientPacketMapChunkHeightmapsElement.Data, PlayToClientPacketMapChunkHeightmapsElementDataElement)
		}
		ret.Heightmaps = append(ret.Heightmaps, PlayToClientPacketMapChunkHeightmapsElement)
	}
	ret.ChunkData, err = ret.ChunkData.Decode(r)
	if err != nil {
		return
	}
	var lPlayToClientPacketMapChunkBlockEntities queser.VarInt
	lPlayToClientPacketMapChunkBlockEntities, err = lPlayToClientPacketMapChunkBlockEntities.Decode(r)
	if err != nil {
		return
	}
	ret.BlockEntities = []ChunkBlockEntity{}
	for range lPlayToClientPacketMapChunkBlockEntities {
		var PlayToClientPacketMapChunkBlockEntitiesElement ChunkBlockEntity
		PlayToClientPacketMapChunkBlockEntitiesElement, err = PlayToClientPacketMapChunkBlockEntitiesElement.Decode(r)
		if err != nil {
			return
		}
		ret.BlockEntities = append(ret.BlockEntities, PlayToClientPacketMapChunkBlockEntitiesElement)
	}
	var lPlayToClientPacketMapChunkSkyLightMask queser.VarInt
	lPlayToClientPacketMapChunkSkyLightMask, err = lPlayToClientPacketMapChunkSkyLightMask.Decode(r)
	if err != nil {
		return
	}
	ret.SkyLightMask = []int64{}
	for range lPlayToClientPacketMapChunkSkyLightMask {
		var PlayToClientPacketMapChunkSkyLightMaskElement int64
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapChunkSkyLightMaskElement)
		if err != nil {
			return
		}
		ret.SkyLightMask = append(ret.SkyLightMask, PlayToClientPacketMapChunkSkyLightMaskElement)
	}
	var lPlayToClientPacketMapChunkBlockLightMask queser.VarInt
	lPlayToClientPacketMapChunkBlockLightMask, err = lPlayToClientPacketMapChunkBlockLightMask.Decode(r)
	if err != nil {
		return
	}
	ret.BlockLightMask = []int64{}
	for range lPlayToClientPacketMapChunkBlockLightMask {
		var PlayToClientPacketMapChunkBlockLightMaskElement int64
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapChunkBlockLightMaskElement)
		if err != nil {
			return
		}
		ret.BlockLightMask = append(ret.BlockLightMask, PlayToClientPacketMapChunkBlockLightMaskElement)
	}
	var lPlayToClientPacketMapChunkEmptySkyLightMask queser.VarInt
	lPlayToClientPacketMapChunkEmptySkyLightMask, err = lPlayToClientPacketMapChunkEmptySkyLightMask.Decode(r)
	if err != nil {
		return
	}
	ret.EmptySkyLightMask = []int64{}
	for range lPlayToClientPacketMapChunkEmptySkyLightMask {
		var PlayToClientPacketMapChunkEmptySkyLightMaskElement int64
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapChunkEmptySkyLightMaskElement)
		if err != nil {
			return
		}
		ret.EmptySkyLightMask = append(ret.EmptySkyLightMask, PlayToClientPacketMapChunkEmptySkyLightMaskElement)
	}
	var lPlayToClientPacketMapChunkEmptyBlockLightMask queser.VarInt
	lPlayToClientPacketMapChunkEmptyBlockLightMask, err = lPlayToClientPacketMapChunkEmptyBlockLightMask.Decode(r)
	if err != nil {
		return
	}
	ret.EmptyBlockLightMask = []int64{}
	for range lPlayToClientPacketMapChunkEmptyBlockLightMask {
		var PlayToClientPacketMapChunkEmptyBlockLightMaskElement int64
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapChunkEmptyBlockLightMaskElement)
		if err != nil {
			return
		}
		ret.EmptyBlockLightMask = append(ret.EmptyBlockLightMask, PlayToClientPacketMapChunkEmptyBlockLightMaskElement)
	}
	var lPlayToClientPacketMapChunkSkyLight queser.VarInt
	lPlayToClientPacketMapChunkSkyLight, err = lPlayToClientPacketMapChunkSkyLight.Decode(r)
	if err != nil {
		return
	}
	ret.SkyLight = [][]uint8{}
	for range lPlayToClientPacketMapChunkSkyLight {
		var PlayToClientPacketMapChunkSkyLightElement []uint8
		var lPlayToClientPacketMapChunkSkyLightElement queser.VarInt
		lPlayToClientPacketMapChunkSkyLightElement, err = lPlayToClientPacketMapChunkSkyLightElement.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketMapChunkSkyLightElement = []uint8{}
		for range lPlayToClientPacketMapChunkSkyLightElement {
			var PlayToClientPacketMapChunkSkyLightElementElement uint8
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapChunkSkyLightElementElement)
			if err != nil {
				return
			}
			PlayToClientPacketMapChunkSkyLightElement = append(PlayToClientPacketMapChunkSkyLightElement, PlayToClientPacketMapChunkSkyLightElementElement)
		}
		ret.SkyLight = append(ret.SkyLight, PlayToClientPacketMapChunkSkyLightElement)
	}
	var lPlayToClientPacketMapChunkBlockLight queser.VarInt
	lPlayToClientPacketMapChunkBlockLight, err = lPlayToClientPacketMapChunkBlockLight.Decode(r)
	if err != nil {
		return
	}
	ret.BlockLight = [][]uint8{}
	for range lPlayToClientPacketMapChunkBlockLight {
		var PlayToClientPacketMapChunkBlockLightElement []uint8
		var lPlayToClientPacketMapChunkBlockLightElement queser.VarInt
		lPlayToClientPacketMapChunkBlockLightElement, err = lPlayToClientPacketMapChunkBlockLightElement.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketMapChunkBlockLightElement = []uint8{}
		for range lPlayToClientPacketMapChunkBlockLightElement {
			var PlayToClientPacketMapChunkBlockLightElementElement uint8
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMapChunkBlockLightElementElement)
			if err != nil {
				return
			}
			PlayToClientPacketMapChunkBlockLightElement = append(PlayToClientPacketMapChunkBlockLightElement, PlayToClientPacketMapChunkBlockLightElementElement)
		}
		ret.BlockLight = append(ret.BlockLight, PlayToClientPacketMapChunkBlockLightElement)
	}
	return
}

var PlayToClientPacketMapChunkHeightmapsTypeReverseMap = map[string]queser.VarInt{"world_surface_wg": 0, "world_surface": 1, "ocean_floor_wg": 2, "ocean_floor": 3, "motion_blocking": 4, "motion_blocking_no_leaves": 5}
var PlayToClientPacketMapChunkHeightmapsInnerTypeReverseMap = map[string]queser.VarInt{"world_surface_wg": 0, "world_surface": 1, "ocean_floor_wg": 2, "ocean_floor": 3, "motion_blocking": 4, "motion_blocking_no_leaves": 5}

func (ret PlayToClientPacketMapChunk) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Heightmaps)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketMapChunkHeightmaps := range len(ret.Heightmaps) {
		var vPlayToClientPacketMapChunkHeightmapsInnerType queser.VarInt
		vPlayToClientPacketMapChunkHeightmapsInnerType, err = queser.ErroringIndex(PlayToClientPacketMapChunkHeightmapsInnerTypeReverseMap, ret.Heightmaps[iPlayToClientPacketMapChunkHeightmaps].Type)
		if err != nil {
			return
		}
		err = vPlayToClientPacketMapChunkHeightmapsInnerType.Encode(w)
		if err != nil {
			return
		}
		err = queser.VarInt(len(ret.Heightmaps[iPlayToClientPacketMapChunkHeightmaps].Data)).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketMapChunkHeightmapsInnerData := range len(ret.Heightmaps[iPlayToClientPacketMapChunkHeightmaps].Data) {
			err = binary.Write(w, binary.BigEndian, ret.Heightmaps[iPlayToClientPacketMapChunkHeightmaps].Data[iPlayToClientPacketMapChunkHeightmapsInnerData])
			if err != nil {
				return
			}
		}
	}
	err = ret.ChunkData.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.BlockEntities)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketMapChunkBlockEntities := range len(ret.BlockEntities) {
		err = ret.BlockEntities[iPlayToClientPacketMapChunkBlockEntities].Encode(w)
		if err != nil {
			return
		}
	}
	err = queser.VarInt(len(ret.SkyLightMask)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketMapChunkSkyLightMask := range len(ret.SkyLightMask) {
		err = binary.Write(w, binary.BigEndian, ret.SkyLightMask[iPlayToClientPacketMapChunkSkyLightMask])
		if err != nil {
			return
		}
	}
	err = queser.VarInt(len(ret.BlockLightMask)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketMapChunkBlockLightMask := range len(ret.BlockLightMask) {
		err = binary.Write(w, binary.BigEndian, ret.BlockLightMask[iPlayToClientPacketMapChunkBlockLightMask])
		if err != nil {
			return
		}
	}
	err = queser.VarInt(len(ret.EmptySkyLightMask)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketMapChunkEmptySkyLightMask := range len(ret.EmptySkyLightMask) {
		err = binary.Write(w, binary.BigEndian, ret.EmptySkyLightMask[iPlayToClientPacketMapChunkEmptySkyLightMask])
		if err != nil {
			return
		}
	}
	err = queser.VarInt(len(ret.EmptyBlockLightMask)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketMapChunkEmptyBlockLightMask := range len(ret.EmptyBlockLightMask) {
		err = binary.Write(w, binary.BigEndian, ret.EmptyBlockLightMask[iPlayToClientPacketMapChunkEmptyBlockLightMask])
		if err != nil {
			return
		}
	}
	err = queser.VarInt(len(ret.SkyLight)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketMapChunkSkyLight := range len(ret.SkyLight) {
		err = queser.VarInt(len(ret.SkyLight[iPlayToClientPacketMapChunkSkyLight])).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketMapChunkSkyLightInner := range len(ret.SkyLight[iPlayToClientPacketMapChunkSkyLight]) {
			err = binary.Write(w, binary.BigEndian, ret.SkyLight[iPlayToClientPacketMapChunkSkyLight][iPlayToClientPacketMapChunkSkyLightInner])
			if err != nil {
				return
			}
		}
	}
	err = queser.VarInt(len(ret.BlockLight)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketMapChunkBlockLight := range len(ret.BlockLight) {
		err = queser.VarInt(len(ret.BlockLight[iPlayToClientPacketMapChunkBlockLight])).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketMapChunkBlockLightInner := range len(ret.BlockLight[iPlayToClientPacketMapChunkBlockLight]) {
			err = binary.Write(w, binary.BigEndian, ret.BlockLight[iPlayToClientPacketMapChunkBlockLight][iPlayToClientPacketMapChunkBlockLightInner])
			if err != nil {
				return
			}
		}
	}
	return
}

type PlayToClientPacketMoveMinecart struct {
	EntityId queser.VarInt
	Steps    []struct {
		Position Vec3f
		Movement Vec3f
		Yaw      float32
		Pitch    float32
		Weight   float32
	}
}

func (_ PlayToClientPacketMoveMinecart) Decode(r io.Reader) (ret PlayToClientPacketMoveMinecart, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	var lPlayToClientPacketMoveMinecartSteps queser.VarInt
	lPlayToClientPacketMoveMinecartSteps, err = lPlayToClientPacketMoveMinecartSteps.Decode(r)
	if err != nil {
		return
	}
	ret.Steps = []struct {
		Position Vec3f
		Movement Vec3f
		Yaw      float32
		Pitch    float32
		Weight   float32
	}{}
	for range lPlayToClientPacketMoveMinecartSteps {
		var PlayToClientPacketMoveMinecartStepsElement struct {
			Position Vec3f
			Movement Vec3f
			Yaw      float32
			Pitch    float32
			Weight   float32
		}
		PlayToClientPacketMoveMinecartStepsElement.Position, err = PlayToClientPacketMoveMinecartStepsElement.Position.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketMoveMinecartStepsElement.Movement, err = PlayToClientPacketMoveMinecartStepsElement.Movement.Decode(r)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMoveMinecartStepsElement.Yaw)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMoveMinecartStepsElement.Pitch)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketMoveMinecartStepsElement.Weight)
		if err != nil {
			return
		}
		ret.Steps = append(ret.Steps, PlayToClientPacketMoveMinecartStepsElement)
	}
	return
}
func (ret PlayToClientPacketMoveMinecart) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Steps)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketMoveMinecartSteps := range len(ret.Steps) {
		err = ret.Steps[iPlayToClientPacketMoveMinecartSteps].Position.Encode(w)
		if err != nil {
			return
		}
		err = ret.Steps[iPlayToClientPacketMoveMinecartSteps].Movement.Encode(w)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Steps[iPlayToClientPacketMoveMinecartSteps].Yaw)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Steps[iPlayToClientPacketMoveMinecartSteps].Pitch)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Steps[iPlayToClientPacketMoveMinecartSteps].Weight)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketMultiBlockChange struct {
	ChunkCoordinates uint64
	Records          []queser.VarInt
}

func (_ PlayToClientPacketMultiBlockChange) Decode(r io.Reader) (ret PlayToClientPacketMultiBlockChange, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.ChunkCoordinates)
	if err != nil {
		return
	}
	var lPlayToClientPacketMultiBlockChangeRecords queser.VarInt
	lPlayToClientPacketMultiBlockChangeRecords, err = lPlayToClientPacketMultiBlockChangeRecords.Decode(r)
	if err != nil {
		return
	}
	ret.Records = []queser.VarInt{}
	for range lPlayToClientPacketMultiBlockChangeRecords {
		var PlayToClientPacketMultiBlockChangeRecordsElement queser.VarInt
		PlayToClientPacketMultiBlockChangeRecordsElement, err = PlayToClientPacketMultiBlockChangeRecordsElement.Decode(r)
		if err != nil {
			return
		}
		ret.Records = append(ret.Records, PlayToClientPacketMultiBlockChangeRecordsElement)
	}
	return
}
func (ret PlayToClientPacketMultiBlockChange) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.ChunkCoordinates)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Records)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketMultiBlockChangeRecords := range len(ret.Records) {
		err = ret.Records[iPlayToClientPacketMultiBlockChangeRecords].Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketNbtQueryResponse struct {
	TransactionId queser.VarInt
	Nbt           nbt.Anon
}

func (_ PlayToClientPacketNbtQueryResponse) Decode(r io.Reader) (ret PlayToClientPacketNbtQueryResponse, err error) {
	ret.TransactionId, err = ret.TransactionId.Decode(r)
	if err != nil {
		return
	}
	ret.Nbt, err = ret.Nbt.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketNbtQueryResponse) Encode(w io.Writer) (err error) {
	err = ret.TransactionId.Encode(w)
	if err != nil {
		return
	}
	err = ret.Nbt.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketOpenBook struct {
	Hand queser.VarInt
}

func (_ PlayToClientPacketOpenBook) Decode(r io.Reader) (ret PlayToClientPacketOpenBook, err error) {
	ret.Hand, err = ret.Hand.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketOpenBook) Encode(w io.Writer) (err error) {
	err = ret.Hand.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketOpenHorseWindow struct {
	WindowId ContainerID
	NbSlots  queser.VarInt
	EntityId int32
}

func (_ PlayToClientPacketOpenHorseWindow) Decode(r io.Reader) (ret PlayToClientPacketOpenHorseWindow, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	ret.NbSlots, err = ret.NbSlots.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.EntityId)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketOpenHorseWindow) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	err = ret.NbSlots.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.EntityId)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketOpenSignEntity struct {
	Location    Position
	IsFrontText bool
}

func (_ PlayToClientPacketOpenSignEntity) Decode(r io.Reader) (ret PlayToClientPacketOpenSignEntity, err error) {
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IsFrontText)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketOpenSignEntity) Encode(w io.Writer) (err error) {
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IsFrontText)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketOpenWindow struct {
	WindowId      queser.VarInt
	InventoryType queser.VarInt
	WindowTitle   nbt.Anon
}

func (_ PlayToClientPacketOpenWindow) Decode(r io.Reader) (ret PlayToClientPacketOpenWindow, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	ret.InventoryType, err = ret.InventoryType.Decode(r)
	if err != nil {
		return
	}
	ret.WindowTitle, err = ret.WindowTitle.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketOpenWindow) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	err = ret.InventoryType.Encode(w)
	if err != nil {
		return
	}
	err = ret.WindowTitle.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketPing struct {
	Id int32
}

func (_ PlayToClientPacketPing) Decode(r io.Reader) (ret PlayToClientPacketPing, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Id)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketPing) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Id)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketPingResponse struct {
	Id int64
}

func (_ PlayToClientPacketPingResponse) Decode(r io.Reader) (ret PlayToClientPacketPingResponse, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Id)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketPingResponse) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Id)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketPlayerChat struct {
	GlobalIndex         queser.VarInt
	SenderUuid          uuid.UUID
	Index               queser.VarInt
	Signature           *[256]byte
	PlainMessage        string
	Timestamp           int64
	Salt                int64
	PreviousMessages    PreviousMessages
	UnsignedChatContent *nbt.Anon
	FilterType          queser.VarInt
	FilterTypeMask      any
	Type                PlayToClientChatTypesHolder
	NetworkName         nbt.Anon
	NetworkTargetName   *nbt.Anon
}

func (_ PlayToClientPacketPlayerChat) Decode(r io.Reader) (ret PlayToClientPacketPlayerChat, err error) {
	ret.GlobalIndex, err = ret.GlobalIndex.Decode(r)
	if err != nil {
		return
	}
	_, err = io.ReadFull(r, ret.SenderUuid[:])
	if err != nil {
		return
	}
	ret.Index, err = ret.Index.Decode(r)
	if err != nil {
		return
	}
	var PlayToClientPacketPlayerChatSignaturePresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketPlayerChatSignaturePresent)
	if err != nil {
		return
	}
	if PlayToClientPacketPlayerChatSignaturePresent {
		var PlayToClientPacketPlayerChatSignaturePresentValue [256]byte
		_, err = r.Read(PlayToClientPacketPlayerChatSignaturePresentValue[:])
		if err != nil {
			return
		}
		ret.Signature = &PlayToClientPacketPlayerChatSignaturePresentValue
	}
	ret.PlainMessage, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Timestamp)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Salt)
	if err != nil {
		return
	}
	ret.PreviousMessages, err = ret.PreviousMessages.Decode(r)
	if err != nil {
		return
	}
	var PlayToClientPacketPlayerChatUnsignedChatContentPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketPlayerChatUnsignedChatContentPresent)
	if err != nil {
		return
	}
	if PlayToClientPacketPlayerChatUnsignedChatContentPresent {
		var PlayToClientPacketPlayerChatUnsignedChatContentPresentValue nbt.Anon
		PlayToClientPacketPlayerChatUnsignedChatContentPresentValue, err = PlayToClientPacketPlayerChatUnsignedChatContentPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.UnsignedChatContent = &PlayToClientPacketPlayerChatUnsignedChatContentPresentValue
	}
	ret.FilterType, err = ret.FilterType.Decode(r)
	if err != nil {
		return
	}
	switch ret.FilterType {
	case 2:
		var PlayToClientPacketPlayerChatFilterTypeMaskTmp []int64
		var lPlayToClientPacketPlayerChatFilterTypeMask queser.VarInt
		lPlayToClientPacketPlayerChatFilterTypeMask, err = lPlayToClientPacketPlayerChatFilterTypeMask.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketPlayerChatFilterTypeMaskTmp = []int64{}
		for range lPlayToClientPacketPlayerChatFilterTypeMask {
			var PlayToClientPacketPlayerChatFilterTypeMaskElement int64
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketPlayerChatFilterTypeMaskElement)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerChatFilterTypeMaskTmp = append(PlayToClientPacketPlayerChatFilterTypeMaskTmp, PlayToClientPacketPlayerChatFilterTypeMaskElement)
		}
		ret.FilterTypeMask = PlayToClientPacketPlayerChatFilterTypeMaskTmp
	default:
		var PlayToClientPacketPlayerChatFilterTypeMaskTmp queser.Void
		PlayToClientPacketPlayerChatFilterTypeMaskTmp, err = PlayToClientPacketPlayerChatFilterTypeMaskTmp.Decode(r)
		if err != nil {
			return
		}
		ret.FilterTypeMask = PlayToClientPacketPlayerChatFilterTypeMaskTmp
	}
	ret.Type, err = ret.Type.Decode(r)
	if err != nil {
		return
	}
	ret.NetworkName, err = ret.NetworkName.Decode(r)
	if err != nil {
		return
	}
	var PlayToClientPacketPlayerChatNetworkTargetNamePresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketPlayerChatNetworkTargetNamePresent)
	if err != nil {
		return
	}
	if PlayToClientPacketPlayerChatNetworkTargetNamePresent {
		var PlayToClientPacketPlayerChatNetworkTargetNamePresentValue nbt.Anon
		PlayToClientPacketPlayerChatNetworkTargetNamePresentValue, err = PlayToClientPacketPlayerChatNetworkTargetNamePresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.NetworkTargetName = &PlayToClientPacketPlayerChatNetworkTargetNamePresentValue
	}
	return
}
func (ret PlayToClientPacketPlayerChat) Encode(w io.Writer) (err error) {
	err = ret.GlobalIndex.Encode(w)
	if err != nil {
		return
	}
	_, err = w.Write(ret.SenderUuid[:])
	if err != nil {
		return
	}
	err = ret.Index.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Signature != nil)
	if err != nil {
		return
	}
	if ret.Signature != nil {
		arr := *ret.Signature
		_, err = w.Write(arr[:])
		if err != nil {
			return
		}
	}
	err = queser.EncodeString(w, ret.PlainMessage)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Timestamp)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Salt)
	if err != nil {
		return
	}
	err = ret.PreviousMessages.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.UnsignedChatContent != nil)
	if err != nil {
		return
	}
	if ret.UnsignedChatContent != nil {
		err = (*ret.UnsignedChatContent).Encode(w)
		if err != nil {
			return
		}
	}
	err = ret.FilterType.Encode(w)
	if err != nil {
		return
	}
	switch ret.FilterType {
	case 2:
		PlayToClientPacketPlayerChatFilterTypeMask, ok := ret.FilterTypeMask.([]int64)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.VarInt(len(PlayToClientPacketPlayerChatFilterTypeMask)).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketPlayerChatFilterTypeMask := range len(PlayToClientPacketPlayerChatFilterTypeMask) {
			err = binary.Write(w, binary.BigEndian, PlayToClientPacketPlayerChatFilterTypeMask[iPlayToClientPacketPlayerChatFilterTypeMask])
			if err != nil {
				return
			}
		}
	default:
		_, ok := ret.FilterTypeMask.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.FilterTypeMask.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	err = ret.Type.Encode(w)
	if err != nil {
		return
	}
	err = ret.NetworkName.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.NetworkTargetName != nil)
	if err != nil {
		return
	}
	if ret.NetworkTargetName != nil {
		err = (*ret.NetworkTargetName).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketPlayerInfo struct {
	Action uint8
	Data   []struct {
		Uuid         uuid.UUID
		Player       any
		ChatSession  any
		Gamemode     any
		Listed       any
		Latency      any
		DisplayName  any
		ListPriority any
		ShowHat      any
	}
}

func (_ PlayToClientPacketPlayerInfo) Decode(r io.Reader) (ret PlayToClientPacketPlayerInfo, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Action)
	if err != nil {
		return
	}
	var lPlayToClientPacketPlayerInfoData queser.VarInt
	lPlayToClientPacketPlayerInfoData, err = lPlayToClientPacketPlayerInfoData.Decode(r)
	if err != nil {
		return
	}
	ret.Data = []struct {
		Uuid         uuid.UUID
		Player       any
		ChatSession  any
		Gamemode     any
		Listed       any
		Latency      any
		DisplayName  any
		ListPriority any
		ShowHat      any
	}{}
	for range lPlayToClientPacketPlayerInfoData {
		var PlayToClientPacketPlayerInfoDataElement struct {
			Uuid         uuid.UUID
			Player       any
			ChatSession  any
			Gamemode     any
			Listed       any
			Latency      any
			DisplayName  any
			ListPriority any
			ShowHat      any
		}
		_, err = io.ReadFull(r, PlayToClientPacketPlayerInfoDataElement.Uuid[:])
		if err != nil {
			return
		}
		switch ret.Action&1 != 0 {
		case true:
			var PlayToClientPacketPlayerInfoDataElementPlayerTmp GameProfile
			PlayToClientPacketPlayerInfoDataElementPlayerTmp, err = PlayToClientPacketPlayerInfoDataElementPlayerTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.Player = PlayToClientPacketPlayerInfoDataElementPlayerTmp
		default:
			var PlayToClientPacketPlayerInfoDataElementPlayerTmp queser.Void
			PlayToClientPacketPlayerInfoDataElementPlayerTmp, err = PlayToClientPacketPlayerInfoDataElementPlayerTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.Player = PlayToClientPacketPlayerInfoDataElementPlayerTmp
		}
		switch ret.Action&2 != 0 {
		case true:
			var PlayToClientPacketPlayerInfoDataElementChatSessionTmp ChatSession
			PlayToClientPacketPlayerInfoDataElementChatSessionTmp, err = PlayToClientPacketPlayerInfoDataElementChatSessionTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.ChatSession = PlayToClientPacketPlayerInfoDataElementChatSessionTmp
		default:
			var PlayToClientPacketPlayerInfoDataElementChatSessionTmp queser.Void
			PlayToClientPacketPlayerInfoDataElementChatSessionTmp, err = PlayToClientPacketPlayerInfoDataElementChatSessionTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.ChatSession = PlayToClientPacketPlayerInfoDataElementChatSessionTmp
		}
		switch ret.Action&4 != 0 {
		case true:
			var PlayToClientPacketPlayerInfoDataElementGamemodeTmp queser.VarInt
			PlayToClientPacketPlayerInfoDataElementGamemodeTmp, err = PlayToClientPacketPlayerInfoDataElementGamemodeTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.Gamemode = PlayToClientPacketPlayerInfoDataElementGamemodeTmp
		default:
			var PlayToClientPacketPlayerInfoDataElementGamemodeTmp queser.Void
			PlayToClientPacketPlayerInfoDataElementGamemodeTmp, err = PlayToClientPacketPlayerInfoDataElementGamemodeTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.Gamemode = PlayToClientPacketPlayerInfoDataElementGamemodeTmp
		}
		switch ret.Action&8 != 0 {
		case true:
			var PlayToClientPacketPlayerInfoDataElementListedTmp queser.VarInt
			PlayToClientPacketPlayerInfoDataElementListedTmp, err = PlayToClientPacketPlayerInfoDataElementListedTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.Listed = PlayToClientPacketPlayerInfoDataElementListedTmp
		default:
			var PlayToClientPacketPlayerInfoDataElementListedTmp queser.Void
			PlayToClientPacketPlayerInfoDataElementListedTmp, err = PlayToClientPacketPlayerInfoDataElementListedTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.Listed = PlayToClientPacketPlayerInfoDataElementListedTmp
		}
		switch ret.Action&10 != 0 {
		case true:
			var PlayToClientPacketPlayerInfoDataElementLatencyTmp queser.VarInt
			PlayToClientPacketPlayerInfoDataElementLatencyTmp, err = PlayToClientPacketPlayerInfoDataElementLatencyTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.Latency = PlayToClientPacketPlayerInfoDataElementLatencyTmp
		default:
			var PlayToClientPacketPlayerInfoDataElementLatencyTmp queser.Void
			PlayToClientPacketPlayerInfoDataElementLatencyTmp, err = PlayToClientPacketPlayerInfoDataElementLatencyTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.Latency = PlayToClientPacketPlayerInfoDataElementLatencyTmp
		}
		switch ret.Action&20 != 0 {
		case true:
			var PlayToClientPacketPlayerInfoDataElementDisplayNameTmp *nbt.Anon
			var PlayToClientPacketPlayerInfoDataElementDisplayNamePresent bool
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketPlayerInfoDataElementDisplayNamePresent)
			if err != nil {
				return
			}
			if PlayToClientPacketPlayerInfoDataElementDisplayNamePresent {
				var PlayToClientPacketPlayerInfoDataElementDisplayNamePresentValue nbt.Anon
				PlayToClientPacketPlayerInfoDataElementDisplayNamePresentValue, err = PlayToClientPacketPlayerInfoDataElementDisplayNamePresentValue.Decode(r)
				if err != nil {
					return
				}
				PlayToClientPacketPlayerInfoDataElementDisplayNameTmp = &PlayToClientPacketPlayerInfoDataElementDisplayNamePresentValue
			}
			PlayToClientPacketPlayerInfoDataElement.DisplayName = PlayToClientPacketPlayerInfoDataElementDisplayNameTmp
		default:
			var PlayToClientPacketPlayerInfoDataElementDisplayNameTmp queser.Void
			PlayToClientPacketPlayerInfoDataElementDisplayNameTmp, err = PlayToClientPacketPlayerInfoDataElementDisplayNameTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.DisplayName = PlayToClientPacketPlayerInfoDataElementDisplayNameTmp
		}
		switch ret.Action&80 != 0 {
		case true:
			var PlayToClientPacketPlayerInfoDataElementListPriorityTmp queser.VarInt
			PlayToClientPacketPlayerInfoDataElementListPriorityTmp, err = PlayToClientPacketPlayerInfoDataElementListPriorityTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.ListPriority = PlayToClientPacketPlayerInfoDataElementListPriorityTmp
		default:
			var PlayToClientPacketPlayerInfoDataElementListPriorityTmp queser.Void
			PlayToClientPacketPlayerInfoDataElementListPriorityTmp, err = PlayToClientPacketPlayerInfoDataElementListPriorityTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.ListPriority = PlayToClientPacketPlayerInfoDataElementListPriorityTmp
		}
		switch ret.Action&40 != 0 {
		case true:
			var PlayToClientPacketPlayerInfoDataElementShowHatTmp bool
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketPlayerInfoDataElementShowHatTmp)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.ShowHat = PlayToClientPacketPlayerInfoDataElementShowHatTmp
		default:
			var PlayToClientPacketPlayerInfoDataElementShowHatTmp queser.Void
			PlayToClientPacketPlayerInfoDataElementShowHatTmp, err = PlayToClientPacketPlayerInfoDataElementShowHatTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketPlayerInfoDataElement.ShowHat = PlayToClientPacketPlayerInfoDataElementShowHatTmp
		}
		ret.Data = append(ret.Data, PlayToClientPacketPlayerInfoDataElement)
	}
	return
}
func (ret PlayToClientPacketPlayerInfo) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Action)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Data)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketPlayerInfoData := range len(ret.Data) {
		_, err = w.Write(ret.Data[iPlayToClientPacketPlayerInfoData].Uuid[:])
		if err != nil {
			return
		}
		switch ret.Action&1 != 0 {
		case true:
			PlayToClientPacketPlayerInfoDataInnerPlayer, ok := ret.Data[iPlayToClientPacketPlayerInfoData].Player.(GameProfile)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketPlayerInfoDataInnerPlayer.Encode(w)
			if err != nil {
				return
			}
		default:
			_, ok := ret.Data[iPlayToClientPacketPlayerInfoData].Player.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = ret.Data[iPlayToClientPacketPlayerInfoData].Player.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
		switch ret.Action&2 != 0 {
		case true:
			PlayToClientPacketPlayerInfoDataInnerChatSession, ok := ret.Data[iPlayToClientPacketPlayerInfoData].ChatSession.(ChatSession)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketPlayerInfoDataInnerChatSession.Encode(w)
			if err != nil {
				return
			}
		default:
			_, ok := ret.Data[iPlayToClientPacketPlayerInfoData].ChatSession.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = ret.Data[iPlayToClientPacketPlayerInfoData].ChatSession.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
		switch ret.Action&4 != 0 {
		case true:
			PlayToClientPacketPlayerInfoDataInnerGamemode, ok := ret.Data[iPlayToClientPacketPlayerInfoData].Gamemode.(queser.VarInt)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketPlayerInfoDataInnerGamemode.Encode(w)
			if err != nil {
				return
			}
		default:
			_, ok := ret.Data[iPlayToClientPacketPlayerInfoData].Gamemode.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = ret.Data[iPlayToClientPacketPlayerInfoData].Gamemode.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
		switch ret.Action&8 != 0 {
		case true:
			PlayToClientPacketPlayerInfoDataInnerListed, ok := ret.Data[iPlayToClientPacketPlayerInfoData].Listed.(queser.VarInt)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketPlayerInfoDataInnerListed.Encode(w)
			if err != nil {
				return
			}
		default:
			_, ok := ret.Data[iPlayToClientPacketPlayerInfoData].Listed.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = ret.Data[iPlayToClientPacketPlayerInfoData].Listed.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
		switch ret.Action&10 != 0 {
		case true:
			PlayToClientPacketPlayerInfoDataInnerLatency, ok := ret.Data[iPlayToClientPacketPlayerInfoData].Latency.(queser.VarInt)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketPlayerInfoDataInnerLatency.Encode(w)
			if err != nil {
				return
			}
		default:
			_, ok := ret.Data[iPlayToClientPacketPlayerInfoData].Latency.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = ret.Data[iPlayToClientPacketPlayerInfoData].Latency.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
		switch ret.Action&20 != 0 {
		case true:
			PlayToClientPacketPlayerInfoDataInnerDisplayName, ok := ret.Data[iPlayToClientPacketPlayerInfoData].DisplayName.(*nbt.Anon)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = binary.Write(w, binary.BigEndian, PlayToClientPacketPlayerInfoDataInnerDisplayName != nil)
			if err != nil {
				return
			}
			if PlayToClientPacketPlayerInfoDataInnerDisplayName != nil {
				err = (*PlayToClientPacketPlayerInfoDataInnerDisplayName).Encode(w)
				if err != nil {
					return
				}
			}
		default:
			_, ok := ret.Data[iPlayToClientPacketPlayerInfoData].DisplayName.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = ret.Data[iPlayToClientPacketPlayerInfoData].DisplayName.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
		switch ret.Action&80 != 0 {
		case true:
			PlayToClientPacketPlayerInfoDataInnerListPriority, ok := ret.Data[iPlayToClientPacketPlayerInfoData].ListPriority.(queser.VarInt)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketPlayerInfoDataInnerListPriority.Encode(w)
			if err != nil {
				return
			}
		default:
			_, ok := ret.Data[iPlayToClientPacketPlayerInfoData].ListPriority.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = ret.Data[iPlayToClientPacketPlayerInfoData].ListPriority.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
		switch ret.Action&40 != 0 {
		case true:
			PlayToClientPacketPlayerInfoDataInnerShowHat, ok := ret.Data[iPlayToClientPacketPlayerInfoData].ShowHat.(bool)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = binary.Write(w, binary.BigEndian, PlayToClientPacketPlayerInfoDataInnerShowHat)
			if err != nil {
				return
			}
		default:
			_, ok := ret.Data[iPlayToClientPacketPlayerInfoData].ShowHat.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = ret.Data[iPlayToClientPacketPlayerInfoData].ShowHat.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
	}
	return
}

type PlayToClientPacketPlayerRemove struct {
	Players []uuid.UUID
}

func (_ PlayToClientPacketPlayerRemove) Decode(r io.Reader) (ret PlayToClientPacketPlayerRemove, err error) {
	var lPlayToClientPacketPlayerRemovePlayers queser.VarInt
	lPlayToClientPacketPlayerRemovePlayers, err = lPlayToClientPacketPlayerRemovePlayers.Decode(r)
	if err != nil {
		return
	}
	ret.Players = []uuid.UUID{}
	for range lPlayToClientPacketPlayerRemovePlayers {
		var PlayToClientPacketPlayerRemovePlayersElement uuid.UUID
		_, err = io.ReadFull(r, PlayToClientPacketPlayerRemovePlayersElement[:])
		if err != nil {
			return
		}
		ret.Players = append(ret.Players, PlayToClientPacketPlayerRemovePlayersElement)
	}
	return
}
func (ret PlayToClientPacketPlayerRemove) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Players)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketPlayerRemovePlayers := range len(ret.Players) {
		_, err = w.Write(ret.Players[iPlayToClientPacketPlayerRemovePlayers][:])
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketPlayerRotation struct {
	Yaw   float32
	Pitch float32
}

func (_ PlayToClientPacketPlayerRotation) Decode(r io.Reader) (ret PlayToClientPacketPlayerRotation, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketPlayerRotation) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketPlayerlistHeader struct {
	Header nbt.Anon
	Footer nbt.Anon
}

func (_ PlayToClientPacketPlayerlistHeader) Decode(r io.Reader) (ret PlayToClientPacketPlayerlistHeader, err error) {
	ret.Header, err = ret.Header.Decode(r)
	if err != nil {
		return
	}
	ret.Footer, err = ret.Footer.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketPlayerlistHeader) Encode(w io.Writer) (err error) {
	err = ret.Header.Encode(w)
	if err != nil {
		return
	}
	err = ret.Footer.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketPosition struct {
	TeleportId queser.VarInt
	X          float64
	Y          float64
	Z          float64
	Dx         float64
	Dy         float64
	Dz         float64
	Yaw        float32
	Pitch      float32
	Flags      PlayToClientPositionUpdateRelatives
}

func (_ PlayToClientPacketPosition) Decode(r io.Reader) (ret PlayToClientPacketPosition, err error) {
	ret.TeleportId, err = ret.TeleportId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Dx)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Dy)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Dz)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	ret.Flags, err = ret.Flags.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketPosition) Encode(w io.Writer) (err error) {
	err = ret.TeleportId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Dx)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Dy)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Dz)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	err = ret.Flags.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketProfilelessChat struct {
	Message nbt.Anon
	Type    PlayToClientChatTypesHolder
	Name    nbt.Anon
	Target  *nbt.Anon
}

func (_ PlayToClientPacketProfilelessChat) Decode(r io.Reader) (ret PlayToClientPacketProfilelessChat, err error) {
	ret.Message, err = ret.Message.Decode(r)
	if err != nil {
		return
	}
	ret.Type, err = ret.Type.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = ret.Name.Decode(r)
	if err != nil {
		return
	}
	var PlayToClientPacketProfilelessChatTargetPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketProfilelessChatTargetPresent)
	if err != nil {
		return
	}
	if PlayToClientPacketProfilelessChatTargetPresent {
		var PlayToClientPacketProfilelessChatTargetPresentValue nbt.Anon
		PlayToClientPacketProfilelessChatTargetPresentValue, err = PlayToClientPacketProfilelessChatTargetPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.Target = &PlayToClientPacketProfilelessChatTargetPresentValue
	}
	return
}
func (ret PlayToClientPacketProfilelessChat) Encode(w io.Writer) (err error) {
	err = ret.Message.Encode(w)
	if err != nil {
		return
	}
	err = ret.Type.Encode(w)
	if err != nil {
		return
	}
	err = ret.Name.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Target != nil)
	if err != nil {
		return
	}
	if ret.Target != nil {
		err = (*ret.Target).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketRecipeBookAdd struct {
	Entries []struct {
		Recipe struct {
			DisplayId            queser.VarInt
			Display              PlayToClientRecipeDisplay
			Group                Optvarint
			Category             string
			CraftingRequirements *[]IDSet
		}
		Flags uint8
	}
	Replace bool
}

var PlayToClientPacketRecipeBookAddEntriesElementRecipeCategoryMap = map[queser.VarInt]string{0: "crafting_building_blocks", 1: "crafting_redstone", 10: "stonecutter", 11: "smithing", 12: "campfire", 2: "crafting_equipment", 3: "crafting_misc", 4: "furnace_food", 5: "furnace_blocks", 6: "furnace_misc", 7: "blast_furnace_blocks", 8: "blast_furnace_misc", 9: "smoker_food"}

func (_ PlayToClientPacketRecipeBookAdd) Decode(r io.Reader) (ret PlayToClientPacketRecipeBookAdd, err error) {
	var lPlayToClientPacketRecipeBookAddEntries queser.VarInt
	lPlayToClientPacketRecipeBookAddEntries, err = lPlayToClientPacketRecipeBookAddEntries.Decode(r)
	if err != nil {
		return
	}
	ret.Entries = []struct {
		Recipe struct {
			DisplayId            queser.VarInt
			Display              PlayToClientRecipeDisplay
			Group                Optvarint
			Category             string
			CraftingRequirements *[]IDSet
		}
		Flags uint8
	}{}
	for range lPlayToClientPacketRecipeBookAddEntries {
		var PlayToClientPacketRecipeBookAddEntriesElement struct {
			Recipe struct {
				DisplayId            queser.VarInt
				Display              PlayToClientRecipeDisplay
				Group                Optvarint
				Category             string
				CraftingRequirements *[]IDSet
			}
			Flags uint8
		}
		PlayToClientPacketRecipeBookAddEntriesElement.Recipe.DisplayId, err = PlayToClientPacketRecipeBookAddEntriesElement.Recipe.DisplayId.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketRecipeBookAddEntriesElement.Recipe.Display, err = PlayToClientPacketRecipeBookAddEntriesElement.Recipe.Display.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketRecipeBookAddEntriesElement.Recipe.Group, err = PlayToClientPacketRecipeBookAddEntriesElement.Recipe.Group.Decode(r)
		if err != nil {
			return
		}
		var PlayToClientPacketRecipeBookAddEntriesElementRecipeCategoryKey queser.VarInt
		PlayToClientPacketRecipeBookAddEntriesElementRecipeCategoryKey, err = PlayToClientPacketRecipeBookAddEntriesElementRecipeCategoryKey.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketRecipeBookAddEntriesElement.Recipe.Category, err = queser.ErroringIndex(PlayToClientPacketRecipeBookAddEntriesElementRecipeCategoryMap, PlayToClientPacketRecipeBookAddEntriesElementRecipeCategoryKey)
		if err != nil {
			return
		}
		var PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsPresent bool
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsPresent)
		if err != nil {
			return
		}
		if PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsPresent {
			var PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsPresentValue []IDSet
			var lPlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirements queser.VarInt
			lPlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirements, err = lPlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirements.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsPresentValue = []IDSet{}
			for range lPlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirements {
				var PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsElement IDSet
				PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsElement, err = PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsElement.Decode(r)
				if err != nil {
					return
				}
				PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsPresentValue = append(PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsPresentValue, PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsElement)
			}
			PlayToClientPacketRecipeBookAddEntriesElement.Recipe.CraftingRequirements = &PlayToClientPacketRecipeBookAddEntriesElementRecipeCraftingRequirementsPresentValue
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketRecipeBookAddEntriesElement.Flags)
		if err != nil {
			return
		}
		ret.Entries = append(ret.Entries, PlayToClientPacketRecipeBookAddEntriesElement)
	}
	err = binary.Read(r, binary.BigEndian, &ret.Replace)
	if err != nil {
		return
	}
	return
}

var PlayToClientPacketRecipeBookAddEntriesRecipeCategoryReverseMap = map[string]queser.VarInt{"crafting_building_blocks": 0, "crafting_redstone": 1, "stonecutter": 10, "smithing": 11, "campfire": 12, "crafting_equipment": 2, "crafting_misc": 3, "furnace_food": 4, "furnace_blocks": 5, "furnace_misc": 6, "blast_furnace_blocks": 7, "blast_furnace_misc": 8, "smoker_food": 9}
var PlayToClientPacketRecipeBookAddEntriesInnerRecipeCategoryReverseMap = map[string]queser.VarInt{"crafting_building_blocks": 0, "crafting_redstone": 1, "stonecutter": 10, "smithing": 11, "campfire": 12, "crafting_equipment": 2, "crafting_misc": 3, "furnace_food": 4, "furnace_blocks": 5, "furnace_misc": 6, "blast_furnace_blocks": 7, "blast_furnace_misc": 8, "smoker_food": 9}

func (ret PlayToClientPacketRecipeBookAdd) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Entries)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketRecipeBookAddEntries := range len(ret.Entries) {
		err = ret.Entries[iPlayToClientPacketRecipeBookAddEntries].Recipe.DisplayId.Encode(w)
		if err != nil {
			return
		}
		err = ret.Entries[iPlayToClientPacketRecipeBookAddEntries].Recipe.Display.Encode(w)
		if err != nil {
			return
		}
		err = ret.Entries[iPlayToClientPacketRecipeBookAddEntries].Recipe.Group.Encode(w)
		if err != nil {
			return
		}
		var vPlayToClientPacketRecipeBookAddEntriesInnerRecipeCategory queser.VarInt
		vPlayToClientPacketRecipeBookAddEntriesInnerRecipeCategory, err = queser.ErroringIndex(PlayToClientPacketRecipeBookAddEntriesInnerRecipeCategoryReverseMap, ret.Entries[iPlayToClientPacketRecipeBookAddEntries].Recipe.Category)
		if err != nil {
			return
		}
		err = vPlayToClientPacketRecipeBookAddEntriesInnerRecipeCategory.Encode(w)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Entries[iPlayToClientPacketRecipeBookAddEntries].Recipe.CraftingRequirements != nil)
		if err != nil {
			return
		}
		if ret.Entries[iPlayToClientPacketRecipeBookAddEntries].Recipe.CraftingRequirements != nil {
			err = queser.VarInt(len(*ret.Entries[iPlayToClientPacketRecipeBookAddEntries].Recipe.CraftingRequirements)).Encode(w)
			if err != nil {
				return
			}
			for iPlayToClientPacketRecipeBookAddEntriesInnerRecipeCraftingRequirements := range len(*ret.Entries[iPlayToClientPacketRecipeBookAddEntries].Recipe.CraftingRequirements) {
				err = (*ret.Entries[iPlayToClientPacketRecipeBookAddEntries].Recipe.CraftingRequirements)[iPlayToClientPacketRecipeBookAddEntriesInnerRecipeCraftingRequirements].Encode(w)
				if err != nil {
					return
				}
			}
		}
		err = binary.Write(w, binary.BigEndian, ret.Entries[iPlayToClientPacketRecipeBookAddEntries].Flags)
		if err != nil {
			return
		}
	}
	err = binary.Write(w, binary.BigEndian, ret.Replace)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketRecipeBookRemove struct {
	RecipeIds []queser.VarInt
}

func (_ PlayToClientPacketRecipeBookRemove) Decode(r io.Reader) (ret PlayToClientPacketRecipeBookRemove, err error) {
	var lPlayToClientPacketRecipeBookRemoveRecipeIds queser.VarInt
	lPlayToClientPacketRecipeBookRemoveRecipeIds, err = lPlayToClientPacketRecipeBookRemoveRecipeIds.Decode(r)
	if err != nil {
		return
	}
	ret.RecipeIds = []queser.VarInt{}
	for range lPlayToClientPacketRecipeBookRemoveRecipeIds {
		var PlayToClientPacketRecipeBookRemoveRecipeIdsElement queser.VarInt
		PlayToClientPacketRecipeBookRemoveRecipeIdsElement, err = PlayToClientPacketRecipeBookRemoveRecipeIdsElement.Decode(r)
		if err != nil {
			return
		}
		ret.RecipeIds = append(ret.RecipeIds, PlayToClientPacketRecipeBookRemoveRecipeIdsElement)
	}
	return
}
func (ret PlayToClientPacketRecipeBookRemove) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.RecipeIds)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketRecipeBookRemoveRecipeIds := range len(ret.RecipeIds) {
		err = ret.RecipeIds[iPlayToClientPacketRecipeBookRemoveRecipeIds].Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketRecipeBookSettings struct {
	Crafting PlayToClientRecipeBookSetting
	Furnace  PlayToClientRecipeBookSetting
	Blast    PlayToClientRecipeBookSetting
	Smoker   PlayToClientRecipeBookSetting
}

func (_ PlayToClientPacketRecipeBookSettings) Decode(r io.Reader) (ret PlayToClientPacketRecipeBookSettings, err error) {
	ret.Crafting, err = ret.Crafting.Decode(r)
	if err != nil {
		return
	}
	ret.Furnace, err = ret.Furnace.Decode(r)
	if err != nil {
		return
	}
	ret.Blast, err = ret.Blast.Decode(r)
	if err != nil {
		return
	}
	ret.Smoker, err = ret.Smoker.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketRecipeBookSettings) Encode(w io.Writer) (err error) {
	err = ret.Crafting.Encode(w)
	if err != nil {
		return
	}
	err = ret.Furnace.Encode(w)
	if err != nil {
		return
	}
	err = ret.Blast.Encode(w)
	if err != nil {
		return
	}
	err = ret.Smoker.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketRelEntityMove struct {
	EntityId queser.VarInt
	DX       int16
	DY       int16
	DZ       int16
	OnGround bool
}

func (_ PlayToClientPacketRelEntityMove) Decode(r io.Reader) (ret PlayToClientPacketRelEntityMove, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.DX)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.DY)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.DZ)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OnGround)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketRelEntityMove) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.DX)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.DY)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.DZ)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OnGround)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketRemoveEntityEffect struct {
	EntityId queser.VarInt
	EffectId queser.VarInt
}

func (_ PlayToClientPacketRemoveEntityEffect) Decode(r io.Reader) (ret PlayToClientPacketRemoveEntityEffect, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	ret.EffectId, err = ret.EffectId.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketRemoveEntityEffect) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = ret.EffectId.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketResetScore struct {
	EntityName    string
	ObjectiveName *string
}

func (_ PlayToClientPacketResetScore) Decode(r io.Reader) (ret PlayToClientPacketResetScore, err error) {
	ret.EntityName, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var PlayToClientPacketResetScoreObjectiveNamePresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketResetScoreObjectiveNamePresent)
	if err != nil {
		return
	}
	if PlayToClientPacketResetScoreObjectiveNamePresent {
		var PlayToClientPacketResetScoreObjectiveNamePresentValue string
		PlayToClientPacketResetScoreObjectiveNamePresentValue, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.ObjectiveName = &PlayToClientPacketResetScoreObjectiveNamePresentValue
	}
	return
}
func (ret PlayToClientPacketResetScore) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.EntityName)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.ObjectiveName != nil)
	if err != nil {
		return
	}
	if ret.ObjectiveName != nil {
		err = queser.EncodeString(w, *ret.ObjectiveName)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketRespawn struct {
	WorldState   PlayToClientSpawnInfo
	CopyMetadata uint8
}

func (_ PlayToClientPacketRespawn) Decode(r io.Reader) (ret PlayToClientPacketRespawn, err error) {
	ret.WorldState, err = ret.WorldState.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.CopyMetadata)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketRespawn) Encode(w io.Writer) (err error) {
	err = ret.WorldState.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.CopyMetadata)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketScoreboardDisplayObjective struct {
	Position queser.VarInt
	Name     string
}

func (_ PlayToClientPacketScoreboardDisplayObjective) Decode(r io.Reader) (ret PlayToClientPacketScoreboardDisplayObjective, err error) {
	ret.Position, err = ret.Position.Decode(r)
	if err != nil {
		return
	}
	ret.Name, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketScoreboardDisplayObjective) Encode(w io.Writer) (err error) {
	err = ret.Position.Encode(w)
	if err != nil {
		return
	}
	err = queser.EncodeString(w, ret.Name)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketScoreboardObjective struct {
	Name         string
	Action       int8
	DisplayText  any
	Type         any
	NumberFormat any
	Styling      any
}

func (_ PlayToClientPacketScoreboardObjective) Decode(r io.Reader) (ret PlayToClientPacketScoreboardObjective, err error) {
	ret.Name, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Action)
	if err != nil {
		return
	}
	switch ret.Action {
	case 0:
		var PlayToClientPacketScoreboardObjectiveDisplayTextTmp nbt.Anon
		PlayToClientPacketScoreboardObjectiveDisplayTextTmp, err = PlayToClientPacketScoreboardObjectiveDisplayTextTmp.Decode(r)
		if err != nil {
			return
		}
		ret.DisplayText = PlayToClientPacketScoreboardObjectiveDisplayTextTmp
	case 2:
		var PlayToClientPacketScoreboardObjectiveDisplayTextTmp nbt.Anon
		PlayToClientPacketScoreboardObjectiveDisplayTextTmp, err = PlayToClientPacketScoreboardObjectiveDisplayTextTmp.Decode(r)
		if err != nil {
			return
		}
		ret.DisplayText = PlayToClientPacketScoreboardObjectiveDisplayTextTmp
	default:
		var PlayToClientPacketScoreboardObjectiveDisplayTextTmp queser.Void
		PlayToClientPacketScoreboardObjectiveDisplayTextTmp, err = PlayToClientPacketScoreboardObjectiveDisplayTextTmp.Decode(r)
		if err != nil {
			return
		}
		ret.DisplayText = PlayToClientPacketScoreboardObjectiveDisplayTextTmp
	}
	switch ret.Action {
	case 0:
		var PlayToClientPacketScoreboardObjectiveTypeTmp queser.VarInt
		PlayToClientPacketScoreboardObjectiveTypeTmp, err = PlayToClientPacketScoreboardObjectiveTypeTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Type = PlayToClientPacketScoreboardObjectiveTypeTmp
	case 2:
		var PlayToClientPacketScoreboardObjectiveTypeTmp queser.VarInt
		PlayToClientPacketScoreboardObjectiveTypeTmp, err = PlayToClientPacketScoreboardObjectiveTypeTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Type = PlayToClientPacketScoreboardObjectiveTypeTmp
	default:
		var PlayToClientPacketScoreboardObjectiveTypeTmp queser.Void
		PlayToClientPacketScoreboardObjectiveTypeTmp, err = PlayToClientPacketScoreboardObjectiveTypeTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Type = PlayToClientPacketScoreboardObjectiveTypeTmp
	}
	switch ret.Action {
	case 0:
		var PlayToClientPacketScoreboardObjectiveNumberFormatTmp *queser.VarInt
		var PlayToClientPacketScoreboardObjectiveNumberFormatPresent bool
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketScoreboardObjectiveNumberFormatPresent)
		if err != nil {
			return
		}
		if PlayToClientPacketScoreboardObjectiveNumberFormatPresent {
			var PlayToClientPacketScoreboardObjectiveNumberFormatPresentValue queser.VarInt
			PlayToClientPacketScoreboardObjectiveNumberFormatPresentValue, err = PlayToClientPacketScoreboardObjectiveNumberFormatPresentValue.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketScoreboardObjectiveNumberFormatTmp = &PlayToClientPacketScoreboardObjectiveNumberFormatPresentValue
		}
		ret.NumberFormat = PlayToClientPacketScoreboardObjectiveNumberFormatTmp
	case 2:
		var PlayToClientPacketScoreboardObjectiveNumberFormatTmp *queser.VarInt
		var PlayToClientPacketScoreboardObjectiveNumberFormatPresent bool
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketScoreboardObjectiveNumberFormatPresent)
		if err != nil {
			return
		}
		if PlayToClientPacketScoreboardObjectiveNumberFormatPresent {
			var PlayToClientPacketScoreboardObjectiveNumberFormatPresentValue queser.VarInt
			PlayToClientPacketScoreboardObjectiveNumberFormatPresentValue, err = PlayToClientPacketScoreboardObjectiveNumberFormatPresentValue.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketScoreboardObjectiveNumberFormatTmp = &PlayToClientPacketScoreboardObjectiveNumberFormatPresentValue
		}
		ret.NumberFormat = PlayToClientPacketScoreboardObjectiveNumberFormatTmp
	default:
		var PlayToClientPacketScoreboardObjectiveNumberFormatTmp queser.Void
		PlayToClientPacketScoreboardObjectiveNumberFormatTmp, err = PlayToClientPacketScoreboardObjectiveNumberFormatTmp.Decode(r)
		if err != nil {
			return
		}
		ret.NumberFormat = PlayToClientPacketScoreboardObjectiveNumberFormatTmp
	}
	switch ret.Action {
	case 0:
		var PlayToClientPacketScoreboardObjectiveStylingTmp any
		switch ret.NumberFormat {
		case 1:
			var PlayToClientPacketScoreboardObjectiveStylingTmp nbt.Anon
			PlayToClientPacketScoreboardObjectiveStylingTmp, err = PlayToClientPacketScoreboardObjectiveStylingTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketScoreboardObjectiveStylingTmp = PlayToClientPacketScoreboardObjectiveStylingTmp
		case 2:
			var PlayToClientPacketScoreboardObjectiveStylingTmp nbt.Anon
			PlayToClientPacketScoreboardObjectiveStylingTmp, err = PlayToClientPacketScoreboardObjectiveStylingTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketScoreboardObjectiveStylingTmp = PlayToClientPacketScoreboardObjectiveStylingTmp
		default:
			var PlayToClientPacketScoreboardObjectiveStylingTmp queser.Void
			PlayToClientPacketScoreboardObjectiveStylingTmp, err = PlayToClientPacketScoreboardObjectiveStylingTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketScoreboardObjectiveStylingTmp = PlayToClientPacketScoreboardObjectiveStylingTmp
		}
		ret.Styling = PlayToClientPacketScoreboardObjectiveStylingTmp
	case 2:
		var PlayToClientPacketScoreboardObjectiveStylingTmp any
		switch ret.NumberFormat {
		case 1:
			var PlayToClientPacketScoreboardObjectiveStylingTmp nbt.Anon
			PlayToClientPacketScoreboardObjectiveStylingTmp, err = PlayToClientPacketScoreboardObjectiveStylingTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketScoreboardObjectiveStylingTmp = PlayToClientPacketScoreboardObjectiveStylingTmp
		case 2:
			var PlayToClientPacketScoreboardObjectiveStylingTmp nbt.Anon
			PlayToClientPacketScoreboardObjectiveStylingTmp, err = PlayToClientPacketScoreboardObjectiveStylingTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketScoreboardObjectiveStylingTmp = PlayToClientPacketScoreboardObjectiveStylingTmp
		default:
			var PlayToClientPacketScoreboardObjectiveStylingTmp queser.Void
			PlayToClientPacketScoreboardObjectiveStylingTmp, err = PlayToClientPacketScoreboardObjectiveStylingTmp.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketScoreboardObjectiveStylingTmp = PlayToClientPacketScoreboardObjectiveStylingTmp
		}
		ret.Styling = PlayToClientPacketScoreboardObjectiveStylingTmp
	default:
		var PlayToClientPacketScoreboardObjectiveStylingTmp queser.Void
		PlayToClientPacketScoreboardObjectiveStylingTmp, err = PlayToClientPacketScoreboardObjectiveStylingTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Styling = PlayToClientPacketScoreboardObjectiveStylingTmp
	}
	return
}
func (ret PlayToClientPacketScoreboardObjective) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Name)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Action)
	if err != nil {
		return
	}
	switch ret.Action {
	case 0:
		PlayToClientPacketScoreboardObjectiveDisplayText, ok := ret.DisplayText.(nbt.Anon)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketScoreboardObjectiveDisplayText.Encode(w)
		if err != nil {
			return
		}
	case 2:
		PlayToClientPacketScoreboardObjectiveDisplayText, ok := ret.DisplayText.(nbt.Anon)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketScoreboardObjectiveDisplayText.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.DisplayText.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.DisplayText.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Action {
	case 0:
		PlayToClientPacketScoreboardObjectiveType, ok := ret.Type.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketScoreboardObjectiveType.Encode(w)
		if err != nil {
			return
		}
	case 2:
		PlayToClientPacketScoreboardObjectiveType, ok := ret.Type.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketScoreboardObjectiveType.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Type.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Type.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Action {
	case 0:
		PlayToClientPacketScoreboardObjectiveNumberFormat, ok := ret.NumberFormat.(*queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToClientPacketScoreboardObjectiveNumberFormat != nil)
		if err != nil {
			return
		}
		if PlayToClientPacketScoreboardObjectiveNumberFormat != nil {
			err = (*PlayToClientPacketScoreboardObjectiveNumberFormat).Encode(w)
			if err != nil {
				return
			}
		}
	case 2:
		PlayToClientPacketScoreboardObjectiveNumberFormat, ok := ret.NumberFormat.(*queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToClientPacketScoreboardObjectiveNumberFormat != nil)
		if err != nil {
			return
		}
		if PlayToClientPacketScoreboardObjectiveNumberFormat != nil {
			err = (*PlayToClientPacketScoreboardObjectiveNumberFormat).Encode(w)
			if err != nil {
				return
			}
		}
	default:
		_, ok := ret.NumberFormat.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.NumberFormat.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Action {
	case 0:
		PlayToClientPacketScoreboardObjectiveStyling, ok := ret.Styling.(any)
		if !ok {
			err = queser.BadTypeError
			return
		}
		switch ret.NumberFormat {
		case 1:
			PlayToClientPacketScoreboardObjectiveStyling, ok := PlayToClientPacketScoreboardObjectiveStyling.(nbt.Anon)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketScoreboardObjectiveStyling.Encode(w)
			if err != nil {
				return
			}
		case 2:
			PlayToClientPacketScoreboardObjectiveStyling, ok := PlayToClientPacketScoreboardObjectiveStyling.(nbt.Anon)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketScoreboardObjectiveStyling.Encode(w)
			if err != nil {
				return
			}
		default:
			_, ok := PlayToClientPacketScoreboardObjectiveStyling.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketScoreboardObjectiveStyling.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
	case 2:
		PlayToClientPacketScoreboardObjectiveStyling, ok := ret.Styling.(any)
		if !ok {
			err = queser.BadTypeError
			return
		}
		switch ret.NumberFormat {
		case 1:
			PlayToClientPacketScoreboardObjectiveStyling, ok := PlayToClientPacketScoreboardObjectiveStyling.(nbt.Anon)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketScoreboardObjectiveStyling.Encode(w)
			if err != nil {
				return
			}
		case 2:
			PlayToClientPacketScoreboardObjectiveStyling, ok := PlayToClientPacketScoreboardObjectiveStyling.(nbt.Anon)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketScoreboardObjectiveStyling.Encode(w)
			if err != nil {
				return
			}
		default:
			_, ok := PlayToClientPacketScoreboardObjectiveStyling.(queser.Void)
			if !ok {
				err = queser.BadTypeError
				return
			}
			err = PlayToClientPacketScoreboardObjectiveStyling.(queser.Void).Encode(w)
			if err != nil {
				return
			}
		}
	default:
		_, ok := ret.Styling.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Styling.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketScoreboardScore struct {
	Val queser.ToDo
}

func (_ PlayToClientPacketScoreboardScore) Decode(r io.Reader) (ret PlayToClientPacketScoreboardScore, err error) {
	err = queser.ToDoError
	return
}
func (ret PlayToClientPacketScoreboardScore) Encode(w io.Writer) (err error) {
	err = queser.ToDoError
	return
}

type PlayToClientPacketSelectAdvancementTab struct {
	Id *string
}

func (_ PlayToClientPacketSelectAdvancementTab) Decode(r io.Reader) (ret PlayToClientPacketSelectAdvancementTab, err error) {
	var PlayToClientPacketSelectAdvancementTabIdPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketSelectAdvancementTabIdPresent)
	if err != nil {
		return
	}
	if PlayToClientPacketSelectAdvancementTabIdPresent {
		var PlayToClientPacketSelectAdvancementTabIdPresentValue string
		PlayToClientPacketSelectAdvancementTabIdPresentValue, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Id = &PlayToClientPacketSelectAdvancementTabIdPresentValue
	}
	return
}
func (ret PlayToClientPacketSelectAdvancementTab) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Id != nil)
	if err != nil {
		return
	}
	if ret.Id != nil {
		err = queser.EncodeString(w, *ret.Id)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketServerData struct {
	Motd      nbt.Anon
	IconBytes *ByteArray
}

func (_ PlayToClientPacketServerData) Decode(r io.Reader) (ret PlayToClientPacketServerData, err error) {
	ret.Motd, err = ret.Motd.Decode(r)
	if err != nil {
		return
	}
	var PlayToClientPacketServerDataIconBytesPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketServerDataIconBytesPresent)
	if err != nil {
		return
	}
	if PlayToClientPacketServerDataIconBytesPresent {
		var PlayToClientPacketServerDataIconBytesPresentValue ByteArray
		PlayToClientPacketServerDataIconBytesPresentValue, err = PlayToClientPacketServerDataIconBytesPresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.IconBytes = &PlayToClientPacketServerDataIconBytesPresentValue
	}
	return
}
func (ret PlayToClientPacketServerData) Encode(w io.Writer) (err error) {
	err = ret.Motd.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IconBytes != nil)
	if err != nil {
		return
	}
	if ret.IconBytes != nil {
		err = (*ret.IconBytes).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketSetCooldown struct {
	CooldownGroup string
	CooldownTicks queser.VarInt
}

func (_ PlayToClientPacketSetCooldown) Decode(r io.Reader) (ret PlayToClientPacketSetCooldown, err error) {
	ret.CooldownGroup, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	ret.CooldownTicks, err = ret.CooldownTicks.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSetCooldown) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.CooldownGroup)
	if err != nil {
		return
	}
	err = ret.CooldownTicks.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSetCursorItem struct {
	Contents Slot
}

func (_ PlayToClientPacketSetCursorItem) Decode(r io.Reader) (ret PlayToClientPacketSetCursorItem, err error) {
	ret.Contents, err = ret.Contents.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSetCursorItem) Encode(w io.Writer) (err error) {
	err = ret.Contents.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSetPassengers struct {
	EntityId   queser.VarInt
	Passengers []queser.VarInt
}

func (_ PlayToClientPacketSetPassengers) Decode(r io.Reader) (ret PlayToClientPacketSetPassengers, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	var lPlayToClientPacketSetPassengersPassengers queser.VarInt
	lPlayToClientPacketSetPassengersPassengers, err = lPlayToClientPacketSetPassengersPassengers.Decode(r)
	if err != nil {
		return
	}
	ret.Passengers = []queser.VarInt{}
	for range lPlayToClientPacketSetPassengersPassengers {
		var PlayToClientPacketSetPassengersPassengersElement queser.VarInt
		PlayToClientPacketSetPassengersPassengersElement, err = PlayToClientPacketSetPassengersPassengersElement.Decode(r)
		if err != nil {
			return
		}
		ret.Passengers = append(ret.Passengers, PlayToClientPacketSetPassengersPassengersElement)
	}
	return
}
func (ret PlayToClientPacketSetPassengers) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Passengers)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketSetPassengersPassengers := range len(ret.Passengers) {
		err = ret.Passengers[iPlayToClientPacketSetPassengersPassengers].Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketSetPlayerInventory struct {
	SlotId   queser.VarInt
	Contents Slot
}

func (_ PlayToClientPacketSetPlayerInventory) Decode(r io.Reader) (ret PlayToClientPacketSetPlayerInventory, err error) {
	ret.SlotId, err = ret.SlotId.Decode(r)
	if err != nil {
		return
	}
	ret.Contents, err = ret.Contents.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSetPlayerInventory) Encode(w io.Writer) (err error) {
	err = ret.SlotId.Encode(w)
	if err != nil {
		return
	}
	err = ret.Contents.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSetProjectilePower struct {
	Id                queser.VarInt
	AccelerationPower float64
}

func (_ PlayToClientPacketSetProjectilePower) Decode(r io.Reader) (ret PlayToClientPacketSetProjectilePower, err error) {
	ret.Id, err = ret.Id.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.AccelerationPower)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSetProjectilePower) Encode(w io.Writer) (err error) {
	err = ret.Id.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.AccelerationPower)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSetSlot struct {
	WindowId ContainerID
	StateId  queser.VarInt
	Slot     int16
	Item     Slot
}

func (_ PlayToClientPacketSetSlot) Decode(r io.Reader) (ret PlayToClientPacketSetSlot, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	ret.StateId, err = ret.StateId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Slot)
	if err != nil {
		return
	}
	ret.Item, err = ret.Item.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSetSlot) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	err = ret.StateId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Slot)
	if err != nil {
		return
	}
	err = ret.Item.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSetTickingState struct {
	TickRate float32
	IsFrozen bool
}

func (_ PlayToClientPacketSetTickingState) Decode(r io.Reader) (ret PlayToClientPacketSetTickingState, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.TickRate)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IsFrozen)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSetTickingState) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.TickRate)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IsFrozen)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSetTitleSubtitle struct {
	Text nbt.Anon
}

func (_ PlayToClientPacketSetTitleSubtitle) Decode(r io.Reader) (ret PlayToClientPacketSetTitleSubtitle, err error) {
	ret.Text, err = ret.Text.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSetTitleSubtitle) Encode(w io.Writer) (err error) {
	err = ret.Text.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSetTitleText struct {
	Text nbt.Anon
}

func (_ PlayToClientPacketSetTitleText) Decode(r io.Reader) (ret PlayToClientPacketSetTitleText, err error) {
	ret.Text, err = ret.Text.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSetTitleText) Encode(w io.Writer) (err error) {
	err = ret.Text.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSetTitleTime struct {
	FadeIn  int32
	Stay    int32
	FadeOut int32
}

func (_ PlayToClientPacketSetTitleTime) Decode(r io.Reader) (ret PlayToClientPacketSetTitleTime, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.FadeIn)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Stay)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.FadeOut)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSetTitleTime) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.FadeIn)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Stay)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.FadeOut)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketShowDialog struct {
	Dialog any
}

func (_ PlayToClientPacketShowDialog) Decode(r io.Reader) (ret PlayToClientPacketShowDialog, err error) {
	var PlayToClientPacketShowDialogDialogId queser.VarInt
	PlayToClientPacketShowDialogDialogId, err = PlayToClientPacketShowDialogDialogId.Decode(r)
	if err != nil {
		return
	}
	if PlayToClientPacketShowDialogDialogId != 0 {
		ret.Dialog = PlayToClientPacketShowDialogDialogId
		return
	}
	var PlayToClientPacketShowDialogDialogResult nbt.Anon
	PlayToClientPacketShowDialogDialogResult, err = PlayToClientPacketShowDialogDialogResult.Decode(r)
	if err != nil {
		return
	}
	ret.Dialog = PlayToClientPacketShowDialogDialogResult
	return
}
func (ret PlayToClientPacketShowDialog) Encode(w io.Writer) (err error) {
	switch PlayToClientPacketShowDialogDialogKnownType := ret.Dialog.(type) {
	case queser.VarInt:
		err = PlayToClientPacketShowDialogDialogKnownType.Encode(w)
		if err != nil {
			return
		}
	case nbt.Anon:
		err = PlayToClientPacketShowDialogDialogKnownType.Encode(w)
		if err != nil {
			return
		}
	default:
		err = queser.BadTypeError
	}
	return
}

type PlayToClientPacketSimulationDistance struct {
	Distance queser.VarInt
}

func (_ PlayToClientPacketSimulationDistance) Decode(r io.Reader) (ret PlayToClientPacketSimulationDistance, err error) {
	ret.Distance, err = ret.Distance.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSimulationDistance) Encode(w io.Writer) (err error) {
	err = ret.Distance.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSoundEffect struct {
	Sound         ItemSoundHolder
	SoundCategory SoundSource
	X             int32
	Y             int32
	Z             int32
	Volume        float32
	Pitch         float32
	Seed          int64
}

func (_ PlayToClientPacketSoundEffect) Decode(r io.Reader) (ret PlayToClientPacketSoundEffect, err error) {
	ret.Sound, err = ret.Sound.Decode(r)
	if err != nil {
		return
	}
	ret.SoundCategory, err = ret.SoundCategory.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Volume)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Seed)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSoundEffect) Encode(w io.Writer) (err error) {
	err = ret.Sound.Encode(w)
	if err != nil {
		return
	}
	err = ret.SoundCategory.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Volume)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Seed)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSpawnEntity struct {
	EntityId   queser.VarInt
	ObjectUUID uuid.UUID
	Type       queser.VarInt
	X          float64
	Y          float64
	Z          float64
	Pitch      int8
	Yaw        int8
	HeadPitch  int8
	ObjectData queser.VarInt
	VelocityX  int16
	VelocityY  int16
	VelocityZ  int16
}

func (_ PlayToClientPacketSpawnEntity) Decode(r io.Reader) (ret PlayToClientPacketSpawnEntity, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	_, err = io.ReadFull(r, ret.ObjectUUID[:])
	if err != nil {
		return
	}
	ret.Type, err = ret.Type.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.HeadPitch)
	if err != nil {
		return
	}
	ret.ObjectData, err = ret.ObjectData.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.VelocityX)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.VelocityY)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.VelocityZ)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSpawnEntity) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	_, err = w.Write(ret.ObjectUUID[:])
	if err != nil {
		return
	}
	err = ret.Type.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.HeadPitch)
	if err != nil {
		return
	}
	err = ret.ObjectData.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.VelocityX)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.VelocityY)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.VelocityZ)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSpawnPosition struct {
	Location Position
	Angle    float32
}

func (_ PlayToClientPacketSpawnPosition) Decode(r io.Reader) (ret PlayToClientPacketSpawnPosition, err error) {
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Angle)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSpawnPosition) Encode(w io.Writer) (err error) {
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Angle)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketStartConfiguration struct {
}

func (_ PlayToClientPacketStartConfiguration) Decode(r io.Reader) (ret PlayToClientPacketStartConfiguration, err error) {
	return
}
func (ret PlayToClientPacketStartConfiguration) Encode(w io.Writer) (err error) {
	return
}

type PlayToClientPacketStatistics struct {
	Entries []struct {
		CategoryId  queser.VarInt
		StatisticId queser.VarInt
		Value       queser.VarInt
	}
}

func (_ PlayToClientPacketStatistics) Decode(r io.Reader) (ret PlayToClientPacketStatistics, err error) {
	var lPlayToClientPacketStatisticsEntries queser.VarInt
	lPlayToClientPacketStatisticsEntries, err = lPlayToClientPacketStatisticsEntries.Decode(r)
	if err != nil {
		return
	}
	ret.Entries = []struct {
		CategoryId  queser.VarInt
		StatisticId queser.VarInt
		Value       queser.VarInt
	}{}
	for range lPlayToClientPacketStatisticsEntries {
		var PlayToClientPacketStatisticsEntriesElement struct {
			CategoryId  queser.VarInt
			StatisticId queser.VarInt
			Value       queser.VarInt
		}
		PlayToClientPacketStatisticsEntriesElement.CategoryId, err = PlayToClientPacketStatisticsEntriesElement.CategoryId.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketStatisticsEntriesElement.StatisticId, err = PlayToClientPacketStatisticsEntriesElement.StatisticId.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketStatisticsEntriesElement.Value, err = PlayToClientPacketStatisticsEntriesElement.Value.Decode(r)
		if err != nil {
			return
		}
		ret.Entries = append(ret.Entries, PlayToClientPacketStatisticsEntriesElement)
	}
	return
}
func (ret PlayToClientPacketStatistics) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Entries)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketStatisticsEntries := range len(ret.Entries) {
		err = ret.Entries[iPlayToClientPacketStatisticsEntries].CategoryId.Encode(w)
		if err != nil {
			return
		}
		err = ret.Entries[iPlayToClientPacketStatisticsEntries].StatisticId.Encode(w)
		if err != nil {
			return
		}
		err = ret.Entries[iPlayToClientPacketStatisticsEntries].Value.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketStepTick struct {
	TickSteps queser.VarInt
}

func (_ PlayToClientPacketStepTick) Decode(r io.Reader) (ret PlayToClientPacketStepTick, err error) {
	ret.TickSteps, err = ret.TickSteps.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketStepTick) Encode(w io.Writer) (err error) {
	err = ret.TickSteps.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketStopSound struct {
	Flags  int8
	Source any
	Sound  any
}

func (_ PlayToClientPacketStopSound) Decode(r io.Reader) (ret PlayToClientPacketStopSound, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Flags)
	if err != nil {
		return
	}
	switch ret.Flags {
	case 1:
		var PlayToClientPacketStopSoundSourceTmp queser.VarInt
		PlayToClientPacketStopSoundSourceTmp, err = PlayToClientPacketStopSoundSourceTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Source = PlayToClientPacketStopSoundSourceTmp
	case 3:
		var PlayToClientPacketStopSoundSourceTmp queser.VarInt
		PlayToClientPacketStopSoundSourceTmp, err = PlayToClientPacketStopSoundSourceTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Source = PlayToClientPacketStopSoundSourceTmp
	default:
		var PlayToClientPacketStopSoundSourceTmp queser.Void
		PlayToClientPacketStopSoundSourceTmp, err = PlayToClientPacketStopSoundSourceTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Source = PlayToClientPacketStopSoundSourceTmp
	}
	switch ret.Flags {
	case 2:
		var PlayToClientPacketStopSoundSoundTmp string
		PlayToClientPacketStopSoundSoundTmp, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Sound = PlayToClientPacketStopSoundSoundTmp
	case 3:
		var PlayToClientPacketStopSoundSoundTmp string
		PlayToClientPacketStopSoundSoundTmp, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Sound = PlayToClientPacketStopSoundSoundTmp
	default:
		var PlayToClientPacketStopSoundSoundTmp queser.Void
		PlayToClientPacketStopSoundSoundTmp, err = PlayToClientPacketStopSoundSoundTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Sound = PlayToClientPacketStopSoundSoundTmp
	}
	return
}
func (ret PlayToClientPacketStopSound) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Flags)
	if err != nil {
		return
	}
	switch ret.Flags {
	case 1:
		PlayToClientPacketStopSoundSource, ok := ret.Source.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketStopSoundSource.Encode(w)
		if err != nil {
			return
		}
	case 3:
		PlayToClientPacketStopSoundSource, ok := ret.Source.(queser.VarInt)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketStopSoundSource.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Source.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Source.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Flags {
	case 2:
		PlayToClientPacketStopSoundSound, ok := ret.Sound.(string)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.EncodeString(w, PlayToClientPacketStopSoundSound)
		if err != nil {
			return
		}
	case 3:
		PlayToClientPacketStopSoundSound, ok := ret.Sound.(string)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.EncodeString(w, PlayToClientPacketStopSoundSound)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Sound.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Sound.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketSyncEntityPosition struct {
	EntityId queser.VarInt
	X        float64
	Y        float64
	Z        float64
	Dx       float64
	Dy       float64
	Dz       float64
	Yaw      float32
	Pitch    float32
	OnGround bool
}

func (_ PlayToClientPacketSyncEntityPosition) Decode(r io.Reader) (ret PlayToClientPacketSyncEntityPosition, err error) {
	ret.EntityId, err = ret.EntityId.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Dx)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Dy)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Dz)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OnGround)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSyncEntityPosition) Encode(w io.Writer) (err error) {
	err = ret.EntityId.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Dx)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Dy)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Dz)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OnGround)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketSystemChat struct {
	Content     nbt.Anon
	IsActionBar bool
}

func (_ PlayToClientPacketSystemChat) Decode(r io.Reader) (ret PlayToClientPacketSystemChat, err error) {
	ret.Content, err = ret.Content.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IsActionBar)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketSystemChat) Encode(w io.Writer) (err error) {
	err = ret.Content.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IsActionBar)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketTabComplete struct {
	TransactionId queser.VarInt
	Start         queser.VarInt
	Length        queser.VarInt
	Matches       []struct {
		Match   string
		Tooltip *nbt.Anon
	}
}

func (_ PlayToClientPacketTabComplete) Decode(r io.Reader) (ret PlayToClientPacketTabComplete, err error) {
	ret.TransactionId, err = ret.TransactionId.Decode(r)
	if err != nil {
		return
	}
	ret.Start, err = ret.Start.Decode(r)
	if err != nil {
		return
	}
	ret.Length, err = ret.Length.Decode(r)
	if err != nil {
		return
	}
	var lPlayToClientPacketTabCompleteMatches queser.VarInt
	lPlayToClientPacketTabCompleteMatches, err = lPlayToClientPacketTabCompleteMatches.Decode(r)
	if err != nil {
		return
	}
	ret.Matches = []struct {
		Match   string
		Tooltip *nbt.Anon
	}{}
	for range lPlayToClientPacketTabCompleteMatches {
		var PlayToClientPacketTabCompleteMatchesElement struct {
			Match   string
			Tooltip *nbt.Anon
		}
		PlayToClientPacketTabCompleteMatchesElement.Match, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		var PlayToClientPacketTabCompleteMatchesElementTooltipPresent bool
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTabCompleteMatchesElementTooltipPresent)
		if err != nil {
			return
		}
		if PlayToClientPacketTabCompleteMatchesElementTooltipPresent {
			var PlayToClientPacketTabCompleteMatchesElementTooltipPresentValue nbt.Anon
			PlayToClientPacketTabCompleteMatchesElementTooltipPresentValue, err = PlayToClientPacketTabCompleteMatchesElementTooltipPresentValue.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketTabCompleteMatchesElement.Tooltip = &PlayToClientPacketTabCompleteMatchesElementTooltipPresentValue
		}
		ret.Matches = append(ret.Matches, PlayToClientPacketTabCompleteMatchesElement)
	}
	return
}
func (ret PlayToClientPacketTabComplete) Encode(w io.Writer) (err error) {
	err = ret.TransactionId.Encode(w)
	if err != nil {
		return
	}
	err = ret.Start.Encode(w)
	if err != nil {
		return
	}
	err = ret.Length.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Matches)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketTabCompleteMatches := range len(ret.Matches) {
		err = queser.EncodeString(w, ret.Matches[iPlayToClientPacketTabCompleteMatches].Match)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Matches[iPlayToClientPacketTabCompleteMatches].Tooltip != nil)
		if err != nil {
			return
		}
		if ret.Matches[iPlayToClientPacketTabCompleteMatches].Tooltip != nil {
			err = (*ret.Matches[iPlayToClientPacketTabCompleteMatches].Tooltip).Encode(w)
			if err != nil {
				return
			}
		}
	}
	return
}

type PlayToClientPacketTags struct {
	Tags []struct {
		TagType string
		Tags    Tags
	}
}

func (_ PlayToClientPacketTags) Decode(r io.Reader) (ret PlayToClientPacketTags, err error) {
	var lPlayToClientPacketTagsTags queser.VarInt
	lPlayToClientPacketTagsTags, err = lPlayToClientPacketTagsTags.Decode(r)
	if err != nil {
		return
	}
	ret.Tags = []struct {
		TagType string
		Tags    Tags
	}{}
	for range lPlayToClientPacketTagsTags {
		var PlayToClientPacketTagsTagsElement struct {
			TagType string
			Tags    Tags
		}
		PlayToClientPacketTagsTagsElement.TagType, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		PlayToClientPacketTagsTagsElement.Tags, err = PlayToClientPacketTagsTagsElement.Tags.Decode(r)
		if err != nil {
			return
		}
		ret.Tags = append(ret.Tags, PlayToClientPacketTagsTagsElement)
	}
	return
}
func (ret PlayToClientPacketTags) Encode(w io.Writer) (err error) {
	err = queser.VarInt(len(ret.Tags)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketTagsTags := range len(ret.Tags) {
		err = queser.EncodeString(w, ret.Tags[iPlayToClientPacketTagsTags].TagType)
		if err != nil {
			return
		}
		err = ret.Tags[iPlayToClientPacketTagsTags].Tags.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketTeams struct {
	Team    string
	Mode    string
	Anon    any
	Players any
}

var PlayToClientPacketTeamsModeMap = map[int8]string{0: "add", 1: "remove", 2: "change", 3: "join", 4: "leave"}
var PlayToClientPacketTeamsAnonNameTagVisibilityMap = map[queser.VarInt]string{0: "always", 1: "never", 2: "hide_for_other_teams", 3: "hide_for_own_team"}
var PlayToClientPacketTeamsAnonCollisionRuleMap = map[queser.VarInt]string{0: "always", 1: "never", 2: "push_other_teams", 3: "push_own_team"}

func (_ PlayToClientPacketTeams) Decode(r io.Reader) (ret PlayToClientPacketTeams, err error) {
	ret.Team, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var PlayToClientPacketTeamsModeKey int8
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTeamsModeKey)
	if err != nil {
		return
	}
	ret.Mode, err = queser.ErroringIndex(PlayToClientPacketTeamsModeMap, PlayToClientPacketTeamsModeKey)
	if err != nil {
		return
	}
	switch ret.Mode {
	case "add":
		var PlayToClientPacketTeamsAnonTmp struct {
			Name              nbt.Anon
			Flags             uint8
			NameTagVisibility string
			CollisionRule     string
			Formatting        queser.VarInt
			Prefix            nbt.Anon
			Suffix            nbt.Anon
		}
		PlayToClientPacketTeamsAnonTmp.Name, err = PlayToClientPacketTeamsAnonTmp.Name.Decode(r)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTeamsAnonTmp.Flags)
		if err != nil {
			return
		}
		var PlayToClientPacketTeamsAnonNameTagVisibilityKey queser.VarInt
		PlayToClientPacketTeamsAnonNameTagVisibilityKey, err = PlayToClientPacketTeamsAnonNameTagVisibilityKey.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsAnonTmp.NameTagVisibility, err = queser.ErroringIndex(PlayToClientPacketTeamsAnonNameTagVisibilityMap, PlayToClientPacketTeamsAnonNameTagVisibilityKey)
		if err != nil {
			return
		}
		var PlayToClientPacketTeamsAnonCollisionRuleKey queser.VarInt
		PlayToClientPacketTeamsAnonCollisionRuleKey, err = PlayToClientPacketTeamsAnonCollisionRuleKey.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsAnonTmp.CollisionRule, err = queser.ErroringIndex(PlayToClientPacketTeamsAnonCollisionRuleMap, PlayToClientPacketTeamsAnonCollisionRuleKey)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsAnonTmp.Formatting, err = PlayToClientPacketTeamsAnonTmp.Formatting.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsAnonTmp.Prefix, err = PlayToClientPacketTeamsAnonTmp.Prefix.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsAnonTmp.Suffix, err = PlayToClientPacketTeamsAnonTmp.Suffix.Decode(r)
		if err != nil {
			return
		}
		ret.Anon = PlayToClientPacketTeamsAnonTmp
	case "change":
		var PlayToClientPacketTeamsAnonTmp struct {
			Name              nbt.Anon
			Flags             uint8
			NameTagVisibility string
			CollisionRule     string
			Formatting        queser.VarInt
			Prefix            nbt.Anon
			Suffix            nbt.Anon
		}
		PlayToClientPacketTeamsAnonTmp.Name, err = PlayToClientPacketTeamsAnonTmp.Name.Decode(r)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTeamsAnonTmp.Flags)
		if err != nil {
			return
		}
		var PlayToClientPacketTeamsAnonNameTagVisibilityKey queser.VarInt
		PlayToClientPacketTeamsAnonNameTagVisibilityKey, err = PlayToClientPacketTeamsAnonNameTagVisibilityKey.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsAnonTmp.NameTagVisibility, err = queser.ErroringIndex(PlayToClientPacketTeamsAnonNameTagVisibilityMap, PlayToClientPacketTeamsAnonNameTagVisibilityKey)
		if err != nil {
			return
		}
		var PlayToClientPacketTeamsAnonCollisionRuleKey queser.VarInt
		PlayToClientPacketTeamsAnonCollisionRuleKey, err = PlayToClientPacketTeamsAnonCollisionRuleKey.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsAnonTmp.CollisionRule, err = queser.ErroringIndex(PlayToClientPacketTeamsAnonCollisionRuleMap, PlayToClientPacketTeamsAnonCollisionRuleKey)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsAnonTmp.Formatting, err = PlayToClientPacketTeamsAnonTmp.Formatting.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsAnonTmp.Prefix, err = PlayToClientPacketTeamsAnonTmp.Prefix.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsAnonTmp.Suffix, err = PlayToClientPacketTeamsAnonTmp.Suffix.Decode(r)
		if err != nil {
			return
		}
		ret.Anon = PlayToClientPacketTeamsAnonTmp
	default:
		var PlayToClientPacketTeamsAnonTmp queser.Void
		PlayToClientPacketTeamsAnonTmp, err = PlayToClientPacketTeamsAnonTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Anon = PlayToClientPacketTeamsAnonTmp
	}
	switch ret.Mode {
	case "add":
		var PlayToClientPacketTeamsPlayersTmp []string
		var lPlayToClientPacketTeamsPlayers queser.VarInt
		lPlayToClientPacketTeamsPlayers, err = lPlayToClientPacketTeamsPlayers.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsPlayersTmp = []string{}
		for range lPlayToClientPacketTeamsPlayers {
			var PlayToClientPacketTeamsPlayersElement string
			PlayToClientPacketTeamsPlayersElement, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			PlayToClientPacketTeamsPlayersTmp = append(PlayToClientPacketTeamsPlayersTmp, PlayToClientPacketTeamsPlayersElement)
		}
		ret.Players = PlayToClientPacketTeamsPlayersTmp
	case "join":
		var PlayToClientPacketTeamsPlayersTmp []string
		var lPlayToClientPacketTeamsPlayers queser.VarInt
		lPlayToClientPacketTeamsPlayers, err = lPlayToClientPacketTeamsPlayers.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsPlayersTmp = []string{}
		for range lPlayToClientPacketTeamsPlayers {
			var PlayToClientPacketTeamsPlayersElement string
			PlayToClientPacketTeamsPlayersElement, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			PlayToClientPacketTeamsPlayersTmp = append(PlayToClientPacketTeamsPlayersTmp, PlayToClientPacketTeamsPlayersElement)
		}
		ret.Players = PlayToClientPacketTeamsPlayersTmp
	case "leave":
		var PlayToClientPacketTeamsPlayersTmp []string
		var lPlayToClientPacketTeamsPlayers queser.VarInt
		lPlayToClientPacketTeamsPlayers, err = lPlayToClientPacketTeamsPlayers.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTeamsPlayersTmp = []string{}
		for range lPlayToClientPacketTeamsPlayers {
			var PlayToClientPacketTeamsPlayersElement string
			PlayToClientPacketTeamsPlayersElement, err = queser.DecodeString(r)
			if err != nil {
				return
			}
			PlayToClientPacketTeamsPlayersTmp = append(PlayToClientPacketTeamsPlayersTmp, PlayToClientPacketTeamsPlayersElement)
		}
		ret.Players = PlayToClientPacketTeamsPlayersTmp
	default:
		var PlayToClientPacketTeamsPlayersTmp queser.Void
		PlayToClientPacketTeamsPlayersTmp, err = PlayToClientPacketTeamsPlayersTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Players = PlayToClientPacketTeamsPlayersTmp
	}
	return
}

var PlayToClientPacketTeamsModeReverseMap = map[string]int8{"add": 0, "remove": 1, "change": 2, "join": 3, "leave": 4}
var PlayToClientPacketTeamsAnonNameTagVisibilityReverseMap = map[string]queser.VarInt{"always": 0, "never": 1, "hide_for_other_teams": 2, "hide_for_own_team": 3}
var PlayToClientPacketTeamsAnonCollisionRuleReverseMap = map[string]queser.VarInt{"always": 0, "never": 1, "push_other_teams": 2, "push_own_team": 3}

func (ret PlayToClientPacketTeams) Encode(w io.Writer) (err error) {
	err = queser.EncodeString(w, ret.Team)
	if err != nil {
		return
	}
	var vPlayToClientPacketTeamsMode int8
	vPlayToClientPacketTeamsMode, err = queser.ErroringIndex(PlayToClientPacketTeamsModeReverseMap, ret.Mode)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, vPlayToClientPacketTeamsMode)
	if err != nil {
		return
	}
	switch ret.Mode {
	case "add":
		PlayToClientPacketTeamsAnon, ok := ret.Anon.(struct {
			Name              nbt.Anon
			Flags             uint8
			NameTagVisibility string
			CollisionRule     string
			Formatting        queser.VarInt
			Prefix            nbt.Anon
			Suffix            nbt.Anon
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketTeamsAnon.Name.Encode(w)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToClientPacketTeamsAnon.Flags)
		if err != nil {
			return
		}
		var vPlayToClientPacketTeamsAnonNameTagVisibility queser.VarInt
		vPlayToClientPacketTeamsAnonNameTagVisibility, err = queser.ErroringIndex(PlayToClientPacketTeamsAnonNameTagVisibilityReverseMap, PlayToClientPacketTeamsAnon.NameTagVisibility)
		if err != nil {
			return
		}
		err = vPlayToClientPacketTeamsAnonNameTagVisibility.Encode(w)
		if err != nil {
			return
		}
		var vPlayToClientPacketTeamsAnonCollisionRule queser.VarInt
		vPlayToClientPacketTeamsAnonCollisionRule, err = queser.ErroringIndex(PlayToClientPacketTeamsAnonCollisionRuleReverseMap, PlayToClientPacketTeamsAnon.CollisionRule)
		if err != nil {
			return
		}
		err = vPlayToClientPacketTeamsAnonCollisionRule.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientPacketTeamsAnon.Formatting.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientPacketTeamsAnon.Prefix.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientPacketTeamsAnon.Suffix.Encode(w)
		if err != nil {
			return
		}
	case "change":
		PlayToClientPacketTeamsAnon, ok := ret.Anon.(struct {
			Name              nbt.Anon
			Flags             uint8
			NameTagVisibility string
			CollisionRule     string
			Formatting        queser.VarInt
			Prefix            nbt.Anon
			Suffix            nbt.Anon
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketTeamsAnon.Name.Encode(w)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToClientPacketTeamsAnon.Flags)
		if err != nil {
			return
		}
		var vPlayToClientPacketTeamsAnonNameTagVisibility queser.VarInt
		vPlayToClientPacketTeamsAnonNameTagVisibility, err = queser.ErroringIndex(PlayToClientPacketTeamsAnonNameTagVisibilityReverseMap, PlayToClientPacketTeamsAnon.NameTagVisibility)
		if err != nil {
			return
		}
		err = vPlayToClientPacketTeamsAnonNameTagVisibility.Encode(w)
		if err != nil {
			return
		}
		var vPlayToClientPacketTeamsAnonCollisionRule queser.VarInt
		vPlayToClientPacketTeamsAnonCollisionRule, err = queser.ErroringIndex(PlayToClientPacketTeamsAnonCollisionRuleReverseMap, PlayToClientPacketTeamsAnon.CollisionRule)
		if err != nil {
			return
		}
		err = vPlayToClientPacketTeamsAnonCollisionRule.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientPacketTeamsAnon.Formatting.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientPacketTeamsAnon.Prefix.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientPacketTeamsAnon.Suffix.Encode(w)
		if err != nil {
			return
		}
	default:
		_, ok := ret.Anon.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Anon.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	switch ret.Mode {
	case "add":
		PlayToClientPacketTeamsPlayers, ok := ret.Players.([]string)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.VarInt(len(PlayToClientPacketTeamsPlayers)).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketTeamsPlayers := range len(PlayToClientPacketTeamsPlayers) {
			err = queser.EncodeString(w, PlayToClientPacketTeamsPlayers[iPlayToClientPacketTeamsPlayers])
			if err != nil {
				return
			}
		}
	case "join":
		PlayToClientPacketTeamsPlayers, ok := ret.Players.([]string)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.VarInt(len(PlayToClientPacketTeamsPlayers)).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketTeamsPlayers := range len(PlayToClientPacketTeamsPlayers) {
			err = queser.EncodeString(w, PlayToClientPacketTeamsPlayers[iPlayToClientPacketTeamsPlayers])
			if err != nil {
				return
			}
		}
	case "leave":
		PlayToClientPacketTeamsPlayers, ok := ret.Players.([]string)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.VarInt(len(PlayToClientPacketTeamsPlayers)).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketTeamsPlayers := range len(PlayToClientPacketTeamsPlayers) {
			err = queser.EncodeString(w, PlayToClientPacketTeamsPlayers[iPlayToClientPacketTeamsPlayers])
			if err != nil {
				return
			}
		}
	default:
		_, ok := ret.Players.(queser.Void)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = ret.Players.(queser.Void).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketTestInstanceBlockStatus struct {
	Status nbt.Anon
	Size   *Vec3i
}

func (_ PlayToClientPacketTestInstanceBlockStatus) Decode(r io.Reader) (ret PlayToClientPacketTestInstanceBlockStatus, err error) {
	ret.Status, err = ret.Status.Decode(r)
	if err != nil {
		return
	}
	var PlayToClientPacketTestInstanceBlockStatusSizePresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTestInstanceBlockStatusSizePresent)
	if err != nil {
		return
	}
	if PlayToClientPacketTestInstanceBlockStatusSizePresent {
		var PlayToClientPacketTestInstanceBlockStatusSizePresentValue Vec3i
		PlayToClientPacketTestInstanceBlockStatusSizePresentValue, err = PlayToClientPacketTestInstanceBlockStatusSizePresentValue.Decode(r)
		if err != nil {
			return
		}
		ret.Size = &PlayToClientPacketTestInstanceBlockStatusSizePresentValue
	}
	return
}
func (ret PlayToClientPacketTestInstanceBlockStatus) Encode(w io.Writer) (err error) {
	err = ret.Status.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Size != nil)
	if err != nil {
		return
	}
	if ret.Size != nil {
		err = (*ret.Size).Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketTileEntityData struct {
	Location Position
	Action   queser.VarInt
	NbtData  nbt.Anon
}

func (_ PlayToClientPacketTileEntityData) Decode(r io.Reader) (ret PlayToClientPacketTileEntityData, err error) {
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	ret.Action, err = ret.Action.Decode(r)
	if err != nil {
		return
	}
	ret.NbtData, err = ret.NbtData.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketTileEntityData) Encode(w io.Writer) (err error) {
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = ret.Action.Encode(w)
	if err != nil {
		return
	}
	err = ret.NbtData.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketTrackedWaypoint struct {
	Operation string
	Waypoint  struct {
		HasUUID bool
		Anon    any
		Icon    struct {
			Style string
			Color *int32
		}
		Type string
		Data any
	}
}

var PlayToClientPacketTrackedWaypointOperationMap = map[queser.VarInt]string{0: "track", 1: "untrack", 2: "update"}
var PlayToClientPacketTrackedWaypointWaypointTypeMap = map[queser.VarInt]string{0: "empty", 1: "vec3i", 2: "chunk", 3: "azimuth"}

func (_ PlayToClientPacketTrackedWaypoint) Decode(r io.Reader) (ret PlayToClientPacketTrackedWaypoint, err error) {
	var PlayToClientPacketTrackedWaypointOperationKey queser.VarInt
	PlayToClientPacketTrackedWaypointOperationKey, err = PlayToClientPacketTrackedWaypointOperationKey.Decode(r)
	if err != nil {
		return
	}
	ret.Operation, err = queser.ErroringIndex(PlayToClientPacketTrackedWaypointOperationMap, PlayToClientPacketTrackedWaypointOperationKey)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Waypoint.HasUUID)
	if err != nil {
		return
	}
	switch ret.Waypoint.HasUUID {
	case false:
		var PlayToClientPacketTrackedWaypointWaypointAnonTmp struct {
			Id string
		}
		PlayToClientPacketTrackedWaypointWaypointAnonTmp.Id, err = queser.DecodeString(r)
		if err != nil {
			return
		}
		ret.Waypoint.Anon = PlayToClientPacketTrackedWaypointWaypointAnonTmp
	case true:
		var PlayToClientPacketTrackedWaypointWaypointAnonTmp struct {
			Uuid uuid.UUID
		}
		_, err = io.ReadFull(r, PlayToClientPacketTrackedWaypointWaypointAnonTmp.Uuid[:])
		if err != nil {
			return
		}
		ret.Waypoint.Anon = PlayToClientPacketTrackedWaypointWaypointAnonTmp
	}
	ret.Waypoint.Icon.Style, err = queser.DecodeString(r)
	if err != nil {
		return
	}
	var PlayToClientPacketTrackedWaypointWaypointIconColorPresent bool
	err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTrackedWaypointWaypointIconColorPresent)
	if err != nil {
		return
	}
	if PlayToClientPacketTrackedWaypointWaypointIconColorPresent {
		var PlayToClientPacketTrackedWaypointWaypointIconColorPresentValue int32
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTrackedWaypointWaypointIconColorPresentValue)
		if err != nil {
			return
		}
		ret.Waypoint.Icon.Color = &PlayToClientPacketTrackedWaypointWaypointIconColorPresentValue
	}
	var PlayToClientPacketTrackedWaypointWaypointTypeKey queser.VarInt
	PlayToClientPacketTrackedWaypointWaypointTypeKey, err = PlayToClientPacketTrackedWaypointWaypointTypeKey.Decode(r)
	if err != nil {
		return
	}
	ret.Waypoint.Type, err = queser.ErroringIndex(PlayToClientPacketTrackedWaypointWaypointTypeMap, PlayToClientPacketTrackedWaypointWaypointTypeKey)
	if err != nil {
		return
	}
	switch ret.Waypoint.Type {
	case "azimuth":
		var PlayToClientPacketTrackedWaypointWaypointDataTmp float32
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTrackedWaypointWaypointDataTmp)
		if err != nil {
			return
		}
		ret.Waypoint.Data = PlayToClientPacketTrackedWaypointWaypointDataTmp
	case "chunk":
		var PlayToClientPacketTrackedWaypointWaypointDataTmp struct {
			ChunkX queser.VarInt
			ChunkZ queser.VarInt
		}
		PlayToClientPacketTrackedWaypointWaypointDataTmp.ChunkX, err = PlayToClientPacketTrackedWaypointWaypointDataTmp.ChunkX.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTrackedWaypointWaypointDataTmp.ChunkZ, err = PlayToClientPacketTrackedWaypointWaypointDataTmp.ChunkZ.Decode(r)
		if err != nil {
			return
		}
		ret.Waypoint.Data = PlayToClientPacketTrackedWaypointWaypointDataTmp
	case "vec3i":
		var PlayToClientPacketTrackedWaypointWaypointDataTmp Vec3i
		PlayToClientPacketTrackedWaypointWaypointDataTmp, err = PlayToClientPacketTrackedWaypointWaypointDataTmp.Decode(r)
		if err != nil {
			return
		}
		ret.Waypoint.Data = PlayToClientPacketTrackedWaypointWaypointDataTmp
	}
	return
}

var PlayToClientPacketTrackedWaypointOperationReverseMap = map[string]queser.VarInt{"track": 0, "untrack": 1, "update": 2}
var PlayToClientPacketTrackedWaypointWaypointTypeReverseMap = map[string]queser.VarInt{"empty": 0, "vec3i": 1, "chunk": 2, "azimuth": 3}

func (ret PlayToClientPacketTrackedWaypoint) Encode(w io.Writer) (err error) {
	var vPlayToClientPacketTrackedWaypointOperation queser.VarInt
	vPlayToClientPacketTrackedWaypointOperation, err = queser.ErroringIndex(PlayToClientPacketTrackedWaypointOperationReverseMap, ret.Operation)
	if err != nil {
		return
	}
	err = vPlayToClientPacketTrackedWaypointOperation.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Waypoint.HasUUID)
	if err != nil {
		return
	}
	switch ret.Waypoint.HasUUID {
	case false:
		PlayToClientPacketTrackedWaypointWaypointAnon, ok := ret.Waypoint.Anon.(struct {
			Id string
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = queser.EncodeString(w, PlayToClientPacketTrackedWaypointWaypointAnon.Id)
		if err != nil {
			return
		}
	case true:
		PlayToClientPacketTrackedWaypointWaypointAnon, ok := ret.Waypoint.Anon.(struct {
			Uuid uuid.UUID
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		_, err = w.Write(PlayToClientPacketTrackedWaypointWaypointAnon.Uuid[:])
		if err != nil {
			return
		}
	}
	err = queser.EncodeString(w, ret.Waypoint.Icon.Style)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Waypoint.Icon.Color != nil)
	if err != nil {
		return
	}
	if ret.Waypoint.Icon.Color != nil {
		err = binary.Write(w, binary.BigEndian, *ret.Waypoint.Icon.Color)
		if err != nil {
			return
		}
	}
	var vPlayToClientPacketTrackedWaypointWaypointType queser.VarInt
	vPlayToClientPacketTrackedWaypointWaypointType, err = queser.ErroringIndex(PlayToClientPacketTrackedWaypointWaypointTypeReverseMap, ret.Waypoint.Type)
	if err != nil {
		return
	}
	err = vPlayToClientPacketTrackedWaypointWaypointType.Encode(w)
	if err != nil {
		return
	}
	switch ret.Waypoint.Type {
	case "azimuth":
		PlayToClientPacketTrackedWaypointWaypointData, ok := ret.Waypoint.Data.(float32)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = binary.Write(w, binary.BigEndian, PlayToClientPacketTrackedWaypointWaypointData)
		if err != nil {
			return
		}
	case "chunk":
		PlayToClientPacketTrackedWaypointWaypointData, ok := ret.Waypoint.Data.(struct {
			ChunkX queser.VarInt
			ChunkZ queser.VarInt
		})
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketTrackedWaypointWaypointData.ChunkX.Encode(w)
		if err != nil {
			return
		}
		err = PlayToClientPacketTrackedWaypointWaypointData.ChunkZ.Encode(w)
		if err != nil {
			return
		}
	case "vec3i":
		PlayToClientPacketTrackedWaypointWaypointData, ok := ret.Waypoint.Data.(Vec3i)
		if !ok {
			err = queser.BadTypeError
			return
		}
		err = PlayToClientPacketTrackedWaypointWaypointData.Encode(w)
		if err != nil {
			return
		}
	}
	return
}

type PlayToClientPacketTradeList struct {
	WindowId ContainerID
	Trades   []struct {
		InputItem1 struct {
			ItemId     queser.VarInt
			ItemCount  queser.VarInt
			Components ExactComponentMatcher
		}
		OutputItem Slot
		InputItem2 *struct {
			ItemId     queser.VarInt
			ItemCount  queser.VarInt
			Components ExactComponentMatcher
		}
		TradeDisabled      bool
		NbTradeUses        int32
		MaximumNbTradeUses int32
		Xp                 int32
		SpecialPrice       int32
		PriceMultiplier    float32
		Demand             int32
	}
	VillagerLevel     queser.VarInt
	Experience        queser.VarInt
	IsRegularVillager bool
	CanRestock        bool
}

func (_ PlayToClientPacketTradeList) Decode(r io.Reader) (ret PlayToClientPacketTradeList, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	var lPlayToClientPacketTradeListTrades queser.VarInt
	lPlayToClientPacketTradeListTrades, err = lPlayToClientPacketTradeListTrades.Decode(r)
	if err != nil {
		return
	}
	ret.Trades = []struct {
		InputItem1 struct {
			ItemId     queser.VarInt
			ItemCount  queser.VarInt
			Components ExactComponentMatcher
		}
		OutputItem Slot
		InputItem2 *struct {
			ItemId     queser.VarInt
			ItemCount  queser.VarInt
			Components ExactComponentMatcher
		}
		TradeDisabled      bool
		NbTradeUses        int32
		MaximumNbTradeUses int32
		Xp                 int32
		SpecialPrice       int32
		PriceMultiplier    float32
		Demand             int32
	}{}
	for range lPlayToClientPacketTradeListTrades {
		var PlayToClientPacketTradeListTradesElement struct {
			InputItem1 struct {
				ItemId     queser.VarInt
				ItemCount  queser.VarInt
				Components ExactComponentMatcher
			}
			OutputItem Slot
			InputItem2 *struct {
				ItemId     queser.VarInt
				ItemCount  queser.VarInt
				Components ExactComponentMatcher
			}
			TradeDisabled      bool
			NbTradeUses        int32
			MaximumNbTradeUses int32
			Xp                 int32
			SpecialPrice       int32
			PriceMultiplier    float32
			Demand             int32
		}
		PlayToClientPacketTradeListTradesElement.InputItem1.ItemId, err = PlayToClientPacketTradeListTradesElement.InputItem1.ItemId.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTradeListTradesElement.InputItem1.ItemCount, err = PlayToClientPacketTradeListTradesElement.InputItem1.ItemCount.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTradeListTradesElement.InputItem1.Components, err = PlayToClientPacketTradeListTradesElement.InputItem1.Components.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketTradeListTradesElement.OutputItem, err = PlayToClientPacketTradeListTradesElement.OutputItem.Decode(r)
		if err != nil {
			return
		}
		var PlayToClientPacketTradeListTradesElementInputItem2Present bool
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTradeListTradesElementInputItem2Present)
		if err != nil {
			return
		}
		if PlayToClientPacketTradeListTradesElementInputItem2Present {
			var PlayToClientPacketTradeListTradesElementInputItem2PresentValue struct {
				ItemId     queser.VarInt
				ItemCount  queser.VarInt
				Components ExactComponentMatcher
			}
			PlayToClientPacketTradeListTradesElementInputItem2PresentValue.ItemId, err = PlayToClientPacketTradeListTradesElementInputItem2PresentValue.ItemId.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketTradeListTradesElementInputItem2PresentValue.ItemCount, err = PlayToClientPacketTradeListTradesElementInputItem2PresentValue.ItemCount.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketTradeListTradesElementInputItem2PresentValue.Components, err = PlayToClientPacketTradeListTradesElementInputItem2PresentValue.Components.Decode(r)
			if err != nil {
				return
			}
			PlayToClientPacketTradeListTradesElement.InputItem2 = &PlayToClientPacketTradeListTradesElementInputItem2PresentValue
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTradeListTradesElement.TradeDisabled)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTradeListTradesElement.NbTradeUses)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTradeListTradesElement.MaximumNbTradeUses)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTradeListTradesElement.Xp)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTradeListTradesElement.SpecialPrice)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTradeListTradesElement.PriceMultiplier)
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketTradeListTradesElement.Demand)
		if err != nil {
			return
		}
		ret.Trades = append(ret.Trades, PlayToClientPacketTradeListTradesElement)
	}
	ret.VillagerLevel, err = ret.VillagerLevel.Decode(r)
	if err != nil {
		return
	}
	ret.Experience, err = ret.Experience.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.IsRegularVillager)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.CanRestock)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketTradeList) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Trades)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketTradeListTrades := range len(ret.Trades) {
		err = ret.Trades[iPlayToClientPacketTradeListTrades].InputItem1.ItemId.Encode(w)
		if err != nil {
			return
		}
		err = ret.Trades[iPlayToClientPacketTradeListTrades].InputItem1.ItemCount.Encode(w)
		if err != nil {
			return
		}
		err = ret.Trades[iPlayToClientPacketTradeListTrades].InputItem1.Components.Encode(w)
		if err != nil {
			return
		}
		err = ret.Trades[iPlayToClientPacketTradeListTrades].OutputItem.Encode(w)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Trades[iPlayToClientPacketTradeListTrades].InputItem2 != nil)
		if err != nil {
			return
		}
		if ret.Trades[iPlayToClientPacketTradeListTrades].InputItem2 != nil {
			err = (*ret.Trades[iPlayToClientPacketTradeListTrades].InputItem2).ItemId.Encode(w)
			if err != nil {
				return
			}
			err = (*ret.Trades[iPlayToClientPacketTradeListTrades].InputItem2).ItemCount.Encode(w)
			if err != nil {
				return
			}
			err = (*ret.Trades[iPlayToClientPacketTradeListTrades].InputItem2).Components.Encode(w)
			if err != nil {
				return
			}
		}
		err = binary.Write(w, binary.BigEndian, ret.Trades[iPlayToClientPacketTradeListTrades].TradeDisabled)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Trades[iPlayToClientPacketTradeListTrades].NbTradeUses)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Trades[iPlayToClientPacketTradeListTrades].MaximumNbTradeUses)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Trades[iPlayToClientPacketTradeListTrades].Xp)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Trades[iPlayToClientPacketTradeListTrades].SpecialPrice)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Trades[iPlayToClientPacketTradeListTrades].PriceMultiplier)
		if err != nil {
			return
		}
		err = binary.Write(w, binary.BigEndian, ret.Trades[iPlayToClientPacketTradeListTrades].Demand)
		if err != nil {
			return
		}
	}
	err = ret.VillagerLevel.Encode(w)
	if err != nil {
		return
	}
	err = ret.Experience.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.IsRegularVillager)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.CanRestock)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketUnloadChunk struct {
	ChunkZ int32
	ChunkX int32
}

func (_ PlayToClientPacketUnloadChunk) Decode(r io.Reader) (ret PlayToClientPacketUnloadChunk, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.ChunkZ)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.ChunkX)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketUnloadChunk) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.ChunkZ)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.ChunkX)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketUpdateHealth struct {
	Health         float32
	Food           queser.VarInt
	FoodSaturation float32
}

func (_ PlayToClientPacketUpdateHealth) Decode(r io.Reader) (ret PlayToClientPacketUpdateHealth, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Health)
	if err != nil {
		return
	}
	ret.Food, err = ret.Food.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.FoodSaturation)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketUpdateHealth) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Health)
	if err != nil {
		return
	}
	err = ret.Food.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.FoodSaturation)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketUpdateLight struct {
	ChunkX              queser.VarInt
	ChunkZ              queser.VarInt
	SkyLightMask        []int64
	BlockLightMask      []int64
	EmptySkyLightMask   []int64
	EmptyBlockLightMask []int64
	SkyLight            [][]uint8
	BlockLight          [][]uint8
}

func (_ PlayToClientPacketUpdateLight) Decode(r io.Reader) (ret PlayToClientPacketUpdateLight, err error) {
	ret.ChunkX, err = ret.ChunkX.Decode(r)
	if err != nil {
		return
	}
	ret.ChunkZ, err = ret.ChunkZ.Decode(r)
	if err != nil {
		return
	}
	var lPlayToClientPacketUpdateLightSkyLightMask queser.VarInt
	lPlayToClientPacketUpdateLightSkyLightMask, err = lPlayToClientPacketUpdateLightSkyLightMask.Decode(r)
	if err != nil {
		return
	}
	ret.SkyLightMask = []int64{}
	for range lPlayToClientPacketUpdateLightSkyLightMask {
		var PlayToClientPacketUpdateLightSkyLightMaskElement int64
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketUpdateLightSkyLightMaskElement)
		if err != nil {
			return
		}
		ret.SkyLightMask = append(ret.SkyLightMask, PlayToClientPacketUpdateLightSkyLightMaskElement)
	}
	var lPlayToClientPacketUpdateLightBlockLightMask queser.VarInt
	lPlayToClientPacketUpdateLightBlockLightMask, err = lPlayToClientPacketUpdateLightBlockLightMask.Decode(r)
	if err != nil {
		return
	}
	ret.BlockLightMask = []int64{}
	for range lPlayToClientPacketUpdateLightBlockLightMask {
		var PlayToClientPacketUpdateLightBlockLightMaskElement int64
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketUpdateLightBlockLightMaskElement)
		if err != nil {
			return
		}
		ret.BlockLightMask = append(ret.BlockLightMask, PlayToClientPacketUpdateLightBlockLightMaskElement)
	}
	var lPlayToClientPacketUpdateLightEmptySkyLightMask queser.VarInt
	lPlayToClientPacketUpdateLightEmptySkyLightMask, err = lPlayToClientPacketUpdateLightEmptySkyLightMask.Decode(r)
	if err != nil {
		return
	}
	ret.EmptySkyLightMask = []int64{}
	for range lPlayToClientPacketUpdateLightEmptySkyLightMask {
		var PlayToClientPacketUpdateLightEmptySkyLightMaskElement int64
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketUpdateLightEmptySkyLightMaskElement)
		if err != nil {
			return
		}
		ret.EmptySkyLightMask = append(ret.EmptySkyLightMask, PlayToClientPacketUpdateLightEmptySkyLightMaskElement)
	}
	var lPlayToClientPacketUpdateLightEmptyBlockLightMask queser.VarInt
	lPlayToClientPacketUpdateLightEmptyBlockLightMask, err = lPlayToClientPacketUpdateLightEmptyBlockLightMask.Decode(r)
	if err != nil {
		return
	}
	ret.EmptyBlockLightMask = []int64{}
	for range lPlayToClientPacketUpdateLightEmptyBlockLightMask {
		var PlayToClientPacketUpdateLightEmptyBlockLightMaskElement int64
		err = binary.Read(r, binary.BigEndian, &PlayToClientPacketUpdateLightEmptyBlockLightMaskElement)
		if err != nil {
			return
		}
		ret.EmptyBlockLightMask = append(ret.EmptyBlockLightMask, PlayToClientPacketUpdateLightEmptyBlockLightMaskElement)
	}
	var lPlayToClientPacketUpdateLightSkyLight queser.VarInt
	lPlayToClientPacketUpdateLightSkyLight, err = lPlayToClientPacketUpdateLightSkyLight.Decode(r)
	if err != nil {
		return
	}
	ret.SkyLight = [][]uint8{}
	for range lPlayToClientPacketUpdateLightSkyLight {
		var PlayToClientPacketUpdateLightSkyLightElement []uint8
		var lPlayToClientPacketUpdateLightSkyLightElement queser.VarInt
		lPlayToClientPacketUpdateLightSkyLightElement, err = lPlayToClientPacketUpdateLightSkyLightElement.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketUpdateLightSkyLightElement = []uint8{}
		for range lPlayToClientPacketUpdateLightSkyLightElement {
			var PlayToClientPacketUpdateLightSkyLightElementElement uint8
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketUpdateLightSkyLightElementElement)
			if err != nil {
				return
			}
			PlayToClientPacketUpdateLightSkyLightElement = append(PlayToClientPacketUpdateLightSkyLightElement, PlayToClientPacketUpdateLightSkyLightElementElement)
		}
		ret.SkyLight = append(ret.SkyLight, PlayToClientPacketUpdateLightSkyLightElement)
	}
	var lPlayToClientPacketUpdateLightBlockLight queser.VarInt
	lPlayToClientPacketUpdateLightBlockLight, err = lPlayToClientPacketUpdateLightBlockLight.Decode(r)
	if err != nil {
		return
	}
	ret.BlockLight = [][]uint8{}
	for range lPlayToClientPacketUpdateLightBlockLight {
		var PlayToClientPacketUpdateLightBlockLightElement []uint8
		var lPlayToClientPacketUpdateLightBlockLightElement queser.VarInt
		lPlayToClientPacketUpdateLightBlockLightElement, err = lPlayToClientPacketUpdateLightBlockLightElement.Decode(r)
		if err != nil {
			return
		}
		PlayToClientPacketUpdateLightBlockLightElement = []uint8{}
		for range lPlayToClientPacketUpdateLightBlockLightElement {
			var PlayToClientPacketUpdateLightBlockLightElementElement uint8
			err = binary.Read(r, binary.BigEndian, &PlayToClientPacketUpdateLightBlockLightElementElement)
			if err != nil {
				return
			}
			PlayToClientPacketUpdateLightBlockLightElement = append(PlayToClientPacketUpdateLightBlockLightElement, PlayToClientPacketUpdateLightBlockLightElementElement)
		}
		ret.BlockLight = append(ret.BlockLight, PlayToClientPacketUpdateLightBlockLightElement)
	}
	return
}
func (ret PlayToClientPacketUpdateLight) Encode(w io.Writer) (err error) {
	err = ret.ChunkX.Encode(w)
	if err != nil {
		return
	}
	err = ret.ChunkZ.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.SkyLightMask)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketUpdateLightSkyLightMask := range len(ret.SkyLightMask) {
		err = binary.Write(w, binary.BigEndian, ret.SkyLightMask[iPlayToClientPacketUpdateLightSkyLightMask])
		if err != nil {
			return
		}
	}
	err = queser.VarInt(len(ret.BlockLightMask)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketUpdateLightBlockLightMask := range len(ret.BlockLightMask) {
		err = binary.Write(w, binary.BigEndian, ret.BlockLightMask[iPlayToClientPacketUpdateLightBlockLightMask])
		if err != nil {
			return
		}
	}
	err = queser.VarInt(len(ret.EmptySkyLightMask)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketUpdateLightEmptySkyLightMask := range len(ret.EmptySkyLightMask) {
		err = binary.Write(w, binary.BigEndian, ret.EmptySkyLightMask[iPlayToClientPacketUpdateLightEmptySkyLightMask])
		if err != nil {
			return
		}
	}
	err = queser.VarInt(len(ret.EmptyBlockLightMask)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketUpdateLightEmptyBlockLightMask := range len(ret.EmptyBlockLightMask) {
		err = binary.Write(w, binary.BigEndian, ret.EmptyBlockLightMask[iPlayToClientPacketUpdateLightEmptyBlockLightMask])
		if err != nil {
			return
		}
	}
	err = queser.VarInt(len(ret.SkyLight)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketUpdateLightSkyLight := range len(ret.SkyLight) {
		err = queser.VarInt(len(ret.SkyLight[iPlayToClientPacketUpdateLightSkyLight])).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketUpdateLightSkyLightInner := range len(ret.SkyLight[iPlayToClientPacketUpdateLightSkyLight]) {
			err = binary.Write(w, binary.BigEndian, ret.SkyLight[iPlayToClientPacketUpdateLightSkyLight][iPlayToClientPacketUpdateLightSkyLightInner])
			if err != nil {
				return
			}
		}
	}
	err = queser.VarInt(len(ret.BlockLight)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketUpdateLightBlockLight := range len(ret.BlockLight) {
		err = queser.VarInt(len(ret.BlockLight[iPlayToClientPacketUpdateLightBlockLight])).Encode(w)
		if err != nil {
			return
		}
		for iPlayToClientPacketUpdateLightBlockLightInner := range len(ret.BlockLight[iPlayToClientPacketUpdateLightBlockLight]) {
			err = binary.Write(w, binary.BigEndian, ret.BlockLight[iPlayToClientPacketUpdateLightBlockLight][iPlayToClientPacketUpdateLightBlockLightInner])
			if err != nil {
				return
			}
		}
	}
	return
}

type PlayToClientPacketUpdateTime struct {
	Age         int64
	Time        int64
	TickDayTime bool
}

func (_ PlayToClientPacketUpdateTime) Decode(r io.Reader) (ret PlayToClientPacketUpdateTime, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Age)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Time)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.TickDayTime)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketUpdateTime) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Age)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Time)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.TickDayTime)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketUpdateViewDistance struct {
	ViewDistance queser.VarInt
}

func (_ PlayToClientPacketUpdateViewDistance) Decode(r io.Reader) (ret PlayToClientPacketUpdateViewDistance, err error) {
	ret.ViewDistance, err = ret.ViewDistance.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketUpdateViewDistance) Encode(w io.Writer) (err error) {
	err = ret.ViewDistance.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketUpdateViewPosition struct {
	ChunkX queser.VarInt
	ChunkZ queser.VarInt
}

func (_ PlayToClientPacketUpdateViewPosition) Decode(r io.Reader) (ret PlayToClientPacketUpdateViewPosition, err error) {
	ret.ChunkX, err = ret.ChunkX.Decode(r)
	if err != nil {
		return
	}
	ret.ChunkZ, err = ret.ChunkZ.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketUpdateViewPosition) Encode(w io.Writer) (err error) {
	err = ret.ChunkX.Encode(w)
	if err != nil {
		return
	}
	err = ret.ChunkZ.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketVehicleMove struct {
	X     float64
	Y     float64
	Z     float64
	Yaw   float32
	Pitch float32
}

func (_ PlayToClientPacketVehicleMove) Decode(r io.Reader) (ret PlayToClientPacketVehicleMove, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Pitch)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketVehicleMove) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Yaw)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Pitch)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketWindowItems struct {
	WindowId    ContainerID
	StateId     queser.VarInt
	Items       []Slot
	CarriedItem Slot
}

func (_ PlayToClientPacketWindowItems) Decode(r io.Reader) (ret PlayToClientPacketWindowItems, err error) {
	ret.WindowId, err = ret.WindowId.Decode(r)
	if err != nil {
		return
	}
	ret.StateId, err = ret.StateId.Decode(r)
	if err != nil {
		return
	}
	var lPlayToClientPacketWindowItemsItems queser.VarInt
	lPlayToClientPacketWindowItemsItems, err = lPlayToClientPacketWindowItemsItems.Decode(r)
	if err != nil {
		return
	}
	ret.Items = []Slot{}
	for range lPlayToClientPacketWindowItemsItems {
		var PlayToClientPacketWindowItemsItemsElement Slot
		PlayToClientPacketWindowItemsItemsElement, err = PlayToClientPacketWindowItemsItemsElement.Decode(r)
		if err != nil {
			return
		}
		ret.Items = append(ret.Items, PlayToClientPacketWindowItemsItemsElement)
	}
	ret.CarriedItem, err = ret.CarriedItem.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketWindowItems) Encode(w io.Writer) (err error) {
	err = ret.WindowId.Encode(w)
	if err != nil {
		return
	}
	err = ret.StateId.Encode(w)
	if err != nil {
		return
	}
	err = queser.VarInt(len(ret.Items)).Encode(w)
	if err != nil {
		return
	}
	for iPlayToClientPacketWindowItemsItems := range len(ret.Items) {
		err = ret.Items[iPlayToClientPacketWindowItemsItems].Encode(w)
		if err != nil {
			return
		}
	}
	err = ret.CarriedItem.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketWorldBorderCenter struct {
	X float64
	Z float64
}

func (_ PlayToClientPacketWorldBorderCenter) Decode(r io.Reader) (ret PlayToClientPacketWorldBorderCenter, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketWorldBorderCenter) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketWorldBorderLerpSize struct {
	OldDiameter float64
	NewDiameter float64
	Speed       queser.VarInt
}

func (_ PlayToClientPacketWorldBorderLerpSize) Decode(r io.Reader) (ret PlayToClientPacketWorldBorderLerpSize, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.OldDiameter)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.NewDiameter)
	if err != nil {
		return
	}
	ret.Speed, err = ret.Speed.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketWorldBorderLerpSize) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.OldDiameter)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.NewDiameter)
	if err != nil {
		return
	}
	err = ret.Speed.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketWorldBorderSize struct {
	Diameter float64
}

func (_ PlayToClientPacketWorldBorderSize) Decode(r io.Reader) (ret PlayToClientPacketWorldBorderSize, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.Diameter)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketWorldBorderSize) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.Diameter)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketWorldBorderWarningDelay struct {
	WarningTime queser.VarInt
}

func (_ PlayToClientPacketWorldBorderWarningDelay) Decode(r io.Reader) (ret PlayToClientPacketWorldBorderWarningDelay, err error) {
	ret.WarningTime, err = ret.WarningTime.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketWorldBorderWarningDelay) Encode(w io.Writer) (err error) {
	err = ret.WarningTime.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketWorldBorderWarningReach struct {
	WarningBlocks queser.VarInt
}

func (_ PlayToClientPacketWorldBorderWarningReach) Decode(r io.Reader) (ret PlayToClientPacketWorldBorderWarningReach, err error) {
	ret.WarningBlocks, err = ret.WarningBlocks.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketWorldBorderWarningReach) Encode(w io.Writer) (err error) {
	err = ret.WarningBlocks.Encode(w)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketWorldEvent struct {
	EffectId int32
	Location Position
	Data     int32
	Global   bool
}

func (_ PlayToClientPacketWorldEvent) Decode(r io.Reader) (ret PlayToClientPacketWorldEvent, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.EffectId)
	if err != nil {
		return
	}
	ret.Location, err = ret.Location.Decode(r)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Data)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Global)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketWorldEvent) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.EffectId)
	if err != nil {
		return
	}
	err = ret.Location.Encode(w)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Data)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Global)
	if err != nil {
		return
	}
	return
}

type PlayToClientPacketWorldParticles struct {
	LongDistance   bool
	AlwaysShow     bool
	X              float64
	Y              float64
	Z              float64
	OffsetX        float32
	OffsetY        float32
	OffsetZ        float32
	VelocityOffset float32
	Amount         int32
	Particle       Particle
}

func (_ PlayToClientPacketWorldParticles) Decode(r io.Reader) (ret PlayToClientPacketWorldParticles, err error) {
	err = binary.Read(r, binary.BigEndian, &ret.LongDistance)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.AlwaysShow)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.X)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Y)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Z)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OffsetX)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OffsetY)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.OffsetZ)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.VelocityOffset)
	if err != nil {
		return
	}
	err = binary.Read(r, binary.BigEndian, &ret.Amount)
	if err != nil {
		return
	}
	ret.Particle, err = ret.Particle.Decode(r)
	if err != nil {
		return
	}
	return
}
func (ret PlayToClientPacketWorldParticles) Encode(w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, ret.LongDistance)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.AlwaysShow)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.X)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Y)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Z)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OffsetX)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OffsetY)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.OffsetZ)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.VelocityOffset)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, ret.Amount)
	if err != nil {
		return
	}
	err = ret.Particle.Encode(w)
	if err != nil {
		return
	}
	return
}
