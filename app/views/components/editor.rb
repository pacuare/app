class Views::Components::Editor < Views::Base
  def view_template
    div(data_controller: 'editor') {
      textarea(data_editor_target: 'editor', data_action: 'input->editor#rerender')
      pre { code(data_editor_target: 'overlay') }
      select(data_editor_target: 'language', data_action: 'change->editor#rerender') {
        option(value: 'sql') { 'SQL' }
        option(value: "python") { "Python" }
      }
    }
  end
end
