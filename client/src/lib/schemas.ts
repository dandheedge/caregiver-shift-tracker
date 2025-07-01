import { z } from 'zod';

export const LocationSchema = z.object({
  latitude: z.number(),
  longitude: z.number(),
});

export const ScheduleSchema = z.object({
  id: z.number(),
  client_name: z.string(),
  shift_start: z.string(),
  shift_end: z.string(),
  latitude: z.number(),
  longitude: z.number(),
  status: z.enum(['upcoming', 'in_progress', 'completed', 'missed', 'cancelled']),
  created_at: z.string(),
  updated_at: z.string(),
}).transform((data) => {
  return {
    ...data,
    start_time: data.shift_start,
    end_time: data.shift_end,
    caregiver_name: 'Caregiver TBD',
    client_address: 'Address not provided',
    status: data.status === 'upcoming' ? 'scheduled' as const : data.status,
  };
});

export const TaskSchema = z.object({
  id: z.number(),
  schedule_id: z.number(),
  description: z.string(),
  status: z.enum(['pending', 'completed', 'not_completed']),
  reason: z.string().optional(),
  created_at: z.string(),
  updated_at: z.string(),
}).transform((data) => {
  return {
    ...data,
    title: data.description,
  };
});

export const ActivitySchema = z.object({
  id: z.number(),
  schedule_id: z.number(),
  title: z.string(),
  description: z.string(),
  is_resolved: z.boolean(),
  reason: z.string().optional(),
  created_at: z.string(),
  updated_at: z.string(),
}).transform((data) => {
  return {
    ...data,
    name: data.title,
    status: data.is_resolved ? 'completed' as const : 'pending' as const,
    progress: data.is_resolved ? 100 : 0,
    notes: data.reason,
  };
});

export const VisitSchema = z.object({
  id: z.number(),
  schedule_id: z.number(),
  start_time: z.string().optional(),
  end_time: z.string().optional(),
  start_lat: z.number().optional(),
  start_lng: z.number().optional(),
  end_lat: z.number().optional(),
  end_lng: z.number().optional(),
  created_at: z.string(),
  updated_at: z.string(),
}).transform((data) => {
  // Transform for frontend compatibility
  return {
    ...data,
    start_location: data.start_lat && data.start_lng ? {
      latitude: data.start_lat,
      longitude: data.start_lng,
    } : undefined,
    end_location: data.end_lat && data.end_lng ? {
      latitude: data.end_lat,
      longitude: data.end_lng,
    } : undefined,
    status: data.start_time ? (data.end_time ? 'completed' as const : 'in_progress' as const) : 'not_started' as const,
  };
});

export const StatsSchema = z.object({
  total_schedules: z.number(),
  missed_schedules: z.number(),
  upcoming_today: z.number(),
  completed_today: z.number(),
}).transform((data) => {
  return {
    ...data,
    upcoming_today_schedules: data.upcoming_today,
    completed_today_schedules: data.completed_today,
  };
});

export const StartVisitRequestSchema = z.object({
  latitude: z.number(),
  longitude: z.number(),
});

export const EndVisitRequestSchema = z.object({
  latitude: z.number(),
  longitude: z.number(),
});

export const UpdateTaskRequestSchema = z.object({
  status: z.enum(['completed', 'not_completed']),
  reason: z.string().optional(),
});

export const CreateActivityRequestSchema = z.object({
  title: z.string(),
  description: z.string(),
});

export const UpdateActivityRequestSchema = z.object({
  is_resolved: z.boolean(),
  reason: z.string().optional(),
});

export type Location = z.infer<typeof LocationSchema>;
export type Schedule = z.infer<typeof ScheduleSchema>;
export type Task = z.infer<typeof TaskSchema>;
export type Activity = z.infer<typeof ActivitySchema>;
export type Visit = z.infer<typeof VisitSchema>;
export type Stats = z.infer<typeof StatsSchema>;
export type StartVisitRequest = z.infer<typeof StartVisitRequestSchema>;
export type EndVisitRequest = z.infer<typeof EndVisitRequestSchema>;
export type UpdateTaskRequest = z.infer<typeof UpdateTaskRequestSchema>;
export type CreateActivityRequest = z.infer<typeof CreateActivityRequestSchema>;
export type UpdateActivityRequest = z.infer<typeof UpdateActivityRequestSchema>; 