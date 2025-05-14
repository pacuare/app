import "prismjs";

export function rerender(evt) {
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

window.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll("[data-component=editor]").forEach((component) => {
    component
      .querySelector("[data-editor=editor]")
      .addEventListener("input", rerender);
    component
      .querySelector("[data-editor=language]")
      .addEventListener("change", rerender);
  });
});
