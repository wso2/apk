#!/bin/bash


set -e
script_dir=$(dirname "$0")
# Change directory to make sure logs directory is created inside $script_dir
cd "$script_dir"

test_name="cpu-1"
heap_size="1g"
duration=1200
remote_hosts=""
tokens_path=""

declare -a user_counts_array=(1000 500 200 100 50 10)
declare -a payloads_array=("102400B" "10240B" "1024B" "50B")



function usage() {
    echo ""
    echo "Usage: "
    echo "$0 [-n test_name>] [-r <remote_hosts>] "
    echo ""
    echo "-n: Test Name."
    echo "-r: Remote Servers."
    echo "-t: Tokens csv file path."
    echo "-h: Display this help and exit."
    echo ""
}

while getopts "n:r:d:t:h" opts; do
    case $opts in
    n)
        test_name=${OPTARG}
        ;;
    r)
        remote_hosts=${OPTARG}
        ;;
    d)
        duration=${OPTARG}
        ;;
    t)
        tokens_path=${OPTARG}
        ;;
    h)
        usage
        exit 0
        ;;
    \?)
        usage
        exit 1
        ;;
    esac
done

if [[ -z $remote_hosts ]]; then
    echo "Please specify remote hosts."
    exit 1
fi


# Run Test
for user_count in "${user_counts_array[@]}"; do
    for payload_size in "${payloads_array[@]}"; do
        echo "############################### Start Test ###############################"
        echo "Users Count: $user_count"
        echo "Payload Size: $payload_size"

        results_dir="${HOME}/results/${test_name}/passthrough/${heap_size}_heap/${user_count}_users/${payload_size}/0ms_sleep"
        echo "Results Dir: ${results_dir}"
        mkdir -p "${results_dir}"

        kubectl top po --containers -A > "${results_dir}/resources-start.txt"
        kubectl get po -owide -A > "${results_dir}/pods-distribution.txt"
        nohup sh -c "sleep $((duration/4)) && kubectl top po --containers -A > ${results_dir}/resources-$((duration/4/60))min.txt" >/dev/null &
        nohup sh -c "sleep $((duration/2)) && kubectl top po --containers -A > ${results_dir}/resources-$((duration/2/60))min.txt" >/dev/null &
        nohup sh -c "sleep $((duration*3/4)) && kubectl top po --containers -A > ${results_dir}/resources-$((duration*3/4/60))min.txt" >/dev/null &
        
        pod_name=$(kubectl get pods -n apk-perf-test -l app.kubernetes.io/app=router -o jsonpath='{.items[0].metadata.name}')
        nohup sh -c "kubectl -n apk-perf-test logs $pod_name -c enforcer --since 1s -f > ${results_dir}/enforcer.log"  >/dev/null &
        enforcer_log_pid=$!
        nohup sh -c "echo "Killing enforcer log process" && sleep $duration && kill -9 $enforcer_log_pid" &
        nohup sh -c "kubectl -n apk-perf-test logs $pod_name -c router --since 1s -f > ${results_dir}/router.log"  >/dev/null &
        router_log_pid=$!
        nohup sh -c "echo "Killing enforcer log process" && sleep $duration && kill -9 $router_log_pid" &
        echo echo "Process ids: $enforcer_log_pid, $router_log_pid"
        kubectl describe  pods   $pod_name -n apk-perf-test > ${results_dir}/pod-describe-start.txt
        nohup sh -c "sleep $((duration/4)) && kubectl describe  pods   $pod_name -n apk-perf-test > ${results_dir}/pod-describe-$((duration/4/60))min.txt" >/dev/null &
        nohup sh -c "sleep $((duration/2)) && kubectl describe  pods   $pod_name -n apk-perf-test > ${results_dir}/pod-describe-$((duration/2/60))min.txt" >/dev/null &
        nohup sh -c "sleep $((duration*3/4)) && kubectl describe  pods   $pod_name -n apk-perf-test > ${results_dir}/pod-describe-$((duration*3/4/60))min.txt" >/dev/null &
        nohup sh -c "sleep $duration && kubectl describe  pods   $pod_name -n apk-perf-test > ${results_dir}/pod-describe-end.txt" >/dev/null &
        
        ./run-jmeter.sh -m "$heap_size" -u "$user_count" -p "$payload_size" -d "$duration" -s "$remote_hosts" -r "$results_dir" -t "$tokens_path"

        kubectl top po --containers -A >"${results_dir}/resources-end.txt"

        echo "################################ End Test ################################"
    done
done




