import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/lib/api';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { MoreHorizontal } from 'lucide-react';
import { StatsCards } from './StatsCards';
import { ScheduleItem } from './ScheduleItem';
import { ActiveVisit } from './ActiveVisit';

export function Dashboard() {
  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ['stats'],
    queryFn: () => apiClient.getStats(),
  });

  const { data: todaySchedules, isLoading: schedulesLoading } = useQuery({
    queryKey: ['schedules', 'today'],
    queryFn: () => apiClient.getTodaySchedules(),
  });

  // Find active visit (schedule with status 'in_progress')
  const activeVisit = todaySchedules?.find(schedule => schedule.status === 'in_progress');

  return (
    <div className="flex flex-col min-h-screen bg-slate-50">
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

      {activeVisit && <ActiveVisit schedule={activeVisit} />}

      <div className="flex-1 p-6">
        <div className="max-w-7xl mx-auto space-y-6">
          <h2 className="text-2xl font-bold text-slate-900">Dashboard</h2>

          <StatsCards 
            stats={stats} 
            isLoading={statsLoading} 
          />

          <div className="bg-white rounded-lg border border-slate-200">
            <div className="px-6 py-4 border-b border-slate-200">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <h3 className="text-lg font-semibold text-slate-900">Schedule</h3>
                  <Badge variant="secondary" className="bg-blue-100 text-blue-800">
                    {todaySchedules?.length || 0}
                  </Badge>
                </div>
                <Button variant="ghost" size="sm">
                  <MoreHorizontal className="w-4 h-4" />
                </Button>
              </div>
            </div>

            <div className="p-6 space-y-4">
              {schedulesLoading ? (
                <div className="space-y-4">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className="animate-pulse">
                      <div className="flex items-center space-x-4">
                        <div className="w-12 h-12 bg-slate-200 rounded-full"></div>
                        <div className="flex-1 space-y-2">
                          <div className="h-4 bg-slate-200 rounded w-1/4"></div>
                          <div className="h-3 bg-slate-200 rounded w-1/2"></div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : todaySchedules && todaySchedules.length > 0 ? (
                todaySchedules.map((schedule) => (
                  <ScheduleItem key={schedule.id} schedule={schedule} />
                ))
              ) : (
                <div className="text-center py-8 text-slate-500">
                  No schedules for today
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
} 