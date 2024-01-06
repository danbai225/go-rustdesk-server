protoc --go_out=model/model_proto/ model/protos/*
go install github.com/99designs/gqlgen@latest
gqlgen