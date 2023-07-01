import {Participant} from './participant';

export interface Event {
  id: number;
  google_form_id: string;
  name: string;
  location: string;
  preregister_date: number;
  event_date: number;
  total_participants: number;
  total_approved_participant: number;
  total_waiting_approval_participant: number;
  total_declined_participant: number;
  weekly_overview: EventWeeklyOverview[];
  latest_respondents: Participant[];
  is_active: boolean;
}

export interface EventWeeklyOverview {
  name: string;
  total: number;
}
