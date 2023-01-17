#!/bin/bash
mkdir empty
openssl verify -x509_strict -verbose -CAfile ca.crt -CApath empty -purpose sslclient -verify_hostname node1 cert.pem
openssl verify -x509_strict -verbose -CAfile ca.crt -CApath empty -purpose sslserver -verify_hostname node1 cert.pem
rmdir empty
