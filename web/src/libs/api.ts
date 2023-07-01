import {EventExportType} from './enums/event-export-type';

// `${window.location.host}/api/v1`
const BaseUrl = `http://localhost:8000/api/v1`

const Endpoint = {
  Auth: {
    Validate: "auth/validate",
    Verify: "auth/verify",
    Profile: "auth/profile",
    Logout: "auth/logout",
  },
  User: {
    List: "users",
    Invite: "users/invite",
    Remove: "users/remove",
  },
  Events: {
    List: "events",
    Store: "events",
    Validate: "events/validate",
    Overview: (eventId: string) => `events/${eventId}/overview`,
    Participants: (eventId: string) => `events/${eventId}/participants`,
    Sync: (eventId: string) => `events/${eventId}/sync`,
    GenerateTicket: (eventId: string, participantId: number) => `events/${eventId}/participants/${participantId}/ticket`,
    Status: (eventId: string, participantId: number) => `events/${eventId}/participants/${participantId}/status`,
    Export: (eventId: string, type: EventExportType) => `events/${eventId}/export/${type}`,
  }
}

const ErrorMessage = {
  Auth: {
    ToManyRequest:  "For security purposes, you can only request this once every 60 seconds"
  }
}

export {BaseUrl, Endpoint, ErrorMessage}