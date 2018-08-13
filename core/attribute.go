package core

import (
	"errors"
	"io"

	. "github.com/ioeX/ioeX.Utility/common"
)

type AttributeUsage byte

const (
	Nonce          AttributeUsage = 0x00
	Script         AttributeUsage = 0x20
	Memo           AttributeUsage = 0x81
	Description    AttributeUsage = 0x90
	DescriptionUrl AttributeUsage = 0x91
	Confirmations  AttributeUsage = 0x92
)

func (it AttributeUsage) Name() string {
	switch it {
	case Nonce:
		return "Nonce"
	case Script:
		return "Script"
	case Memo:
		return "Memo"
	case Description:
		return "Description"
	case DescriptionUrl:
		return "DescriptionUrl"
	case Confirmations:
		return "Confirmations"
	default:
		return "Unknown"
	}
}

func IsValidAttributeType(usage AttributeUsage) bool {
	switch usage {
	case Nonce, Script, Memo, Description, DescriptionUrl, Confirmations:
		return true
	}
	return false
}

type Attribute struct {
	Usage AttributeUsage
	Data  []byte
}

func (attr Attribute) String() string {
	return "Attribute: {\n\t\t" +
		"Usage: " + attr.Usage.Name() + "\n\t\t" +
		"Data: " + BytesToHexString(attr.Data) + "\n\t\t" +
		"}"
}

func NewAttribute(u AttributeUsage, d []byte) Attribute {
	return Attribute{Usage: u, Data: d}
}

func (attr *Attribute) Serialize(w io.Writer) error {
	if err := WriteUint8(w, byte(attr.Usage)); err != nil {
		return errors.New("Transaction attribute Usage serialization error.")
	}
	if !IsValidAttributeType(attr.Usage) {
		return errors.New("[Attribute error] Unsupported attribute Description.")
	}
	if err := WriteVarBytes(w, attr.Data); err != nil {
		return errors.New("Transaction attribute Data serialization error.")
	}
	return nil
}

func (attr *Attribute) Deserialize(r io.Reader) error {
	val, err := ReadBytes(r, 1)
	if err != nil {
		return errors.New("Transaction attribute Usage deserialization error.")
	}
	attr.Usage = AttributeUsage(val[0])
	if !IsValidAttributeType(attr.Usage) {
		return errors.New("[Attribute error] Unsupported attribute Description.")
	}
	attr.Data, err = ReadVarBytes(r)
	if err != nil {
		return errors.New("Transaction attribute Data deserialization error.")
	}
	return nil
}
