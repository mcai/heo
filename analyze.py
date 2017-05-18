#!/usr/bin/env python3
from common import bench_and_trace_file_name_range, working_directory, max_cycles, aco_selection_alpha_range, \
    reinforcement_factor_range, data_packet_injection_rate_range
from utils import parse_result, to_csv, generate_plot

results = []

for bench, trace_file_name in bench_and_trace_file_name_range:
    for data_packet_injection_rate in data_packet_injection_rate_range:
        results.append(
            parse_result(working_directory(bench, 64, 'OddEven', 'BufferLevel', max_cycles, data_packet_injection_rate, -1, -1),
                         bench=bench)
        )

        for aco_selection_alpha in aco_selection_alpha_range:
            for reinforcement_factor in reinforcement_factor_range:
                results.append(
                    parse_result(working_directory(bench, 64, 'OddEven', 'ACO', max_cycles, data_packet_injection_rate, aco_selection_alpha, reinforcement_factor),
                                 bench=bench)
                )


def algorithm(r):
    if r.configs['Selection'] == 'ACO':
        return r.configs['Routing'] + '+' + r.configs['Selection'] + '(' + str(r.configs['AcoSelectionAlpha']) + ', ' + str(r.configs['ReinforcementFactor']) + ')'
    else:
        return r.configs['Routing'] + '+' + r.configs['Selection']

to_csv('results/general.csv', results, [
    ('Bench', lambda r: r.props['bench']),

    ('Routing', lambda r: r.configs['Routing']),
    ('Selection', lambda r: r.configs['Selection']),

    ('Data Packet Rate', lambda r: r.configs['DataPacketRate']),

    ('ACO Selection Alpha', lambda r: r.configs['AcoSelectionAlpha']),
    ('Reinforcement Factor', lambda r: r.configs['ReinforcementFactor']),

    ('Algorithm', algorithm),

    ('Max Cycles', lambda r: r.configs['MaxCycles']),
    ('Simulation Time (Seconds)', lambda r: r.stats['SimulationTimeInSeconds']),

    ('Throughput', lambda r: r.stats['Throughput']),
    ('Average Packet Delay', lambda r: r.stats['AveragePacketDelay']),

    ('Payload Throughput', lambda r: r.stats['PayloadThroughput']),
    ('Average Payload Packet Delay', lambda r: r.stats['AveragePayloadPacketDelay']),
])

# generate_plot('results/general.csv',
#               'results/simulation_time', 'Bench', 'Simulation Time (Seconds)',
#               'Algorithm', 'Simulation Time (Seconds)')
#
# generate_plot('results/general.csv',
#               'results/throughput', 'Bench', 'Throughput',
#               'Algorithm', 'Throughput')
#
# generate_plot('results/general.csv',
#               'results/average_packet_delay', 'Bench', 'Average Packet Delay',
#               'Algorithm', 'Average Packet Delay')
#
# generate_plot('results/general.csv',
#               'results/payload_throughput', 'Bench', 'Payload Throughput',
#               'Algorithm', 'Payload Throughput')
#
# generate_plot('results/general.csv',
#               'results/average_payload_packet_delay', 'Bench', 'Average Payload Packet Delay',
#               'Algorithm', 'Average Payload Packet Delay')
