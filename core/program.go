package core

import (
	"errors"
	"io"

	. "github.com/ioeXNetwork/ioeX.Utility/common"
)

type Program struct {
	//the contract program code,which will be run on VM or specific envrionment
	Code []byte

	//the program code's parameter
	Parameter []byte
}

func (self Program) String() string {
	return "Program: {\n\t\t" +
		"Code: " + BytesToHexString(self.Code) + "\n\t\t" +
		"Parameter: " + BytesToHexString(self.Parameter) + "\n\t\t" +
		"}"
}

//Serialize the Program
func (p *Program) Serialize(w io.Writer) error {
	if err := WriteVarBytes(w, p.Parameter); err != nil {
		return errors.New("Execute Program Serialize Parameter failed.")
	}

	if err := WriteVarBytes(w, p.Code); err != nil {
		return errors.New("Execute Program Serialize Code failed.")
	}

	return nil
}

//Deserialize the Program
func (p *Program) Deserialize(w io.Reader) error {
	parameter, err := ReadVarBytes(w)
	if err != nil {
		return errors.New("Execute Program Deserialize Parameter failed.")
	}
	p.Parameter = parameter

	code, err := ReadVarBytes(w)
	if err != nil {
		return errors.New("Execute Program Deserialize Code failed.")
	}
	p.Code = code

	return nil
}
