// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package builtin_test

import (
	. "gopkg.in/check.v1"

	"github.com/snapcore/snapd/interfaces"
	"github.com/snapcore/snapd/interfaces/builtin"
	"github.com/snapcore/snapd/snap"
	"github.com/snapcore/snapd/testutil"
)

type UDisks2InterfaceSuite struct {
	iface interfaces.Interface
	slot  *interfaces.Slot
	plug  *interfaces.Plug
}

var _ = Suite(&UDisks2InterfaceSuite{
	iface: &builtin.UDisks2Interface{},
	slot: &interfaces.Slot{
		SlotInfo: &snap.SlotInfo{
			Snap: &snap.Info{
				SuggestedName: "udisks2",
				SideInfo:      snap.SideInfo{Developer: "canonical"}},
			Name:      "udisks2",
			Interface: "udisks2",
		},
	},
	plug: &interfaces.Plug{
		PlugInfo: &snap.PlugInfo{
			Snap:      &snap.Info{SuggestedName: "udisks2"},
			Name:      "udisks2",
			Interface: "udisks2",
		},
	},
})

func (s *UDisks2InterfaceSuite) TestName(c *C) {
	c.Assert(s.iface.Name(), Equals, "udisks2")
}

func (s *UDisks2InterfaceSuite) TestSanitizeSlot(c *C) {
	err := s.iface.SanitizeSlot(s.slot)
	c.Assert(err, IsNil)
}

func (s *UDisks2InterfaceSuite) TestSanitizeSlotNotUdisks2FromCanonical(c *C) {
	err := s.iface.SanitizeSlot(&interfaces.Slot{SlotInfo: &snap.SlotInfo{
		Snap: &snap.Info{
			SuggestedName: "notudisks2",
			SideInfo:      snap.SideInfo{Developer: "canonical"},
		},
		Name:      "udisks2",
		Interface: "udisks2",
	}})
	c.Assert(err, ErrorMatches, "udisks2 slot reserved \\(snap name 'notudisks2' != 'udisks2'\\)")
}

func (s *UDisks2InterfaceSuite) TestSanitizeSlotUdisks2NotFromCanonical(c *C) {
	err := s.iface.SanitizeSlot(&interfaces.Slot{SlotInfo: &snap.SlotInfo{
		Snap: &snap.Info{
			SuggestedName: "udisks2",
			SideInfo:      snap.SideInfo{Developer: "foo"},
		},
		Name:      "udisks2",
		Interface: "udisks2",
	}})
	c.Assert(err, ErrorMatches, "udisks2 slot is reserved for Canonical")
}

func (s *UDisks2InterfaceSuite) TestSanitizeSlotNotUdisks2NotFromCanonical(c *C) {
	err := s.iface.SanitizeSlot(&interfaces.Slot{SlotInfo: &snap.SlotInfo{
		Snap: &snap.Info{
			SuggestedName: "notudisks2",
			SideInfo:      snap.SideInfo{Developer: "foo"},
		},
		Name:      "udisks2",
		Interface: "udisks2",
	}})
	c.Assert(err, ErrorMatches, "udisks2 slot reserved \\(snap name 'notudisks2' != 'udisks2'\\)")
}

func (s *UDisks2InterfaceSuite) TestSanitizeSlotUdisks2Sideload(c *C) {
	err := s.iface.SanitizeSlot(&interfaces.Slot{SlotInfo: &snap.SlotInfo{
		Snap: &snap.Info{
			SuggestedName: "udisks2",
			SideInfo:      snap.SideInfo{Developer: ""},
		},
		Name:      "udisks2",
		Interface: "udisks2",
	}})
	c.Assert(err, ErrorMatches, "udisks2 slot is reserved for Canonical")
}

// The label glob when all apps are bound to the udisks2 slot
func (s *UDisks2InterfaceSuite) TestConnectedPlugSnippetUsesSlotLabelAll(c *C) {
	app1 := &snap.AppInfo{Name: "app1"}
	app2 := &snap.AppInfo{Name: "app2"}
	slot := &interfaces.Slot{
		SlotInfo: &snap.SlotInfo{
			Snap: &snap.Info{
				SuggestedName: "udisks2",
				Apps:          map[string]*snap.AppInfo{"app1": app1, "app2": app2},
			},
			Name:      "udisks2",
			Interface: "udisks2",
			Apps:      map[string]*snap.AppInfo{"app1": app1, "app2": app2},
		},
	}
	snippet, err := s.iface.ConnectedPlugSnippet(s.plug, slot, interfaces.SecurityAppArmor)
	c.Assert(err, IsNil)
	c.Assert(string(snippet), testutil.Contains, `peer=(label="snap.udisks2.*"),`)
}

