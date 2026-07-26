<?php

namespace Optivor;

class OptivorClient
{
    private string $baseUrl;
    private string $defaultBucket;

    public function __construct(string $baseUrl = 'http://localhost:8080', string $defaultBucket = '')
    {
        $this->baseUrl = rtrim($baseUrl, '/');
        $this->defaultBucket = $defaultBucket;
    }

    public function url(string $key, array $params = []): string
    {
        $cleanKey = ltrim($key, '/');
        $fullPath = ($this->defaultBucket && !str_contains($cleanKey, '/'))
            ? "{$this->defaultBucket}/{$cleanKey}"
            : $cleanKey;

        $queryParams = [];
        if (isset($params['width']) || isset($params['w'])) {
            $queryParams['w'] = $params['width'] ?? $params['w'];
        }
        if (isset($params['height']) || isset($params['h'])) {
            $queryParams['h'] = $params['height'] ?? $params['h'];
        }
        if (isset($params['fit'])) {
            $queryParams['fit'] = $params['fit'];
        }
        if (isset($params['format'])) {
            $queryParams['format'] = $params['format'];
        }
        if (isset($params['focal'])) {
            $queryParams['focal'] = is_array($params['focal']) ? implode(',', $params['focal']) : $params['focal'];
        }
        if (isset($params['overlay'])) {
            $queryParams['overlay'] = $params['overlay'];
        }
        if (isset($params['gravity'])) {
            $queryParams['gravity'] = $params['gravity'];
        }
        if (isset($params['opacity'])) {
            $queryParams['opacity'] = $params['opacity'];
        }
        if (isset($params['blur'])) {
            $queryParams['blur'] = $params['blur'];
        }
        if (!empty($params['grayscale'])) {
            $queryParams['grayscale'] = 'true';
        }
        if (isset($params['pixelate'])) {
            $queryParams['pixelate'] = $params['pixelate'];
        }

        $query = http_build_query($queryParams);
        return "{$this.baseUrl}/image/{$fullPath}" . ($query ? "?{$query}" : '');
    }
}
