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
    component
      .querySelector("[data-editor=editor]")
      .addEventListener("input", rerender);
    component
      .querySelector("[data-editor=language]")
      .addEventListener("change", rerender);
  });
});
