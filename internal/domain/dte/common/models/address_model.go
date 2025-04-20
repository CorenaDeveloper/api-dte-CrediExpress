package models

import "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/location"

// Address es una estructura que representa un Department, Municipality y Complement de un DTE
type Address struct {
	Department   location.Department   `json:"department"`
	Municipality location.Municipality `json:"municipality"`
	Complement   location.Address      `json:"complement"`
}

func (a *Address) GetDepartment() string {
	return a.Department.GetValue()
}

func (a *Address) GetMunicipality() string {
	return a.Municipality.GetValue()
}

func (a *Address) GetComplement() string {
	return a.Complement.GetValue()
}

func (a *Address) SetDepartment(department string) error {
	deptObj, err := location.NewDepartment(department)
	if err != nil {
		return err
	}
	a.Department = *deptObj
	return nil
}

func (a *Address) SetMunicipality(municipality string) error {
	munObj, err := location.NewMunicipality(municipality, a.Department)
	if err != nil {
		return err
	}
	a.Municipality = *munObj
	return nil
}

func (a *Address) SetComplement(complement string) error {
	compObj, err := location.NewAddress(complement)
	if err != nil {
		return err
	}
	a.Complement = *compObj
	return nil
}
