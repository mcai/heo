#!/usr/bin/env python3

import os

from common import bench_and_trace_file_name_range, working_directory, max_cycles, num_cores, aco_selection_alpha_range, \
    reinforcement_factor_range
from utils import add_experiment, run_experiments


def run(bench, trace_file_name, num_nodes, routing, selection, max_cycles, aco_selection_alpha, reinforcement_factor):
    dir = working_directory(bench, num_nodes, routing, selection, max_cycles, aco_selection_alpha, reinforcement_factor)

    os.system('rm -fr ' + dir)
    os.system('mkdir -p ' + dir)

    cmd_run = '~/GoProjects/bin/heo -d=' + dir + ' -b=' + bench + ' -f=' + trace_file_name \
              + ' -n=' + str(num_nodes) + ' -r=' + routing + ' -s=' + selection \
              + ' -c=' + str(max_cycles) \
              + ' -sa=' + str(aco_selection_alpha) + ' -rf=' + str(reinforcement_factor)
    print(cmd_run)
    os.system(cmd_run)


def run_experiment(args):
    run(*args)


experiments = []


for bench, trace_file_name in bench_and_trace_file_name_range:
    add_experiment(experiments, bench, trace_file_name + '.combined', num_cores, 'OddEven', 'BufferLevel', max_cycles, -1, -1)

    for aco_selection_alpha in aco_selection_alpha_range:
        for reinforcement_factor in reinforcement_factor_range:
            add_experiment(experiments, bench, trace_file_name + '.combined', num_cores, 'OddEven', 'ACO', max_cycles, aco_selection_alpha, reinforcement_factor)

run_experiments(experiments, run_experiment)
