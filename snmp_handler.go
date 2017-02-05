package main

import (
	"sort"

	agentx "github.com/posteo/go-agentx"
	"github.com/posteo/go-agentx/pdu"
	"github.com/posteo/go-agentx/value"
)

type SNMPHandler struct {
	oids  OIDSorter
	items map[string]*agentx.ListItem
}

// Add adds a list item
func (self *SNMPHandler) Add(oid value.OID, oid_type pdu.VariableType, value interface{}) {
	log.Infof("Adding oid %v", oid)
	if self.oids == nil {
		self.items = make(map[string]*agentx.ListItem)
		self.oids = OIDSorter{}
	}

	item := &agentx.ListItem{Type: oid_type, Value: value}
	if _, ok := self.items[oid.String()]; ok {
		self.items[oid.String()] = item
	} else {
		self.oids = append(self.oids, oid)
		sort.Sort(self.oids)
		self.items[oid.String()] = item
	}
}

// Remove removes a list item
func (self *SNMPHandler) Remove(oid value.OID) {
	log.Infof("Removing oid %v", oid)
	if self.oids == nil {
		return
	}

	delete(self.items, oid.String())
	for i, curr_oid := range self.oids {
		if oid.String() == curr_oid.String() {
			self.oids.Swap(i, len(self.oids)-1)
			self.oids = self.oids[:len(self.oids)-2]
		}
	}
}

// RemoveAll removes all oids from the list
func (self *SNMPHandler) RemoveAll() {
	log.Info("Removing all oids")
	self.items = make(map[string]*agentx.ListItem)
	self.oids = OIDSorter{}
}

// Get tries to find the provided oid and returns the corresponding value.
func (self *SNMPHandler) Get(oid value.OID) (value.OID, pdu.VariableType, interface{}, error) {
	log.Debugf("SNMP GET %v", oid)
	return self.doGet(oid)
}

// GetNext tries to find the value that follows the provided oid and returns it.
func (self *SNMPHandler) GetNext(from value.OID, includeFrom bool, to value.OID) (value.OID, pdu.VariableType, interface{}, error) {
	log.Debugf("SNMP GETNEXT %v", from)
	if self.oids == nil {
		return nil, pdu.VariableTypeNoSuchObject, nil, nil
	}

	for _, oid := range self.oids {
		if OIDLessThan(oid, from) {
			continue
		} else if OIDGreaterThan(oid, from) {
			return self.doGet(oid)
		} else {
			// Not less than and not greater than means equal
			if includeFrom {
				return self.doGet(oid)
			}
		}
	}
	return nil, pdu.VariableTypeNoSuchObject, nil, nil
}

// doGet retrieves the value
func (self *SNMPHandler) doGet(oid value.OID) (value.OID, pdu.VariableType, interface{}, error) {
	if self.items == nil {
		return nil, pdu.VariableTypeNoSuchObject, nil, nil
	}

	item, ok := self.items[oid.String()]
	if ok {
		return oid, item.Type, item.Value, nil
	}
	return nil, pdu.VariableTypeNoSuchObject, nil, nil
}
