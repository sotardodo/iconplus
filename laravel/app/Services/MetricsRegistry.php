<?php

namespace App\Services;

class MetricsRegistry
{
    private static array $requestCounts = [];

    public static function increment(string $method, string $path, string $status): void
    {
        $key = $method . '|' . $path . '|' . $status;
        self::$requestCounts[$key] = (self::$requestCounts[$key] ?? 0) + 1;
    }

    public static function render(): string
    {
        $lines = [
            '# HELP up Laravel service status',
            '# TYPE up gauge',
            'up{service="laravel"} 1',
            '# HELP http_requests_total Total HTTP requests per endpoint',
            '# TYPE http_requests_total counter',
        ];

        foreach (self::$requestCounts as $key => $count) {
            [$method, $path, $status] = explode('|', $key, 3);
            $lines[] = sprintf(
                'http_requests_total{method="%s",path="%s",status="%s"} %d',
                $method,
                self::escapeLabel($path),
                $status,
                $count
            );
        }

        return implode("\n", $lines) . "\n";
    }

    private static function escapeLabel(string $value): string
    {
        return str_replace(['\\', '"'], ['\\\\', '\\"'], $value);
    }
}
