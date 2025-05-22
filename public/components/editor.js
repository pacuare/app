Alpine.data('editor', () => ({
  editorContent: '',
  output: '',
  language: 'sql',
  overlayWidth: null,
  overlayHeight: null,

  overlay: {
    [':style']() {

    }
  },

  rerender() {
    this.output = hljs.highlight(
      this.editorContent,
      {language: this.language}
    ).value;
  },

  init() {
    new ResizeObserver(() => {
      this.overlayWidth = textArea.offsetWidth + 'px';
      this.overlayHeight = textArea.offsetHeight + 'px';
    }).observe(document.querySelector(this.$id('editor')));
    this.overlayWidth = textArea.offsetWidth + 'px';
    this.overlayHeight = textArea.offsetHeight + 'px';
  }
}));

export function language(component) {
  return component.querySelector("[data-editor=language]").value
}
export function query(component) {
  return component.querySelector("[data-editor=editor]").value
}

window.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll("[data-component=editor]").forEach((component) => {
    const textArea = component.querySelector("[data-editor=editor]");
    const overlay = component.querySelector("[data-editor=overlay]");
    textArea.addEventListener("input", rerender);

    textArea.addEventListener("scroll", e => {
      overlay.scrollTop = textArea.scrollTop;
      overlay.scrollLeft = textArea.scrollLeft;
    });

    component
      .querySelector("[data-editor=language]")
      .addEventListener("change", rerender);
  });
});

document.addEventListener('resume', () => {
  // in case a restored browser tab prefills it
  rerender({target: component});
})
