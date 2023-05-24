package custom

type Profile struct {
	ProfileId     int    `gorm:"column:profile_id"`
	ProfileStatus string `gorm:"column:profile_status"`
	FirstName     string `gorm:"column:first_name"`
	LastName      string `gorm:"column:last_name"`
	PositionId    int    `gorm:"column:position_id"`
	PositionName  string `gorm:"column:position_name"`
	University    string `gorm:"column:university"`
	Email         string `gorm:"column:email"`
	Phone_number  string `gorm:"column:phone_number"`
}
