const { h, defineComponent } = require('vue');
const { OptivorClient } = require('@optivor/js');

const OptivorImage = defineComponent({
  name: 'OptivorImage',
  props: {
    src: { type: String, required: true },
    width: { type: Number, default: undefined },
    height: { type: Number, default: undefined },
    fit: { type: String, default: 'cover' },
    format: { type: String, default: 'webp' },
    focal: { type: [Array, String], default: undefined },
    overlay: { type: String, default: undefined },
    gravity: { type: String, default: undefined },
    opacity: { type: Number, default: undefined },
    blur: { type: Number, default: undefined },
    grayscale: { type: Boolean, default: false },
    pixelate: { type: Number, default: undefined },
    baseUrl: { type: String, default: undefined },
    alt: { type: String, default: '' }
  },
  setup(props, { attrs }) {
    const client = new OptivorClient({ baseUrl: props.baseUrl });
    return () => {
      const url = client.buildUrl(props.src, {
        width: props.width,
        height: props.height,
        fit: props.fit,
        format: props.format,
        focal: props.focal,
        overlay: props.overlay,
        gravity: props.gravity,
        opacity: props.opacity,
        blur: props.blur,
        grayscale: props.grayscale,
        pixelate: props.pixelate
      });

      return h('img', {
        src: url,
        width: props.width,
        height: props.height,
        alt: props.alt,
        loading: 'lazy',
        ...attrs
      });
    };
  }
});

module.exports = {
  OptivorImage,
  default: OptivorImage
};
