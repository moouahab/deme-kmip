import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/kmip": "http://localhost:8080",
      "/metrics": "http://localhost:8080",
      "/health": "http://localhost:8080",
      "/keys": "http://localhost:8080",
      "/audit": "http://localhost:8080",
    },
  },
});