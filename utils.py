#!/usr/bin/python
# -*- coding: UTF8 -*-
import csv
import json

import collections
from objectpath import *
from os import path
from pyparsing import Word, Optional, ParseException, printables, nums, restOfLine
import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns
import multiprocessing as mp


class ExperimentResults:
    def __init__(self, configs, stats, props):
        self.configs = configs
        self.stats = stats
        self.props = props


class ExperimentConfigs:
    def __init__(self, raw_configs):
        self.raw_configs = raw_configs

    def __getitem__(self, index):
        return self.raw_configs.execute('$.' + index)


class ExperimentStats:
    def __init__(self, raw_stats):
        self.raw_stats = raw_stats

    def __getitem__(self, index):
        return self.raw_stats.execute('$.' + index)


def read_configs(result_dir, config_json_file_name):
    try:
        with open(path.join(result_dir, config_json_file_name)) as config_json_file:
            configs = Tree(json.load(config_json_file))
    except Exception as e:
        print(e)
        return None
    else:
        return configs


def read_stats(result_dir, stats_file_name):
    try:
        with open(path.join(result_dir, stats_file_name)) as stats_json_file:
            configs = Tree(json.load(stats_json_file))
    except Exception as e:
        print(e)
        return None
    else:
        return configs


def parse_result(result_dir, config_json_file_name='config.noc.json', stats_json_file_name='stats.json', **props):
    return ExperimentResults(ExperimentConfigs(read_configs(result_dir, config_json_file_name)),
                             ExperimentStats(read_stats(result_dir, stats_json_file_name)), props)


def to_csv(output_file_name, results, fields):
    with open(output_file_name, 'w') as output_file:
        writer = csv.writer(output_file, delimiter=',', quotechar='"', quoting=csv.QUOTE_ALL)

        writer.writerow([field[0] for field in fields])

        for result in results:
            writer.writerow([field[1](result) for field in fields])


def generate_plot(csv_file_name, plot_file_name, x, y, hue, y_title, xticklabels_rotation=90):
    sns.set(font_scale=1.5)

    sns.set_style("white", {"legend.frameon": True})

    df = pd.read_csv(csv_file_name)

    ax = sns.barplot(data=df, x=x, y=y, hue=hue, palette=sns.color_palette("Paired"))
    ax.set_xlabel('')
    ax.set_ylabel(y_title)

    labels = ax.get_xticklabels()
    ax.set_xticklabels(labels, rotation=xticklabels_rotation)

    fig = ax.get_figure()

    if hue:
        legend = ax.legend(bbox_to_anchor=(1.05, 1), loc='upper left', borderaxespad=0.)
        legend.set_label('')

        fig.savefig(plot_file_name, bbox_extra_artists=(legend,), bbox_inches='tight')
        fig.savefig(plot_file_name + '.jpg', bbox_extra_artists=(legend,), bbox_inches='tight')
    else:
        fig.tight_layout()

        fig.savefig(plot_file_name)
        fig.savefig(plot_file_name + '.jpg')

    plt.clf()
    plt.close('all')


def run_experiments(experiments, run_experiment):
    num_processes = mp.cpu_count()
    pool = mp.Pool(num_processes)
    pool.map(run_experiment, experiments)

    pool.close()
    pool.join()


def add_experiment(experiments, *args):
    experiments.append(args)