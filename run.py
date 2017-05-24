#!/usr/bin/env python3

import os

from common import bench_and_trace_file_name_range, working_directory, max_cycles, num_cores, aco_selection_alpha_range, \
    reinforcement_factor_range, ant_packet_injection_rate, data_packet_injection_rate_range
from utils import add_experiment, run_experiments


def run(bench, trace_file_name, num_nodes, routing, selection, max_cycles, data_packet_injection_rate, aco_selection_alpha, reinforcement_factor):
    dir = working_directory(bench, num_nodes, routing, selection, max_cycles, data_packet_injection_rate, aco_selection_alpha, reinforcement_factor)

    os.system('rm -fr ' + dir)
    os.system('mkdir -p ' + dir)

    cmd_run = '~/GoProjects/bin/heo -d=' + dir + ' -b=' + bench + ' -f=' + trace_file_name \
              + ' -n=' + str(num_nodes) + ' -r=' + routing + ' -s=' + selection \
              + ' -c=' + str(max_cycles) \
              + ' -di=' + str(data_packet_injection_rate) + ' -ai=' + str(ant_packet_injection_rate) \
              + ' -sa=' + str(aco_selection_alpha) + ' -rf=' + str(reinforcement_factor)
    print(cmd_run)
    os.system(cmd_run)


def run_experiment(args):
    run(*args)


experiments = []


for bench, trace_file_name in bench_and_trace_file_name_range:
    for data_packet_injection_rate in data_packet_injection_rate_range:
        add_experiment(experiments, bench, trace_file_name + '.combined', num_cores, 'XY', 'Random', max_cycles, data_packet_injection_rate, -1, -1)
        add_experiment(experiments, bench, trace_file_name + '.combined', num_cores, 'OddEven', 'BufferLevel', max_cycles, data_packet_injection_rate, -1, -1)

        for aco_selection_alpha in aco_selection_alpha_range:
            for reinforcement_factor in reinforcement_factor_range:
                add_experiment(experiments, bench, trace_file_name + '.combined', num_cores, 'OddEven', 'ACO', max_cycles, data_packet_injection_rate, aco_selection_alpha, reinforcement_factor)

run_experiments(experiments, run_experiment)
