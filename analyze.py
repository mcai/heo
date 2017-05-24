#!/usr/bin/env python3
from common import bench_and_trace_file_name_range, working_directory, max_cycles, aco_selection_alpha_range, \
    reinforcement_factor_range, data_packet_injection_rate_range, num_cores
from utils import parse_result, to_csv, generate_plot


def results_to_csv(results, csv_file_name):
    def algorithm(r):
        if r.configs['Selection'] == 'ACO':
            return r.configs['Routing'] + '+' + r.configs['Selection'] + '(' + str(r.configs['AcoSelectionAlpha']) + ', ' + str(r.configs['ReinforcementFactor']) + ')'
        else:
            return r.configs['Routing'] + '+' + r.configs['Selection']

    to_csv(csv_file_name, results, [
        ('Bench', lambda r: r.props['bench']),

        ('Routing', lambda r: r.configs['Routing']),
        ('Selection', lambda r: r.configs['Selection']),

        ('Data Packet Injection Rate', lambda r: r.configs['DataPacketInjectionRate']),

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


for data_packet_injection_rate in data_packet_injection_rate_range:
    results = []

    for bench, _ in bench_and_trace_file_name_range:
            results.append(
                parse_result(working_directory(bench, num_cores, 'XY', 'Random', max_cycles, data_packet_injection_rate, -1, -1),
                             bench=bench)
            )
            results.append(
                parse_result(working_directory(bench, num_cores, 'OddEven', 'BufferLevel', max_cycles, data_packet_injection_rate, -1, -1),
                             bench=bench)
            )

            for aco_selection_alpha in aco_selection_alpha_range:
                for reinforcement_factor in reinforcement_factor_range:
                    results.append(
                        parse_result(working_directory(bench, num_cores, 'OddEven', 'ACO', max_cycles, data_packet_injection_rate, aco_selection_alpha, reinforcement_factor),
                                     bench=bench)
                    )

    csv_file_name = 'results/general_' + str(data_packet_injection_rate) + '.csv'

    results_to_csv(results, csv_file_name)

    generate_plot(csv_file_name,
                  'results/simulation_time_' + str(data_packet_injection_rate), 'Bench', 'Simulation Time (Seconds)',
                  'Algorithm', 'Simulation Time (Seconds)')

    generate_plot(csv_file_name,
                  'results/throughput_' + str(data_packet_injection_rate), 'Bench', 'Throughput',
                  'Algorithm', 'Throughput')

    generate_plot(csv_file_name,
                  'results/average_packet_delay_' + str(data_packet_injection_rate), 'Bench', 'Average Packet Delay',
                  'Algorithm', 'Average Packet Delay')

    generate_plot(csv_file_name,
                  'results/payload_throughput_' + str(data_packet_injection_rate), 'Bench', 'Payload Throughput',
                  'Algorithm', 'Payload Throughput')

    generate_plot(csv_file_name,
                  'results/average_payload_packet_delay_' + str(data_packet_injection_rate), 'Bench', 'Average Payload Packet Delay',
                  'Algorithm', 'Average Payload Packet Delay')

for bench, __ in bench_and_trace_file_name_range:
    results = []

    for data_packet_injection_rate in data_packet_injection_rate_range:
        results.append(
            parse_result(working_directory(bench, num_cores, 'XY', 'Random', max_cycles, data_packet_injection_rate, -1, -1),
                         bench=bench)
        )
        results.append(
            parse_result(working_directory(bench, num_cores, 'OddEven', 'BufferLevel', max_cycles, data_packet_injection_rate, -1, -1),
                         bench=bench)
        )

        for aco_selection_alpha in aco_selection_alpha_range:
            for reinforcement_factor in reinforcement_factor_range:
                results.append(
                    parse_result(working_directory(bench, num_cores, 'OddEven', 'ACO', max_cycles, data_packet_injection_rate, aco_selection_alpha, reinforcement_factor),
                                 bench=bench)
                )

    csv_file_name = 'results/general_' + bench + '.csv'

    results_to_csv(results, csv_file_name)

    generate_plot(csv_file_name,
                  'results/simulation_time_' + bench, 'Data Packet Injection Rate', 'Simulation Time (Seconds)',
                  'Algorithm', 'Simulation Time (Seconds)')

    generate_plot(csv_file_name,
                  'results/throughput_' + bench, 'Data Packet Injection Rate', 'Throughput',
                  'Algorithm', 'Throughput')

    generate_plot(csv_file_name,
                  'results/average_packet_delay_' + bench, 'Data Packet Injection Rate', 'Average Packet Delay',
                  'Algorithm', 'Average Packet Delay')

    generate_plot(csv_file_name,
                  'results/payload_throughput_' + bench, 'Data Packet Injection Rate', 'Payload Throughput',
                  'Algorithm', 'Payload Throughput')

    generate_plot(csv_file_name,
                  'results/average_payload_packet_delay_' + bench, 'Data Packet Injection Rate', 'Average Payload Packet Delay',
                  'Algorithm', 'Average Payload Packet Delay')
