import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useParams, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ArrowLeft, Check, X, MapPin } from 'lucide-react';
import { apiClient } from '@/lib/api';
import { ClockOutSuccessModal } from './ClockOutSuccessModal';

interface ClockOutParams {
  id: string;
  [key: string]: string | undefined;
}

interface ActivityState {
  id: number;
  isResolved: boolean;
  reason: string;
}

export function ClockOut() {
  const { id } = useParams<ClockOutParams>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const scheduleId = parseInt(id || '0');
  
  const [activityStates, setActivityStates] = useState<Record<number, ActivityState>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [timer, setTimer] = useState({ hours: 1, minutes: 35, seconds: 40 });
  const [showSuccessModal, setShowSuccessModal] = useState(false);

  const { data: schedule, isLoading: scheduleLoading } = useQuery({
    queryKey: ['schedule', scheduleId],
    queryFn: () => apiClient.getScheduleById(scheduleId),
    enabled: !!scheduleId,
  });

  const { data: activities, isLoading: activitiesLoading } = useQuery({
    queryKey: ['activities', scheduleId],
    queryFn: () => apiClient.getActivitiesBySchedule(scheduleId),
    enabled: !!scheduleId,
  });

  // Initialize activity states when activities load
  useEffect(() => {
    if (activities) {
      const initialStates: Record<number, ActivityState> = {};
      activities.forEach(activity => {
        initialStates[activity.id] = {
          id: activity.id,
          isResolved: activity.is_resolved,
          reason: activity.reason || ''
        };
      });
      setActivityStates(initialStates);
    }
  }, [activities]);

  // Timer effect
  useEffect(() => {
    const interval = setInterval(() => {
      setTimer(prev => {
        let { hours, minutes, seconds } = prev;
        seconds++;
        if (seconds >= 60) {
          seconds = 0;
          minutes++;
          if (minutes >= 60) {
            minutes = 0;
            hours++;
          }
        }
        return { hours, minutes, seconds };
      });
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  const updateActivityMutation = useMutation({
    mutationFn: async (activityUpdate: { id: number; is_resolved: boolean; reason?: string }) => {
      return apiClient.updateActivity(activityUpdate.id, {
        is_resolved: activityUpdate.is_resolved,
        reason: activityUpdate.reason || ''
      });
    },
  });

  const handleActivityChange = (activityId: number, isResolved: boolean) => {
    setActivityStates(prev => ({
      ...prev,
      [activityId]: {
        ...prev[activityId],
        isResolved,
        reason: isResolved ? '' : prev[activityId]?.reason || ''
      }
    }));
  };

  const handleReasonChange = (activityId: number, reason: string) => {
    setActivityStates(prev => ({
      ...prev,
      [activityId]: {
        ...prev[activityId],
        reason
      }
    }));
  };

  const handleClockOut = async () => {
    // Validate that all activities marked as "No" have reasons
    const invalidActivities = Object.values(activityStates).filter(
      state => !state.isResolved && !state.reason.trim()
    );

    if (invalidActivities.length > 0) {
      alert('Please provide a reason for all activities marked as "No"');
      return;
    }

    setIsSubmitting(true);
    try {
      // Submit all activity updates
      await Promise.all(
        Object.values(activityStates).map(state =>
          updateActivityMutation.mutateAsync({
            id: state.id,
            is_resolved: state.isResolved,
            reason: state.reason
          })
        )
      );

      // Invalidate queries to refresh data
      queryClient.invalidateQueries({ queryKey: ['activities', scheduleId] });
      queryClient.invalidateQueries({ queryKey: ['schedules'] });
      queryClient.invalidateQueries({ queryKey: ['stats'] });

      // Show success modal instead of navigating immediately
      setShowSuccessModal(true);
    } catch (error) {
      console.error('Failed to update activities:', error);
      alert('Failed to clock out. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCancelClockIn = () => {
    navigate(`/schedule/${scheduleId}`);
  };

  const handleGoHome = () => {
    setShowSuccessModal(false);
    navigate('/');
  };

  const formatDuration = () => {
    return `${timer.hours} hour${timer.hours !== 1 ? 's' : ''} ${timer.minutes} minute${timer.minutes !== 1 ? 's' : ''}`;
  };

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(word => word[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  const formatTimer = () => {
    return `${timer.hours.toString().padStart(2, '0')} : ${timer.minutes.toString().padStart(2, '0')} : ${timer.seconds.toString().padStart(2, '0')}`;
  };

  if (scheduleLoading || activitiesLoading) {
    return (
      <div className="flex flex-col min-h-screen bg-slate-50">
        <div className="flex-1 p-6">
          <div className="max-w-md mx-auto">
            <div className="animate-pulse space-y-4">
              <div className="h-8 bg-slate-200 rounded w-1/3"></div>
              <div className="h-16 bg-slate-200 rounded"></div>
              <div className="h-48 bg-slate-200 rounded"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!schedule || !activities) {
    return (
      <div className="flex flex-col min-h-screen bg-slate-50">
        <div className="flex-1 p-6">
          <div className="max-w-md mx-auto text-center">
            <p className="text-slate-500">Schedule or activities not found</p>
            <Button onClick={() => navigate('/')} className="mt-4">
              Back to Dashboard
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen bg-slate-50">
      <header className="bg-white border-b border-slate-200 px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <Button
              variant="ghost"
              size="sm"
              onClick={handleCancelClockIn}
              className="p-0 hover:bg-transparent"
            >
              <ArrowLeft className="w-5 h-5" />
            </Button>
            <div className="w-8 h-8 bg-teal-600 rounded flex items-center justify-center text-white font-bold">
              C
            </div>
            <h1 className="text-xl font-semibold text-slate-900">Careviah</h1>
          </div>
          <div className="flex items-center space-x-3">
            <span className="text-sm text-slate-600">Admin A</span>
            <Avatar className="w-8 h-8">
              <AvatarImage src="/api/placeholder/32/32" />
              <AvatarFallback>AA</AvatarFallback>
            </Avatar>
          </div>
        </div>
      </header>

      <div className="flex-1 p-6">
        <div className="max-w-md mx-auto space-y-6">
          <div className="flex items-center space-x-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigate('/')}
              className="p-0 hover:bg-transparent"
            >
              <ArrowLeft className="w-5 h-5 text-slate-600" />
            </Button>
            <h2 className="text-lg font-semibold text-slate-900">Clock-Out</h2>
          </div>

          <div className="text-center">
            <div className="text-3xl font-mono font-bold text-slate-900">
              {formatTimer()}
            </div>
          </div>

          <div className="text-center space-y-3">
            <h3 className="text-xl font-semibold text-slate-900">Service Name A</h3>
            
            <div className="flex items-center justify-center space-x-3">
              <Avatar className="w-10 h-10">
                <AvatarImage src="/api/placeholder/40/40" />
                <AvatarFallback className="bg-slate-100 text-slate-600">
                  {getInitials(schedule.client_name)}
                </AvatarFallback>
              </Avatar>
              <span className="font-medium text-slate-900">{schedule.client_name}</span>
            </div>
          </div>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base">Tasks:</CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              {activities.map((activity) => (
                <div key={activity.id} className="space-y-3">
                  <h4 className="font-medium text-slate-900">{activity.title}</h4>
                  <p className="text-sm text-slate-600 leading-relaxed">
                    {activity.description}
                  </p>
                  
                  <div className="flex items-center space-x-4">
                    <button
                      onClick={() => handleActivityChange(activity.id, true)}
                      className={`flex items-center space-x-2 px-3 py-2 rounded-md border transition-colors cursor-pointer ${
                        activityStates[activity.id]?.isResolved
                          ? 'bg-green-100 border-green-300 text-green-800'
                          : 'bg-white border-slate-300 text-slate-600 hover:bg-slate-50'
                      }`}
                    >
                      <Check className="w-4 h-4" />
                      <span className="text-sm font-medium">Yes</span>
                    </button>
                    
                    <button
                      onClick={() => handleActivityChange(activity.id, false)}
                      className={`flex items-center space-x-2 px-3 py-2 rounded-md border transition-colors cursor-pointer ${
                        activityStates[activity.id]?.isResolved === false
                          ? 'bg-red-100 border-red-300 text-red-800'
                          : 'bg-white border-slate-300 text-slate-600 hover:bg-slate-50'
                      }`}
                    >
                      <X className="w-4 h-4" />
                      <span className="text-sm font-medium">No</span>
                    </button>
                  </div>

                  {/* Reason field - shown when "No" is selected */}
                  {activityStates[activity.id]?.isResolved === false && (
                    <div className="mt-3">
                      <textarea
                        value={activityStates[activity.id]?.reason || ''}
                        onChange={(e) => handleReasonChange(activity.id, e.target.value)}
                        placeholder="Add reason..."
                        className="w-full px-3 py-2 border border-slate-300 rounded-md text-sm placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-teal-500"
                        rows={2}
                      />
                    </div>
                  )}
                </div>
              ))}
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base">Clock-In Location</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-start space-x-3">
                <div className="w-16 h-16 bg-slate-200 rounded-lg flex-shrink-0">
                  <div className="w-full h-full bg-slate-300 rounded-lg flex items-center justify-center">
                    <MapPin className="w-6 h-6 text-slate-500" />
                  </div>
                </div>
                <div className="text-sm text-slate-600">
                 <div>{schedule.latitude}</div>
                 <div>{schedule.longitude}</div>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base">Service Notes</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-slate-600 leading-relaxed">
                Lorem ipsum dolor sit amet consectetur. Praesent adipiscing malesuada est vestibulum leo tempus sociis. Sodales libero mauris eu donec tempor in sagittis urna turpis. Vitae et vestibulum convallis volutpat commodo blandit in fusce viverra. Semper magna amet ipsum massa turpis non tortor. Etiam diam neque tristique nulla. Ipsum duis praesent sed a mattis morbi morbi aliquam. Enim quam amet cras nibh. Amet quis malesuada ac in ultrices. Viverra sagittis aenean vulputate at orci aliquam enim.
              </p>
            </CardContent>
          </Card>

          <div className="flex space-x-3">
            <Button
              onClick={handleCancelClockIn}
              variant="outline"
              className="flex-1 text-red-600 border-red-300 hover:bg-red-50"
            >
              Cancel Clock-In
            </Button>
            <Button
              onClick={handleClockOut}
              disabled={isSubmitting}
              className="flex-1 bg-teal-600 hover:bg-teal-700 text-white"
            >
              {isSubmitting ? 'Clocking Out...' : 'Clock-Out'}
            </Button>
          </div>

          <div className="text-center text-xs text-slate-400 pt-4">
            @2025 Careviah, Inc. All rights reserved.
          </div>
        </div>
      </div>

      <ClockOutSuccessModal
        open={showSuccessModal}
        onOpenChange={setShowSuccessModal}
        onGoHome={handleGoHome}
        duration={formatDuration()}
      />
    </div>
  );
} 