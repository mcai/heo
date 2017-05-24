synthesized_traffic_range = [
    'Uniform',
    'Transpose1',
    'Transpose2',
]

bench_and_trace_file_name_range = [
    # ('simple_pthread', 'test_traces/simple_pthread.trace.21454.0'),

    ('blackscholes', 'traces/blackscholes.trace.4.txt'),
    ('bodytrack', 'traces/bodytrack.trace.4.txt'),
    ('canneal', 'traces/canneal.trace.4.txt'),
    ('fluidanimate', 'traces/fluidanimate.trace.4.txt'),
    ('freqmine', 'traces/freqmine.trace.4.txt'),
    ('streamcluster', 'traces/streamcluster.trace.4.txt'),
    ('x264', 'traces/x264.trace.4.txt'),
]

# max_cycles = 100000000
# max_cycles = 10000000
max_cycles = 1000
num_nodes = 64

synthesized_data_packet_injection_rate_range = [
    0.015,
    0.030,
    0.045,
    0.060,
    0.075,
    0.090,
    0.105,
    0.120,
]

trace_driven_data_packet_injection_rate_range = [
    0.05,
    0.1,
    0.15,
    0.2,
    0.21,
    0.22,
    0.23,
    0.24,
    0.25,
]

ant_packet_injection_rate = 0.0002

aco_selection_alpha_range = [
    0.30,
    # 0.35,
    0.40,
    # 0.45,
    0.50,
    # 0.55,
    0.60,
    # 0.65,
    0.70,
]

reinforcement_factor_range = [
    0.0005,
    # 0.001,
    0.002,
    # 0.004,
    0.008,
    # 0.016,
    0.032,
    # 0.064,
]


def working_directory(traffic, bench, max_cycles, num_nodes, routing, selection, data_packet_injection_rate, aco_selection_alpha, reinforcement_factor):
    return 'results/t_' + traffic + '/b_' + bench + '/c_' + str(max_cycles) + '/n_' + str(num_nodes) + '/r_' + routing + '/s_' + selection \
           + '/di_' + str(data_packet_injection_rate) + '/ai_' + str(aco_selection_alpha) + '/rf_' + str(reinforcement_factor)
