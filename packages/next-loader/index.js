/**
 * Optivor custom loader function for Next.js image optimization
 */
function optivorLoader({ src, width, quality }) {
  const baseUrl = (typeof process !== 'undefined' && process.env && process.env.NEXT_PUBLIC_OPTIVOR_URL) || 'http://localhost:8080';
  const cleanSrc = src.startsWith('/') ? src.slice(1) : src;

  const params = new URLSearchParams();
  if (width) params.set('w', width);
  if (quality) params.set('q', quality);
  params.set('format', 'webp');

  const queryString = params.toString();
  return `${baseUrl.replace(/\/$/, '')}/image/${cleanSrc}${queryString ? '?' + queryString : ''}`;
}

/**
 * Optivor React Image component wrapping next/image for zero-config integration
 */
function Image(props) {
  let React, NextImage;
  try {
    React = require('react');
    NextImage = require('next/image').default;
  } catch (e) {
    throw new Error('@optivor/next: react and next/image must be installed in your project to use the Image component.');
  }

  return React.createElement(NextImage, Object.assign({}, props, { loader: optivorLoader }));
}

module.exports = {
  optivorLoader,
  Image,
  default: Image
};
