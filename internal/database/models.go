package database

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
	// ID        uint      `json:"id,omitempty" gorm:"primarykey"`
	// CreatedAt time.Time `json:"created,omitempty"`
	Type string `json:"type,omitempty"`
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
	Amt  int    `json:"amt"`
}

type UserBalance struct {
	Username string `json:"username" gorm:"primaryKey"`
	Balance  int    `json:"balance,omitempty"`
}
