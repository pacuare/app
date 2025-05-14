class Views::Index < Views::Base
  def view_template
    div {
      h1 { "Hello" }
      ::Views::Components::Editor()
    }
  end
end
