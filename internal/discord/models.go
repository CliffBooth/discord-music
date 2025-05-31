package discord

import "time"

// generated here https://transform.tools/json-to-go
// TODO: change all interface{} to types
// TODO: take care of all repeating structs (like User)

// event resposne types
const (
	CHANNEL_MESSAGE_WITH_SOURCE = 4
)

type BaseEvent struct {
	T  string `json:"t,omitempty"`
	S  int    `json:"s,omitempty"`
	Op int    `json:"op,omitempty"`
}

type ReadyEvent struct {
	BaseEvent
	D struct {
		V            int `json:"v"`
		UserSettings struct {
		} `json:"user_settings"`
		User struct {
			Verified      bool        `json:"verified"`
			Username      string      `json:"username"`
			PrimaryGuild  interface{} `json:"primary_guild"`
			MfaEnabled    bool        `json:"mfa_enabled"`
			ID            string      `json:"id"`
			GlobalName    interface{} `json:"global_name"`
			Flags         int         `json:"flags"`
			Email         interface{} `json:"email"`
			Discriminator string      `json:"discriminator"`
			Clan          interface{} `json:"clan"`
			Bot           bool        `json:"bot"`
			Avatar        string      `json:"avatar"`
		} `json:"user"`
		SessionType      string        `json:"session_type"`
		SessionID        string        `json:"session_id"`
		ResumeGatewayURL string        `json:"resume_gateway_url"`
		Relationships    []interface{} `json:"relationships"`
		PrivateChannels  []interface{} `json:"private_channels"`
		Presences        []interface{} `json:"presences"`
		Guilds           []struct {
			Unavailable bool   `json:"unavailable"`
			ID          string `json:"id"`
		} `json:"guilds"`
		GuildJoinRequests    []interface{} `json:"guild_join_requests"`
		GeoOrderedRtcRegions []string      `json:"geo_ordered_rtc_regions"`
		GameRelationships    []interface{} `json:"game_relationships"`
		Auth                 struct {
		} `json:"auth"`
		Application struct {
			ID    string `json:"id"`
			Flags int    `json:"flags"`
		} `json:"application"`
		Trace []string `json:"_trace"`
	} `json:"d"`
}

