# import matplotlib.pyplot as plt
# import plotcommon
# import pandas as pd

# df = pd.read_csv('summary.csv')
# print (df)
# x = df[['Concurrent Users']]
# y = df[['Throughput (Requests/sec)']]
# print (x)
# plt.plot(x, y)
# plt.show()
#plotcommon.save_line_plot('sample.png', 'Throughput (Requests/sec)', 'Throughput vs Concurrent users', 'Payload size', df)


import pandas as pd
import argparse
# import plotcommon
import matplotlib
# matplotlib.use('TkAgg')
import matplotlib.pyplot as plt
import numpy as np
import os

parser = argparse.ArgumentParser(description='Create plots')
parser.add_argument('-f', '--file', required=False, help='The summary CSV file name.', type=str, default='summary.csv')
parser.add_argument('-n', '--name', required=False, help='Test name.', type=str, default='cpu-1')
args = parser.parse_args()

path = os.path.join('plots', args.name)
if not os.path.exists(path):
    os.makedirs(path)

# PLOT_COLUMN_RANGE_START = plotcommon.PLOT_COLUMN_RANGE_START
# PLOT_COLUMN_RANGE_END = plotcommon.PLOT_COLUMN_RANGE_END

print("Reading " + args.file + " file...")
df = pd.read_csv(args.file)

