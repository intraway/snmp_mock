package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/posteo/go-agentx/pdu"
	"github.com/posteo/go-agentx/value"
)

func LoadOids(snmp_handler *SNMPHandler, oid_files ...string) error {
	log.Info("Loading oids...")
	for _, oid_file := range oid_files {
		log.Info("Loading oids from file", oid_file)
		file_data, err := ioutil.ReadFile(oid_file)
		if err != nil {
			return err
		}

		err = loadOids(snmp_handler, string(file_data))
		if err != nil {
			return err
		}
	}
	return nil
}

func loadOids(snmp_handler *SNMPHandler, data string) error {
	r := csv.NewReader(strings.NewReader(data))
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	for _, oid_data := range records {
		if len(oid_data) != 3 {
			return fmt.Errorf("Expected 3 elements in line (oid;type;value). Got: %v", len(oid_data))
		}

		oid, err := value.ParseOID(oid_data[0])
		if err != nil {
			return err
		}

		oid_type, err := stringToVariableType(oid_data[1])
		if err != nil {
			return err
		}

		value, err := convertVariable(oid_data[2], oid_type)
		if err != nil {
			return err
		}

		snmp_handler.Add(oid, oid_type, value)
	}

	return nil
}

func stringToVariableType(s string) (pdu.VariableType, error) {
	switch s {
	case "Integer":
		return pdu.VariableTypeInteger, nil
	case "OctetString":
		return pdu.VariableTypeOctetString, nil
	case "Null":
		return pdu.VariableTypeNull, nil
	case "ObjectIdentifier":
		return pdu.VariableTypeObjectIdentifier, nil
	case "IPAddress":
		return pdu.VariableTypeIPAddress, nil
	case "Counter32":
		return pdu.VariableTypeCounter32, nil
	case "Gauge32":
		return pdu.VariableTypeGauge32, nil
	case "TimeTicks":
		return pdu.VariableTypeTimeTicks, nil
	//case "Opaque":
	//	return pdu.VariableTypeOpaque
	case "Counter64":
		return pdu.VariableTypeCounter64, nil
	default:
		return pdu.VariableTypeNoSuchObject, fmt.Errorf("Unknown type '%v'", s)
	}
}

func convertVariable(val string, oid_type pdu.VariableType) (interface{}, error) {
	switch oid_type {
	case pdu.VariableTypeInteger:
		if i, err := strconv.ParseInt(val, 10, 32); err != nil {
			return nil, err
		} else {
			return int32(i), nil
		}
	case pdu.VariableTypeOctetString:
		return val, nil
	case pdu.VariableTypeObjectIdentifier:
		_, err := value.ParseOID(val)
		return val, err
	case pdu.VariableTypeIPAddress:
		ip := net.ParseIP(val).To4()
		if ip == nil {
			return nil, fmt.Errorf("Wrong IP format (%v)", val)
		} else {
			return ip, nil
		}
	case pdu.VariableTypeCounter32, pdu.VariableTypeGauge32:
		if i, err := strconv.ParseUint(val, 10, 32); err != nil {
			return nil, err
		} else {
			return uint32(i), nil
		}
	case pdu.VariableTypeTimeTicks:
		if d, err := strconv.ParseUint(val, 10, 32); err != nil {
			return nil, err
		} else {
			return time.Duration(d) * time.Millisecond * 10, nil
		}
	//case VariableTypeOpaque:
	case pdu.VariableTypeCounter64:
		if i, err := strconv.ParseUint(val, 10, 64); err != nil {
			return nil, err
		} else {
			return uint64(i), nil
		}
	default:
		return nil, fmt.Errorf("Unknown type %v", oid_type)
	}
}

/*
func (v *Variable) UnmarshalBinary(data []byte) error {
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.LittleEndian, &v.Type); err != nil {
		return errgo.Mask(err)
	}
	offset := 4

	if err := v.Name.UnmarshalBinary(data[offset:]); err != nil {
		return errgo.Mask(err)
	}
	offset += v.Name.ByteSize()

	switch v.Type {
	case VariableTypeInteger:
		value := int32(0)
		if err := binary.Read(buffer, binary.LittleEndian, &value); err != nil {
			return errgo.Mask(err)
		}
		v.Value = value
	case VariableTypeOctetString:
		octetString := &OctetString{}
		if err := octetString.UnmarshalBinary(data[offset:]); err != nil {
			return errgo.Mask(err)
		}
		v.Value = octetString.Text
	case VariableTypeNull, VariableTypeNoSuchObject, VariableTypeNoSuchInstance, VariableTypeEndOfMIBView:
		v.Value = nil
	case VariableTypeObjectIdentifier:
		oid := &ObjectIdentifier{}
		if err := oid.UnmarshalBinary(data[offset:]); err != nil {
			return errgo.Mask(err)
		}
		v.Value = oid.GetIdentifier()
	case VariableTypeIPAddress:
		octetString := &OctetString{}
		if err := octetString.UnmarshalBinary(data[offset:]); err != nil {
			return errgo.Mask(err)
		}
		v.Value = net.IP(octetString.Text)
	case VariableTypeCounter32, VariableTypeGauge32:
		value := uint32(0)
		if err := binary.Read(buffer, binary.LittleEndian, &value); err != nil {
			return errgo.Mask(err)
		}
		v.Value = value
	case VariableTypeTimeTicks:
		value := uint32(0)
		if err := binary.Read(buffer, binary.LittleEndian, &value); err != nil {
			return errgo.Mask(err)
		}
		v.Value = time.Duration(value) * time.Second / 100
	case VariableTypeOpaque:
		octetString := &OctetString{}
		if err := octetString.UnmarshalBinary(data[offset:]); err != nil {
			return errgo.Mask(err)
		}
		v.Value = []byte(octetString.Text)
	case VariableTypeCounter64:
		value := uint64(0)
		if err := binary.Read(buffer, binary.LittleEndian, &value); err != nil {
			return errgo.Mask(err)
		}
		v.Value = value
	default:
		return errgo.Newf("unhandled variable type %s", v.Type)
	}

	return nil
}
*/
