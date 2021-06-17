package validate

import "github.com/go-playground/validator/v10"

type UserAuth struct {
	Code string `validate:"required"`
}

type GetForksBody struct {
	Warehouse string `validate:"min=0,max=4"`
	Category  string `validate:"min=0,max=4"`
}

type InsertForksBody struct {
	Warehouse string `validate:"required,min=4,max=4"`
	Category  string `validate:"required,min=4,max=4"`
	No        string `validate:"required,min=3,max=5"`
}

type InsertWarehouseBody struct {
	Name   string `validate:"required"`
	Number string `validate:"required,min=4,max=4"`
}

type InsertForkCatBody struct {
	Name   string `validate:"required"`
	Number string `validate:"required,min=3,max=4"`
}

type InsertBatteryBody struct {
	Number      string `validate:"required"`
	Forkliftcat string `validate:"required,min=4,max=4"`
	Warehouse   string `validate:"required,min=4,max=4"`
}

type SwitchBatBody struct {
	No        string `validate:"required"`
	ForkNo    string `validate:"required,min=4,max=4"`
	Warehouse string `validate:"required,min=4,max=4"`
	UserID    string `validate:"required"`
	Switch    string `validate:"required"`
}

func SwichBatValidator(i *SwitchBatBody) error {
	validate := validator.New()
	err := validate.Struct(i)
	return err
}

func UserAuthValidator(i *UserAuth) error {
	validate := validator.New()
	err := validate.Struct(i)
	return err
}

func GetForksBodyValidator(i *GetForksBody) error {
	validate := validator.New()
	err := validate.Struct(i)
	return err
}

func InsertForksBodyValidator(i *InsertForksBody) error {
	validate := validator.New()
	err := validate.Struct(i)
	return err
}

func InsertWarehouseValidator(i *InsertWarehouseBody) error {
	validate := validator.New()
	err := validate.Struct(i)
	return err

}

func InsertForkCatValidator(i *InsertForkCatBody) error {
	validate := validator.New()
	err := validate.Struct(i)
	return err
}

func InsertBatteryValidator(i *InsertBatteryBody) error {
	validate := validator.New()
	err := validate.Struct(i)
	return err
}
