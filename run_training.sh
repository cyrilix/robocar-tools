#! /bin/bash

set +e
set +x

RECORDS_PATH=~/robocar/record-sim4-2
#TRAINING_OPTS="--horizon=20"
TRAINING_OPTS=""
#MODEL_TYPE="categorical"
MODEL_TYPE="linear"
IMG_WIDTH=160
IMG_HEIGHT=120
HORIZON=20

TRAINING_DATA_DIR=/tmp/data
TRAINING_OUTPUT_DIR=/tmp/output
TRAIN_ARCHIVE=${TRAINING_DATA_DIR}/train.zip

#######################

rm -rf ${TRAINING_DATA_DIR} ${TRAINING_OUTPUT_DIR}
mkdir -p ${TRAINING_DATA_DIR}
mkdir -p ${TRAINING_OUTPUT_DIR}

printf "Build archive %s\n\n" "${TRAIN_ARCHIVE}"
go run ./cmd/rc-tools training archive \
            -record-path ${RECORDS_PATH} \
            -output ${TRAIN_ARCHIVE} \
            -image-height ${IMG_HEIGHT} \
            -image-width ${IMG_WIDTH}

printf "\n\nRun training\n\n"
podman run --rm -it \
            -v /tmp/data:/opt/ml/input/data/train \
            -v /tmp/output:/opt/ml/model/ \
            localhost/tensorflow_without_gpu \
                python /opt/ml/code/train.py \
                      --img_height=${IMG_HEIGHT} \
                      --img_width=${IMG_WIDTH} \
                      --batch_size=32 \
                      --model_type=${MODEL_TYPE} \
                      --horizon=${HORIZON} \
                      ${TRAINING_OPTS}

printf "\n\nConvert model\n\n"
edgetpu_compiler -o ${TRAINING_OUTPUT_DIR} ${TRAINING_OUTPUT_DIR}/model_*.tflite

