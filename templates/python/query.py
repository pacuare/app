import requests
import pandas as pd

def query(sql: str, params: list = [], base = 'https://app.pacuare.dev') -> pd.DataFrame:
        """Query your database using SQL, returning a Pandas DataFrame of the results."""
        res = requests.post(base + '/api/query',
                            headers={
                                'Content-Type': 'application/json',
                                'Accept': 'application/json'
                            },
                            json={
                                'query': sql,
                                'params': params
                            })

        return pd.DataFrame.from_records(res.json())