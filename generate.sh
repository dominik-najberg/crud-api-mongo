#!/bin/bash

protoc --go_out=plugins=grpc:. blog/blogpb/blog.proto