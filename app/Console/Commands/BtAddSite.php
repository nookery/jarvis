<?php

/**
 * API文档：https://www.bt.cn/api-doc.pdf
 */

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\Http;

class BtAddSite extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'bt:add-site {form}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '操作宝塔面板，增加网站';

    /**
     * Execute the console command.
     *
     * @return int
     */
    public function handle()
    {
        parse_str($this->argument('form'), $form);

        $url = config('bt.bt_url') . config('bt.WebAddSite');
        $data = $this->patchSign([
            'webname' => json_encode([
                'domain' => $form['domain'],
                'domainlist' => [],
                'count' => 0,
            ]),
            'path' => $form['path'],
            'type_id' => 0,
            'type' => 'PHP',
            'version' => '80',
            'ps' => $form['ps'] ?? $form['domain'],
            'ftp' => false,
            'sql' => false,
            'port' => 80,
        ]);

        $this->info('请求链接：' . $url);
        $this->info("请求数据：\r\n" . json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));

        $response = Http::asForm()->post($url, $data);
        $body = $response->body();
        $output = $body;
        if (json_decode($output)) {
            $output = json_encode(json_decode($output), JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE);
        }

        $this->info('HTTP状态码：' . $response->status());
        $this->info("返回的内容：\r\n" . $output);
    }

    /**
     * 补充签名
     *
     * @param array $data
     * @return void
     */
    public function patchSign($data = [])
    {
        $time = time();
        $sign = [
            'request_token' => md5($time . '' . md5(config('bt.bt_key'))),
            'request_time' => $time,
        ];

        return array_merge($data, $sign);
    }
}