type GuildCreateEvent struct {
	BaseEvent
	D struct {
		LatestOnboardingQuestionID interface{}   `json:"latest_onboarding_question_id"`
		Region                     string        `json:"region"`
		ActivityInstances          []interface{} `json:"activity_instances"`
		Profile                    interface{}   `json:"profile"`
		DiscoverySplash            interface{}   `json:"discovery_splash"`
		MaxVideoChannelUsers       int           `json:"max_video_channel_users"`
		RulesChannelID             interface{}   `json:"rules_channel_id"`
		PremiumProgressBarEnabled  bool          `json:"premium_progress_bar_enabled"`
		Presences                  []interface{} `json:"presences"`
		MfaLevel                   int           `json:"mfa_level"`
		SystemChannelFlags         int           `json:"system_channel_flags"`
		SoundboardSounds           []interface{} `json:"soundboard_sounds"`
		StageInstances             []interface{} `json:"stage_instances"`
		Banner                     interface{}   `json:"banner"`
		Nsfw                       bool          `json:"nsfw"`
		Stickers                   []interface{} `json:"stickers"`
		ID                         string        `json:"id"`
		VanityURLCode              interface{}   `json:"vanity_url_code"`
		ApplicationCommandCounts   struct {
		} `json:"application_command_counts"`
		Members []struct {
			User struct {
				Username             string      `json:"username"`
				PublicFlags          int         `json:"public_flags"`
				PrimaryGuild         interface{} `json:"primary_guild"`
				ID                   string      `json:"id"`
				GlobalName           interface{} `json:"global_name"`
				DisplayName          interface{} `json:"display_name"`
				Discriminator        string      `json:"discriminator"`
				Clan                 interface{} `json:"clan"`
				Bot                  bool        `json:"bot"`
				AvatarDecorationData interface{} `json:"avatar_decoration_data"`
				Avatar               string      `json:"avatar"`
			} `json:"user"`
			Roles                      []string    `json:"roles"`
			PremiumSince               interface{} `json:"premium_since"`
			Pending                    bool        `json:"pending"`
			Nick                       interface{} `json:"nick"`
			Mute                       bool        `json:"mute"`
			JoinedAt                   time.Time   `json:"joined_at"`
			Flags                      int         `json:"flags"`
			Deaf                       bool        `json:"deaf"`
			CommunicationDisabledUntil interface{} `json:"communication_disabled_until"`
			Banner                     interface{} `json:"banner"`
			Avatar                     interface{} `json:"avatar"`
		} `json:"members"`
		Name       string      `json:"name"`
		HomeHeader interface{} `json:"home_header"`
		AfkTimeout int         `json:"afk_timeout"`
		HubType    interface{} `json:"hub_type"`
		Roles      []struct {
			Version      int64       `json:"version"`
			UnicodeEmoji interface{} `json:"unicode_emoji"`
			Tags         struct {
			} `json:"tags"`
			Position    int         `json:"position"`
			Permissions int         `json:"permissions"`
			Name        string      `json:"name"`
			Mentionable bool        `json:"mentionable"`
			Managed     bool        `json:"managed"`
			ID          string      `json:"id"`
			Icon        interface{} `json:"icon"`
			Hoist       bool        `json:"hoist"`
			Flags       int         `json:"flags"`
			Color       int         `json:"color"`
		} `json:"roles"`
		Features                    []interface{} `json:"features"`
		PreferredLocale             string        `json:"preferred_locale"`
		AfkChannelID                interface{}   `json:"afk_channel_id"`
		MaxStageVideoChannelUsers   int           `json:"max_stage_video_channel_users"`
		Threads                     []interface{} `json:"threads"`
		Emojis                      []interface{} `json:"emojis"`
		PublicUpdatesChannelID      interface{}   `json:"public_updates_channel_id"`
		Splash                      interface{}   `json:"splash"`
		InventorySettings           interface{}   `json:"inventory_settings"`
		Description                 interface{}   `json:"description"`
		Lazy                        bool          `json:"lazy"`
		DefaultMessageNotifications int           `json:"default_message_notifications"`
		IncidentsData               interface{}   `json:"incidents_data"`
		EmbeddedActivities          []interface{} `json:"embedded_activities"`
		ExplicitContentFilter       int           `json:"explicit_content_filter"`
		VoiceStates                 []interface{} `json:"voice_states"`
		ApplicationID               interface{}   `json:"application_id"`
		PremiumSubscriptionCount    int           `json:"premium_subscription_count"`
		MaxMembers                  int           `json:"max_members"`
		PremiumTier                 int           `json:"premium_tier"`
		OwnerID                     string        `json:"owner_id"`
		Unavailable                 bool          `json:"unavailable"`
		JoinedAt                    time.Time     `json:"joined_at"`
		Large                       bool          `json:"large"`
		Clan                        interface{}   `json:"clan"`
		NsfwLevel                   int           `json:"nsfw_level"`
		Icon                        interface{}   `json:"icon"`
		Version                     int64         `json:"version"`
		SystemChannelID             string        `json:"system_channel_id"`
		Channels                    []struct {
			Version              int64         `json:"version"`
			Type                 int           `json:"type"`
			Position             int           `json:"position"`
			PermissionOverwrites []interface{} `json:"permission_overwrites"`
			Name                 string        `json:"name"`
			ID                   string        `json:"id"`
			Flags                int           `json:"flags"`
			Topic                interface{}   `json:"topic,omitempty"`
			RateLimitPerUser     int           `json:"rate_limit_per_user,omitempty"`
			ParentID             string        `json:"parent_id,omitempty"`
			LastMessageID        string        `json:"last_message_id,omitempty"`
			IconEmoji            struct {
				Name string      `json:"name"`
				ID   interface{} `json:"id"`
			} `json:"icon_emoji,omitempty"`
			UserLimit int         `json:"user_limit,omitempty"`
			RtcRegion interface{} `json:"rtc_region,omitempty"`
			Bitrate   int         `json:"bitrate,omitempty"`
		} `json:"channels"`
		MemberCount           int           `json:"member_count"`
		GuildScheduledEvents  []interface{} `json:"guild_scheduled_events"`
		VerificationLevel     int           `json:"verification_level"`
		SafetyAlertsChannelID interface{}   `json:"safety_alerts_channel_id"`
	} `json:"d"`
}

