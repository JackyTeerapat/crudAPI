package models

type ResearcherRequest struct {
	ProfileID    int                     `json:"profile_id"`
	FirstName    string                  `json:"first_name"`
	LastName     string                  `json:"last_name"`
	AddressHome  string                  `json:"address_home"`
	PositionName string                  `json:"position_name"`
	University   string                  `json:"university"`
	AddressWork  string                  `json:"address_work"`
	Email        string                  `json:"email"`
	PhoneNumber  string                  `json:"phone_number"`
	Degree       []TempDegree_create     `json:"degree"`
	Program      []TempProgram_create    `json:"program"`
	Experience   []TempExperience_create `json:"experience"`
	Explore      []TempExplore_create    `json:"explore"`
	// Attach      []TempAttach_create     `json:"attach"`
}

type TempDegree_create struct {
	DegreeType       string `json:"degree_type"`
	DegreeProgram    string `json:"degree_program"`
	DegreeUniversity string `json:"degree_university"`
	Activated        bool   `json:"activated" gorm:"default:true"`
}

type TempProgram_create struct {
	ProgramName string `json:"program_name"`
	Activated   bool   `json:"activated" gorm:"default:true"`
}

type TempExperience_create struct {
	ExperienceType       string  `json:"experience_type"`
	ExperienceStart      string  `json:"experience_start"`
	ExperienceEnd        *string `json:"experience_end"`
	ExperienceUniversity string  `json:"experience_university"`
	ExperienceRemark     string  `json:"experience_remark"`
	Activated            bool    `json:"activated" gorm:"default:true"`
}

type TempExplore_create struct {
	ExploreName   string `json:"explore_name"`
	ExploreYear   string `json:"explore_year"`
	ExploreDetail string `json:"explore_detail"`
	Activated     bool   `json:"activated" gorm:"default:true"`
}

type TempAttach_create struct {
	FileName     string `json:"file_name"`
	FileAction   string `json:"file_action"`
	File_storage string `json:"file_storage"`
	Activated    bool   `json:"activated" gorm:"default:true"`
}
