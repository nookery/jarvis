<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    return response(config('app.name'));
});

Route::get('/why', function (Request $request) {
    $result = collect([
        'getClientIp' => $request->getClientIp(),
        'getClientIps' => $request->getClientIps(),
        'getBaseUrl' => $request->getBaseUrl(),
        'getHost' => $request->getHost(),
        'getPathInfo' => $request->getPathInfo(),
        'getScheme' => $request->getScheme(),
        'getSchemeAndHttpHost' => $request->getSchemeAndHttpHost(),
        'HTTP_X_FORWARDED_PROTO' => $request->server('HTTP_X_FORWARDED_PROTO'),
        'HTTP_X_FORWARDED_FOR' => $request->server('HTTP_X_FORWARDED_FOR'),
        'HTTP_X_FORWARDED_PORT' => $request->server('HTTP_X_FORWARDED_PORT'),
        'HTTP_X_REAL_IP' => $request->server('HTTP_X_REAL_IP'),
        'SERVER_SOFTWARE' => $request->server('SERVER_SOFTWARE'),
        'REQUEST_SCHEME' => $request->server('REQUEST_SCHEME'),
        'SERVER_PORT' => $request->server('SERVER_PORT'),
        'SERVER_ADDR' => $request->server('SERVER_ADDR'),
        'SERVER' => $request->server(),
    ]);

    return response()->json($result);
});
