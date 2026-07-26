<?php

return [
    /*
    |--------------------------------------------------------------------------
    | Optivor Server Base URL
    |--------------------------------------------------------------------------
    |
    | The base URL of your deployed Optivor engine or Cloudflare Worker edge proxy.
    |
    */
    'base_url' => env('OPTIVOR_BASE_URL', 'http://localhost:8080'),

    /*
    |--------------------------------------------------------------------------
    | Security HMAC Key
    |--------------------------------------------------------------------------
    |
    | Optional secret key used to generate signed URLs and prevent tampering.
    |
    */
    'security_key' => env('OPTIVOR_SECURITY_KEY', null),

    /*
    |--------------------------------------------------------------------------
    | Default Bucket Identifier
    |--------------------------------------------------------------------------
    |
    | Default S3 / storage bucket identifier used when generating image URLs.
    |
    */
    'default_bucket' => env('OPTIVOR_DEFAULT_BUCKET', 'default'),
];
