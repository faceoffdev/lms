import Aura from '@primevue/themes/aura';

export default defineNuxtConfig({
  app: {
    head: {
      title: "LMS",
    },
    link: [{ rel: "icon", type: "image/x-icon", href: "/favicon.ico" }],
  },

  runtimeConfig: {
    public: {
      API_URL: process.env.API_URL,
    },
  },

  devtools: { enabled: true },
  ssr: false,

  css: [
    "@/assets/main.css",
    "primeicons/primeicons.css",
  ],

  modules: ["@primevue/nuxt-module"],

  primevue: {
    options: {
      theme: {
        preset: Aura,
        options: {
          darkModeSelector: 'system',
        }
      },
      ripple: true
    },
    autoImport: true
  },

  compatibilityDate: "2025-01-06"
});