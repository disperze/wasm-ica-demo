package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
)

type CapabilityKeeper interface {
	GetCapability(ctx sdk.Context, name string) (*capabilitytypes.Capability, bool)
	ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error
	AuthenticateCapability(ctx sdk.Context, capability *capabilitytypes.Capability, name string) bool
}
