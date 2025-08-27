export TEST=$(pwd)/tests
export IDATABASE=$(pwd)/id_16054

test() {
    exist imgui.ini
    if [ $? ]; then
        rm imgui.ini
    fi
    go run .
}
