<?php

namespace App\Http\Middleware;

use App\Services\MetricsRegistry;
use Closure;
use Illuminate\Http\Request;

class PrometheusMiddleware
{
    public function handle(Request $request, Closure $next)
    {
        $response = $next($request);

        if ($request->path() !== 'metrics') {
            MetricsRegistry::increment(
                $request->method(),
                '/' . ltrim($request->path(), '/'),
                (string) $response->getStatusCode()
            );
        }

        return $response;
    }
}
