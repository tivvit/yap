.PHONY: a

all: file a b

file:
	echo "a" > file

a: file # this depends on the file but it does not react to it because it is phony
	@echo "a"

b: file  # this depends on the file but it does not react to it because b does never exist
	@echo "b"