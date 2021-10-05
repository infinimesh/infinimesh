#!/bin/bash
#--------------------------------------------------------------------------
# Copyright 2018-2021
# www.infinimesh.io
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#--------------------------------------------------------------------------
FILENAME=$1

if [ ! -z $FILENAME ] 
then
    openssl genrsa -out $FILENAME.key 4096
    openssl req -new -x509 -sha256 -key $FILENAME.key -out $FILENAME.crt -days 365 -subj "/C=/ST=/L=/O=/CN=/emailAddress=/"
else
    echo "Please set <filename> without extension for first argument without blank spaces.";
    echo "./cert.sh <filename>";
fi
