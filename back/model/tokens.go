package model

// RefreshToken stores token properties that
// are accessed in multiple application layers
type RefreshToken struct {
    ID string `json:"-"`
    SS string `json:"refreshToken"`
}

// IDToken stores token properties that
// are accessed in multiple application layers
type IDToken struct {
    SS string `json:"idToken"`
}

// TokenPair used for returning pairs of id and refresh tokens
type TokenPair struct {
    IDToken
    RefreshToken
}
