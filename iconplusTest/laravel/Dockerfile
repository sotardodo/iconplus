# =========================
# Stage 1 - Build dependencies
# =========================
FROM composer:2 AS build

WORKDIR /app

COPY composer.json composer.lock ./
COPY . .

RUN composer install --no-dev --no-interaction --prefer-dist --optimize-autoloader

RUN if [ ! -f ".env" ]; then \
      echo "⚙️ Creating .env from .env.example..."; \
      cp .env.example .env; \
    fi

RUN set -eux; \
    if ! grep -q "APP_KEY=" .env || [ -z "$(grep 'APP_KEY=' .env | cut -d '=' -f2)" ]; then \
      echo "🔑 Generating Laravel APP_KEY ..."; \
      php artisan key:generate --force; \
    else \
      echo "✅ APP_KEY already exists, skipping generation."; \
    fi

# =========================
# Stage 2 - Runtime
# =========================
FROM php:8.2-fpm-alpine

WORKDIR /var/www/html

RUN apk add --no-cache bash git unzip libzip-dev icu-dev oniguruma-dev \
    && docker-php-ext-install pdo_mysql intl opcache zip

COPY --from=build /app .

RUN chown -R www-data:www-data /var/www/html \
    && chmod -R 775 storage bootstrap/cache

ENV APP_ENV=production \
    APP_DEBUG=false

EXPOSE 8001

CMD ["sh", "-c", "php artisan migrate --force || true && php artisan serve --host=0.0.0.0 --port=8001"]
