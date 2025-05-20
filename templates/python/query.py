from pyodide.http import pyfetch
import json
import pandas as pd

async def sql(query: str, base: str = 'https://app.pacuare.dev') -> pd.DataFrame:
    res = await pyfetch(base + '/api/query?language=sql',
                        method='post',
                        headers={
                            'Content-Type': 'text/sql',
                            'Accept': 'application/json'
                        },
                        body=query)
    
    return pd.DataFrame.from_records(await res.json())