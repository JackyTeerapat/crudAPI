package models

type Researcher_get struct {
	ProfileID      int               `json:"profile_id"`
	Profile_status bool              `json:"profile_status"`
	PrefixName     string            `json:"prefix_name"`
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	University     string            `json:"university"`
	AddressHome    string            `json:"address_home"`
	AddressWork    string            `json:"address_work"`
	Email          string            `json:"email"`
	PhoneNumber    string            `json:"phone_number"`
	Degree         []TempDegree_get  `gorm:"-" json:"degree"`
	Program        []TempProgram_get `gorm:"-" json:"program"`
}

type TempDegree_get struct {
	DegreeID         int    `json:"id"`
	DegreeType       string `json:"degree_type"`
	DegreeProgram    string `json:"degree_program"`
	DegreeUniversity string `json:"degree_university"`
}

type TempProgram_get struct {
	ProgramID   int    `json:"program_id"`
	ProgramName string `json:"program_name"`
}
