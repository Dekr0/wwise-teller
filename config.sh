test() {
    exist imgui.ini
    if [ $? ]; then
        rm imgui.ini
    fi
    go run .
}
