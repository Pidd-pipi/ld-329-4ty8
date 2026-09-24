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

export type AppointmentStatus = 'pending' | 'confirmed' | 'completed' | 'cancelled';
export type MeetingMode = 'online' | 'offline';

export interface Appointment {
  id: number;
  matchId: number;
  initiator: string;
  confirmer: string;
  provider: string;
  learner: string;
  pair: string;
  offerSkill: string;
  wantedSkill: string;
  slot: string;
  meetingMode: MeetingMode;
  location: string;
  agenda: string;
  status: AppointmentStatus;
  cancelledBy?: string;
  createdAt: string;
  confirmedAt?: string;
  completedAt?: string;
  cancelledAt?: string;
}

export interface CreateAppointmentPayload {
  matchId: number;
  initiator: string;
  slot: string;
  meetingMode: MeetingMode;
  location: string;
  agenda?: string;
}

export interface AppointmentActionPayload {
  actor: string;
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
