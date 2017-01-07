PREFIX = /usr
TARGET=multi-display-session

all: build

build: ${TARGET}

${TARGET}:
	go build -o ${TARGET}

clean:
	rm ${TARGET}

install:
	mkdir -p ${PREFIX}/bin/
	cp -f ${TARGET} ${PREFIX}/bin/
	mkdir -p ${PREFIX}/share/xsessions/
	cp multi-display-session.desktop ${PREFIX}/share/xsessions/
