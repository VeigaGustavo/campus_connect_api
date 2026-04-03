package models

// DiscoverItem um elemento do feed Descobrir.
type DiscoverItem struct {
	ID            string       `json:"id"`
	Kind          DiscoverKind `json:"kind"`
	Title         string       `json:"title"`
	Subtitle      string       `json:"subtitle,omitempty"`
	Excerpt       string       `json:"excerpt,omitempty"`
	MetaPrimary   string       `json:"meta_primary,omitempty"`
	MetaSecondary string       `json:"meta_secondary,omitempty"`
	ReferenceID   string       `json:"reference_id"`
}

// Opportunity lista e detalhe.
type Opportunity struct {
	ID               string       `json:"id"`
	Title            string       `json:"title"`
	CompanyName      string       `json:"company_name"`
	ShortDescription string       `json:"short_description"`
	FullDescription  string       `json:"full_description,omitempty"`
	ApplyDeadline    string       `json:"apply_deadline"` // ISO 8601
	WorkLocation     WorkLocation `json:"work_location"`
	TypeLabel        string       `json:"type_label,omitempty"`
	Requirements     []string     `json:"requirements,omitempty"`
}

// CampusEvent evento no campus.
type CampusEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	StartAt     string `json:"start_at"` // ISO 8601
	Location    string `json:"location,omitempty"`
	Organizer   string `json:"organizer,omitempty"`
}

// StudyGroup grupo de estudo.
type StudyGroup struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	FieldOfStudy  string     `json:"field_of_study"`
	Description   string     `json:"description,omitempty"`
	Level         GroupLevel `json:"level"`
	MemberCount   int        `json:"member_count"`
	ScheduleLabel string     `json:"schedule_label,omitempty"`
}

// Interest no perfil (objeto com label).
type Interest struct {
	Label string `json:"label"`
}

// ProfileActivity item de atividade recente.
type ProfileActivity struct {
	Kind           ProfileActivityKind `json:"kind"`
	TitleHighlight string              `json:"title_highlight"`
	Subtitle       string              `json:"subtitle,omitempty"`
	TimeAgoLabel   string              `json:"time_ago_label,omitempty"`
	OccurredAt     string              `json:"occurred_at,omitempty"` // ISO 8601; evoluir para timestamp
}

// UserProfile resposta de GET /me ou /users/me.
type UserProfile struct {
	Name              string            `json:"name"`
	Initials          string            `json:"initials,omitempty"`
	CourseAndSemester string            `json:"course_and_semester,omitempty"`
	Email             string            `json:"email"`
	CityState         string            `json:"city_state,omitempty"`
	ApplicationsCount int               `json:"applications_count"`
	GroupsCount       int               `json:"groups_count"`
	EventsCount       int               `json:"events_count"`
	Interests         []Interest        `json:"interests"`
	RecentActivity    []ProfileActivity `json:"recent_activity"`
}

// APIError corpo JSON padrão para erros HTTP.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
