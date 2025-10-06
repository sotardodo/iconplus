# Stage 1: composer install
FROM php:8.2-fpm-alpine AS vendor
WORKDIR /var/www/html
RUN apk add --no-cache bash git unzip libzip-dev icu-dev oniguruma-dev \
    && docker-php-ext-install pdo_mysql intl opcache zip
COPY composer.json composer.lock ./
RUN php -r "copy('https://getcomposer.org/installer','composer-setup.php');" \
    && php composer-setup.php --install-dir=/usr/local/bin --filename=composer \
    && rm composer-setup.php
ENV COMPOSER_ALLOW_SUPERUSER=1
RUN composer install --no-dev --no-interaction --optimize-autoloader --prefer-dist --ignore-platform-reqs

# Stage 2: app image
FROM php:8.2-fpm-alpine
WORKDIR /var/www/html
RUN apk add --no-cache bash libzip
# copy vendor
COPY --from=vendor /var/www/html/vendor ./vendor
# copy app
COPY . .
# set permissions (adjust as needed)
RUN chown -R www-data:www-data /var/www/html
EXPOSE 9000
CMD ["php-fpm"]
