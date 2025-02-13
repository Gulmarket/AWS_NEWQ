package model

import (
    "encoding/json"
    "errors"
)

type DaySchedule struct {
    Start string `json:"start"`
    End   string `json:"end"`
}

type Schedule struct {
    Days *map[string]DaySchedule `json:"days"`
}

func (s *Schedule) Scan(value interface{}) error {
    b, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }

    return json.Unmarshal(b, &s)
}

type Plantation struct {
    Id           int       `db:"id" json:"id"`
    Email        string    `db:"email" json:"email"`
    Password     string    `db:"password" json:"-"`
    Name         *string   `db:"name" json:"name"`
    Description  *string   `db:"description" json:"description"`
    Country      *string   `db:"country" json:"country"`
    City         *string   `db:"city" json:"city"`
    LogoUrl      *string   `db:"logo_url" json:"logo_url"`
    Phone        string    `db:"phone" json:"phone"`
    WorkSchedule *Schedule `db:"work_schedule" json:"work_schedule"`
}

type Delivery struct {
    Id           int     `db:"id" json:"id"`
    FarmBox      string  `db:"farm_box" json:"farm_box"`
    BoxSize      string  `db:"box_size" json:"box_size"`
    Mixed        bool    `db:"mixed" json:"mixed"`
    Species      string  `db:"species" json:"species"`
    Product      string  `db:"product" json:"product"`
    Color        string  `db:"color" json:"color"`
    Length       string  `db:"length" json:"length"`
    Price        float64 `db:"price" json:"price"`
    Boxes        int     `db:"boxes" json:"boxes"`
    Packing      int     `db:"packing" json:"packing"`
    PlantationId int     `db:"plantation_id" json:"plantation_id"`
}

type Filter struct {
    FlowerType    []string `json:"flower_type,omitempty"`
    FlowerSpecies []string `json:"flower_species"`
    SizeStart     int      `json:"size_start"`
    SizeEnd       int      `json:"size_end"`
    BoxType       []string `json:"box_type,omitempty"`
    Color         []string `json:"color,omitempty"`
}

type Deliveries struct {
    PlantationName string  `db:"plantation_name" json:"plantation_name"`
    TengePrice     float64 `json:"tenge_price"`
    IsFavorite     bool    `json:"is_favorite"`
    Delivery
}

type FlowersFilter struct {
    FlowerTypes   []string            `json:"flower_types"`
    FlowerSpecies map[string][]string `json:"flower_species"`
    Deliveries    []Deliveries        `json:"delivery"`
}

type PlantationFlowersFilter struct {
    FlowerTypes   []string            `json:"flower_types"`
    FlowerSpecies map[string][]string `json:"flower_species"`
    Deliveries    []Delivery          `json:"delivery"`
}

type Card struct {
    Delivery    Delivery  `json:"delivery"`
    FlowersLeft int       `json:"flowers_left"`
    TengePrice  float64   `json:"tenge_price"`
    Addresses   []Address `json:"addresses"`
}
