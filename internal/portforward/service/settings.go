package service

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gosettings"
)

type Settings struct {
	UserSettings  settings.PortForwarding
	PortForwarder PortForwarder
	Interface     string // needed for PIA and ProtonVPN, tun0 for example
	ServerName    string // needed for PIA
	VPNProvider   string // used to validate new settings
}

// UpdateWith deep copies the receiving settings, overrides the copy with
// fields set in the partialUpdate argument, validates the new settings
// and returns them if they are valid, or returns an error otherwise.
// In all cases, the receiving settings are unmodified.
func (s Settings) UpdateWith(partialUpdate Settings) (updatedSettings Settings, err error) {
	updatedSettings = s.copy()
	updatedSettings.overrideWith(partialUpdate)
	err = updatedSettings.validate()
	if err != nil {
		return updatedSettings, fmt.Errorf("validating new settings: %w", err)
	}
	return updatedSettings, nil
}

func (s Settings) copy() (copied Settings) {
	copied.UserSettings = s.UserSettings.Copy()
	copied.PortForwarder = s.PortForwarder
	copied.Interface = s.Interface
	copied.ServerName = s.ServerName
	copied.VPNProvider = s.VPNProvider
	return copied
}

func (s *Settings) overrideWith(update Settings) {
	s.UserSettings.OverrideWith(update.UserSettings)
	s.PortForwarder = gosettings.OverrideWithInterface(s.PortForwarder, update.PortForwarder)
	s.Interface = gosettings.OverrideWithString(s.Interface, update.Interface)
	s.ServerName = gosettings.OverrideWithString(s.ServerName, update.ServerName)
	s.VPNProvider = gosettings.OverrideWithString(s.VPNProvider, update.VPNProvider)
}

var (
	ErrVPNProviderNotSet   = errors.New("VPN provider not set")
	ErrServerNameNotSet    = errors.New("server name not set")
	ErrPortForwarderNotSet = errors.New("port forwarder not set")
	ErrGatewayNotSet       = errors.New("gateway not set")
	ErrInterfaceNotSet     = errors.New("interface not set")
)

func (s *Settings) validate() (err error) {
	switch {
	case s.VPNProvider == "":
		return fmt.Errorf("%w", ErrVPNProviderNotSet)
	case s.VPNProvider == providers.PrivateInternetAccess && s.ServerName == "":
		return fmt.Errorf("%w", ErrServerNameNotSet)
	case s.PortForwarder == nil:
		return fmt.Errorf("%w", ErrPortForwarderNotSet)
	case s.Interface == "":
		return fmt.Errorf("%w", ErrInterfaceNotSet)
	}

	return s.UserSettings.Validate(s.VPNProvider)
}
