defmodule RateCalculator.Repo do
  use Ecto.Repo,
    otp_app: :rate_calculator,
    adapter: Ecto.Adapters.Postgres
end
