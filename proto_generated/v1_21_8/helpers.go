package v1_21_8

import (
	"github.com/admin-else/strom/proto_base"
	"io"
)

func HandshakingToServerPacketIdentifierToType(s string) (t any) {
	switch s {
	case "legacy_server_list_ping":
		t = HandshakingToServerPacketLegacyServerListPing{}
	case "set_protocol":
		t = HandshakingToServerPacketSetProtocol{}
	default:
		t = nil
	}
	return
}
func HandshakingToServerTypeToPacketIdentifier(t any) (s string) {
	switch t.(type) {
	case HandshakingToServerPacketLegacyServerListPing:
		s = "legacy_server_list_ping"
	case HandshakingToServerPacketSetProtocol:
		s = "set_protocol"
	}
	return
}
func HandshakingToClientPacketIdentifierToType(s string) (t any) {
	switch s {
	default:
		t = nil
	}
	return
}
func HandshakingToClientTypeToPacketIdentifier(t any) (s string) {
	switch t.(type) {
	}
	return
}
func StatusToServerPacketIdentifierToType(s string) (t any) {
	switch s {
	case "ping":
		t = StatusToServerPacketPing{}
	case "ping_start":
		t = StatusToServerPacketPingStart{}
	default:
		t = nil
	}
	return
}
func StatusToServerTypeToPacketIdentifier(t any) (s string) {
	switch t.(type) {
	case StatusToServerPacketPing:
		s = "ping"
	case StatusToServerPacketPingStart:
		s = "ping_start"
	}
	return
}
func StatusToClientPacketIdentifierToType(s string) (t any) {
	switch s {
	case "ping":
		t = StatusToClientPacketPing{}
	case "server_info":
		t = StatusToClientPacketServerInfo{}
	default:
		t = nil
	}
	return
}
func StatusToClientTypeToPacketIdentifier(t any) (s string) {
	switch t.(type) {
	case StatusToClientPacketPing:
		s = "ping"
	case StatusToClientPacketServerInfo:
		s = "server_info"
	}
	return
}
func LoginToServerPacketIdentifierToType(s string) (t any) {
	switch s {
	case "cookie_response":
		t = PacketCommonCookieResponse{}
	case "encryption_begin":
		t = LoginToServerPacketEncryptionBegin{}
	case "login_acknowledged":
		t = LoginToServerPacketLoginAcknowledged{}
	case "login_plugin_response":
		t = LoginToServerPacketLoginPluginResponse{}
	case "login_start":
		t = LoginToServerPacketLoginStart{}
	default:
		t = nil
	}
	return
}
func LoginToServerTypeToPacketIdentifier(t any) (s string) {
	switch t.(type) {
	case LoginToServerPacketEncryptionBegin:
		s = "encryption_begin"
	case LoginToServerPacketLoginAcknowledged:
		s = "login_acknowledged"
	case LoginToServerPacketLoginPluginResponse:
		s = "login_plugin_response"
	case LoginToServerPacketLoginStart:
		s = "login_start"
	case PacketCommonCookieResponse:
		s = "cookie_response"
	}
	return
}
func LoginToClientPacketIdentifierToType(s string) (t any) {
	switch s {
	case "compress":
		t = LoginToClientPacketCompress{}
	case "cookie_request":
		t = PacketCommonCookieRequest{}
	case "disconnect":
		t = LoginToClientPacketDisconnect{}
	case "encryption_begin":
		t = LoginToClientPacketEncryptionBegin{}
	case "login_plugin_request":
		t = LoginToClientPacketLoginPluginRequest{}
	case "success":
		t = LoginToClientPacketSuccess{}
	default:
		t = nil
	}
	return
}
func LoginToClientTypeToPacketIdentifier(t any) (s string) {
	switch t.(type) {
	case LoginToClientPacketCompress:
		s = "compress"
	case LoginToClientPacketDisconnect:
		s = "disconnect"
	case LoginToClientPacketEncryptionBegin:
		s = "encryption_begin"
	case LoginToClientPacketLoginPluginRequest:
		s = "login_plugin_request"
	case LoginToClientPacketSuccess:
		s = "success"
	case PacketCommonCookieRequest:
		s = "cookie_request"
	}
	return
}
func ConfigurationToServerPacketIdentifierToType(s string) (t any) {
	switch s {
	case "cookie_response":
		t = PacketCommonCookieResponse{}
	case "custom_click_action":
		t = PacketCommonCustomClickAction{}
	case "custom_payload":
		t = ConfigurationToServerPacketCustomPayload{}
	case "finish_configuration":
		t = ConfigurationToServerPacketFinishConfiguration{}
	case "keep_alive":
		t = ConfigurationToServerPacketKeepAlive{}
	case "pong":
		t = ConfigurationToServerPacketPong{}
	case "resource_pack_receive":
		t = ConfigurationToServerPacketResourcePackReceive{}
	case "select_known_packs":
		t = PacketCommonSelectKnownPacks{}
	case "settings":
		t = PacketCommonSettings{}
	default:
		t = nil
	}
	return
}
func ConfigurationToServerTypeToPacketIdentifier(t any) (s string) {
	switch t.(type) {
	case ConfigurationToServerPacketCustomPayload:
		s = "custom_payload"
	case ConfigurationToServerPacketFinishConfiguration:
		s = "finish_configuration"
	case ConfigurationToServerPacketKeepAlive:
		s = "keep_alive"
	case ConfigurationToServerPacketPong:
		s = "pong"
	case ConfigurationToServerPacketResourcePackReceive:
		s = "resource_pack_receive"
	case PacketCommonCookieResponse:
		s = "cookie_response"
	case PacketCommonCustomClickAction:
		s = "custom_click_action"
	case PacketCommonSelectKnownPacks:
		s = "select_known_packs"
	case PacketCommonSettings:
		s = "settings"
	}
	return
}
func ConfigurationToClientPacketIdentifierToType(s string) (t any) {
	switch s {
	case "add_resource_pack":
		t = PacketCommonAddResourcePack{}
	case "clear_dialog":
		t = PacketCommonClearDialog{}
	case "cookie_request":
		t = PacketCommonCookieRequest{}
	case "custom_payload":
		t = ConfigurationToClientPacketCustomPayload{}
	case "custom_report_details":
		t = PacketCommonCustomReportDetails{}
	case "disconnect":
		t = ConfigurationToClientPacketDisconnect{}
	case "feature_flags":
		t = ConfigurationToClientPacketFeatureFlags{}
	case "finish_configuration":
		t = ConfigurationToClientPacketFinishConfiguration{}
	case "keep_alive":
		t = ConfigurationToClientPacketKeepAlive{}
	case "ping":
		t = ConfigurationToClientPacketPing{}
	case "registry_data":
		t = ConfigurationToClientPacketRegistryData{}
	case "remove_resource_pack":
		t = PacketCommonRemoveResourcePack{}
	case "reset_chat":
		t = ConfigurationToClientPacketResetChat{}
	case "select_known_packs":
		t = PacketCommonSelectKnownPacks{}
	case "server_links":
		t = PacketCommonServerLinks{}
	case "show_dialog":
		t = ConfigurationToClientPacketShowDialog{}
	case "store_cookie":
		t = PacketCommonStoreCookie{}
	case "tags":
		t = ConfigurationToClientPacketTags{}
	case "transfer":
		t = PacketCommonTransfer{}
	default:
		t = nil
	}
	return
}
func ConfigurationToClientTypeToPacketIdentifier(t any) (s string) {
	switch t.(type) {
	case ConfigurationToClientPacketCustomPayload:
		s = "custom_payload"
	case ConfigurationToClientPacketDisconnect:
		s = "disconnect"
	case ConfigurationToClientPacketFeatureFlags:
		s = "feature_flags"
	case ConfigurationToClientPacketFinishConfiguration:
		s = "finish_configuration"
	case ConfigurationToClientPacketKeepAlive:
		s = "keep_alive"
	case ConfigurationToClientPacketPing:
		s = "ping"
	case ConfigurationToClientPacketRegistryData:
		s = "registry_data"
	case ConfigurationToClientPacketResetChat:
		s = "reset_chat"
	case ConfigurationToClientPacketShowDialog:
		s = "show_dialog"
	case ConfigurationToClientPacketTags:
		s = "tags"
	case PacketCommonAddResourcePack:
		s = "add_resource_pack"
	case PacketCommonClearDialog:
		s = "clear_dialog"
	case PacketCommonCookieRequest:
		s = "cookie_request"
	case PacketCommonCustomReportDetails:
		s = "custom_report_details"
	case PacketCommonRemoveResourcePack:
		s = "remove_resource_pack"
	case PacketCommonSelectKnownPacks:
		s = "select_known_packs"
	case PacketCommonServerLinks:
		s = "server_links"
	case PacketCommonStoreCookie:
		s = "store_cookie"
	case PacketCommonTransfer:
		s = "transfer"
	}
	return
}
func PlayToServerPacketIdentifierToType(s string) (t any) {
	switch s {
	case "abilities":
		t = PlayToServerPacketAbilities{}
	case "advancement_tab":
		t = PlayToServerPacketAdvancementTab{}
	case "arm_animation":
		t = PlayToServerPacketArmAnimation{}
	case "block_dig":
		t = PlayToServerPacketBlockDig{}
	case "block_place":
		t = PlayToServerPacketBlockPlace{}
	case "change_gamemode":
		t = PlayToServerPacketChangeGamemode{}
	case "chat_command":
		t = PlayToServerPacketChatCommand{}
	case "chat_command_signed":
		t = PlayToServerPacketChatCommandSigned{}
	case "chat_message":
		t = PlayToServerPacketChatMessage{}
	case "chat_session_update":
		t = PlayToServerPacketChatSessionUpdate{}
	case "chunk_batch_received":
		t = PlayToServerPacketChunkBatchReceived{}
	case "client_command":
		t = PlayToServerPacketClientCommand{}
	case "close_window":
		t = PlayToServerPacketCloseWindow{}
	case "configuration_acknowledged":
		t = PlayToServerPacketConfigurationAcknowledged{}
	case "cookie_response":
		t = PacketCommonCookieResponse{}
	case "craft_recipe_request":
		t = PlayToServerPacketCraftRecipeRequest{}
	case "custom_click_action":
		t = PacketCommonCustomClickAction{}
	case "custom_payload":
		t = PlayToServerPacketCustomPayload{}
	case "debug_sample_subscription":
		t = PlayToServerPacketDebugSampleSubscription{}
	case "displayed_recipe":
		t = PlayToServerPacketDisplayedRecipe{}
	case "edit_book":
		t = PlayToServerPacketEditBook{}
	case "enchant_item":
		t = PlayToServerPacketEnchantItem{}
	case "entity_action":
		t = PlayToServerPacketEntityAction{}
	case "flying":
		t = PlayToServerPacketFlying{}
	case "generate_structure":
		t = PlayToServerPacketGenerateStructure{}
	case "held_item_slot":
		t = PlayToServerPacketHeldItemSlot{}
	case "keep_alive":
		t = PlayToServerPacketKeepAlive{}
	case "lock_difficulty":
		t = PlayToServerPacketLockDifficulty{}
	case "look":
		t = PlayToServerPacketLook{}
	case "message_acknowledgement":
		t = PlayToServerPacketMessageAcknowledgement{}
	case "name_item":
		t = PlayToServerPacketNameItem{}
	case "pick_item_from_block":
		t = PlayToServerPacketPickItemFromBlock{}
	case "pick_item_from_entity":
		t = PlayToServerPacketPickItemFromEntity{}
	case "ping_request":
		t = PlayToServerPacketPingRequest{}
	case "player_input":
		t = PlayToServerPacketPlayerInput{}
	case "player_loaded":
		t = PlayToServerPacketPlayerLoaded{}
	case "pong":
		t = PlayToServerPacketPong{}
	case "position":
		t = PlayToServerPacketPosition{}
	case "position_look":
		t = PlayToServerPacketPositionLook{}
	case "query_block_nbt":
		t = PlayToServerPacketQueryBlockNbt{}
	case "query_entity_nbt":
		t = PlayToServerPacketQueryEntityNbt{}
	case "recipe_book":
		t = PlayToServerPacketRecipeBook{}
	case "resource_pack_receive":
		t = PlayToServerPacketResourcePackReceive{}
	case "select_bundle_item":
		t = PlayToServerPacketSelectBundleItem{}
	case "select_trade":
		t = PlayToServerPacketSelectTrade{}
	case "set_beacon_effect":
		t = PlayToServerPacketSetBeaconEffect{}
	case "set_creative_slot":
		t = PlayToServerPacketSetCreativeSlot{}
	case "set_difficulty":
		t = PlayToServerPacketSetDifficulty{}
	case "set_slot_state":
		t = PlayToServerPacketSetSlotState{}
	case "set_test_block":
		t = PlayToServerPacketSetTestBlock{}
	case "settings":
		t = PacketCommonSettings{}
	case "spectate":
		t = PlayToServerPacketSpectate{}
	case "steer_boat":
		t = PlayToServerPacketSteerBoat{}
	case "tab_complete":
		t = PlayToServerPacketTabComplete{}
	case "teleport_confirm":
		t = PlayToServerPacketTeleportConfirm{}
	case "test_instance_block_action":
		t = PlayToServerPacketTestInstanceBlockAction{}
	case "tick_end":
		t = PlayToServerPacketTickEnd{}
	case "update_command_block":
		t = PlayToServerPacketUpdateCommandBlock{}
	case "update_command_block_minecart":
		t = PlayToServerPacketUpdateCommandBlockMinecart{}
	case "update_jigsaw_block":
		t = PlayToServerPacketUpdateJigsawBlock{}
	case "update_sign":
		t = PlayToServerPacketUpdateSign{}
	case "update_structure_block":
		t = PlayToServerPacketUpdateStructureBlock{}
	case "use_entity":
		t = PlayToServerPacketUseEntity{}
	case "use_item":
		t = PlayToServerPacketUseItem{}
	case "vehicle_move":
		t = PlayToServerPacketVehicleMove{}
	case "window_click":
		t = PlayToServerPacketWindowClick{}
	default:
		t = nil
	}
	return
}
func PlayToServerTypeToPacketIdentifier(t any) (s string) {
	switch t.(type) {
	case PacketCommonCookieResponse:
		s = "cookie_response"
	case PacketCommonCustomClickAction:
		s = "custom_click_action"
	case PacketCommonSettings:
		s = "settings"
	case PlayToServerPacketAbilities:
		s = "abilities"
	case PlayToServerPacketAdvancementTab:
		s = "advancement_tab"
	case PlayToServerPacketArmAnimation:
		s = "arm_animation"
	case PlayToServerPacketBlockDig:
		s = "block_dig"
	case PlayToServerPacketBlockPlace:
		s = "block_place"
	case PlayToServerPacketChangeGamemode:
		s = "change_gamemode"
	case PlayToServerPacketChatCommand:
		s = "chat_command"
	case PlayToServerPacketChatCommandSigned:
		s = "chat_command_signed"
	case PlayToServerPacketChatMessage:
		s = "chat_message"
	case PlayToServerPacketChatSessionUpdate:
		s = "chat_session_update"
	case PlayToServerPacketChunkBatchReceived:
		s = "chunk_batch_received"
	case PlayToServerPacketClientCommand:
		s = "client_command"
	case PlayToServerPacketCloseWindow:
		s = "close_window"
	case PlayToServerPacketConfigurationAcknowledged:
		s = "configuration_acknowledged"
	case PlayToServerPacketCraftRecipeRequest:
		s = "craft_recipe_request"
	case PlayToServerPacketCustomPayload:
		s = "custom_payload"
	case PlayToServerPacketDebugSampleSubscription:
		s = "debug_sample_subscription"
	case PlayToServerPacketDisplayedRecipe:
		s = "displayed_recipe"
	case PlayToServerPacketEditBook:
		s = "edit_book"
	case PlayToServerPacketEnchantItem:
		s = "enchant_item"
	case PlayToServerPacketEntityAction:
		s = "entity_action"
	case PlayToServerPacketFlying:
		s = "flying"
	case PlayToServerPacketGenerateStructure:
		s = "generate_structure"
	case PlayToServerPacketHeldItemSlot:
		s = "held_item_slot"
	case PlayToServerPacketKeepAlive:
		s = "keep_alive"
	case PlayToServerPacketLockDifficulty:
		s = "lock_difficulty"
	case PlayToServerPacketLook:
		s = "look"
	case PlayToServerPacketMessageAcknowledgement:
		s = "message_acknowledgement"
	case PlayToServerPacketNameItem:
		s = "name_item"
	case PlayToServerPacketPickItemFromBlock:
		s = "pick_item_from_block"
	case PlayToServerPacketPickItemFromEntity:
		s = "pick_item_from_entity"
	case PlayToServerPacketPingRequest:
		s = "ping_request"
	case PlayToServerPacketPlayerInput:
		s = "player_input"
	case PlayToServerPacketPlayerLoaded:
		s = "player_loaded"
	case PlayToServerPacketPong:
		s = "pong"
	case PlayToServerPacketPosition:
		s = "position"
	case PlayToServerPacketPositionLook:
		s = "position_look"
	case PlayToServerPacketQueryBlockNbt:
		s = "query_block_nbt"
	case PlayToServerPacketQueryEntityNbt:
		s = "query_entity_nbt"
	case PlayToServerPacketRecipeBook:
		s = "recipe_book"
	case PlayToServerPacketResourcePackReceive:
		s = "resource_pack_receive"
	case PlayToServerPacketSelectBundleItem:
		s = "select_bundle_item"
	case PlayToServerPacketSelectTrade:
		s = "select_trade"
	case PlayToServerPacketSetBeaconEffect:
		s = "set_beacon_effect"
	case PlayToServerPacketSetCreativeSlot:
		s = "set_creative_slot"
	case PlayToServerPacketSetDifficulty:
		s = "set_difficulty"
	case PlayToServerPacketSetSlotState:
		s = "set_slot_state"
	case PlayToServerPacketSetTestBlock:
		s = "set_test_block"
	case PlayToServerPacketSpectate:
		s = "spectate"
	case PlayToServerPacketSteerBoat:
		s = "steer_boat"
	case PlayToServerPacketTabComplete:
		s = "tab_complete"
	case PlayToServerPacketTeleportConfirm:
		s = "teleport_confirm"
	case PlayToServerPacketTestInstanceBlockAction:
		s = "test_instance_block_action"
	case PlayToServerPacketTickEnd:
		s = "tick_end"
	case PlayToServerPacketUpdateCommandBlock:
		s = "update_command_block"
	case PlayToServerPacketUpdateCommandBlockMinecart:
		s = "update_command_block_minecart"
	case PlayToServerPacketUpdateJigsawBlock:
		s = "update_jigsaw_block"
	case PlayToServerPacketUpdateSign:
		s = "update_sign"
	case PlayToServerPacketUpdateStructureBlock:
		s = "update_structure_block"
	case PlayToServerPacketUseEntity:
		s = "use_entity"
	case PlayToServerPacketUseItem:
		s = "use_item"
	case PlayToServerPacketVehicleMove:
		s = "vehicle_move"
	case PlayToServerPacketWindowClick:
		s = "window_click"
	}
	return
}
func PlayToClientPacketIdentifierToType(s string) (t any) {
	switch s {
	case "abilities":
		t = PlayToClientPacketAbilities{}
	case "acknowledge_player_digging":
		t = PlayToClientPacketAcknowledgePlayerDigging{}
	case "action_bar":
		t = PlayToClientPacketActionBar{}
	case "add_resource_pack":
		t = PacketCommonAddResourcePack{}
	case "advancements":
		t = PlayToClientPacketAdvancements{}
	case "animation":
		t = PlayToClientPacketAnimation{}
	case "attach_entity":
		t = PlayToClientPacketAttachEntity{}
	case "block_action":
		t = PlayToClientPacketBlockAction{}
	case "block_break_animation":
		t = PlayToClientPacketBlockBreakAnimation{}
	case "block_change":
		t = PlayToClientPacketBlockChange{}
	case "boss_bar":
		t = PlayToClientPacketBossBar{}
	case "bundle_delimiter":
		t = proto_base.Void{}
	case "camera":
		t = PlayToClientPacketCamera{}
	case "chat_suggestions":
		t = PlayToClientPacketChatSuggestions{}
	case "chunk_batch_finished":
		t = PlayToClientPacketChunkBatchFinished{}
	case "chunk_batch_start":
		t = PlayToClientPacketChunkBatchStart{}
	case "chunk_biomes":
		t = PlayToClientPacketChunkBiomes{}
	case "clear_dialog":
		t = PacketCommonClearDialog{}
	case "clear_titles":
		t = PlayToClientPacketClearTitles{}
	case "close_window":
		t = PlayToClientPacketCloseWindow{}
	case "collect":
		t = PlayToClientPacketCollect{}
	case "cookie_request":
		t = PacketCommonCookieRequest{}
	case "craft_progress_bar":
		t = PlayToClientPacketCraftProgressBar{}
	case "craft_recipe_response":
		t = PlayToClientPacketCraftRecipeResponse{}
	case "custom_payload":
		t = PlayToClientPacketCustomPayload{}
	case "custom_report_details":
		t = PacketCommonCustomReportDetails{}
	case "damage_event":
		t = PlayToClientPacketDamageEvent{}
	case "death_combat_event":
		t = PlayToClientPacketDeathCombatEvent{}
	case "debug_sample":
		t = PlayToClientPacketDebugSample{}
	case "declare_commands":
		t = PlayToClientPacketDeclareCommands{}
	case "declare_recipes":
		t = PlayToClientPacketDeclareRecipes{}
	case "difficulty":
		t = PlayToClientPacketDifficulty{}
	case "end_combat_event":
		t = PlayToClientPacketEndCombatEvent{}
	case "enter_combat_event":
		t = PlayToClientPacketEnterCombatEvent{}
	case "entity_destroy":
		t = PlayToClientPacketEntityDestroy{}
	case "entity_effect":
		t = PlayToClientPacketEntityEffect{}
	case "entity_equipment":
		t = PlayToClientPacketEntityEquipment{}
	case "entity_head_rotation":
		t = PlayToClientPacketEntityHeadRotation{}
	case "entity_look":
		t = PlayToClientPacketEntityLook{}
	case "entity_metadata":
		t = PlayToClientPacketEntityMetadata{}
	case "entity_move_look":
		t = PlayToClientPacketEntityMoveLook{}
	case "entity_sound_effect":
		t = PlayToClientPacketEntitySoundEffect{}
	case "entity_status":
		t = PlayToClientPacketEntityStatus{}
	case "entity_teleport":
		t = PlayToClientPacketEntityTeleport{}
	case "entity_update_attributes":
		t = PlayToClientPacketEntityUpdateAttributes{}
	case "entity_velocity":
		t = PlayToClientPacketEntityVelocity{}
	case "experience":
		t = PlayToClientPacketExperience{}
	case "explosion":
		t = PlayToClientPacketExplosion{}
	case "face_player":
		t = PlayToClientPacketFacePlayer{}
	case "game_state_change":
		t = PlayToClientPacketGameStateChange{}
	case "held_item_slot":
		t = PlayToClientPacketHeldItemSlot{}
	case "hide_message":
		t = PlayToClientPacketHideMessage{}
	case "hurt_animation":
		t = PlayToClientPacketHurtAnimation{}
	case "initialize_world_border":
		t = PlayToClientPacketInitializeWorldBorder{}
	case "keep_alive":
		t = PlayToClientPacketKeepAlive{}
	case "kick_disconnect":
		t = PlayToClientPacketKickDisconnect{}
	case "login":
		t = PlayToClientPacketLogin{}
	case "map":
		t = PlayToClientPacketMap{}
	case "map_chunk":
		t = PlayToClientPacketMapChunk{}
	case "move_minecart":
		t = PlayToClientPacketMoveMinecart{}
	case "multi_block_change":
		t = PlayToClientPacketMultiBlockChange{}
	case "nbt_query_response":
		t = PlayToClientPacketNbtQueryResponse{}
	case "open_book":
		t = PlayToClientPacketOpenBook{}
	case "open_horse_window":
		t = PlayToClientPacketOpenHorseWindow{}
	case "open_sign_entity":
		t = PlayToClientPacketOpenSignEntity{}
	case "open_window":
		t = PlayToClientPacketOpenWindow{}
	case "ping":
		t = PlayToClientPacketPing{}
	case "ping_response":
		t = PlayToClientPacketPingResponse{}
	case "player_chat":
		t = PlayToClientPacketPlayerChat{}
	case "player_info":
		t = PlayToClientPacketPlayerInfo{}
	case "player_remove":
		t = PlayToClientPacketPlayerRemove{}
	case "player_rotation":
		t = PlayToClientPacketPlayerRotation{}
	case "playerlist_header":
		t = PlayToClientPacketPlayerlistHeader{}
	case "position":
		t = PlayToClientPacketPosition{}
	case "profileless_chat":
		t = PlayToClientPacketProfilelessChat{}
	case "recipe_book_add":
		t = PlayToClientPacketRecipeBookAdd{}
	case "recipe_book_remove":
		t = PlayToClientPacketRecipeBookRemove{}
	case "recipe_book_settings":
		t = PlayToClientPacketRecipeBookSettings{}
	case "rel_entity_move":
		t = PlayToClientPacketRelEntityMove{}
	case "remove_entity_effect":
		t = PlayToClientPacketRemoveEntityEffect{}
	case "remove_resource_pack":
		t = PacketCommonRemoveResourcePack{}
	case "reset_score":
		t = PlayToClientPacketResetScore{}
	case "respawn":
		t = PlayToClientPacketRespawn{}
	case "scoreboard_display_objective":
		t = PlayToClientPacketScoreboardDisplayObjective{}
	case "scoreboard_objective":
		t = PlayToClientPacketScoreboardObjective{}
	case "scoreboard_score":
		t = PlayToClientPacketScoreboardScore{}
	case "select_advancement_tab":
		t = PlayToClientPacketSelectAdvancementTab{}
	case "server_data":
		t = PlayToClientPacketServerData{}
	case "server_links":
		t = PacketCommonServerLinks{}
	case "set_cooldown":
		t = PlayToClientPacketSetCooldown{}
	case "set_cursor_item":
		t = PlayToClientPacketSetCursorItem{}
	case "set_passengers":
		t = PlayToClientPacketSetPassengers{}
	case "set_player_inventory":
		t = PlayToClientPacketSetPlayerInventory{}
	case "set_projectile_power":
		t = PlayToClientPacketSetProjectilePower{}
	case "set_slot":
		t = PlayToClientPacketSetSlot{}
	case "set_ticking_state":
		t = PlayToClientPacketSetTickingState{}
	case "set_title_subtitle":
		t = PlayToClientPacketSetTitleSubtitle{}
	case "set_title_text":
		t = PlayToClientPacketSetTitleText{}
	case "set_title_time":
		t = PlayToClientPacketSetTitleTime{}
	case "show_dialog":
		t = PlayToClientPacketShowDialog{}
	case "simulation_distance":
		t = PlayToClientPacketSimulationDistance{}
	case "sound_effect":
		t = PlayToClientPacketSoundEffect{}
	case "spawn_entity":
		t = PlayToClientPacketSpawnEntity{}
	case "spawn_position":
		t = PlayToClientPacketSpawnPosition{}
	case "start_configuration":
		t = PlayToClientPacketStartConfiguration{}
	case "statistics":
		t = PlayToClientPacketStatistics{}
	case "step_tick":
		t = PlayToClientPacketStepTick{}
	case "stop_sound":
		t = PlayToClientPacketStopSound{}
	case "store_cookie":
		t = PacketCommonStoreCookie{}
	case "sync_entity_position":
		t = PlayToClientPacketSyncEntityPosition{}
	case "system_chat":
		t = PlayToClientPacketSystemChat{}
	case "tab_complete":
		t = PlayToClientPacketTabComplete{}
	case "tags":
		t = PlayToClientPacketTags{}
	case "teams":
		t = PlayToClientPacketTeams{}
	case "test_instance_block_status":
		t = PlayToClientPacketTestInstanceBlockStatus{}
	case "tile_entity_data":
		t = PlayToClientPacketTileEntityData{}
	case "tracked_waypoint":
		t = PlayToClientPacketTrackedWaypoint{}
	case "trade_list":
		t = PlayToClientPacketTradeList{}
	case "transfer":
		t = PacketCommonTransfer{}
	case "unload_chunk":
		t = PlayToClientPacketUnloadChunk{}
	case "update_health":
		t = PlayToClientPacketUpdateHealth{}
	case "update_light":
		t = PlayToClientPacketUpdateLight{}
	case "update_time":
		t = PlayToClientPacketUpdateTime{}
	case "update_view_distance":
		t = PlayToClientPacketUpdateViewDistance{}
	case "update_view_position":
		t = PlayToClientPacketUpdateViewPosition{}
	case "vehicle_move":
		t = PlayToClientPacketVehicleMove{}
	case "window_items":
		t = PlayToClientPacketWindowItems{}
	case "world_border_center":
		t = PlayToClientPacketWorldBorderCenter{}
	case "world_border_lerp_size":
		t = PlayToClientPacketWorldBorderLerpSize{}
	case "world_border_size":
		t = PlayToClientPacketWorldBorderSize{}
	case "world_border_warning_delay":
		t = PlayToClientPacketWorldBorderWarningDelay{}
	case "world_border_warning_reach":
		t = PlayToClientPacketWorldBorderWarningReach{}
	case "world_event":
		t = PlayToClientPacketWorldEvent{}
	case "world_particles":
		t = PlayToClientPacketWorldParticles{}
	default:
		t = nil
	}
	return
}
func PlayToClientTypeToPacketIdentifier(t any) (s string) {
	switch t.(type) {
	case PacketCommonAddResourcePack:
		s = "add_resource_pack"
	case PacketCommonClearDialog:
		s = "clear_dialog"
	case PacketCommonCookieRequest:
		s = "cookie_request"
	case PacketCommonCustomReportDetails:
		s = "custom_report_details"
	case PacketCommonRemoveResourcePack:
		s = "remove_resource_pack"
	case PacketCommonServerLinks:
		s = "server_links"
	case PacketCommonStoreCookie:
		s = "store_cookie"
	case PacketCommonTransfer:
		s = "transfer"
	case PlayToClientPacketAbilities:
		s = "abilities"
	case PlayToClientPacketAcknowledgePlayerDigging:
		s = "acknowledge_player_digging"
	case PlayToClientPacketActionBar:
		s = "action_bar"
	case PlayToClientPacketAdvancements:
		s = "advancements"
	case PlayToClientPacketAnimation:
		s = "animation"
	case PlayToClientPacketAttachEntity:
		s = "attach_entity"
	case PlayToClientPacketBlockAction:
		s = "block_action"
	case PlayToClientPacketBlockBreakAnimation:
		s = "block_break_animation"
	case PlayToClientPacketBlockChange:
		s = "block_change"
	case PlayToClientPacketBossBar:
		s = "boss_bar"
	case PlayToClientPacketCamera:
		s = "camera"
	case PlayToClientPacketChatSuggestions:
		s = "chat_suggestions"
	case PlayToClientPacketChunkBatchFinished:
		s = "chunk_batch_finished"
	case PlayToClientPacketChunkBatchStart:
		s = "chunk_batch_start"
	case PlayToClientPacketChunkBiomes:
		s = "chunk_biomes"
	case PlayToClientPacketClearTitles:
		s = "clear_titles"
	case PlayToClientPacketCloseWindow:
		s = "close_window"
	case PlayToClientPacketCollect:
		s = "collect"
	case PlayToClientPacketCraftProgressBar:
		s = "craft_progress_bar"
	case PlayToClientPacketCraftRecipeResponse:
		s = "craft_recipe_response"
	case PlayToClientPacketCustomPayload:
		s = "custom_payload"
	case PlayToClientPacketDamageEvent:
		s = "damage_event"
	case PlayToClientPacketDeathCombatEvent:
		s = "death_combat_event"
	case PlayToClientPacketDebugSample:
		s = "debug_sample"
	case PlayToClientPacketDeclareCommands:
		s = "declare_commands"
	case PlayToClientPacketDeclareRecipes:
		s = "declare_recipes"
	case PlayToClientPacketDifficulty:
		s = "difficulty"
	case PlayToClientPacketEndCombatEvent:
		s = "end_combat_event"
	case PlayToClientPacketEnterCombatEvent:
		s = "enter_combat_event"
	case PlayToClientPacketEntityDestroy:
		s = "entity_destroy"
	case PlayToClientPacketEntityEffect:
		s = "entity_effect"
	case PlayToClientPacketEntityEquipment:
		s = "entity_equipment"
	case PlayToClientPacketEntityHeadRotation:
		s = "entity_head_rotation"
	case PlayToClientPacketEntityLook:
		s = "entity_look"
	case PlayToClientPacketEntityMetadata:
		s = "entity_metadata"
	case PlayToClientPacketEntityMoveLook:
		s = "entity_move_look"
	case PlayToClientPacketEntitySoundEffect:
		s = "entity_sound_effect"
	case PlayToClientPacketEntityStatus:
		s = "entity_status"
	case PlayToClientPacketEntityTeleport:
		s = "entity_teleport"
	case PlayToClientPacketEntityUpdateAttributes:
		s = "entity_update_attributes"
	case PlayToClientPacketEntityVelocity:
		s = "entity_velocity"
	case PlayToClientPacketExperience:
		s = "experience"
	case PlayToClientPacketExplosion:
		s = "explosion"
	case PlayToClientPacketFacePlayer:
		s = "face_player"
	case PlayToClientPacketGameStateChange:
		s = "game_state_change"
	case PlayToClientPacketHeldItemSlot:
		s = "held_item_slot"
	case PlayToClientPacketHideMessage:
		s = "hide_message"
	case PlayToClientPacketHurtAnimation:
		s = "hurt_animation"
	case PlayToClientPacketInitializeWorldBorder:
		s = "initialize_world_border"
	case PlayToClientPacketKeepAlive:
		s = "keep_alive"
	case PlayToClientPacketKickDisconnect:
		s = "kick_disconnect"
	case PlayToClientPacketLogin:
		s = "login"
	case PlayToClientPacketMap:
		s = "map"
	case PlayToClientPacketMapChunk:
		s = "map_chunk"
	case PlayToClientPacketMoveMinecart:
		s = "move_minecart"
	case PlayToClientPacketMultiBlockChange:
		s = "multi_block_change"
	case PlayToClientPacketNbtQueryResponse:
		s = "nbt_query_response"
	case PlayToClientPacketOpenBook:
		s = "open_book"
	case PlayToClientPacketOpenHorseWindow:
		s = "open_horse_window"
	case PlayToClientPacketOpenSignEntity:
		s = "open_sign_entity"
	case PlayToClientPacketOpenWindow:
		s = "open_window"
	case PlayToClientPacketPing:
		s = "ping"
	case PlayToClientPacketPingResponse:
		s = "ping_response"
	case PlayToClientPacketPlayerChat:
		s = "player_chat"
	case PlayToClientPacketPlayerInfo:
		s = "player_info"
	case PlayToClientPacketPlayerRemove:
		s = "player_remove"
	case PlayToClientPacketPlayerRotation:
		s = "player_rotation"
	case PlayToClientPacketPlayerlistHeader:
		s = "playerlist_header"
	case PlayToClientPacketPosition:
		s = "position"
	case PlayToClientPacketProfilelessChat:
		s = "profileless_chat"
	case PlayToClientPacketRecipeBookAdd:
		s = "recipe_book_add"
	case PlayToClientPacketRecipeBookRemove:
		s = "recipe_book_remove"
	case PlayToClientPacketRecipeBookSettings:
		s = "recipe_book_settings"
	case PlayToClientPacketRelEntityMove:
		s = "rel_entity_move"
	case PlayToClientPacketRemoveEntityEffect:
		s = "remove_entity_effect"
	case PlayToClientPacketResetScore:
		s = "reset_score"
	case PlayToClientPacketRespawn:
		s = "respawn"
	case PlayToClientPacketScoreboardDisplayObjective:
		s = "scoreboard_display_objective"
	case PlayToClientPacketScoreboardObjective:
		s = "scoreboard_objective"
	case PlayToClientPacketScoreboardScore:
		s = "scoreboard_score"
	case PlayToClientPacketSelectAdvancementTab:
		s = "select_advancement_tab"
	case PlayToClientPacketServerData:
		s = "server_data"
	case PlayToClientPacketSetCooldown:
		s = "set_cooldown"
	case PlayToClientPacketSetCursorItem:
		s = "set_cursor_item"
	case PlayToClientPacketSetPassengers:
		s = "set_passengers"
	case PlayToClientPacketSetPlayerInventory:
		s = "set_player_inventory"
	case PlayToClientPacketSetProjectilePower:
		s = "set_projectile_power"
	case PlayToClientPacketSetSlot:
		s = "set_slot"
	case PlayToClientPacketSetTickingState:
		s = "set_ticking_state"
	case PlayToClientPacketSetTitleSubtitle:
		s = "set_title_subtitle"
	case PlayToClientPacketSetTitleText:
		s = "set_title_text"
	case PlayToClientPacketSetTitleTime:
		s = "set_title_time"
	case PlayToClientPacketShowDialog:
		s = "show_dialog"
	case PlayToClientPacketSimulationDistance:
		s = "simulation_distance"
	case PlayToClientPacketSoundEffect:
		s = "sound_effect"
	case PlayToClientPacketSpawnEntity:
		s = "spawn_entity"
	case PlayToClientPacketSpawnPosition:
		s = "spawn_position"
	case PlayToClientPacketStartConfiguration:
		s = "start_configuration"
	case PlayToClientPacketStatistics:
		s = "statistics"
	case PlayToClientPacketStepTick:
		s = "step_tick"
	case PlayToClientPacketStopSound:
		s = "stop_sound"
	case PlayToClientPacketSyncEntityPosition:
		s = "sync_entity_position"
	case PlayToClientPacketSystemChat:
		s = "system_chat"
	case PlayToClientPacketTabComplete:
		s = "tab_complete"
	case PlayToClientPacketTags:
		s = "tags"
	case PlayToClientPacketTeams:
		s = "teams"
	case PlayToClientPacketTestInstanceBlockStatus:
		s = "test_instance_block_status"
	case PlayToClientPacketTileEntityData:
		s = "tile_entity_data"
	case PlayToClientPacketTrackedWaypoint:
		s = "tracked_waypoint"
	case PlayToClientPacketTradeList:
		s = "trade_list"
	case PlayToClientPacketUnloadChunk:
		s = "unload_chunk"
	case PlayToClientPacketUpdateHealth:
		s = "update_health"
	case PlayToClientPacketUpdateLight:
		s = "update_light"
	case PlayToClientPacketUpdateTime:
		s = "update_time"
	case PlayToClientPacketUpdateViewDistance:
		s = "update_view_distance"
	case PlayToClientPacketUpdateViewPosition:
		s = "update_view_position"
	case PlayToClientPacketVehicleMove:
		s = "vehicle_move"
	case PlayToClientPacketWindowItems:
		s = "window_items"
	case PlayToClientPacketWorldBorderCenter:
		s = "world_border_center"
	case PlayToClientPacketWorldBorderLerpSize:
		s = "world_border_lerp_size"
	case PlayToClientPacketWorldBorderSize:
		s = "world_border_size"
	case PlayToClientPacketWorldBorderWarningDelay:
		s = "world_border_warning_delay"
	case PlayToClientPacketWorldBorderWarningReach:
		s = "world_border_warning_reach"
	case PlayToClientPacketWorldEvent:
		s = "world_event"
	case PlayToClientPacketWorldParticles:
		s = "world_particles"
	case proto_base.Void:
		s = "bundle_delimiter"
	}
	return
}

