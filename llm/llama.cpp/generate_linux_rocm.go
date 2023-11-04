//go:build rocm

package llm

//go:generate git submodule init

//go:generate git submodule update --force ggml
//go:generate git -C ggml apply ../patches/0001-add-detokenize-endpoint.patch
//go:generate git -C ggml apply ../patches/0002-34B-model-support.patch
//go:generate git -C ggml apply ../patches/0005-ggml-support-CUDA-s-half-type-for-aarch64-1455-2670.patch
//go:generate git -C ggml apply ../patches/0001-copy-cuda-runtime-libraries.patch

//go:generate git submodule update --force gguf
//go:generate git -C gguf apply ../patches/0001-copy-cuda-runtime-libraries.patch
//go:generate git -C gguf apply ../patches/0001-update-default-log-target.patch

//go:generate cmake -S ggml -B ggml/build/rocm -DLLAMA_CLBLAST=on -DLLAMA_ACCELERATE=on -DLLAMA_K_QUANTS=on
//go:generate cmake --build ggml/build/rocm --target server --config Release
//go:generate mv ggml/build/rocm/bin/server ggml/build/rocm/bin/ollama-runner

//go:generate cmake -S gguf -B gguf/build/rocm -DLLAMA_HIPBLAS=on -DLLAMA_ACCELERATE=on -DLLAMA_K_QUANTS=on -DCMAKE_C_COMPILER=$ROCM_PATH/llvm/bin/clang -DCMAKE_CXX_COMPILER=$ROCM_PATH/llvm/bin/clang++ -DAMDGPU_TARGETS='gfx803;gfx900;gfx906:xnack-;gfx908:xnack-;gfx90a:xnack+;gfx90a:xnack-;gfx1010;gfx1012;gfx1030;gfx1100;gfx1101;gfx1102' -DGPU_TARGETS='gfx803;gfx900;gfx906:xnack-;gfx908:xnack-;gfx90a:xnack+;gfx90a:xnack-;gfx1010;gfx1012;gfx1030;gfx1100;gfx1101;gfx1102'
//go:generate cmake --build gguf/build/rocm --target server --config Release
//go:generate mv gguf/build/rocm/bin/server gguf/build/rocm/bin/ollama-runner