import gradio as gr
import query

async def process():
  return (
    (await query.sql("select injuries, turtle_occurrences from unique_turtles"))
    .apply(
      lambda row: {**row, 'n_injuries': len(row['injuries'].split(','))},
      axis=1,
      result_type='expand'
    )
  )

with gr.Blocks() as demo:
    gr.ScatterPlot(
        value=await process(),
        x='turtle_occurrences',
        y='n_injuries'
    )

demo.launch()