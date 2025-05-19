import "prismjs";

function rerender(evt) {
  const me = evt.target.closest("[data-component=editor]"),
    editor = me.querySelector("[data-editor=editor]"),
    overlay = me.querySelector("[data-editor=overlay]"),
    language = me.querySelector("[data-editor=language]");

  overlay.innerHTML = Prism.highlight(
    editor.value,
    Prism.languages[language.value],
    language.value,
  );
}

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

    new ResizeObserver(() => {
      overlay.style.width = textArea.offsetWidth + 'px';
      overlay.style.height = textArea.offsetHeight + 'px';
    }).observe(textArea);

    textArea.addEventListener("scroll", e => {
      overlay.scrollTop = textArea.scrollTop;
      overlay.scrollLeft = textArea.scrollLeft;
    });

    component
      .querySelector("[data-editor=language]")
      .addEventListener("change", rerender);
  });
});
