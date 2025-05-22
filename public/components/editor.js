import Alpine from "alpinejs";

Alpine.data('editor', () => ({
  editorContent: '',
  output: '',
  language: 'sql',
  overlayWidth: null,
  overlayHeight: null,
  scrollTop: 0,
  scrollLeft: 0,

  overlay: {
    ['x-ref']: 'overlay',
  },

  runBtn: {
    ['@click']() {
      this.$store.results.runQuery(this.editorContent);
    }
  },

  editor: {
    ['x-model']: 'editorContent',
    ['x-ref']: 'editor',
    ['@scroll']() {
      const editor = this.$refs.editor;
      this.scrollTop = editor.scrollTop;
      this.scrollLeft = editor.scrollLeft;
    },
    ['@input']() {
      this.rerender()
    }
  },

  rerender() {
    this.output = window.hljs.highlight(
      this.editorContent,
      {language: this.language}
    ).value;
  },

  resize() {
    const {editor, overlay} = this.$refs;
    overlay.style.width = editor.offsetWidth + 'px';
    overlay.style.height = editor.offsetHeight + 'px';
  },

  init() {
    new ResizeObserver(this.resize.bind(this)).observe(editor);
    this.resize();
    this.rerender();
  }
}));