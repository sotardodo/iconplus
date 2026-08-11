# ==========================================
# STAGE 1: Builder 
# ==========================================
FROM php:8.2-cli-alpine AS build

RUN apk add --no-cache zip unzip git

COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

WORKDIR /app

COPY composer.json composer.lock ./
RUN composer install --no-dev --no-interaction --prefer-dist --no-scripts --no-autoloader

COPY . .

RUN composer dump-autoload --optimize && \
    composer run-script post-root-package-install || true && \
    composer run-script post-create-project-cmd || true

# ==========================================
# STAGE 2: Final Image (FPM & Supervisor)
# ==========================================
FROM php:8.2-fpm-alpine

WORKDIR /var/www/html

RUN apk add --no-cache \
    bash git unzip curl supervisor tzdata mysql-client \
    icu-dev libzip-dev oniguruma-dev && \
    docker-php-ext-install pdo_mysql intl opcache zip mbstring bcmath

COPY --from=build --chown=www-data:www-data /app ./

COPY ./supervisord.conf /etc/supervisord.conf

RUN chmod -R 775 storage bootstrap/cache

ENV APP_ENV=production \
    APP_DEBUG=false \
    PORT=8001

EXPOSE 8001

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]