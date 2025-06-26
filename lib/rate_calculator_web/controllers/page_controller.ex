defmodule RateCalculatorWeb.PageController do
  use RateCalculatorWeb, :controller

  def home(conn, _params) do
    render(conn, :home)
  end
end
