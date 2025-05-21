import gradio as gr
from query import query

def process(id):
  sql = "select injuries, turtle_occurrences from unique_turtles where turtle_id = $1"
  params = [id]
  if id == '':
    sql = "select injuries, turtle_occurrences from unique_turtles"
    params = []
  return (
    query(sql, params)
    .apply(
      lambda row: {
        **row, 
        'n_injuries': len(
          [injury for injury in row['injuries'].split(',') if injury.strip() != '']
        )
      },
      axis=1,
      result_type='expand'
    )
  )

demo = gr.Interface(
  fn=process,
  inputs=["text"],
  outputs=[gr.ScatterPlot(x='turtle_occurrences', y='n_injuries')]
)

demo.launch()