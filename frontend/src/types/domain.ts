export interface Skill {
  id: number;
  owner: string;
  title: string;
  category: string;
  level: number;
  campus: string;
  description: string;
  timeSlots: string[];
  rewards: string[];
  portfolio: string;
}

export interface Need {
  id: number;
  requester: string;
  title: string;
  category: string;
  campus: string;
  expectTime: string;
  budgetType: string;
  description: string;
  responses: number;
}

export interface Match {
  id: number;
  provider: string;
  learner: string;
  offerSkill: string;
  wantedSkill: string;
  score: number;
  commonSlots: string[];
  recommendation: string;
}

// 预约状态机：PENDING 待对方确认 → CONFIRMED 待完成 → COMPLETED 已完成；CONFIRMED/PENDING 可被 CANCELLED。
export type AppointmentStatus = 'PENDING' | 'CONFIRMED' | 'COMPLETED' | 'CANCELLED';

export type LocationType = 'ONLINE' | 'OFFLINE';

export interface Appointment {
  id: number;
  matchId: number;
  initiator: string;
  confirmer: string;
  provider: string;
  learner: string;
  offerSkill: string;
  wantedSkill: string;
  pair: string;
  slot: string;
  locationType: LocationType;
  place: string;
  agenda: string;
  status: AppointmentStatus;
  statusLabel: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateAppointmentPayload {
  matchId: number;
  slot: string;
  locationType: LocationType;
  place: string;
  agenda: string;
}

export interface AppointmentMutationResult {
  appointment: Appointment;
  appointments: Appointment[];
}

export interface Review {
  id: number;
  from: string;
  to: string;
  rating: number;
  content: string;
}

export interface Conversation {
  id: number;
  withUser: string;
  unread: number;
  messages: string[];
}

export interface Profile {
  name: string;
  major: string;
  creditScore: number;
  creditLevel: string;
  skillWall: Skill[];
  radar: Record<string, number>;
  history: string[];
  reviews: Review[];
}

export interface Overview {
  service: string;
  categories: string[];
  metrics: Record<string, number>;
  skills: Skill[];
  needs: Need[];
  matches: Match[];
  appointments: Appointment[];
  reviews: Review[];
  messages: Conversation[];
  profile: Profile;
}
