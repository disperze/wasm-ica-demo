package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/interchain-accounts/x/inter-tx/types"
)

type WasmCapabitilyKeeper struct {
	capabilityKeeper types.CapabilityKeeper
}

func NewWasmCapabitilyKeeper(keeper types.CapabilityKeeper) *WasmCapabitilyKeeper {
	return &WasmCapabitilyKeeper{
		capabilityKeeper: keeper,
	}
}

func (k *WasmCapabitilyKeeper) GetCapability(ctx sdk.Context, name string) (*capabilitytypes.Capability, bool) {
	return k.capabilityKeeper.GetCapability(ctx, name)
}

func (k *WasmCapabitilyKeeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	// ica auth calls ClaimCapability
	if ctx.Context().Value("ica") != nil {
		return nil
	}

	return k.capabilityKeeper.ClaimCapability(ctx, cap, name)
}

func (k *WasmCapabitilyKeeper) AuthenticateCapability(ctx sdk.Context, capability *capabilitytypes.Capability, name string) bool {
	return k.capabilityKeeper.AuthenticateCapability(ctx, capability, name)
}
