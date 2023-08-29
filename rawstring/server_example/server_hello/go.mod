module test

go 1.13

replace trpc.group/trpc-go/trpc-codec/rawstring => ../../

require (
	trpc.group/trpc-go/trpc-codec/cmd v0.0.0-20200909043214-7ed5bf81693c
	trpc.group/trpc-go/trpc-codec/rawstring v0.0.0-00010101000000-000000000000
	trpc.group/trpc-go/trpc-go v0.0.0-20230824091938-4699a10e2f35
)
