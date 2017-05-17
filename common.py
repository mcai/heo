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

max_cycles = 100000000

num_cores = 64


def working_directory(bench, num_nodes, routing, selection, max_cycles):
    return 'results/' + str(num_nodes) + '/' + routing + '/' + selection + '/' + bench + '/' + str(max_cycles)
