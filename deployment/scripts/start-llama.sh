#!/bin/bash

# Directory where llama.cpp will be built
LLAMA_DIR="$HOME/llama.cpp"
MODEL_PATH="/Users/manishsv/Library/Caches/llama.cpp/Qwen_Qwen2-7B-Instruct-GGUF_qwen2-7b-instruct-q5_k_m.gguf"

# Check if llama.cpp directory exists
if [ ! -d "$LLAMA_DIR" ]; then
    echo "Cloning llama.cpp..."
    git clone https://github.com/ggerganov/llama.cpp.git "$LLAMA_DIR"
fi

# Build llama.cpp with Metal support
cd "$LLAMA_DIR"
git pull
make clean
LLAMA_METAL=1 make -j

# Start the server with optimized parameters for M1
./server \
    -m "$MODEL_PATH" \
    -c 2048 \
    -ngl 35 \
    -t 8 \
    -cb 512 \
    --host 0.0.0.0 \
    --port 8082 