func TypeToPacketIdentifier(d proto_base.Direction, s proto_base.State, t any) (i string) {
	switch d {
	case proto_base.ToServer:
		switch s {
		case proto_base.Handshaking:
			i = HandshakingToServerTypeToPacketIdentifier(t)
		case proto_base.Status:
			i = StatusToServerTypeToPacketIdentifier(t)
		case proto_base.Login:
			i = LoginToServerTypeToPacketIdentifier(t)
		case proto_base.Configuration:
			i = ConfigurationToServerTypeToPacketIdentifier(t)
		case proto_base.Play:
			i = PlayToServerTypeToPacketIdentifier(t)
		}
	case proto_base.ToClient:
		switch s {
		case proto_base.Handshaking:
			i = HandshakingToClientTypeToPacketIdentifier(t)
		case proto_base.Status:
			i = StatusToClientTypeToPacketIdentifier(t)
		case proto_base.Login:
			i = LoginToClientTypeToPacketIdentifier(t)
		case proto_base.Configuration:
			i = ConfigurationToClientTypeToPacketIdentifier(t)
		case proto_base.Play:
			i = PlayToClientTypeToPacketIdentifier(t)
		}
	}
	return
}

func PacketIdentifierToType(d proto_base.Direction, s proto_base.State, i string) (t any) {
	switch d {
	case proto_base.ToServer:
		switch s {
		case proto_base.Handshaking:
			t = HandshakingToServerPacketIdentifierToType(i)
		case proto_base.Status:
			t = StatusToServerPacketIdentifierToType(i)
		case proto_base.Login:
			t = LoginToServerPacketIdentifierToType(i)
		case proto_base.Configuration:
			t = ConfigurationToServerPacketIdentifierToType(i)
		case proto_base.Play:
			t = PlayToServerPacketIdentifierToType(i)
		}
	case proto_base.ToClient:
		switch s {
		case proto_base.Handshaking:
			t = HandshakingToClientPacketIdentifierToType(i)
		case proto_base.Status:
			t = StatusToClientPacketIdentifierToType(i)
		case proto_base.Login:
			t = LoginToClientPacketIdentifierToType(i)
		case proto_base.Configuration:
			t = ConfigurationToClientPacketIdentifierToType(i)
		case proto_base.Play:
			t = PlayToClientPacketIdentifierToType(i)
		}
	}
	return
}

