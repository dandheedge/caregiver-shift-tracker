import { useState, useEffect } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { MapPin, Clock } from 'lucide-react';
import { apiClient, getCurrentPosition, GeolocationError, GeolocationErrorType } from '@/lib/api';
import type { Schedule } from '@/lib/schemas';

interface ActiveVisitProps {
  schedule: Schedule;
}

export function ActiveVisit({ schedule }: ActiveVisitProps) {
  const [elapsedTime, setElapsedTime] = useState('00:00:00');
  const [isEndingVisit, setIsEndingVisit] = useState(false);
  const queryClient = useQueryClient();
  const navigate = useNavigate();

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

  // Calculate elapsed time from start_time
  useEffect(() => {
    const updateElapsedTime = () => {
      if (!schedule.start_time) return;

      const startTime = new Date(schedule.start_time);
      const now = new Date();
      const diffMs = now.getTime() - startTime.getTime();
      
      const hours = Math.floor(diffMs / (1000 * 60 * 60));
      const minutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60));
      const seconds = Math.floor((diffMs % (1000 * 60)) / 1000);

      setElapsedTime(
        `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
      );
    };

    updateElapsedTime();
    const interval = setInterval(updateElapsedTime, 1000);

    return () => clearInterval(interval);
  }, [schedule.start_time]);

  const handleGeolocationError = (error: unknown) => {
    if (error instanceof GeolocationError) {
      switch (error.type) {
        case GeolocationErrorType.PERMISSION_DENIED:
          alert('Location access is required to clock out. Please enable location permissions in your browser settings and reload the page.');
          break;
        case GeolocationErrorType.NOT_SUPPORTED:
          alert('Your browser does not support location services. Please use a modern browser to access this feature.');
          break;
        default:
          alert(`Location error: ${error.message}`);
      }
    } else {
      console.error('Failed to access location:', error);
      alert('Unable to access your location. Please ensure location services are enabled and try again.');
    }
  };

  const handleEndVisit = async (e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent navigation when clicking clock-out
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

  const handleNavigateToDetails = () => {
    navigate(`/schedule/${schedule.id}`);
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
    <div className="bg-emerald-600 text-white">
      <div className="max-w-7xl mx-auto px-6 py-4">
        <div 
          className="flex items-center justify-between cursor-pointer"
          onClick={handleNavigateToDetails}
        >
          <div className="text-center">
            <div className="text-2xl font-mono font-bold">
              {elapsedTime}
            </div>
          </div>

          <div className="flex items-center space-x-4 flex-1 ml-8">
            <Avatar className="w-12 h-12">
              <AvatarImage src="/api/placeholder/48/48" />
              <AvatarFallback className="bg-white text-emerald-600">
                {getInitials(schedule.client_name)}
              </AvatarFallback>
            </Avatar>
            
            <div className="flex-1">
              <h3 className="font-semibold text-lg">{schedule.client_name}</h3>
              <div className="flex items-center text-sm text-emerald-100 space-x-4">
                <div className="flex items-center">
                  <MapPin className="w-4 h-4 mr-1" />
                  <span>{schedule.latitude}, {schedule.longitude}</span>
                </div>
                <div className="flex items-center">
                  <Clock className="w-4 h-4 mr-1" />
                  <span>
                    {formatTime(schedule.start_time || '')} - {formatTime(schedule.end_time || '')} SGT
                  </span>
                </div>
              </div>
            </div>
          </div>

          <Button
            onClick={handleEndVisit}
            disabled={isEndingVisit}
            variant="secondary"
            className="bg-white text-emerald-600 hover:bg-emerald-50"
          >
            {isEndingVisit ? 'Clock-Out...' : 'Clock-Out'}
          </Button>
        </div>
      </div>
    </div>
  );
} 