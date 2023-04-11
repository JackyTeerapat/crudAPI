package models

type AssessmentRequest struct {
	ProfileID       int             `json:"profile_id"`
	AssessmentStart string          `json:"assessment_start"`
	AssessmentEnd   string          `json:"assessment_end"`
	Project         ProjectDetails  `json:"assessment_project"`
	Progress        ProgressDetails `json:"assessment_progress"`
	Report          ReportDetails   `json:"assessment_report"`
	Article         ArticleDetails  `json:"assessment_article"`
}

type ProjectDetails struct {
	ProjectYear      string `json:"project_year"`
	ProjectTitle     string `json:"project_title"`
	ProjectPoint     int    `json:"project_point"`
	ProjectEstimate  bool   `json:"project_estimate"`
	ProjectRecommend bool   `json:"project_recommend"`
	ProjectPeriod    bool   `json:"project_period"`
}

type ProgressDetails struct {
	ProgressYear      string `json:"progress_year"`
	ProgressTitle     string `json:"progress_title"`
	ProgressEstimate  bool   `json:"progress_estimate"`
	ProgressRecommend bool   `json:"progress_recommend"`
	ProgressPeriod    bool   `json:"progress_period"`
}

type ReportDetails struct {
	ReportYear      string `json:"report_year"`
	ReportTitle     string `json:"report_title"`
	ReportEstimate  bool   `json:"report_estimate"`
	ReportRecommend bool   `json:"report_recommend"`
	ReportPeriod    bool   `json:"report_period"`
}

type ArticleDetails struct {
	ArticleYear      string `json:"article_year"`
	ArticleTitle     string `json:"article_title"`
	ArticleEstimate  bool   `json:"article_estimate"`
	ArticleRecommend bool   `json:"article_recommend"`
	ArticlePeriod    bool   `json:"article_period"`
}
