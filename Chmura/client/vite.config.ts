import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    host: true,
    allowedHosts: [
      'frontend-app-c2a5baggb6euhehq.polandcentral-01.azurewebsites.net'
    ]
  }
});
