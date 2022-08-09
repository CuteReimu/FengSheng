@echo off
protoc --proto_path=. --go_out=. fengsheng.proto
