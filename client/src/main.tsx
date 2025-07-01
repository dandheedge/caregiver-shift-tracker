import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)

// Alternative: Remove StrictMode to eliminate duplicate API calls in development
// Note: StrictMode helps detect side effects and is recommended for development
// In production builds, StrictMode doesn't cause duplicate calls
// 
// createRoot(document.getElementById('root')!).render(
//   <App />
// )
