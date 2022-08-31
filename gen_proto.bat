@echo off
for /f %%I in ('dir protos\proto\*.proto /A-D /B /ON') do (
  echo generating %%I...
  protoc --proto_path=protos/proto --go_out=. %%I
)