type InteractionCreateEvent struct {
	BaseEvent
	D struct {
		Version int    `json:"version"`
		Type    int    `json:"type"`
		Token   string `json:"token"`
		Member  struct {
			User struct {
				Username             string      `json:"username"`
				PublicFlags          int         `json:"public_flags"`
				PrimaryGuild         interface{} `json:"primary_guild"`
				ID                   string      `json:"id"`
				GlobalName           string      `json:"global_name"`
				Discriminator        string      `json:"discriminator"`
				Clan                 interface{} `json:"clan"`
				AvatarDecorationData interface{} `json:"avatar_decoration_data"`
				Avatar               string      `json:"avatar"`
			} `json:"user"`
			UnusualDmActivityUntil     interface{}   `json:"unusual_dm_activity_until"`
			Roles                      []interface{} `json:"roles"`
			PremiumSince               interface{}   `json:"premium_since"`
			Permissions                string        `json:"permissions"`
			Pending                    bool          `json:"pending"`
			Nick                       interface{}   `json:"nick"`
			Mute                       bool          `json:"mute"`
			JoinedAt                   time.Time     `json:"joined_at"`
			Flags                      int           `json:"flags"`
			Deaf                       bool          `json:"deaf"`
			CommunicationDisabledUntil interface{}   `json:"communication_disabled_until"`
			Banner                     interface{}   `json:"banner"`
			Avatar                     interface{}   `json:"avatar"`
		} `json:"member"`
		Locale      string `json:"locale"`
		ID          string `json:"id"`
		GuildLocale string `json:"guild_locale"`
		GuildID     string `json:"guild_id"`
		Guild       struct {
			Locale   string        `json:"locale"`
			ID       string        `json:"id"`
			Features []interface{} `json:"features"`
		} `json:"guild"`
		Entitlements      []interface{} `json:"entitlements"`
		EntitlementSkuIds []interface{} `json:"entitlement_sku_ids"`
		Data              struct {
			Type int    `json:"type"`
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"data"`
		Context   int    `json:"context"`
		ChannelID string `json:"channel_id"`
		Channel   struct {
			Type             int         `json:"type"`
			Topic            interface{} `json:"topic"`
			ThemeColor       interface{} `json:"theme_color"`
			RateLimitPerUser int         `json:"rate_limit_per_user"`
			Position         int         `json:"position"`
			Permissions      string      `json:"permissions"`
			ParentID         string      `json:"parent_id"`
			Nsfw             bool        `json:"nsfw"`
			Name             string      `json:"name"`
			LastMessageID    string      `json:"last_message_id"`
			ID               string      `json:"id"`
			IconEmoji        struct {
				Name string      `json:"name"`
				ID   interface{} `json:"id"`
			} `json:"icon_emoji"`
			GuildID string `json:"guild_id"`
			Flags   int    `json:"flags"`
		} `json:"channel"`
		AuthorizingIntegrationOwners struct {
			Num0 string `json:"0"`
		} `json:"authorizing_integration_owners"`
		ApplicationID  string `json:"application_id"`
		AppPermissions string `json:"app_permissions"`
	} `json:"d"`
}

//TODO: fill the field types
type InteractionResponse struct {
	Type int                 `json:"type"`
	Data InteractionRespData `json:"data,omitempty"`
}

type InteractionRespData struct {
	Tts             bool   `json:"tts,omitempty"`
	Content         string `json:"content,omitempty"`
	Embeds          interface{}
	AllowedMentions interface{}
	Flags           int `json:"falgs,omitempty"`
	Components      interface{}
	Attachments     interface{}
	Poll            interface{}
}

type UserVoiceStateResponse struct {
	ChannelID string `json:"channel_id"`
}

type VoiceStateUpdateRequest struct {
	OP int                  `json:"op"`
	D  VoiceStateUpdateData `json:"d"`
}

type VoiceStateUpdateData struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	SelfMute  bool   `json:"self_mute"`
	SelfDeaf  bool   `json:"slef_deaf"`
}

type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	//TODO: int64?
	Type              int   `json:"type"`
	Integration_types []int `json:"integration_types"`
	Contexts          []int `json:"contexts"`

	F CommandFunc `json:"-"`
}
