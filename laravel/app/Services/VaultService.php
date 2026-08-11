<?php

namespace App\Services;

use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Log;

class VaultService
{
    private string $addr;

    private string $role;

    public function __construct()
    {
        $this->addr = rtrim((string) env('VAULT_ADDR', ''), '/');
        $this->role = (string) env('VAULT_ROLE', 'laravel-app');
    }

    public function isEnabled(): bool
    {
        return $this->addr !== '';
    }

    public function getSecrets(string $path): array
    {
        if (!$this->isEnabled()) {
            return [];
        }

        try {
            $token = $this->resolveToken();
            if ($token === '') {
                return [];
            }

            $response = Http::timeout(5)
                ->withToken($token)
                ->get($this->addr . '/v1/secret/data/' . $path);

            if (!$response->successful()) {
                Log::warning('Vault secret read failed', [
                    'path' => $path,
                    'status' => $response->status(),
                ]);

                return [];
            }

            return $response->json('data.data') ?? [];
        } catch (\Throwable $e) {
            Log::warning('Vault secret read error', [
                'path' => $path,
                'error' => $e->getMessage(),
            ]);

            return [];
        }
    }

    private function resolveToken(): string
    {
        $staticToken = (string) env('VAULT_TOKEN', '');
        if ($staticToken !== '') {
            return $staticToken;
        }

        $jwtPath = '/var/run/secrets/kubernetes.io/serviceaccount/token';
        if (!is_readable($jwtPath)) {
            return '';
        }

        $response = Http::timeout(5)->post($this->addr . '/v1/auth/kubernetes/login', [
            'role' => $this->role,
            'jwt' => file_get_contents($jwtPath),
        ]);

        if (!$response->successful()) {
            Log::warning('Vault kubernetes login failed', [
                'status' => $response->status(),
            ]);

            return '';
        }

        return (string) $response->json('auth.client_token', '');
    }
}