func DecodePacket(d proto_base.Direction, s proto_base.State, r io.Reader) (params any, err error) {
	switch d {
	case proto_base.ToServer:
		switch s {
		case proto_base.Handshaking:
			var packet HandshakingToServerPacket
			packet, err = HandshakingToServerPacket{}.Decode(r)
			if err != nil {
				return
			}
			params = packet.Params
		case proto_base.Status:
			var packet StatusToServerPacket
			packet, err = StatusToServerPacket{}.Decode(r)
			if err != nil {
				return
			}
			params = packet.Params
		case proto_base.Login:
			var packet LoginToServerPacket
			packet, err = LoginToServerPacket{}.Decode(r)
			if err != nil {
				return
			}
			params = packet.Params
		case proto_base.Configuration:
			var packet ConfigurationToServerPacket
			packet, err = ConfigurationToServerPacket{}.Decode(r)
			if err != nil {
				return
			}
			params = packet.Params
		case proto_base.Play:
			var packet PlayToServerPacket
			packet, err = packet.Decode(r)
			if err != nil {
				return
			}
			params = packet.Params
		}
	case proto_base.ToClient:
		switch s {
		case proto_base.Handshaking:
			var packet HandshakingToClientPacket
			packet, err = HandshakingToClientPacket{}.Decode(r)
			if err != nil {
				return
			}
			params = packet.Params
		case proto_base.Status:
			var packet StatusToClientPacket
			packet, err = StatusToClientPacket{}.Decode(r)
			if err != nil {
				return
			}
			params = packet.Params
		case proto_base.Login:
			var packet LoginToClientPacket
			packet, err = LoginToClientPacket{}.Decode(r)
			if err != nil {
				return
			}
			params = packet.Params
		case proto_base.Configuration:
			var packet ConfigurationToClientPacket
			packet, err = ConfigurationToClientPacket{}.Decode(r)
			if err != nil {
				return
			}
			params = packet.Params
		case proto_base.Play:
			var packet PlayToClientPacket
			packet, err = PlayToClientPacket{}.Decode(r)
			if err != nil {
				return
			}
			params = packet.Params
		}
	}
	return
}

func EncodePacket(d proto_base.Direction, s proto_base.State, i string, p any, w io.Writer) (err error) {
	switch d {
	case proto_base.ToServer:
		switch s {
		case proto_base.Handshaking:
			err = HandshakingToServerPacket{i, p}.Encode(w)
		case proto_base.Status:
			err = StatusToServerPacket{i, p}.Encode(w)
		case proto_base.Login:
			err = LoginToServerPacket{i, p}.Encode(w)
		case proto_base.Configuration:
			err = ConfigurationToServerPacket{i, p}.Encode(w)
		case proto_base.Play:
			err = PlayToServerPacket{i, p}.Encode(w)
		}
	case proto_base.ToClient:
		switch s {
		case proto_base.Handshaking:
			err = HandshakingToClientPacket{i, p}.Encode(w)
		case proto_base.Status:
			err = StatusToClientPacket{i, p}.Encode(w)
		case proto_base.Login:
			err = LoginToClientPacket{i, p}.Encode(w)
		case proto_base.Configuration:
			err = ConfigurationToClientPacket{i, p}.Encode(w)
		case proto_base.Play:
			err = PlayToClientPacket{i, p}.Encode(w)
		}
	}
	return
}
