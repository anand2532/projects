APP_NAME    := tamper-app
APP_VERSION := 0.1.5

TARGET_OS   := linux
TARGET_ARCH := arm64

OUT_DIR  := bin
EXE_NAME := ${APP_NAME}-${APP_VERSION}-${TARGET_ARCH}
EXE_FILE := ${OUT_DIR}/${EXE_NAME}

TARGET_CLIENT  := two@two.local
TARGET_DIR     := ~

install-deps:
	sudo apt install sshpass

clean:
	rm ${OUT_DIR}/* ||:

build:
	GOOS=${TARGET_OS} GOARCH=${TARGET_ARCH} go build -o ${EXE_FILE}

deploy: clean build
	sshpass -p "two" scp ${EXE_FILE} ${TARGET_CLIENT}:${TARGET_DIR}

deploy-and-run: deploy
	sshpass -p "two" ssh ${TARGET_CLIENT} ./${EXE_NAME}

