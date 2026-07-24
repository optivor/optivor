/**
 * Custom loader for Next.js image optimization component (next/image) targeting Optivor.
 *
 * @param {Object} options
 * @param {string} options.src - Image path or URL
 * @param {number} options.width - Requested image width
 * @param {number} [options.quality] - Requested image quality
 * @returns {string} Fully qualified Optivor image URL
 */
module.exports = function optivorLoader({ src, width, quality }) {
  const baseUrl = process.env.NEXT_PUBLIC_OPTIVOR_URL || 'http://localhost:8080';
  const cleanSrc = src.startsWith('/') ? src.slice(1) : src;
  
  const params = new URLSearchParams();
  if (width) params.set('w', width);
  if (quality) params.set('q', quality);
  params.set('format', 'webp');

  const queryString = params.toString();
  return `${baseUrl.replace(/\/$/, '')}/image/${cleanSrc}${queryString ? '?' + queryString : ''}`;
};
