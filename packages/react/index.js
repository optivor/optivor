const React = require('react');
const { OptivorClient } = require('@optivor/js');

function OptivorImage({
  src,
  width,
  height,
  fit = 'cover',
  format = 'webp',
  focal,
  overlay,
  gravity,
  opacity,
  blur,
  grayscale,
  pixelate,
  baseUrl,
  alt = '',
  className,
  style,
  ...rest
}) {
  const client = new OptivorClient({ baseUrl });
  const imageUrl = client.buildUrl(src, {
    width,
    height,
    fit,
    format,
    focal,
    overlay,
    gravity,
    opacity,
    blur,
    grayscale,
    pixelate
  });

  const srcset = width ? [
    `${client.buildUrl(src, { width, height, fit, format, blur, grayscale })} 1x`,
    `${client.buildUrl(src, { width: width * 2, height: height ? height * 2 : undefined, fit, format, blur, grayscale })} 2x`
  ].join(', ') : undefined;

  return React.createElement('img', {
    src: imageUrl,
    srcSet: srcset,
    width,
    height,
    alt,
    className,
    style,
    loading: 'lazy',
    ...rest
  });
}

module.exports = {
  OptivorImage,
  default: OptivorImage
};
