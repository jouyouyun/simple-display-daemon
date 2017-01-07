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
	cp ${TARGET}.desktop ${PREFIX}/share/xsessions/
	mkdir -p ${PREFIX}/share/${TARGET}
	cp data/* ${PREFIX}/share/${TARGET}/
