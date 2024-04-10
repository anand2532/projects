APP_NAME    := OA
APP_VERSION := 0.1.1


TARGET_OS   := linux
TARGET_ARCH := arm64

OUT_DIR  := bin
EXE_NAME := ${APP_NAME}-${APP_VERSION}-${TARGET_ARCH}
EXE_FILE := ${OUT_DIR}/${EXE_NAME}

TARGET_CLIENT  := cplid@192.168.68.59
TARGET_DIR     := ~

install-deps:
	sudo apt install sshpass

clean:
	rm ${OUT_DIR}/* ||:

build:
	GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH} go build -o ${EXE_FILE}

deploy: clean build
	sshpass -p "78963" scp ${EXE_FILE} ${TARGET_CLIENT}:${TARGET_DIR}

deploy-and-run: deploy
	sshpass -p "78963" ssh -v ${TARGET_CLIENT} ./${EXE_NAME}





