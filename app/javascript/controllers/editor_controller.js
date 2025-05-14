import "prismjs";
import { Controller } from "@hotwired/stimulus";

Prism.manual = true;

export default class extends Controller {
  static targets = ["language", "editor", "overlay"];

  rerender() {
    this.overlayTarget.innerHTML = Prism.highlight(
      this.editorTarget.value,
      Prism.languages[this.languageTarget.value],
      this.languageTarget.value,
    );
  }
}
