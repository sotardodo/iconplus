# =========================
# Stage 1 - Build dependencies
# =========================
FROM composer:2 AS build

WORKDIR /app

# Copy composer files and source
COPY composer.json composer.lock ./
COPY . .

# Install dependencies
RUN composer install --no-dev --no-interaction --prefer-dist --optimize-autoloader

# Generate APP_KEY if missing
# (we need PHP CLI, so run artisan command after vendor installed)
RUN if ! grep -q "APP_KEY=" .env || [ -z "$(grep 'APP_KEY=' .env | cut -d '=' -f2)" ]; then \
      echo "🛠 Generating Laravel APP_KEY ..."; \
      php artisan key:generate --force; \
    else \
      echo "✅ APP_KEY already exists, skipping generation."; \
    fi

# =========================
# Stage 2 - Runtime
# =========================
FROM php:8.2-fpm-alpine

WORKDIR /var/www/html

# Install PHP extensions
RUN apk add --no-cache bash git unzip libzip-dev icu-dev oniguruma-dev \
    && docker-php-ext-install pdo_mysql intl opcache zip

# Copy app (including vendor from previous stage)
COPY --from=build /app ./ 

# Set permission
RUN chown -R www-data:www-data /var/www/html

# Expose port and start Laravel
EXPOSE 8001

# Run Laravel server
CMD ["php", "artisan", "serve", "--host=0.0.0.0", "--port=8001"]