def tps():
    #Throughput
    df1 = df[['Concurrent Users', 'Message Size (Bytes)', 'Throughput (Requests/sec)']]

    x_values = list(df1.groupby(['Concurrent Users']).groups.keys())
    print(x_values)
    y_50_values = [df1.iloc[j]['Throughput (Requests/sec)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['50B']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    y_1024_values = [df1.iloc[j]['Throughput (Requests/sec)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['1KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    y_10240_values = [df1.iloc[j]['Throughput (Requests/sec)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['10KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    y_102400_values = [df1.iloc[j]['Throughput (Requests/sec)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['100KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    print (list(x_values))
    print (y_50_values)

    plt.rcParams.update({'font.size': 20})
    plt.figure(figsize=(3304/144, 1726/144), dpi=144)
    plt.plot(x_values, y_50_values, label = "50B", marker='o')
    plt.plot(x_values, y_1024_values, label = "1KiB", marker='o')
    plt.plot(x_values, y_10240_values, label = "10KiB", marker='o')
    plt.plot(x_values, y_102400_values, label = "100KiB", marker='o')
    plt.xlabel('Concurrent Users')
    plt.ylabel('Throughput (Requests/sec)')
    
    plt.legend(title="Message Size")
    plt.title('Throughput vs Concurrent Users for 0ms backend delay')
    plt.savefig(os.path.join(path, 'tps_0ms.png'))
    plt.clf()


def response_time():
    df1 = df[['Concurrent Users', 'Message Size (Bytes)', 'Average Response Time (ms)']]

    x_values = list(df1.groupby(['Concurrent Users']).groups.keys())
    y_50_values = [df1.iloc[j]['Average Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['50B']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    y_1024_values = [df1.iloc[j]['Average Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['1KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    y_10240_values = [df1.iloc[j]['Average Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['10KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    y_102400_values = [df1.iloc[j]['Average Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['100KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    print (list(x_values))
    print (y_50_values)

    plt.rcParams.update({'font.size': 20})
    plt.figure(figsize=(3304/144, 1726/144), dpi=144)
    plt.plot(x_values, y_50_values, label = "50B", marker='o')
    plt.plot(x_values, y_1024_values, label = "1KiB", marker='o')
    plt.plot(x_values, y_10240_values, label = "10KiB", marker='o')
    plt.plot(x_values, y_102400_values, label = "100KiB", marker='o')
    plt.xlabel('Concurrent Users')
    plt.ylabel('Average Response Time (ms)')
    
    plt.legend(title="Message Size")
    plt.title('Average Response Time (ms) vs Concurrent Users for 0ms backend delay')
    plt.savefig(os.path.join(path, 'response_time_0ms.png'))
    plt.clf()


    # ##GC throughput enforcer
    # df1 = df[['Concurrent Users', 'Message Size (Bytes)', 'WSO2 API Microgateway GC Throughput (%)']]

    # x_values = list(df1.groupby(['Concurrent Users']).groups.keys())
    # y_50_values = [df1.iloc[j]['WSO2 API Microgateway GC Throughput (%)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups[50]).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_1024_values = [df1.iloc[j]['WSO2 API Microgateway GC Throughput (%)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups[1024]).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_10240_values = [df1.iloc[j]['WSO2 API Microgateway GC Throughput (%)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups[10240]).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_102400_values = [df1.iloc[j]['WSO2 API Microgateway GC Throughput (%)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['100KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # print (list(x_values))
    # print (y_50_values)

    # plt.rcParams.update({'font.size': 20})
    # plt.figure(figsize=(3304/144, 1726/144), dpi=144)
    # plt.plot(x_values, y_50_values, label = "50B", marker='o')
    # plt.plot(x_values, y_1024_values, label = "1KiB", marker='o')
    # plt.plot(x_values, y_10240_values, label = "10KiB", marker='o')
    # plt.xlabel('Concurrent Users')
    # plt.ylabel('WSO2 Choreo-Connect enforcer GC Throughput (%)')
    
    # plt.legend(title="Message Size")
    # plt.title('GC Throughput % (enforcer) vs Concurrent Users for 0ms backend delay')
    # plt.savefig('gc_tps_0ms.png')
    # plt.clf()


    #percentiles 90
    # df1 = df[['Concurrent Users', 'Message Size (Bytes)', '90th Percentile of Response Time (ms)']]

    # x_values = np.array(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # y_50_values = [df1.iloc[j]['90th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['50B']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_1024_values = [df1.iloc[j]['90th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['1KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_10240_values = [df1.iloc[j]['90th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['10KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_102400_values = [df1.iloc[j]['90th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['100KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # print (list(x_values))
    # print (y_50_values)

    # plt.rcParams.update({'font.size': 20})
    # plt.figure(figsize=(3304/144, 1726/144), dpi=144)
    # # ax = fig.add_axes([0,0,1,1])
    # fig, ax = plt.subplots(figsize=(3304/144, 1726/144))
    # ax.bar(x_values-10, y_50_values, label = "50B", width=10)
    # ax.bar(x_values, y_1024_values, label = "1KiB", width=10)
    # ax.bar(x_values+10, y_10240_values, label = "10KiB", width=10)    
    # ax.bar(x_values+20, y_102400_values, label = "100KiB", width=10)
    # ax.set_xticks(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # ax.set_xticklabels(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # plt.xlabel('Concurrent Users')
    # plt.ylabel('90th Percentile of Response Time (ms)')

    
    # plt.legend(title="Message Size")
    # plt.title('90th percentile of Response time vs Concurrent Users for 0ms backend delay')
    # plt.savefig('p90_0ms.png')
    # #plt.show()
    # plt.clf()


    #percentiles 95
    # df1 = df[['Concurrent Users', 'Message Size (Bytes)', '95th Percentile of Response Time (ms)']]

    # x_values = np.array(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # y_50_values = [df1.iloc[j]['95th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['50B']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_1024_values = [df1.iloc[j]['95th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['1KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_10240_values = [df1.iloc[j]['95th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['10KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_102400_values = [df1.iloc[j]['95th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['100KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # print (list(x_values))
    # print (y_50_values)

    # plt.rcParams.update({'font.size': 20})
    # plt.figure(figsize=(3304/144, 1726/144), dpi=144)
    # # ax = fig.add_axes([0,0,1,1])
    # fig, ax = plt.subplots(figsize=(3304/144, 1726/144))
    # ax.bar(x_values-10, y_50_values, label = "50B", width=10)
    # ax.bar(x_values, y_1024_values, label = "1KiB", width=10)
    # ax.bar(x_values+10, y_10240_values, label = "10KiB", width=10)
    # ax.bar(x_values+20, y_10240_values, label = "100KiB", width=10)
    # ax.set_xticks(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # ax.set_xticklabels(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # plt.xlabel('Concurrent Users')
    # plt.ylabel('95th Percentile of Response Time (ms)')

    
    # plt.legend(title="Message Size")
    # plt.title('95th percentile of Response time vs Concurrent Users for 0ms backend delay')
    # plt.savefig('p95_0ms.png')
    # #plt.show()
    # plt.clf()


    # #percentiles 99
    # df1 = df[['Concurrent Users', 'Message Size (Bytes)', '99th Percentile of Response Time (ms)']]

    # x_values = np.array(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # y_50_values = [df1.iloc[j]['99th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['50B']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_1024_values = [df1.iloc[j]['99th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['1KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_10240_values = [df1.iloc[j]['99th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['10KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_102400_values = [df1.iloc[j]['99th Percentile of Response Time (ms)'] for j in [set(df1.groupby(['Message Size (Bytes)']).groups['100KiB']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # print (list(x_values))
    # print (y_50_values)

    # plt.rcParams.update({'font.size': 20})
    # plt.figure(figsize=(3304/144, 1726/144), dpi=144)
    # # ax = fig.add_axes([0,0,1,1])
    # fig, ax = plt.subplots(figsize=(3304/144, 1726/144))
    # ax.bar(x_values-10, y_50_values, label = "50B", width=10)
    # ax.bar(x_values, y_1024_values, label = "1KiB", width=10)
    # ax.bar(x_values+10, y_10240_values, label = "10KiB", width=10)
    # ax.bar(x_values+20, y_102400_values, label = "100KiB", width=10)
    # ax.set_xticks(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # ax.set_xticklabels(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # plt.xlabel('Concurrent Users')
    # plt.ylabel('99th Percentile of Response Time (ms)')

    
    # plt.legend(title="Message Size")
    # plt.title('99th percentile of Response time vs Concurrent Users for 0ms backend delay')
    # plt.savefig('p99_0ms.png')
    # #plt.show()
    # plt.clf()

    # x_values = np.array(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # print(x_values)
    # y_50_values = [df1.iloc[j]['Message Size (Bytes)'] for j in [set(df1.groupby(['90th Percentile of Response Time (ms)']).groups['90th Percentile of Response Time (ms)']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_1024_values = [df1.iloc[j]['Message Size (Bytes)'] for j in [set(df1.groupby(['95th Percentile of Response Time (ms)']).groups['50B']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # y_10240_values = [df1.iloc[j]['Message Size (Bytes)'] for j in [set(df1.groupby(['99th Percentile of Response Time (ms)']).groups['50B']).intersection(set(df1.groupby(['Concurrent Users']).groups[i])).pop()  for i in x_values]]
    # print (list(x_values))
    # print (y_50_values)

    # plt.rcParams.update({'font.size': 20})
    # plt.figure(figsize=(3304/144, 1726/144), dpi=144)
    # # ax = fig.add_axes([0,0,1,1])
    # fig, ax = plt.subplots(figsize=(3304/144, 1726/144))
    # ax.bar(x_values-10, y_50_values, label = "50B", width=10)
    # ax.bar(x_values, y_1024_values, label = "1KiB", width=10)
    # ax.bar(x_values+10, y_10240_values, label = "10KiB", width=10)
    # ax.set_xticks(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # ax.set_xticklabels(list(df1.groupby(['Concurrent Users']).groups.keys()))
    # plt.xlabel('Concurrent Users')
    # plt.ylabel('99th Percentile of Response Time (ms)')

    
    # plt.legend(title="Message Size")
    # plt.title('99th percentile of Response time vs Concurrent Users for 0ms backend delay')
    # plt.savefig('test.png')
    # #plt.show()
    # plt.clf()

    
    # Save lmplots first
    # unique_heap_sizes = df['Heap Size'].unique()
    # for heap_size in unique_heap_sizes:
    #     df_heap = df.loc[df['Heap Size'] == heap_size]
    #     xcolumns = ["Concurrent Users", "Message Size (Bytes)", "Back-end Service Delay (ms)"]
    #     for xcolumn in xcolumns:
    #         for ycolumn in df.columns[PLOT_COLUMN_RANGE_START:PLOT_COLUMN_RANGE_END]:
    #             file_suffix = "-" + str(heap_size) + ".png"
    #             plotcommon.save_lm_plot(
    #                 "lmplot-" + plotcommon.get_filename(ycolumn) + "-" + plotcommon.get_filename(xcolumn) + file_suffix,
    #                 xcolumn, ycolumn, ycolumn + " vs " + xcolumn, df_heap)

    # df.rename(columns={'Message Size (Bytes)': 'Message Size',
    #                    'Back-end Service Delay (ms)': 'Back-end Service Delay'},
    #           inplace=True)
    # # Format message size values
    # df['Message Size'] = df['Message Size'].map(plotcommon.format_bytes)
    # # Format time
    # df['Back-end Service Delay'] = df['Back-end Service Delay'].map(plotcommon.format_time)

    # unique_backend_delays = df['Back-end Service Delay'].unique()
    # unique_message_sizes = df['Message Size'].unique()

    # for heap_size in unique_heap_sizes:
    #     df_heap = df.loc[df['Heap Size'] == heap_size]

    #     for backend_delay in unique_backend_delays:
    #         # Plot individual charts for message sizes
    #         for message_size in unique_message_sizes:
    #             df_data = df_heap.loc[
    #                 (df_heap['Message Size'] == message_size) & (df_heap['Back-end Service Delay'] == backend_delay)]
    #             print(type(message_size))
    #             file_suffix = "-" + str(heap_size) + "-" + message_size + "-" + backend_delay + ".png"
    #             subtitle = "Memory = " + str(heap_size) + ", Message Size = " + message_size + ", Back-end Service Delay = " + backend_delay
    #             for ycolumn in df.columns[PLOT_COLUMN_RANGE_START:PLOT_COLUMN_RANGE_END]:
    #                 plotcommon.save_line_plot("lineplot-" + plotcommon.get_filename(ycolumn) + file_suffix, ycolumn,
    #                                           ycolumn + " vs Concurrent Users", subtitle, df_data)

    #         # Categorical plots by message size
    #         for ycolumn in df.columns[PLOT_COLUMN_RANGE_START:PLOT_COLUMN_RANGE_END]:
    #             # Cat plot for each backend delay, with message size as column
    #             df_data = df_heap.loc[df_heap['Back-end Service Delay'] == backend_delay]
    #             file_suffix = "-" + str(heap_size) + "-" + backend_delay + ".png"
    #             subtitle = "Memory = " + str(heap_size) + ", Back-end Service Delay = " + backend_delay
    #             plotcommon.save_cat_plot("catplot-" + plotcommon.get_filename(ycolumn) + file_suffix, ycolumn,
    #                                      ycolumn + " vs Concurrent Users", subtitle, df_data, "Message Size")

    # print("Done")



def percentiles():
    plt.rcParams.update(matplotlib.rcParamsDefault)
    plt.clf()
    
    figure, axes = plt.subplots(2, 2)
    figure.tight_layout()
    axes[0, 0].set_title('Message size = 50B')
    axes[0, 1].set_title('Message size = 1KiB')
    axes[1, 0].set_title('Message size = 10KiB')
    axes[1, 1].set_title('Message size = 100KiB')
    #percentiles 50B message
    df1 = df[['Concurrent Users', 'Message Size (Bytes)', '90th Percentile of Response Time (ms)','95th Percentile of Response Time (ms)','99th Percentile of Response Time (ms)']]
    df1 = df1.iloc[:6]
    print(df1)

    df1.plot(ylabel='Response Time (ms)', x="Concurrent Users", y=["90th Percentile of Response Time (ms)", "95th Percentile of Response Time (ms)", "99th Percentile of Response Time (ms)"], kind="bar",figsize=(19,18), ax=axes[0,0])
    plt.title('Message size = 50B')
    # plt.savefig('pfor50B.png')
    # plt.show()


    # #percentiles 1KiB message
    df2 = df[['Concurrent Users', 'Message Size (Bytes)', '90th Percentile of Response Time (ms)','95th Percentile of Response Time (ms)','99th Percentile of Response Time (ms)']]
    df2 = df2.iloc[6:12]
    print(df2)

    df2.plot(ylabel='Response Time (ms)', x="Concurrent Users", y=["90th Percentile of Response Time (ms)", "95th Percentile of Response Time (ms)", "99th Percentile of Response Time (ms)"], kind="bar",figsize=(19,18) , ax=axes[0,1])
    plt.title('Message size = 1KiB')
    # plt.savefig('pfor1KiB.png')
    # plt.show()

    # #percentiles 10KiB message
    df3 = df[['Concurrent Users', 'Message Size (Bytes)', '90th Percentile of Response Time (ms)','95th Percentile of Response Time (ms)','99th Percentile of Response Time (ms)']]
    df3 = df3.iloc[12:18]
    print(df3)

    df3.plot(ylabel='Response Time (ms)', x="Concurrent Users", y=["90th Percentile of Response Time (ms)", "95th Percentile of Response Time (ms)", "99th Percentile of Response Time (ms)"], kind="bar",figsize=(19,18), ax=axes[1,0])
    plt.title('Message size = 10KiB')
    # plt.savefig('pfor10KiB.png')
    # plt.show()


    #percentiles 100KiB message
    df4 = df[['Concurrent Users', 'Message Size (Bytes)', '90th Percentile of Response Time (ms)','95th Percentile of Response Time (ms)','99th Percentile of Response Time (ms)']]
    df4 = df4.iloc[18:24]
    print(df4)

    df4.plot(ylabel='Response Time (ms)', x="Concurrent Users", y=["90th Percentile of Response Time (ms)", "95th Percentile of Response Time (ms)", "99th Percentile of Response Time (ms)"], kind="bar",figsize=(19,18), ax=axes[1,1])
    plt.title('Message size = 100KiB')
    plt.savefig(os.path.join(path, 'percentiles_0ms.png'))
    # plt.show()


tps()
response_time()
percentiles()
