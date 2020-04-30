all: build
	
install: 
	git clone https://github.com/jedisct1/libsodium --branch stable
	cd libsodium/
	./configure
	make && make check
	sudo make install
	cd ..

build:
	go build github-secret
