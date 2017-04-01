#!/bin/bash
nex kvl.nex && go tool yacc -o kvl.yy.go -v '' kvl.y && echo "OK"