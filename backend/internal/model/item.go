package model

type Skill struct {
	ID          int      `json:"id"`
	Owner       string   `json:"owner"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Level       int      `json:"level"`
	Campus      string   `json:"campus"`
	Description string   `json:"description"`
	TimeSlots   []string `json:"timeSlots"`
	Rewards     []string `json:"rewards"`
	Portfolio   string   `json:"portfolio"`
}

type Need struct {
	ID          int    `json:"id"`
	Requester   string `json:"requester"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Campus      string `json:"campus"`
	ExpectTime  string `json:"expectTime"`
	BudgetType  string `json:"budgetType"`
	Description string `json:"description"`
	Responses   int    `json:"responses"`
}

type Match struct {
	ID             int      `json:"id"`
	Provider       string   `json:"provider"`
	Learner        string   `json:"learner"`
	OfferSkill     string   `json:"offerSkill"`
	WantedSkill    string   `json:"wantedSkill"`
	Score          int      `json:"score"`
	CommonSlots    []string `json:"commonSlots"`
	Recommendation string   `json:"recommendation"`
}

type Appointment struct {
	ID           int    `json:"id"`
	MatchID      int    `json:"matchId"`
	Initiator    string `json:"initiator"`
	Confirmer    string `json:"confirmer"`
	Provider     string `json:"provider"`
	Learner      string `json:"learner"`
	OfferSkill   string `json:"offerSkill"`
	WantedSkill  string `json:"wantedSkill"`
	Pair         string `json:"pair"`
	Slot         string `json:"slot"`
	LocationType string `json:"locationType"`
	Place        string `json:"place"`
	Agenda       string `json:"agenda"`
	Status       string `json:"status"`
	StatusLabel  string `json:"statusLabel"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// CreateAppointmentRequest 发起预约请求：匹配卡片、共同时段、线上会议链接或线下地点。
type CreateAppointmentRequest struct {
	MatchID      int    `json:"matchId"`
	Slot         string `json:"slot"`
	LocationType string `json:"locationType"`
	Place        string `json:"place"`
	Agenda       string `json:"agenda"`
}

// AppointmentMutationResult 预约变更结果：返回最新预约与全量列表，供前端立即刷新。
type AppointmentMutationResult struct {
	Appointment  Appointment   `json:"appointment"`
	Appointments []Appointment `json:"appointments"`
}

type Review struct {
	ID      int    `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Rating  int    `json:"rating"`
	Content string `json:"content"`
}

type Conversation struct {
	ID       int      `json:"id"`
	WithUser string   `json:"withUser"`
	Unread   int      `json:"unread"`
	Messages []string `json:"messages"`
}

type Profile struct {
	Name        string         `json:"name"`
	Major       string         `json:"major"`
	CreditScore int            `json:"creditScore"`
	CreditLevel string         `json:"creditLevel"`
	SkillWall   []Skill        `json:"skillWall"`
	Radar       map[string]int `json:"radar"`
	History     []string       `json:"history"`
	Reviews     []Review       `json:"reviews"`
}

type Overview struct {
	Service      string         `json:"service"`
	Categories   []string       `json:"categories"`
	Metrics      map[string]int `json:"metrics"`
	Skills       []Skill        `json:"skills"`
	Needs        []Need         `json:"needs"`
	Matches      []Match        `json:"matches"`
	Appointments []Appointment  `json:"appointments"`
	Reviews      []Review       `json:"reviews"`
	Messages     []Conversation `json:"messages"`
	Profile      Profile        `json:"profile"`
}
