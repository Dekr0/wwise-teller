export TEST=$(pwd)/tests
export IDATABASE=$(pwd)/id_15314

test() {
    exist imgui.ini
    if [ $? ]; then
        rm imgui.ini
    fi
    go run .
}
