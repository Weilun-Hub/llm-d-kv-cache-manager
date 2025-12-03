
cd ..

export PYTHONPATH=${PYTHONPATH}:$(pwd)/pkg/preprocessing/chat_completions

export HF_TOKEN="/home/relay/liujiacheng06/models/kat-coder-pro-v1-0-1/"
export ZMQ_ENDPOINT="tcp://*:5557"
export ZMQ_TOPIC="sglang@public-wlf3-ge103-kce-node142@kat-coder-pro-v1-0-1"
export PYTHONHASHSEED=""
export POOL_CONCURRENCY=4
export BLOCK_SIZE=64

./bin/kv-cache-manager

