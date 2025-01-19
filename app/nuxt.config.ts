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

  modules: ["@primevue/nuxt-module", "@sentry/nuxt/module"],

  sentry: {
    sourceMapsUploadOptions: {
      org: "faceoff",
      project: "lms",
      authToken: "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",
    },
  },

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