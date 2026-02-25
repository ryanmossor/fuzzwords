#!/usr/bin/env bash

if [ -z "$TMUX" ]; then
    echo "tmux not running. Exiting..."
    exit 0
fi

if [[ $(uname -s) == "Darwin" ]]; then
    cache_dir="$HOME/Library/Caches"
else
    cache_dir="$XDG_CACHE_HOME"
fi

current_dir=$(pwd)

tmux new-window -a
tmux rename-window $(basename $current_dir)
main_pane=$(tmux display-message -p '#{pane_id}')

tmux split-window -h -p 40 -c "$current_dir"
log_pane=$(tmux display-message -p '#{pane_id}')
tmux select-pane -t "$log_pane"

log_cmd="tail -f "$cache_dir/fuzzwords/log.json" | jq"
tmux send-keys -t "$log_pane" "$log_cmd" C-m

tmux send-keys -t "$main_pane" "go run ." C-m
tmux select-pane -t "$main_pane"
