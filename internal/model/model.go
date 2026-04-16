package model

type StaffMapping struct {
	StaffPassID string
	TeamName    string
	CreatedAt   int64
}

type Redemption struct {
	TeamName   string
	RedeemedAt int64
}
