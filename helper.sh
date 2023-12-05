#!/bin/bash

set_color_env()
{
    OK="$(tput setaf 2)[OK]$(tput sgr0)"
    ERROR="$(tput setaf 1)[ERROR]$(tput sgr0)"
    NOTE="$(tput setaf 3)[NOTE]$(tput sgr0)"
    WARN="$(tput setaf 166)[WARN]$(tput sgr0)"
    CAT="$(tput setaf 6)[ACTION]$(tput sgr0)"
    RESET=$(tput sgr0)

    RED=$(tput setaf 1)
    GREN=$(tput setaf 2)
    YELLOW=$(tput setaf 3)
    BLUE=$(tput setaf 4)
    LIGHTBLUE=$(tput setaf 6)
    ORANGE=$(tput setaf 166)
    WIGHT=$(tput setaf 255)

    BOLD=$(tput bold)
}

set_basic_env(){
    set_color_env
    set -e
    TERM_WIDTH=$(tput cols)
}

print_help(){
    echo "帮助(todo)"
}

new_content(){
    CONTENTPATH=$1
    mkdir -p $CONTENTPATH
    mkdir -p $CONTENTPATH/media
    mkdir -p $CONTENTPATH/content
    TITLE=$(echo "$CONTENTPATH" | awk -F'/' '{print $NF}')
    echo "# $TITLE" > $CONTENTPATH/README.md
    echo "" >> $CONTENTPATH/README.md
    echo "简介:" >> $CONTENTPATH/README.md
    echo "" >> $CONTENTPATH/README.md
    echo "# 目录" >> $CONTENTPATH/README.md
    echo "" >> $CONTENTPATH/README.md
    echo "## 二级标题" >> $CONTENTPATH/README.md
    echo "" >> $CONTENTPATH/README.md
    echo "二级简介" >> $CONTENTPATH/README.md
    echo "" >> $CONTENTPATH/README.md
    echo "- [xx](xx)" >> $CONTENTPATH/README.md
}

set_basic_env

ACTION=$1
case "$ACTION" in
"new")
    new_content $2
;;
*)
    print_help
;;
esac
