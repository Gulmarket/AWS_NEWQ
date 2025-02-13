package model

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

type Address struct {
    City   string `json:"city"`
    Street string `json:"street"`
    House  string `json:"house"`
}

type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("type assertion to []byte failed")
    }
    return json.Unmarshal(bytes, a)
}

func (a StringArray) Value() (driver.Value, error) {
    return json.Marshal(a)
}

type Shop struct {
    Id           int         `db:"id" json:"id"`
    Name         string      `db:"name" json:"name"`
    Description  string      `db:"description" json:"description"`
    Addresses    StringArray `db:"addresses" json:"addresses"`
    ShopLogo     *string     `db:"shop_logo" json:"shop_logo"`
    WorkSchedule Schedule    `db:"work_schedule" json:"work_schedule"`
    ProUserId    int         `db:"pro_user_id" json:"pro_user_id"`
}

type PlantationInfo struct {
    Name         string   `db:"name" json:"name"`
    Description  string   `db:"description" json:"description"`
    Country      string   `db:"country" json:"country"`
    City         string   `db:"city" json:"city"`
    LogoUrl      string   `db:"logo_url" json:"logo_url"`
    WorkSchedule Schedule `db:"work_schedule" json:"work_schedule"`
}
