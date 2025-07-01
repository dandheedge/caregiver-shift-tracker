import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { MapPin, Clock, Calendar } from 'lucide-react';
import { apiClient, getCurrentPosition, GeolocationError, GeolocationErrorType } from '@/lib/api';
import type { Schedule } from '@/lib/schemas';

interface ScheduleItemProps {
  schedule: Schedule;
}

export function ScheduleItem({ schedule }: ScheduleItemProps) {
  const [isStartingVisit, setIsStartingVisit] = useState(false);
  const [isEndingVisit, setIsEndingVisit] = useState(false);
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const startVisitMutation = useMutation({
    mutationFn: async (scheduleId: number) => {
      const position = await getCurrentPosition();
      return apiClient.startVisit(scheduleId, {
        latitude: position.coords.latitude,
        longitude: position.coords.longitude,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['schedules', 'today'] });
      queryClient.invalidateQueries({ queryKey: ['stats'] });
    },
  });

  const endVisitMutation = useMutation({
    mutationFn: async (scheduleId: number) => {
      const position = await getCurrentPosition();
      return apiClient.endVisit(scheduleId, {
        latitude: position.coords.latitude,
        longitude: position.coords.longitude,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['schedules', 'today'] });
      queryClient.invalidateQueries({ queryKey: ['stats'] });
    },
  });

  const handleGeolocationError = (error: unknown) => {
    if (error instanceof GeolocationError) {
      switch (error.type) {
        case GeolocationErrorType.PERMISSION_DENIED:
          alert('Location access is required to clock in/out. Please enable location permissions in your browser settings and reload the page.');
          break;
        case GeolocationErrorType.NOT_SUPPORTED:
          alert('Your browser does not support location services. Please use a modern browser to access this feature.');
          break;
        case GeolocationErrorType.POSITION_UNAVAILABLE:
          alert('Unable to access your location. Please ensure location services are enabled and try again.');
          break;
        case GeolocationErrorType.TIMEOUT:
          alert('Location request timed out. Please check your connection and try again.');
          break;
        default:
          alert(`Location error: ${error.message}`);
      }
    } else {
      console.error('Failed to access location:', error);
      alert('Unable to access your location. Please ensure location services are enabled and try again.');
    }
  };

  const handleStartVisit = async () => {
    setIsStartingVisit(true);
    try {
      await startVisitMutation.mutateAsync(schedule.id);
    } catch (error) {
      if (error instanceof GeolocationError) {
        handleGeolocationError(error);
      } else {
        console.error('Failed to start visit:', error);
        alert('Failed to start visit. Please try again.');
      }
    } finally {
      setIsStartingVisit(false);
    }
  };

  const handleEndVisit = async () => {
    setIsEndingVisit(true);
    try {
      await endVisitMutation.mutateAsync(schedule.id);
    } catch (error) {
      if (error instanceof GeolocationError) {
        handleGeolocationError(error);
      } else {
        console.error('Failed to end visit:', error);
        alert('Failed to end visit. Please try again.');
      }
    } finally {
      setIsEndingVisit(false);
    }
  };

  const getStatusBadge = (status: Schedule['status']) => {
    switch (status) {
      case 'scheduled':
        return <Badge className="bg-blue-100 text-blue-800 hover:bg-blue-100">Scheduled</Badge>;
      case 'in_progress':
        return <Badge className="bg-orange-100 text-orange-800 hover:bg-orange-100">In Progress</Badge>;
      case 'completed':
        return <Badge className="bg-green-100 text-green-800 hover:bg-green-100">Completed</Badge>;
      case 'missed':
        return <Badge className="bg-red-100 text-red-800 hover:bg-red-100">Missed</Badge>;
      case 'cancelled':
        return <Badge variant="secondary">Cancelled</Badge>;
      default:
        return <Badge variant="secondary">{status}</Badge>;
    }
  };

  const formatTime = (timeString: string) => {
    try {
      const date = new Date(timeString);
      return date.toLocaleTimeString('en-US', {
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
      });
    } catch {
      return timeString;
    }
  };

  const formatDate = (timeString: string) => {
    try {
      const date = new Date(timeString);
      return date.toLocaleDateString('en-US', {
        weekday: 'short',
        day: 'numeric',
        month: 'short',
        year: 'numeric'
      });
    } catch {
      return '';
    }
  };

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(word => word[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  return (
    <div className="flex flex-col space-y-4 p-4 rounded-lg border border-slate-200 bg-white">
      <div className="flex-shrink-0">
        {getStatusBadge(schedule.status)}
      </div>

      <div className="flex flex-row space-x-4">
        <div className="flex-shrink-0">
          <Avatar className="w-12 h-12">
            <AvatarImage src="/api/placeholder/48/48" />
            <AvatarFallback className="bg-slate-100 text-slate-600">
              {getInitials(schedule.client_name)}
            </AvatarFallback>
          </Avatar>
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between">
            <div>
              <h4 className="font-medium text-slate-900">{schedule.client_name}</h4>
              <p className="text-sm text-slate-600">{schedule.caregiver_name || 'Caregiver TBD'}</p>
            </div>
          </div>
        </div>
      </div>
      <div className="flex items-center text-sm text-slate-600">
        <MapPin className="w-4 h-4 mr-1" />
        <span className="truncate">{schedule.latitude}, {schedule.longitude}</span>
      </div>

      <div className="flex flex-row space-x-4">
        <div className="flex-1 flex items-center text-sm bg-sky-200 text-gray-700 px-3 py-2 rounded-md border">
          
          <Calendar className="w-4 h-4 mr-2" />
          <span>{formatDate(schedule.start_time || '')}</span>
        </div>
        <div className="flex-1 flex items-center text-sm bg-sky-200 text-gray-700 px-3 py-2 rounded-md border">
          <Clock color='' className="w-4 h-4 mr-2" />
          <span>{formatTime(schedule.start_time || '')} - {formatTime(schedule.end_time || '')}</span>
        </div>
      </div>

      <div className="flex items-center space-x-2">
        {schedule.status === 'scheduled' && (
          <Button
            onClick={handleStartVisit}
            disabled={isStartingVisit}
            className="bg-emerald-600 hover:bg-emerald-700 text-white w-full"
          >
            {isStartingVisit ? 'Clock-In Now...' : 'Clock-In Now'}
          </Button>
        )}

        {schedule.status === 'in_progress' && (
          <>
            <Button 
              variant="outline" 
              size="lg" 
              className="w-1/2"
              onClick={() => navigate(`/schedule/${schedule.id}`)}
            >
              View Progress
            </Button>
            <Button
              size="lg"
              onClick={handleEndVisit}
              disabled={isEndingVisit}
              className="bg-emerald-600 hover:bg-emerald-700 text-white w-1/2"
            >
              {isEndingVisit ? 'Clock-Out Now...' : 'Clock-Out Now'}
            </Button>
          </>
        )}

        {schedule.status === 'completed' && (
          <Button 
            variant="outline" 
            size="sm"
            onClick={() => navigate(`/schedule/${schedule.id}`)}
          >
            View Report
          </Button>
        )}

        {schedule.status === 'scheduled' && (
          <Button 
            variant="outline" 
            size="sm"
            onClick={() => navigate(`/schedule/${schedule.id}`)}
            className="w-full mt-2"
          >
            View Details
          </Button>
        )}

      </div>
    </div>
  );
} 