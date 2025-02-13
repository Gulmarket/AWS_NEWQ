package model

type Delivery struct {
    Id           int     `db:"id" json:"id"`
    FarmBox      string  `db:"farm_box" json:"farm_box"`
    BoxType      string  `db:"box_type" json:"box_type"`
    BoxSize      string  `db:"box_size" json:"box_size"`
    Mixed        string  `db:"mixed" json:"mixed"`
    Species      string  `db:"species" json:"species"`
    Product      string  `db:"product" json:"product"`
    Color        string  `db:"color" json:"color"`
    Length       string  `db:"length" json:"length"`
    Price        float64 `db:"price" json:"price"`
    Boxes        int     `db:"boxes" json:"boxes"`
    Packing      int     `db:"packing" json:"packing"`
    PlantationId string  `db:"plantation_id" json:"plantation_id"`
}
