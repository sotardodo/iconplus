# =========================
# Stage 1 - Build dependencies
# =========================
FROM composer:2 AS build

WORKDIR /app

RUN apk add --no-cache \
    zip unzip git libzip-dev icu-dev zlib-dev oniguruma-dev libxml2-dev && \
    docker-php-ext-install zip intl

COPY . .

RUN cp .env.example .env || true

# Install dependencies Laravel
RUN composer install --no-dev --no-interaction --prefer-dist --optimize-autoloader

# (Opsional) Tambahan package renoki exporter
RUN composer require renoki-co/laravel-exporter:^3.0 --no-interaction --prefer-dist || true

# Set permission folder Laravel
RUN chmod -R 775 storage bootstrap/cache


# =========================
# Stage 2 - PHP Runtime
# =========================
FROM php:8.2-fpm-alpine

WORKDIR /var/www/html

RUN apk add --no-cache \
    bash git unzip curl supervisor \
    icu-dev libzip-dev zlib-dev oniguruma-dev \
    libpng-dev libjpeg-turbo-dev libxpm-dev libwebp-dev libxml2-dev \
    mysql-client

RUN docker-php-ext-install pdo_mysql intl opcache zip mbstring bcmath

COPY --from=build /app ./


COPY ./supervisord.conf /etc/supervisord.conf

# Set environment
ENV APP_ENV=production \
    APP_DEBUG=false \
    APP_URL=http://localhost:8001 \
    PORT=8001

EXPOSE 8001

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
