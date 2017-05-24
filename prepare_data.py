import sys

from common import bench_and_trace_file_name_range, num_nodes

results = []

num_groups = 6
cores_per_group = 10

for bench, trace_file_name in bench_and_trace_file_name_range:
    trace_file_name_combined = trace_file_name + '.combined'
    print("Combining " + str(num_groups) + " " + trace_file_name + " into " + trace_file_name_combined)

    with open(trace_file_name, "r") as f_original:
        with open(trace_file_name_combined, "w") as f_combined:
            for line in f_original:
                parts = str.split(line, ',')

                if parts[0] == '':
                    continue

                thread_id = int(parts[0])

                if thread_id >= num_nodes - 2:
                    print('threadId is out of range, corresponding line.\n')
                    sys.exit(-1)

                for i in range(num_groups):
                    thread_id_combined = cores_per_group * i + thread_id
                    line_combined = str(thread_id_combined) + ',' + parts[1] + ',' + parts[2] + ',' + parts[3]
                    f_combined.write(line_combined)
