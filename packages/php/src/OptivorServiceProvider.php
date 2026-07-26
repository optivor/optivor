<?php

namespace Optivor;

use Illuminate\Support\ServiceProvider;
use Illuminate\Support\Facades\Blade;

class OptivorServiceProvider extends ServiceProvider
{
    /**
     * Register any application services.
     */
    public function register()
    {
        $this->mergeConfigFrom(__DIR__ . '/../config/optivor.php', 'optivor');

        $this->app->singleton(OptivorClient::class, function ($app) {
            $config = $app['config']['optivor'] ?? [];
            return new OptivorClient([
                'baseUrl' => $config['base_url'] ?? env('OPTIVOR_BASE_URL', 'http://localhost:8080'),
                'securityKey' => $config['security_key'] ?? env('OPTIVOR_SECURITY_KEY', null),
                'defaultBucket' => $config['default_bucket'] ?? env('OPTIVOR_DEFAULT_BUCKET', 'default'),
            ]);
        });
    }

    /**
     * Bootstrap any application services.
     */
    public function boot()
    {
        if ($this->app->runningInConsole()) {
            $this->publishes([
                __DIR__ . '/../config/optivor.php' => config_path('optivor.php'),
            ], 'optivor-config');
        }

        // Register @optivor Blade directive
        Blade::directive('optivor', function ($expression) {
            return "<?php echo app(\\Optivor\\OptivorClient::class)->url($expression); ?>";
        });
    }
}