// The label uses alternation when some, but not all, apps is bound to the udisks2 slot
func (s *UDisks2InterfaceSuite) TestConnectedPlugSnippetUsesSlotLabelSome(c *C) {
	app1 := &snap.AppInfo{Name: "app1"}
	app2 := &snap.AppInfo{Name: "app2"}
	app3 := &snap.AppInfo{Name: "app3"}
	slot := &interfaces.Slot{
		SlotInfo: &snap.SlotInfo{
			Snap: &snap.Info{
				SuggestedName: "udisks2",
				Apps:          map[string]*snap.AppInfo{"app1": app1, "app2": app2, "app3": app3},
			},
			Name:      "udisks2",
			Interface: "udisks2",
			Apps:      map[string]*snap.AppInfo{"app1": app1, "app2": app2},
		},
	}
	snippet, err := s.iface.ConnectedPlugSnippet(s.plug, slot, interfaces.SecurityAppArmor)
	c.Assert(err, IsNil)
	c.Assert(string(snippet), testutil.Contains, `peer=(label="snap.udisks2.{app1,app2}"),`)
}

// The label uses short form when exactly one app is bound to the udisks2 slot
func (s *UDisks2InterfaceSuite) TestConnectedPlugSnippetUsesSlotLabelOne(c *C) {
	app := &snap.AppInfo{Name: "app"}
	slot := &interfaces.Slot{
		SlotInfo: &snap.SlotInfo{
			Snap: &snap.Info{
				SuggestedName: "udisks2",
				Apps:          map[string]*snap.AppInfo{"app": app},
			},
			Name:      "udisks2",
			Interface: "udisks2",
			Apps:      map[string]*snap.AppInfo{"app": app},
		},
	}
	snippet, err := s.iface.ConnectedPlugSnippet(s.plug, slot, interfaces.SecurityAppArmor)
	c.Assert(err, IsNil)
	c.Assert(string(snippet), testutil.Contains, `peer=(label="snap.udisks2.app"),`)
}

// The label glob when all apps are bound to the udisks2 plug
func (s *UDisks2InterfaceSuite) TestConnectedSlotSnippetUsesPlugLabelAll(c *C) {
	app1 := &snap.AppInfo{Name: "app1"}
	app2 := &snap.AppInfo{Name: "app2"}
	plug := &interfaces.Plug{
		PlugInfo: &snap.PlugInfo{
			Snap: &snap.Info{
				SuggestedName: "udisks2",
				Apps:          map[string]*snap.AppInfo{"app1": app1, "app2": app2},
			},
			Name:      "udisks2",
			Interface: "udisks2",
			Apps:      map[string]*snap.AppInfo{"app1": app1, "app2": app2},
		},
	}
	snippet, err := s.iface.ConnectedSlotSnippet(plug, s.slot, interfaces.SecurityAppArmor)
	c.Assert(err, IsNil)
	c.Assert(string(snippet), testutil.Contains, `peer=(label="snap.udisks2.*"),`)
}

// The label uses alternation when some, but not all, apps is bound to the udisks2 plug
func (s *UDisks2InterfaceSuite) TestConnectedSlotSnippetUsesPlugLabelSome(c *C) {
	app1 := &snap.AppInfo{Name: "app1"}
	app2 := &snap.AppInfo{Name: "app2"}
	app3 := &snap.AppInfo{Name: "app3"}
	plug := &interfaces.Plug{
		PlugInfo: &snap.PlugInfo{
			Snap: &snap.Info{
				SuggestedName: "udisks2",
				Apps:          map[string]*snap.AppInfo{"app1": app1, "app2": app2, "app3": app3},
			},
			Name:      "udisks2",
			Interface: "udisks2",
			Apps:      map[string]*snap.AppInfo{"app1": app1, "app2": app2},
		},
	}
	snippet, err := s.iface.ConnectedSlotSnippet(plug, s.slot, interfaces.SecurityAppArmor)
	c.Assert(err, IsNil)
	c.Assert(string(snippet), testutil.Contains, `peer=(label="snap.udisks2.{app1,app2}"),`)
}

