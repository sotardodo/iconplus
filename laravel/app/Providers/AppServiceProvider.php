<?php

namespace App\Providers;

use App\Services\VaultService;
use Illuminate\Support\ServiceProvider;

class AppServiceProvider extends ServiceProvider
{
    /**
     * Register any application services.
     *
     * @return void
     */
    public function register()
    {
        //
    }

    /**
     * Bootstrap any application services.
     *
     * @return void
     */
    public function boot()
    {
        $vault = new VaultService();
        if (!$vault->isEnabled()) {
            return;
        }

        $secrets = $vault->getSecrets('laravel');
        if (isset($secrets['db_password']) && $secrets['db_password'] !== '') {
            config(['database.connections.mysql.password' => $secrets['db_password']]);
        }
        if (isset($secrets['app_key']) && $secrets['app_key'] !== '') {
            config(['app.key' => $secrets['app_key']]);
        }
    }
}
