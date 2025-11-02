<?php

use Illuminate\Support\Facades\Route;
use RenokiCo\LaravelExporter\Http\Controllers\ExporterController;

Route::get('/health', function () {
    return response()->json(['status' => 'ok'], 200);
});

Route::get('/', function () {
    return view('welcome');
});

// Endpoint metrics untuk Prometheus
Route::get('/metrics', [ExporterController::class, 'metrics']);
