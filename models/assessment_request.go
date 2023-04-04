package models

type AssessmentRequest struct {
	AssessmentStart string `json:"assessment_start"`
	AssessmentEnd   string `json:"assessment_end"`
	// AssessmentFile  AssessmentFileDetails `gorm:"-" json:"assessment_file"`
	AssessmentFileName    string          `json:"assessment_file_name"`
	AssessmentFileStorage string          `json:"assessment_file_storage"`
	Project               ProjectDetails  `json:"assessment_project"`
	Progress              ProgressDetails `json:"assessment_progress"`
	Report                ReportDetails   `json:"assessment_report"`
	Article               ArticleDetails  `json:"assessment_article"`
}

type ProjectDetails struct {
	ProjectYear      string      `json:"project_year"`
	ProjectTitle     string      `json:"project_title"`
	ProjectPoint     int         `json:"project_point"`
	ProjectEstimate  bool        `json:"project_estimate"`
	ProjectRecommend bool        `json:"project_recommend"`
	ProjectFile      FileDetails `json:"project_file"`
	ProjectPeriod    bool        `json:"period"`
}

type ProgressDetails struct {
	ProgressYear      string      `json:"progress_year"`
	ProgressTitle     string      `json:"progress_title"`
	ProgressEstimate  bool        `json:"progress_estimate"`
	ProgressRecommend bool        `json:"progress_recommend"`
	ProgressFile      FileDetails `json:"progress_file"`
	ProgressPeriod    bool        `json:"period"`
}

type ReportDetails struct {
	ReportYear      string      `json:"report_year"`
	ReportTitle     string      `json:"report_title"`
	ReportEstimate  bool        `json:"report_estimate"`
	ReportRecommend bool        `json:"report_recommend"`
	ReportFile      FileDetails `json:"report_file"`
	ReportPeriod    bool        `json:"period"`
}

type ArticleDetails struct {
	ArticleYear      string      `json:"article_year"`
	ArticleTitle     string      `json:"article_title"`
	ArticleEstimate  bool        `json:"article_estimate"`
	ArticleRecommend bool        `json:"article_recommend"`
	ArticleFile      FileDetails `json:"article_file"`
	ArticlePeriod    bool        `json:"period"`
}
type AssessmentFileDetails struct {
	FileName    string `json:"assessment_file_name"`
	FileStorage string `json:"assessment_file_storage"`
}
type FileDetails struct {
	FileName    string `json:"file_name"`
	FileStorage string `json:"file_storage"`
}
