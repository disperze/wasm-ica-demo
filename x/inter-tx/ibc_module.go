package inter_tx

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
	"github.com/cosmos/interchain-accounts/x/inter-tx/keeper"
)

var _ porttypes.IBCModule = IBCModule{}
var icaVal string = "ica"

// IBCModule implements the ICS26 interface for interchain accounts controller chains
type IBCModule struct {
	app    porttypes.IBCModule
	keeper keeper.Keeper
}

// NewIBCModule creates a new IBCModule given the keeper
func NewIBCModule(app porttypes.IBCModule, k keeper.Keeper) IBCModule {
	return IBCModule{
		app:    app,
		keeper: k,
	}
}

// OnChanOpenInit implements the IBCModule interface
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) error {
	// id for ica calls for skip wasm.ClaimCapability
	icaCtx := ctx.WithContext(context.WithValue(ctx.Context(), "ica", &icaVal))
	wasmPortID := strings.Replace(portID, icatypes.PortKeyPrefix, "wasm.", 1)
	if err := im.app.OnChanOpenInit(icaCtx, order, connectionHops, wasmPortID, channelID, chanCap, counterparty, version); err != nil {
		return err
	}

	// Claim channel capability passed back by IBC module
	if err := im.keeper.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return err
	}

	return nil
}

// OnChanOpenTry implements the IBCModule interface
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	// never called
	return "", nil
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	// wasm call SetChannelID
	icaCtx := ctx.WithContext(context.WithValue(ctx.Context(), "ica", &icaVal))
	wasmPortID := strings.Replace(portID, icatypes.PortKeyPrefix, "wasm.", 1)
	return im.app.OnChanOpenAck(icaCtx, wasmPortID, channelID, counterpartyChannelID, counterpartyVersion)
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// never called
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// never called
	return nil
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// never called
	return nil
}

// OnRecvPacket implements the IBCModule interface. A successful acknowledgement
// is returned if the packet data is succesfully decoded and the receive application
// logic returns without error.
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	// never called
	return channeltypes.NewErrorAcknowledgement("cannot receive packet via interchain accounts authentication module")
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	icaCtx := ctx.WithContext(context.WithValue(ctx.Context(), "ica", &icaVal))
	packet.SourcePort = strings.Replace(packet.SourcePort, icatypes.PortKeyPrefix, "wasm.", 1)
	return im.app.OnAcknowledgementPacket(icaCtx, packet, acknowledgement, relayer)
}

// OnTimeoutPacket implements the IBCModule interface.
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	icaCtx := ctx.WithContext(context.WithValue(ctx.Context(), "ica", &icaVal))
	packet.SourcePort = strings.Replace(packet.SourcePort, icatypes.PortKeyPrefix, "wasm.", 1)
	return im.app.OnTimeoutPacket(icaCtx, packet, relayer)
}
