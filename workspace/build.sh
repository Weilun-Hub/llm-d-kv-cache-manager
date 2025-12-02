#apt install -y build-essential clang
#wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
#tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

cd ..

export CGO_CFLAGS="$(python3.12-config --cflags) -I$(pwd)/lib"
export CGO_LDFLAGS="$(python3.12-config --ldflags --embed) -L$(pwd)/lib -ltokenizers -ldl -lm"
export PYTHON=python3.12
export PYTHONPATH=$(pwd)/pkg/preprocessing/chat_completions
export CGO_ENABLED=1

go mod tidy
go mod download

mkdir -p bin
go build -o \
    bin/kv-cache-manager \
    examples/kv_events/online/main.go

