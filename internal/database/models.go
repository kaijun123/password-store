package database

import "gorm.io/gorm"

// Credentials provided by the user during sign-up and sign-in
type RawUserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct for data saved in the db
type StoredUserCredentials struct {
	Username string `gorm:"primaryKey"`
	Salt     string
	Hash     []byte
}

// Struct for all types of transactions
// Type can be "Transfer", "Deposit", "Withdrawal"
type UserTransaction struct {
	gorm.Model
	Type string  `json:"type,omitempty"`
	From string  `json:"from,omitempty"`
	To   string  `json:"to,omitempty"`
	Amt  float32 `json:"amt"`
}

type UserBalance struct {
	gorm.Model
	Username string  `json:"username"`
	Balance  float32 `json:"balance"`
}
