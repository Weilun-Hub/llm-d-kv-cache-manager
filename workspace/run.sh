
cd ..

export PYTHONPATH=${PYTHONPATH}:$(pwd)/pkg/preprocessing/chat_completions

export ZMQ_ENDPOINT="tcp://*:5557"
export ZMQ_TOPIC="kv@public-wlf3-ge103-kce-node142@kat-coder"
export PYTHONHASHSEED=0
export POOL_CONCURRENCY=4
export BLOCK_SIZE=64

./bin/kv-cache-manager

