cmake_minimum_required(VERSION 3.1)
project(ed25519)

add_library(ed25519 STATIC
	collective.c
	fe.c
	ge.c
	keypair.c
	sc.c
	sha512.c
	sign.c
	verify.c
)
