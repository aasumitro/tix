export interface Participant {
  id: number;
  event_id: number;
  name: string;
  email: string;
  phone: string;
  job: string;
  prof_of_payment: string;
  date_of_birth: string;
  approved_at?: number | null;
  declined_at?: number | null;
  status: string;
  declined_reason: string;
}