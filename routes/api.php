<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| is assigned the "api" middleware group. Enjoy building your API!
|
*/

Route::get('/instruction', function () {
    return response()->json([
        'code' => 0,
        'message' => '查询成功',
        'data' => [
            'code' => config('app.code'),
            'name' => config('app.name'),
            'description' => config('app.description'),
        ],
    ]);
});

Route::get('/', function (Request $request) {
    $set = $request->input('set', $request->getClientIp());

    return response()->json([
        'code' => 0,
        'message' => '查询成功',
        'data' => [
            'ip' => $set,
        ],
    ]);
});
