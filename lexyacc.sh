#!/bin/bash
nex line.nex && go tool yacc -o line.yy.go -v '' line.y && echo "OK"