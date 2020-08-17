package types

// staking module event types
const (
	EventTypeSendCoinsFromModuleToAccount       = "SendCoinsFromModuleToAccount"
	EventTypeSendCoinsFromModuleToModule        = "SendCoinsFromModuleToModule"
	EventTypeSendCoinsFromAccountToModule       = "SendCoinsFromAccountToModule"
	EventTypeDelegateCoinsFromAccountToModule   = "DelegateCoinsFromAccountToModule"
	EventTypeUndelegateCoinsFromModuleToAccount = "UndelegateCoinsFromModuleToAccount"
	EventTypeModuleMintCoins                    = "ModuleMintCoins"
	EventTypeModuleBurnCoins                    = "ModuleBurnCoins"
	EventTypeInitModuleAccount                  = "initModuleAccount"

	AttributeKeyFrom    = "from"
	AttributeKeyTo      = "to"
	AttributeKeyAmount  = "amount"
	AttributeKeyAccount = "account"
	AttributeKeyAuth    = "auth"
	AttributeKeyCreator = "creator"
)
