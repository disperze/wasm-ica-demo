package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
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

type WasmChannelKeeper struct {
	channelKeeper types.ChannelKeeper
}

func NewWasmChannelKeeper(keeper types.ChannelKeeper) *WasmChannelKeeper {
	return &WasmChannelKeeper{
		channelKeeper: keeper,
	}
}

func (k *WasmChannelKeeper) GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool) {
	if ctx.Context().Value("ica") != nil {
		// TODO: check port start with ".wasm" and contract is ica-Account
		srcPort = strings.Replace(srcPort, "wasm.", icatypes.PortPrefix, 1)
	}
	return k.channelKeeper.GetChannel(ctx, srcPort, srcChan)
}

func (k *WasmChannelKeeper) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	return k.channelKeeper.GetNextSequenceSend(ctx, portID, channelID)
}

func (k *WasmChannelKeeper) SendPacket(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet ibcexported.PacketI) error {
	return k.channelKeeper.SendPacket(ctx, channelCap, packet)
}

func (k *WasmChannelKeeper) ChanCloseInit(ctx sdk.Context, portID, channelID string, chanCap *capabilitytypes.Capability) error {
	return k.channelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

func (k *WasmChannelKeeper) GetAllChannels(ctx sdk.Context) (channels []channeltypes.IdentifiedChannel) {
	return k.channelKeeper.GetAllChannels(ctx)
}

func (k *WasmChannelKeeper) IterateChannels(ctx sdk.Context, cb func(channeltypes.IdentifiedChannel) bool) {
	k.channelKeeper.IterateChannels(ctx, cb)
}

func (k *WasmChannelKeeper) SetChannel(ctx sdk.Context, portID, channelID string, channel channeltypes.Channel) {
	if ctx.Context().Value("ica") != nil {
		portID = strings.Replace(portID, "wasm.", icatypes.PortPrefix, 1)
	}
	k.channelKeeper.SetChannel(ctx, portID, channelID, channel)
}
