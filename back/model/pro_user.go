package model

type ProUser struct {
    Id       int     `db:"id" json:"id"`
    Email    string  `db:"email" json:"email"`
    Password string  `db:"password" json:"-"`
    Name     *string `db:"name" json:"name"`
    Surname  *string `db:"surname" json:"surname"`
    Patronym *string `db:"patronym" json:"patronym"`
    Role     string  `db:"role" json:"role"`
    City     *string `db:"city" json:"city"`
    Phone    string  `db:"phone" json:"phone"`
    Wallet   float64 `db:"wallet" json:"wallet"`
}

type ProUserInfo struct {
    Name     *string `db:"name" json:"name"`
    Surname  *string `db:"surname" json:"surname"`
    Patronym *string `db:"patronym" json:"patronym"`
    Role     string  `db:"role" json:"role"`
    City     *string `db:"city" json:"city"`
}

type PrivateOffice struct {
    ProUserInfo ProUserInfo `json:"pro_user"`
    Shops       []Shop      `json:"shops"`
}
