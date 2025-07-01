import { Card, CardContent } from '@/components/ui/card';
import type { Stats } from '@/lib/schemas';

interface StatsCardsProps {
  stats?: Stats;
  isLoading: boolean;
}

export function StatsCards({ stats, isLoading }: StatsCardsProps) {
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[1, 2, 3].map((i) => (
          <Card key={i} className="animate-pulse">
            <CardContent className="p-6">
              <div className="space-y-2">
                <div className="h-8 bg-slate-200 rounded w-16"></div>
                <div className="h-4 bg-slate-200 rounded w-32"></div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  const statsData = [
    {
      value: stats?.missed_schedules || 0,
      label: 'Missed Scheduled',
      color: 'text-red-600',
      bgColor: 'bg-red-50',
    },
    {
      value: stats?.upcoming_today_schedules || 0,
      label: "Upcoming Today's Schedule",
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
    },
    {
      value: stats?.completed_today_schedules || 0,
      label: "Today's Completed Schedule",
      color: 'text-green-600',
      bgColor: 'bg-green-50',
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      {statsData.map((stat, index) => (
        <Card key={index} className={`${stat.bgColor} border-0`}>
          <CardContent className="p-6">
            <div className="space-y-2">
              <div className={`text-3xl font-bold ${stat.color}`}>
                {stat.value}
              </div>
              <div className="text-sm text-slate-600">
                {stat.label}
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
} 