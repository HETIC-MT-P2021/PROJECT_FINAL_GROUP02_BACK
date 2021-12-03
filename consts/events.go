package consts

import "github.com/SteakBarbare/RPGBot/game"

var DamageCategoryId = 0

var PermanentDamageTrapEvent = game.Event{
	EventModelId: 0,
	CategoryId: DamageCategoryId,
	Name: "Permanent Trap",
	Description: "The room is poisoned! You lose 1 HP",
	IsAlwaysActive:	true,
	WasActivated: false,
}

var OneTimeDamageTrapEvent = game.Event{
	EventModelId: 1,
	CategoryId: DamageCategoryId,
	Name: "Trap",
	Description: "It's a trap! You lose 3 HP",
	IsAlwaysActive:	false,
	WasActivated: false,
}

var DamageCategory = game.Category{
	Id: DamageCategoryId,
	Name: "Damage",
	Events: []game.Event{PermanentDamageTrapEvent, OneTimeDamageTrapEvent},
}

var HealCategoryId = 1

var SmallHealEvent = game.Event{
	EventModelId: 2,
	CategoryId: HealCategoryId,
	Name: "Small heal",
	Description: "You pass by a weird bush and grab a few berries ! You recover 1 HP",
	IsAlwaysActive:	false,
	WasActivated: false,
}

var MediumHealEvent = game.Event{
	EventModelId: 3,
	CategoryId: HealCategoryId,
	Name: "Medium heal",
	Description: "Lucky find, a potion ! You recover 3 HP",
	IsAlwaysActive:	false,
	WasActivated: false,
}

var HealCategory = game.Category{
	Id: HealCategoryId,
	Name: "Heal",
	Events: []game.Event{SmallHealEvent, MediumHealEvent},
}

var BuffCategoryId = 2

var BuffPrecisionTriggerEvent = game.Event{
	EventModelId: 4,
	CategoryId: BuffCategoryId,
	Name: "Precision Buff Trigger",
	Description: "A curious looking fairy grants you her blessing, You gain +15 precision in this tile and close ones",
	IsAlwaysActive:	false,
	WasActivated: false,
}

var BuffPrecisionEffectEvent = game.Event{
	EventModelId: 5,
	CategoryId: BuffCategoryId,
	Name: "Precision Buff Effect",
	Description: "The blessing of the Fairy, You gain +15 precision here",
	IsAlwaysActive:	true,
	WasActivated: false,
}

var BuffCategory = game.Category{
	Id: BuffCategoryId,
	Name: "Buff",
	Events: []game.Event{BuffPrecisionTriggerEvent, BuffPrecisionEffectEvent},
}

var DebuffCategoryId = 3

var DebuffStrengthTriggerEvent = game.Event{
	EventModelId: 6,
	CategoryId: DebuffCategoryId,
	Name: "Strength Debuff Trigger",
	Description: "In a corner, you find a pile of dead crows, you feel weaker, You lose -15 strength in this tile and close ones",
	IsAlwaysActive:	false,
	WasActivated: false,
}

var DebuffStrengthEffectEvent = game.Event{
	EventModelId: 7,
	CategoryId: DebuffCategoryId,
	Name: "Strength Debuff Effect",
	Description: "The curse of the crows, You lose -15 strength here",
	IsAlwaysActive:	true,
	WasActivated: false,
}

var DebuffCategory = game.Category{
	Id: DebuffCategoryId,
	Name: "Debuff",
	Events: []game.Event{DebuffStrengthTriggerEvent, DebuffStrengthEffectEvent},
}

var EventModels = []game.Event{
	PermanentDamageTrapEvent,
	OneTimeDamageTrapEvent,
	SmallHealEvent,
	MediumHealEvent,
	BuffPrecisionTriggerEvent,
	DebuffStrengthTriggerEvent,
}

var Categories = []*game.Category{
	&DamageCategory,
	&HealCategory,
	&BuffCategory,
	&DebuffCategory,
}
