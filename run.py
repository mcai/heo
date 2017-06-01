#!/usr/bin/env python3

import os

from common import bench_and_trace_file_name_range, working_directory, max_cycles, num_nodes, aco_selection_alpha_range, \
    reinforcement_factor_range, ant_packet_injection_rate, trace_driven_data_packet_injection_rate_range, \
    synthesized_data_packet_injection_rate_range, synthesized_traffic_range
from utils import add_experiment, run_experiments


def run(traffic, bench, trace_file_name, max_cycles, num_nodes, routing, selection, data_packet_injection_rate, aco_selection_alpha, reinforcement_factor):
    dir = working_directory(traffic, bench, max_cycles, num_nodes, routing, selection, data_packet_injection_rate, aco_selection_alpha, reinforcement_factor)

    stats_file_name = dir + '/stats.json'
    if os.path.isfile(stats_file_name):
        print('Stats file exists, skipped: ' + stats_file_name)
        return

    os.system('rm -fr ' + dir)
    os.system('mkdir -p ' + dir)

    cmd_run = '~/go/bin/heo -d=' + dir + ' -t=' + traffic + ' -tf=' + trace_file_name \
              + ' -c=' + str(max_cycles) + ' -n=' + str(num_nodes) + ' -r=' + routing + ' -s=' + selection \
              + ' -di=' + str(data_packet_injection_rate) + ' -ai=' + str(ant_packet_injection_rate) \
              + ' -sa=' + str(aco_selection_alpha) + ' -rf=' + str(reinforcement_factor)
    print(cmd_run)
    os.system(cmd_run)


def run_experiment(args):
    run(*args)


experiments = []


for traffic in synthesized_traffic_range:
    for data_packet_injection_rate in synthesized_data_packet_injection_rate_range:
        add_experiment(experiments, traffic, '', '', max_cycles, num_nodes, 'XY', 'Random', data_packet_injection_rate, -1, -1)
        add_experiment(experiments, traffic, '', '', max_cycles, num_nodes, 'OddEven', 'BufferLevel', data_packet_injection_rate, -1, -1)

        for aco_selection_alpha in aco_selection_alpha_range:
            for reinforcement_factor in reinforcement_factor_range:
                add_experiment(experiments, traffic, '', '', max_cycles, num_nodes, 'OddEven', 'ACO', data_packet_injection_rate, aco_selection_alpha, reinforcement_factor)

for bench, trace_file_name in bench_and_trace_file_name_range:
    for data_packet_injection_rate in trace_driven_data_packet_injection_rate_range:
        add_experiment(experiments, 'Trace', bench, trace_file_name + '.combined', max_cycles, num_nodes, 'XY', 'Random', data_packet_injection_rate, -1, -1)
        add_experiment(experiments, 'Trace', bench, trace_file_name + '.combined', max_cycles, num_nodes, 'OddEven', 'BufferLevel', data_packet_injection_rate, -1, -1)

        for aco_selection_alpha in aco_selection_alpha_range:
            for reinforcement_factor in reinforcement_factor_range:
                add_experiment(experiments, 'Trace', bench, trace_file_name + '.combined', max_cycles, num_nodes, 'OddEven', 'ACO', data_packet_injection_rate, aco_selection_alpha, reinforcement_factor)

run_experiments(experiments, run_experiment)
