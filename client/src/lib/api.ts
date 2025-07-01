import { z } from 'zod';
import {
  ScheduleSchema,
  TaskSchema,
  ActivitySchema,
  StatsSchema,
  StartVisitRequestSchema,
  EndVisitRequestSchema,
  UpdateTaskRequestSchema,
  CreateActivityRequestSchema,
  UpdateActivityRequestSchema,
  type Schedule,
  type Task,
  type Activity,
  type Stats,
  type StartVisitRequest,
  type EndVisitRequest,
  type UpdateTaskRequest,
  type CreateActivityRequest,
  type UpdateActivityRequest,
} from '@/lib/schemas';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://server.31.97.179.158.sslip.io/api/v1';



class ApiClient {
  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
    schema?: z.ZodSchema<T>
  ): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`);
    }

    const data = await response.json();
    
    if (schema) {
      return schema.parse(data);
    }
    
    return data;
  }

  // Schedule endpoints
  async getAllSchedules(): Promise<Schedule[]> {
    const data = await this.request('/schedules');
    return z.array(ScheduleSchema).parse(data);
  }

  async getTodaySchedules(): Promise<Schedule[]> {
    const data = await this.request('/schedules/today');
    return z.array(ScheduleSchema).parse(data);
  }

  async getScheduleById(id: number): Promise<Schedule> {
    const data = await this.request(`/schedules/${id}`);
    return ScheduleSchema.parse(data);
  }

  async getTasksBySchedule(scheduleId: number): Promise<Task[]> {
    const data = await this.request(`/schedules/${scheduleId}/tasks`);
    return z.array(TaskSchema).parse(data);
  }

  // Visit endpoints
  async startVisit(scheduleId: number, location: StartVisitRequest): Promise<void> {
    // Validate the request data
    StartVisitRequestSchema.parse(location);
    
    return this.request(`/schedules/${scheduleId}/start`, {
      method: 'POST',
      body: JSON.stringify(location),
    });
  }

  async endVisit(scheduleId: number, location: EndVisitRequest): Promise<void> {
    // Validate the request data
    EndVisitRequestSchema.parse(location);
    
    return this.request(`/schedules/${scheduleId}/end`, {
      method: 'POST',
      body: JSON.stringify(location),
    });
  }

  // Task endpoints
  async updateTask(taskId: number, update: UpdateTaskRequest): Promise<void> {
    // Validate the request data
    UpdateTaskRequestSchema.parse(update);
    
    return this.request(`/tasks/${taskId}/update`, {
      method: 'POST',
      body: JSON.stringify(update),
    });
  }

  // Activity endpoints
  async getActivityById(id: number): Promise<Activity> {
    const data = await this.request(`/activities/${id}`);
    return ActivitySchema.parse(data);
  }

  async getActivitiesBySchedule(scheduleId: number): Promise<Activity[]> {
    const data = await this.request(`/schedules/${scheduleId}/activities`);
    return z.array(ActivitySchema).parse(data);
  }

  async createActivity(scheduleId: number, activity: CreateActivityRequest): Promise<Activity> {
    // Validate the request data
    CreateActivityRequestSchema.parse(activity);
    
    const data = await this.request(`/schedules/${scheduleId}/activities`, {
      method: 'POST',
      body: JSON.stringify(activity),
    });
    return ActivitySchema.parse(data);
  }

  async updateActivity(id: number, update: UpdateActivityRequest): Promise<Activity> {
    // Validate the request data
    UpdateActivityRequestSchema.parse(update);
    
    const data = await this.request(`/activities/${id}`, {
      method: 'PUT',
      body: JSON.stringify(update),
    });
    return ActivitySchema.parse(data);
  }

  // Stats endpoint
  async getStats(): Promise<Stats> {
    const data = await this.request('/stats');
    return StatsSchema.parse(data);
  }
}

export const apiClient = new ApiClient();

// Enhanced Geolocation helper with mandatory access and fallback handling
export const GeolocationErrorType = {
  NOT_SUPPORTED: 'GEOLOCATION_NOT_SUPPORTED',
  PERMISSION_DENIED: 'GEOLOCATION_PERMISSION_DENIED',
  POSITION_UNAVAILABLE: 'GEOLOCATION_POSITION_UNAVAILABLE',
  TIMEOUT: 'GEOLOCATION_TIMEOUT',
  UNKNOWN: 'GEOLOCATION_UNKNOWN_ERROR'
} as const;

export type GeolocationErrorTypeValues = typeof GeolocationErrorType[keyof typeof GeolocationErrorType];

export class GeolocationError extends Error {
  type: GeolocationErrorTypeValues;
  originalError?: GeolocationPositionError;

  constructor(
    type: GeolocationErrorTypeValues,
    message: string,
    originalError?: GeolocationPositionError
  ) {
    super(message);
    this.name = 'GeolocationError';
    this.type = type;
    this.originalError = originalError;
  }
}

export const getCurrentPosition = (): Promise<GeolocationPosition> => {
  return new Promise((resolve, reject) => {
    // Check if geolocation is supported
    if (!navigator.geolocation) {
      reject(new GeolocationError(
        GeolocationErrorType.NOT_SUPPORTED,
        'Geolocation is not supported by this browser. Please use a modern browser that supports location services to access this application.',
      ));
      return;
    }

    // Request current position with high accuracy
    navigator.geolocation.getCurrentPosition(
      (position) => {
        // Validate position data
        if (!position.coords || 
            typeof position.coords.latitude !== 'number' || 
            typeof position.coords.longitude !== 'number') {
          reject(new GeolocationError(
            GeolocationErrorType.POSITION_UNAVAILABLE,
            'Invalid location data received. Please try again.',
          ));
          return;
        }
        resolve(position);
      },
      (error) => {
        let errorType: GeolocationErrorTypeValues;
        let message: string;

        switch (error.code) {
          case error.PERMISSION_DENIED:
            errorType = GeolocationErrorType.PERMISSION_DENIED;
            message = 'Location access is required for this application. Please enable location permissions in your browser settings and reload the page.';
            break;
          case error.POSITION_UNAVAILABLE:
            errorType = GeolocationErrorType.POSITION_UNAVAILABLE;
            message = 'Location information is unavailable. Please ensure location services are enabled on your device and try again.';
            break;
          case error.TIMEOUT:
            errorType = GeolocationErrorType.TIMEOUT;
            message = 'Location request timed out. Please check your connection and try again.';
            break;
          default:
            errorType = GeolocationErrorType.UNKNOWN;
            message = 'An unknown error occurred while accessing your location. Please try again.';
        }

        reject(new GeolocationError(errorType, message, error));
      },
      {
        enableHighAccuracy: true,
        timeout: 15000, // Increased timeout for better reliability
        maximumAge: 30000, // Reduced max age for more accurate positioning
      }
    );
  });
};

// Check if geolocation is available and permissions are granted
export const checkGeolocationAvailability = async (): Promise<boolean> => {
  if (!navigator.geolocation) {
    return false;
  }

  try {
    // Try to get permission status if available
    if ('permissions' in navigator) {
      const permission = await navigator.permissions.query({ name: 'geolocation' });
      return permission.state !== 'denied';
    }
    
    // Fallback: try to get position with short timeout
    await new Promise<GeolocationPosition>((resolve, reject) => {
      navigator.geolocation.getCurrentPosition(resolve, reject, {
        timeout: 1000,
        maximumAge: 60000,
      });
    });
    return true;
  } catch {
    return false;
  }
};

// Fallback location (to be used only when geolocation fails and is explicitly acknowledged by user)
export const FALLBACK_LOCATION = {
  latitude: 0,
  longitude: 0,
  accuracy: 0,
  description: 'Fallback location - actual location unavailable'
} as const; 