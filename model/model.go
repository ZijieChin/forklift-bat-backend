package model

import (
	"forklift-bat-backend/database"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID     string `gorm:"unique"`
	UserName   string
	Status     bool
	Group      string
	LastLogin  string
	LoginCount string
}

type Group struct {
	GroupID       string `gorm:"unique"`
	GroupDescribe string
	Authority     string
}

type Forklift struct {
	gorm.Model
	Number    string
	Category  string
	Warehouse string
	FullNo    string `gorm:"unique"`
}

type ForkCat struct {
	gorm.Model
	Name   string `gorm:"unique"`
	Number string
}

type Warehouse struct {
	gorm.Model
	Number string `gorm:"unique"`
	Name   string
}

type Battery struct {
	gorm.Model
	Number     string
	Forklift   string // forklift type
	Warehouse  string
	FullNo     string `gorm:"unique"` // warehouse + forklift type + battery SN
	Status     string // on for busy, off for idle
	LastSeenAt string // forklift number last time battery is on
}

type SwitchBattery struct {
	gorm.Model
	BatNumber  string
	ForkNumber string //forklift SN
	UserID     string
	Operation  string // on for busy, off for idle
	FullNo     string
}

var db *gorm.DB

func SyncDatabase() {
	db = database.InitDatabase()
	db.AutoMigrate(&Forklift{}, &Warehouse{}, &Battery{}, &ForkCat{}, &User{}, &Group{}, &SwitchBattery{})
}

func (User) GetUser(userid string) (user []User) {
	db.Where(&User{UserID: userid}).Find(&user)
	return
}

func (Forklift) GetForkliftNo(dc, cat string) (forklifts []Forklift) {
	db.Where(&Forklift{Warehouse: dc, Category: cat}).Find(&forklifts)
	return
}

func (Warehouse) GetWarehouse(dc string) (warehouse []Warehouse) {
	db.Where(&Warehouse{Number: dc}).First(&warehouse)
	return
}

func (ForkCat) GetForkCat(no string) (forkcat []ForkCat) {
	db.Where(&ForkCat{Number: no}).First(&forkcat)
	return
}

func (Battery) GetBatStatus(batNo string, forklift string, dc string) (battery []Battery) {
	db.Where(&Battery{Number: batNo, Forklift: forklift, Warehouse: dc}).First(&battery)
	return
}

func (SwitchBattery) GetBatSwitch(fullNo string) (switchbattery []SwitchBattery) {
	db.Where(&SwitchBattery{FullNo: fullNo}).Last(&switchbattery)
	return
}

func InsertForks(insertbody *Forklift) error {
	result := db.Create(&insertbody)
	return result.Error
}

func InsertWarehouse(insertbody *Warehouse) error {
	result := db.Create(&insertbody)
	return result.Error
}

func InsertForkCat(insertbody *ForkCat) error {
	result := db.Create(&insertbody)
	return result.Error
}

func InsertBattery(insertbody *Battery) error {
	result := db.Create(&insertbody)
	return result.Error
}

func InsertBatSwitch(insertbody *SwitchBattery) error {
	result := db.Create(&insertbody)
	return result.Error
}

func UpdateBatStatus(updateBat *Battery, updateBatSwitch *SwitchBattery, fullNo string, fullnocat string) error {
	battery := Battery{}
	worker := db.Begin()
	batResult := worker.Model(&battery).Where("full_no=?", fullnocat).Updates(updateBat)
	switchResult := worker.Create(&updateBatSwitch)
	if batResult.Error == nil && switchResult.Error == nil {
		worker.Commit()
		return nil
	} else {
		worker.Rollback()
		return batResult.Error
	}
}
