database {
  host = "127.0.0.1"
  port = 5432
  user = "postgres"
  password = "postgres"
  database = "api"
  sslmode = "disable"
}

app {
  allowed_origins = ["*"]

  http_addr = ":5555"
  grpc_addr = ":5556"
}
