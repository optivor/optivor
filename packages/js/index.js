class OptivorClient {
  constructor(options = {}) {
    this.baseUrl = (options.baseUrl || 'http://localhost:8080').replace(/\/$/, '');
    this.defaultBucket = options.defaultBucket || '';
    this.securityKey = options.securityKey || '';
  }

  buildUrl(key, params = {}) {
    const cleanKey = key.startsWith('/') ? key.slice(1) : key;
    const fullPath = this.defaultBucket && !cleanKey.includes('/') 
      ? `${this.defaultBucket}/${cleanKey}` 
      : cleanKey;

    const query = new URLSearchParams();
    if (params.width || params.w) query.set('w', params.width || params.w);
    if (params.height || params.h) query.set('h', params.height || params.h);
    if (params.fit) query.set('fit', params.fit);
    if (params.format) query.set('format', params.format);
    if (params.focal) {
      const focalStr = Array.isArray(params.focal) ? params.focal.join(',') : params.focal;
      query.set('focal', focalStr);
    }
    if (params.overlay) query.set('overlay', params.overlay);
    if (params.gravity) query.set('gravity', params.gravity);
    if (params.opacity) query.set('opacity', params.opacity);
    if (params.overlayScale) query.set('overlay_scale', params.overlayScale);
    if (params.blur) query.set('blur', params.blur);
    if (params.grayscale) query.set('grayscale', 'true');
    if (params.pixelate) query.set('pixelate', params.pixelate);

    const queryString = query.toString();
    return `${this.baseUrl}/image/${fullPath}${queryString ? '?' + queryString : ''}`;
  }
}

module.exports = {
  OptivorClient,
  default: OptivorClient
};
