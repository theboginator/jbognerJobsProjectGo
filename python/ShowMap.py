import plotly.express as px
import pandas as pd

us_cities = pd.read_csv("https://raw.githubusercontent.com/plotly/datasets/master/us-cities-top-1k.csv")

df = px.data.gapminder().query("year == 2007")
fig = px.scatter_geo(us_cities,
                     lat="lat", lon="lon", hover_name="City", hover_data=["State", "Population"],
                     color_discrete_sequence=["fuchsia"]
                     )
fig.show()