// The label uses short form when exactly one app is bound to the udisks2 plug
func (s *UDisks2InterfaceSuite) TestConnectedSlotSnippetUsesPlugLabelOne(c *C) {
	app := &snap.AppInfo{Name: "app"}
	plug := &interfaces.Plug{
		PlugInfo: &snap.PlugInfo{
			Snap: &snap.Info{
				SuggestedName: "udisks2",
				Apps:          map[string]*snap.AppInfo{"app": app},
			},
			Name:      "udisks2",
			Interface: "udisks2",
			Apps:      map[string]*snap.AppInfo{"app": app},
		},
	}
	snippet, err := s.iface.ConnectedSlotSnippet(plug, s.slot, interfaces.SecurityAppArmor)
	c.Assert(err, IsNil)
	c.Assert(string(snippet), testutil.Contains, `peer=(label="snap.udisks2.app"),`)
}

func (s *UDisks2InterfaceSuite) TestUnusedSecuritySystems(c *C) {
	ppSystems := [...]interfaces.SecuritySystem{interfaces.SecuritySecComp,
		interfaces.SecurityDBus, interfaces.SecurityUDev, interfaces.SecurityMount,
		interfaces.SecurityAppArmor}
	for _, system := range ppSystems {
		snippet, err := s.iface.PermanentPlugSnippet(s.plug, system)
		c.Assert(err, IsNil)
		c.Assert(snippet, IsNil)
	}

	csSystems := [...]interfaces.SecuritySystem{interfaces.SecuritySecComp,
		interfaces.SecurityDBus, interfaces.SecurityUDev, interfaces.SecurityMount}
	for _, system := range csSystems {
		snippet, err := s.iface.ConnectedSlotSnippet(s.plug, s.slot, system)
		c.Assert(err, IsNil)
		c.Assert(snippet, IsNil)
	}

	cpSystems := [...]interfaces.SecuritySystem{interfaces.SecurityUDev,
		interfaces.SecurityMount}
	for _, system := range cpSystems {
		snippet, err := s.iface.ConnectedPlugSnippet(s.plug, s.slot, system)
		c.Assert(err, IsNil)
		c.Assert(snippet, IsNil)
	}

	snippet, err := s.iface.PermanentSlotSnippet(s.slot, interfaces.SecurityMount)
	c.Assert(err, IsNil)
	c.Assert(snippet, IsNil)
}

func (s *UDisks2InterfaceSuite) TestUsedSecuritySystems(c *C) {
	systems := [...]interfaces.SecuritySystem{interfaces.SecurityAppArmor,
		interfaces.SecuritySecComp, interfaces.SecurityDBus}
	for _, system := range systems {
		snippet, err := s.iface.ConnectedPlugSnippet(s.plug, s.slot, system)
		c.Assert(err, IsNil)
		c.Assert(snippet, Not(IsNil))
		snippet, err = s.iface.PermanentSlotSnippet(s.slot, system)
		c.Assert(err, IsNil)
		c.Assert(snippet, Not(IsNil))
	}
	snippet, err := s.iface.ConnectedSlotSnippet(s.plug, s.slot, interfaces.SecurityAppArmor)
	c.Assert(err, IsNil)
	c.Assert(snippet, Not(IsNil))
	snippet, err = s.iface.PermanentSlotSnippet(s.slot, interfaces.SecurityUDev)
	c.Assert(err, IsNil)
	c.Assert(snippet, Not(IsNil))
}

func (s *UDisks2InterfaceSuite) TestUnexpectedSecuritySystems(c *C) {
	snippet, err := s.iface.PermanentPlugSnippet(s.plug, "foo")
	c.Assert(err, Equals, interfaces.ErrUnknownSecurity)
	c.Assert(snippet, IsNil)
	snippet, err = s.iface.ConnectedPlugSnippet(s.plug, s.slot, "foo")
	c.Assert(err, Equals, interfaces.ErrUnknownSecurity)
	c.Assert(snippet, IsNil)
	snippet, err = s.iface.PermanentSlotSnippet(s.slot, "foo")
	c.Assert(err, Equals, interfaces.ErrUnknownSecurity)
	c.Assert(snippet, IsNil)
	snippet, err = s.iface.ConnectedSlotSnippet(s.plug, s.slot, "foo")
	c.Assert(err, Equals, interfaces.ErrUnknownSecurity)
	c.Assert(snippet, IsNil)
}

func (s *UDisks2InterfaceSuite) TestAutoConnect(c *C) {
	c.Check(s.iface.AutoConnect(), Equals, false)
}
