import * as Sentry from '@sentry/nuxt';

Sentry.init({
    dsn: 'https://6b4b0223d29887555b3c5450ccd1a731@sentry.np-internal.ru/115',
    tracesSampleRate: 1.0,
});
