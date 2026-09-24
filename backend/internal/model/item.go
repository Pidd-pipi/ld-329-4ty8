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

// Appointment 交换预约：记录发起人、确认人、时段、线上会议链接或线下地点及状态流转
type Appointment struct {
	ID           int64  `json:"id"`
	MatchID      int    `json:"matchId"`
	Initiator    string `json:"initiator"`
	Confirmer    string `json:"confirmer"`
	Provider     string `json:"provider"`
	Learner      string `json:"learner"`
	Pair         string `json:"pair"`
	OfferSkill   string `json:"offerSkill"`
	WantedSkill  string `json:"wantedSkill"`
	Slot         string `json:"slot"`
	MeetingMode  string `json:"meetingMode"`
	Location     string `json:"location"`
	Agenda       string `json:"agenda"`
	Status       string `json:"status"`
	CancelledBy  string `json:"cancelledBy,omitempty"`
	CreatedAt    string `json:"createdAt"`
	ConfirmedAt  string `json:"confirmedAt,omitempty"`
	CompletedAt  string `json:"completedAt,omitempty"`
	CancelledAt  string `json:"cancelledAt,omitempty"`
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
