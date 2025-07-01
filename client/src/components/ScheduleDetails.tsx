import { useQuery } from '@tanstack/react-query';
import { useParams, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ArrowLeft, Calendar, Clock, Mail, Phone, MapPin } from 'lucide-react';
import { apiClient } from '@/lib/api';

interface ScheduleDetailsParams {
  id: string;
  [key: string]: string | undefined;
}

export function ScheduleDetails() {
  const { id } = useParams<ScheduleDetailsParams>();
  const navigate = useNavigate();
  const scheduleId = parseInt(id || '0');

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

  const handleClockIn = () => {
    navigate(`/schedule/${scheduleId}/clock-out`);
  };

  if (scheduleLoading || activitiesLoading) {
    return (
      <div className="flex flex-col min-h-screen bg-slate-50">
        <div className="flex-1 p-6">
          <div className="max-w-md mx-auto">
            <div className="animate-pulse space-y-4">
              <div className="h-8 bg-slate-200 rounded w-1/3"></div>
              <div className="h-32 bg-slate-200 rounded"></div>
              <div className="h-24 bg-slate-200 rounded"></div>
              <div className="h-48 bg-slate-200 rounded"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!schedule) {
    return (
      <div className="flex flex-col min-h-screen bg-slate-50">
        <div className="flex-1 p-6">
          <div className="max-w-md mx-auto text-center">
            <p className="text-slate-500">Schedule not found</p>
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
      {/* Header */}
      <header className="bg-white border-b border-slate-200 px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            
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
        <div className="max-w-2xl mx-auto space-y-6">
          <div className="flex items-center space-x-2 cursor-pointer">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigate('/')}
              className="p-0 hover:bg-transparent"
            >
              <ArrowLeft className="w-5 h-5 text-slate-600" />
            </Button>
            <h2 className="text-lg font-semibold text-slate-900">Schedule Details</h2>
          </div>

          <div className="text-center space-y-4">
            <h3 className="text-xl font-semibold text-slate-900">Service Name A</h3>
            
            <div className="flex items-center justify-center space-x-3">
              <Avatar className="w-12 h-12">
                <AvatarImage src="/api/placeholder/48/48" />
                <AvatarFallback className="bg-slate-100 text-slate-600">
                  {getInitials(schedule.client_name)}
                </AvatarFallback>
              </Avatar>
              <span className="font-medium text-slate-900">{schedule.client_name}</span>
            </div>

            <div className="flex items-center justify-center space-x-8 bg-blue-50 py-3 px-4 rounded-lg">
              <div className="flex items-center space-x-2 text-blue-700">
                <Calendar className="w-4 h-4" />
                <span className="text-sm font-medium">{formatDate(schedule.start_time)}</span>
              </div>
              <div className="w-px h-4 bg-blue-200"></div>
              <div className="flex items-center space-x-2 text-blue-700">
                <Clock className="w-4 h-4" />
                <span className="text-sm font-medium">
                  {formatTime(schedule.start_time)} - {formatTime(schedule.end_time)}
                </span>
              </div>
            </div>
          </div>

          {/* Client Contact */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base">Client Contact:</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex items-center space-x-3">
                <Mail className="w-4 h-4 text-slate-500" />
                <span className="text-sm text-slate-600">melisa@gmail.com</span>
              </div>
              <div className="flex items-center space-x-3">
                <Phone className="w-4 h-4 text-slate-500" />
                <span className="text-sm text-slate-600">+44 1232 212 3233</span>
              </div>
            </CardContent>
          </Card>

          {/* Address */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base">Address:</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-start space-x-3">
                <MapPin className="w-4 h-4 text-slate-500 mt-0.5" />
                <div className="text-sm text-slate-600">
                  <div>{schedule.latitude}, {schedule.longitude}</div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Tasks */}
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base">Tasks:</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {activities && activities.length > 0 ? (
                activities.map((activity) => (
                  <div key={activity.id} className="space-y-2">
                    <h4 className="font-medium text-slate-900">{activity.title}</h4>
                    <p className="text-sm text-slate-600 leading-relaxed">
                      {activity.description}
                    </p>
                  </div>
                ))
              ) : (
                <p className="text-sm text-slate-500">No activities assigned</p>
              )}
            </CardContent>
          </Card>

          {/* Service Notes */}
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

          {/* Clock-In Button */}
          <Button
            onClick={handleClockIn}
            className="w-full bg-teal-600 hover:bg-teal-700 text-white py-3"
            size="lg"
          >
            Clock-In Now
          </Button>

          <div className="text-center text-xs text-slate-400 pt-4">
            @2025 Careviah, Inc. All rights reserved.
          </div>
        </div>
      </div>
    </div>
  );
} 