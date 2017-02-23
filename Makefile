PREFIX = /usr
TARGET=multi-display-session

all: build

build: ${TARGET}

${TARGET}:
	go build -o ${TARGET}

clean:
	rm -rf ${TARGET}

install:
	install -Dm755 ${TARGET} ${DESTDIR}/${PREFIX}/bin/${TARGET}
	mkdir -p ${DESTDIR}/${PREFIX}/share/xsessions/
	install -Dm644 ${TARGET}.desktop ${DESTDIR}/${PREFIX}/share/xsessions/${TARGET}.desktop
	mkdir -p ${DESTDIR}/${PREFIX}/share/${TARGET}
	cp data/* ${DESTDIR}/${PREFIX}/share/${TARGET}/
