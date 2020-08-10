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

	AttributeKeyFrom   = "from"
	AttributeKeyTo     = "to"
	AttributeKeyAmount = "amount"
)
