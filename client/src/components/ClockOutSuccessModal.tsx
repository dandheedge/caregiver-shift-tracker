import { Dialog, DialogContent } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Check, Calendar, Clock } from 'lucide-react';

interface ClockOutSuccessModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onGoHome: () => void;
  duration?: string;
}

export function ClockOutSuccessModal({ 
  open, 
  onOpenChange, 
  onGoHome,
  duration = "1 hour"
}: ClockOutSuccessModalProps) {
  const formatDate = () => {
    const now = new Date();
    return now.toLocaleDateString('en-US', {
      weekday: 'short',
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    });
  };

  const formatTime = () => {
    const now = new Date();
    return now.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    }) + ' SGT';
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="w-full max-w-sm mx-auto text-center p-8" onClose={() => onOpenChange(false)}>
        <div className="flex flex-col items-center space-y-6">
          {/* Success Icon */}
          <div className="relative">
            <div className="w-16 h-16 bg-orange-100 rounded-full flex items-center justify-center">
              <div className="w-12 h-12 bg-orange-500 rounded-full flex items-center justify-center">
                <Check className="w-6 h-6 text-white" />
              </div>
            </div>
            {/* Decorative elements */}
            <div className="absolute -top-2 -right-2 w-3 h-3 bg-orange-300 rounded-full"></div>
            <div className="absolute -bottom-1 -left-3 w-2 h-2 bg-orange-200 rounded-full"></div>
            <div className="absolute -top-3 left-1 w-1 h-1 bg-orange-400 rounded-full"></div>
          </div>

          {/* Title */}
          <h2 className="text-xl font-semibold text-slate-900">
            Schedule Completed
          </h2>

          {/* Date and Time Info */}
          <div className="space-y-3">
            <div className="flex items-center justify-center space-x-2 text-slate-600">
              <Calendar className="w-4 h-4" />
              <span className="text-sm">{formatDate()}</span>
            </div>
            
            <div className="flex items-center justify-center space-x-2 text-slate-600">
              <Clock className="w-4 h-4" />
              <span className="text-sm">{formatTime()}</span>
            </div>
            
            <div className="text-sm text-slate-500">
              ({duration})
            </div>
          </div>

          {/* Go to Home Button */}
          <Button
            onClick={onGoHome}
            className="w-full bg-teal-600 hover:bg-teal-700 text-white py-3 mt-6"
            size="lg"
          >
            Go to Home
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
} 