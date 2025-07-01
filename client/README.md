# Visit Tracker - Client Application

A React-based client application for tracking caregiver visits and schedules with mandatory geolocation support.

## 🌍 Geolocation Requirements

**⚠️ IMPORTANT: Geolocation access is mandatory for this application to function properly.**

This application requires access to your device's location services for:
- Clocking in/out of visits
- Location verification during visits
- Accurate visit tracking and reporting

### Geolocation Support

The application uses the browser's built-in Geolocation API with the following features:
- **High accuracy positioning** for precise location tracking
- **Comprehensive error handling** for all geolocation scenarios
- **User-friendly error messages** explaining how to resolve location issues

### Browser Requirements

- **Modern browser support**: Chrome 5+, Firefox 3.5+, Safari 5+, Edge 12+
- **HTTPS required**: Geolocation API requires secure context (except localhost)
- **Location permissions**: User must grant location access when prompted

### Fallback Handling

If geolocation is unavailable, the application will:

1. **Permission Denied**: Display instructions to enable location permissions in browser settings
2. **Not Supported**: Alert user to use a modern browser with geolocation support
3. **Position Unavailable**: Guide user to enable location services on their device
4. **Timeout**: Suggest checking internet connection and trying again

**No manual location input is available as a fallback - geolocation access is required.**

## 🚀 Setup Instructions

### 1. Environment Configuration

Create environment files:

```bash
# Copy the example environment file
cp .env.example .env

# The .env file contains:
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 2. Install Dependencies

```bash
npm install
```

### 3. Development Server

```bash
npm run dev
```

### 4. Enable Location Services

**Before using the application:**

1. **Browser Permissions**: When prompted, click "Allow" to grant location access
2. **Device Settings**: Ensure location services are enabled on your device
3. **HTTPS**: For production deployment, ensure the app is served over HTTPS

## 🔧 Troubleshooting Geolocation Issues

### Common Issues and Solutions

**"Location access denied"**
- Go to browser settings → Privacy & Security → Site Settings → Location
- Find your app's URL and change permission to "Allow"
- Reload the page

**"Geolocation not supported"**
- Update to a modern browser version
- Ensure you're not using a very old browser or privacy-focused browser with location disabled

**"Position unavailable"**
- Check if location services are enabled on your device
- Try moving to an area with better GPS/network signal
- Restart your browser and try again

**"Location request timed out"**
- Check your internet connection
- Move to an area with better signal
- Try refreshing the page

### Testing Geolocation

For development/testing purposes:
- Chrome DevTools: Use "Sensors" tab to simulate location
- Firefox: Use "Developer Tools" → "Responsive Design Mode" → location simulation

## 🛠 Technical Stack

- **React 18** with TypeScript
- **Vite** for fast development and building
- **TanStack Query** for server state management
- **React Router** for navigation
- **Tailwind CSS** for styling
- **Zod** for runtime type validation
- **Lucide React** for icons

## 📱 Features

- Real-time visit tracking with geolocation
- Schedule management and navigation
- Activity/task completion tracking
- Dashboard with visit statistics
- Responsive mobile-first design

## 🔒 Security Notes

- Geolocation data is only used for visit verification
- Location coordinates are sent securely to the backend
- No location data is stored locally or shared with third parties
- All API calls include geolocation validation

## 🏗️ Build for Production

```bash
npm run build
```

## 🧪 Development Tools

- ESLint for code linting
- TypeScript for type checking
- Vite for fast HMR and building
